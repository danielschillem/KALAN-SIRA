package httpapi

import (
	"net/http"

	"github.com/danielschillem/KALAN-SIRA/internal/admission"
)

func (r *Router) admissionCatalog(w http.ResponseWriter, q *http.Request) {
	sid, ok := schoolPrincipal(q)
	if !ok { writeError(w, 403, "forbidden", "school context required"); return }
	o, err := r.admissions.Catalog(q.Context(), sid)
	if err != nil { writeError(w, 500, "admission_catalog_failed", err.Error()); return }
	writeJSON(w, 200, o)
}

func (r *Router) admissionPreview(w http.ResponseWriter, q *http.Request) {
	sid, ok := schoolPrincipal(q)
	if !ok { writeError(w, 403, "forbidden", "school context required"); return }
	o, err := r.admissions.Preview(q.Context(), sid, q.PathValue("classID"))
	if err != nil { writeError(w, 404, "fee_configuration_not_found", "no fee configuration found for selected class"); return }
	writeJSON(w, 200, o)
}

func (r *Router) createAdmission(w http.ResponseWriter, q *http.Request) {
	sid, ok := schoolPrincipal(q)
	if !ok { writeError(w, 403, "forbidden", "school context required"); return }
	var in admission.CreateInput
	if !decode(w, q, &in) { return }
	o, err := r.admissions.CreateAndActivate(q.Context(), sid, in)
	if err != nil { writeError(w, 400, "admission_failed", err.Error()); return }
	writeJSON(w, 201, o)
}
