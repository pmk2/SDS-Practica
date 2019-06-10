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

	"github.com/zserge/lorca"
)

func (c2 *usuario) setCuentaInsertar(login string, pass string, urlC string, notes string, credit string) {
	//c2.Lock()
	//defer c2.Unlock()
	var account cuenta
	account.User = login
	account.Pass = pass
	account.URL = urlC
	account.Notes = notes
	account.Credit = credit

	c2.cuentaInsertar = account
	c2.mensaje = "Login: " + login + " Pass: " + pass + " URL: " + urlC + " Notes: " + notes + " Credit Card: " + credit
}

func (c2 *usuario) obtenerCuentas() {
	var cuentas []cuenta
	cuentas = obtenerCuentasUser(c2)
	c2.cuentas = cuentas
	//fmt.Println(c2.cuentas)
}

func (c2 *usuario) getMensaje() string {
	//fmt.Println("Entra")
	return c2.mensaje
}

func (c2 *usuario) getCuentas() string {
	if c2.cuentas != nil {
		contador := 1
		cuentasUnidas := "*****Cuentas de " + c2.user + "*****\n"
		for i := 0; i < len(c2.cuentas); i++ {
			contStr := strconv.Itoa(contador)
			cuentasUnidas += contStr + ". Usuario: " + c2.cuentas[i].User + " || Password: " + c2.cuentas[i].Pass + " || URL: " + c2.cuentas[i].URL + " || Notes: " + c2.cuentas[i].Notes + " || Credit Card: " + c2.cuentas[i].Credit + "\n"
			contador++
		}
		c2.mensaje = cuentasUnidas
		fmt.Println(c2.mensaje)
	} else {
		c2.mensaje = "Token de sesión incorrecto o expirado. Cierre sesión y vuelva a conectarse."
	}

	return c2.mensaje
}

func (c2 *usuario) insertarCuenta() {
	resul := resp{}
	resul = insertCuenta(c2)
	c2.mensaje = resul.Msg
}

func (c2 *usuario) getRandomPassCuentas() string {
	random := randomPass()
	return random
}

func cuentas(user *usuario) {
	args := []string{}
	if runtime.GOOS == "linux" {
		args = append(args, "--class=Lorca")
	}
	ui, err := lorca.New("", "", 600, 830, args...)
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
	ui.Bind("getRandomPass", c2.getRandomPassCuentas)

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

/*
func transformarCuentas(cuentasString string) []cuenta {
	var cuentas []cuenta
	var cuenta cuenta
	var cuentasSplit []string

	cuentasSplit = strings.Split(cuentasString, "#")
	//fmt.Println(cuentasSplit[0])

	for i := 0; i < len(cuentasSplit)-1; i++ {
		userPassURL := strings.Split(cuentasSplit[i], "|")
		//fmt.Println(userPassURL)
		cuenta.User = userPassURL[0]
		cuenta.Pass = userPassURL[1]
		cuenta.URL = userPassURL[2]
		cuentas = append(cuentas, cuenta)
	}

	//fmt.Println(cuentas)

	return cuentas
}*/
