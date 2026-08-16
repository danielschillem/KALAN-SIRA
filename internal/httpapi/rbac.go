package httpapi

import "net/http"

func hasRole(role string, allowed ...string) bool { for _,a:=range allowed { if role==a{return true} }; return false }
func (r *Router) requireRoles(next http.Handler, roles ...string) http.Handler { return http.HandlerFunc(func(w http.ResponseWriter,q *http.Request){p,ok:=principal(q);if !ok||!hasRole(p.Role,roles...){writeError(w,403,"forbidden","insufficient permissions");return};next.ServeHTTP(w,q)}) }
func (r *Router) requireSchool(next http.Handler) http.Handler { return http.HandlerFunc(func(w http.ResponseWriter,q *http.Request){p,ok:=principal(q);if !ok{writeError(w,401,"unauthorized","authentication required");return};if p.Role=="SUPER_ADMIN"{next.ServeHTTP(w,q);return};if p.SchoolPublicID==""{writeError(w,403,"forbidden","school context required");return};next.ServeHTTP(w,q)}) }
func schoolPrincipal(q *http.Request)(string,bool){p,ok:=principal(q);if !ok||p.SchoolPublicID==""{return "",false};return p.SchoolPublicID,true}
