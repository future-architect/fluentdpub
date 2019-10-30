package fluentdpub

import (
	"context"
	"errors"
	"fmt"

	"github.com/fluent/fluent-logger-golang/fluent"
	"gocloud.dev/gcerrors"
	"gocloud.dev/pubsub"
	"gocloud.dev/pubsub/driver"
)

// TopicOptions sets options for constructing a *pubsub.Topic backed by Fluentd.
//
// There is no option now.
type TopicOptions struct {
	// BodyKey is a key name to store message body. Default is "message".
	BodyKey string
	// TagKey is a key name of additional tag. Default is "tag".
	TagKey string
}

type topic struct {
	f         *fluent.Fluent
	tagPrefix string
	bodyKey   string
	tagKey    string
}

// OpenTopic returns a *pubsub.Topic for use with Fluentd.
// The subject is the Fluentd's tag; for more info, see
// https://github.com/fluent/fluent-logger-golang
func OpenTopic(f *fluent.Fluent, tagPrefix string, opt TopicOptions) (*pubsub.Topic, error) {
	if f == nil {
		return nil, errors.New("fluentdpub: fluent.Fluent is required")
	}
	if opt.BodyKey == "" {
		opt.BodyKey = "message"
	}
	if opt.TagKey == "" {
		opt.TagKey = "tag"
	}
	return pubsub.NewTopic(&topic{
		f:         f,
		tagPrefix: tagPrefix,
		bodyKey:   opt.BodyKey,
		tagKey:    opt.TagKey,
	}, nil), nil
}

func (t topic) SendBatch(ctx context.Context, ms []*driver.Message) error {
	for _, msg := range ms {
		var fullTag string
		tag, ok := msg.Metadata[t.tagKey]
		if ok {
			delete(msg.Metadata, t.tagKey)
			if t.tagPrefix != "" {
				fullTag = t.tagPrefix + "." + tag
			} else {
				fullTag = tag
			}
		} else {
			fullTag = t.tagPrefix
		}
		if fullTag == "" {
			return fmt.Errorf("Message %v doesn't have tag", msg.AckID)
		}
		msg.Metadata[t.bodyKey] = string(msg.Body)
		err := t.f.Post(fullTag, msg.Metadata)
		if err != nil {
			return err
		}
	}
	return nil
}

// IsRetryable implements driver.Topic.IsRetryable.
func (t topic) IsRetryable(err error) bool {
	return false
}

// As implements driver.Topic.As.
func (t *topic) As(i interface{}) bool {
	f, ok := i.(**fluent.Fluent)
	if !ok {
		return false
	}
	*f = t.f
	return true
}

// ErrorAs implements driver.Topic.ErrorAs.
func (t topic) ErrorAs(error, interface{}) bool {
	return false
}

// ErrorCode implements driver.Topic.ErrorCode.
func (t topic) ErrorCode(err error) gcerrors.ErrorCode {
	switch err {
	case nil:
		return gcerrors.OK
	case context.Canceled:
		return gcerrors.Canceled
	}
	return gcerrors.Unknown
}

// Close implements driver.Topic.Close.
func (t topic) Close() error {
	return t.f.Close()
}
