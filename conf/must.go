package conf

var (
	ServerID uint32
	Groups   []string
)

func mustInit() {
	ServerID = Uint32("must.serverID")
	Groups = List[string]("must.groups")
}
