package kafkatesting

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStartKafka(t *testing.T) {

	ctx := context.Background()

	kafkaUrl, terminateKafka, err := StartKafka(ctx)
	if err != nil {
		assert.Fail(t, err.Error())
	} else {
		defer terminateKafka(ctx)
	}

	assert.Regexp(t,"([0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3})\\:?([0-9]{1,5})?", kafkaUrl)
}