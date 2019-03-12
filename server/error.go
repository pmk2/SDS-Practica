package main

import (
	_ "github.com/go-sql-driver/mysql"
)

// función para comprobar errores (ahorra escritura)
func chk(e error) {
	if e != nil {
		panic(e)
	}
}
