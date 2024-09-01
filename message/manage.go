package message

import "github.com/tnnmigga/corev2/log"

func Start() error {
	err := createOrUpdateStream()
	if err != nil {
		return err
	}
	err = subscribeMsg()
	if err != nil {
		return err
	}
	return nil
}

func Stop() {
	streamConsCtx.Drain()
	for _, sub := range msgSubscriptions {
		err := sub.Drain()
		if err != nil {
			log.Error(err)
		}
	}
}
