package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/signal"
	"runtime"
	"strconv"

	"github.com/sethvargo/go-password/password"
	"github.com/zserge/lorca"
)

type generatePass struct {
	characters     int
	numDigits      int
	numSymbols     int
	upperLower     bool
	repeatChar     bool
	randomPassword string
}

func (generate *generatePass) crearRandomPass() {
	pass, err := password.Generate(generate.characters, generate.numDigits, generate.numSymbols, generate.upperLower, generate.repeatChar)
	chk(err)

	generate.randomPassword = pass
}

func (generate *generatePass) setDatos(numChar int, numDig string, numSym string, altern string, repeat string) {
	generate.characters = numChar
	numDigInt, _ := strconv.Atoi(numDig)
	generate.numDigits = numDigInt
	numSymInt, _ := strconv.Atoi(numSym)
	generate.numSymbols = numSymInt

	if altern == "s" {
		generate.upperLower = false
	} else {
		generate.upperLower = true
	}
	if repeat == "s" {
		generate.repeatChar = true
	} else {
		generate.repeatChar = false
	}

	//fmt.Println(string(numChar) + " " + numDig + " " + numSym + " " + altern + " " + repeat)

}

func (generate *generatePass) getRandomPass() string {
	return generate.randomPassword
}

func randomPass() string {
	args := []string{}
	if runtime.GOOS == "linux" {
		args = append(args, "--class=Lorca")
	}
	ui, err := lorca.New("", "", 600, 755, args...)
	if err != nil {
		log.Fatal(err)
	}
	defer ui.Close()

	// A simple way to know when UI is ready (uses body.onload event in JS)
	ui.Bind("start", func() {
		log.Println("UI is ready")
	})

	// Create and bind Go object to the UI
	generate := &generatePass{}
	//c2 = user

	ui.Bind("crearRandomPass", generate.crearRandomPass)
	ui.Bind("getRandomPass", generate.getRandomPass)
	ui.Bind("setDatos", generate.setDatos)

	// Load HTML.
	b, err := ioutil.ReadFile("./www/indexRandom.html") // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	html := string(b) // convert content to a 'string'
	ui.Load("data:text/html," + url.PathEscape(html))

	// You may use console.log to debug your JS code, it will be printed via
	// log.Println(). Also exceptions are printed in a similar manner.
	ui.Eval(`
		console.log("Hello, world!");
		console.log('Multiple values:', [1, false, {"x":5}]);
	`)

	// Wait until the interrupt signal arrives or browser window is closed
	sigc := make(chan os.Signal)
	signal.Notify(sigc, os.Interrupt)
	select {
	case <-sigc:
	case <-ui.Done():
	}

	log.Println("exiting...")

	return generate.randomPassword
}
