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
	"strings"

	"github.com/zserge/lorca"
)

func (c2 *usuario) setCuentaInsertar(login string, pass string, urlC string) {
	//c2.Lock()
	//defer c2.Unlock()
	var account cuenta
	account.user = login
	account.pass = pass
	account.url = urlC

	c2.cuentaInsertar = account
	c2.mensaje = "Login: " + login + " Pass: " + pass + " URL: " + urlC
}

func (c2 *usuario) obtenerCuentas() {
	var cuentas []cuenta
	cuentas = obtenerCuentasUser(c2)
	c2.cuentas = cuentas
	fmt.Println(c2.cuentas)
}

func (c2 *usuario) getMensaje() string {
	//fmt.Println("Entra")
	return c2.mensaje
}

func (c2 *usuario) getCuentas() string {
	//FALTA DEVOLVER STRING CON TODAS LAS CUENTAS
	contador := 1
	cuentasUnidas := "*****Cuentas de " + c2.user + "*****\n"
	for i := 0; i < len(c2.cuentas); i++ {
		contStr := strconv.Itoa(contador)
		cuentasUnidas += contStr + ". Usuario: " + c2.cuentas[i].user + "|| Password: " + c2.cuentas[i].pass + "|| URL: " + c2.cuentas[i].url + "\n"
		contador++
	}
	c2.mensaje = cuentasUnidas
	fmt.Println(c2.mensaje)
	return c2.mensaje
}

func (c2 *usuario) insertarCuenta() {
	resul := resp{}
	resul = insertCuenta(c2)
	c2.mensaje = resul.Msg
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

	ui.Bind("setCuentaInsertar", c2.setCuentaInsertar)
	ui.Bind("obtenerCuentas", c2.obtenerCuentas)
	ui.Bind("insertarCuenta", c2.insertarCuenta)
	ui.Bind("getMensaje", c2.getMensaje)
	ui.Bind("getCuentas", c2.getCuentas)

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

func transformarCuentas(cuentasString string) []cuenta {
	var cuentas []cuenta
	var cuenta cuenta
	var cuentasSplit []string

	cuentasSplit = strings.Split(cuentasString, "#")
	//fmt.Println(cuentasSplit[0])

	for i := 0; i < len(cuentasSplit)-1; i++ {
		userPassURL := strings.Split(cuentasSplit[i], "|")
		//fmt.Println(userPassURL)
		cuenta.user = userPassURL[0]
		cuenta.pass = userPassURL[1]
		cuenta.url = userPassURL[2]
		cuentas = append(cuentas, cuenta)
	}

	//fmt.Println(cuentas)

	return cuentas
}
