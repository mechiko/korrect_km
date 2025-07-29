package spaserver

import (
	"context"
	"korrectkm/sse"
	"strings"
)

type MessageType int

const (
	FlushError MessageType = iota
	FlushInfo
)

func (s MessageType) String() string {
	switch s {
	case FlushError:
		return "error"
	case FlushInfo:
		return "info"
	}
	return "unknown"
}

// Type - error message
type FlushMsg struct {
	Msg  string
	Type string
}

// msg сообщение
// msgType тип сообщения error info spaserver\templates\index\flush.html
func (s *Server) SetFlush(msg string, msgType string) {
	s.flushMu.Lock()
	defer s.flushMu.Unlock()

	msg = strings.ReplaceAll(msg, "\n", " ")
	msg = strings.ReplaceAll(msg, "\r", " ")
	fl := &FlushMsg{
		Msg:  msg,
		Type: msgType,
	}
	s.flush = fl
	// switch fl.Type {
	// case "error":
	// 	s.streamError.Eventlog.Clear()
	// case "info":
	// 	s.streamInfo.Eventlog.Clear()
	// }
	s.PublishFlush(fl)
}

func (s *Server) GetFlush() *FlushMsg {
	s.flushMu.Lock()
	defer s.flushMu.Unlock()
	return s.flush
}

func (s *Server) PublishFlush(msg *FlushMsg) {
	ev := &sse.Event{
		Data: []byte(msg.Msg),
	}
	s.sseManager.Publish(msg.Type, ev)
}

func (s *Server) Publish(ctx context.Context) {
	s.flushMu.Lock()
	defer s.flushMu.Unlock()

	if s.flush != nil {
		s.PublishFlush(s.flush)
	}
}
