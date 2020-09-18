# kafka-testing-go
Library enabling easy testing of your applications interacting with kafka.

Uses testcontainers-go to start a docker container running zookeeper and a kafka broker and provides easy means to start, stop and interact with it.

## Usage

    ctx := context.Background()
    
    kafkaURL, terminateKafka, err := kafkatesting.StartKafka(ctx)

    if err != nil {
        // handle error
    }
    
    // use kafkaURL to access kafka
    ...
    
    // stop containers after tests are done
    terminateKafka(ctx)
