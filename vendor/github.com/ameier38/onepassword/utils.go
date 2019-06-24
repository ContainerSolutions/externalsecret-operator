package onepassword

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func getExtension() string {
	if runtime.GOOS == "windows" {
		return ".exe"
	}
	return ""
}

func buildMockOnePassword() (string, error) {
	cmd := exec.Command(
		"go",
		"install",
		"github.com/ameier38/onepassword/cmd/mock-op")

	if output, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build mock op program: %s\n%s", err, output)
	}

	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = filepath.Join(os.Getenv("HOME"), "go")
	}

	programPath := filepath.Join(
		filepath.SplitList(gopath)[0],
		"bin",
		"mock-op"+getExtension())

	return programPath, nil
}
