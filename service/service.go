package service

import (
	"fmt"
	"net"
	"os"

	"github.com/claudioluciano/goutils/db"
	"github.com/claudioluciano/goutils/logger"
	"google.golang.org/grpc"
)

const (
	PORT int = 50051
)

type Service struct {
	port       int
	grpcServer *grpc.Server
	*logger.Logger
	*db.DB
}

type NewServiceOpts struct {
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
}

func NewService(opts *NewServiceOpts, migrations ...interface{}) (*Service, error) {
	lg := logger.NewLogger(nil)

	svc := &Service{
		port:       PORT,
		Logger:     lg,
		grpcServer: grpc.NewServer(),
	}

	err := svc.dbInitialize(&db.NewPostgressOpts{
		Host:     GetEnv("POSTGRES_HOST", "localhost"),
		Port:     GetEnv("POSTGRES_POST", "5432"),
		DbName:   GetEnv("POSTGRES_DBNAME", "MyDB"),
		User:     GetEnv("POSTGRES_DBNAME", "root"),
		Password: GetEnv("POSTGRES_DBNAME", "qwerty"),
	}, migrations)
	if err != nil {
		return nil, err
	}

	return svc, nil
}

func (s *Service) dbInitialize(opts *db.NewPostgressOpts, migrations ...interface{}) error {
	db, err := db.NewPostgress(opts)
	if err != nil {
		return err
	}

	mErr := s.DB.AutoMigrate(migrations)
	if mErr != nil {
		return mErr
	}

	s.DB = db
	return nil
}

func (s *Service) getPortAsString() string {
	return fmt.Sprintf(":%v", s.port)
}

func (s *Service) GRPCServer() *grpc.Server {
	return s.grpcServer
}

func (s *Service) ListenAndServe() error {
	lis, err := net.Listen("tcp", s.getPortAsString())
	if err != nil {
		s.Logger.ErrorWithError("failed to listen", err)
		return err
	}

	if err := s.grpcServer.Serve(lis); err != nil {
		s.Logger.ErrorWithError("failed to serve", err)
		return err
	}

	s.Logger.Info(fmt.Sprintf("Listen port %v", s.port))

	return nil
}

func (s *Service) Stop() {
	s.grpcServer.GracefulStop()
}

func GetEnv(name string, defaultValue string) string {
	value := os.Getenv(name)

	if defaultValue != "" && value == "" {
		return defaultValue
	}

	return value
}
