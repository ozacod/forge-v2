package cli

import (
	"github.com/ozacod/cpx/internal/pkg/quality"
	"github.com/ozacod/cpx/internal/pkg/vcpkg"
	"github.com/spf13/cobra"
)

// AnalyzeCmd creates the analyze command
func AnalyzeCmd(client *vcpkg.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "analyze",
		Short: "Run comprehensive code analysis and generate HTML report",
		Long:  "Run comprehensive code analysis using cppcheck, clang-tidy, and flawfinder. Generates a combined HTML report (analyze.html).",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAnalyze(cmd, args, client)
		},
		Args: cobra.ArbitraryArgs,
	}

	cmd.Flags().String("output", "analyze.html", "Output HTML file path")
	cmd.Flags().Bool("skip-cppcheck", false, "Skip Cppcheck analysis")
	cmd.Flags().Bool("skip-lint", false, "Skip clang-tidy analysis")
	cmd.Flags().Bool("skip-flawfinder", false, "Skip Flawfinder analysis")

	return cmd
}

func runAnalyze(cmd *cobra.Command, args []string, client *vcpkg.Client) error {
	output, _ := cmd.Flags().GetString("output")
	skipCppcheck, _ := cmd.Flags().GetBool("skip-cppcheck")
	skipLint, _ := cmd.Flags().GetBool("skip-lint")
	skipFlawfinder, _ := cmd.Flags().GetBool("skip-flawfinder")

	// Get remaining args as target directories (default to current directory)
	targets := args
	if len(targets) == 0 {
		targets = []string{"."}
	}

	return quality.RunComprehensiveAnalysis(output, skipCppcheck, skipLint, skipFlawfinder, targets, client)
}
