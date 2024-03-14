package apiserver

import (
	"context"
	"htmx-go/database"
	"htmx-go/models/permissions"
	"htmx-go/utils"
	"log/slog"
	"net/http"

	"github.com/ProteaTech/bok"
)

func authMiddleware() bok.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Validate User Authentication Status and set data.User if user is authenticated
			uid, err := utils.ValidateJWT(r)

			if err != nil {
				// If the JWT is invalid, check the refresh token
				uid, err = utils.ValidateRefreshToken(r)
				// If the refresh token is valid.
				// Create a new JWT pair and set the cookies
				if err != nil {
					next(w, r)
					return
				}

				t, rt, err := utils.GenerateJWTokens(uid)
				if err != nil {
					slog.Error("could not generate JWT", "err", err)
				}
				utils.AddAuthCookies(&w, t, rt)
			}

			ctx := context.WithValue(r.Context(), "uid", uid)
			r = r.WithContext(ctx)

			next(w, r)
		}
	}
}

func adminRouteGuard() bok.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			uid := r.Context().Value("uid")

			if uid == nil {
				http.Redirect(w, r, "/404", http.StatusSeeOther)
				return
			}

			db := r.Context().Value("db").(database.Database)
			hasPermission, err := db.HasPermission(uid.(string), permissions.ViewAdmin)
			if err != nil {
				slog.Error("adminRouteGuard: db.HasPermission", "err", err)
				http.Redirect(w, r, "/500", http.StatusInternalServerError)
				return
			}

			if !hasPermission {
				slog.Warn("adminRouteGuard: user does not have permission to view admin panel")
				http.Redirect(w, r, "/404", http.StatusSeeOther)
				return
			}
			next(w, r)
		}
	}
}
