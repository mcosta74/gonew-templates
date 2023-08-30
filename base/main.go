package main

import (
	"flag"
	"fmt"
	"os"
)

const AppName = "go-template"

// build info
var (
	version = "development"
	commit  = "N.A."
	date    = "N.A."
)

// flags
var (
	fs *flag.FlagSet

	showVersion   bool
	showBuildInfo bool
)

func init() {
	fs = flag.NewFlagSet(AppName, flag.ExitOnError)

	fs.BoolVar(&showVersion, "v", false, "Print version and exit")
	fs.BoolVar(&showBuildInfo, "V", false, "Print build information and exit")
}

func main() {
	fs.Parse(os.Args[1:])

	if showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	if showBuildInfo {
		fmt.Printf("Version:%s, GitCommit:%s, BuildDate:%s\n", version, commit, date)
		os.Exit(0)
	}
	fmt.Println("Hello")
	defer fmt.Println("Bye")
}
