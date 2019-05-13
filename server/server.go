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

//Estructura de cuenta
type cuenta struct {
	User   string `json:"user"`
	Pass   string `json:"pass"`
	URL    string `json:"url"`
	Notes  string `json:"notes"`
	Credit string `json:"credit"`
}

// respuesta del servidor para validar user
type resp struct {
	Ok  bool   // true -> correcto, false -> error
	Msg string // mensaje adicional
	ID  int    //id del user
}

// respuesta del servidor de cuenta
type respCuenta struct {
	Ok      bool     `json:"Ok"`
	Cuentas []cuenta `json:"Cuentas"`
}

// función para escribir una respuesta del servidor
func response(w io.Writer, ok bool, msg string, id int) {
	r := resp{Ok: ok, Msg: msg, ID: id} // formateamos respuesta
	rJSON, err := json.Marshal(&r)      // codificamos en JSON
	chk(err)                            // comprobamos error
	w.Write(rJSON)                      // escribimos el JSON resultante
}

// función para escribir una respuesta del servidor
func responseCuentas(w io.Writer, ok bool, cuentas []cuenta) {
	r := respCuenta{Ok: ok, Cuentas: cuentas} // formateamos respuesta

	rJSON, err := json.Marshal(&r) // codificamos en JSON
	chk(err)                       // comprobamos error
	w.Write(rJSON)                 // escribimos el JSON resultante

	//------------Pruebas-----------------
	//rPrueba := respCuenta{Ok: ok, Cuentas: cuentas}
	//rPruebaJSON, _ := json.Marshal(rPrueba)
	//prueba := string(rPruebaJSON)
	//fmt.Println(prueba)
	//------------------------------------

	//var res respCuenta
	//accounts := make([]cuenta, len(cuentas))
	//json.Unmarshal(rPruebaJSON, &res)
	//fmt.Println(res)
}

// gestiona el modo servidor
func server() {

	http.HandleFunc("/", handler) // asignamos un handler global

	// escuchamos el puerto 10443 con https y comprobamos el error
	// Para generar certificados autofirmados con openssl usar:
	//    openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes -subj "/C=ES/ST=Alicante/L=Alicante/O=UA/OU=Org/CN=www.ua.com"
	chk(http.ListenAndServeTLS(":10443", "cert.pem", "key.pem", nil))
}

func handler(w http.ResponseWriter, req *http.Request) {
	//req.ParseForm()                              // es necesario parsear el formulario
	req.ParseMultipartForm(1024)
	w.Header().Set("Content-Type", "text/plain") // cabecera estándar

	//switch req.Form.Get("cmd") { // comprobamos comando desde el cliente
	switch req.PostFormValue("cmd") { // comprobamos comando desde el cliente
	case "register": // ** registro
		//var idUser int
		var existeUser bool
		//var passUser, saltUser string

		u := user{}
		u.Name = req.Form.Get("user") // nombre
		u.Salt = make([]byte, 16)     // sal (16 bytes == 128 bits)
		rand.Read(u.Salt)             // la sal es aleatoria
		//u.Data = make(map[string]string)           // reservamos mapa de datos de usuario
		//u.Data["private"] = req.Form.Get("prikey") // clave privada
		//u.Data["public"] = req.Form.Get("pubkey")  // clave pública
		password := decode64(req.Form.Get("pass")) // contraseña (keyLogin)

		// "hasheamos" la contraseña con scrypt
		u.Hash, _ = scrypt.Key(password, u.Salt, 16384, 8, 1, 32)
		//fmt.Println(string(encode64(u.Hash)))

		existeUser, _, _, _ = devolverUser(u.Name)

		//numUser 0 no registrado //-1 registrado pass incorrecta
		//Otro numero id del user
		if existeUser {
			msg := "Usuario ya registrado"
			response(w, false, msg, 0)
		} else {
			insertado := insertUser(u.Name, string(encode64(u.Hash)), string(encode64(u.Salt)))
			if insertado {
				response(w, insertado, "Usuario registrado correctamente", 0)
			} else {
				response(w, insertado, "Fallo al insertar user", 0)
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
		var idUser int
		password := decode64(req.Form.Get("pass")) // obtenemos la contraseña

		existeUser, idUser, passUser, saltUser = devolverUser(req.Form.Get("user"))

		if existeUser {
			hash, _ := scrypt.Key(password, []byte(decode64(saltUser)), 16384, 8, 1, 32) // scrypt(contraseña)
			stringHash := string(encode64(hash))                                         // Pasamos el hash a string para compararlo con el de la bd

			if comprobarPass(passUser, stringHash) { // Comparamos las pass
				response(w, true, "Usuario válido", idUser)
			} else {
				response(w, false, "Contraseña inválida", 0)
			}

		} else {
			response(w, false, "El usuario no existe", 0)
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

	case "addAccount": // ** login
		//var idUser int
		var idUser, userCuenta, passCuenta, urlCuenta, notasCuenta, tarjetaCuenta string
		idUser = req.Form.Get("id")
		userCuenta = req.Form.Get("user")
		passCuenta = req.Form.Get("pass")
		urlCuenta = req.Form.Get("url")
		notasCuenta = req.Form.Get("notes")
		tarjetaCuenta = req.Form.Get("credit")

		insertado := insertCuenta(idUser, userCuenta, passCuenta, urlCuenta, notasCuenta, tarjetaCuenta)

		if insertado {
			response(w, true, "Cuenta insertada", 0)

		} else {
			response(w, false, "Error al insertar", 0)
		}

	case "getAccounts": // ** Obtener cuentas de usuario
		//var idUser int
		var idUser string
		idUser = req.Form.Get("id")

		//insertado := insertCuenta(idUser, userCuenta, passCuenta, urlCuenta)

		cuentasUser := mostrarCuentas(idUser)

		//fmt.Println(cuentasUser)

		if len(cuentasUser) > 0 {
			responseCuentas(w, true, cuentasUser)

		} else {
			responseCuentas(w, false, cuentasUser)
		}

	default:
		response(w, false, "Comando inválido", 0)
	}

}
