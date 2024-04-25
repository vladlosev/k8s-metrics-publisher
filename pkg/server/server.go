package server

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/vladlosev/k8s-metrics-publisher/pkg/client"
)

// Server is the type implementing the
type Server struct {
	http.Server
	mux    *http.ServeMux
	client *client.Client
}

// New returns a new instance of Server.
func New(client *client.Client, port uint32, endpointPath string) *Server {
	mux := http.NewServeMux()
	server := &Server{
		Server: http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: mux,
		},
		mux:    mux,
		client: client,
	}
	server.mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("Me be healthy. Me be smart."))
		if err != nil {
			logrus.WithError(err).Error("Error writing healthz response")
		}
	})
	server.mux.HandleFunc(endpointPath, server.handleMetrics)
	return server
}

// handleMetrics writes the contents of the current API server's /metrics page
// to the response.
func (s *Server) handleMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.Header().Add("Allow", "GET")
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	body, err := s.client.GetMetrics(r.Context())
	if err != nil {
		logrus.
			WithError(err).
			WithField("body", string(body)).
			Error("error polling API server")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	_, err = w.Write(body)
	if err != nil {
		logrus.WithError(err).Error("Error writing metrics response")
	}
}
