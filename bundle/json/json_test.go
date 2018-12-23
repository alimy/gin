package json

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEnableJsonDecoderUseNumber(t *testing.T) {
	assert.False(t, EnableDecoderUseNumber)
	EnableJsonDecoderUseNumber()
	assert.True(t, EnableDecoderUseNumber)
}