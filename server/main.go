// Main de server
package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func main() {
	//PARA DEPURACION COMENTO LAS LINEAS DE PASSWORD//
	/*var granted bool
	granted = false
	for !granted {
		fmt.Println("Enter password to access the server: ")
		password, _ := terminal.ReadPassword(int(syscall.Stdin))
		//fmt.Println()
		granted = comprobarClave(string(password))

		if !granted {
			fmt.Println("Clave incorrecta")
		}
	}
	//Si la clave es correcta lanzamos el server
	fmt.Println("Clave correcta")*/
	fmt.Println("Modo servidor: Esperando peticiones...")
	server()
}

func comprobarClave(clave string) bool {
	var iguales bool
	iguales = false

	b, err := ioutil.ReadFile("clave.txt") // just pass the file name
	if err != nil {
		fmt.Print(err)
	}

	str := string(b) // convert content to a 'string'
	if strings.Compare(str, clave) == 0 {
		iguales = true
	}

	return iguales
}
