package redirector

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/reimirno/golinks/pkg/logging"
	"github.com/reimirno/golinks/pkg/mapper"
	"github.com/reimirno/golinks/pkg/types"
)

const redirectorServiceName = "redirector"

type Server struct {
	server  *http.Server
	logger  *zap.SugaredLogger
	manager *mapper.MapperManager
	port    string
}

var _ types.Service = (*Server)(nil)

func (s *Server) GetName() string {
	return redirectorServiceName
}

func (s *Server) Start(errChan chan<- error) {
	go func() {
		s.logger.Infof("Service %s starting on %s...", s.GetName(), s.server.Addr)
		errChan <- s.server.ListenAndServe()
	}()
}

func (s *Server) Stop() error {
	s.logger.Infof("Shutting down service %s...", s.GetName())
	s.server.Shutdown(context.Background())
	s.logger.Infof("Service %s shutdown complete", s.GetName())
	return nil
}

func NewServer(m *mapper.MapperManager, port string) (*Server, error) {
	r := mux.NewRouter()
	l := logging.NewLogger(redirectorServiceName)
	r.Use(logging.HttpMiddleware(l))

	r.HandleFunc("/{path}", func(rw http.ResponseWriter, r *http.Request) {
		path := mux.Vars(r)["path"]
		pair, err := m.GetUrl(path, true)
		handleError := func(rw http.ResponseWriter, msg string, err error, statusCode int) {
			l.Errorf("%s: %v", msg, err)
			http.Error(rw, msg, statusCode)
		}
		if err != nil {
			handleError(rw, fmt.Sprintf("Error occured when resolving path: %v", err), err, http.StatusInternalServerError)
			return
		}
		if pair != nil {
			l.Infof("Mapping found: %s -> %s", path, pair.Url)
			http.Redirect(rw, r, pair.Url, http.StatusFound)
			return
		}
		handleError(rw, fmt.Sprintf("Mapping not found: %s", path), nil, http.StatusNotFound)
	}).Methods("GET")

	s := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: r,
	}

	return &Server{
		server:  s,
		logger:  l,
		manager: m,
		port:    port,
	}, nil
}
