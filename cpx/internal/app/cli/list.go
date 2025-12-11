package cli

import (
	"fmt"

	"github.com/ozacod/cpx/internal/pkg/vcpkg"
	"github.com/spf13/cobra"
)

// ListCmd creates the list command
func ListCmd(client *vcpkg.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List available libraries",
		Long:  "List available libraries. Passes through to vcpkg list command.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(cmd, args, client)
		},
		Args: cobra.ArbitraryArgs,
	}

	return cmd
}

func runList(_ *cobra.Command, args []string, client *vcpkg.Client) error {
	// Directly pass all arguments to vcpkg list command
	vcpkgArgs := []string{"list"}
	vcpkgArgs = append(vcpkgArgs, args...)

	if client == nil {
		return fmt.Errorf("vcpkg client not initialized")
	}
	return client.RunCommand(vcpkgArgs)
}
