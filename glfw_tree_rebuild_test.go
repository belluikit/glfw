package glfw

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestUpstreamTreeSHAIsCurrent fails when upstreamTreeSHA no longer matches
// the glfw/ tree, so a change to the C sources can't land without the Go
// build cache key that tracks them (glfw#11). It computes the hash the way
// scripts/glfw_tree_rebuild.sh does; run `go generate` to refresh it.
//
// It needs a git checkout and skips without one (in the module cache, say).
func TestUpstreamTreeSHAIsCurrent(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not found")
	}
	if err := exec.Command("git", "rev-parse", "--show-toplevel").Run(); err != nil {
		t.Skip("not a git checkout")
	}

	// A throwaway index, so the real one is untouched.
	index := filepath.Join(t.TempDir(), "index")
	git := func(args ...string) string {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Env = append(os.Environ(), "GIT_INDEX_FILE="+index)
		out, err := cmd.Output()
		if err != nil {
			t.Fatalf("git %s: %v", strings.Join(args, " "), err)
		}
		return strings.TrimSpace(string(out))
	}
	git("add", "glfw")
	got := git("write-tree")

	if got != upstreamTreeSHA {
		t.Errorf("upstreamTreeSHA = %s, but the glfw/ tree hashes to %s; run `go generate` and commit glfw_tree_rebuild.go", upstreamTreeSHA, got)
	}
}
