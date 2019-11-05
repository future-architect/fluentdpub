# fluentdpub - fluentd backend for gocloud.dev's pubsub

[![GoDoc](https://godoc.org/github.com/future-architect/fluentdpub?status.svg)](https://godoc.org/github.com/future-architect/fluentdpub)

`fluentdpub` is a backend for gocloud.dev. It only supports publisher (topic) side.

It supports two style API like gocloud.dev, common URL constructor and Fluentd specific constructor.

## Common URL Constructor

After importing `fluentdpub` package, you can open bia `pubsub.OpenTopic()` function.
`fluentd://` URL only contains tag prefix.

Upstream Fluentd server is specified by ``FLUENTD_UPSTREAM_URL`` env var.

```go
import (
	"context"

	"gocloud.dev/pubsub"
	_ "github.com/future-architect/fluentdpub"
)

// pubsub.OpenTopic creates a *pubsub.Topic from a URL.
// This URL will Dial the Fluentd server at the URL in the environment variable
// FLUENTD_UPSTREAM_URL and send messages with tag "example.tag".
topic, err := pubsub.OpenTopic(ctx, "fluentd://stg.my-app")
if err != nil {
	return err
}
err = topic.Send(ctx, &pubsub.Message{
    // Fluentd doesn't have main contant named body
    // fluentdpub embeds body into metadata.
    // Default key of body is "message"
	Body: []byte("Hello, World!\n"),
	// Metadata is optional and can be nil.
	Metadata: map[string]string{
		// These are examples of metadata.
		// There is nothing special about the key names.
		"language":   "en",
		"importance": "high",
	},
})

defer topic.Shutdown(ctx)
```

### Option

URL can have two queries:

- `bodykey` (default: `"message"`):
      Message.Body content is stored into `Metadata` of this key.
- `tagkey` (default: `"tag"`):
      The Value in this key of `Metadata` is used as tag prefix.

### Tag Name

There three locations you can specify tag.

- `FLUENTD_UPSTREAM_URL`'s path name
- URL's host name
- Metadata's tag value (specified by `tagkey` parameter or `"tag"` by default)

If `FLUENTD_UPSTREAM_URL` is `"tcp://localhost:24224/prod"` URL is `fluentd://my-app` and the message to send contains `"tag"` Metadata with the value `"error"`, log will be sent with `prod.my-app.error` tag. Empty element will be ignored (if `FLUENTD_UPSTREAM_URL` is `"tcp://localhost:24224"`, final tag will be `my-app.error`).

## Fluentd Specific Constructor

The `fluentdpub.OpenTopic` constructor opens a Fluentd subject as a topic. You must first create an *fluent.Fluent to your Fluentd server instance.

```go
import (
	"context"

	"github.com/future-architect/fluentdpub"
)

// pubsub.OpenTopic creates a *pubsub.Topic from a URL.
// This URL will Dial the Fluentd server at the URL in the environment variable
// FLUENTD_UPSTREAM_URL and send messages with tag "example.tag".
f, err := fluent.New(fluent.Config{
	FluentHost: "fluentd.example.com",
})
topic, err := fluentdpub.OpenTopic(f, "prod.myapp", fluentdpub.TopicOptions{})
if err != nil {
	return err
}
err = topic.Send(context.Background(), &pubsub.Message{
    // Fluentd doesn't have main contant named body
    // fluentdpub embeds body into metadata.
    // Default key of body is "message"
	Body: []byte("Hello, World!\n"),
	// Metadata is optional and can be nil.
	Metadata: map[string]string{
		// These are examples of metadata.
		// There is nothing special about the key names.
		"language":   "en",
		"importance": "high",
	},
})
```

### Option

TopicOptions has two fields:

- `BodyKey` (default: `"message"`):
      Message.Body content is stored into `Metadata` of this key.
- `TagKey` (default: `"tag"`):
      The Value in this key of `Metadata` is used as tag prefix.

### Tag Name

There two location you can specify tag.

- URL's host name
- Metadata's tag value (specified by `tagkey` parameter or `"tag"` by default)

If URL is `fluentd://prod.my-app` and the message to send contains `"tag"` Metadata with the value `"error"`, log will be sent with `prod.my-app.error` tag. If one of them is empty, non empty one is used.

## License

Apache 2
