package srv

import (
	"github.com/meschbach/elevatinator/pkg/ipc/grpc/telepathy/pb"
	"github.com/meschbach/elevatinator/pkg/simulator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
)

type config struct {
	//listenAt is the address to bind the gRPC service too
	listenAt string
	//publishHealthService , when true, will export the gRPC health protocol
	// See https://github.com/grpc/grpc/blob/master/doc/health-checking.md for more details
	publishHealthService bool
}

type Option func(c *config)

// ListenAt overrides the default address, 'localhost:9998', to the supplied target address.
func ListenAt(address string) Option {
	return func(c *config) {
		c.listenAt = address
	}
}

func DisableHealthService() Option {
	return func(c *config) {
		c.publishHealthService = false
	}
}

// Network is an abstraction for binding to a listener.  Closes the gap between real world usage of the code and testing
// paths.
type Network interface {
	Listener() (net.Listener, error)
}

// RunControllerService exports the given controller on port tcp/9998 over gRPC
func RunControllerService(builder simulator.ControllerFunc, withOptions ...Option) error {
	c := &config{
		listenAt:             "localhost:9998",
		publishHealthService: true,
	}
	for _, o := range withOptions {
		o(c)
	}

	l := &tcp{listenAt: c.listenAt}
	return RunControllerOn(builder, l, func(server *grpc.Server) error {
		if !c.publishHealthService {
			return nil
		}
		healthService := health.NewServer()
		healthService.SetServingStatus("ai", grpc_health_v1.HealthCheckResponse_SERVING)
		grpc_health_v1.RegisterHealthServer(server, healthService)
		return nil
	})
}

// RunControllerOn exports the given controller on the specified network address
func RunControllerOn(builder simulator.ControllerFunc, on Network, otherServices ...func(server *grpc.Server) error) error {
	s := grpc.NewServer()
	pb.RegisterControllerServiceServer(s, newRemoteController(builder))
	for _, otherService := range otherServices {
		if err := otherService(s); err != nil {
			return err
		}
	}

	listener, err := on.Listener()
	if err != nil {
		return err
	}

	if err := s.Serve(listener); err != nil {
		return err
	}
	return nil
}

type tcp struct {
	listenAt string
}

func (r *tcp) Listener() (net.Listener, error) {
	lis, err := net.Listen("tcp", r.listenAt)
	if err != nil {
		return nil, err
	}
	return lis, nil
}
