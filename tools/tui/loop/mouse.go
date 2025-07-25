// License: GPLv3 Copyright: 2023, Kovid Goyal, <kovid at kovidgoyal.net>

package loop

import (
	"fmt"
	"strconv"
	"strings"
)

var _ = fmt.Print

type MouseEventType uint
type MouseButtonFlag uint

const (
	MOUSE_PRESS MouseEventType = iota
	MOUSE_RELEASE
	MOUSE_MOVE
	MOUSE_CLICK
	MOUSE_LEAVE
)

func (e MouseEventType) String() string {
	switch e {
	case MOUSE_PRESS:
		return "press"
	case MOUSE_RELEASE:
		return "release"
	case MOUSE_MOVE:
		return "move"
	case MOUSE_CLICK:
		return "click"
	}
	return strconv.Itoa(int(e))
}

type PointerShape uint8

const (
	// start pointer shape enum (auto generated by gen-key-constants.py do not edit)
	DEFAULT_POINTER       PointerShape = 0
	TEXT_POINTER          PointerShape = 1
	POINTER_POINTER       PointerShape = 2
	HELP_POINTER          PointerShape = 3
	WAIT_POINTER          PointerShape = 4
	PROGRESS_POINTER      PointerShape = 5
	CROSSHAIR_POINTER     PointerShape = 6
	CELL_POINTER          PointerShape = 7
	VERTICAL_TEXT_POINTER PointerShape = 8
	MOVE_POINTER          PointerShape = 9
	E_RESIZE_POINTER      PointerShape = 10
	NE_RESIZE_POINTER     PointerShape = 11
	NW_RESIZE_POINTER     PointerShape = 12
	N_RESIZE_POINTER      PointerShape = 13
	SE_RESIZE_POINTER     PointerShape = 14
	SW_RESIZE_POINTER     PointerShape = 15
	S_RESIZE_POINTER      PointerShape = 16
	W_RESIZE_POINTER      PointerShape = 17
	EW_RESIZE_POINTER     PointerShape = 18
	NS_RESIZE_POINTER     PointerShape = 19
	NESW_RESIZE_POINTER   PointerShape = 20
	NWSE_RESIZE_POINTER   PointerShape = 21
	ZOOM_IN_POINTER       PointerShape = 22
	ZOOM_OUT_POINTER      PointerShape = 23
	ALIAS_POINTER         PointerShape = 24
	COPY_POINTER          PointerShape = 25
	NOT_ALLOWED_POINTER   PointerShape = 26
	NO_DROP_POINTER       PointerShape = 27
	GRAB_POINTER          PointerShape = 28
	GRABBING_POINTER      PointerShape = 29

// end pointer shape enum
)

func (e PointerShape) String() string {
	switch e {
	// start pointer shape tostring (auto generated by gen-key-constants.py do not edit)
	case DEFAULT_POINTER:
		return "default"
	case TEXT_POINTER:
		return "text"
	case POINTER_POINTER:
		return "pointer"
	case HELP_POINTER:
		return "help"
	case WAIT_POINTER:
		return "wait"
	case PROGRESS_POINTER:
		return "progress"
	case CROSSHAIR_POINTER:
		return "crosshair"
	case CELL_POINTER:
		return "cell"
	case VERTICAL_TEXT_POINTER:
		return "vertical-text"
	case MOVE_POINTER:
		return "move"
	case E_RESIZE_POINTER:
		return "e-resize"
	case NE_RESIZE_POINTER:
		return "ne-resize"
	case NW_RESIZE_POINTER:
		return "nw-resize"
	case N_RESIZE_POINTER:
		return "n-resize"
	case SE_RESIZE_POINTER:
		return "se-resize"
	case SW_RESIZE_POINTER:
		return "sw-resize"
	case S_RESIZE_POINTER:
		return "s-resize"
	case W_RESIZE_POINTER:
		return "w-resize"
	case EW_RESIZE_POINTER:
		return "ew-resize"
	case NS_RESIZE_POINTER:
		return "ns-resize"
	case NESW_RESIZE_POINTER:
		return "nesw-resize"
	case NWSE_RESIZE_POINTER:
		return "nwse-resize"
	case ZOOM_IN_POINTER:
		return "zoom-in"
	case ZOOM_OUT_POINTER:
		return "zoom-out"
	case ALIAS_POINTER:
		return "alias"
	case COPY_POINTER:
		return "copy"
	case NOT_ALLOWED_POINTER:
		return "not-allowed"
	case NO_DROP_POINTER:
		return "no-drop"
	case GRAB_POINTER:
		return "grab"
	case GRABBING_POINTER:
		return "grabbing"
		// end pointer shape tostring
	}
	return strconv.Itoa(int(e))
}

const (
	SHIFT_INDICATOR         int = 1 << 2
	ALT_INDICATOR               = 1 << 3
	CTRL_INDICATOR              = 1 << 4
	MOTION_INDICATOR            = 1 << 5
	SCROLL_BUTTON_INDICATOR     = 1 << 6
	EXTRA_BUTTON_INDICATOR      = 1 << 7
	LEAVE_INDICATOR             = 1 << 8
)

