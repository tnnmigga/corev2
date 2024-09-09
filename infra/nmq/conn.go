package nmq

import (
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/tnnmigga/corev2/conf"
	"github.com/tnnmigga/corev2/log"
)

func init() {
	if err := initFromConf(); err != nil {
		panic(err)
	}
}

var conns = map[string]*natsConn{}

type config struct {
	URL string
}

type natsConn struct {
	*nats.Conn
	stream jetstream.JetStream
}

func (conn *natsConn) Stream() jetstream.JetStream {
	return conn.stream
}

func initFromConf() error {
	data := conf.Map[config]("nmq", nil)
	for key, item := range data {
		conn, err := newConn(item)
		if err != nil {
			panic(err)
		}
		conns[key] = conn
	}
	return nil
}

func Use(name string) *natsConn {
	return conns[name]
}

func Default() *natsConn {
	return conns["default"]
}

func newConn(c config) (*natsConn, error) {
	conn, err := nats.Connect(
		c.URL,
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(10),
		nats.ReconnectWait(time.Second),
		nats.ReconnectHandler(func(_ *nats.Conn) {
			log.Errorf("nats retry connect")
		}),
	)
	if err != nil {
		return nil, err
	}
	js, err := jetstream.New(conn)
	if err != nil {
		return nil, err
	}

	return &natsConn{
		Conn:   conn,
		stream: js,
	}, nil
}
