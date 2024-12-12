package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetContextValue(t *testing.T) {
	// Test case 1: Value is set in the context

	t.Run("value is not set in the context", func(t *testing.T) {
		// Set value in the current goroutine's local storage

		// Call the function
		result := GetContextValue("myKey")

		// Assert that the returned value is correct
		assert.Equal(t, "", result)
	})
}
