package telepathy

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"time"
)

func CheckHealth(ctx context.Context, serviceAddress string) (bool, error) {
	conn, err := grpc.Dial(serviceAddress, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(1*time.Second))
	if err != nil {
		return false, &ConnectionError{
			Target:     serviceAddress,
			Underlying: err,
		}
	}

	defer conn.Close()

	service := grpc_health_v1.NewHealthClient(conn)
	response, err := service.Check(ctx, &grpc_health_v1.HealthCheckRequest{Service: "ai"})
	if err != nil {
		return false, err
	}

	switch response.Status {
	case grpc_health_v1.HealthCheckResponse_SERVING:
		return true, nil
	case grpc_health_v1.HealthCheckResponse_NOT_SERVING:
		return false, nil
	default:
		return false, fmt.Errorf("Unknown health status %q\n", response.Status.String())
	}
}
