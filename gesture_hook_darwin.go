//go:build darwin && glfwcocoatest

package glfw

/*
#import <Cocoa/Cocoa.h>
#define GLFW_INCLUDE_NONE
#include "glfw/include/GLFW/glfw3.h"
#define GLFW_EXPOSE_NATIVE_COCOA
#include "glfw/include/GLFW/glfw3native.h"

// An NSEvent whose scroll and magnify properties are set by the test
// (belluikit/glfw#13). AppKit gives no way to make a scroll event with a
// phase, or a magnify event, so the test answers the properties GLFW's
// content view reads. It is never initialised as a real event, and is never
// released, since NSEvent's dealloc expects one.
@interface GLFWTestGestureEvent : NSEvent
@property CGFloat testDeltaX, testDeltaY, testMagnification;
@property BOOL testPrecise;
@property NSEventPhase testPhase, testMomentumPhase;
@property NSEventModifierFlags testFlags;
@end

@implementation GLFWTestGestureEvent
- (CGFloat)scrollingDeltaX { return self.testDeltaX; }
- (CGFloat)scrollingDeltaY { return self.testDeltaY; }
- (BOOL)hasPreciseScrollingDeltas { return self.testPrecise; }
- (NSEventPhase)phase { return self.testPhase; }
- (NSEventPhase)momentumPhase { return self.testMomentumPhase; }
- (NSEventModifierFlags)modifierFlags { return self.testFlags; }
- (CGFloat)magnification { return self.testMagnification; }
@end

static void glfwTestCocoaScroll(GLFWwindow* window, double dx, double dy, int precise,
                                unsigned long phase, unsigned long momentum, unsigned long flags)
{
    GLFWTestGestureEvent* e = [GLFWTestGestureEvent alloc];
    e.testDeltaX = dx;
    e.testDeltaY = dy;
    e.testPrecise = precise ? YES : NO;
    e.testPhase = (NSEventPhase) phase;
    e.testMomentumPhase = (NSEventPhase) momentum;
    e.testFlags = (NSEventModifierFlags) flags;
    [[glfwGetCocoaWindow(window) contentView] scrollWheel:e];
}

static void glfwTestCocoaMagnify(GLFWwindow* window, double magnification, unsigned long phase)
{
    GLFWTestGestureEvent* e = [GLFWTestGestureEvent alloc];
    e.testMagnification = magnification;
    e.testPhase = (NSEventPhase) phase;
    [[glfwGetCocoaWindow(window) contentView] magnifyWithEvent:e];
}
*/
import "C"

// CocoaEventPhase is an NSEventPhase, for SimulateCocoaScroll and
// SimulateCocoaMagnify.
type CocoaEventPhase uint

// The NSEventPhase values.
const (
	CocoaPhaseNone       CocoaEventPhase = 0
	CocoaPhaseBegan      CocoaEventPhase = 1 << 0
	CocoaPhaseStationary CocoaEventPhase = 1 << 1
	CocoaPhaseChanged    CocoaEventPhase = 1 << 2
	CocoaPhaseEnded      CocoaEventPhase = 1 << 3
	CocoaPhaseCancelled  CocoaEventPhase = 1 << 4
	CocoaPhaseMayBegin   CocoaEventPhase = 1 << 5
)

// CocoaModifierControl and CocoaModifierShift are NSEventModifierFlags bits.
const (
	CocoaModifierShift   uint = 1 << 17
	CocoaModifierControl uint = 1 << 18
)

// SimulateCocoaScroll hands GLFW's content view a scroll event with the given
// deltas, precision, phases and modifier flags, as AppKit would
// (belluikit/glfw#13). It exists only in glfwcocoatest builds, for the
// tests; call it on the main thread.
func SimulateCocoaScroll(w *Window, dx, dy float64, precise bool, phase, momentum CocoaEventPhase, flags uint) {
	p := C.int(0)
	if precise {
		p = 1
	}
	C.glfwTestCocoaScroll(w.data, C.double(dx), C.double(dy), p, C.ulong(phase), C.ulong(momentum), C.ulong(flags))
}

// SimulateCocoaMagnify hands GLFW's content view a magnify (pinch) event, as
// AppKit would (belluikit/glfw#13). It exists only in glfwcocoatest builds,
// for the tests; call it on the main thread.
func SimulateCocoaMagnify(w *Window, magnification float64, phase CocoaEventPhase) {
	C.glfwTestCocoaMagnify(w.data, C.double(magnification), C.ulong(phase))
}
