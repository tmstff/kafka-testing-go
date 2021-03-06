package kafkatesting

const DockerComposeFile =
	`version: '3'

networks:
  kafka-net:
    driver: bridge

services:
  zookeeper:
    image: confluentinc/cp-zookeeper:5.4.0
    networks:
      - kafka-net
    environment:
      ZOOKEEPER_SERVER_ID: 1
      ZOOKEEPER_CLIENT_PORT: 2181
      ZOOKEEPER_SERVERS: "zookeeper:2888:3888"

  kafka:
    image: confluentinc/cp-kafka:5.4.0
    ports:
      - ${KAFKA_PORT_EXPOSITION}
    networks:
      - kafka-net
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_AUTO_CREATE_TOPICS_ENABLE: ${KAFKA_AUTO_CREATE_TOPICS_ENABLE}
      KAFKA_DELETE_TOPIC_ENABLE: 'true'
      KAFKA_AUTO_LEADER_REBALANCE_ENABLE: 'true'
      KAFKA_LISTENERS: LISTENER_DOCKER_INTERNAL://kafka:29092,LISTENER_DOCKER_EXTERNAL://0.0.0.0:9092
      KAFKA_ADVERTISED_LISTENERS: LISTENER_DOCKER_INTERNAL://kafka:29092,LISTENER_DOCKER_EXTERNAL://${KAFKA_EXTERNAL_ADVERTISED_LISTENER}
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: LISTENER_DOCKER_INTERNAL:PLAINTEXT,LISTENER_DOCKER_EXTERNAL:PLAINTEXT
      KAFKA_INTER_BROKER_LISTENER_NAME: LISTENER_DOCKER_INTERNAL
      KAFKA_ZOOKEEPER_CONNECT: "zookeeper:2181"
      KAFKA_LOG4J_LOGGERS: "kafka.controller=INFO,kafka.producer.async.DefaultEventHandler=INFO,state.change.logger=INFO"
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
      KAFKA_NUM_PARTITIONS: 3
`
