package middleware

import "net/http"

func OauthMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("x-api-key")
		if header != "divya-zs" {
			w.WriteHeader(http.StatusUnauthorized)

			return
		}
		h.ServeHTTP(w, r)
	})
}
