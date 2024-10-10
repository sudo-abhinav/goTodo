package middlewares

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"github.com/sudo-abhinav/go-todo/Database/dbHelper"
	"github.com/sudo-abhinav/go-todo/utils/response"
	"net/http"
	"os"
	"strings"
)

// import (
//
//	"context"
//	"fmt"
//	"github.com/sudo-abhinav/go-todo/utils/encryption"
//	"net/http"
//
// )
//
//	func Authentication(next http.Handler) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			cookie, err := r.Cookie("Token")
//			if err != nil {
//				http.Error(w, "Unauthorized: missing cookie", http.StatusUnauthorized)
//				fmt.Println(err)
//				return
//			}
//			fmt.Println(cookie)
//
//			tokenString := cookie.Value
//			claims, err := encryption.VerifyJWT(tokenString)
//			if err != nil {
//				http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
//				return
//			}
//
//			ctx := r.Context()
//			ctx = context.WithValue(ctx, "claims", claims)
//			r = r.WithContext(ctx)
//
//			next.ServeHTTP(w, r)
//		})
//	}
type ContextKeys string

const (
	userContext ContextKeys = "userContext"
)

type UserCtx struct {
	UserID    string
	Email     string
	SessionID string
}

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.RespondWithError(w, http.StatusUnauthorized, nil, "authorization header missing")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			response.RespondWithError(w, http.StatusUnauthorized, nil, "bearer token missing")
			return
		}

		token, parseErr := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method") // Invalid signing method error
			}
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		})

		if parseErr != nil || !token.Valid {
			response.RespondWithError(w, http.StatusUnauthorized, parseErr, "invalid claims")
			return
		}

		claimValues, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			response.RespondWithError(w, http.StatusUnauthorized, nil, "invalid token ")
			return
		}

		sessionID := claimValues["sessionID"].(string)

		archivedAt, err := dbHelper.GetArchivedAt(sessionID)
		if err != nil {
			response.RespondWithError(w, http.StatusInternalServerError, err, "internal server error")
			return
		}

		if archivedAt != nil {
			response.RespondWithError(w, http.StatusUnauthorized, nil, "invalid token")
			return
		}

		user := &UserCtx{
			UserID:    claimValues["userID"].(string),
			Email:     claimValues["email"].(string),
			SessionID: sessionID,
		}

		ctx := context.WithValue(r.Context(), userContext, user)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func UserContext(r *http.Request) *UserCtx {
	if user, ok := r.Context().Value(userContext).(*UserCtx); ok {
		return user
	}
	return nil
}
