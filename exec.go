package whaler

import (
	"bufio"
	"context"
	"errors"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func RunCommand(containerID string, commands ...string) (*bufio.Reader, error) {
	client, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	if _, err := client.ContainerInspect(ctx, containerID); err != nil {
		return nil, err
	}
	response, err := client.ContainerExecCreate(ctx, containerID, types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Detach:       false,
		Tty:          false,
		Cmd:          commands,
	})
	if err != nil {
		return nil, err
	}

	execID := response.ID
	if execID == "" {
		return nil, errors.New("exec ID empty")
	}

	resp, err := client.ContainerExecAttach(ctx, execID, types.ExecConfig{
		AttachStdout: true,
		AttachStdin:  false,
		AttachStderr: true,
	})
	if err != nil {
		return nil, err
	}
	return resp.Reader, nil
}
