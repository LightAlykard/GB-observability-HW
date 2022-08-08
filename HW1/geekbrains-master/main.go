package main

import (
	"geekbrains/middleware"
	"html/template"
	"net/http"
	"path"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	r := mux.NewRouter()

	metricsMiddleware := middleware.NewMetricsMiddleware()

	r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/alert", alertHandler).Methods(http.MethodGet)
	r.HandleFunc("/simple", simpleHandler).Methods(http.MethodPost)
	r.HandleFunc("/hard", hardHandler).Methods(http.MethodPut)

	r.Use(metricsMiddleware.Metrics)

	http.ListenAndServe(":8080", r)
}

func alertHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("Alert"))
}

func simpleHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("Simple"))
}

func hardHandler(w http.ResponseWriter, r *http.Request) {
	fp := path.Join("templates", "hard.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
