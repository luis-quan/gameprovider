package gameprovider

import (
	"github.com/luis-quan/cellnet/timer"
)

// IServerLogic 服务器逻辑类
type IServerLogic interface {
	OnInit()
	OnTimer(loop *timer.Loop)
	OnNetMessage(sescontext *SesContext, id int, msg interface{})
	OnConnectSuccess(sescontext *SesContext)
	OnConnectAccept(sescontext *SesContext)
	OnConnectClosed(sescontext *SesContext)
	OnUserOffline(sescontext *SesContext)
	OnUserGoBack(sescontext *SesContext)
}
