package crud_http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/reimirno/golinks/pkg/logging"
	"github.com/reimirno/golinks/pkg/mapper"
	"github.com/reimirno/golinks/pkg/types"
	"github.com/reimirno/golinks/pkg/utils"
)

const crudHttpServiceName = "crud_http"

type Server struct {
	manager *mapper.MapperManager
	logger  *zap.SugaredLogger
	server  *http.Server
	port    string
}

var _ types.Service = (*Server)(nil)

func (s *Server) GetName() string {
	return crudHttpServiceName
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
	l := logging.NewLogger(crudHttpServiceName)
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
	r.HandleFunc("/go/{path}/", svr.handleGetUrl).Methods("GET")
	r.HandleFunc("/go/", svr.handleListUrls).Methods("GET")
	r.HandleFunc("/go/", svr.handlePutUrl).Methods("PUT")
	r.HandleFunc("/go/{path}/", svr.handleDeleteUrl).Methods("DELETE")
	return svr, nil
}

func (s *Server) handleGetUrl(rw http.ResponseWriter, r *http.Request) {
	fmt.Println("handleGetUrl")
	vars := mux.Vars(r)
	path := vars["path"]
	pair, err := s.manager.GetUrl(path, false)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	if pair == nil {
		http.Error(rw, fmt.Sprintf("path %s not found", path), http.StatusNotFound)
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(pair)
}

func (s *Server) handleListUrls(rw http.ResponseWriter, r *http.Request) {
	var err error
	pagination := utils.DefaultPagination
	offset := r.URL.Query().Get("offset")
	if offset != "" {
		pagination.Offset, err = strconv.Atoi(offset)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
	}
	limit := r.URL.Query().Get("limit")
	if limit != "" {
		pagination.Limit, err = strconv.Atoi(limit)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
	}
	pairs, err := s.manager.ListUrls(pagination)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(pairs)
}

func (s *Server) handlePutUrl(rw http.ResponseWriter, r *http.Request) {
	var pair types.PathUrlPair
	err := json.NewDecoder(r.Body).Decode(&pair)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	pairPut, err := s.manager.PutUrl(&pair)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusAccepted)
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(pairPut)
}

func (s *Server) handleDeleteUrl(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := vars["path"]
	err := s.manager.DeleteUrl(path)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusNoContent)
}
