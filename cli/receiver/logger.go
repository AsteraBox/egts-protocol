package main

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

type PlainFormatter struct {
	LevelDesc []string
}

func (f *PlainFormatter) Format(entry *log.Entry) ([]byte, error) {
	timestamp := fmt.Sprint(entry.Time.Format(time.RFC3339))
	return []byte(fmt.Sprintf("%s %s %s\n", f.LevelDesc[entry.Level], timestamp, entry.Message)), nil
}

func ConfigureLogger(level log.Level) {
	log.SetLevel(level)
	customFormatter := new(PlainFormatter)
	customFormatter.LevelDesc = []string{"PANC", "FATL", "ERRO", "WARN", "INFO", "DEBG"}
	log.SetFormatter(customFormatter)
}
