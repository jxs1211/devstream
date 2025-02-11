package zentao

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/goclient"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
)

func Update(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	// Initialize Operator with Operations
	operator := &plugininstaller.Operator{
		PreExecuteOperations: plugininstaller.PreExecuteOperations{
			goclient.Validate,
		},
		ExecuteOperations: plugininstaller.ExecuteOperations{
			goclient.DeleteApp,
			goclient.CreateDeploymentWrapperLabelAndContainerPorts(defaultZentaolabels, &defaultZentaoPorts, defaultName),
			goclient.CreateServiceWrapperLabelAndPorts(defaultZentaolabels, &defaultSVCPort),
			goclient.WaitForReady(retryTimes),
		},
		TerminateOperations: plugininstaller.TerminateOperations{
			goclient.DealWithErrWhenInstall,
		},
		GetStatusOperation: goclient.GetStatus,
	}

	// Execute all Operations in Operator
	status, err := operator.Execute(options)
	if err != nil {
		return nil, err
	}
	return status, nil
}
