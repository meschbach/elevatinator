package systest

import (
	"context"
	"github.com/meschbach/elevatinator/pkg/controllers/queue"
	"github.com/meschbach/elevatinator/pkg/ipc/grpc/telepathy"
	"github.com/meschbach/elevatinator/pkg/ipc/grpc/telepathy/srv"
	"github.com/meschbach/elevatinator/pkg/junk/grpctest"
	"github.com/meschbach/elevatinator/pkg/scenarios"
	"github.com/stretchr/testify/require"
	"net"
	"testing"
	"time"
)

func TestSystem(t *testing.T) {
	ctx, done := context.WithTimeout(context.Background(), 2*time.Second)
	t.Cleanup(done)

	virtualNetwork := &testNetwork{transport: grpctest.NewBufferTransport()}
	go func() {
		if err := srv.RunControllerOn(queue.NewController, virtualNetwork); err != nil {
			require.NoError(t, err)
		}
	}()

	conn, err := virtualNetwork.transport.GRPCClient(ctx)
	require.NoError(t, err)

	landing := telepathy.LandingWithConnection(conn)
	controller := landing.ControllerAdapter()

	t.Run("Single Person Up", func(t *testing.T) {
		scenarios.RunScenario(controller, scenarios.SinglePersonUp)
	})

	t.Run("Single Person Down", func(t *testing.T) {
		scenarios.RunScenario(controller, scenarios.SinglePersonDown)
	})

	t.Run("Multiple players and back", func(t *testing.T) {
		scenarios.RunScenario(controller, scenarios.MultipleUpAndBack)
	})
}

type testNetwork struct {
	transport *grpctest.BufferTransport
}

func (t *testNetwork) Listener() (net.Listener, error) {
	return t.transport.Listener, nil
}
