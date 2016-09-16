package com

import (
	"errors"
	"fmt"
	"testing"

	"bytes"

	"github.com/arteev/td/driver"
	"golang.org/x/text/encoding/charmap"
)

type mockSerialer struct {
	CloseFn      func() error
	CloseInvoked bool

	WriteFn      func([]byte) (int, error)
	WriteInvoked bool

	ReadFn      func(b []byte) (n int, err error)
	ReadInvoked bool
}

func (m *mockSerialer) Write(b []byte) (n int, err error) {
	m.WriteInvoked = true
	if m.WriteFn != nil {
		return m.WriteFn(b)
	}
	return 0, nil
}
func (m *mockSerialer) Close() error {
	m.CloseInvoked = true
	if m.CloseFn != nil {
		return m.CloseFn()
	}
	return nil
}

func (m *mockSerialer) Read(b []byte) (n int, err error) {
	m.ReadInvoked = true
	if m.ReadFn != nil {
		return m.ReadFn(b)
	}
	return 0, nil
}

type mockProtocol struct {
	InitCmdFn      func() []byte
	InitCmdInvoked bool

	ClearCmdFn      func() []byte
	ClearCmdInvoked bool

	TestCmdFn      func() []byte
	TestCmdInvoked bool

	ClearRowCmdFn      func() []byte
	ClearRowCmdInvoked bool

	CursorVisibleCmdFn      func(bool) []byte
	CursorVisibleCmdInvoked bool

	ModeRewriteCmdFn      func() []byte
	ModeRewriteCmdInvoked bool

	ModeVScrollCmdFn      func() []byte
	ModeVScrollCmdInvoked bool

	ModeHScrollCmdFn      func() []byte
	ModeHScrollCmdInvoked bool

	BrightnessCmdFn      func(byte) []byte
	BrightnessCmdInvoked bool

	PrintRowCmdFn      func(byte, string) []byte
	PrintRowCmdInvoked bool

	CursorMoveUpCmdFn      func() []byte
	CursorMoveUpCmdInvoked bool

	CursorMoveDownCmdFn      func() []byte
	CursorMoveDownCmdInvoked bool

	CursorMoveRightCmdFn      func() []byte
	CursorMoveRightCmdInvoked bool

	CursorMoveLeftCmdFn      func() []byte
	CursorMoveLeftCmdInvoked bool

	CursorMoveLeftTopCmdFn      func() []byte
	CursorMoveLeftTopCmdInvoked bool

	CursorMoveBeginInRowCmdFn      func() []byte
	CursorMoveBeginInRowCmdInvoked bool

	CursorMoveEndInRowCmdFn      func() []byte
	CursorMoveEndInRowCmdInvoked bool

	CursorMoveBottomCmdFn      func() []byte
	CursorMoveBottomCmdInvoked bool

	CursorMoveCmdFn      func(row, col byte) []byte
	CursorMoveCmdInvoked bool

	FlagEnableCmdFn      func(enabled bool, num byte) []byte
	FlagEnableCmdInvoked bool

	FlagsDisableCmdFn      func() []byte
	FlagsDisableCmdInvoked bool
}

func (m *mockProtocol) InitCmd() []byte {
	m.InitCmdInvoked = true
	return m.InitCmdFn()
}
func (m *mockProtocol) ClearCmd() []byte {
	m.ClearCmdInvoked = true
	return m.ClearCmdFn()
}
func (m *mockProtocol) TestCmd() []byte {
	m.TestCmdInvoked = true
	return m.TestCmdFn()
}
func (m *mockProtocol) ClearRowCmd() []byte {
	m.ClearRowCmdInvoked = true
	return m.ClearRowCmdFn()
}
func (m *mockProtocol) CursorVisibleCmd(visible bool) []byte {
	m.CursorVisibleCmdInvoked = true
	return m.CursorVisibleCmdFn(visible)
}
func (m *mockProtocol) ModeRewriteCmd() []byte {
	m.ModeRewriteCmdInvoked = true
	return m.ModeRewriteCmdFn()
}
func (m *mockProtocol) ModeVScrollCmd() []byte {
	m.ModeVScrollCmdInvoked = true
	return m.ModeVScrollCmdFn()
}
func (m *mockProtocol) ModeHScrollCmd() []byte {
	m.ModeHScrollCmdInvoked = true
	return m.ModeHScrollCmdFn()
}

