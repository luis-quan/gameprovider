package main

import (
	"fmt"
	"reflect"

	"github.com/luis-quan/gameprovider"

	"github.com/luis-quan/cellnet"
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
	s.serverPeer = gameprovider.Ggameprovider.CreateServer("server", "127.0.0.1:18803")
}

func (s *Serverlogic) OnTimer(loop *timer.Loop) {

}
func (s *Serverlogic) OnNetMessage(sescontext *gameprovider.SesContext, id int, msg interface{}) {
	if sescontext.Peer() == s.serverPeer {
		fmt.Println("OnNetMessage")
	}
}
func (s *Serverlogic) OnConnectSuccess(context *gameprovider.SesContext) {

}
func (s *Serverlogic) OnConnectAccept(context *gameprovider.SesContext) {

}
func (s *Serverlogic) OnConnectClosed(context *gameprovider.SesContext) {

}
func (s *Serverlogic) OnUserOffline(context *gameprovider.SesContext) {

}

func (s *Serverlogic) OnUserGoBack(context *gameprovider.SesContext) {

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
