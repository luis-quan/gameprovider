package gameprovider

import (
	"container/list"
	"os"
)

type sescontextmgr struct {
	mapSesContext map[int64]*SesContext
	listFree      list.List
}

func (s *sescontextmgr) init() {
	s.mapSesContext = make(map[int64]*SesContext)
}

func (s *sescontextmgr) getFree() *SesContext {
	log.Debugf("free sesContext len:%d", s.listFree.Len())

	var ses *SesContext
	if s.listFree.Len() == 0 {
		log.Debugln("create new sesContext empty!")
		ses = new(SesContext)
	} else {
		element := s.listFree.Back()
		value, ok := element.Value.(*SesContext)

		if ok {
			log.Debugln("create by pool sesContext!")
			ses = value
		} else {
			log.Debugln("create new sesContext full!")
			ses = new(SesContext)
		}
		//不管符合不符合的都删掉
		s.listFree.Remove(element)
	}

	return ses
}

func (s *sescontextmgr) addFree(ses *SesContext) {
	log.Debugf("add to free pool sesContext len:%d!", s.listFree.Len())
	s.listFree.PushBack(ses)
}

func (s *sescontextmgr) addContext(ses *SesContext) {
	log.Debugf("add to use pool sesContext len:%d!", s.listFree.Len())
	if id, err := ses.ID(); err == nil {
		s.mapSesContext[id] = ses
	} else {
		os.Exit(100)
	}
}

func (s *sescontextmgr) removeContext(ses *SesContext) {
	log.Debugf("remove from user pool sesContext len:%d!", s.listFree.Len())
	if id, err := ses.ID(); err == nil {
		delete(s.mapSesContext, id)
		s.addFree(ses)
		ses.reset()
	} else {
		os.Exit(101)
	}
}

func (s *sescontextmgr) findContext(sesID int64) (*SesContext, bool) {
	ses, bf := s.mapSesContext[sesID]
	return ses, bf
}
