package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Note: Project configuration is done through TUI only (cpx new).
// There is no cpx.yaml file - all settings come from the interactive prompts.

// ProjectConfig represents the project configuration structure
type ProjectConfig struct {
	Package struct {
		Name        string   `yaml:"name"`
		Version     string   `yaml:"version"`
		CppStandard int      `yaml:"cpp_standard"`
		Authors     []string `yaml:"authors,omitempty"`
		Description string   `yaml:"description,omitempty"`
	} `yaml:"package"`
	Build struct {
		SharedLibs  bool   `yaml:"shared_libs"`
		ClangFormat string `yaml:"clang_format"`
		BuildType   string `yaml:"build_type,omitempty"`
		CxxFlags    string `yaml:"cxx_flags,omitempty"`
	} `yaml:"build"`
	VCS struct {
		Type string `yaml:"type,omitempty"` // "git" or "none"
	} `yaml:"vcs,omitempty"`
	PackageManager struct {
		Type string `yaml:"type,omitempty"` // "vcpkg" or "none"
	} `yaml:"package_manager,omitempty"`
	Testing struct {
		Framework string `yaml:"framework"`
	} `yaml:"testing"`
	Hooks struct {
		PreCommit []string `yaml:"precommit,omitempty"` // e.g., ["fmt", "lint"]
		PrePush   []string `yaml:"prepush,omitempty"`   // e.g., ["test"]
	} `yaml:"hooks,omitempty"`
}

// CIConfig represents the cpx.ci structure for cross-compilation
type CIConfig struct {
	Targets []CITarget `yaml:"targets"`
	Build   CIBuild    `yaml:"build"`
	Output  string     `yaml:"output"`
}

// CITarget represents a cross-compilation target
type CITarget struct {
	Name       string `yaml:"name"`
	Dockerfile string `yaml:"dockerfile"`
	Image      string `yaml:"image"`
	Triplet    string `yaml:"triplet"`
	Platform   string `yaml:"platform"`
}

// CIBuild represents CI build configuration
type CIBuild struct {
	Type         string   `yaml:"type"`
	Optimization string   `yaml:"optimization"`
	Jobs         int      `yaml:"jobs"`
	CMakeArgs    []string `yaml:"cmake_args"`
	BuildArgs    []string `yaml:"build_args"`
}

// LoadCI loads the CI configuration from cpx.ci
func LoadCI(path string) (*CIConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config CIConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse cpx.ci: %w", err)
	}

	// Set defaults
	if config.Output == "" {
		config.Output = "out"
	}
	if config.Build.Type == "" {
		config.Build.Type = "Release"
	}
	if config.Build.Optimization == "" {
		config.Build.Optimization = "2"
	}

	return &config, nil
}
