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
	defaultPORT int = 50051
)

type Service struct {
	port       int
	grpcServer *grpc.Server
	*logger.Logger
	*db.DB
}

type NewServiceOpts struct {
	migrations []interface{}
}

func NewService(opts ...NewServiceOpts) (*Service, error) {
	opt := NewServiceOpts{
		migrations: nil,
	}

	if len(opts) != 0 {
		opt = opts[0]
	}

	lg := logger.NewLogger(nil)

	svc := &Service{
		port:       defaultPORT,
		Logger:     lg,
		grpcServer: grpc.NewServer(),
	}

	err := svc.dbInitialize(&db.NewPostgresOpts{
		Logger:   lg,
		Host:     getEnv("POSTGRES_HOST", "localhost"),
		Port:     getEnv("POSTGRES_PORT", "5432"),
		DbName:   getEnv("POSTGRES_DBNAME", "MyDB"),
		User:     getEnv("POSTGRES_DBUSER", "root"),
		Password: getEnv("POSTGRES_DBPASSWORD", "qwerty"),
	}, opt.migrations...)
	if err != nil {
		return nil, err
	}

	return svc, nil
}

func (s *Service) dbInitialize(opts *db.NewPostgresOpts, migrations ...interface{}) error {
	db, err := db.NewPostgres(opts)
	if err != nil {
		return err
	}

	if migrations != nil {
		mErr := db.AutoMigrate(migrations)
		if mErr != nil {
			return mErr
		}
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

func getEnv(name string, defaultValue string) string {
	value := os.Getenv(name)

	if defaultValue != "" && value == "" {
		return defaultValue
	}

	return value
}
