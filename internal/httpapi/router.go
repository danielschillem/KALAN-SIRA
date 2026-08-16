package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/danielschillem/KALAN-SIRA/internal/school"
)

type Router struct{ schools *school.Service }

func NewRouter(schools *school.Service) http.Handler {
	r := &Router{schools: schools}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", r.health)
	mux.HandleFunc("POST /api/v1/schools", r.createSchool)
	mux.HandleFunc("GET /api/v1/schools/{publicID}", r.getSchool)
	return mux
}

func (r *Router) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status":"ok","service":"kalan-sira-api","version":"0.2.0"})
}

func (r *Router) createSchool(w http.ResponseWriter, req *http.Request) {
	var in school.CreateInput
	dec := json.NewDecoder(http.MaxBytesReader(w, req.Body, 1<<20))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "invalid JSON body")
		return
	}
	out, err := r.schools.Create(req.Context(), in)
	if err != nil {
		writeError(w, http.StatusBadRequest, "school_create_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, out)
}

func (r *Router) getSchool(w http.ResponseWriter, req *http.Request) {
	publicID := strings.TrimSpace(req.PathValue("publicID"))
	out, err := r.schools.GetByPublicID(req.Context(), publicID)
	if errors.Is(err, school.ErrNotFound) {
		writeError(w, http.StatusNotFound, "school_not_found", "school not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "unable to load school")
		return
	}
	writeJSON(w, http.StatusOK, out)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, map[string]any{"error": map[string]string{"code": code, "message": message}})
}
