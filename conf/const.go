package conf

var (
	serverID   int
	serverType string
)

func ServerID() int {
	return serverID
}

func ServerType() string {
	return serverType
}
