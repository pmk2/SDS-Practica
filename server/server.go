package main

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/scrypt"
)

// Create the JWT key used to create the signature
var jwtKey = []byte("my_secret_key")

type tokenSession struct {
	Token   string
	Expires time.Time
}

var usersToken map[string]tokenSession

// Claims Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

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
	Ok    bool   // true -> correcto, false -> error
	Msg   string // mensaje adicional
	ID    int    // id del user
	Token string // token de sesion
}

// respuesta del servidor de cuenta
type respCuenta struct {
	Ok      bool     `json:"Ok"`
	Cuentas []cuenta `json:"Cuentas"`
}

// función para escribir una respuesta del servidor
func response(w io.Writer, ok bool, msg string, id int, token string) {
	r := resp{Ok: ok, Msg: msg, ID: id, Token: token} // formateamos respuesta
	rJSON, err := json.Marshal(&r)                    // codificamos en JSON
	chk(err)                                          // comprobamos error
	w.Write(rJSON)                                    // escribimos el JSON resultante
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
	usersToken = make(map[string]tokenSession) // inicializamos mapa de usuarios

	http.HandleFunc("/", handler) // asignamos un handler global

	// escuchamos el puerto 10443 con https y comprobamos el error
	// Para generar certificados autofirmados con openssl usar:
	//    openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes -subj "/C=ES/ST=Alicante/L=Alicante/O=UA/OU=Org/CN=www.ua.com"
	chk(http.ListenAndServeTLS(":10443", "cert.pem", "key.pem", nil))
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "*")
}

func handler(w http.ResponseWriter, req *http.Request) {

	//req.ParseForm()                              // es necesario parsear el formulario
	req.ParseMultipartForm(1024)
	w.Header().Set("Content-Type", "text/plain") // cabecera estándar

	//switch req.Form.Get("cmd") { // comprobamos comando desde el cliente
	switch req.PostFormValue("cmd") { // comprobamos comando desde el cliente
	case "register": // ** registro
		var idUser int
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
			response(w, false, msg, 0, "")
		} else {
			insertado := insertUser(u.Name, string(encode64(u.Hash)), string(encode64(u.Salt)))
			if insertado {
				_, idUser, _, _ = devolverUser(u.Name)
				token := crearTokenSesion(strconv.Itoa(idUser)) // Creamos el token de sesion del nuevo usuario
				fmt.Println(token)
				response(w, insertado, "Usuario registrado correctamente", idUser, token)
			} else {
				response(w, insertado, "Fallo al insertar user", 0, "")
			}
		}

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
				token := crearTokenSesion(strconv.Itoa(idUser)) // Creamos token de sesión
				//fmt.Println(usersToken)

				response(w, true, "Usuario válido", idUser, token)

			} else {
				response(w, false, "Contraseña inválida", 0, "")
			}

		} else {
			if idUser == 0 {
				response(w, false, "La base de datos no está disponible, inténtelo de nuevo más tarde", 0, "")
			} else if idUser == -1 {
				response(w, false, "El usuario no existe", 0, "")
			}
		}

	case "addAccount": // ** Anyadir cuenta
		//var idUser int
		var idUser, userCuenta, passCuenta, urlCuenta, notasCuenta, tarjetaCuenta, tokenUser string
		idUser = req.Form.Get("id")
		userCuenta = req.Form.Get("user")
		passCuenta = req.Form.Get("pass")
		urlCuenta = req.Form.Get("url")
		notasCuenta = req.Form.Get("notes")
		tarjetaCuenta = req.Form.Get("credit")
		tokenUser = req.Form.Get("token")

		if comprobarToken(tokenUser, idUser) {
			insertado := insertCuenta(idUser, userCuenta, passCuenta, urlCuenta, notasCuenta, tarjetaCuenta)

			if insertado {
				response(w, true, "Cuenta insertada", 0, "")

			} else {
				response(w, false, "Error al insertar", 0, "")
			}
		} else {
			response(w, false, "Token de sesión incorrecto o expirado. Cierre sesión y vuelva a conectarse.", 0, "")
		}

	case "getAccounts": // ** Obtener cuentas de usuario
		var idUser string
		idUser = req.Form.Get("id")
		tokenUser := req.Form.Get("token")

		//fmt.Println(usersToken)

		if comprobarToken(tokenUser, idUser) { //Si el token es correcto y no ha expirado
			cuentasUser := mostrarCuentas(idUser)
			if len(cuentasUser) > 0 {
				responseCuentas(w, true, cuentasUser)
			} else {
				responseCuentas(w, true, cuentasUser)
			}
		} else {
			responseCuentas(w, false, nil)
		}

	case "getAccountss":

		var idUser string
		idUser = req.Form.Get("user")
		if idUser == "" {
			response(w, false, "No hay usuario activo", 0, "")
		} else {
			cuentasUser := mostrarCuentas(idUser)
			var cuentasDec []cuenta
			keyClient := sha512.Sum512([]byte(req.Form.Get("pass")))
			key := keyClient[32:64]
			cuentasDec = decryptCuentas(cuentasUser, key)
			cuentas := getCuentas(cuentasDec)

			if len(cuentasDec) > 0 {
				response(w, true, cuentas, 0, "")

			} else {
				response(w, false, "No dispone de cuentas", 0, "")
			}
		}

	case "loginext":

		keyClient := sha512.Sum512([]byte(req.Form.Get("pass")))
		keyLogin := keyClient[:32]
		var existeUser bool
		var passUser, saltUser string
		var idUser int
		password := keyLogin
		existeUser, idUser, passUser, saltUser = devolverUser(req.Form.Get("user"))

		if existeUser {
			hash, _ := scrypt.Key(password, []byte(decode64(saltUser)), 16384, 8, 1, 32)
			stringHash := string(encode64(hash))

			if comprobarPass(passUser, stringHash) { // Comparamos las pass
				token := crearTokenSesion(strconv.Itoa(idUser)) // Creamos token de sesión
				//fmt.Println(usersToken)

				response(w, true, "Usuario válido", idUser, token)

			} else {
				response(w, false, "Contraseña inválida", 0, "")
			}

		} else {
			if idUser == 0 {
				response(w, false, "La base de datos no está disponible, inténtelo de nuevo más tarde", 0, "")
			} else if idUser == -1 {
				response(w, false, "El usuario no existe", 0, "")
			}
		}

	default:
		response(w, false, "Comando inválido", 0, "")
	}

}

