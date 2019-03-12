package main

import (
	"bufio"
	"fmt"
	"net"
)

// gestiona el modo servidor
func server() {

	var user string
	var pass string

	ln, err := net.Listen("tcp", "localhost:1337") // escucha en espera de conexión
	chk(err)
	defer ln.Close() // nos aseguramos que cerramos las conexiones aunque el programa falle

	for { // búcle infinito, se sale con ctrl+c
		conn, err := ln.Accept() // para cada nueva petición de conexión
		chk(err)
		go func() { // lanzamos un cierre (lambda, función anónima) en concurrencia

			_, port, err := net.SplitHostPort(conn.RemoteAddr().String()) // obtenemos el puerto remoto para identificar al cliente (decorativo)
			chk(err)

			fmt.Println("conexión: ", conn.LocalAddr(), " <--> ", conn.RemoteAddr())

			scanner := bufio.NewScanner(conn) // el scanner nos permite trabajar con la entrada línea a línea (por defecto)

			for {
				//fmt.Fprintln(conn, "Escriba user: ") //Solicitamos al cliente su user
				for scanner.Scan() { // escaneamos la conexión para obtener user
					user = scanner.Text()
					//fmt.Println("cliente[", port, "]: ", scanner.Text()) // mostramos el mensaje del cliente
					//fmt.Fprintln(conn, "ack: ", scanner.Text())          // enviamos ack al cliente
					break
				}
				//fmt.Fprintln(conn, "Escriba pass: ")
				for scanner.Scan() { // escaneamos la conexión para obtener pass
					pass = scanner.Text()
					//fmt.Println("cliente[", port, "]: ", scanner.Text()) // mostramos el mensaje del cliente
					//fmt.Fprintln(conn, "ack: ", scanner.Text())          // enviamos ack al cliente
					break
				}
				break
			}

			passEncrypt := encryptPass(pass)
			var numUser int
			numUser = comprobarUser(user, passEncrypt)

			var cuentas string
			if numUser != 0 && numUser != -1 {
				cuentas = "*****Cuentas pertenecientes a " + user + "*****\n"
				cuentas += mostrarCuentas(numUser)
			} else if numUser == 0 {
				cuentas = "El usuario no existe"
			} else {
				cuentas = "La contraseña es incorrecta"
			}

			//Enviamos las cuentas al cliente
			fmt.Fprintln(conn, cuentas)

			/*
				//Probando insert
				//--------------------------------------
				ok := insertUser(user, pass)
				if ok {
					fmt.Fprintln(conn, "Usuario registrado correctamente")
				} else {
					fmt.Fprintln(conn, "Usuario ya existente")
				}
				//--------------------------------------
			*/
			conn.Close() // cerramos al finalizar el cliente (EOF se envía con ctrl+d o ctrl+z según el sistema)
			fmt.Println("cierre[", port, "]")

		}()
	}
}
