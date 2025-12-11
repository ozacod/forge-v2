package cli

import (
	"fmt"

	"github.com/ozacod/cpx/internal/app/cli/tui"
	"github.com/ozacod/cpx/internal/pkg/vcpkg"
	"github.com/spf13/cobra"
)

// SearchCmd creates the search command
func SearchCmd(client *vcpkg.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search [query]",
		Short: "Search for libraries interactively",
		Long:  "Search for libraries using an interactive TUI. Select packages to add them to your project.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSearch(cmd, args, client)
		},
		Args: cobra.MaximumNArgs(1),
	}

	return cmd
}

func runSearch(_ *cobra.Command, args []string, client *vcpkg.Client) error {
	query := ""
	if len(args) > 0 {
		query = args[0]
	}

	if err := requireVcpkgProject("cpx search"); err != nil {
		return err
	}

	if client == nil {
		return fmt.Errorf("vcpkg client not initialized")
	}

	vcpkgPath, err := client.GetPath()
	if err != nil {
		return fmt.Errorf("failed to get vcpkg path: %w", err)
	}

	return tui.RunSearch(query, vcpkgPath, client.RunCommand)
}
