package boot

import (
	"github.com/kenpusney/cra/core/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewContext(t *testing.T) {
	context := NewContext(&common.Opts{})

	assert.NotNil(t, context)
}