func (m *mockProtocol) BrightnessCmd(value byte) []byte {
	m.BrightnessCmdInvoked = true
	return m.BrightnessCmdFn(value)
}

func (m *mockProtocol) PrintRowCmd(row byte, text string) []byte {
	m.PrintRowCmdInvoked = true
	return m.PrintRowCmdFn(row, text)
}

func (m *mockProtocol) CursorMoveUpCmd() []byte {
	m.CursorMoveUpCmdInvoked = true
	return m.CursorMoveUpCmdFn()
}
func (m *mockProtocol) CursorMoveDownCmd() []byte {
	m.CursorMoveDownCmdInvoked = true
	return m.CursorMoveDownCmdFn()
}
func (m *mockProtocol) CursorMoveRightCmd() []byte {
	m.CursorMoveRightCmdInvoked = true
	return m.CursorMoveRightCmdFn()
}
func (m *mockProtocol) CursorMoveLeftCmd() []byte {
	m.CursorMoveLeftCmdInvoked = true
	return m.CursorMoveLeftCmdFn()
}
func (m *mockProtocol) CursorMoveLeftTopCmd() []byte {
	m.CursorMoveLeftTopCmdInvoked = true
	return m.CursorMoveLeftTopCmdFn()
}
func (m *mockProtocol) CursorMoveBeginInRowCmd() []byte {
	m.CursorMoveBeginInRowCmdInvoked = true
	return m.CursorMoveBeginInRowCmdFn()
}
func (m *mockProtocol) CursorMoveEndInRowCmd() []byte {
	m.CursorMoveEndInRowCmdInvoked = true
	return m.CursorMoveEndInRowCmdFn()
}
func (m *mockProtocol) CursorMoveBottomCmd() []byte {
	m.CursorMoveBottomCmdInvoked = true
	return m.CursorMoveBottomCmdFn()
}
func (m *mockProtocol) CursorMoveCmd(row, col byte) []byte {
	m.CursorMoveCmdInvoked = true
	return m.CursorMoveCmdFn(row, col)
}

func (m *mockProtocol) FlagEnableCmd(enabled bool, num byte) []byte {
	m.FlagEnableCmdInvoked = true
	return m.FlagEnableCmdFn(enabled, num)
}
func (m *mockProtocol) FlagsDisableCmd() []byte {
	m.FlagsDisableCmdInvoked = true
	return m.FlagsDisableCmdFn()
}

func TestCreateAndClose(t *testing.T) {
	mprot := &mockProtocol{}
	mser := &mockSerialer{}
	s := MustSerial(mprot)

	notinit := "The device is not initialized"
	if err := s.check(); err == nil || err.Error() != notinit {
		t.Fatalf("Excepted:The device is not initialized. Got: %v", err)
	}
	s.CreatePort(mser)
	if err := s.check(); err != nil {
		t.Fatalf("Excepted:check()=nil,got:%q", err)
	}
	if got := s.Protocol(); got != mprot {
		t.Fatalf("Excepted protocol %q, got %q", mprot, got)
	}
	if got := s.Serialer(); got != mser {
		t.Fatalf("Excepted Serialer %q, got %q", mser, got)
	}

	mser.CloseFn = func() error {
		return nil
	}
	err := s.Close()
	if err != nil {
		t.Fatal(err)
	}
	if !mser.CloseInvoked {
		t.Fatal("Excepted Close() to be invoked")
	}
	if s.opened {
		t.Error("Serial must be are not opened")
	}

	if err := s.Close(); err == nil || err.Error() != notinit {
		t.Fatalf("Excepted:The device is not initialized. Got: %v", err)
	}
}

func TestUnsupportProtocol(t *testing.T) {
	mprot := &mockProtocol{}
	mser := &mockSerialer{}
	s := &Serial{proto: mprot}
	s.CreatePort(mser)
	mprot.InitCmdFn = func() []byte {
		return nil
	}
	if err := s.Init(); err == nil || err != driver.ErrNotSupported {
		t.Errorf("Excepted %q\n", driver.ErrNotSupported)

	}
}

func TestEncoding(t *testing.T) {
	s := &Serial{proto: nil}
	casestr := "Россия"
	casestr1251 := []byte{0xd0, 0xee, 0xf1, 0xf1, 0xe8, 0xff}
	if got, err := s.encodetext(casestr); err != nil || got != casestr {
		if err != nil {
			t.Fatal(err)
		}
		t.Errorf("Excepted encoding %q, go %q", casestr, got)
	}
	s.SetEncoding(charmap.Windows1251)

	if got, err := s.encodetext(casestr); err != nil || bytes.Compare([]byte(got), casestr1251) != 0 {
		if err != nil {
			t.Fatal(err)
		}
		t.Errorf("Excepted encoding %q, go %q", casestr1251, got)
	}
}

