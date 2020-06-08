package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/KishoreCM/TodoApp/models"
	"github.com/dgrijalva/jwt-go"
)

type Exception struct {
	Message string `json:"message"`
}

type RefreshToken struct {
	Token string `json: "token"`
}

func generateRefreshToken(claims *models.Claims) (string, error) {
	// Now, create a new token for the current use, with a renewed expiration time
	expirationTime := time.Now().Add(30 * time.Second)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		fmt.Println("JWT Signing Failed")
		return "", err
	}
	return tokenString, nil
}

// JwtVerify Middleware function
func JwtVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		/*var header = r.Header.Get("x-access-token") //Grab the token from the header

		header = strings.TrimSpace(header)

		if header == "" {
			//Token is missing, returns with error code 403 Unauthorized
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(Exception{Message: "Missing auth token"})
			return
		}
		*/
		cookie, cErr := r.Cookie("todoToken")
		if cErr != nil {
			if cErr == http.ErrNoCookie {
				// If the cookie is not set, return an unauthorized status
				//w.WriteHeader(http.StatusUnauthorized)
				//json.NewEncoder(w).Encode(Exception{Message: cErr.Error()})
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
		}
		token := cookie.Value

		claims := &models.Claims{}

		_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})

		if err != nil {
			if strings.Contains(err.Error(), "expired") {
				fmt.Println(err.Error())
				refreshToken, tokenErr := generateRefreshToken(claims)
				if tokenErr != nil {
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(Exception{Message: tokenErr.Error()})
					return
				}
				fmt.Println("Refresh token generated!")

				http.SetCookie(w, &http.Cookie{
					Name:  "todoToken",
					Value: refreshToken,
				})
				ctx := context.WithValue(r.Context(), "userId", claims.UserID)
				next.ServeHTTP(w, r.WithContext(ctx))
				//next.ServeHTTP(w, r)
				return
			}

			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(Exception{Message: err.Error()})
			return
		}

		ctx := context.WithValue(r.Context(), "userId", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
		//next.ServeHTTP(w, r)
	})
}
