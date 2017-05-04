package main

import (
	"log"
	"os"
	"runtime"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()

	// Set the working directory to the root of Go package, so that its assets can be accessed.
	dir, err := importPathToDir("github.com/stojg/cspace")
	if err != nil {
		log.Fatalln("Unable to find Go package in your GOPATH, it's needed to load assets:", err)
	}
	err = os.Chdir(dir)
	if err != nil {
		log.Panicln("os.Chdir:", err)
	}
}
