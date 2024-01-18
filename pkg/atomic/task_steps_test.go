package atomic

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecStepWithoutArgs(t *testing.T) {
	step, err := NewExecStep("whoami", "bash")
	assert.Nil(t, err)

	ctx := context.Background()
	r := step.Exec(ctx)
	assert.Nil(t, r.Error)
	assert.NotNil(t, r.Data)

	result := r.Data.(*ExecStepResult)
	assert.Equal(t, *result.Subprocess.ExitCode, 0)
}
