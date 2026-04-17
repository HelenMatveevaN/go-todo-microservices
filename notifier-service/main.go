package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "github.com/HelenMatveevaN/notifier-proj/api/proto" // Путь к сгенерированным файлам
)

// server — структура, которая реализует интерфейс Notifier
type server struct {
	pb.UnimplementedNotifierServer
}

// SendNotification — та самая логика, которую мы описали в .proto
func (s *server) SendNotification(ctx context.Context, req *pb.NotificationRequest) (*pb.NotificationResponse, error) {
	log.Printf("ПОЛУЧЕНО УВЕДОМЛЕНИЕ: Задача '%s' — %s", req.TaskTitle, req.Message)
	return &pb.NotificationResponse{Success: true}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051") // Слушаем порт 50051
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterNotifierServer(s, &server{})

	log.Println("gRPC Server Notifier запущен на порту :50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}