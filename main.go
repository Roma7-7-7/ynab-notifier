package main

import (
	"flag"
	stdLog "log"

	"go.uber.org/zap"
)

var verbose = flag.Bool("debug", false, "debug output")

func main() {
	flag.Parse()

	var err error
	var l *zap.Logger
	if *verbose {
		l, err = zap.NewDevelopment()
	} else {
		l, err = zap.NewProduction()
	}
	if err != nil {
		stdLog.Fatalf("can't initialize zap logger: %v", err)
	}
	log := l.Sugar()

	log.Infow("hello world", "foo", "bar")
}
