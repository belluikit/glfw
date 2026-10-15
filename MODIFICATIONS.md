# Modifications to the bundled GLFW C sources

The `glfw/` directory contains the C sources of [GLFW](https://www.glfw.org),
taken from [glfw/glfw](https://github.com/glfw/glfw) at commit
`7b6aead9fb88b3623e3b3725ebb42670cbe4c579` (GLFW 3.4, recorded in
[GLFW_C_REVISION.txt](GLFW_C_REVISION.txt)), as vendored by
[go-gl/glfw](https://github.com/go-gl/glfw) under `v3.4/glfw/glfw`.

GLFW is distributed under the zlib/libpng license
([glfw/LICENSE.md](glfw/LICENSE.md), unchanged), whose clause 2 requires that
altered source versions be plainly marked as such and not be misrepresented as
the original software. This file is that record. **The bundled tree is not
pristine upstream GLFW.**

## Baseline

Files were compared against upstream glfw/glfw at the revision above, not
against go-gl/glfw. A file counts as altered if it differs from that upstream
file in any way, including changes that came from go-gl/glfw rather than from
this fork. Files present upstream but not bundled here (docs, examples, tests,
CMake helpers) are omitted, not altered.

Every altered file carries a one-line notice right after its original license
block, pointing here. The license blocks themselves are untouched.

## Altered files

Alterations made by go-gl/glfw (carried unchanged here):

- `glfw/src/posix_poll.c`: `#define _GNU_SOURCE` wrapped in `#ifndef _GNU_SOURCE`, so it is not redefined when the build already defines it. This is the file's only change.
- `glfw/src/wl_window.c`: the same `_GNU_SOURCE` guard (plus this fork's changes below).

Alterations made by this fork (github.com/belluikit/glfw):

- `glfw/include/GLFW/glfw3.h`: public API for IME/preedit, generic MIME clipboard flavors, drag-motion callback (`glfwSetDragCallback`) and touch input (`glfwSetTouchCallback`).
- `glfw/include/GLFW/glfw3native.h`: native accessors for window parenting (X11/Wayland), window modality (X11, Wayland xdg-dialog-v1), xdg-foreign export handles and Wayland drag-out (`glfwStartWaylandDrag`).
- `glfw/src/CMakeLists.txt`: generate the xdg-dialog-v1 and xdg-foreign-unstable-v2 Wayland protocol headers.
- `glfw/src/cocoa_init.m`: IME/preedit support; MIME clipboard flavors.
- `glfw/src/cocoa_platform.h`: IME/preedit support; MIME clipboard flavors.
- `glfw/src/cocoa_window.m`: IME/preedit support and the correct input-context discard; MIME clipboard flavors; drag-motion callback.
- `glfw/src/init.c`: IME/preedit support.
- `glfw/src/input.c`: IME/preedit support; MIME clipboard flavors; drag-motion callback; touch input.
- `glfw/src/internal.h`: IME/preedit support; MIME clipboard flavors; drag-motion callback; touch input.
- `glfw/src/null_init.c`: IME/preedit and MIME clipboard stubs for the null platform.
- `glfw/src/null_platform.h`: IME/preedit and MIME clipboard stubs for the null platform.
- `glfw/src/null_window.c`: IME/preedit and MIME clipboard stubs for the null platform.
- `glfw/src/win32_init.c`: IME/preedit support; MIME clipboard flavors; OLE drop-target drag-motion; touch input.
- `glfw/src/win32_platform.h`: IME/preedit support; MIME clipboard flavors; OLE drop-target drag-motion; touch input.
- `glfw/src/win32_window.c`: IME/preedit support; MIME clipboard flavors; OLE drop-target drag-motion; touch input, including cancelling touches on teardown.
- `glfw/src/window.c`: IME/preedit support.
- `glfw/src/wl_init.c`: IME (text-input protocols); MIME clipboard flavors; xdg-dialog modality; xdg-foreign export; drag-out; survive compositors without a seat; touch input.
- `glfw/src/wl_platform.h`: IME (text-input protocols) and the preedit commit/done feedback-loop fix; MIME clipboard flavors; modality; xdg-foreign; drag-out; touch input.
- `glfw/src/wl_window.c`: IME and the preedit feedback-loop fix; MIME clipboard flavors; window parenting and modality; xdg-foreign; drag-out; seatless-compositor and data-offer listener fixes (a crash in `glfwStartWaylandDrag`); NULL-window guard on drag enter; touch input.
- `glfw/src/x11_init.c`: IME/preedit support; MIME clipboard flavors; window modality; drag-motion callback.
- `glfw/src/x11_platform.h`: IME/preedit support; MIME clipboard flavors; window modality; drag-motion callback.
- `glfw/src/x11_window.c`: IME/preedit support; MIME clipboard flavors; window parenting and modality; reaping INCR clipboard transfers to a requestor that went away; drag-motion callback.

The IME/preedit support is a port of the unmerged upstream pull request
[glfw/glfw#2130](https://github.com/glfw/glfw/pull/2130).

## Added files

These files are not in upstream GLFW. They need no modification notice, but
are listed so the tree is fully accounted for.

Added by this fork, Wayland protocol definitions (from wayland-protocols, each
under its own license in the file) and their `wayland-scanner` output:

- `glfw/deps/wayland/text-input-unstable-v1.xml`, `glfw/deps/wayland/text-input-unstable-v3.xml` (IME)
- `glfw/deps/wayland/xdg-dialog-v1.xml` (window modality)
- `glfw/deps/wayland/xdg-foreign-unstable-v2.xml` (portal parenting)
- `glfw/include/text-input-unstable-v1-client-protocol.h`, `glfw/include/text-input-unstable-v1-client-protocol-code.h`
- `glfw/include/text-input-unstable-v3-client-protocol.h`, `glfw/include/text-input-unstable-v3-client-protocol-code.h`
- `glfw/include/xdg-dialog-v1-client-protocol.h`, `glfw/include/xdg-dialog-v1-client-protocol-code.h`
- `glfw/include/xdg-foreign-unstable-v2-client-protocol.h`, `glfw/include/xdg-foreign-unstable-v2-client-protocol-code.h`

Added by go-gl/glfw: `wayland-scanner` output for the protocols upstream GLFW
ships, generated at vendoring time instead of by CMake:

- `glfw/include/{wayland,viewporter,xdg-shell,idle-inhibit-unstable-v1,pointer-constraints-unstable-v1,relative-pointer-unstable-v1,fractional-scale-v1,xdg-activation-v1,xdg-decoration-unstable-v1}-client-protocol.h`
- the matching `*-client-protocol-code.h` file for each of those protocols

Added by go-gl/glfw and adapted here: a `dummy.go` in each directory under
`glfw/`, so Go tooling keeps the C files when vendoring. `glfw/deps/dummy.go`
imports this module's path (`github.com/belluikit/glfw/glfw/deps/...`).
