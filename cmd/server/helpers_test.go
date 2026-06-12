package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func requireString(t *testing.T, values map[string]interface{}, key string) string {
	t.Helper()

	value, ok := values[key].(string)
	require.Truef(t, ok, "expected %s to be string, got %T", key, values[key])
	require.NotEmpty(t, value)
	return value
}

func stringSlice(t *testing.T, value interface{}) []string {
	t.Helper()

	raw, ok := value.([]interface{})
	require.Truef(t, ok, "expected string slice, got %T", value)
	result := make([]string, 0, len(raw))
	for _, item := range raw {
		str, ok := item.(string)
		require.Truef(t, ok, "expected string item, got %T", item)
		result = append(result, str)
	}
	return result
}
