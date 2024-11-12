package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{

		{
			name: "Проверка нормальной конфигурации c БД",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, err := configurate()
			require.NoError(t, err)
			err = profile(server.ProfileType)
			require.NoError(t, err)
		})
	}
}
