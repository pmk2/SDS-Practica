//go:generate go run -tags generate gen.go
//Es necesario hacer un go get github.com/zserge/lorca para cargarlo

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/url"
	"os"
	"os/signal"
	"runtime"
	"sync"

	"github.com/zserge/lorca"
)

// Go types that are bound to the UI must be thread-safe, because each binding
// is executed in its own goroutine. In this simple case we may use atomic
// operations, but for more complex cases one should use proper synchronization.
type usuario struct {
	sync.Mutex
	user     string
	pass     string
	register bool
	cuentas  string
}

func (c *usuario) setUser(us string) {
	c.Lock()
	defer c.Unlock()
	c.user = us
}

func (c *usuario) setPass(us string) {
	c.Lock()
	defer c.Unlock()
	c.pass = us
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

func (c *usuario) getRegister() bool {
	c.Lock()
	defer c.Unlock()
	return c.register
}

func (c *usuario) getCuentas() string {
	c.Lock()
	defer c.Unlock()
	//fmt.Println(c.cuentas)
	return c.cuentas
}

func (c *usuario) validarUser() {
	c.Lock()
	defer c.Unlock()

	//Conectamos con el server
	conn, err := net.Dial("tcp", "localhost:1337") // llamamos al servidor
	chk(err)
	defer conn.Close() // es importante cerrar la conexión al finalizar

	fmt.Println("conectado a ", conn.RemoteAddr())
	netscan := bufio.NewScanner(conn) // scanner para la conexión (datos desde el servidor)

	var cuentas string
	fmt.Fprintln(conn, c.user)
	fmt.Fprintln(conn, c.pass)
	for netscan.Scan() { // escaneamos la conexión
		cuentas += netscan.Text() // guardamos las cuentas desde el servidor
		cuentas += "\r\n"
		//fmt.Println(cuentas)
		//fmt.Println(netscan.Text())
		//break
	}
	c.cuentas = cuentas
	fmt.Println(cuentas)
	/*
		if strings.Compare(c.cuentas, "El usuario no existe") == 0 {
			c.register = false
		} else {
			c.register = true
		}*/

}

func main() {
	args := []string{}
	if runtime.GOOS == "linux" {
		args = append(args, "--class=Lorca")
	}
	ui, err := lorca.New("", "", 600, 500, args...)
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

	ui.Bind("setUser", c.setUser)
	ui.Bind("setPass", c.setPass)
	ui.Bind("getUser", c.getUser)
	ui.Bind("getPass", c.getPass)
	ui.Bind("getValidado", c.getRegister)
	ui.Bind("validarUser", c.validarUser)
	ui.Bind("getCuentas", c.getCuentas)

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
