//go:build darwin && glfwcocoatest

// These tests hand GLFW's content view synthesized AppKit events, so they
// need a Mac with a display and are kept out of a plain `go test`:
//
//	go test -tags glfwcocoatest .
//
// They have their own tag and main-thread harness because the
// glfwintegration suite expects X11 or Wayland.

package glfw_test

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"testing"

	glfw "github.com/belluikit/glfw"
)

var callQ = make(chan func())

// do runs fn on the main thread, where Cocoa requires GLFW calls, and waits.
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

// newHiddenWindow creates a hidden, context-less window for one test.
func newHiddenWindow(t *testing.T) *glfw.Window {
	t.Helper()
	var (
		w   *glfw.Window
		err error
	)
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
	return w
}

type scrollEx struct {
	x, y   float64
	source glfw.ScrollSource
	phase  glfw.GesturePhase
	mods   glfw.ModifierKey
}

func (s scrollEx) String() string {
	return fmt.Sprintf("{%.2f %.2f source=%d phase=%d mods=%d}", s.x, s.y, s.source, s.phase, s.mods)
}

// A trackpad scroll reaches the extended callback with its source, every
// phase, the zero-movement start and end the plain callback never saw, and
// the coast marked as momentum; a wheel step as a wheel with no phase; and
// the event's modifiers (belluikit/glfw#13).
func TestCocoaScrollExReportsTheGesture(t *testing.T) {
	w := newHiddenWindow(t)
	var got []scrollEx
	plain := 0
	do(func() {
		w.SetScrollCallback(func(*glfw.Window, float64, float64) { plain++ })
		w.SetScrollExCallback(func(_ *glfw.Window, x, y float64, s glfw.ScrollSource, p glfw.GesturePhase, m glfw.ModifierKey) {
			got = append(got, scrollEx{x, y, s, p, m})
		})

		none := glfw.CocoaPhaseNone
		// Fingers rest, then move, then lift.
		glfw.SimulateCocoaScroll(w, 0, 0, true, glfw.CocoaPhaseMayBegin, none, 0)
		glfw.SimulateCocoaScroll(w, 0, 0, true, glfw.CocoaPhaseBegan, none, 0)
		glfw.SimulateCocoaScroll(w, 0, 20, true, glfw.CocoaPhaseChanged, none, 0)
		glfw.SimulateCocoaScroll(w, 3, 0, true, glfw.CocoaPhaseEnded, none, 0)
		// The coast: its own begin and end carry no movement.
		glfw.SimulateCocoaScroll(w, 0, 0, true, none, glfw.CocoaPhaseBegan, 0)
		glfw.SimulateCocoaScroll(w, 0, 10, true, none, glfw.CocoaPhaseChanged, 0)
		glfw.SimulateCocoaScroll(w, 0, 0, true, none, glfw.CocoaPhaseEnded, 0)
		// A wheel step with ctrl held, and an empty wheel event.
		glfw.SimulateCocoaScroll(w, 0, -1, false, none, none, glfw.CocoaModifierControl)
		glfw.SimulateCocoaScroll(w, 0, 0, false, none, none, 0)
		// Fingers that lift without moving cancel.
		glfw.SimulateCocoaScroll(w, 0, 0, true, glfw.CocoaPhaseMayBegin, none, 0)
		glfw.SimulateCocoaScroll(w, 0, 0, true, glfw.CocoaPhaseCancelled, none, glfw.CocoaModifierShift)
	})

	tp, wheel := glfw.ScrollSourceTouchpad, glfw.ScrollSourceWheel
	want := []scrollEx{
		{0, 0, tp, glfw.GesturePhaseBegan, 0},
		{0, 2, tp, glfw.GesturePhaseChanged, 0},
		{0, 0, tp, glfw.GesturePhaseEnded, 0},
		{0, 1, tp, glfw.GesturePhaseMomentum, 0},
		{0, -1, wheel, glfw.GesturePhaseNone, glfw.ModControl},
		{0, 0, tp, glfw.GesturePhaseCancelled, glfw.ModShift},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("extended scrolls:\n got %v\nwant %v", got, want)
	}
	if plain != 0 {
		t.Errorf("the plain callback was called %d times while the extended one was set", plain)
	}
}

// Without an extended callback the plain one gets the scrolls that move, as
// before, and none of the zero-movement gesture events.
func TestCocoaPlainScrollUnchanged(t *testing.T) {
	w := newHiddenWindow(t)
	var got [][2]float64
	do(func() {
		w.SetScrollCallback(func(_ *glfw.Window, x, y float64) { got = append(got, [2]float64{x, y}) })
		none := glfw.CocoaPhaseNone
		glfw.SimulateCocoaScroll(w, 0, 0, true, glfw.CocoaPhaseBegan, none, 0)
		glfw.SimulateCocoaScroll(w, 0, 20, true, glfw.CocoaPhaseChanged, none, 0)
		glfw.SimulateCocoaScroll(w, 0, 0, true, glfw.CocoaPhaseEnded, none, 0)
		glfw.SimulateCocoaScroll(w, 0, 10, true, none, glfw.CocoaPhaseChanged, 0)
		glfw.SimulateCocoaScroll(w, 0, -1, false, none, none, 0)
	})
	want := [][2]float64{{0, 2}, {0, 1}, {0, -1}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("plain scrolls: got %v, want %v", got, want)
	}
}

// A trackpad pinch reaches the pinch callback, each step relative to the
// last, and a scale of 1 when it ends (belluikit/glfw#13).
func TestCocoaPinchReportsMagnify(t *testing.T) {
	w := newHiddenWindow(t)
	type step struct {
		scale float64
		phase glfw.GesturePhase
	}
	var got []step
	do(func() {
		w.SetPinchCallback(func(_ *glfw.Window, s float64, p glfw.GesturePhase) { got = append(got, step{s, p}) })
		glfw.SimulateCocoaMagnify(w, 0.10, glfw.CocoaPhaseBegan)
		glfw.SimulateCocoaMagnify(w, -0.05, glfw.CocoaPhaseChanged)
		glfw.SimulateCocoaMagnify(w, 0.02, glfw.CocoaPhaseEnded)
		glfw.SimulateCocoaMagnify(w, 0.30, glfw.CocoaPhaseNone)
	})
	want := []step{{1.10, glfw.GesturePhaseBegan}, {0.95, glfw.GesturePhaseChanged}, {1, glfw.GesturePhaseEnded}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("pinch steps: got %v, want %v", got, want)
	}
}
