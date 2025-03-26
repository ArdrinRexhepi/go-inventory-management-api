package middleware

import (
	"go-inventory-management-api/utils"
	"log"
	"net/http"
)
func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value("claims").(*utils.Claims)
		if !ok || !claims.IsAdmin { 
		log.Println(claims.IsAdmin)
			http.Error(w, "Forbidden: Admin access required", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