const (
	NO_MOUSE_BUTTON   MouseButtonFlag = 0
	LEFT_MOUSE_BUTTON MouseButtonFlag = 1 << iota
	MIDDLE_MOUSE_BUTTON
	RIGHT_MOUSE_BUTTON
	FOURTH_MOUSE_BUTTON
	FIFTH_MOUSE_BUTTON
	SIXTH_MOUSE_BUTTON
	SEVENTH_MOUSE_BUTTON
	MOUSE_WHEEL_UP
	MOUSE_WHEEL_DOWN
	MOUSE_WHEEL_LEFT
	MOUSE_WHEEL_RIGHT
)

var bmap = [...]MouseButtonFlag{LEFT_MOUSE_BUTTON, MIDDLE_MOUSE_BUTTON, RIGHT_MOUSE_BUTTON}
var ebmap = [...]MouseButtonFlag{FOURTH_MOUSE_BUTTON, FIFTH_MOUSE_BUTTON, SIXTH_MOUSE_BUTTON, SEVENTH_MOUSE_BUTTON}
var wbmap = [...]MouseButtonFlag{MOUSE_WHEEL_UP, MOUSE_WHEEL_DOWN, MOUSE_WHEEL_LEFT, MOUSE_WHEEL_RIGHT}

func (b MouseButtonFlag) String() string {
	ans := ""
	switch {
	case b&LEFT_MOUSE_BUTTON != 0:
		ans += "|LEFT"
	case b&MIDDLE_MOUSE_BUTTON != 0:
		ans += "|MIDDLE"
	case b&RIGHT_MOUSE_BUTTON != 0:
		ans += "|RIGHT"
	case b&FOURTH_MOUSE_BUTTON != 0:
		ans += "|FOURTH"
	case b&FIFTH_MOUSE_BUTTON != 0:
		ans += "|FIFTH"
	case b&SIXTH_MOUSE_BUTTON != 0:
		ans += "|SIXTH"
	case b&SEVENTH_MOUSE_BUTTON != 0:
		ans += "|SEVENTH"
	case b&MOUSE_WHEEL_UP != 0:
		ans += "|WHEEL_UP"
	case b&MOUSE_WHEEL_DOWN != 0:
		ans += "|WHEEL_DOWN"
	case b&MOUSE_WHEEL_LEFT != 0:
		ans += "|WHEEL_LEFT"
	case b&MOUSE_WHEEL_RIGHT != 0:
		ans += "|WHEEL_RIGHT"
	}
	ans = strings.TrimLeft(ans, "|")
	if ans == "" {
		ans = "NONE"
	}
	return ans
}

type MouseEvent struct {
	Event_type  MouseEventType
	Buttons     MouseButtonFlag
	Mods        KeyModifiers
	Cell, Pixel struct{ X, Y int }
}

func (e MouseEvent) String() string {
	return fmt.Sprintf("MouseEvent{%s %s %s Cell:%v Pixel:%v}", e.Event_type, e.Buttons, e.Mods, e.Cell, e.Pixel)
}

func pixel_to_cell(px, length, cell_length int) int {
	px = max(0, min(px, length-1))
	if cell_length > 0 {
		return px / cell_length
	}
	return 0
}

func decode_sgr_mouse(text string, screen_size ScreenSize, last_letter byte) *MouseEvent {
	text = text[:len(text)-1]
	parts := strings.Split(text, ";")
	if len(parts) != 3 {
		return nil
	}
	cb, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil
	}
	ans := MouseEvent{}
	ans.Pixel.X, err = strconv.Atoi(parts[1])
	if err != nil {
		return nil
	}
	if len(parts[2]) < 1 {
		return nil
	}
	if ans.Pixel.Y, err = strconv.Atoi(parts[2]); err != nil {
		return nil
	}
	if last_letter == 'm' {
		ans.Event_type = MOUSE_RELEASE
	} else if cb&MOTION_INDICATOR != 0 {
		ans.Event_type = MOUSE_MOVE
	}
	cb3 := cb & 3
	switch {
	case cb&LEAVE_INDICATOR != 0:
		ans.Event_type = MOUSE_LEAVE
	case cb&EXTRA_BUTTON_INDICATOR != 0:
		ans.Buttons |= ebmap[cb3]
	case cb&SCROLL_BUTTON_INDICATOR != 0:
		ans.Buttons |= wbmap[cb3]
	case cb3 < 3:
		ans.Buttons |= bmap[cb3]
	}
	if cb&SHIFT_INDICATOR != 0 {
		ans.Mods |= SHIFT
	}
	if cb&ALT_INDICATOR != 0 {
		ans.Mods |= ALT
	}
	if cb&CTRL_INDICATOR != 0 {
		ans.Mods |= CTRL
	}
	ans.Cell.X = pixel_to_cell(ans.Pixel.X, int(screen_size.WidthPx), int(screen_size.CellWidth))
	ans.Cell.Y = pixel_to_cell(ans.Pixel.Y, int(screen_size.HeightPx), int(screen_size.CellHeight))

	return &ans
}

func MouseEventFromCSI(csi string, screen_size ScreenSize) *MouseEvent {
	if len(csi) == 0 {
		return nil
	}
	last_char := csi[len(csi)-1]
	if last_char != 'm' && last_char != 'M' {
		return nil
	}
	switch csi[0] {
	case '<':
		return decode_sgr_mouse(csi[1:], screen_size, last_char)
	default:
		return nil
	}
}
