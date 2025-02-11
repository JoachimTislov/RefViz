package lsp

import "os/exec"

// runs the gopls command with the given arguments
func RunGopls(args ...string) ([]byte, error) {
	a := []string{"-vv", "-rpc.trace"}
	return exec.Command("gopls", append(a, args...)...).Output()
}
