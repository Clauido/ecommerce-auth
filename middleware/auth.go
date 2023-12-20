package middleware

import (
	"strings"
)

var (
	NO_AUTH_NEEDED = []string{
		"/login",
		"/signup",
	}
)

func shouldChekToken(route string) bool {
	for _, p := range NO_AUTH_NEEDED {
		if strings.Contains(route, p) {
			return false
		}
	}
	return true
}
// func CheckAuthMiddleWare(s server.Server) func(h http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler{
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
// 			if !shouldChekToken(r.URL.Path){
// 				log.Println(r.URL.Path)
// 				next.ServeHTTP(w, r)
// 				return
// 			}
// 			tokenString := strings.TrimSpace(r.Header.Get("Authorization"))
// 			_, err:= jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(token *jwt.Token) (interface{}, error){
// 				return []byte(s.Config().JWTSecret),nil
// 			})
// 			if err!=nil{
// 				http.Error(w,err.Error(),http.StatusUnauthorized)
// 				return 
// 			}
// 			next.ServeHTTP(w, r)
// 		})
// 	}
// }