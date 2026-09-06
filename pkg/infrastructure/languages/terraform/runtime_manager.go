package terraform

import (
	"regexp"

	"github.com/rios0rios0/langforge/pkg/infrastructure/languages/toolchain"
	"github.com/rios0rios0/langforge/pkg/support/cmdexec"
)

// terraformCLIVersionRe matches the version the installed binary reports, whose
// first line reads `Terraform v1.7.0`. It is deliberately distinct from
// terraformVersionRe, which reads the required_version constraint a repository
// declares.
var terraformCLIVersionRe = regexp.MustCompile(`Terraform\s+v(\d+\.\d+(?:\.\d+)?)`)

// RuntimeManager provides SDK and runtime information for Terraform projects.
type RuntimeManager struct {
	*toolchain.RuntimeManager
}

// NewRuntimeManager creates a RuntimeManager with the given runner.
func NewRuntimeManager(runner cmdexec.Runner) *RuntimeManager {
	return &RuntimeManager{toolchain.NewRuntimeManager(runner, terraformSDK())}
}

// terraformSDK describes the Terraform CLI: tfenv installs it, `terraform apply`
// stands in for starting a project, there is no standard stop command, and
// `terraform version` reports the version.
func terraformSDK() toolchain.SDK {
	return toolchain.SDK{
		Name:           "Terraform",
		VersionManager: "tfenv",
		InstallCommand: "tfenv install %s",
		StartCommand:   "terraform apply",
		VersionCommand: cmdexec.CommandSpec{Name: terraformCLI, Args: []string{"version"}},
		VersionPattern: terraformCLIVersionRe,
	}
}
