package main

import (
	"database/sql"
	"fmt"
	"runtime"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql" //Libreria para mysql
)

//Funcion para mostrar cuentas de un user
func mostrarCuentas(user int) string {
	var cuentas string
	cuentas = ""

	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/gestorpass")
	chk(err)
	defer db.Close() //Para que se cierre la bd al finalizar

	idUser := strconv.Itoa(user) // Convertimos el int a string
	dato, err2 := db.Query("select * from cuentas where id_user = " + idUser)
	chk(err2)
	defer dato.Close() //Para que se cierre la query al finalizar

	contador := 1
	for dato.Next() {
		var user, password, sitio_web string
		var id_cuentas, id_user int
		err3 := dato.Scan(&id_cuentas, &id_user, &user, &password, &sitio_web)
		chk(err3)

		//1. Usuario: pmk2@alu.ua.es || Pass: pass || Sitio: ua.es
		//fmt.Print(id_cuentas)
		//fmt.Print(" ")
		//fmt.Print(id_user)
		cuentas += fmt.Sprint(contador)
		cuentas += fmt.Sprintln(". User: " + user + " || Pass: " + password + " || Sitio: " + sitio_web)
		contador++ // Incrementamos en uno el contador
	}

	return cuentas
}

// Funcion que comprueba user y pass, devuelve un int
// 0 si el usuario no existe
// -1 si la pass es incorrecta
// otro int con el id del user
func comprobarUser(user string, pass string) int {
	//var fout *os.File // fichero para almacenar pass
	//var err error
	var existe int
	existe = 0

	if user != "" && pass != "" {
		//Abrimos la base de datos
		db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/gestorpass")
		if err != nil {
			fmt.Print("Error abriendo BD: ")
			fmt.Println(err)
		}
		defer db.Close() //Para que se cierre la bd al finalizar

		//Query para contar users y comprobar que existe user
		num, _ := db.Query("select count(id_user) from usuario where user = '" + user + "'")
		defer num.Close() //Para que se cierre la query al finalizar

		var count int
		for num.Next() {
			err5 := num.Scan(&count)
			if err5 != nil {
				fmt.Println(err5)
			}
		}

		if count != 0 {
			//Query para obtener pass
			dato, err2 := db.Query("select * from usuario where user = '" + user + "'")
			if err2 != nil {
				fmt.Print("Error en la query: ")
				fmt.Println(err2)
			}
			defer dato.Close() //Para que se cierre la query al finalizar

			var id_user int
			var user, password string
			for dato.Next() {
				err3 := dato.Scan(&id_user, &user, &password)
				if err3 != nil {
					fmt.Print("Error escaneando fila: ")
					fmt.Println(err3)
				}
			}

			//Para igualar \n y \r\n de windows
			if runtime.GOOS == "windows" {
				password = strings.TrimRight(password, "\r\n")
				pass = strings.TrimRight(pass, "\r\n")
			} else {
				password = strings.TrimRight(password, "\n")
				pass = strings.TrimRight(pass, "\n")
			}

			fmt.Println("Esta es la pass BD: " + password)
			fmt.Println("Esta es la pass pa: " + pass)

			if strings.Compare(password, pass) == 0 {
				existe = id_user
			} else {
				existe = -1
			}
		}
	}

	return existe
}

//Funcion para hacer inserts en la base de datos
func insertUser(user string, pass string) bool {
	var ok bool
	ok = false

	//Comprobamos que user y pass no esten vacios ni null
	if user != "" && pass != "" {
		//Abrimos la base de datos
		db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/gestorpass")
		chk(err)
		defer db.Close() //Para que se cierre la bd al finalizar

		//Query para contar users y comprobar que existe user
		num, _ := db.Query("select count(id_user) from usuario where user = '" + user + "'")
		defer num.Close() //Para que se cierre la query al finalizar

		var count int
		for num.Next() {
			err5 := num.Scan(&count)
			chk(err5)
		}

		if count == 0 {
			//Query para obtener pass
			insert, err2 := db.Query("INSERT INTO usuario (user, password) VALUES ('" + user + "',SHA1('" + pass + "'))")
			if err2 != nil {
				fmt.Print("Error en la query: ")
				fmt.Println(err2)
			} else {
				ok = true
			}
			defer insert.Close() //Para que se cierre la query al finalizar
		}
	}

	return ok
}

//Funcion para hacer inserts en la base de datos
func insertCuenta(id int, user string, pass string, url string) bool {
	var ok bool
	ok = false

	//Comprobamos que user y pass no esten vacios ni null
	if user != "" && pass != "" && url != "" {
		//Abrimos la base de datos
		db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/gestorpass")
		chk(err)
		defer db.Close() //Para que se cierre la bd al finalizar

		//Query para obtener pass
		idUser := strconv.Itoa(id) // Convertimos el int a string
		insert, err2 := db.Query("INSERT INTO cuentas (id_user, user, password, sitio_web) VALUES (" + idUser + ",'" + user + "','" + pass + "','" + url + "')")
		if err2 != nil {
			fmt.Print("Error en la query: ")
			fmt.Println(err2)
		} else {
			ok = true
		}
		defer insert.Close() //Para que se cierre la query al finalizar
	}

	return ok
}
