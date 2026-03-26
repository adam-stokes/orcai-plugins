package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if err := run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr, os.Getenv); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// modelEntry is a single model entry for --list-models JSON output.
type modelEntry struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

// listModels runs `ollama list`, parses its output, and writes JSON to stdout.
func listModels(stdout io.Writer, stderr io.Writer) error {
	out, err := exec.Command("ollama", "list").Output()
	if err != nil {
		return fmt.Errorf("ollama list: %w", err)
	}
	var models []modelEntry
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for i, line := range lines {
		if i == 0 { // skip header
			continue
		}
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		id := fields[0]
		models = append(models, modelEntry{ID: id, Label: id})
	}
	if models == nil {
		models = []modelEntry{}
	}
	data, err := json.Marshal(models)
	if err != nil {
		return fmt.Errorf("marshalling models: %w", err)
	}
	fmt.Fprintln(stdout, string(data))
	return nil
}

// run is the testable entry point.
// args are the CLI arguments (os.Args[1:] in production).
// getenv is the environment resolver (os.Getenv in production).
//
// Model resolution priority: --model flag > ORCAI_MODEL env var.
// URL resolution priority:   ORCAI_OLLAMA_URL env var > default http://localhost:11434.
func run(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer, getenv func(string) string) error {
	// Handle --list-models before flag parsing.
	for _, arg := range args {
		if arg == "--list-models" {
			return listModels(stdout, stderr)
		}
	}

	fs := flag.NewFlagSet("orcai-ollama", flag.ContinueOnError)
	fs.SetOutput(stderr)
	modelFlag := fs.String("model", "", "Ollama model name (e.g. llama3.2)")
	if err := fs.Parse(args); err != nil {
		return err
	}

	// Read prompt from stdin.
	promptBytes, err := io.ReadAll(stdin)
	if err != nil {
		return fmt.Errorf("reading stdin: %w", err)
	}
	prompt := strings.TrimSpace(string(promptBytes))
	if prompt == "" {
		return fmt.Errorf("prompt is required: no input received on stdin")
	}

	// Resolve model: --model flag takes precedence over ORCAI_MODEL env var.
	model := *modelFlag
	if model == "" {
		model = getenv("ORCAI_MODEL")
	}
	if model == "" {
		return fmt.Errorf("model is required: set --model flag or ORCAI_MODEL environment variable")
	}

	// Resolve base URL from ORCAI_OLLAMA_URL, defaulting to localhost.
	baseURL := getenv("ORCAI_OLLAMA_URL")
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}
	baseURL = strings.TrimRight(baseURL, "/")

	result, err := callOllama(baseURL, model, prompt)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			fmt.Fprintf(stderr, "Model %q not found locally — pulling...\n", model)
			pull := exec.Command("ollama", "pull", model)
			pull.Stdout = stderr
			pull.Stderr = stderr
			if pullErr := pull.Run(); pullErr != nil {
				return fmt.Errorf("ollama pull %s: %w", model, pullErr)
			}
			result, err = callOllama(baseURL, model, prompt)
		}
		if err != nil {
			return err
		}
	}

	fmt.Fprint(stdout, result)
	return nil
}

// generateRequest is the JSON body for POST /api/generate.
type generateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// generateResponse is the JSON body returned by Ollama when stream=false.
type generateResponse struct {
	Response string `json:"response"`
	Error    string `json:"error,omitempty"`
}

// callOllama sends a prompt to the Ollama /api/generate endpoint and returns the completion.
func callOllama(baseURL, model, prompt string) (string, error) {
	body, err := json.Marshal(generateRequest{
		Model:  model,
		Prompt: prompt,
		Stream: false,
	})
	if err != nil {
		return "", fmt.Errorf("marshalling request: %w", err)
	}

	resp, err := http.Post(baseURL+"/api/generate", "application/json", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("connecting to Ollama at %s: %w", baseURL, err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("Ollama returned HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(respBytes)))
	}

	var genResp generateResponse
	if err := json.Unmarshal(respBytes, &genResp); err != nil {
		return "", fmt.Errorf("parsing response: %w", err)
	}
	if genResp.Error != "" {
		return "", fmt.Errorf("Ollama error: %s", genResp.Error)
	}

	return genResp.Response, nil
}
