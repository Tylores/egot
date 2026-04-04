package testhelpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AssertHelper wraps common assertion patterns
type AssertHelper struct {
	t *testing.T
}

// New creates a new assertion helper
func New(t *testing.T) *AssertHelper {
	return &AssertHelper{t: t}
}

// NoError asserts that an error is nil
func (h *AssertHelper) NoError(err error) {
	require.NoError(h.t, err)
}

// Error asserts that an error occurred
func (h *AssertHelper) Error(err error) {
	require.Error(h.t, err)
}

// Equal asserts two values are equal
func (h *AssertHelper) Equal(expected, actual interface{}, msgAndArgs ...interface{}) {
	assert.Equal(h.t, expected, actual, msgAndArgs...)
}

// NotNil asserts a value is not nil
func (h *AssertHelper) NotNil(obj interface{}, msgAndArgs ...interface{}) {
	assert.NotNil(h.t, obj, msgAndArgs...)
}

// Nil asserts a value is nil
func (h *AssertHelper) Nil(obj interface{}, msgAndArgs ...interface{}) {
	assert.Nil(h.t, obj, msgAndArgs...)
}

// True asserts a condition is true
func (h *AssertHelper) True(value bool, msgAndArgs ...interface{}) {
	assert.True(h.t, value, msgAndArgs...)
}

// False asserts a condition is false
func (h *AssertHelper) False(value bool, msgAndArgs ...interface{}) {
	assert.False(h.t, value, msgAndArgs...)
}

// Contains asserts that a string contains a substring
func (h *AssertHelper) Contains(s, contains string, msgAndArgs ...interface{}) {
	assert.Contains(h.t, s, contains, msgAndArgs...)
}
