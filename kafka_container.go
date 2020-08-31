package kafkatesting

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/sirupsen/logrus"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"os"
)

func StartKafkaContainer(ctx context.Context) (kafkaUrl string, terminateFunction func(context.Context) error, err error) {

	var containerHost string
	var kafkaPortExpositionHost string
	if val, exists := os.LookupEnv("HOST_IP"); exists {
		containerHost = val
		kafkaPortExpositionHost = val
	} else {
		containerHost = GetOutboundIP()
		kafkaPortExpositionHost = "0.0.0.0"
	}

	hostPortInt, err := GetFreePort(kafkaPortExpositionHost)
	if err != nil {
		return "", nil, err
	}

	logrus.Infof("container host: %s, kafkaPortExpositionHost: %s ", containerHost, kafkaPortExpositionHost)

	hostPort := fmt.Sprintf("%d", hostPortInt)

	logrus.Infof("kafka host port: %s", hostPort)

	natContainerPort := "9092/tcp"

	// hostPort:containerPort/protocol
	portExposition := kafkaPortExpositionHost + ":" + hostPort + ":" + natContainerPort

	req := testcontainers.ContainerRequest{
		Image:        "spotify/kafka",
		ExposedPorts: []string{portExposition},
		WaitingFor:   wait.ForListeningPort(nat.Port(natContainerPort)),
		Env: map[string]string{
			"ADVERTISED_HOST": containerHost,
			"ADVERTISED_PORT": hostPort,
		},
	}

	logrus.Info("Starting up kafka container")
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return "", nil, err
	}

	address := fmt.Sprintf("%s:%s", containerHost, hostPort)

	logrus.Infof("kafka container started at %s", address)

	return address, container.Terminate, nil
}

