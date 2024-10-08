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
	err := s.server.Shutdown(context.Background())
	if err != nil {
		s.logger.Errorf("Error shutting down service %s: %v", s.GetName(), err)
	}
	s.logger.Infof("Service %s shutdown complete", s.GetName())
	return err
}

func NewServer(m *mapper.MapperManager, port string) (*Server, error) {
	r := mux.NewRouter()
	l := logging.NewLogger(redirectorServiceName)
	r.Use(logging.HttpMiddleware(l))
	s := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: r,
	}
	svr := &Server{
		server:  s,
		logger:  l,
		manager: m,
		port:    port,
	}
	r.HandleFunc("/{path}", svr.handleRedirect).Methods("GET")
	return svr, nil
}

func (s *Server) handleRedirect(rw http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)["path"]
	pair, err := s.manager.GetUrl(path, true)

	handleError := func(rw http.ResponseWriter, msg string, err error, statusCode int) {
		s.logger.Errorf("%s: %v", msg, err)
		http.Error(rw, msg, statusCode)
	}
	if err != nil {
		handleError(rw, fmt.Sprintf("Error occurred when resolving path: %v", err), err, http.StatusInternalServerError)
		return
	}
	if pair != nil {
		s.logger.Infof("Mapping found: %s -> %s", path, pair.Url)
		http.Redirect(rw, r, pair.Url, http.StatusFound)
		return
	}
	handleError(rw, fmt.Sprintf("Mapping not found: %s", path), nil, http.StatusNotFound)
}
