package docker

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

type Client struct {
	client   client.APIClient
	ctx      context.Context
	basePath string
	logger   zerolog.Logger
}

func NewClient() *Client {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	basePath := "docker/buildContext"
	executionLog, err := os.Create(filepath.Join(basePath, "execution.log"))
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
	return &Client{client: cli, ctx: context.Background(), basePath: basePath, logger: logger}
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
	execResp, err := cli.inspectExecResp(createResp.ID)
	if err != nil {
		return ExecResult{}, err
	}
	cli.logger.Info().
		Str("containerID", containerID).
		Msg(execResp.StdOut)
	return execResp, nil
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

func (cli *Client) WriteZoneFile(containerID, srcPath string) {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		panic(err)
	}
	dstFile, err := os.Create(filepath.Join(cli.basePath, "nameserver", containerID, "zones", "active.zone"))
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		panic(err)
	}
	if err != nil {
		err = fmt.Errorf("could not write zone file: %w", err)
		panic(err)
	}
}

func (cli *Client) ReadLog(containerID, containerType, fileName string) []byte {
	logFilePath := filepath.Join(cli.basePath, containerType, containerID, "logs", fileName)
	content, err := os.ReadFile(logFilePath)
	if err != nil {
		panic(err)
	}
	var cleanedByteSlice []uint8
	for _, byte := range content {
		if byte != 0 {
			cleanedByteSlice = append(cleanedByteSlice, byte)
		}
	}
	return cleanedByteSlice
}

func (cli *Client) FlushLog(containerID, containerType, fileName string) {
	logFilePath := filepath.Join(cli.basePath, containerType, containerID, "logs", fileName)
	_, err := os.Create(logFilePath)
	if err != nil {
		panic(err)
	}
}
