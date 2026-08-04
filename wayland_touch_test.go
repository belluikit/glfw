//go:build (linux && !x11 && !wayland) || (linux && wayland)
// +build linux,!x11,!wayland linux,wayland

package glfw

import "testing"

// TestWaylandTouchDownNullWindowGuard covers the belluikit/glfw#1 touch-down
// guard: a wl_touch "down" carrying a foreign surface (one not tagged as a
// glfw window) resolves to a NULL window, and the guard must reject it without
// dereferencing the window — the same shape as the glfw#17 drag-enter crash,
// which SIGSEGVed at the wl.surface field offset.
//
// It calls the real touchDownTargetsWindow helper through a C test hook, so
// removing the NULL guard re-crashes this test. No glfwInit or compositor is
// needed — it is pure pointer logic, which is why it can cover a path no
// headless compositor emits.
func TestWaylandTouchDownNullWindowGuard(t *testing.T) {
	if testTouchDownNullWindowTargets() {
		t.Fatal("NULL-window touch down treated as a target — belluikit/glfw#1 guard missing")
	}
}
