package srv

import (
	"github.com/meschbach/elevatinator/ipc/grpc/telepathy/pb"
	"github.com/meschbach/elevatinator/simulator"
	"google.golang.org/grpc"
	"log"
	"net"
)

func RunControllerService(builder simulator.ControllerFunc)  {
	port := ":9998"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterControllerServiceServer(s, newRemoteController(builder))
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
