package devlake

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/helm"
)

func Delete(options configmanager.RawOptions) (bool, error) {
	// Initialize Operator with Operations
	operator := &plugininstaller.Operator{
		PreExecuteOperations: plugininstaller.PreExecuteOperations{
			helm.SetDefaultConfig(&defaultHelmConfig),
			helm.Validate,
		},
		ExecuteOperations: helm.DefaultDeleteOperations,
	}

	// Execute all Operations in Operator
	_, err := operator.Execute(options)
	if err != nil {
		return false, err
	}

	return true, nil
}
