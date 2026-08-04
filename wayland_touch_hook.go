//go:build (linux && !x11 && !wayland) || (linux && wayland)
// +build linux,!x11,!wayland linux,wayland

package glfw

//#include <stdlib.h>
//extern int _glfwTestTouchDownTargetsWindowWayland(void* window, void* surface);
import "C"

// testTouchDownNullWindowTargets bridges to the C touchDownTargetsWindow
// helper for the belluikit/glfw#1 touch-down guard test. It reports whether a
// wl_touch "down" resolving to a NULL window (a foreign surface — one not
// tagged as a glfw window) is treated as targeting a window; a correct guard
// makes this false without dereferencing the NULL window (cf. the glfw#17
// drag-enter crash). cgo is not permitted in _test.go files, so the bridge
// lives here.
func testTouchDownNullWindowTargets() bool {
	// A non-nil, never-dereferenced surface pointer stands in for the foreign
	// surface; the guard must short-circuit on the NULL window before touching
	// it. malloc(1) is a valid distinct address without a real wl_surface
	// (which would need a compositor).
	surface := C.malloc(1)
	defer C.free(surface)

	return C._glfwTestTouchDownTargetsWindowWayland(nil, surface) != 0
}
