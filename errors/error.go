package errors

import (
	"context"
	"fmt"

	"github.com/claudioluciano/goutils/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type errorKind string

var (
	kindValidation    errorKind = "ValidationError"
	kindInternal      errorKind = "InternalError"
	kindNotFound      errorKind = "NotFoundError"
	kindAlreadyExists errorKind = "AlreadyExistsError"
)

type Client struct {
	serviceName string
	logger      *logger.Client
}

type ClientOptions struct {
	ServiceName string
	Logger      *logger.Client
}

func NewClient(opts *ClientOptions) *Client {
	return &Client{
		serviceName: opts.ServiceName,
		logger:      opts.Logger,
	}
}

func (c *Client) NotFound(fields ...interface{}) error {
	msg := fmt.Sprintf("%s not found", c.serviceName)

	logFields := append(fields, fmt.Sprintf("%s error.kind", kindNotFound))
	c.logger.Warn(msg, logFields...)

	return status.Error(codes.NotFound, msg)
}

func (c *Client) Internal(err error, fields ...interface{}) error {
	msg := fmt.Sprintf("%s got an internal error", c.serviceName)

	logFields := append(fields, fmt.Sprintf("%s error.kind", kindInternal))
	logFields = append(logFields, "error")
	logFields = append(logFields, err)

	c.logger.Error(msg, logFields...)

	return status.Error(codes.Internal, err.Error())
}

func (c *Client) AlreadyExists(fields ...interface{}) error {
	msg := fmt.Sprintf("%s already exists", c.serviceName)

	logFields := append(fields, fmt.Sprintf("%s error.kind", kindAlreadyExists))

	c.logger.Warn(msg, logFields...)

	return status.Error(codes.Internal, msg)
}

func (c *Client) InvalidArgument(ctx context.Context, err error, fields ...interface{}) error {
	msg := "request validation failed"

	logFields := append(fields, fmt.Sprintf("%s error.kind", kindValidation))
	c.logger.Warn(msg, logFields...)

	return status.Error(codes.InvalidArgument, err.Error())
}
