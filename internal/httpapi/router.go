package httpapi

import (
 "encoding/json"
 "errors"
 "net/http"
 "strings"
 "github.com/danielschillem/KALAN-SIRA/internal/academic"
 "github.com/danielschillem/KALAN-SIRA/internal/billing"
 "github.com/danielschillem/KALAN-SIRA/internal/enrollment"
 "github.com/danielschillem/KALAN-SIRA/internal/school"
 "github.com/danielschillem/KALAN-SIRA/internal/student"
)
type Router struct{schools *school.Service;academics *academic.Service;students *student.Service;enrollments *enrollment.Service;billing *billing.Service}
func NewRouter(s *school.Service,a *academic.Service,st *student.Service,e *enrollment.Service,b *billing.Service)http.Handler{r:=&Router{s,a,st,e,b};m:=http.NewServeMux();m.HandleFunc("GET /health",r.health);m.HandleFunc("POST /api/v1/schools",r.createSchool);m.HandleFunc("GET /api/v1/schools/{publicID}",r.getSchool);m.HandleFunc("POST /api/v1/school-years",r.createSchoolYear);m.HandleFunc("POST /api/v1/levels",r.createLevel);m.HandleFunc("POST /api/v1/classes",r.createClass);m.HandleFunc("POST /api/v1/students",r.createStudent);m.HandleFunc("POST /api/v1/guardians",r.createGuardian);m.HandleFunc("POST /api/v1/student-guardians",r.linkGuardian);m.HandleFunc("POST /api/v1/enrollments",r.createEnrollment);m.HandleFunc("POST /api/v1/fee-schedules",r.createFeeSchedule);m.HandleFunc("POST /api/v1/fee-items",r.createFeeItem);m.HandleFunc("POST /api/v1/installment-plans",r.createPlan);m.HandleFunc("POST /api/v1/installments",r.createInstallment);m.HandleFunc("POST /api/v1/enrollments/activate",r.activateEnrollment);return m}
func(r *Router)health(w http.ResponseWriter,_ *http.Request){writeJSON(w,200,map[string]string{"status":"ok","service":"kalan-sira-api","version":"0.4.0"})}
func decode(w http.ResponseWriter,q *http.Request,v any)bool{d:=json.NewDecoder(http.MaxBytesReader(w,q.Body,1<<20));d.DisallowUnknownFields();if d.Decode(v)!=nil{writeError(w,400,"invalid_request","invalid JSON body");return false};return true}
func(r *Router)createSchool(w http.ResponseWriter,q *http.Request){var in school.CreateInput;if !decode(w,q,&in){return};o,e:=r.schools.Create(q.Context(),in);created(w,o,e,"school_create_failed")}
func(r *Router)getSchool(w http.ResponseWriter,q *http.Request){o,e:=r.schools.GetByPublicID(q.Context(),strings.TrimSpace(q.PathValue("publicID")));if errors.Is(e,school.ErrNotFound){writeError(w,404,"school_not_found","school not found");return};if e!=nil{writeError(w,500,"internal_error","unable to load school");return};writeJSON(w,200,o)}
func(r *Router)createSchoolYear(w http.ResponseWriter,q *http.Request){var in academic.CreateSchoolYearInput;if !decode(w,q,&in){return};o,e:=r.academics.CreateSchoolYear(q.Context(),in);created(w,o,e,"school_year_create_failed")}
func(r *Router)createLevel(w http.ResponseWriter,q *http.Request){var in academic.CreateLevelInput;if !decode(w,q,&in){return};o,e:=r.academics.CreateLevel(q.Context(),in);created(w,o,e,"level_create_failed")}
func(r *Router)createClass(w http.ResponseWriter,q *http.Request){var in academic.CreateClassInput;if !decode(w,q,&in){return};o,e:=r.academics.CreateClass(q.Context(),in);created(w,o,e,"class_create_failed")}
func(r *Router)createStudent(w http.ResponseWriter,q *http.Request){var in student.CreateStudentInput;if !decode(w,q,&in){return};o,e:=r.students.CreateStudent(q.Context(),in);created(w,o,e,"student_create_failed")}
func(r *Router)createGuardian(w http.ResponseWriter,q *http.Request){var in student.CreateGuardianInput;if !decode(w,q,&in){return};o,e:=r.students.CreateGuardian(q.Context(),in);created(w,o,e,"guardian_create_failed")}
func(r *Router)linkGuardian(w http.ResponseWriter,q *http.Request){var in student.LinkGuardianInput;if !decode(w,q,&in){return};if e:=r.students.LinkGuardian(q.Context(),in);e!=nil{writeError(w,400,"guardian_link_failed",e.Error());return};writeJSON(w,201,map[string]string{"status":"linked"})}
func(r *Router)createEnrollment(w http.ResponseWriter,q *http.Request){var in enrollment.CreateInput;if !decode(w,q,&in){return};o,e:=r.enrollments.Create(q.Context(),in);created(w,o,e,"enrollment_create_failed")}
func(r *Router)createFeeSchedule(w http.ResponseWriter,q *http.Request){var in billing.CreateFeeScheduleInput;if !decode(w,q,&in){return};o,e:=r.billing.CreateFeeSchedule(q.Context(),in);created(w,o,e,"fee_schedule_create_failed")}
func(r *Router)createFeeItem(w http.ResponseWriter,q *http.Request){var in billing.CreateFeeItemInput;if !decode(w,q,&in){return};o,e:=r.billing.CreateFeeItem(q.Context(),in);created(w,o,e,"fee_item_create_failed")}
func(r *Router)createPlan(w http.ResponseWriter,q *http.Request){var in billing.CreatePlanInput;if !decode(w,q,&in){return};o,e:=r.billing.CreatePlan(q.Context(),in);created(w,o,e,"installment_plan_create_failed")}
func(r *Router)createInstallment(w http.ResponseWriter,q *http.Request){var in billing.CreateInstallmentInput;if !decode(w,q,&in){return};o,e:=r.billing.CreateInstallment(q.Context(),in);created(w,o,e,"installment_create_failed")}
func(r *Router)activateEnrollment(w http.ResponseWriter,q *http.Request){var in billing.ActivateEnrollmentInput;if !decode(w,q,&in){return};o,e:=r.billing.ActivateEnrollment(q.Context(),in);if e!=nil{writeError(w,400,"enrollment_activation_failed",e.Error());return};writeJSON(w,200,o)}
func created(w http.ResponseWriter,o any,e error,c string){if e!=nil{writeError(w,400,c,e.Error());return};writeJSON(w,201,o)}
func writeJSON(w http.ResponseWriter,s int,v any){w.Header().Set("Content-Type","application/json; charset=utf-8");w.WriteHeader(s);_ = json.NewEncoder(w).Encode(v)}
func writeError(w http.ResponseWriter,s int,c,m string){writeJSON(w,s,map[string]any{"error":map[string]string{"code":c,"message":m}})}
