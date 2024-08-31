package natsmq

func Use(name string) *natsConn {
	return conns[name]
}

func Default() *natsConn {
	return conns["default"]
}
