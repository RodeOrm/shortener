package main

import (
	"testing"

	"golang.org/x/tools/go/analysis"
)

func TestAddStatisticCheckAnalyzersOther(t *testing.T) {
	analyzers := []*analysis.Analyzer{}
	analyzers = addStatisticCheckAnalyzersOther(analyzers)

	expectedNames := map[string]bool{"S1000": true, "ST1000": true, "QF1001": true}
	for _, a := range analyzers {
		if !expectedNames[a.Name] {
			t.Errorf("Unexpected analyzer: %s", a.Name)
		}
	}
}

func TestAddStatisticCheckAnalyzersSA(t *testing.T) {
	analyzers := []*analysis.Analyzer{}
	analyzers = addStatisticCheckAnalyzersSA(analyzers)

	if len(analyzers) == 0 {
		t.Errorf("не удалось добавить анализаторы")
	}
}

func TestAddStandardAnalyzers(t *testing.T) {
	analyzers := []*analysis.Analyzer{}
	analyzers = addStandardAnalyzers(analyzers)

	if len(analyzers) == 0 {
		t.Errorf("не удалось добавить анализаторы")
	}
}
