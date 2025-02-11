package golang

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/github"
)

// Delete remove GitHub Actions workflows.
func Delete(options configmanager.RawOptions) (bool, error) {
	operator := &plugininstaller.Operator{
		PreExecuteOperations: plugininstaller.PreExecuteOperations{
			validate,
			github.BuildWorkFlowsWrapper(workflows),
		},
		ExecuteOperations: plugininstaller.ExecuteOperations{
			deleteDockerHubInfoForPush,
			github.ProcessAction(github.ActionDelete),
		},
	}

	_, err := operator.Execute(options)
	if err != nil {
		return false, err
	}
	return true, nil
}
