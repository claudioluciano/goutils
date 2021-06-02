package service

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/claudioluciano/goutils/database"
	"github.com/claudioluciano/goutils/logger"
	"google.golang.org/grpc"
)

type ServiceEnvironment string

const (
	defaultPORT     int32              = 50051
	ENV_PRODUCTION  ServiceEnvironment = "PRODUCTION"
	ENV_TEST        ServiceEnvironment = "TEST"
	ENV_DEVELOPMENT ServiceEnvironment = "DEVELOPMENT"
)

type Service struct {
	name        string
	port        int32
	environment ServiceEnvironment
	grpcServer  *grpc.Server
	logger      *logger.Logger
	db          *database.DB
}

type NewServiceOpts struct {
	ServiceName string
	Environment ServiceEnvironment
	Database    *DatabaseOpts
}

type DatabaseOpts struct {
	Disabled      bool
	AutoMigration bool
	Migrations    []interface{}
}

func NewService(opts ...*NewServiceOpts) (*Service, error) {
	opt := &NewServiceOpts{
		Environment: ENV_DEVELOPMENT,
		Database: &DatabaseOpts{
			Disabled: true,
		},
	}

	if len(opts) > 0 {
		opt = opts[0]
	}

	lg := logger.NewLogger(&logger.NewLoggerOpts{
		Name:  opt.ServiceName,
		Level: logger.LevelInfo(),
	})

	svc := &Service{
		name:        opt.ServiceName,
		port:        defaultPORT,
		environment: opt.Environment,
		logger:      lg,
		grpcServer:  grpc.NewServer(),
	}

	if !opt.Database.Disabled {
		if err := svc.dbInitialize(opt); err != nil {
			return nil, err
		}
	}

	return svc, nil
}

func (s *Service) dbInitialize(opt *NewServiceOpts) error {
	var db *database.DB
	if opt.Environment == ENV_PRODUCTION {
		idb, err := database.NewPostgres(&database.NewPostgresOpts{
			Table:    strings.ToLower(s.name),
			Logger:   s.Logger(),
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     getEnv("POSTGRES_PORT", "5432"),
			DBName:   getEnv("POSTGRES_DBNAME", "MyDB"),
			User:     getEnv("POSTGRES_DBUSER", "root"),
			Password: getEnv("POSTGRES_DBPASSWORD", "qwerty"),
		})
		if err != nil {
			return err
		}

		db = idb
	}

	if opt.Environment != ENV_PRODUCTION {
		idb, err := database.NewSqlite(&database.NewSqliteOpts{
			Table:  strings.ToLower(s.name),
			Logger: s.Logger(),
			DBName: strings.ToLower(fmt.Sprintf("%v.db", s.name)),
		})
		if err != nil {
			return err
		}

		db = idb
	}

	if opt.Database.AutoMigration && opt.Database.Migrations != nil {
		err := db.AutoMigrate(opt.Database.Migrations...)
		if err != nil {
			return err
		}
	}

	s.db = db
	return nil
}

func (s *Service) getPortAsString() string {
	return fmt.Sprintf(":%v", s.port)
}

func (s *Service) ClientConnection(name string) *grpc.ClientConn {
	address := fmt.Sprintf("%s:%d", name, defaultPORT)
	conn, err := grpc.Dial(
		address,
		grpc.WithInsecure(),
	)
	if err != nil {
		s.logger.Fatal("could not create client connection", "name", name, "address", address, "port", defaultPORT)
		return nil
	}

	return conn
}

func (s *Service) GRPCServer() *grpc.Server {
	return s.grpcServer
}

func (s *Service) Logger() *logger.Logger {
	return s.logger
}

func (s *Service) ListenAndServe() error {
	lis, err := net.Listen("tcp", s.getPortAsString())
	if err != nil {
		s.logger.Error("failed to listen", err)
		return err
	}

	if err := s.grpcServer.Serve(lis); err != nil {
		s.logger.Error("failed to serve", err)
		return err
	}

	s.logger.Info(fmt.Sprintf("Listen port %v", s.port))

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
