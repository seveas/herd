package ssh

import (
	"context"
	"fmt"
	"os/user"
	"time"

	"github.com/seveas/herd"
)

type PingExecutor struct {
	Executor
}

func NewPingExecutor(agentTimeout time.Duration, user user.User) (herd.Executor, error) {
	executor, err := NewExecutor(agentTimeout, user)
	if err != nil {
		return nil, err
	}
	return &PingExecutor{Executor: *(executor.(*Executor))}, nil
}

func (e *PingExecutor) Run(ctx context.Context, host *herd.Host, command string, oc chan herd.OutputLine) *herd.Result {
	now := time.Now()
	r := &herd.Result{Host: host.Name, StartTime: now, EndTime: now, ElapsedTime: 0, ExitStatus: -1}
	defer func() {
		r.EndTime = time.Now()
		r.ElapsedTime = r.EndTime.Sub(r.StartTime).Seconds()
	}()

	if err := ctx.Err(); err != nil {
		r.Err = err
		return r
	}
	connection, err := e.connect(ctx, host)
	if err != nil {
		r.Err = err
		return r
	}
	defer connection.Close()
	r.ExitStatus = 0
	pong := []byte(fmt.Sprintf("connection successful in %s\n", time.Since(now).Truncate(time.Millisecond)))
	if oc != nil {
		oc <- herd.OutputLine{Host: host, Data: pong}
	}
	r.Stdout = pong
	return r
}
