package middleware

import "net/http"

const (
	accessControlAllowOrigin      = "Access-Control-Allow-Origin"
	accessControlAllowMethods     = "Access-Control-Allow-Methods"
	accessControlAllowHeaders     = "Access-Control-Allow-Headers"
	accessControlAllowCredentials = "Access-Control-Allow-Credentials"
	varyHeader                    = "Vary"
	originHeader                  = "Origin"
	requestMethodHeader           = "Access-Control-Request-Method"
	requestHeadersHeader          = "Access-Control-Request-Headers"
	defaultAllowedMethods         = "GET,POST,PUT,PATCH,DELETE,OPTIONS"
	defaultAllowedHeaders         = "Accept,Authorization,Content-Type,Origin,X-Requested-With"
)

// CORS enables browser access from the local desktop dev server and Electron renderer.
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get(originHeader)
		if origin != "" {
			w.Header().Set(accessControlAllowOrigin, origin)
			w.Header().Set(accessControlAllowCredentials, "true")
			w.Header().Set(accessControlAllowMethods, defaultAllowedMethods)
			w.Header().Set(accessControlAllowHeaders, defaultAllowedHeaders)
			w.Header().Add(varyHeader, originHeader)
			w.Header().Add(varyHeader, requestMethodHeader)
			w.Header().Add(varyHeader, requestHeadersHeader)
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
