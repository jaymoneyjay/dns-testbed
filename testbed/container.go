package testbed

import (
	"bytes"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/rs/zerolog"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Container struct {
	client   client.APIClient
	ctx      context.Context
	logger   zerolog.Logger
	ID       string
	dir      string
	queryLog string
	ip       string
}

func NewContainer(id, dir, ip string) *Container {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	executionLog, err := os.Create(filepath.Join(dir, "execution.log"))
	if err != nil {
		panic(err)
	}
	output := zerolog.ConsoleWriter{Out: executionLog, TimeFormat: time.RFC3339}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("***\n%s****", i)
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
	}
	logger := zerolog.New(output).With().Timestamp().Logger()
	queryLog := filepath.Join(dir, "query.log")
	_, err = os.Create(queryLog)
	if err != nil {
		panic(err)
	}
	return &Container{
		client:   cli,
		ctx:      context.Background(),
		logger:   logger,
		ID:       id,
		dir:      dir,
		ip:       ip,
		queryLog: queryLog,
	}
}

type ExecResult struct {
	StdOut   string
	StdErr   string
	ExitCode int
}

func (c *Container) Exec(cmd []string) (ExecResult, error) {
	execConfig := types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          cmd,
	}
	createResp, err := c.client.ContainerExecCreate(c.ctx, c.ID, execConfig)
	if err != nil {
		return ExecResult{}, err
	}
	execResp, err := c.inspectExecResp(createResp.ID)
	if err != nil {
		return ExecResult{}, err
	}
	c.logger.Info().
		Str("containerID", c.ID).
		Msg(execResp.StdOut)
	return execResp, nil
}

func (c *Container) inspectExecResp(execID string) (ExecResult, error) {
	attachResp, err := c.client.ContainerExecAttach(c.ctx, execID, types.ExecStartCheck{})
	if err != nil {
		return ExecResult{}, err
	}
	defer attachResp.Close()
	var outBuf, errBuf bytes.Buffer
	outputDone := make(chan error)

	go func() {
		_, err = stdcopy.StdCopy(&outBuf, &errBuf, attachResp.Reader)
		outputDone <- err
	}()
	select {
	case err := <-outputDone:
		if err != nil {
			return ExecResult{}, err
		}
		break
	case <-c.ctx.Done():
		return ExecResult{}, c.ctx.Err()
	}
	stdout, err := io.ReadAll(&outBuf)
	if err != nil {
		return ExecResult{}, err
	}
	stderr, err := io.ReadAll(&errBuf)
	if err != nil {
		return ExecResult{}, err
	}
	res, err := c.client.ContainerExecInspect(c.ctx, execID)
	if err != nil {
		return ExecResult{}, err
	}
	return ExecResult{
		ExitCode: res.ExitCode,
		StdOut:   string(stdout),
		StdErr:   string(stderr),
	}, nil
}

func (c *Container) ReadQueryLog(minTimeout time.Duration) []byte {
	var lines []string
	numberOfCurrentLines := 0
	for true {
		time.Sleep(minTimeout)
		queryLog, err := os.ReadFile(c.queryLog)
		queryLog = bytes.ReplaceAll(queryLog, []byte{'\x00'}, []byte{})
		if err != nil {
			panic(err)
		}
		lines = strings.Split(string(queryLog), "\n")
		if len(lines) == numberOfCurrentLines {
			break
		}
		numberOfCurrentLines = len(lines)
	}
	return []byte(strings.Join(lines, "\n"))
}

func (c *Container) FlushQueryLog() {
	_, err := os.Create(c.queryLog)
	if err != nil {
		panic(err)
	}
}
