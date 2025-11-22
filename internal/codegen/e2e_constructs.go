package codegen

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func runEndToEnd(t *testing.T, bytecode []int, expectedExitCode int, message string) {
	t.Helper()

	cg := NewCodeGen()
	asm := cg.Generate(bytecode)

	tmpDir := t.TempDir()
	asmFile := tmpDir + "/test.s"
	objFile := tmpDir + "/test.o"
	exeFile := tmpDir + "/test"

	err := os.WriteFile(asmFile, []byte(asm), 0644)
	assert.NoError(t, err)

	cmd := exec.Command("as", "-o", objFile, asmFile)
	output, err := cmd.CombinedOutput()
	assert.NoError(t, err, "assembler failed: %s", string(output))

	cmd = exec.Command("ld", "-o", exeFile, objFile)
	output, err = cmd.CombinedOutput()
	assert.NoError(t, err, "linker failed: %s", string(output))

	cmd = exec.Command(exeFile)
	cmd.Run()

	exitCode := cmd.ProcessState.ExitCode()
	assert.Equal(t, expectedExitCode, exitCode, message)
}
