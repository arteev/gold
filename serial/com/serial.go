package com

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/arteev/gold/driver"

	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
)

type Serial struct {
	mu       sync.Mutex
	opened   bool
	encoding encoding.Encoding

	port  Serialer
	proto driver.Protocol
}
type Serialer interface {
	Write(b []byte) (n int, err error)
	Close() error
	Read(b []byte) (n int, err error)
}

func MustSerial(protocol driver.Protocol) *Serial {
	return &Serial{
		proto: protocol,
	}
}

func (s *Serial) check() error {
	if !s.opened {
		return errors.New("The device is not initialized")
	}
	return nil
}

func (s *Serial) send(data []byte) error {
	if err := s.check(); err != nil {
		return err
	}
	n, err := s.port.Write(data)
	if err != nil {
		return err
	}
	if len(data) != n {
		return fmt.Errorf("Error write: must %d byte(s), but %d byte(s)", n, len(data))
	}
	return nil
}
func (s *Serial) sendFromProtocol(fn func() []byte) error {
	data := fn()
	if len(data) == 0 {
		return driver.ErrNotSupported
	}
	return s.send(data)
}

func (s *Serial) encodetext(text string) (string, error) {
	if s.encoding == nil {
		return text, nil
	}
	sr := strings.NewReader(text)
	tr := transform.NewReader(sr, s.encoding.NewEncoder())
	buf, err := ioutil.ReadAll(tr)
	if err != err {
		return "", err
	}
	return string(buf), nil
}

/////
func (s *Serial) CreatePort(port Serialer) {
	s.port = port
	s.opened = true
}

/////

func (s *Serial) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.check(); err != nil {
		return err
	}
	err := s.port.Close()
	if err == nil {
		s.opened = false
	}
	return err
}
func (s *Serial) SetEncoding(encoding encoding.Encoding) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.encoding = encoding
}
func (s *Serial) Init() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.sendFromProtocol(s.proto.InitCmd)
}
func (s *Serial) Test() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.sendFromProtocol(s.proto.TestCmd)
}
func (s *Serial) Clear() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.sendFromProtocol(s.proto.ClearCmd)
}
func (s *Serial) ClearRow() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.sendFromProtocol(s.proto.ClearRowCmd)
}

func (s *Serial) CursorVisible(visible bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	fn := func() []byte {
		return s.proto.CursorVisibleCmd(visible)
	}
	return s.sendFromProtocol(fn)
}
func (s *Serial) ModeRewrite() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.sendFromProtocol(s.proto.ModeRewriteCmd)
}
func (s *Serial) ModeVScroll() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.sendFromProtocol(s.proto.ModeVScrollCmd)
}
func (s *Serial) ModeHScroll() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.sendFromProtocol(s.proto.ModeHScrollCmd)
}
func (s *Serial) Brightness(value byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	fn := func() []byte {
		return s.proto.BrightnessCmd(value)
	}
	return s.sendFromProtocol(fn)
}

func (s *Serial) PrintRow(row byte, text string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	outtext, err := s.encodetext(text)
	if err != nil {
		return err
	}
	fn := func() []byte {
		return s.proto.PrintRowCmd(row, outtext)
	}
	return s.sendFromProtocol(fn)
}

func (s *Serial) CursorMoveUp() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.sendFromProtocol(s.proto.CursorMoveUpCmd)
}
func (s *Serial) CursorMoveDown() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.sendFromProtocol(s.proto.CursorMoveDownCmd)
}
func (s *Serial) CursorMoveRight() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.sendFromProtocol(s.proto.CursorMoveRightCmd)
}
func (s *Serial) CursorMoveLeft() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.sendFromProtocol(s.proto.CursorMoveLeftCmd)
}

func (s *Serial) CursorMoveLeftTop() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.sendFromProtocol(s.proto.CursorMoveLeftTopCmd)
}
func (s *Serial) CursorMoveBeginInRow() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.sendFromProtocol(s.proto.CursorMoveBeginInRowCmd)
}
func (s *Serial) CursorMoveEndInRow() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.sendFromProtocol(s.proto.CursorMoveEndInRowCmd)
}
func (s *Serial) CursorMoveBottom() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.sendFromProtocol(s.proto.CursorMoveBottomCmd)
}
func (s *Serial) CursorMove(row, col byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	fn := func() []byte {
		return s.proto.CursorMoveCmd(row, col)
	}
	return s.sendFromProtocol(fn)
}

func (s *Serial) FlagEnable(enabled bool, num byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	fn := func() []byte {
		return s.proto.FlagEnableCmd(enabled, num)
	}
	return s.sendFromProtocol(fn)
}
func (s *Serial) FlagsDisable() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.sendFromProtocol(s.proto.FlagsDisableCmd)
}

func (s *Serial) Send(data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.send(data)
}

func (s *Serial) Receive(b []byte) (n int, err error) {
	n, err = s.port.Read(b)
	if err != nil {
		return 0, err
	}
	if n < len(b) {

		nn, err := s.Receive(b[n:])
		return n + nn, err
	}
	return n, nil
}

func (s *Serial) Protocol() driver.Protocol {
	return s.proto
}

func (s *Serial) Serialer() Serialer {
	return s.port
}
