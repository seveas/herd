package common

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"time"

	"github.com/seveas/herd"

	plugin "github.com/hashicorp/go-plugin"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCClient struct {
	broker *plugin.GRPCBroker
	client ProviderPluginClient
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

func (c *GRPCClient) SetDataDir(dir string) error {
	resp, err := c.client.SetDataDir(c.ctx, &SetDataDirRequest{Dir: dir})
	if err != nil {
		return err
	}
	if resp.Err != "" {
		return errors.New(resp.Err)
	}
	return nil
}

func (c *GRPCClient) SetCacheDir(dir string) {
	_, _ = c.client.SetCacheDir(c.ctx, &SetCacheDirRequest{Dir: dir})
}

func (c *GRPCClient) Invalidate() {
	_, _ = c.client.Invalidate(c.ctx, &Empty{})
}

func (c *GRPCClient) Keep() {
	_, _ = c.client.Keep(c.ctx, &Empty{})
}

func (c *GRPCClient) SetLogger(logger Logger) error {
	loggerServer := &GRPCLoggerServer{Impl: logger}
	var s *grpc.Server
	serverFunc := func(opts []grpc.ServerOption) *grpc.Server {
		s = grpc.NewServer(opts...)
		RegisterLoggerServer(s, loggerServer)
		return s
	}

	id := c.broker.NextId()
	go c.broker.AcceptAndServe(id, serverFunc)

	resp, err := c.client.SetLogger(c.ctx, &SetLoggerRequest{Logger: id})
	if err != nil {
		return err
	}
	if resp.Err != "" {
		return errors.New(resp.Err)
	}
	return nil
}

func (c *GRPCClient) Load(ctx context.Context) (*herd.HostSet, error) {
	deadline, ok := ctx.Deadline()
	if !ok {
		deadline = time.Now().Add(10 * time.Minute)
	}
	ts := timestamppb.New(deadline)
	resp, err := c.client.Load(c.ctx, &LoadRequest{Deadline: ts})
	if err != nil {
		return nil, err
	}
	if resp.Err != "" {
		return nil, errors.New(resp.Err)
	}
	hosts := herd.NewHostSet()
	if err := json.Unmarshal(resp.Data, &hosts); err != nil {
		return nil, err
	}
	/* FIXME set maxhostlen? Maybe have a custom json unserializer? */
	return hosts, nil
}

type GRPCServer struct {
	UnimplementedProviderPluginServer
	Impl   ProviderPluginImpl
	broker *plugin.GRPCBroker
}

func (s *GRPCServer) SetLogger(ctx context.Context, req *SetLoggerRequest) (*SetLoggerResponse, error) {
	conn, err := s.broker.Dial(req.Logger)
	if err != nil {
		return &SetLoggerResponse{Err: err.Error()}, nil //nolint:nilerr // The error is returned in the response
	}
	logger := &GRPCLoggerClient{NewLoggerClient(conn)}
	logrus.SetOutput(io.Discard)
	logrus.AddHook(logger)
	if err := s.Impl.SetLogger(logger); err != nil {
		return &SetLoggerResponse{Err: err.Error()}, nil //nolint:nilerr // The error is returned in the response
	}
	return &SetLoggerResponse{}, nil
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

func (s *GRPCServer) SetDataDir(ctx context.Context, req *SetDataDirRequest) (*SetDataDirResponse, error) {
	if err := s.Impl.SetDataDir(req.Dir); err != nil {
		return &SetDataDirResponse{Err: err.Error()}, nil //nolint:nilerr // The error is returned in the response
	}
	return &SetDataDirResponse{}, nil
}

func (s *GRPCServer) SetCacheDir(ctx context.Context, req *SetCacheDirRequest) (*Empty, error) {
	s.Impl.SetCacheDir(req.Dir)
	return &Empty{}, nil
}

func (s *GRPCServer) Invalidate(ctx context.Context, req *Empty) (*Empty, error) {
	s.Impl.Invalidate()
	return &Empty{}, nil
}

func (s *GRPCServer) Keep(ctx context.Context, req *Empty) (*Empty, error) {
	s.Impl.Keep()
	return &Empty{}, nil
}

func (s *GRPCServer) Load(ctx context.Context, req *LoadRequest) (*LoadResponse, error) {
	ts := req.Deadline.AsTime()
	ctx, cancel := context.WithDeadline(ctx, ts)
	defer cancel()
	hosts, err := s.Impl.Load(ctx)
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

func (c *GRPCLoggerClient) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (c *GRPCLoggerClient) EmitLogMessage(level logrus.Level, message string) {
	_, _ = c.client.EmitLogMessage(context.Background(), &EmitLogMessageRequest{Level: uint32(level), Message: message})
}

func (c *GRPCLoggerClient) Fire(entry *logrus.Entry) error {
	c.EmitLogMessage(entry.Level, entry.Message)
	return nil
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

var (
	_ Logger             = &GRPCLoggerClient{}
	_ ProviderPluginImpl = &GRPCClient{}
)
