package templates

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateModuleBazel(t *testing.T) {
	tests := []struct {
		name               string
		projectName        string
		version            string
		testFramework      string
		benchmarkFramework string
		shouldContain      []string
		shouldNotContain   []string
	}{
		{
			name:               "Basic project",
			projectName:        "myproject",
			version:            "1.0.0",
			testFramework:      "",
			benchmarkFramework: "",
			shouldContain:      []string{"module(", `name = "myproject"`, `version = "1.0.0"`, "rules_cc"},
			shouldNotContain:   []string{"googletest", "catch2", "google_benchmark"},
		},
		{
			name:               "With googletest",
			projectName:        "testproject",
			version:            "0.1.0",
			testFramework:      "googletest",
			benchmarkFramework: "",
			shouldContain:      []string{"googletest"},
			shouldNotContain:   []string{"catch2", "google_benchmark"},
		},
		{
			name:               "With catch2",
			projectName:        "catchproject",
			version:            "0.1.0",
			testFramework:      "catch2",
			benchmarkFramework: "",
			shouldContain:      []string{"catch2"},
			shouldNotContain:   []string{"googletest"},
		},
		{
			name:               "With google benchmark",
			projectName:        "benchproject",
			version:            "0.1.0",
			testFramework:      "",
			benchmarkFramework: "google-benchmark",
			shouldContain:      []string{"google_benchmark"},
			shouldNotContain:   []string{"googletest"},
		},
		{
			name:               "Default version",
			projectName:        "project",
			version:            "",
			testFramework:      "",
			benchmarkFramework: "",
			shouldContain:      []string{`version = "0.1.0"`},
			shouldNotContain:   []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateModuleBazel(tt.projectName, tt.version, tt.testFramework, tt.benchmarkFramework)

			for _, s := range tt.shouldContain {
				assert.Contains(t, result, s, "Expected to contain: %s", s)
			}
			for _, s := range tt.shouldNotContain {
				assert.NotContains(t, result, s, "Expected NOT to contain: %s", s)
			}
		})
	}
}

func TestGenerateBuildBazelRoot(t *testing.T) {
	tests := []struct {
		name          string
		projectName   string
		isExe         bool
		shouldContain []string
	}{
		{
			name:        "Executable project",
			projectName: "myapp",
			isExe:       true,
			shouldContain: []string{
				"alias(",
				`name = "myapp"`,
				`actual = "//src:myapp"`,
				`name = "myapp_lib"`,
			},
		},
		{
			name:        "Library project",
			projectName: "mylib",
			isExe:       false,
			shouldContain: []string{
				"alias(",
				`name = "mylib"`,
				`actual = "//src:mylib"`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateBuildBazelRoot(tt.projectName, tt.isExe)

			for _, s := range tt.shouldContain {
				assert.Contains(t, result, s, "Expected to contain: %s", s)
			}
		})
	}
}

func TestGenerateBuildBazelTests(t *testing.T) {
	tests := []struct {
		name          string
		projectName   string
		testFramework string
		shouldContain []string
	}{
		{
			name:          "GoogleTest",
			projectName:   "myproject",
			testFramework: "googletest",
			shouldContain: []string{"cc_test", "@googletest//:gtest_main", `name = "myproject_test"`},
		},
		{
			name:          "Catch2",
			projectName:   "myproject",
			testFramework: "catch2",
			shouldContain: []string{"cc_test", "@catch2//:catch2_main"},
		},
		{
			name:          "Doctest",
			projectName:   "myproject",
			testFramework: "doctest",
			shouldContain: []string{"cc_test", "@doctest//:doctest"},
		},
		{
			name:          "No framework",
			projectName:   "myproject",
			testFramework: "",
			shouldContain: []string{"cc_test", `"//src:myproject_lib"`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateBuildBazelTests(tt.projectName, tt.testFramework)

			for _, s := range tt.shouldContain {
				assert.Contains(t, result, s, "Expected to contain: %s", s)
			}
		})
	}
}

