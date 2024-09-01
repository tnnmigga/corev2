package mgdb

func Use(name string) *conn {
	return conns[name]
}

func Default() *conn {
	return conns["default"]
}
