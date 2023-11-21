package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"log/slog"
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

	logLevel  slog.Level
	logFormat string
)

func init() {
	fs = flag.NewFlagSet(AppName, flag.ExitOnError)

	fs.BoolVar(&showVersion, "v", false, "Print version and exit")
	fs.BoolVar(&showBuildInfo, "V", false, "Print build information and exit")

	fs.StringVar(&logFormat, "log.format", getEnv("APP_LOG_FORMAT", "text"), "Log format (text, json)")
	fs.TextVar(&logLevel, "log.level", getLogLevelEnv("APP_LOG_LEVEL", slog.LevelInfo), "Log level (debug, info, warn, error)")
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

	logger := initLogger()

	logger.Info("started")
	defer logger.Info("stopped")
}

func getEnv(key string, defaultValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultValue
}

func getLogLevelEnv(key string, defaultLevel slog.Level) slog.Level {
	lvl := defaultLevel
	if val, ok := os.LookupEnv(key); ok {
		if err := lvl.UnmarshalText([]byte(val)); err != nil {
			return defaultLevel
		}
	}
	return lvl
}

func initLogger() *slog.Logger {
	opts := &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.TimeKey:
				// use UTC time
				a.Value = slog.TimeValue(a.Value.Time().UTC())

			case slog.SourceKey:
				// remove directories from File
				source := a.Value.Any().(*slog.Source)
				source.File = filepath.Base(source.File)
			}
			return a
		},
	}

	var h slog.Handler = slog.NewTextHandler(os.Stdout, opts)
	if logFormat == "json" {
		h = slog.NewJSONHandler(os.Stdout, opts)
	}
	return slog.New(h)
}
