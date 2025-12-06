package cli

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

// NewReleaseCmd creates the release command
func NewReleaseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "release",
		Short: "Bump version number",
		Long:  "Bump version number (major, minor, or patch) in CMakeLists.txt. Defaults to patch if not specified.",
		RunE:  runRelease,
		Args:  cobra.MaximumNArgs(1),
	}

	return cmd
}

func runRelease(cmd *cobra.Command, args []string) error {
	bumpType := "patch"
	if len(args) > 0 {
		bumpType = args[0]
	}
	return bumpVersion(bumpType)
}

// Release is kept for backward compatibility (if needed)
func Release(args []string) {
	// This function is deprecated - use NewReleaseCmd instead
	// Kept for compatibility during migration
}

func bumpVersion(bumpType string) error {
	// Read version from CMakeLists.txt
	cmakeContent, err := os.ReadFile("CMakeLists.txt")
	if err != nil {
		return fmt.Errorf("failed to read CMakeLists.txt: %w", err)
	}

	// Find VERSION in project() declaration
	versionRegex := regexp.MustCompile(`project\s*\(\s*\w+\s+VERSION\s+(\d+\.\d+\.\d+)`)
	matches := versionRegex.FindStringSubmatch(string(cmakeContent))

	var version string
	if len(matches) > 1 {
		version = matches[1]
	} else {
		return fmt.Errorf("could not find VERSION in CMakeLists.txt project() declaration")
	}

	// Parse version
	parts := strings.Split(strings.TrimPrefix(version, "v"), ".")
	if len(parts) < 3 {
		parts = append(parts, make([]string, 3-len(parts))...)
	}

	major, minor, patch := 0, 0, 0
	fmt.Sscanf(parts[0], "%d", &major)
	fmt.Sscanf(parts[1], "%d", &minor)
	fmt.Sscanf(parts[2], "%d", &patch)

	switch bumpType {
	case "major":
		major++
		minor = 0
		patch = 0
	case "minor":
		minor++
		patch = 0
	case "patch":
		patch++
	default:
		return fmt.Errorf("invalid bump type: %s (use major, minor, or patch)", bumpType)
	}

	newVersion := fmt.Sprintf("%d.%d.%d", major, minor, patch)

	fmt.Printf("%s Bumping version: %s → %s%s\n", Cyan, version, newVersion, Reset)

	// Replace version in CMakeLists.txt
	newContent := versionRegex.ReplaceAllStringFunc(string(cmakeContent), func(match string) string {
		return strings.Replace(match, version, newVersion, 1)
	})

	if err := os.WriteFile("CMakeLists.txt", []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write CMakeLists.txt: %w", err)
	}

	fmt.Printf("%s Version updated to %s in CMakeLists.txt%s\n", Green, newVersion, Reset)
	return nil
}
