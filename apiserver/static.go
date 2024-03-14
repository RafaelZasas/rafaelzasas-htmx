package apiserver

import "net/http"

func StaticAssetsCacheMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// cache static assets for 1 week
		w.Header().Set("Cache-Control", "public, max-age=604800, immutable")
		next.ServeHTTP(w, r)
	})
}

func (api *api) bindStaticRoutes() {
	// Serve static files
	fileServer := http.FileServer(http.Dir("./public/"))
	cachedFileServer := StaticAssetsCacheMiddleware(fileServer)
	api.router.Handle("GET /public/", http.StripPrefix("/public/", cachedFileServer))
	api.router.HandleFunc("GET /robots.txt",
		func(w http.ResponseWriter, r *http.Request) {
			// Cache for max duration to avoid unnecessary requests by search engines
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
			http.ServeFile(w, r, "public/robots.txt")
		})
}
