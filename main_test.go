// +build integration

package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var (
	cPath  = "build/a.out"
	goPath = "build/go.out"
	stdin  = "7"
	args   = []string{"some", "args"}
)

type programOut struct {
	stdout bytes.Buffer
	stderr bytes.Buffer
	isZero bool
}

// TestIntegrationScripts tests all programs in the tests directory.
//
// Integration tests are not run by default (only unit tests). These are
// indicated by the build flags at the top of the file. To include integration
// tests use:
//
//     go test -tags=integration
//
// You can also run a single file with:
//
//     go test -tags=integration -run=TestIntegrationScripts/tests/ctype/isalnum.c
//
func TestIntegrationScripts(t *testing.T) {
	files, err := filepath.Glob("tests/*.c")
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range files {
		// Create build folder
		os.Mkdir("build/", os.ModePerm)

		t.Run(file, func(t *testing.T) {
			cProgram := programOut{}
			goProgram := programOut{}

			// Compile C.
			out, err := exec.Command("clang", "-lm", "-o", cPath, file).CombinedOutput()
			if err != nil {
				t.Fatalf("error: %s\n%s", err, out)
			}

			// Run C program
			cmd := exec.Command(cPath, args...)
			cmd.Stdin = strings.NewReader(stdin)
			cmd.Stdout = &cProgram.stdout
			cmd.Stderr = &cProgram.stderr
			err = cmd.Run()
			cProgram.isZero = err == nil

			// Compile Go
			goSrc := Start([]string{file})
			// fmt.Println(goSrc)
			ioutil.WriteFile("build/main.go", []byte(goSrc), os.ModePerm)
			buildErr, err := exec.Command("go", "build", "-o", goPath, "build/main.go").CombinedOutput()
			if err != nil {
				t.Fatal(string(buildErr), err)
			}

			// Run Go program
			cmd = exec.Command(goPath, args...)
			cmd.Stdin = strings.NewReader(stdin)
			cmd.Stdout = &goProgram.stdout
			cmd.Stderr = &goProgram.stderr
			err = cmd.Run()
			goProgram.isZero = err == nil

			// Check if both exit codes are zero (or non-zero)
			if cProgram.isZero != goProgram.isZero {
				t.Fatalf("Expected: %t, Got: %t", cProgram.isZero, goProgram.isZero)
			}

			// Check stderr
			if cProgram.stderr.String() != goProgram.stderr.String() {
				t.Fatalf("Expected %q, Got: %q", cProgram.stderr.String(), goProgram.stderr.String())
			}

			// Check stdout
			if cProgram.stdout.String() != goProgram.stdout.String() {
				t.Fatalf("Expected %q, Got: %q", cProgram.stdout.String(), goProgram.stdout.String())
			}
		})
	}
}
