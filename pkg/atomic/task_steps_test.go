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
