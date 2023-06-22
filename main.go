package main

import (
	"flag"
	stdLog "log"

	"go.uber.org/zap"
)

func main() {
	var debug = flag.Bool("debug", false, "debug mode")

	flag.Parse()

	var err error
	var l *zap.Logger
	if *debug {
		l, err = zap.NewDevelopment()
	} else {
		l, err = zap.NewProduction()
	}
	if err != nil {
		stdLog.Fatalf("can't initialize  zap logger: %v", err)
	}
	log := l.Sugar()

	log.Infow("hello world", "foo", "bar")
}
