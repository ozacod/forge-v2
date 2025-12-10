package cli

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/ozacod/cpx/internal/pkg/bazel"
	"github.com/ozacod/cpx/pkg/config"
	"github.com/spf13/cobra"
)

var addRunVcpkgCommandFunc func([]string) error
var addGetBcrPathFunc func() string

// AddCmd creates the add command
func AddCmd(runVcpkgCommand func([]string) error, getBcrPath func() string) *cobra.Command {
	addRunVcpkgCommandFunc = runVcpkgCommand
	addGetBcrPathFunc = getBcrPath

	cmd := &cobra.Command{
		Use:   "add [package]",
		Short: "Add a dependency",
		Long: `Add a dependency to your project.

For vcpkg projects: passes through to 'vcpkg add port' and prints usage info.
For Bazel projects: fetches the latest version from BCR and updates MODULE.bazel.`,
		RunE: runAdd,
		Args: cobra.MinimumNArgs(1),
	}

	return cmd
}

func runAdd(cmd *cobra.Command, args []string) error {
	projectType, err := RequireProject("cpx add")
	if err != nil {
		return err
	}

	switch projectType {
	case ProjectTypeVcpkg:
		return runVcpkgAdd(args)
	case ProjectTypeBazel:
		return runBazelAdd(args)
	case ProjectTypeMeson:
		return runMesonAdd(args)
	default:
		return fmt.Errorf("unsupported project type")
	}
}

func runVcpkgAdd(args []string) error {
	// Directly pass all arguments to vcpkg add command
	// cpx add <pkg> -> vcpkg add port <pkg>
	vcpkgArgs := []string{"add", "port"}
	vcpkgArgs = append(vcpkgArgs, args...)

	if err := addRunVcpkgCommandFunc(vcpkgArgs); err != nil {
		return err
	}

	// Print usage info for the first package
	if len(args) > 0 {
		pkgName := args[0]
		if !strings.HasPrefix(pkgName, "-") {
			printVcpkgUsageInfo(pkgName)
		}
	}

	return nil
}

func runBazelAdd(args []string) error {
	bcrPath := addGetBcrPathFunc()
	if bcrPath == "" {
		return fmt.Errorf("bazel Central Registry not configured\n  hint: run 'cpx config set-bcr-root <path>' or reinstall cpx")
	}

	client := bazel.NewClient(bcrPath)

	for _, pkgName := range args {
		if strings.HasPrefix(pkgName, "-") {
			continue
		}

		// Get latest version
		version, err := client.GetLatestVersion(pkgName)
		if err != nil {
			fmt.Printf("%s✗ Module '%s' not found in BCR%s\n", Red, pkgName, Reset)
			continue
		}

		// Add to MODULE.bazel
		if err := bazel.AddDependency("MODULE.bazel", pkgName, version); err != nil {
			return fmt.Errorf("failed to add dependency: %w", err)
		}

		fmt.Printf("%s✓ Added %s@%s to MODULE.bazel%s\n", Green, pkgName, version, Reset)
		printBazelUsageInfo(pkgName, version)
	}

	return nil
}

func runMesonAdd(args []string) error {
	// Load config for local WrapDB path
	cfg, cfgErr := config.LoadGlobal()
	hasLocalWrapdb := cfgErr == nil && cfg.WrapdbRoot != ""

	// Meson uses WrapDB - copy from local or download .wrap files to subprojects/
	for _, pkgName := range args {
		if strings.HasPrefix(pkgName, "-") {
			continue
		}

		fmt.Printf("%sAdding wrap file for %s...%s\n", Cyan, pkgName, Reset)

		// Ensure subprojects directory exists
		if err := createDirIfNotExists("subprojects"); err != nil {
			return fmt.Errorf("failed to create subprojects directory: %w", err)
		}

		wrapPath := fmt.Sprintf("subprojects/%s.wrap", pkgName)

		// Try local WrapDB first
		if hasLocalWrapdb {
			localWrap := cfg.WrapdbRoot + "/" + pkgName + ".wrap"
			if content, err := os.ReadFile(localWrap); err == nil {
				if err := writeFile(wrapPath, content); err != nil {
					return fmt.Errorf("failed to write wrap file: %w", err)
				}
				fmt.Printf("%s✓ Added %s (from local cache)%s\n", Green, pkgName, Reset)
				printMesonUsageInfo(pkgName)
				continue
			}
		}

		// Fallback: download from WrapDB
		wrapURL := fmt.Sprintf("https://wrapdb.mesonbuild.com/v2/%s.wrap", pkgName)
		resp, err := http.Get(wrapURL)
		if err != nil || resp.StatusCode != 200 {
			fmt.Printf("%s✗ Package '%s' not found in WrapDB%s\n", Red, pkgName, Reset)
			fmt.Printf("  Try searching at: https://wrapdb.mesonbuild.com/\n")
			continue
		}
		defer resp.Body.Close()

		wrapContent, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("%s✗ Failed to download wrap for %s%s\n", Red, pkgName, Reset)
			continue
		}

		if err := writeFile(wrapPath, wrapContent); err != nil {
			return fmt.Errorf("failed to write wrap file: %w", err)
		}

		fmt.Printf("%s✓ Added %s to subprojects/%s.wrap%s\n", Green, pkgName, pkgName, Reset)
		printMesonUsageInfo(pkgName)
	}

	return nil
}

// printVcpkgUsageInfo fetches and prints usage info from GitHub for vcpkg packages
func printVcpkgUsageInfo(pkgName string) {
	resp, err := http.Get(fmt.Sprintf("https://raw.githubusercontent.com/microsoft/vcpkg/master/ports/%s/usage", pkgName))
	if err != nil || resp.StatusCode != 200 {
		return
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	content := strings.TrimSpace(string(bytes))
	if content != "" {
		fmt.Printf("\n%sUSAGE INFO FOR %s:%s\n", Cyan, pkgName, Reset)
		fmt.Println(content)
		fmt.Println()
	}

	// Print link to cpx website for more info
	fmt.Printf("%s📦 Find sample usage and more info at:%s\n", Cyan, Reset)
	fmt.Printf("   https://cpx-dev.vercel.app/packages#package/%s\n\n", pkgName)
}

// printBazelUsageInfo prints usage info for Bazel modules
func printBazelUsageInfo(pkgName, version string) {
	fmt.Printf("\n%sUSAGE INFO FOR %s:%s\n", Cyan, pkgName, Reset)
	fmt.Printf("Add this to your BUILD.bazel:\n\n")
	fmt.Printf("  deps = [\"@%s//:<target>\"]\n\n", pkgName)
	fmt.Printf("%s📦 Find more info at:%s\n", Cyan, Reset)
	fmt.Printf("   https://registry.bazel.build/modules/%s\n", pkgName)
	fmt.Printf("   https://cpx-dev.vercel.app/bazel#module/%s\n\n", pkgName)
}

// printMesonUsageInfo prints usage info for Meson wraps
func printMesonUsageInfo(pkgName string) {
	fmt.Printf("\n%sUSAGE INFO FOR %s:%s\n", Cyan, pkgName, Reset)
	fmt.Printf("Add this to your meson.build:\n\n")
	fmt.Printf("  %s_dep = dependency('%s')\n\n", pkgName, pkgName)
	fmt.Printf("%s📦 Find more info at:%s\n", Cyan, Reset)
	fmt.Printf("   https://wrapdb.mesonbuild.com/\n\n")
}

func createDirIfNotExists(path string) error {
	return os.MkdirAll(path, 0755)
}

func writeFile(path string, content []byte) error {
	return os.WriteFile(path, content, 0644)
}
