package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

// respuesta del servidor de cuenta
type respCuenta struct {
	Ok      bool   // true -> correcto, false -> error
	Cuentas string // cuentas del user
}

type respPrueba struct {
	Ok      bool     `json:"Ok"`
	Cuentas []cuenta `json:"Cuentas"`
	//Token   string   `json:"Token"`
}

// respuesta del servidor
type resp struct {
	Ok    bool   // true -> correcto, false -> error
	Msg   string // mensaje adicional
	ID    int    //id del user
	Token string // token de sesion
}

// opcion 0 register, 1 login
func client(c *usuario, opc int) resp {
	nameUser := c.user
	passUser := c.pass

	/* creamos un cliente especial que no comprueba la validez de los certificados
	esto es necesario por que usamos certificados autofirmados (para pruebas) */
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	// hash con SHA512 de la contraseña
	keyClient := sha512.Sum512([]byte(passUser))
	keyLogin := keyClient[:32]  // una mitad para el login (256 bits)
	keyData := keyClient[32:64] // la otra para los datos (256 bits)

	//Pruebas para pass
	//fmt.Println(string(encode64(keyLogin)))

	// generamos un par de claves (privada, pública) para el servidor
	pkClient, err := rsa.GenerateKey(rand.Reader, 1024)
	chk(err)
	pkClient.Precompute() // aceleramos su uso con un precálculo

	pkJSON, err := json.Marshal(&pkClient) // codificamos con JSON
	chk(err)

	keyPub := pkClient.Public()           // extraemos la clave pública por separado
	pubJSON, err := json.Marshal(&keyPub) // y codificamos con JSON
	chk(err)

	var respuesta resp
	switch opc {
	case 0:
		// registro
		data := url.Values{}                 // estructura para contener los valores
		data.Set("cmd", "register")          // comando (string)
		data.Set("user", nameUser)           // usuario (string)
		data.Set("pass", encode64(keyLogin)) // "contraseña" a base64

		r, err := client.PostForm("https://localhost:10443", data) // enviamos por POST
		chk(err)

		responseData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}

		json.Unmarshal(responseData, &respuesta)

	case 1:
		// login
		data := url.Values{}
		data.Set("cmd", "login")             // comando (string)
		data.Set("user", nameUser)           // usuario (string)
		data.Set("pass", encode64(keyLogin)) // contraseña (a base64 porque es []byte)

		// comprimimos y codificamos la clave pública
		data.Set("pubkey", encode64(compress(pubJSON)))

		// comprimimos, ciframos y codificamos la clave privada
		data.Set("prikey", encode64(encrypt(compress(pkJSON), keyData)))

		r, err := client.PostForm("https://localhost:10443", data)
		chk(err)

		responseData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}

		json.Unmarshal(responseData, &respuesta)
	}

	return respuesta
}

// opcion 0 register, 1 login
func insertCuenta(c *usuario) resp {
	/* creamos un cliente especial que no comprueba la validez de los certificados
	esto es necesario por que usamos certificados autofirmados (para pruebas) */
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	// hash con SHA512 de la contraseña
	keyClient := sha512.Sum512([]byte(c.pass))
	//keyLogin := keyClient[:32]  // una mitad para el login (256 bits)
	keyData := keyClient[32:64] // la otra para los datos (256 bits)

	var respuesta resp
	// Anyadir cuenta al user
	userJSON, err := json.Marshal(c.cuentaInsertar.User) // Codificamos con JSON el user
	chk(err)
	passJSON, err := json.Marshal(c.cuentaInsertar.Pass) // Codificamos con JSON la pass
	chk(err)
	urlJSON, err := json.Marshal(c.cuentaInsertar.URL) // Codificamos con JSON la url
	chk(err)
	notesJSON, err := json.Marshal(c.cuentaInsertar.Notes) // Codificamos con JSON las notas
	chk(err)
	creditJSON, err := json.Marshal(c.cuentaInsertar.Credit) // Codificamos con JSON la tarjeta
	chk(err)

	data := url.Values{}
	data.Set("cmd", "addAccount")      // comando (string)
	data.Set("id", strconv.Itoa(c.id)) // id usuario (string)
	// comprimimos, ciframos y codificamos el user
	data.Set("user", encode64(encrypt(compress(userJSON), keyData)))
	// comprimimos, ciframos y codificamos la pass
	data.Set("pass", encode64(encrypt(compress(passJSON), keyData)))
	// comprimimos, ciframos y codificamos la url
	data.Set("url", encode64(encrypt(compress(urlJSON), keyData)))
	// comprimimos, ciframos y codificamos las notas
	data.Set("notes", encode64(encrypt(compress(notesJSON), keyData)))
	// comprimimos, ciframos y codificamos la tarjeta
	data.Set("credit", encode64(encrypt(compress(creditJSON), keyData)))
	// Pasamos, tambien, el token de sesion
	data.Set("token", c.token)

	r, err := client.PostForm("https://localhost:10443", data)
	chk(err)

	responseData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(responseData, &respuesta)

	return respuesta
}

//Función para obtener cuentas
func obtenerCuentasUser(c *usuario) []cuenta {
	var cuentas []cuenta

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	// hash con SHA512 de la contraseña
	keyClient := sha512.Sum512([]byte(c.pass))
	keyData := keyClient[32:64] // la otra para los datos (256 bits)

	var respuesta respPrueba
	data := url.Values{}
	data.Set("cmd", "getAccounts")     // comando (string)
	data.Set("id", strconv.Itoa(c.id)) // id usuario (string)
	data.Set("token", c.token)         // token de sesion

	r, err := client.PostForm("https://localhost:10443", data)
	chk(err)

	responseData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(responseData, &respuesta)
	// Comprobamos que el token a true or false
	if respuesta.Ok {
		//transformar String de cuentas a []cuenta
		var cuentasUser []cuenta
		cuentasUser = respuesta.Cuentas
		cuentas = decryptCuentas(cuentasUser, keyData)
	} else {
		cuentas = respuesta.Cuentas
	}
	return cuentas
}
