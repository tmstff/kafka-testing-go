package kafkatesting

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/testcontainers/testcontainers-go"
	"os"
	"strings"
)

// deprecated
func StartKafkaContainer(ctx context.Context) (kafkaUrl string, terminateFunction func(context.Context) error, err error) {
	return StartKafkaWithEnv(ctx, map[string]string{})
}

func StartKafka(ctx context.Context) (kafkaUrl string, terminateFunction func(context.Context) error, err error) {
	return StartKafkaWithEnv(ctx, map[string]string{})
}

func StartKafkaWithEnv(ctx context.Context, env map[string]string) (kafkaUrl string, terminateFunction func(context.Context) error, err error) {
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

	logrus.Infof("containerHost: %s, kafkaPortExpositionHost: %s ", containerHost, kafkaPortExpositionHost)

	hostPort := fmt.Sprintf("%d", hostPortInt)
	kafkaExternalAdvertisedListener := containerHost + ":" + hostPort
	logrus.Debugf("kafkaExternalAdvertisedListener: %s", kafkaExternalAdvertisedListener)

	// hostPort:containerPort/protocol
	portExpositionString := kafkaPortExpositionHost + ":" + hostPort + ":9092"
	logrus.Debugf("kafka port exposition: %s", portExpositionString)

	identifier := strings.ToLower(uuid.New().String())

	composeFilePaths := []string{"docker-compose.yml"}

	compose := testcontainers.NewLocalDockerCompose(composeFilePaths, identifier)

	envWithDefaults := map[string]string{
		"KAFKA_PORT_EXPOSITION":              portExpositionString,
		"KAFKA_EXTERNAL_ADVERTISED_LISTENER": kafkaExternalAdvertisedListener,
	}

	for k,v := range env {
		envWithDefaults[k] = v
	}

	execError := compose.
		WithCommand([]string{"up", "-d"}).
		WithEnv(envWithDefaults).
		Invoke()
	if execError.Error != nil {
		logrus.Errorf("%s failed with\nstdout:\n%v\nstderr:\n%v", execError.Command, execError.Stdout, execError.Stderr)
		return "", nil, fmt.Errorf("could not run compose file: %v - %w", composeFilePaths, execError.Error)
	}

	logrus.Infof("kafka started at %s", kafkaExternalAdvertisedListener)

	return kafkaExternalAdvertisedListener, func(ctx context.Context) error {
		execError := compose.Down()
		if execError.Error != nil {
			logrus.Errorf("%s failed with\nstdout:\n%v\nstderr:\n%v", execError.Command, execError.Stdout, execError.Stderr)
			return fmt.Errorf("'docker-compose down' failed for compose file: %v - %w", composeFilePaths, err)
		}
		return nil
	}, nil
}