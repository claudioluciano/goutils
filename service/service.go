package service

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/claudioluciano/goutils/database/gorm"
	"github.com/claudioluciano/goutils/database/mongo"
	"github.com/claudioluciano/goutils/errors"
	"github.com/claudioluciano/goutils/logger"
	"github.com/lithammer/shortuuid"
	"google.golang.org/grpc"
)

type ServiceEnvironment string
type ServiceDatabaseType string

const (
	defaultPORT     int32               = 50051
	ENV_PRODUCTION  ServiceEnvironment  = "PRODUCTION"
	ENV_TEST        ServiceEnvironment  = "TEST"
	ENV_DEVELOPMENT ServiceEnvironment  = "DEVELOPMENT"
	DATABASE_MONGO  ServiceDatabaseType = "MONGO"
	DATABASE_GORM   ServiceDatabaseType = "GORM"
)

type ServiceDatabase struct {
	gorm  *gorm.Client
	mongo *mongo.Client
}

type Service struct {
	name        string
	port        int32
	environment ServiceEnvironment
	grpcServer  *grpc.Server
	logger      *logger.Client
	errors      *errors.Client
	db          *ServiceDatabase
}

type NewServiceOptions struct {
	ServiceName string
	Environment ServiceEnvironment
	Database    *DatabaseOptions
}

type GormOptions struct {
	AutoMigration bool
	Migrations    []interface{}
}

type DatabaseOptions struct {
	Disabled bool
	Type     ServiceDatabaseType
	*GormOptions
}

func NewService(opts ...*NewServiceOptions) (*Service, error) {
	opt := &NewServiceOptions{
		Environment: ENV_DEVELOPMENT,
		Database: &DatabaseOptions{
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

	errs := errors.NewClient(&errors.ClientOptions{
		ServiceName: opt.ServiceName,
		Logger:      lg,
	})

	svc := &Service{
		name:        opt.ServiceName,
		port:        defaultPORT,
		environment: opt.Environment,
		logger:      lg,
		errors:      errs,
		grpcServer:  grpc.NewServer(),
	}

	if !opt.Database.Disabled {
		if opt.Database.Type == DATABASE_GORM {
			if err := svc.gormInitialize(opt); err != nil {
				return nil, err
			}
		}

		if opt.Database.Type == DATABASE_MONGO {
			if err := svc.mongoInitialize(opt); err != nil {
				return nil, err
			}
		}
	}

	return svc, nil
}

func (s *Service) gormInitialize(opt *NewServiceOptions) error {
	var db *gorm.Client
	if opt.Environment == ENV_PRODUCTION {
		idb, err := gorm.NewPostgres(&gorm.NewPostgresOpts{
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
		idb, err := gorm.NewSqlite(&gorm.NewSqliteOpts{
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

	s.db = &ServiceDatabase{
		gorm: db,
	}
	return nil
}

func (s *Service) mongoInitialize(opt *NewServiceOptions) error {
	db, err := mongo.NewClient(&mongo.NewClientOptions{
		Host:           getEnv("MONGO_HOST", "localhost"),
		Port:           getEnv("MONGO_PORT", "27017"),
		User:           getEnv("MONGO_USER", "root"),
		Password:       getEnv("MONGO_PASSWORD", "qwerty"),
		DatabaseName:   getEnv("MONGO_DB_NAME", "mydb"),
		CollectionName: s.name,
		Logger:         s.logger,
		Errors:         s.errors,
	})
	if err != nil {
		return err
	}

	s.db = &ServiceDatabase{
		mongo: db,
	}

	return nil
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

func (s *Service) Gorm() *gorm.Client {
	return s.db.gorm
}

func (s *Service) Mongo() *mongo.Client {
	return s.db.mongo
}

func (s *Service) GRPCServer() *grpc.Server {
	return s.grpcServer
}

func (s *Service) Logger() *logger.Client {
	return s.logger
}

func (s *Service) getPortAsString() string {
	return fmt.Sprintf(":%v", s.port)
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

func (db *Service) NewID(prefix string) string {
	var newID string

	if prefix != "" {
		prefix += "_"
	}

	newID = shortuuid.New()

	return prefix + newID
}
