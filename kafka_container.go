package kafkatesting

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/testcontainers/testcontainers-go"
	"io/ioutil"
	"os"
	"strings"
)

// deprecated
func StartKafkaContainer(ctx context.Context) (kafkaUrl string, terminateFunction func(context.Context), err error) {
	return StartKafkaWithEnv(ctx, map[string]string{})
}

func StartKafka(ctx context.Context) (kafkaUrl string, terminateFunction func(context.Context), err error) {
	return StartKafkaWithEnv(ctx, map[string]string{})
}

func StartKafkaWithEnv(ctx context.Context, env map[string]string) (kafkaUrl string, terminateFunction func(context.Context), err error) {
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

	composeFileName, deleteTempDockerComposeYml, err := createTmpDockerComposeYml()
	if err != nil {
		return "", nil, err
	}

	logrus.Infof("docker-compose.yml: %s", composeFileName)

	compose := testcontainers.NewLocalDockerCompose([]string{composeFileName}, identifier)

	envWithDefaults := map[string]string{
		"KAFKA_PORT_EXPOSITION":              portExpositionString,
		"KAFKA_EXTERNAL_ADVERTISED_LISTENER": kafkaExternalAdvertisedListener,
		"KAFKA_AUTO_CREATE_TOPICS_ENABLE": "true",
	}

	for k, v := range env {
		envWithDefaults[k] = v
	}

	execError := compose.
		WithCommand([]string{"up", "-d"}).
		WithEnv(envWithDefaults).
		Invoke()
	err = execError.Error
	if err != nil {
		return "", nil, fmt.Errorf("could not run compose file: %v - %w\nstdout:\n%v\nstderr:\n%v", []string{composeFileName}, err, execError.Stdout, execError.Stderr)
	}

	logrus.Infof("kafka started at %s", kafkaExternalAdvertisedListener)

	return kafkaExternalAdvertisedListener, func(ctx context.Context) {
		execError := compose.Down()
		err := execError.Error
		if err != nil {
			logrus.Errorf("'docker-compose down' failed for compose file: %v - %w\nstdout:\n%v\nstderr:\n%v", []string{composeFileName}, err, execError.Stdout, execError.Stderr)
		}
		deleteTempDockerComposeYml()
	}, nil
}

func createTmpDockerComposeYml() (composeFileName string, deleteTempFile func(), err error) {
	tmpDirName, err := ioutil.TempDir(os.TempDir(), "kafka-testing-go") // create temp dir to avoid .env confusion
	if err != nil {
		return "", nil, fmt.Errorf("cannot create temporary directory: %w", err)
	}
	tmpFile, err := ioutil.TempFile(tmpDirName, "docker-compose-*.yml")
	if err != nil {
		return "", nil, fmt.Errorf("cannot create temporary file: %w", err)
	}
	defer tmpFile.Close()

	composeFileName = tmpFile.Name()

	content := []byte(DockerComposeFile)
	_, err = tmpFile.Write(content)
	if err != nil {
		return "", nil, fmt.Errorf("cannot write to file '%s': %w", composeFileName, err)
	}

	deleteTempFile = func() {
		err := os.Remove(composeFileName)
		if err != nil {
			logrus.Warnf("could not remove temp file %s", composeFileName)
		}
		err = os.Remove(tmpDirName)
		if err != nil {
			logrus.Warnf("could not remove temp dir %s", tmpDirName)
		}
	}

	return composeFileName, deleteTempFile, nil
}