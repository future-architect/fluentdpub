// +build fluentd

package fluentdpub

import (
	"context"
	"github.com/fluent/fluent-logger-golang/fluent"
	"github.com/stretchr/testify/assert"
	"gocloud.dev/pubsub"
	"os"
	"testing"
)

//  TestOpenTopic needs local fluentd to test. To run this test, run test with "-tags fluend" flag
//
//  Run it via the following command:
//
//   docker run --rm -v ${PWD}/tmp/fluentd:/fluentd/log --name fluentd -p 24224:24224 fluent/fluentd:latest
func TestOpenTopic(t *testing.T) {
	// tcp://localhost:24224
	f, err := fluent.New(fluent.Config{})
	assert.NoError(t, err)
	topic, err := OpenTopic(f, "test.tag", TopicOptions{})
	assert.NoError(t, err)
	assert.NotNil(t, topic)
	err = topic.Send(context.Background(), &pubsub.Message{
		Body: []byte("Hello, World!\n"),
		// Metadata is optional and can be nil.
		Metadata: map[string]string{
			// These are examples of metadata.
			// There is nothing special about the key names.
			"language":   "en",
			"importance": "high",
		},
	})
	assert.NoError(t, err)
}

func TestOpenTopicByURL(t *testing.T) {
	os.Setenv("FLUENTD_UPSTREAM_URL", "tcp://localhost:24224/first")
	topic, err := pubsub.OpenTopic(context.Background(), "fluentd://second")
	assert.NoError(t, err)
	assert.NotNil(t, topic)
	err = topic.Send(context.Background(), &pubsub.Message{
		Body: []byte("Hello, World!\n"),
		// Metadata is optional and can be nil.
		Metadata: map[string]string{
			// These are examples of metadata.
			// There is nothing special about the key names.
			"tag":        "third",
			"language":   "en",
			"importance": "high",
		},
	})
	assert.NoError(t, err)
}
