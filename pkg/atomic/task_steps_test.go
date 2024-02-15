package atomic

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecuteCommandStep(t *testing.T) {
	step, err := NewExecuteCommandStep("whoami", "bash")
	assert.Nil(t, err)

	ctx := context.Background()
	r := step.Run(ctx)
	assert.Nil(t, r.Error)
	assert.NotNil(t, r.Data)

	result := r.Data.(*ExecuteCommandStepResult)
	assert.Equal(t, *result.Process.ExitCode, 0)
}

func TestListProcessesStep(t *testing.T) {
	step, err := NewListProcessesStep()
	assert.Nil(t, err)

	ctx := context.Background()
	r := step.Run(ctx)
	assert.Nil(t, r.Error)
	assert.NotNil(t, r.Data)

	result := r.Data.(*ListProcessesStepResult)
	assert.True(t, len(result.Processes) > 0)
}

func TestListUsersStep(t *testing.T) {
	step, err := NewListUsersStep()
	assert.Nil(t, err)

	ctx := context.Background()
	r := step.Run(ctx)
	assert.Nil(t, r.Error)
	assert.NotNil(t, r.Data)

	result := r.Data.(*ListUsersStepResult)
	assert.True(t, len(result.Users) > 0)
}

func TestListNetworkConnectionsStep(t *testing.T) {
	step, err := NewListNetworkConnectionsStep()
	assert.Nil(t, err)

	ctx := context.Background()
	r := step.Run(ctx)
	assert.Nil(t, r.Error)
	assert.NotNil(t, r.Data)

	result := r.Data.(*ListNetworkConnectionsStepResult)
	assert.True(t, len(result.Connections) > 0)
}

func TestListNetworkInterfacesStep(t *testing.T) {
	step, err := NewListNetworkInterfacesStep()
	assert.Nil(t, err)

	ctx := context.Background()
	r := step.Run(ctx)
	assert.Nil(t, r.Error)
	assert.NotNil(t, r.Data)

	result := r.Data.(*ListNetworkInterfacesStepResult)
	assert.True(t, len(result.Interfaces) > 0)
}
