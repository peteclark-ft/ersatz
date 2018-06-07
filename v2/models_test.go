package v2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionIsCorrect(t *testing.T) {
	f := Fixtures{}
	assert.Equal(t, 2, f.Version())
}
