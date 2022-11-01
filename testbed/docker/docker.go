package docker

import (
	"bytes"
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"io"
	"time"
)

type Client struct {
	client client.APIClient
	ctx    context.Context
}

func NewClient() (*Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	return &Client{client: cli, ctx: context.Background()}, nil
}

type ExecResult struct {
	StdOut   string
	StdErr   string
	ExitCode int
}

func (cli *Client) Exec(containerID string, cmd []string) (ExecResult, error) {
	execConfig := types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          cmd,
	}
	createResp, err := cli.client.ContainerExecCreate(cli.ctx, containerID, execConfig)
	if err != nil {
		return ExecResult{}, err
	}
	return cli.inspectExecResp(createResp.ID)
}

func (cli *Client) inspectExecResp(execID string) (ExecResult, error) {
	attachResp, err := cli.client.ContainerExecAttach(cli.ctx, execID, types.ExecStartCheck{})
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
	case <-cli.ctx.Done():
		return ExecResult{}, cli.ctx.Err()
	}
	stdout, err := io.ReadAll(&outBuf)
	if err != nil {
		return ExecResult{}, err
	}
	stderr, err := io.ReadAll(&errBuf)
	if err != nil {
		return ExecResult{}, err
	}
	res, err := cli.client.ContainerExecInspect(cli.ctx, execID)
	if err != nil {
		return ExecResult{}, err
	}
	return ExecResult{
		ExitCode: res.ExitCode,
		StdOut:   string(stdout),
		StdErr:   string(stderr),
	}, nil
}

func (cli *Client) RestartContainer(containerID string) error {
	duration := 5 * time.Second
	return cli.client.ContainerRestart(cli.ctx, containerID, &duration)
}