func TestSendSerialer(t *testing.T) {
	mprot := &mockProtocol{}
	mser := &mockSerialer{}
	s := Serial{proto: mprot}

	mprot.InitCmdFn = func() []byte {
		return []byte{0x0}
	}
	notinit := "The device is not initialized"
	if err := s.Init(); err == nil || err.Error() != notinit {
		t.Fatalf("Excepted error: %q, got: %q", notinit, err)
	}
	s.CreatePort(mser)

	mser.WriteFn = func(b []byte) (int, error) {
		return len(b), nil
	}
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}

	if !mser.WriteInvoked {
		t.Fatal("Excepted serialer Write() must be invoked")
	}

	errtest := errors.New("fake")
	mser.WriteFn = func(b []byte) (int, error) {
		return 0, errtest
	}
	if err := s.Init(); err != errtest {
		t.Errorf("Excepted %q, got %q", errtest, err)
	}

	mser.WriteFn = func(b []byte) (int, error) {
		return 0, nil
	}
	msg := fmt.Sprintf("Error write: must %d byte(s), but %d byte(s)", 0, 1)
	if err := s.Init(); err == nil || err.Error() != msg {
		t.Errorf("Excepted %q, got %q", msg, err)
	}

}
func TestSendData(t *testing.T) {

	mprot := &mockProtocol{}
	mser := &mockSerialer{}
	s := Serial{proto: mprot}
	s.CreatePort(mser)
	data := []byte{1, 2, 3}
	var retdata []byte
	mser.WriteFn = func(b []byte) (int, error) {
		retdata = b
		return len(b), nil
	}
	err := s.Send(data)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(data, retdata) != 0 {
		t.Errorf("Excepted: %v, got %v\n", data, retdata)
	}

}
func TestReceive(t *testing.T) {
	mprot := &mockProtocol{}
	mser := &mockSerialer{}
	s := Serial{proto: mprot}
	s.CreatePort(mser)
	retdata := []byte{1, 2, 3}
	pos := 0
	mser.ReadFn = func(b []byte) (n int, err error) {
		b[0] = retdata[pos]
		pos++
		return 1, nil
	}
	buf := make([]byte, 3)
	n, err := s.Receive(buf)
	if err != nil {
		t.Fatal(err)
	}
	if n != len(retdata) || bytes.Compare(retdata, buf) != 0 {
		t.Errorf("Excepted: receive %d bytes (%v), but %d(%v)\n", len(retdata), retdata, n, buf)
	}

	mser.ReadFn = func(b []byte) (n int, err error) {
		return 0, fmt.Errorf("fake error")
	}
	if _, err := s.Receive(buf); err == nil || err.Error() != "fake error" {
		t.Errorf("Excepted error,got %q", err)
	}

}