func getCuentas(cuentas []cuenta) string {
	//FALTA DEVOLVER STRING CON TODAS LAS CUENTAS
	contador := 1
	cuentasUnidas := ""
	for i := 0; i < len(cuentas); i++ {
		contStr := strconv.Itoa(contador)
		cuentasUnidas += contStr + ". Usuario: " + cuentas[i].User + "|| Password: " + cuentas[i].Pass + "|| URL: " + cuentas[i].URL + "\n"
		contador++
	}
	return cuentasUnidas
}

// Función para crear el token de sesión
func crearTokenSesion(user string) string {
	var tokenSes tokenSession
	// Declare the expiration time of the token
	// here, we have kept it as 1 minutes
	expirationTime := time.Now().Add(10 * time.Second)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: user,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, _ := token.SignedString(jwtKey)
	//fmt.Println(tokenString)
	tokenSes.Token = tokenString
	tokenSes.Expires = expirationTime
	usersToken[user] = tokenSes
	//fmt.Println(usersToken)

	return tokenString
}

// Función para comprobar el token de sesion recibido y el guardado
func comprobarToken(tokenRecibido string, id string) bool {
	//tokenUser, _ := usersToken[id]
	//fmt.Println(usersToken)
	var ok bool
	ok = true

	// Initialize a new instance of `Claims`
	claims := &Claims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(tokenRecibido, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if !tkn.Valid {
		ok = false
	}
	if err != nil {
		ok = false
	}

	/*if tokenRecibido == tokenUser.Token {
		ok = true
	}*/

	return ok
}