func TestGenerateBazelrc(t *testing.T) {
	tests := []struct {
		name          string
		cppStandard   int
		shouldContain []string
	}{
		{
			name:          "C++17",
			cppStandard:   17,
			shouldContain: []string{"c++17", "--symlink_prefix=.", "build:release", "build:debug"},
		},
		{
			name:          "C++20",
			cppStandard:   20,
			shouldContain: []string{"c++20"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateBazelrc(tt.cppStandard)

			for _, s := range tt.shouldContain {
				assert.Contains(t, result, s, "Expected to contain: %s", s)
			}
		})
	}
}

func TestGenerateBenchmarkSources(t *testing.T) {
	tests := []struct {
		name               string
		projectName        string
		benchmarkFramework string
		expectedDeps       []string
		shouldBeNil        bool
	}{
		{
			name:               "Google Benchmark",
			projectName:        "myproject",
			benchmarkFramework: "google-benchmark",
			expectedDeps:       []string{"benchmark"},
			shouldBeNil:        false,
		},
		{
			name:               "Nanobench",
			projectName:        "myproject",
			benchmarkFramework: "nanobench",
			expectedDeps:       []string{"nanobench"},
			shouldBeNil:        false,
		},
		{
			name:               "Catch2 Benchmark",
			projectName:        "myproject",
			benchmarkFramework: "catch2-benchmark",
			expectedDeps:       []string{"catch2"},
			shouldBeNil:        false,
		},
		{
			name:               "Unknown framework",
			projectName:        "myproject",
			benchmarkFramework: "unknown",
			expectedDeps:       nil,
			shouldBeNil:        true,
		},
		{
			name:               "Empty framework",
			projectName:        "myproject",
			benchmarkFramework: "",
			expectedDeps:       nil,
			shouldBeNil:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sources, deps := GenerateBenchmarkSources(tt.projectName, tt.benchmarkFramework)

			if tt.shouldBeNil {
				assert.Nil(t, sources)
				assert.Nil(t, deps)
			} else {
				assert.NotNil(t, sources)
				assert.NotEmpty(t, sources.Main)
				assert.Equal(t, tt.expectedDeps, deps)
				// Verify the benchmark includes the project header
				assert.Contains(t, sources.Main, tt.projectName)
			}
		})
	}
}

func TestSafeIdent(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple name",
			input:    "myproject",
			expected: "myproject",
		},
		{
			name:     "Name with hyphens",
			input:    "my-project",
			expected: "my_project",
		},
		{
			name:     "Empty name",
			input:    "",
			expected: "project",
		},
		{
			name:     "Starting with digit",
			input:    "123abc",
			expected: "_123abc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SafeIdent(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGenerateVersionHpp(t *testing.T) {
	result := GenerateVersionHpp("myproject", "1.2.3")

	// Should contain version components
	assert.True(t, strings.Contains(result, "1") && strings.Contains(result, "2") && strings.Contains(result, "3"))
	// Should contain project name in include guard
	assert.Contains(t, result, "MYPROJECT")
	// Should be valid C++ header
	assert.Contains(t, result, "#ifndef")
	assert.Contains(t, result, "#define")
	assert.Contains(t, result, "#endif")
}

func TestGenerateGitignore(t *testing.T) {
	result := GenerateGitignore()

	// Should ignore common build artifacts
	assert.Contains(t, result, "build")
	assert.Contains(t, result, ".cmake")
	// Should ignore IDE files
	assert.Contains(t, result, ".vscode")
	assert.Contains(t, result, ".idea")
}

func TestGenerateBazelGitignore(t *testing.T) {
	result := GenerateBazelGitignore()

	// Should ignore Bazel artifacts
	assert.Contains(t, result, "bazel-")
	assert.Contains(t, result, ".bazel-")
	// Should ignore build directory
	assert.Contains(t, result, "build")
}
