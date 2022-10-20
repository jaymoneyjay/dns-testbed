package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type Client struct {
	client client.APIClient
	ctx    context.Context
}

func NewClient() *Client {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	return &Client{client: cli, ctx: context.Background()}
}

func (cli *Client) Exec(containerID string, cmd []string) error {
	execConfig := types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          cmd,
	}
	createResp, err := cli.client.ContainerExecCreate(cli.ctx, containerID, execConfig)
	if err != nil {
		return err
	}
	execID := createResp.ID

	attachResp, err := cli.client.ContainerExecAttach(cli.ctx, execID, types.ExecStartCheck{})
	if err != nil {
		return err
	}
	defer attachResp.Close()
	return nil
}
