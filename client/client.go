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
)

// respuesta del servidor
type resp struct {
	Ok  bool   // true -> correcto, false -> error
	Msg string // mensaje adicional
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
		// ** ejemplo de registro
		data := url.Values{}                 // estructura para contener los valores
		data.Set("cmd", "register")          // comando (string)
		data.Set("user", nameUser)           // usuario (string)
		data.Set("pass", encode64(keyLogin)) // "contraseña" a base64

		// comprimimos y codificamos la clave pública
		data.Set("pubkey", encode64(compress(pubJSON)))

		// comprimimos, ciframos y codificamos la clave privada
		data.Set("prikey", encode64(encrypt(compress(pkJSON), keyData)))

		r, err := client.PostForm("https://localhost:10443", data) // enviamos por POST
		chk(err)

		//io.Copy(os.Stdout, r.Body) // mostramos el cuerpo de la respuesta (es un reader)
		//fmt.Println()

		responseData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}

		json.Unmarshal(responseData, &respuesta)

	case 1:
		// ** ejemplo de login
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
		//io.Copy(os.Stdout, r.Body) // mostramos el cuerpo de la respuesta (es un reader)
		//fmt.Println()

		responseData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}

		json.Unmarshal(responseData, &respuesta)

		//responseString := string(responseData)
		//fmt.Println(responseString)
	}
	return respuesta
}
