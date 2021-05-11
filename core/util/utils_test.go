package util

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestUnmarshall(t *testing.T) {
	simpleObject := ConvertJsonBodyToObject(strings.NewReader("{}"))

	assert.NotNil(t, simpleObject)

	simpleArray := ConvertJsonBodyToObject(strings.NewReader("[1, 2, 3]"))

	assert.NotNil(t, simpleArray)

	malformed := ConvertJsonBodyToObject(strings.NewReader("{1}"))

	assert.Nil(t, malformed)
}

func TestGenerateId(t *testing.T) {
	assert.Equal(t, GenerateId("abc", "def", -1), "abc")
	assert.Equal(t, GenerateId("", "def", -1), "def")
	assert.Equal(t, GenerateId("abc", "def", 1), "abc-1")
	assert.Equal(t, GenerateId("", "def", 1), "def-1")
}
