package glfw

import (
	"bufio"
	stdbytes "bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestVendoredWaylandProtocolsSurviveRegeneration checks that
// scripts/generate-buildable-source-code.sh would reproduce every protocol in
// glfw/deps/wayland. The script rebuilds glfw/ from a pristine upstream tarball,
// so each vendored XML must either be one upstream ships (named in the
// script's protocol_files list) or be kept in scripts/fork-wayland-protocols,
// which the script overlays and scans. A protocol in neither place would be
// deleted by the next regeneration, and the C code using it would stop
// compiling (glfw#4).
func TestVendoredWaylandProtocolsSurviveRegeneration(t *testing.T) {
	upstream := upstreamProtocolList(t, "scripts/generate-buildable-source-code.sh")
	if len(upstream) == 0 {
		t.Fatal("found no protocol_files entries in the regeneration script; has its layout changed?")
	}

	vendored, err := filepath.Glob("glfw/deps/wayland/*.xml")
	if err != nil || len(vendored) == 0 {
		t.Fatalf("no vendored protocols found (err=%v)", err)
	}

	for _, path := range vendored {
		name := filepath.Base(path)
		overlay := filepath.Join("scripts/fork-wayland-protocols", name)
		_, inOverlay := os.Stat(overlay)
		switch {
		case upstream[name] && inOverlay == nil:
			t.Errorf("%s is both in the upstream list and in scripts/fork-wayland-protocols; keep it in one", name)
		case upstream[name]:
			// Regenerated from the upstream tarball.
		case inOverlay == nil:
			want, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			got, err := os.ReadFile(overlay)
			if err != nil {
				t.Fatal(err)
			}
			if !stdbytes.Equal(got, want) {
				t.Errorf("%s differs from %s; a regeneration would vendor the overlay copy", overlay, path)
			}
		default:
			t.Errorf("%s is neither in the script's upstream protocol_files list nor in scripts/fork-wayland-protocols; a regeneration would delete it", path)
		}

		base := strings.TrimSuffix(name, ".xml")
		for _, header := range []string{base + "-client-protocol.h", base + "-client-protocol-code.h"} {
			if _, err := os.Stat(filepath.Join("glfw/include", header)); err != nil {
				t.Errorf("%s has no generated glfw/include/%s", name, header)
			}
		}
	}

	overlays, err := filepath.Glob("scripts/fork-wayland-protocols/*.xml")
	if err != nil {
		t.Fatal(err)
	}
	for _, overlay := range overlays {
		if _, err := os.Stat(filepath.Join("glfw/deps/wayland", filepath.Base(overlay))); err != nil {
			t.Errorf("%s is not vendored in glfw/deps/wayland", overlay)
		}
	}
}

// upstreamProtocolList reads the names in the script's protocol_files=( ... )
// array, one per line.
func upstreamProtocolList(t *testing.T, script string) map[string]bool {
	t.Helper()
	f, err := os.Open(script)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	names := map[string]bool{}
	in := false
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		switch {
		case strings.HasSuffix(line, "protocol_files=("):
			in = true
		case in && line == ")":
			return names
		case in && strings.HasSuffix(line, ".xml"):
			names[line] = true
		}
	}
	if err := sc.Err(); err != nil {
		t.Fatal(err)
	}
	t.Fatal("protocol_files list in the regeneration script is not closed")
	return nil
}
