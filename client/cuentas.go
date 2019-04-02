package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/signal"
	"runtime"

	"github.com/zserge/lorca"
)

func (c2 *usuario) setDatosCuenta(login string, pass string, url string) {
	//c2.Lock()
	//defer c2.Unlock()
	c2.cuentas = "Login: " + login + " Pass: " + pass + " URL: " + url
	//fmt.Println(c2.user)
	//fmt.Println(c2.pass)
	//fmt.Println(c2.cuentas)
}

func (c2 *usuario) obtenerCuentas() string {
	//c2.Lock()
	//defer c2.Unlock()
	return c2.cuentas
}

func cuentas(user *usuario) {
	args := []string{}
	if runtime.GOOS == "linux" {
		args = append(args, "--class=Lorca")
	}
	ui, err := lorca.New("", "", 600, 685, args...)
	if err != nil {
		log.Fatal(err)
	}
	defer ui.Close()

	// A simple way to know when UI is ready (uses body.onload event in JS)
	ui.Bind("start", func() {
		log.Println("UI is ready")
	})

	// Create and bind Go object to the UI
	c2 := &usuario{}
	c2 = user

	ui.Bind("setDatosCuenta", c2.setDatosCuenta)
	ui.Bind("obtenerCuentas", c2.obtenerCuentas)

	// Load HTML.
	b, err := ioutil.ReadFile("./www/indexCuentas.html") // just pass the file name
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
}
