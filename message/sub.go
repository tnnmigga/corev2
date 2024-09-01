package message

import (
	"context"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/tnnmigga/corev2/conc"
	"github.com/tnnmigga/corev2/conf"
	"github.com/tnnmigga/corev2/infra/nmq"
	"github.com/tnnmigga/corev2/log"
	"github.com/tnnmigga/corev2/message/codec"
)

var (
	streamConsCtx    jetstream.ConsumeContext
	msgSubscriptions []*nats.Subscription
)

func createOrUpdateStream() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	js, err := nmq.Default().Stream().CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:         streamName(),
		Subjects:     []string{streamCastSubject(conf.ServerID)},
		MaxConsumers: 1,
		Retention:    jetstream.WorkQueuePolicy,
		MaxMsgs:      defaultMaxMsgs,
		Storage:      jetstream.MemoryStorage,
		MaxAge:       time.Hour * 24 * 3,
	})
	if err != nil {
		return err
	}
	cons, err := js.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:       consumerName(),
		FilterSubject: streamCastSubject(conf.ServerID),
		MaxAckPending: 1000,
		MaxWaiting:    512,
		AckWait:       defaultTimeout,
		MaxDeliver:    5,
	})
	if err != nil {
		return err
	}
	streamConsCtx, err = cons.Consume(consumeMsg)
	if err != nil {
		return err
	}
	return nil
}

func consumeMsg(msg jetstream.Msg) {
	b := msg.Data()
	err := msg.Ack()
	if err != nil {
		log.Error(err)
		return
	}
	data, err := codec.Decode(b)
	if err != nil {
		log.Error(err)
		return
	}
	Delivery(data)
}

func subscribeMsg() error {
	sub, err := nmq.Default().Subscribe(castSubject(conf.ServerID), handleCastMsg)
	if err != nil {
		return err
	}
	msgSubscriptions = append(msgSubscriptions, sub)
	sub, err = nmq.Default().Subscribe(rpcSubject(conf.ServerID), handleRPCMsg)
	if err != nil {
		return err
	}
	msgSubscriptions = append(msgSubscriptions, sub)
	for _, group := range conf.Groups {
		sub, err = nmq.Default().Subscribe(broadcastSubject(group), handleCastMsg)
		if err != nil {
			return err
		}
		msgSubscriptions = append(msgSubscriptions, sub)
		sub, err = nmq.Default().Subscribe(anycastSubject(group), handleCastMsg)
		if err != nil {
			return err
		}
		msgSubscriptions = append(msgSubscriptions, sub)
		sub, err = nmq.Default().Subscribe(anyRPCSubject(group), handleRPCMsg)
		if err != nil {
			return err
		}
		msgSubscriptions = append(msgSubscriptions, sub)
	}
	return nil
}

func handleCastMsg(msg *nats.Msg) {
	b := msg.Data
	data, err := codec.Decode(b)
	if err != nil {
		log.Error(err)
		return
	}
	Delivery(data)
}

func handleRPCMsg(msg *nats.Msg) {
	b := msg.Data
	req, err := codec.Decode(b)
	if err != nil {
		log.Error(err)
		return
	}
	conc.Go(func() {
		ctx := newRPCContext(req)
		data, err := ctx.exec()
		var (
			header nats.Header
			b      []byte
		)
		if err != nil {
			header.Add("err", err.Error())
		} else {
			b = codec.Encode(data)
		}
		err = msg.RespondMsg(&nats.Msg{
			Header: header,
			Data:   b,
		})
		if err != nil {
			log.Error(err)
		}
	})
}
