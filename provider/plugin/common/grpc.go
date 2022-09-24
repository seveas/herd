package common

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/seveas/herd"

	"github.com/golang/protobuf/ptypes"
	plugin "github.com/hashicorp/go-plugin"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type GRPCClient struct {
	broker *plugin.GRPCBroker
	client ProviderClient
	ctx    context.Context
}

func (c *GRPCClient) Configure(settings map[string]interface{}) error {
	data, err := json.Marshal(settings)
	if err != nil {
		return err
	}
	resp, err := c.client.Configure(c.ctx, &ConfigureRequest{Data: data})
	if err != nil {
		return err
	}
	if resp.Err != "" {
		return errors.New(resp.Err)
	}
	return nil
}

func (c *GRPCClient) Load(ctx context.Context, logger Logger) (herd.Hosts, error) {
	loggerServer := &GRPCLoggerServer{Impl: logger}
	var s *grpc.Server
	serverFunc := func(opts []grpc.ServerOption) *grpc.Server {
		s = grpc.NewServer(opts...)
		RegisterLoggerServer(s, loggerServer)
		return s
	}

	id := c.broker.NextId()
	go c.broker.AcceptAndServe(id, serverFunc)
	defer func() {
		if s != nil {
			s.Stop()
		}
	}()

	deadline, ok := ctx.Deadline()
	if !ok {
		deadline = time.Now().Add(10 * time.Minute)
	}
	ts, _ := ptypes.TimestampProto(deadline)
	resp, err := c.client.Load(c.ctx, &LoadRequest{Deadline: ts, Logger: id})
	if err != nil {
		return nil, err
	}
	if resp.Err != "" {
		return nil, errors.New(resp.Err)
	}
	var hosts herd.Hosts
	if err := json.Unmarshal(resp.Data, &hosts); err != nil {
		return nil, err
	}
	return hosts, nil
}

type GRPCServer struct {
	UnimplementedProviderServer
	Impl   Provider
	broker *plugin.GRPCBroker
}

func (s *GRPCServer) Configure(ctx context.Context, req *ConfigureRequest) (*ConfigureResponse, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(req.Data, &data); err != nil {
		return nil, err
	}
	if err := s.Impl.Configure(data); err != nil {
		return &ConfigureResponse{Err: err.Error()}, nil //nolint:nilerr // The error is returned in the response
	}
	return &ConfigureResponse{}, nil
}

func (s *GRPCServer) Load(ctx context.Context, req *LoadRequest) (*LoadResponse, error) {
	conn, err := s.broker.Dial(req.Logger)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	ts, _ := ptypes.Timestamp(req.Deadline)
	ctx, cancel := context.WithDeadline(ctx, ts)
	defer cancel()
	l := &GRPCLoggerClient{NewLoggerClient(conn)}
	hosts, err := s.Impl.Load(ctx, l)
	if err != nil {
		return &LoadResponse{Err: err.Error()}, nil //nolint:nilerr // The error is returned in the response
	}
	data, err := json.Marshal(hosts)
	if err != nil {
		return nil, err
	}
	return &LoadResponse{Data: data}, nil
}

type GRPCLoggerClient struct {
	client LoggerClient
}

func (c *GRPCLoggerClient) LoadingMessage(name string, done bool, err error) {
	errs := ""
	if err != nil {
		errs = err.Error()
	}
	_, _ = c.client.LoadingMessage(context.Background(), &LoadingMessageRequest{Name: name, Done: done, Err: errs})
}

func (c *GRPCLoggerClient) EmitLogMessage(level logrus.Level, message string) {
	_, _ = c.client.EmitLogMessage(context.Background(), &EmitLogMessageRequest{Level: uint32(level), Message: message})
}

type GRPCLoggerServer struct {
	UnimplementedLoggerServer
	Impl Logger
}

func (s *GRPCLoggerServer) LoadingMessage(ctx context.Context, req *LoadingMessageRequest) (*Empty, error) {
	var err error = nil
	if req.Err != "" {
		err = errors.New(req.Err)
	}
	s.Impl.LoadingMessage(req.Name, req.Done, err)
	return &Empty{}, nil
}

func (s *GRPCLoggerServer) EmitLogMessage(ctx context.Context, req *EmitLogMessageRequest) (*Empty, error) {
	s.Impl.EmitLogMessage(logrus.Level(req.Level), req.Message)
	return &Empty{}, nil
}

var _ Logger = &GRPCLoggerClient{}

var _ Provider = &GRPCClient{}
