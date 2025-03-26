package middleware

import (
	"context"
	"go-inventory-management-api/utils"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct{
	Username string `json:"username"`
	UserID string `json:"id"`
	jwt.RegisteredClaims
}

// type App struct{
// 	// DB *sql.DB
// 	JWTKey []byte
// }

func JwtMiddleware(app *utils.App) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
      log.Println(w, "Middleware: JWT validation")
		authHeader  := r.Header.Get("Authorization")
		log.Println(w, authHeader)

		if authHeader == "" {
			http.Error(w, "No token provided", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		log.Println(w, "asdas\n\n\n"+tokenString)
		claims := &Claims{

		}
		token, err :=jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token)(interface{}, error){
			return app.JWTKey, nil
		})
		log.Println(w, "88888888888")
		if err !=nil{
			if err == jwt.ErrSignatureInvalid {
        http.Error(w, "Invalid signature", http.StatusUnauthorized)
        return
      }
      http.Error(w, "Invalid token", http.StatusUnauthorized)
      return
		}
		log.Println(w, "999999999999")
		if !token.Valid{
			http.Error(w, "Invalid token", http.StatusUnauthorized)
        return
		}
		log.Println(w, "88888888888")

		ctx:=context.WithValue(r.Context(), "claims", claims)

    log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
    next.ServeHTTP(w, r.WithContext(ctx))
    })
  }
}
