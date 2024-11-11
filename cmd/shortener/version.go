package main

import "fmt"

var buildVersion, buildDate, buildCommit string

func printBuildInfo() {
	printVersionField("Build version", buildVersion)
	printVersionField("Build date", buildDate)
	printVersionField("Build commit", buildCommit)
}

func printVersionField(name, value string) {
	if value != "" {
		fmt.Printf("%s: %s\n", name, value)
	} else {
		fmt.Printf("%s: N/A\n", name)
	}
}
