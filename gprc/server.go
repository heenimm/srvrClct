package gprc

import (
	"context"
	"log"
	"net"
	"strconv"
	"sync"

	"github.com/Knetic/govaluate"
	"google.golang.org/grpc"
	pb "serverCalc/proto"
)

type server struct {
	pb.UnimplementedTaskServiceServer
	mu          sync.Mutex
	expressions map[string]string
	nextID      int
}

func (s *server) AddExpression(ctx context.Context, req *pb.CalculationRequest) (*pb.CalculationResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := strconv.Itoa(s.nextID)
	s.nextID++
	s.expressions[id] = req.Expression

	return &pb.CalculationResponse{Id: id}, nil
}

func (s *server) GetExpressions(ctx context.Context, _ *pb.Empty) (*pb.ExpressionListResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var list []*pb.Expression
	for id, expr := range s.expressions {
		list = append(list, &pb.Expression{Id: id, Expression: expr})
	}
	return &pb.ExpressionListResponse{Expressions: list}, nil
}

func (s *server) GetExpressionByID(ctx context.Context, req *pb.GetByIDRequest) (*pb.CalculationRequest, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	expr, ok := s.expressions[req.Id]
	if !ok {
		return &pb.CalculationRequest{Expression: ""}, nil
	}
	return &pb.CalculationRequest{Expression: expr}, nil
}

func (s *server) Calculate(ctx context.Context, req *pb.CalculationRequest) (*pb.CalculationResponse, error) {
	expr, err := govaluate.NewEvaluableExpression(req.Expression)
	if err != nil {
		return &pb.CalculationResponse{Error: "Invalid expression: " + err.Error()}, nil
	}
	result, err := expr.Evaluate(nil)
	if err != nil {
		return &pb.CalculationResponse{Error: "Calculation error: " + err.Error()}, nil
	}
	return &pb.CalculationResponse{Result: strconv.FormatFloat(result.(float64), 'f', -1, 64)}, nil
}

func (s *server) GetTask(ctx context.Context, _ *pb.Empty) (*pb.TaskResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for id, expr := range s.expressions {
		return &pb.TaskResponse{Id: id, Expression: expr}, nil
	}
	return &pb.TaskResponse{Error: "no tasks available"}, nil
}

func (s *server) SubmitTaskResult(ctx context.Context, req *pb.TaskResultRequest) (*pb.TaskSubmitResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.expressions[req.Id]; !ok {
		return &pb.TaskSubmitResponse{Status: "failed", Error: "task not found"}, nil
	}
	delete(s.expressions, req.Id)
	return &pb.TaskSubmitResponse{Status: "ok"}, nil
}

func StartGRPCServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterTaskServiceServer(s, &server{
		expressions: make(map[string]string),
		nextID:      1,
	})

	log.Println("gRPC server is running on port 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
