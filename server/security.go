package main

import (
	"crypto/sha1"
	"fmt"
)

//Encriptar pass con sha
func encryptPass(pass string) string {
	h := sha1.New()
	h.Write([]byte(pass))
	bs := h.Sum(nil)
	passEncrypt := fmt.Sprintf("%x\n", bs)
	//fmt.Println(passEncrypt)
	return passEncrypt
}
