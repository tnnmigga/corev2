package conf

var (
	ServerID uint32
)

func mustInit() {
	ServerID = Num[uint32]("must.sid")
}
