package main

import (
	"fmt"
	"os"
	"time"
)

func restartLog() error {
	f, err := os.Create(logFile)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = fmt.Fprintf(f, "%s cspace log file\n", time.Now().Format("15:04:05.000000000"))
	return err
}

func glLog(s string) {
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = fmt.Fprintf(f, "%s %s\n", time.Now().Format("15:04:05.000000000"), s)
	if err != nil {
		panic(err)
	}

}

func glError(inError error) {
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, err = fmt.Fprintf(f, "%s %v\n", time.Now().Format("15:04:05.000000000"), inError)
	fmt.Fprintf(os.Stderr, "%s %v\n", time.Now().Format("15:04:05.000000000"), inError)
	if err != nil {
		panic(err)
	}
}