func TestInvokedProtocol(t *testing.T) {
	mprot := &mockProtocol{}
	mser := &mockSerialer{}
	s := Serial{proto: mprot}

	mser.WriteFn = func(b []byte) (int, error) {
		return len(b), nil
	}

	commonFn := func() []byte {
		return []byte{0x0}
	}
	mprot.InitCmdFn = commonFn
	mprot.ClearCmdFn = commonFn
	mprot.TestCmdFn = commonFn
	mprot.ClearRowCmdFn = commonFn
	mprot.CursorVisibleCmdFn = func(visible bool) []byte {
		return []byte{0x0}
	}
	mprot.ModeRewriteCmdFn = commonFn
	mprot.ModeHScrollCmdFn = commonFn
	mprot.ModeVScrollCmdFn = commonFn
	mprot.BrightnessCmdFn = func(byte) []byte {
		return []byte{0}
	}
	mprot.PrintRowCmdFn = func(num byte, text string) []byte {
		//TODO : Check text
		return []byte{0x0}
	}

	mprot.CursorMoveUpCmdFn = commonFn
	mprot.CursorMoveDownCmdFn = commonFn
	mprot.CursorMoveRightCmdFn = commonFn
	mprot.CursorMoveLeftCmdFn = commonFn
	mprot.CursorMoveLeftTopCmdFn = commonFn
	mprot.CursorMoveBeginInRowCmdFn = commonFn
	mprot.CursorMoveEndInRowCmdFn = commonFn
	mprot.CursorMoveBottomCmdFn = commonFn
	mprot.CursorMoveCmdFn = func(row, col byte) []byte {
		return []byte{0x0}
	}
	mprot.FlagEnableCmdFn = func(bool, byte) []byte {
		return []byte{0x0}
	}
	mprot.FlagsDisableCmdFn = commonFn

	cases := []struct {
		Name    string
		Command func() error
		//CommandInvoked func() []byte
		Invoked *bool
	}{
		{
			Name:    "InitCmd",
			Command: s.Init,
			Invoked: &mprot.InitCmdInvoked,
		},
		{
			Name:    "ClearCmd",
			Command: s.Clear,
			Invoked: &mprot.ClearCmdInvoked,
		},
		{
			Name:    "TestCmd",
			Command: s.Test,
			Invoked: &mprot.TestCmdInvoked,
		},
		{
			Name:    "ClearRowCmd",
			Command: s.ClearRow,
			Invoked: &mprot.ClearRowCmdInvoked,
		},

		{
			Name:    "CursorVisibleCmd",
			Command: func() error { return s.CursorVisible(true) },
			Invoked: &mprot.CursorVisibleCmdInvoked,
		},
		{
			Name:    "ModeRewriteCmd",
			Command: s.ModeRewrite,
			Invoked: &mprot.ModeRewriteCmdInvoked,
		},
		{
			Name:    "ModeVScroll",
			Command: s.ModeVScroll,
			Invoked: &mprot.ModeVScrollCmdInvoked,
		},
		{
			Name:    "ModeHScroll",
			Command: s.ModeHScroll,
			Invoked: &mprot.ModeHScrollCmdInvoked,
		},
		{
			Name:    "BrightnessCmd",
			Command: func() error { return s.Brightness(1) },
			Invoked: &mprot.BrightnessCmdInvoked,
		},
		{
			Name:    "PrintRow",
			Command: func() error { return s.PrintRow(1, "test") },
			Invoked: &mprot.PrintRowCmdInvoked,
		},

		{
			Name:    "CursorMoveUpCmd",
			Command: s.CursorMoveUp,
			Invoked: &mprot.CursorMoveUpCmdInvoked,
		},
		{
			Name:    "CursorMoveDownCmd",
			Command: s.CursorMoveDown,
			Invoked: &mprot.CursorMoveDownCmdInvoked,
		},
		{
			Name:    "CursorMoveRightCmd",
			Command: s.CursorMoveRight,
			Invoked: &mprot.CursorMoveRightCmdInvoked,
		},
		{
			Name:    "CursorMoveLeftCmd",
			Command: s.CursorMoveLeft,
			Invoked: &mprot.CursorMoveLeftCmdInvoked,
		},
		{
			Name:    "CursorMoveLeftTopCmd",
			Command: s.CursorMoveLeftTop,
			Invoked: &mprot.CursorMoveLeftTopCmdInvoked,
		},
		{
			Name:    "CursorMoveBeginInRowCmd",
			Command: s.CursorMoveBeginInRow,
			Invoked: &mprot.CursorMoveBeginInRowCmdInvoked,
		},
		{
			Name:    "CursorMoveEndInRowCmd",
			Command: s.CursorMoveEndInRow,
			Invoked: &mprot.CursorMoveEndInRowCmdInvoked,
		},
		{
			Name:    "CursorMoveBottomCmd",
			Command: s.CursorMoveBottom,
			Invoked: &mprot.CursorMoveBottomCmdInvoked,
		},
		{
			Name:    "CursorMoveCmd",
			Command: func() error { return s.CursorMove(1, 1) },
			Invoked: &mprot.CursorMoveCmdInvoked,
		},
		{
			Name:    "FlagEnableCmd",
			Command: func() error { return s.FlagEnable(true, 0) },
			Invoked: &mprot.FlagEnableCmdInvoked,
		},
		{
			Name:    "FlagsDisableCmd",
			Command: s.FlagsDisable,
			Invoked: &mprot.FlagsDisableCmdInvoked,
		},
	}

	s.CreatePort(mser)

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			err := c.Command()
			if err != nil {
				t.Errorf("%q has error: %q\n", c.Name, err)
			}
			if !*c.Invoked {
				t.Errorf("Excepted %s must be invoked\n", c.Name)
			}
		})

	}
}
