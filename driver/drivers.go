package driver

import "errors"
import "golang.org/x/text/encoding"

//Errors
var (
	ErrNotSupported = errors.New("Command is not supported")
)

type Driver interface {
	GetDisplay(proto Protocol, config map[string]interface{}) (Display, error)
}

type Display interface {
	//System
	Close() error
	SetEncoding(encoding.Encoding)

	//
	Init() error
	Test() error
	Clear() error

	Send([]byte) error
	Receive(b []byte) (n int, err error)

	//Mode
	ModeRewrite() error
	ModeVScroll() error
	ModeHScroll() error

	Brightness(byte) error

	//Rows
	ClearRow() error

	//Cursor
	CursorVisible(bool) error
	CursorMoveUp() error
	CursorMoveDown() error
	CursorMoveRight() error
	CursorMoveLeft() error

	CursorMoveLeftTop() error
	CursorMoveBeginInRow() error
	CursorMoveEndInRow() error
	CursorMoveBottom() error
	CursorMove(row, col byte) error

	//Text
	PrintRow(row byte, text string) error

	//Flags
	FlagEnable(enabled bool, num byte) error
	FlagsDisable() error
}

// Protocol specific to a particular communication protocol
type Protocol interface {
	InitCmd() []byte
	ClearCmd() []byte
	TestCmd() []byte

	//Mode
	ModeRewriteCmd() []byte
	ModeVScrollCmd() []byte
	ModeHScrollCmd() []byte
	BrightnessCmd(byte) []byte

	//Rows
	ClearRowCmd() []byte

	//Cursor
	CursorVisibleCmd(bool) []byte

	CursorMoveUpCmd() []byte
	CursorMoveDownCmd() []byte
	CursorMoveRightCmd() []byte
	CursorMoveLeftCmd() []byte

	CursorMoveLeftTopCmd() []byte
	CursorMoveBeginInRowCmd() []byte
	CursorMoveEndInRowCmd() []byte
	CursorMoveBottomCmd() []byte
	CursorMoveCmd(row, col byte) []byte

	//Text
	PrintRowCmd(row byte, text string) []byte

	//Flags
	FlagEnableCmd(enabled bool, num byte) []byte
	FlagsDisableCmd() []byte
}
