package firich

import (
	"bytes"
)

type FirichProtocol struct {
}

//return []byte{}

func (p FirichProtocol) InitCmd() []byte {
	return []byte{0x1b, 0x40}
}

func (p FirichProtocol) ClearCmd() []byte {
	return []byte{0x0c}
}
func (p FirichProtocol) TestCmd() []byte {
	return []byte{0x1f, 0x40}
}

func (p FirichProtocol) ClearRowCmd() []byte {
	return []byte{0x18}
}

func (p FirichProtocol) CursorVisibleCmd(visible bool) []byte {
	var vbyte byte
	if visible {
		vbyte = 1
	} else {
		vbyte = 0
	}
	return []byte{0x1b, 0x5F, vbyte}
}

func (p FirichProtocol) ModeRewriteCmd() []byte {
	return []byte{0x1b, 0x11}
}
func (p FirichProtocol) ModeVScrollCmd() []byte {
	return []byte{0x1b, 0x12}
}
func (p FirichProtocol) ModeHScrollCmd() []byte {
	return []byte{0x1b, 0x13}
}

func (p FirichProtocol) BrightnessCmd(value byte) []byte {
	return []byte{0x1b, 0x2a, value}
}

func (p FirichProtocol) PrintRowCmd(row byte, text string) []byte {
	var buf bytes.Buffer
	switch row {
	case 1:
		buf.Write([]byte{0x1b, 0x51, 0x41})
		buf.Write([]byte(text))
		buf.WriteByte(0x0d)
		return buf.Bytes()

	case 2:
		buf.Write([]byte{0x1b, 0x51, 0x42})
		buf.Write([]byte(text))
		buf.WriteByte(0x0d)
		return buf.Bytes()

	}
	return nil
}

func (p FirichProtocol) CursorMoveUpCmd() []byte {
	return []byte{0x1b, 0x5b, 0x41}
}
func (p FirichProtocol) CursorMoveDownCmd() []byte {
	return []byte{0x1b, 0x5b, 0x42}
}
func (p FirichProtocol) CursorMoveRightCmd() []byte {
	return []byte{0x1b, 0x5b, 0x43}
}
func (p FirichProtocol) CursorMoveLeftCmd() []byte {
	return []byte{0x1b, 0x5b, 0x44}
}
func (p FirichProtocol) CursorMoveLeftTopCmd() []byte {
	return []byte{0x1b, 0x5b, 0x48}
}
func (p FirichProtocol) CursorMoveBeginInRowCmd() []byte {
	return []byte{0x1b, 0x5b, 0x4c}
}
func (p FirichProtocol) CursorMoveEndInRowCmd() []byte {
	return []byte{0x1b, 0x5b, 0x52}
}
func (p FirichProtocol) CursorMoveBottomCmd() []byte {
	return []byte{0x1b, 0x5b, 0x4B}
}
func (p FirichProtocol) CursorMoveCmd(row, col byte) []byte {
	return []byte{0x1b, 0x6c, col, row}
}
func (p FirichProtocol) FlagEnableCmd(enabled bool, num byte) []byte {
	var vbyte byte
	if enabled {
		vbyte = 1
	} else {
		vbyte = 0
	}
	return []byte{0x1f, 0x23, vbyte, num}
}

func (p FirichProtocol) FlagsDisableCmd() []byte {
	return []byte{0x1b, 0x7a}
}
