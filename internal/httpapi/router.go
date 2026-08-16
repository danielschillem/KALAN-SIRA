package httpapi

import (
 "encoding/json"
 "errors"
 "net/http"
 "strings"

 "github.com/danielschillem/KALAN-SIRA/internal/academic"
 "github.com/danielschillem/KALAN-SIRA/internal/enrollment"
 "github.com/danielschillem/KALAN-SIRA/internal/school"
 "github.com/danielschillem/KALAN-SIRA/internal/student"
)

type Router struct{schools *school.Service; academics *academic.Service; students *student.Service; enrollments *enrollment.Service}
func NewRouter(schools *school.Service,academics *academic.Service,students *student.Service,enrollments *enrollment.Service)http.Handler{r:=&Router{schools:schools,academics:academics,students:students,enrollments:enrollments};mux:=http.NewServeMux();mux.HandleFunc("GET /health",r.health);mux.HandleFunc("POST /api/v1/schools",r.createSchool);mux.HandleFunc("GET /api/v1/schools/{publicID}",r.getSchool);mux.HandleFunc("POST /api/v1/school-years",r.createSchoolYear);mux.HandleFunc("POST /api/v1/levels",r.createLevel);mux.HandleFunc("POST /api/v1/classes",r.createClass);mux.HandleFunc("POST /api/v1/students",r.createStudent);mux.HandleFunc("POST /api/v1/guardians",r.createGuardian);mux.HandleFunc("POST /api/v1/student-guardians",r.linkGuardian);mux.HandleFunc("POST /api/v1/enrollments",r.createEnrollment);return mux}
func(r *Router)health(w http.ResponseWriter,_ *http.Request){writeJSON(w,200,map[string]string{"status":"ok","service":"kalan-sira-api","version":"0.3.0"})}
func decode(w http.ResponseWriter,req *http.Request,v any)bool{d:=json.NewDecoder(http.MaxBytesReader(w,req.Body,1<<20));d.DisallowUnknownFields();if err:=d.Decode(v);err!=nil{writeError(w,400,"invalid_request","invalid JSON body");return false};return true}
func(r *Router)createSchool(w http.ResponseWriter,req *http.Request){var in school.CreateInput;if !decode(w,req,&in){return};out,err:=r.schools.Create(req.Context(),in);respondCreated(w,out,err,"school_create_failed")}
func(r *Router)getSchool(w http.ResponseWriter,req *http.Request){out,err:=r.schools.GetByPublicID(req.Context(),strings.TrimSpace(req.PathValue("publicID")));if errors.Is(err,school.ErrNotFound){writeError(w,404,"school_not_found","school not found");return};if err!=nil{writeError(w,500,"internal_error","unable to load school");return};writeJSON(w,200,out)}
func(r *Router)createSchoolYear(w http.ResponseWriter,req *http.Request){var in academic.CreateSchoolYearInput;if !decode(w,req,&in){return};out,err:=r.academics.CreateSchoolYear(req.Context(),in);respondCreated(w,out,err,"school_year_create_failed")}
func(r *Router)createLevel(w http.ResponseWriter,req *http.Request){var in academic.CreateLevelInput;if !decode(w,req,&in){return};out,err:=r.academics.CreateLevel(req.Context(),in);respondCreated(w,out,err,"level_create_failed")}
func(r *Router)createClass(w http.ResponseWriter,req *http.Request){var in academic.CreateClassInput;if !decode(w,req,&in){return};out,err:=r.academics.CreateClass(req.Context(),in);respondCreated(w,out,err,"class_create_failed")}
func(r *Router)createStudent(w http.ResponseWriter,req *http.Request){var in student.CreateStudentInput;if !decode(w,req,&in){return};out,err:=r.students.CreateStudent(req.Context(),in);respondCreated(w,out,err,"student_create_failed")}
func(r *Router)createGuardian(w http.ResponseWriter,req *http.Request){var in student.CreateGuardianInput;if !decode(w,req,&in){return};out,err:=r.students.CreateGuardian(req.Context(),in);respondCreated(w,out,err,"guardian_create_failed")}
func(r *Router)linkGuardian(w http.ResponseWriter,req *http.Request){var in student.LinkGuardianInput;if !decode(w,req,&in){return};if err:=r.students.LinkGuardian(req.Context(),in);err!=nil{writeError(w,400,"guardian_link_failed",err.Error());return};writeJSON(w,201,map[string]string{"status":"linked"})}
func(r *Router)createEnrollment(w http.ResponseWriter,req *http.Request){var in enrollment.CreateInput;if !decode(w,req,&in){return};out,err:=r.enrollments.Create(req.Context(),in);respondCreated(w,out,err,"enrollment_create_failed")}
func respondCreated(w http.ResponseWriter,out any,err error,code string){if err!=nil{writeError(w,400,code,err.Error());return};writeJSON(w,201,out)}
func writeJSON(w http.ResponseWriter,status int,v any){w.Header().Set("Content-Type","application/json; charset=utf-8");w.WriteHeader(status);_ = json.NewEncoder(w).Encode(v)}
func writeError(w http.ResponseWriter,status int,code,message string){writeJSON(w,status,map[string]any{"error":map[string]string{"code":code,"message":message}})}
