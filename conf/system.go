package conf

var (
	ServerID uint32
	Groups   []string
)

func mustInit() {
	ServerID = Num[uint32]("must.serverID")
	Groups = List[string]("must.groups")
}
