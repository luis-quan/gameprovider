package gameprovider

import (
	"github.com/luis-quan/cellnet"
	"github.com/luis-quan/cellnet/timer"
)

// IServerLogic 服务器逻辑类
type IServerLogic interface {
	OnInit()
	OnTimer(loop *timer.Loop)
	OnNetMessage(ses cellnet.Session, id int, msg interface{})
	OnConnectSuccess(session cellnet.Session)
	OnConnectAccept(session cellnet.Session)
	OnConnectClosed(session cellnet.Session)
	OnUserOffline(ses *SesContext)
	OnUserGoBack(ses *SesContext)
}
