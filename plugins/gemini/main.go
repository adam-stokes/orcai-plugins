package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func main() {
	code, err := run(os.Stdin, os.Stdout, os.Stderr, os.Getenv, execGemini)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	os.Exit(code)
}

// execGemini runs `gemini [--model <model>]` with the prompt on stdin.
func execGemini(model, prompt string, stdout, stderr io.Writer) int {
	var args []string
	if model != "" {
		args = append(args, "--model", model)
	}
	cmd := exec.Command("gemini", args...)
	cmd.Stdin = strings.NewReader(prompt)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		if exit, ok := err.(*exec.ExitError); ok {
			return exit.ExitCode()
		}
		fmt.Fprintln(stderr, err)
		return 1
	}
	return 0
}

// run is the testable entry point.
func run(
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
	getenv func(string) string,
	executor func(model, prompt string, stdout, stderr io.Writer) int,
) (int, error) {
	if _, err := exec.LookPath("gemini"); err != nil {
		return 1, fmt.Errorf("gemini CLI not found in PATH: install it from https://ai.google.dev/gemini-api/docs/gemini-cli")
	}

	promptBytes, err := io.ReadAll(stdin)
	if err != nil {
		return 1, fmt.Errorf("reading stdin: %w", err)
	}
	prompt := strings.TrimSpace(string(promptBytes))
	if prompt == "" {
		return 1, fmt.Errorf("prompt is required: no input received on stdin")
	}

	model := getenv("ORCAI_MODEL")
	return executor(model, prompt, stdout, stderr), nil
}
