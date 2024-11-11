package api

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServerStart(t *testing.T) {
	err := ServerStart(&Server{})
	require.Error(t, err)
}
