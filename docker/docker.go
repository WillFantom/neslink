package docker

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/willfantom/neslink"
)

// NPName returns a netns provider that provides the netns path for the docker
// container with the given name.
func NPName(cli *client.Client, containerName string) neslink.NsProvider {
	return neslink.NPGeneric(
		"docker-container-name",
		func() (neslink.Namespace, error) {
			containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
			if err != nil {
				return neslink.Namespace(""), fmt.Errorf("failed to get container list:  %w", err)
			}
			for _, c := range containers {
				for _, n := range c.Names {
					if strings.EqualFold(strings.TrimPrefix(n, "/"), containerName) {
						return NPID(cli, c.ID).Provide()
					}
				}
			}
			return neslink.Namespace(""), fmt.Errorf("could not find container with given name: %s", containerName)
		},
	)
}

// NPID returns a netns provider that provides the netns path for the docker
// container with the given container id.
func NPID(cli *client.Client, containerID string) neslink.NsProvider {
	return neslink.NPGeneric(
		"docker-container-id",
		func() (neslink.Namespace, error) {
			c, err := cli.ContainerInspect(context.Background(), containerID)
			if err != nil {
				return neslink.Namespace(""), fmt.Errorf("failed to get container infomation:  %w", err)
			}
			return neslink.NPProcess(c.State.Pid).Provide()
		},
	)
}
