package middleware

import (
	"context"
	"go-inventory-management-api/utils"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// type App struct{
// 	// DB *sql.DB
// 	JWTKey []byte
// }

func AuthMiddleware(app *utils.App) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader  := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(w, "No token provided", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &utils.Claims{

		}
		token, err :=jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token)(interface{}, error){
			return app.JWTKey, nil
		})
		if err !=nil{
			if err == jwt.ErrSignatureInvalid {
        http.Error(w, "Invalid signature", http.StatusUnauthorized)
        return
      }
      http.Error(w, "Invalid token", http.StatusUnauthorized)
      return
		}
		if !token.Valid{
			http.Error(w, "Invalid token", http.StatusUnauthorized)
        return
		}

		ctx:=context.WithValue(r.Context(), "claims", claims)

    log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
    next.ServeHTTP(w, r.WithContext(ctx))
    })
  }
}
