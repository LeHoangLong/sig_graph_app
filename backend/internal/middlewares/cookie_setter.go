package middlewares

import (
	"context"
	"net/http"
)

type CookieSetterKeyType struct {
	name string
}

type CookieSetterFunc func(name string, value string, maxAge_s int)

var CookieSetterKey = &CookieSetterKeyType{"cookie-setter"}

func SetCookieMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		SetCookie := func(name string, value string, maxAge_s int) {
			cookie := http.Cookie{
				Name:   name,
				Value:  value,
				MaxAge: maxAge_s,
			}
			http.SetCookie(w, &cookie)
		}

		ctx := context.WithValue(r.Context(), CookieSetterKey, SetCookie)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetCookieSetter(ctx context.Context) CookieSetterFunc {
	if setter, ok := ctx.Value(CookieSetterKey).(CookieSetterFunc); ok {
		return setter
	} else {
		return nil
	}
}
