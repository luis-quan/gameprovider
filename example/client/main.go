package main

import (
	"fmt"
	"gameprovider"
	"reflect"

	"github.com/luis-quan/cellnet"
	_ "github.com/luis-quan/cellnet/peer/gorillaws"
	_ "github.com/luis-quan/cellnet/proc/gorillaws"
	"github.com/luis-quan/cellnet/timer"
)

type playernode struct {
}

func (s *playernode) Init() {
	fmt.Println("playernode init")
}

func (s *playernode) Reset() {
	fmt.Println("playernode reset")
}

type Serverlogic struct {
	serverPeer cellnet.Peer
}

func (s *Serverlogic) OnInit() {
	s.serverPeer = gameprovider.Ggameprovider.ConnectServer("client", "127.0.0.1:18803")
}

func (s *Serverlogic) OnTimer(loop *timer.Loop) {

}
func (s *Serverlogic) OnNetMessage(p cellnet.Peer, id int, msg interface{}) {
	if p == s.serverPeer {
		fmt.Println("OnNetMessage")
	}
}
func (s *Serverlogic) OnConnectSuccess(session cellnet.Session) {

}
func (s *Serverlogic) OnConnectAccept(session cellnet.Session) {

}
func (s *Serverlogic) OnConnectClosed(session cellnet.Session) {

}
func (s *Serverlogic) OnUserOffline(player *gameprovider.SesContext) {

}

func (s *Serverlogic) OnUserGoBack(player *gameprovider.SesContext) {

}

type tableinfo struct {
}

func (s *tableinfo) Init() {
	fmt.Println("tableinfo init")
}

func (s *tableinfo) Reset() {
	fmt.Println("tableinfo reset")
}

func main() {
	gameprovider.Ggameprovider.RegisterSesContext(reflect.TypeOf((*playernode)(nil)).Elem())
	//gameprovider.Ggameprovider.RegisterTablelogicType(reflect.TypeOf(tableinfo{}))

	var server Serverlogic
	gameprovider.Ggameprovider.Initialize(&server)
	gameprovider.Ggameprovider.Start()
}
