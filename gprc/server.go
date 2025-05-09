package gprc

import (
	"context"
	"log"
	"net"

	"github.com/Knetic/govaluate"
	"google.golang.org/grpc"
	pb "serverCalc/proto"
	"strconv"
)

type server struct {
	pb.UnimplementedTaskServiceServer
}

func (s *server) Calculate(ctx context.Context, req *pb.TaskRequest) (*pb.TaskResult, error) {
	expr, err := govaluate.NewEvaluableExpression(req.Expression)
	if err != nil {
		return &pb.TaskResult{Error: err.Error()}, nil
	}
	result, err := expr.Evaluate(nil)
	if err != nil {
		return &pb.TaskResult{Error: err.Error()}, nil
	}

	return &pb.TaskResult{Result: strconv.FormatFloat(result.(float64), 'f', -1, 64)}, nil
}

func StartGRPCServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterTaskServiceServer(s, &server{})

	log.Println("gRPC server is running on port 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
