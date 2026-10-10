package main

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/SerraHarmankaya/Chirpy/system"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {

	mux := http.ServeMux{}

	apiCfg := apiConfig{}

	server := &http.Server{
		Addr:    ":8080",
		Handler: &mux,
	}

	fileServer := http.FileServer(http.Dir("./"))

	fileHandler := http.StripPrefix("/app", fileServer)

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(fileHandler))

	mux.HandleFunc("GET /healthz", system.ReadinessHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.fileserverHitsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetFileserverHitsHandler)

	err := server.ListenAndServe()
	if err != nil {
		return
	}
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) fileserverHitsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("<html>\n  <body>\n    <h1>Welcome, Chirpy Admin</h1>\n    <p>Chirpy has been visited %d times!</p>\n  </body>\n</html>", cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) resetFileserverHitsHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}
