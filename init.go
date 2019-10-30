package fluentdpub

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/fluent/fluent-logger-golang/fluent"
	"gocloud.dev/pubsub"
)

// Scheme is the URL scheme fluentdpub registers its URLOpeners under on pubsub.DefaultMux.
const Scheme = "fluentd"

func init() {
	o := new(defaultDialer)
	pubsub.DefaultURLMux().RegisterTopic(Scheme, o)
}

type defaultDialer struct {
	opener *URLOpener
	err    error
}

func (d defaultDialer) OpenTopicURL(ctx context.Context, u *url.URL) (*pubsub.Topic, error) {
	c, tag, err := parseEnvVar(os.Getenv("FLUENTD_UPSTREAM_URL"))
	if err != nil {
		return nil, err
	}
	conn, err := fluent.New(*c)
	if err != nil {
		return nil, err
	}
	o := URLOpener{
		Connection: conn,
		TagPrefix:  tag,
	}
	return o.OpenTopicURL(ctx, u)
}

// URLOpener opens Fluentd URLs like "fluentd://myapp.tag".
//
// Host part is used as a tag prefix
//
// No query parameters are supported.
type URLOpener struct {
	// Connection to use for communication with the server.
	Connection *fluent.Fluent
	// TagPrefix is prefix of tags of all topic of this URLOpener
	TagPrefix string
}

// OpenTopicURL opens a pubsub.Topic based on u.
func (o *URLOpener) OpenTopicURL(ctx context.Context, u *url.URL) (*pubsub.Topic, error) {
	var opt TopicOptions
	key := u.Query().Get("bodykey")
	if key != "" {
		opt.BodyKey = key
		u.Query().Del("bodykey")
	}
	tagkey := u.Query().Get("tagkey")
	if tagkey != "" {
		opt.TagKey = tagkey
		u.Query().Del("tagkey")
	}
	for param := range u.Query() {
		return nil, fmt.Errorf("open topic %v: invalid query parameter %s", u, param)
	}
	var subject string
	if o.TagPrefix != "" && u.Hostname() != "" {
		subject = o.TagPrefix + "." + u.Hostname()
	} else {
		subject = o.TagPrefix + u.Hostname()
	}
	return OpenTopic(o.Connection, subject, opt)
}

func parseEnvVar(env string) (*fluent.Config, string, error) {
	if strings.HasPrefix(env, "://") {
		env = "tcp" + env
	}
	u, err := url.Parse(env)
	if err != nil {
		return nil, "", err
	}
	result := &fluent.Config{
		FluentPort:    24224,
		FluentHost:    "127.0.0.1",
		FluentNetwork: "tcp",
	}
	switch u.Scheme {
	case "":
		fallthrough
	case "tcp":
		break
	case "udp":
		result.FluentNetwork = u.Scheme
	default:
		return nil, "", fmt.Errorf("Unknown scheme %q. Only(", u.Scheme)
	}
	if u.Port() != "" {
		p, err := strconv.ParseInt(u.Port(), 10, 64)
		if err != nil {
			return nil, "", err
		}
		result.FluentPort = int(p)
	}
	if u.Host != "" {
		result.FluentHost = u.Hostname()
	}
	var tagPrefix string
	if u.Path != "" {
		tagPrefix = strings.TrimLeft(u.Path, "/")
	}
	return result, tagPrefix, nil
}
