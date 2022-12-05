package srv

import (
	"github.com/meschbach/elevatinator/ipc/grpc/telepathy/pb"
	"github.com/meschbach/elevatinator/simulator"
	"google.golang.org/grpc"
	"net"
)

// Network is an abstraction for binding to a listener.  Closes the gap between real world usage of the code and testing
// paths.
type Network interface {
	Listener() (net.Listener, error)
}

// RunControllerService exports the given controller on port tcp/9998 over gRPC
func RunControllerService(builder simulator.ControllerFunc) error {
	l := &tcp{listenAt: ":9998"}
	return RunControllerOn(builder, l)
}

// RunControllerOn exports the given controller on the specified network address
func RunControllerOn(builder simulator.ControllerFunc, on Network) error {
	s := grpc.NewServer()
	pb.RegisterControllerServiceServer(s, newRemoteController(builder))
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
