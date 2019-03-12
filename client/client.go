package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func client(c *usuario) {
	conn, err := net.Dial("tcp", "localhost:1337") // llamamos al servidor
	chk(err)
	defer conn.Close() // es importante cerrar la conexión al finalizar

	fmt.Println("conectado a ", conn.RemoteAddr())

	keyscan := bufio.NewScanner(os.Stdin) // scanner para la entrada estándar (teclado)
	netscan := bufio.NewScanner(conn)     // scanner para la conexión (datos desde el servidor)

	//netscan.Scan() //{ // escaneamos la conexión
	//fmt.Println("Entra")
	//fmt.Println(netscan.Text()) // mostramos mensaje desde el servidor
	//break
	//}

	//for keyscan.Scan() { // escaneamos el user
	fmt.Fprintln(conn, keyscan.Text()) // enviamos la entrada al servidor
	//netscan.Scan()                             // escaneamos la conexión
	//fmt.Println("servidor: " + netscan.Text()) // mostramos mensaje desde el servidor
	//break
	//}

	//netscan.Scan()              // escaneamos la conexión
	//fmt.Println(netscan.Text()) // mostramos mensaje desde el servidor

	//for keyscan.Scan() { // escaneamos el pass
	fmt.Fprintln(conn, keyscan.Text()) // enviamos la entrada al servidor
	//netscan.Scan()                             // escaneamos la conexión
	//fmt.Println("servidor: " + netscan.Text()) // mostramos mensaje desde el servidor
	//break
	//}

	//fmt.Println()

	for netscan.Scan() { // escaneamos la conexión
		fmt.Println(netscan.Text()) // mostramos las cuentas desde el servidor
	}
	//fmt.Println("Todos los datos solicitados")

	for {

	}
}
