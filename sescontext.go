package gameprovider

import (
	"errors"

	"github.com/luis-quan/cellnet"
)

type UserContextInterface interface {
	Init()
	Reset()
}

//SesContext 保存节点数据...
type SesContext struct {
	id          int64
	ses         cellnet.Session
	bCanRelease bool
	userContext UserContextInterface
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

func (s *SesContext) setCanRelease(b bool) {
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
	s.bCanRelease = true
	s.id = 0
	s.userContext.Reset()
}

//Initialize 初始化数据
func (s *SesContext) initialize() {
	s.bCanRelease = true
	s.userContext.Init()
}

func (s *SesContext) setSession(ses cellnet.Session) {
	s.ses = ses
	if ses != nil {
		s.id = ses.ID()
	}
}
