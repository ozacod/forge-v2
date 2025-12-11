package cli

import (
	"fmt"

	"github.com/ozacod/cpx/internal/pkg/vcpkg"
	"github.com/spf13/cobra"
)

// RemoveCmd creates the remove command
func RemoveCmd(client *vcpkg.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove",
		Short:   "Remove a dependency",
		Long:    "Remove a dependency. Passes through to vcpkg remove command.",
		Aliases: []string{"rm"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRemove(cmd, args, client)
		},
		Args: cobra.MinimumNArgs(1),
	}

	return cmd
}

func runRemove(_ *cobra.Command, args []string, client *vcpkg.Client) error {
	// Directly pass all arguments to vcpkg remove command
	// cpx remove <args> -> vcpkg remove <args>
	vcpkgArgs := []string{"remove"}
	vcpkgArgs = append(vcpkgArgs, args...)

	if client == nil {
		return fmt.Errorf("vcpkg client not initialized")
	}
	return client.RunCommand(vcpkgArgs)
}
