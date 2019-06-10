//go:generate go run -tags generate gen.go
//Es necesario hacer un go get github.com/zserge/lorca para cargarlo

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
	"sync"

	"github.com/zserge/lorca"
)

// Go types that are bound to the UI must be thread-safe, because each binding
// is executed in its own goroutine. In this simple case we may use atomic
// operations, but for more complex cases one should use proper synchronization.
type usuario struct {
	sync.Mutex
	id             int
	user           string
	pass           string
	validado       bool
	mensaje        string
	cuentas        []cuenta
	cuentaInsertar cuenta
}

//Estructura de cuenta
type cuenta struct {
	User   string `json:"user"`
	Pass   string `json:"pass"`
	URL    string `json:"url"`
	Notes  string `json:"notes"`
	Credit string `json:"credit"`
}

func (c *usuario) setDatosUser(us string, pas string) {
	c.Lock()
	defer c.Unlock()
	c.user = us
	c.pass = pas
}

func (c *usuario) getPass() string {
	c.Lock()
	defer c.Unlock()
	return c.pass
}

func (c *usuario) getUser() string {
	c.Lock()
	defer c.Unlock()
	return c.user
}

func (c *usuario) getValidado() string {
	c.Lock()
	defer c.Unlock()
	s := strconv.FormatBool(c.validado)
	//fmt.Println("entra")
	return s
}

/*
func (c *usuario) getCuentas() []cuenta {
	c.Lock()
	defer c.Unlock()
	//fmt.Println(c.cuentas)
	return c.cuentas
}*/

func (c *usuario) getMSG() string {
	c.Lock()
	defer c.Unlock()
	//fmt.Println(c.cuentas)
	return c.mensaje
}

func (c *usuario) validarUser() {
	c.Lock()
	defer c.Unlock()
	resul := resp{}
	resul = client(c, 1)
	c.validado = resul.Ok
	c.mensaje = resul.Msg
	c.id = resul.ID
}

func (c *usuario) registerUser() {
	c.Lock()
	defer c.Unlock()
	resul := resp{}
	resul = client(c, 0)
	c.validado = resul.Ok
	c.mensaje = resul.Msg
	c.id = resul.ID
	//c.cuentas = client(c, 0)
}

func (c *usuario) cambiarPantalla() {
	c.Lock()
	defer c.Unlock()
	cuentas(c)
}

func (c *usuario) getRandomPass() string {
	c.Lock()
	defer c.Unlock()
	random := randomPass()
	return random
}

func main() {
	args := []string{}
	if runtime.GOOS == "linux" {
		args = append(args, "--class=Lorca")
	}
	ui, err := lorca.New("", "", 600, 490, args...)
	if err != nil {
		log.Fatal(err)
	}
	defer ui.Close()

	// A simple way to know when UI is ready (uses body.onload event in JS)
	ui.Bind("start", func() {
		log.Println("UI is ready")
	})

	// Create and bind Go object to the UI
	c := &usuario{}

	ui.Bind("setDatosUser", c.setDatosUser)
	ui.Bind("getUser", c.getUser)
	ui.Bind("getPass", c.getPass)
	ui.Bind("getValidado", c.getValidado)
	ui.Bind("validarUser", c.validarUser)
	ui.Bind("registerUser", c.registerUser)
	ui.Bind("getRandomPass", c.getRandomPass)
	//ui.Bind("getCuentas", c.getCuentas)
	ui.Bind("getMSG", c.getMSG)
	ui.Bind("cambiarPantalla", c.cambiarPantalla)

	// Load HTML.
	b, err := ioutil.ReadFile("./www/index.html") // just pass the file name
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
