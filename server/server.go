package server

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/ramin/waypoint/config"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// Server is a structured wrapper for
// running the API server, encapsulating a
// route and error handling
type Server struct {
	http http.Server

	// add this to generic server to different
	// implementations of server units
	// cmd/**/* can bind subsets of routes
	// specific to their output binary
	r *mux.Router

	Close chan os.Signal
	Wait  chan bool
	Error chan error
}

// NewServer returns a new Server with
// defaults set, routes bootstrapped
// and middleware set up
func NewServer() Server {
	r := mux.NewRouter()

	r.HandleFunc("/healthz", StatusHandler) // assimilate to borg
	r.Handle("/metrics", promhttp.Handler())

	return Server{
		http: http.Server{
			Addr:    config.Read().Listen,
			Handler: r,
		},
		r:     r,
		Close: make(chan os.Signal),
		Wait:  make(chan bool),
		Error: make(chan error),
	}
}

func (s *Server) Router() *mux.Router {
	return s.r
}

// Start starts the http server
func (s *Server) Start() chan bool {
	err := s.r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		tpl, err := route.GetPathTemplate()
		if err != nil {
			log.Errorf("route: %s", err)
		}

		log.Info(tpl)
		return nil
	})

	if err != nil {
		fmt.Println(err)
	}

	go func() {
		log.Infof("server starting and listening on %s", s.http.Addr)
		err := s.http.ListenAndServe()
		if err != nil {
			if err != http.ErrServerClosed {
				log.Error(err)
				os.Exit(1)
			}
		}
	}()

	return s.Wait
}

func (s *Server) close() {
	close(s.Error)
	close(s.Wait)
	close(s.Close)
}

// GracefulShutdown safely shuts down the server
func (s *Server) GracefulShutdown() {
	if err := s.http.Shutdown(context.Background()); err != nil {
		log.Info("failed to shutdown properly")
		log.Error(err)
	}

	s.close()
}
