package message

import (
	"fmt"
	"reflect"

	"github.com/tnnmigga/corev2/conf"
	"github.com/tnnmigga/corev2/iface"
)

var (
	local toLocal
)

func Cast(dest iface.IMsgDest, msg any) {

}

func RPC(dest iface.IMsgDest, req any) {

}

func ByServerID(serverID int) iface.IMsgDest {
	return byServerID{
		serverID: serverID,
	}
}

type byServerID struct {
	serverID int
}

func (b byServerID) String() string {
	if b.serverID == conf.ServerID() {
		return "local"
	}
	return fmt.Sprintf("server.id.%d", b.serverID)
}

func ByServerType(serverType string) iface.IMsgDest {
	return byServerType{
		serverType: serverType,
	}
}

type byServerType struct {
	serverType string
}

func (b byServerType) String() string {
	return fmt.Sprintf("server.type.%s", b.serverType)
}

func Local() iface.IMsgDest {
	return local
}

type toLocal struct{}

func (b toLocal) String() string {
	return "local"
}

func castLocal(msg any) {
	mType := reflect.TypeOf(msg)
	for _, m := range recvers[mType] {
		m.Assign(msg)
	}
}
