package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/wibecoderr/storex"
	"github.com/wibecoderr/storex/database/dbhelper"
	"github.com/wibecoderr/storex/model"
)

type contextKey struct{}

var userContextKey = contextKey{}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.RespondError(w, http.StatusUnauthorized, nil, "missing authorization header")
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			utils.RespondError(w, http.StatusUnauthorized, nil, "invalid format, expected Bearer <token>")
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == "" {
			utils.RespondError(w, http.StatusUnauthorized, nil, "empty token")
			return
		}

		userID, sessionID, err := utils.VerifyJWT(tokenStr)
		if err != nil {
			utils.RespondError(w, http.StatusUnauthorized, nil, "invalid or expired token")
			return
		}

		// verify session exists in DB
		dbUserID, err := dbhelper.GetUserIDBySession(sessionID)
		if err != nil || dbUserID != userID {
			utils.RespondError(w, http.StatusUnauthorized, nil, "session not found or expired")
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, &model.UserCxt{
			UserId:    userID,
			SessionId: sessionID,
		})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RoleMiddleware(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := UserContext(r)
			if user == nil {
				utils.RespondError(w, http.StatusUnauthorized, nil, "unauthorized")
				return
			}

			role, err := dbhelper.GetEmployeeRole(user.UserId)
			if err != nil {
				utils.RespondError(w, http.StatusInternalServerError, err, "failed to get role")
				return
			}

			for _, allowed := range allowedRoles {
				if allowed == role {
					next.ServeHTTP(w, r)
					return
				}
			}

			utils.RespondError(w, http.StatusUnauthorized, nil, "role not allowed")
		})
	}
}

func UserContext(r *http.Request) *model.UserCxt {
	user, _ := r.Context().Value(userContextKey).(*model.UserCxt)
	return user
}