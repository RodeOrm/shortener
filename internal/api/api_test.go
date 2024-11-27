package api

import (
	"sync"
	"testing"

	"github.com/rodeorm/shortener/internal/core"
	"github.com/stretchr/testify/require"
)

func TestServerStart(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	err := ServerStart(&core.Server{}, &wg)
	require.Error(t, err)
}
