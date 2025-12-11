package cli

import (
	"fmt"

	"github.com/ozacod/cpx/internal/pkg/vcpkg"
	"github.com/spf13/cobra"
)

// InfoCmd creates the info command
func InfoCmd(client *vcpkg.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info",
		Short: "Show detailed library information",
		Long:  "Show detailed library information. Passes through to vcpkg show command.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInfo(cmd, args, client)
		},
		Args: cobra.MinimumNArgs(1),
	}

	return cmd
}

func runInfo(_ *cobra.Command, args []string, client *vcpkg.Client) error {
	// Directly pass all arguments to vcpkg show command
	// cpx info <package> -> vcpkg show <package>
	vcpkgArgs := []string{"show"}
	vcpkgArgs = append(vcpkgArgs, args...)

	if client == nil {
		return fmt.Errorf("vcpkg client not initialized")
	}
	return client.RunCommand(vcpkgArgs)
}
