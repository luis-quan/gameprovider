package gameprovider

import (
	"container/list"
	"os"
	"reflect"
	"time"

	"github.com/luis-quan/cellnet"
	"github.com/luis-quan/cellnet/peer"
	_ "github.com/luis-quan/cellnet/peer/gorillaws"
	"github.com/luis-quan/cellnet/proc"
	"github.com/luis-quan/cellnet/timer"
)

//网络接口供应商
type gameprovider struct {
	//逻辑管理
	serverLogic IServerLogic
	//消息队列 单线程
	queue cellnet.EventQueue
	//定时器
	timerLoop *timer.Loop
	//空闲节点
	sesContextmgr sescontextmgr
	//玩家数据
	userContextType reflect.Type
	//peer离岸边
	peerList list.List
}

//初始化
func (s *gameprovider) Initialize(serverLogic IServerLogic) {
	s.serverLogic = serverLogic
	s.queue = cellnet.NewEventQueue()
	//10毫秒一次 1秒100次
	s.timerLoop = timer.NewLoop(s.queue, time.Millisecond*10, s.onTimer, nil)
	s.sesContextmgr.init()
}

func (s *gameprovider) CreateServer(name string, addr string, procName string) cellnet.Peer {
	peerType := "gorillaws.Acceptor"
	//Peer初始话
	var p cellnet.GenericPeer
	p = peer.NewGenericPeer(peerType, name, addr, s.queue)
	//注册回调
	proc.BindProcessorHandler(p, procName, s.onEvent)
	//开始监听
	p.Start()

	return p
}

func (s *gameprovider) ConnectServer(name string, addr string, procName string) cellnet.Peer {
	peerType := "gorillaws.Connector"
	//Peer初始话
	var p cellnet.GenericPeer
	p = peer.NewGenericPeer(peerType, name, addr, s.queue)
	p.(cellnet.WSConnector).SetReconnectDuration(time.Second)
	//注册回调
	proc.BindProcessorHandler(p, procName, s.onEvent)
	//开始监听
	p.Start()

	return p
}

func (s *gameprovider) Start() {
	s.serverLogic.OnInit()

	s.timerLoop.Start()
	// 事件队列开始循环
	s.queue.StartLoop()
	// 阻塞等待事件队列结束退出( 在另外的goroutine调用queue.StopLoop() )
	s.queue.Wait()
}

//注册玩家数据
func (s *gameprovider) RegisterSesContext(userContextType reflect.Type) {
	if userContextType.Kind() == reflect.Ptr {
		userContextType = userContextType.Elem()
	}
	value := reflect.New(userContextType).Interface()
	_, ok := value.(UserContextInterface)
	if ok {
		s.userContextType = userContextType
	} else {
		os.Exit(10)
	}
}

//创建玩家信息
func (s *gameprovider) createSesContext() *SesContext {
	context := s.sesContextmgr.getFree()
	if context.userContext == nil {
		co := reflect.New(s.userContextType).Interface()
		if userContext, ok := co.(UserContextInterface); ok {
			userContext.OnCreate()
			context.userContext = userContext
		} else {
			os.Exit(111)
		}
	}

	context.initialize()

	return context
}

//接受到连接 创建连接
func (s *gameprovider) createConnection(session cellnet.Session) {
	context := s.createSesContext()
	context.setSession(session, s)
	session.(cellnet.ContextSet).SetContext(SES_CONTEXT, context)
	s.sesContextmgr.addUseContext(context)
}

//关闭的返回
func (s *gameprovider) onCloseConnection(session cellnet.Session) {
	if context, ok := session.(cellnet.ContextSet).GetContext(SES_CONTEXT); ok == true && context != nil {
		sesContext, ok := context.(*SesContext)
		if ok && sesContext != nil {
			if sesContext.bCanRelease {
				s.sesContextmgr.eraseUseContext(sesContext)
			}
			sesContext.setSession(nil, s)
		}
	}
}

func (s *gameprovider) releaseContext(sescontext *SesContext) {
	s.sesContextmgr.eraseUseContext(sescontext)
}

// SetContextCanRelease是不是能释放节点
func (s *gameprovider) SetContextCanRelease(sescontext *SesContext, b bool) {
	if sescontext.id != 0 {
		if ses, bf := s.sesContextmgr.findContext(sescontext.id); bf == true {
			//已经断开连接了
			if sescontext.ses == nil {
				//释放掉
				s.sesContextmgr.eraseUseContext(ses)
			} else {
				//连接还在
				sescontext.SetCanRelease(true)
			}
		} else {
			//查找出错了
		}
	} else {
		//这个节点 本来就不在使用不需要控制
	}
}

func (s *gameprovider) onEvent(ev cellnet.Event) {
	switch event := ev.(type) {
	case *cellnet.RecvMsgEvent:
		switch msg := event.Message().(type) {
		case *cellnet.SessionConnected:
			log.Debugln("client connected ", msg)
			s.createConnection(event.Session())
			if sescontext, bf := s.sesContextmgr.findContext(event.Session().ID()); bf == true {
				s.serverLogic.OnConnectSuccess(sescontext)
			}
			break
		case *cellnet.SessionAccepted:
			log.Debugln("client accepted ", msg)
			s.createConnection(event.Session())
			if sescontext, bf := s.sesContextmgr.findContext(event.Session().ID()); bf == true {
				s.serverLogic.OnConnectAccept(sescontext)
			}
			break
		case *cellnet.SessionClosed:
			log.Debugln("client error ")
			if sescontext, bf := s.sesContextmgr.findContext(event.Session().ID()); bf == true {
				s.serverLogic.OnConnectClosed(sescontext)
			}
			s.onCloseConnection(event.Session())
			break
		case *cellnet.SessionConnectError:
			log.Debugln("client SessionConnectError ")
			os.Exit(1)
			break
		default:
			if sescontext, bf := s.sesContextmgr.findContext(event.Session().ID()); bf == true {
				s.serverLogic.OnNetMessage(sescontext, event.ID(), event.Message())
			}
			break
		}
	}
}

//定时器
func (s *gameprovider) onTimer(loop *timer.Loop) {
	//log.Debugf("onTimer")

	if s.serverLogic != nil {
		s.serverLogic.OnTimer(loop)
	}
}

//消息队列
func (s *gameprovider) Queue() *cellnet.EventQueue {
	return &s.queue
}

// Ggameprovider 全局管理类
var Ggameprovider gameprovider

func init() {

}
