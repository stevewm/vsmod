package cmd

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

var debug bool

type PlainFormatter struct {
}

func (f *PlainFormatter) Format(entry *log.Entry) ([]byte, error) {
	return fmt.Appendf(nil, "%s", entry.Message), nil
}

func toggleDebug() {
	if debug {
		log.Info("Debug logs enabled")
		log.SetLevel(log.DebugLevel)
		log.SetFormatter(&log.TextFormatter{})
	} else {
		plainFormatter := new(PlainFormatter)
		log.SetFormatter(plainFormatter)
	}
}
