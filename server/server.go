package main

import (
	"crypto/rand"
	"encoding/json"
	"io"
	"net/http"

	"golang.org/x/crypto/scrypt"
)

// ejemplo de tipo para un usuario
type user struct {
	Name string            // nombre de usuario
	Hash []byte            // hash de la contraseña (en bd = string)
	Salt []byte            // sal para la contraseña
	Data map[string]string // datos adicionales del usuario
}

// mapa con todos los usuarios
// (se podría codificar en JSON y escribir/leer de disco para persistencia)
var gUsers map[string]user

// respuesta del servidor
type resp struct {
	Ok  bool   // true -> correcto, false -> error
	Msg string // mensaje adicional
}

// función para escribir una respuesta del servidor
func response(w io.Writer, ok bool, msg string) {
	r := resp{Ok: ok, Msg: msg}    // formateamos respuesta
	rJSON, err := json.Marshal(&r) // codificamos en JSON
	chk(err)                       // comprobamos error
	w.Write(rJSON)                 // escribimos el JSON resultante
}

// gestiona el modo servidor
func server() {

	gUsers = make(map[string]user) // inicializamos mapa de usuarios

	http.HandleFunc("/", handler) // asignamos un handler global

	// escuchamos el puerto 10443 con https y comprobamos el error
	// Para generar certificados autofirmados con openssl usar:
	//    openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes -subj "/C=ES/ST=Alicante/L=Alicante/O=UA/OU=Org/CN=www.ua.com"
	chk(http.ListenAndServeTLS(":10443", "cert.pem", "key.pem", nil))

	/*
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


					//Probando insert
					//--------------------------------------
					ok := insertUser(user, pass)
					if ok {
						fmt.Fprintln(conn, "Usuario registrado correctamente")
					} else {
						fmt.Fprintln(conn, "Usuario ya existente")
					}
					//--------------------------------------

				conn.Close() // cerramos al finalizar el cliente (EOF se envía con ctrl+d o ctrl+z según el sistema)
				fmt.Println("cierre[", port, "]")

			}()
		}*/
}

func handler(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()                              // es necesario parsear el formulario
	w.Header().Set("Content-Type", "text/plain") // cabecera estándar

	switch req.Form.Get("cmd") { // comprobamos comando desde el cliente
	case "register": // ** registro
		//var idUser int
		var existeUser bool
		//var passUser, saltUser string

		u := user{}
		u.Name = req.Form.Get("user")              // nombre
		u.Salt = make([]byte, 16)                  // sal (16 bytes == 128 bits)
		rand.Read(u.Salt)                          // la sal es aleatoria
		u.Data = make(map[string]string)           // reservamos mapa de datos de usuario
		u.Data["private"] = req.Form.Get("prikey") // clave privada
		u.Data["public"] = req.Form.Get("pubkey")  // clave pública
		password := decode64(req.Form.Get("pass")) // contraseña (keyLogin)

		// "hasheamos" la contraseña con scrypt
		u.Hash, _ = scrypt.Key(password, u.Salt, 16384, 8, 1, 32)
		//fmt.Println(string(encode64(u.Hash)))

		existeUser, _, _, _ = devolverUser(u.Name)

		//numUser 0 no registrado //-1 registrado pass incorrecta
		//Otro numero id del user
		if existeUser {
			msg := "Usuario ya registrado"
			response(w, false, msg)
		} else {
			insertado := insertUser(u.Name, string(encode64(u.Hash)), string(encode64(u.Salt)))
			if insertado {
				response(w, insertado, "Usuario registrado correctamente")
			} else {
				response(w, insertado, "Fallo al insertar user")
			}
		}

		/*
			_, ok := gUsers[u.Name] // ¿existe ya el usuario?
			if ok {
				response(w, false, "Usuario ya registrado")
			} else {
				gUsers[u.Name] = u
				response(w, true, "Usuario registrado")
			}*/

	case "login": // ** login
		//var idUser int
		var existeUser bool
		var passUser, saltUser string
		password := decode64(req.Form.Get("pass")) // obtenemos la contraseña

		existeUser, _, passUser, saltUser = devolverUser(req.Form.Get("user"))

		if existeUser {
			hash, _ := scrypt.Key(password, []byte(decode64(saltUser)), 16384, 8, 1, 32) // scrypt(contraseña)
			stringHash := string(encode64(hash))                                         // Pasamos el hash a string para compararlo con el de la bd

			if comprobarPass(passUser, stringHash) { // Comparamos las pass
				response(w, true, "Usuario válido")
			} else {
				response(w, false, "Contraseña inválida")
			}

		} else {
			response(w, false, "El usuario no existe")
		}

		/*
			u, ok := gUsers[req.Form.Get("user")] // ¿existe ya el usuario?
			if !ok {
				response(w, false, "Usuario inexistente")
				return
			}

			password := decode64(req.Form.Get("pass"))               // obtenemos la contraseña
			hash, _ := scrypt.Key(password, u.Salt, 16384, 8, 1, 32) // scrypt(contraseña)
			if bytes.Compare(u.Hash, hash) != 0 {                    // comparamos
				response(w, false, "Credenciales inválidas")
				return
			}
			response(w, true, "Credenciales válidas")
		*/

	default:
		response(w, false, "Comando inválido")
	}

}
