package main

import (
	"bytes"
	"compress/zlib"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"runtime"
	"strings"
)

func decryptCuentas(cuentasEnc []cuenta, key []byte) []cuenta {
	var cuentasDec []cuenta
	var user, pass, url string
	var cuenta cuenta

	for i := 0; i < len(cuentasEnc); i++ {
		//		json.Unmarshal(responseData, &respuesta)
		json.Unmarshal(decompress(decrypt(decode64(cuentasEnc[i].User), key)), &user)
		json.Unmarshal(decompress(decrypt(decode64(cuentasEnc[i].Pass), key)), &pass)
		json.Unmarshal(decompress(decrypt(decode64(cuentasEnc[i].URL), key)), &url)
		//fmt.Println("Entra despues")
		cuenta.User = user
		cuenta.Pass = pass
		cuenta.URL = url

		cuentasDec = append(cuentasDec, cuenta)
	}

	return cuentasDec
}

//Encriptar pass con sha
func encryptPass(pass string) string {
	h := sha1.New()
	h.Write([]byte(pass))
	bs := h.Sum(nil)
	passEncrypt := fmt.Sprintf("%x\n", bs)
	//fmt.Println(passEncrypt)
	return passEncrypt
}

// función para cifrar (con AES en este caso), adjunta el IV al principio
func encrypt(data, key []byte) (out []byte) {
	out = make([]byte, len(data)+16)    // reservamos espacio para el IV al principio
	rand.Read(out[:16])                 // generamos el IV
	blk, err := aes.NewCipher(key)      // cifrador en bloque (AES), usa key
	chk(err)                            // comprobamos el error
	ctr := cipher.NewCTR(blk, out[:16]) // cifrador en flujo: modo CTR, usa IV
	ctr.XORKeyStream(out[16:], data)    // ciframos los datos
	return
}

// función para descifrar (con AES en este caso)
func decrypt(data, key []byte) (out []byte) {
	out = make([]byte, len(data)-16)     // la salida no va a tener el IV
	blk, err := aes.NewCipher(key)       // cifrador en bloque (AES), usa key
	chk(err)                             // comprobamos el error
	ctr := cipher.NewCTR(blk, data[:16]) // cifrador en flujo: modo CTR, usa IV
	ctr.XORKeyStream(out, data[16:])     // desciframos (doble cifrado) los datos
	return
}

// función para comprimir
func compress(data []byte) []byte {
	var b bytes.Buffer      // b contendrá los datos comprimidos (tamaño variable)
	w := zlib.NewWriter(&b) // escritor que comprime sobre b
	w.Write(data)           // escribimos los datos
	w.Close()               // cerramos el escritor (buffering)
	return b.Bytes()        // devolvemos los datos comprimidos
}

// función para descomprimir
func decompress(data []byte) []byte {
	var b bytes.Buffer // b contendrá los datos descomprimidos

	r, err := zlib.NewReader(bytes.NewReader(data)) // lector descomprime al leer

	chk(err)         // comprobamos el error
	io.Copy(&b, r)   // copiamos del descompresor (r) al buffer (b)
	r.Close()        // cerramos el lector (buffering)
	return b.Bytes() // devolvemos los datos descomprimidos
}

// función para codificar de []bytes a string (Base64)
func encode64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data) // sólo utiliza caracteres "imprimibles"
}

// función para decodificar de string a []bytes (Base64)
func decode64(s string) []byte {
	b, err := base64.StdEncoding.DecodeString(s) // recupera el formato original
	chk(err)                                     // comprobamos el error
	return b                                     // devolvemos los datos originales
}

// funcion para comprobar pass
func comprobarPass(pass1 string, pass2 string) bool {
	var iguales bool
	//Para igualar \n y \r\n de windows
	if runtime.GOOS == "windows" {
		pass1 = strings.TrimRight(pass1, "\r\n")
		pass2 = strings.TrimRight(pass2, "\r\n")
	} else {
		pass1 = strings.TrimRight(pass1, "\n")
		pass2 = strings.TrimRight(pass2, "\n")
	}

	//fmt.Println("Esta es la pass BD: " + passUser)
	//fmt.Println("Esta es la pass pa: " + stringHash)

	if strings.Compare(pass1, pass2) == 0 { // Comparamos las pass
		iguales = true
	} else {
		iguales = false
	}

	return iguales
}
