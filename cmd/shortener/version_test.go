package main

import (
	"testing"
)

func TestPrintBuildInfo(t *testing.T) {
	tests := []struct {
		name    string
		value   int
		wantErr bool
	}{

		{
			name: "Проверка печати конфигурации",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			printBuildInfo()
		})
	}
}

func TestPrintVersionField(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{

		{
			name: "Проверка обработки пустого значения", value: "",
		},
		{
			name: "Проверка обработки нормального значения", value: "Нормальное значение",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			printVersionField(tt.name, tt.value)
		})
	}
}
