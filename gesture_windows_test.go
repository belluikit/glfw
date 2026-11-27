//go:build windows

package glfw_test

import (
	"os"
	"runtime"
	"syscall"
	"testing"
	"unsafe"

	glfw "github.com/belluikit/glfw"
)

var callQ = make(chan func())

// do runs fn on the thread that called glfw.Init, as GLFW requires, and waits.
func do(fn func()) {
	done := make(chan struct{})
	callQ <- func() { defer close(done); fn() }
	<-done
}

func TestMain(m *testing.M) {
	runtime.LockOSThread()
	var code int
	go func() {
		do(func() {
			if err := glfw.Init(); err != nil {
				panic("glfw.Init: " + err.Error())
			}
		})
		code = m.Run()
		do(glfw.Terminate)
		close(callQ)
	}()
	for fn := range callQ {
		fn()
	}
	os.Exit(code)
}

var sendMessage = syscall.NewLazyDLL("user32.dll").NewProc("SendMessageW")

const (
	wmMouseWheel  = 0x020A
	wmMouseHWheel = 0x020E
	mkShift       = 0x0004
	mkControl     = 0x0008
	wheelDelta    = 120
)

// A wheel message reaches the extended scroll callback with the control and
// shift it carries, which a key-tracking client never sees: a precision
// touchpad sends a pinch as WM_MOUSEWHEEL with MK_CONTROL set and no key
// pressed (belluikit/glfw#13). Nothing here presses a key.
func TestWin32ScrollExCarriesTheWheelMessageModifiers(t *testing.T) {
	var w *glfw.Window
	var err error
	do(func() {
		glfw.DefaultWindowHints()
		glfw.WindowHint(glfw.ClientAPI, glfw.NoAPI)
		glfw.WindowHint(glfw.Visible, glfw.False)
		w, err = glfw.CreateWindow(320, 240, t.Name(), nil, nil)
	})
	if err != nil {
		t.Fatalf("CreateWindow: %v", err)
	}
	t.Cleanup(func() { do(w.Destroy) })

	type scroll struct {
		x, y  float64
		phase glfw.GesturePhase
		mods  glfw.ModifierKey
	}
	var got []scroll
	plain := 0
	do(func() {
		w.SetScrollCallback(func(*glfw.Window, float64, float64) { plain++ })
		w.SetScrollExCallback(func(_ *glfw.Window, x, y float64, _ glfw.ScrollSource, p glfw.GesturePhase, m glfw.ModifierKey) {
			got = append(got, scroll{x, y, p, m & (glfw.ModControl | glfw.ModShift)})
		})
		hwnd := uintptr(unsafe.Pointer(w.GetWin32Window()))
		send := func(msg uintptr, keys uint16, delta int16) {
			wParam := uintptr(uint32(uint16(delta))<<16 | uint32(keys))
			sendMessage.Call(hwnd, msg, wParam, 0)
		}
		send(wmMouseWheel, mkControl, wheelDelta)
		send(wmMouseWheel, 0, -wheelDelta)
		send(wmMouseHWheel, mkShift, wheelDelta)
	})

	want := []scroll{
		{0, 1, glfw.GesturePhaseNone, glfw.ModControl},
		{0, -1, glfw.GesturePhaseNone, 0},
		{-1, 0, glfw.GesturePhaseNone, glfw.ModShift},
	}
	if len(got) != len(want) {
		t.Fatalf("scrolls: got %+v, want %+v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("scroll %d: got %+v, want %+v", i, got[i], want[i])
		}
	}
	if plain != 0 {
		t.Errorf("the plain callback was called %d times while the extended one was set", plain)
	}
}
