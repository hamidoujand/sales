package docker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os/exec"
	"time"
)

type Container struct {
	Name     string
	HostPort string
}

func StartContainer(image string, name string, port string, dockerArgs []string, imageArgs []string) (Container, error) {

	for i := range 2 {
		c, err := startContainer(image, name, port, dockerArgs, imageArgs)
		if err == nil {
			return c, nil
		}

		time.Sleep(time.Duration(i) * 100 * time.Millisecond)
	}

	//one last try
	return startContainer(image, name, port, dockerArgs, imageArgs)
}

func startContainer(image, name, port string, dockerArgs []string, imageArgs []string) (Container, error) {
	//check to see if container exists we return it
	if c, err := exists(name, port); err == nil {
		return c, nil
	}

	args := []string{"run", "-P", "-d", "--name", name}
	args = append(args, dockerArgs...)
	args = append(args, image)
	args = append(args, imageArgs...)

	var out bytes.Buffer
	cmd := exec.Command("docker", args...)
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return Container{}, fmt.Errorf("could not start the container %s: %w", name, err)
	}

	id := out.String()[:12]

	hostIP, hostPort, err := extractHostIP(id, port)
	if err != nil {
		//stop container if we failed to get ip/port
		StopContainer(id)
		return Container{}, fmt.Errorf("extractHostIP: %w", err)
	}

	c := Container{
		Name:     name,
		HostPort: net.JoinHostPort(hostIP, hostPort),
	}

	return c, nil
}

func exists(name string, port string) (Container, error) {
	hostIP, hostPort, err := extractHostIP(name, port)
	if err != nil {
		return Container{}, fmt.Errorf("container %s is not running", name)
	}

	return Container{
		Name:     name,
		HostPort: net.JoinHostPort(hostIP, hostPort),
	}, nil
}

func StopContainer(name string) error {
	if err := exec.Command("docker", "stop", name).Run(); err != nil {
		return fmt.Errorf("stopping container %s failed: %w", name, err)
	}

	if err := exec.Command("docker", "rm", name, "-v").Run(); err != nil {
		return fmt.Errorf("removing container %s failed: %w", name, err)
	}

	return nil
}

func DumpContainerLogs(name string) []byte {
	out, err := exec.Command("docker", "logs", name).CombinedOutput()
	if err != nil {
		return []byte(err.Error())
	}
	return out
}

//==============================================================================
// HostIP/HostPort

type ContainerInspect struct {
	NetworkSettings struct {
		Ports map[string][]PortBinding `json:"Ports"`
	} `json:"NetworkSettings"`
}

type PortBinding struct {
	HostIP   string `json:"HostIp"`
	HostPort string `json:"HostPort"`
}

func extractHostIP(containerName string, containerPort string) (host string, port string, err error) {
	var out bytes.Buffer

	cmd := exec.Command("docker", "inspect", containerName)
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", "", fmt.Errorf("failed to inspect the container[%s]: %w", containerName, err)
	}

	var inspectData []ContainerInspect
	if err := json.Unmarshal(out.Bytes(), &inspectData); err != nil {
		return "", "", fmt.Errorf("failed to unmarshal into inspect data slice: %w", err)
	}

	portKey := containerPort + "/tcp"
	ports := inspectData[0].NetworkSettings.Ports[portKey]

	for _, binding := range ports {
		//not looking for IPV6
		if binding.HostIP != "::" {
			if binding.HostIP == "" {
				//localhost
				return "localhost", binding.HostPort, nil
			}

			return binding.HostIP, binding.HostPort, nil
		}
	}

	//not-found
	return "", "", fmt.Errorf("could not locate ip/port for container %s", containerName)
}
