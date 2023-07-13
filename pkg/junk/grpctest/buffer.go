package grpctest

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"net"
)

type BufferTransport struct {
	Listener *bufconn.Listener
}

func NewBufferTransport() *BufferTransport {
	const bufSize = 1024 * 1024

	var listener *bufconn.Listener
	listener = bufconn.Listen(bufSize)

	return &BufferTransport{
		Listener: listener,
	}
}

func (b *BufferTransport) GRPCClient(ctx context.Context) (conn *grpc.ClientConn, err error) {
	dialer := func(ctx context.Context, _ string) (net.Conn, error) {
		return b.Listener.DialContext(ctx)
	}
	return grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(dialer), grpc.WithInsecure())
}
