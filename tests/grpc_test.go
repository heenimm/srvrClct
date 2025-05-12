package tests_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	pb "serverCalc/proto"
	"testing"
	"time"
)

func TestGRPCServer(t *testing.T) {
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterTaskServiceServer(grpcServer, &server{})
	reflection.Register(grpcServer)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			t.Fatalf("failed to serve: %v", err)
		}
	}()
	time.Sleep(time.Second)

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewTaskServiceClient(conn)

	t.Run("Test AddExpression", func(t *testing.T) {
		req := &pb.CalculationRequest{
			Expression: "5 + 3",
		}

		resp, err := client.Calculate(context.Background(), req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "8", resp.Result)
	})

	t.Run("Test GetExpressions", func(t *testing.T) {
		resp, err := client.GetExpressions(context.Background(), &pb.Empty{})
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Expressions, 0)
	})

	t.Run("Test Calculate with invalid expression", func(t *testing.T) {
		req := &pb.CalculationRequest{
			Expression: "invalid expression",
		}

		resp, err := client.Calculate(context.Background(), req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Contains(t, resp.Error, "Invalid expression")
	})
}

type server struct {
	pb.UnimplementedTaskServiceServer
}
