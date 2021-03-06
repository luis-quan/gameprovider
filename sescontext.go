package gameprovider

import (
	"errors"

	"github.com/luis-quan/cellnet"
	"github.com/luis-quan/cellnet/serial/binaryserial"
)

type NodeContextInterface interface {
	OnCreate()
	Init()
	Reset()
}

//SesContext 保存节点数据...
type SesContext struct {
	id           int64
	ses          cellnet.Session
	gameProvider *gameprovider
	bCanRelease  bool
	nodeContext  NodeContextInterface
}

func (s *SesContext) NodeContext() NodeContextInterface {
	return s.nodeContext
}

func (s *SesContext) SendRawData(msg interface{}, id int) {
	b, err := binaryserial.BinaryWrite(msg, 4)

	if err == nil {
		data := new(cellnet.RawPacket)
		data.MsgData = b
		data.MsgID = id
		if s.ses != nil {
			s.ses.Send(data)
		}
	} else {
		log.Errorln(msg)
	}
}

func (s *SesContext) SendData(msg interface{}) {
	if s.ses != nil {
		s.ses.Send(msg)
	}
}

func (s *SesContext) Peer() cellnet.Peer {
	if s.ses != nil {
		return s.ses.Peer()
	}

	return nil
}

func (s *SesContext) Close() {
	s.SetCanRelease(true)
	if s.ses != nil {
		s.ses.Close()
	}
}

func (s *SesContext) SetCanRelease(b bool) {
	s.bCanRelease = b
}

func (s *SesContext) ID() (int64, error) {
	if s.ses != nil {
		return s.ses.ID(), nil
	}

	return 0, errors.New("ses is nil")
}

//ResetNode 重置节点数据
func (s *SesContext) reset() {
	s.ses = nil
	s.gameProvider = nil
	s.bCanRelease = true
	s.id = 0
	s.nodeContext.Reset()
}

//Initialize 初始化数据
func (s *SesContext) initialize() {
	s.bCanRelease = true
	s.nodeContext.Init()
}

func (s *SesContext) setSession(ses cellnet.Session, gameProvider *gameprovider) {
	s.ses = ses
	s.gameProvider = gameProvider
	if ses != nil {
		s.id = ses.ID()
	} else {
		s.id = 0
	}
}
