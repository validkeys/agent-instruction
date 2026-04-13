package commands

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestPromptYesNo(t *testing.T) {
	tests := map[string]struct {
		input      string
		defaultYes bool
		want       bool
	}{
		"accepts yes": {
			input:      "yes\n",
			defaultYes: false,
			want:       true,
		},
		"accepts y": {
			input:      "y\n",
			defaultYes: false,
			want:       true,
		},
		"accepts Y uppercase": {
			input:      "Y\n",
			defaultYes: false,
			want:       true,
		},
		"accepts no": {
			input:      "no\n",
			defaultYes: true,
			want:       false,
		},
		"accepts n": {
			input:      "n\n",
			defaultYes: true,
			want:       false,
		},
		"uses default yes on empty": {
			input:      "\n",
			defaultYes: true,
			want:       true,
		},
		"uses default no on empty": {
			input:      "\n",
			defaultYes: false,
			want:       false,
		},
		"uses default on whitespace": {
			input:      "   \n",
			defaultYes: true,
			want:       true,
		},
		"treats invalid as default yes": {
			input:      "invalid\n",
			defaultYes: true,
			want:       true,
		},
		"treats invalid as default no": {
			input:      "invalid\n",
			defaultYes: false,
			want:       false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			cmd := &cobra.Command{}
			in := strings.NewReader(tt.input)
			out := &bytes.Buffer{}
			cmd.SetIn(in)
			cmd.SetOut(out)

			got := promptYesNo(cmd, "Test question", tt.defaultYes)
			if got != tt.want {
				t.Errorf("promptYesNo() = %v, want %v", got, tt.want)
			}

			// Check output contains the question
			output := out.String()
			if !strings.Contains(output, "Test question") {
				t.Errorf("output doesn't contain question: %s", output)
			}
		})
	}
}

func TestPromptFrameworks(t *testing.T) {
	tests := map[string]struct {
		input string
		want  []string
	}{
		"selects claude": {
			input: "1\n",
			want:  []string{"claude"},
		},
		"selects agents": {
			input: "2\n",
			want:  []string{"agents"},
		},
		"selects both": {
			input: "3\n",
			want:  []string{"claude", "agents"},
		},
		"defaults to both on empty": {
			input: "\n",
			want:  []string{"claude", "agents"},
		},
		"defaults to both on invalid": {
			input: "invalid\n",
			want:  []string{"claude", "agents"},
		},
		"defaults to both on out of range": {
			input: "99\n",
			want:  []string{"claude", "agents"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			cmd := &cobra.Command{}
			in := strings.NewReader(tt.input)
			out := &bytes.Buffer{}
			cmd.SetIn(in)
			cmd.SetOut(out)

			got := promptFrameworks(cmd)

			if len(got) != len(tt.want) {
				t.Fatalf("len(got) = %d, want %d", len(got), len(tt.want))
			}

			for i, fw := range got {
				if fw != tt.want[i] {
					t.Errorf("got[%d] = %q, want %q", i, fw, tt.want[i])
				}
			}

			// Check output contains options
			output := out.String()
			if !strings.Contains(output, "claude") {
				t.Error("output doesn't contain claude option")
			}
			if !strings.Contains(output, "agents") {
				t.Error("output doesn't contain agents option")
			}
		})
	}
}

func TestPromptPackages(t *testing.T) {
	tests := map[string]struct {
		input string
		want  []string
	}{
		"selects auto": {
			input: "1\n",
			want:  []string{"auto"},
		},
		"defaults to auto on empty": {
			input: "\n",
			want:  []string{"auto"},
		},
		"selects manual with single package": {
			input: "2\napp\n",
			want:  []string{"app"},
		},
		"selects manual with multiple packages": {
			input: "2\napp,lib,services\n",
			want:  []string{"app", "lib", "services"},
		},
		"selects manual with spaces": {
			input: "2\napp , lib , services\n",
			want:  []string{"app", "lib", "services"},
		},
		"defaults to auto on manual with empty packages": {
			input: "2\n\n",
			want:  []string{"auto"},
		},
		"defaults to auto on manual with only whitespace": {
			input: "2\n   \n",
			want:  []string{"auto"},
		},
		"defaults to auto on invalid choice": {
			input: "99\n",
			want:  []string{"auto"},
		},
		"filters empty package names": {
			input: "2\napp,,lib,  ,services\n",
			want:  []string{"app", "lib", "services"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			cmd := &cobra.Command{}
			in := strings.NewReader(tt.input)
			out := &bytes.Buffer{}
			cmd.SetIn(in)
			cmd.SetOut(out)

			got := promptPackages(cmd)

			if len(got) != len(tt.want) {
				t.Fatalf("len(got) = %d, want %d\ngot: %v\nwant: %v", len(got), len(tt.want), got, tt.want)
			}

			for i, pkg := range got {
				if pkg != tt.want[i] {
					t.Errorf("got[%d] = %q, want %q", i, pkg, tt.want[i])
				}
			}

			// Check output contains options
			output := out.String()
			if !strings.Contains(output, "auto") {
				t.Error("output doesn't contain auto option")
			}
			if !strings.Contains(output, "manual") {
				t.Error("output doesn't contain manual option")
			}
		})
	}
}

func TestPromptYesNoDisplaysCorrectDefault(t *testing.T) {
	tests := map[string]struct {
		defaultYes   bool
		wantInOutput string
	}{
		"shows Y/n for default yes": {
			defaultYes:   true,
			wantInOutput: "Y/n",
		},
		"shows y/N for default no": {
			defaultYes:   false,
			wantInOutput: "y/N",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			cmd := &cobra.Command{}
			in := strings.NewReader("\n")
			out := &bytes.Buffer{}
			cmd.SetIn(in)
			cmd.SetOut(out)

			promptYesNo(cmd, "Question", tt.defaultYes)

			output := out.String()
			if !strings.Contains(output, tt.wantInOutput) {
				t.Errorf("output = %q, want to contain %q", output, tt.wantInOutput)
			}
		})
	}
}

func TestPromptFrameworksOutputFormat(t *testing.T) {
	cmd := &cobra.Command{}
	in := strings.NewReader("3\n")
	out := &bytes.Buffer{}
	cmd.SetIn(in)
	cmd.SetOut(out)

	promptFrameworks(cmd)

	output := out.String()

	// Check that options are numbered
	if !strings.Contains(output, "1)") {
		t.Error("output should contain numbered option 1)")
	}
	if !strings.Contains(output, "2)") {
		t.Error("output should contain numbered option 2)")
	}
	if !strings.Contains(output, "3)") {
		t.Error("output should contain numbered option 3)")
	}

	// Check that default is shown
	if !strings.Contains(output, "default") {
		t.Error("output should mention default")
	}
}

func TestPromptPackagesOutputFormat(t *testing.T) {
	cmd := &cobra.Command{}
	in := strings.NewReader("1\n")
	out := &bytes.Buffer{}
	cmd.SetIn(in)
	cmd.SetOut(out)

	promptPackages(cmd)

	output := out.String()

	// Check that options are numbered
	if !strings.Contains(output, "1)") {
		t.Error("output should contain numbered option 1)")
	}
	if !strings.Contains(output, "2)") {
		t.Error("output should contain numbered option 2)")
	}

	// Check that descriptions are present
	if !strings.Contains(output, "Automatically discover") {
		t.Error("output should describe auto mode")
	}
	if !strings.Contains(output, "Specify package paths") {
		t.Error("output should describe manual mode")
	}
}
