package main

import "fmt"

// función para comprobar errores (ahorra escritura)
func chk(e error) {
	if e != nil {
		//panic(e)
		fmt.Println("PRobando")
		fmt.Println(e)
	}
}
