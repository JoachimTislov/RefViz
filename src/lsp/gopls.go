package lsp

import (
	"os/exec"
)

// runs the gopls command with the given arguments
func RunGopls(projectPath string, args ...string) ([]byte, error) {
	a := []string{"-vv", "-rpc.trace"}
	c := exec.Command("gopls", append(a, args...)...)
	c.Dir = projectPath
	return c.Output()
}
