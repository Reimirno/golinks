package crud

import (
	"context"
	"fmt"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"reimirno.com/golinks/pkg/logging"
	"reimirno.com/golinks/pkg/mapper"
	"reimirno.com/golinks/pkg/pb"
	"reimirno.com/golinks/pkg/types"
)

const crudServiceName = "crud"

type Server struct {
	manager *mapper.MapperManager
	logger  *zap.SugaredLogger
	port    string
	server  *grpc.Server
	pb.UnimplementedGolinksServer
}

var _ types.Service = (*Server)(nil)

func (s *Server) GetName() string {
	return crudServiceName
}

func (s *Server) Start(errChan chan<- error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", s.port))
	if err != nil {
		s.logger.Errorf("Failed to start service %s", s.GetName())
		errChan <- err
		return
	}

	go func() {
		s.logger.Infof("Service %s starting on %s...", s.GetName(), lis.Addr())
		if err := s.server.Serve(lis); err != nil {
			s.logger.Errorf("Failed to start service %s", s.GetName())
			errChan <- err
		}
	}()
}

func (s *Server) Stop() error {
	s.logger.Infof("Shutting down service %s...", s.GetName())
	s.server.GracefulStop()
	s.logger.Infof("Service %s shutdown complete", s.GetName())
	return nil
}

func NewServer(m *mapper.MapperManager, port string, debug bool) (*Server, error) {
	logger := logging.NewLogger(crudServiceName)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(logging.GrpcInterceptor(logger)),
	)
	service := &Server{
		manager: m,
		logger:  logger,
		port:    port,
		server:  server,
	}
	pb.RegisterGolinksServer(server, service)
	if debug {
		reflection.Register(server)
	}

	return service, nil
}

func (s *Server) GetUrl(ctx context.Context, req *pb.GetUrlRequest) (*pb.PathUrlPair, error) {
	pair, err := s.manager.GetUrl(req.Path, false)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get url: %v", err)
	}
	if pair == nil {
		return nil, status.Errorf(codes.NotFound, "path %s not found", req.Path)
	}
	return getProto(pair), nil
}

func (s *Server) PutUrl(ctx context.Context, req *pb.PathUrlPair) (*pb.PathUrlPair, error) {
	pair := getStruct(req)
	pair, err := s.manager.PutUrl(pair)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to put url: %v", err)
	}
	return getProto(pair), nil
}

func (s *Server) DeleteUrl(ctx context.Context, req *pb.DeleteUrlRequest) (*emptypb.Empty, error) {
	err := s.manager.DeleteUrl(req.Path)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete url: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) ListUrls(ctx context.Context, req *emptypb.Empty) (*pb.ListUrlsResponse, error) {
	pairs, err := s.manager.ListUrls()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list urls: %v", err)
	}
	result := make([]*pb.PathUrlPair, 0, len(pairs))
	for _, pair := range pairs {
		result = append(result, getProto(pair))
	}
	return &pb.ListUrlsResponse{
		Pairs: result,
	}, nil
}
