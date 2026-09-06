package builders

import (
	"regexp"

	"github.com/rios0rios0/langforge/pkg/infrastructure/languages/toolchain"
	"github.com/rios0rios0/langforge/pkg/support/cmdexec"
	testkit "github.com/rios0rios0/testkit/pkg/test"
)

// SDKBuilder builds toolchain.SDK descriptions using the builder pattern.
type SDKBuilder struct {
	*testkit.BaseBuilder

	name           string
	versionManager string
	installCommand string
	startCommand   string
	stopCommand    string
	versionCommand cmdexec.CommandSpec
	versionPattern *regexp.Regexp
}

// NewSDKBuilder creates a new builder describing a fictional SDK whose `tool
// --version` output reads `tool 1.2.3 (stable)`.
func NewSDKBuilder() *SDKBuilder {
	return &SDKBuilder{
		BaseBuilder:    testkit.NewBaseBuilder(),
		name:           "Tool",
		versionManager: "tvm",
		installCommand: "tvm install %s",
		startCommand:   "tool start",
		versionCommand: cmdexec.CommandSpec{Name: "tool", Args: []string{"--version"}},
		versionPattern: regexp.MustCompile(`tool\s+(\d+\.\d+\.\d+)`),
	}
}

func (b *SDKBuilder) WithName(name string) *SDKBuilder {
	b.name = name
	return b
}

func (b *SDKBuilder) WithVersionManager(versionManager string) *SDKBuilder {
	b.versionManager = versionManager
	return b
}

func (b *SDKBuilder) WithInstallCommand(installCommand string) *SDKBuilder {
	b.installCommand = installCommand
	return b
}

func (b *SDKBuilder) WithStartCommand(startCommand string) *SDKBuilder {
	b.startCommand = startCommand
	return b
}

func (b *SDKBuilder) WithStopCommand(stopCommand string) *SDKBuilder {
	b.stopCommand = stopCommand
	return b
}

func (b *SDKBuilder) WithVersionCommand(name string, args ...string) *SDKBuilder {
	b.versionCommand = cmdexec.CommandSpec{Name: name, Args: args}
	return b
}

func (b *SDKBuilder) WithVersionPattern(pattern *regexp.Regexp) *SDKBuilder {
	b.versionPattern = pattern
	return b
}

func (b *SDKBuilder) Build() any {
	return toolchain.SDK{
		Name:           b.name,
		VersionManager: b.versionManager,
		InstallCommand: b.installCommand,
		StartCommand:   b.startCommand,
		StopCommand:    b.stopCommand,
		VersionCommand: b.versionCommand,
		VersionPattern: b.versionPattern,
	}
}
