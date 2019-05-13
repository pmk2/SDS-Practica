package main

import (
	"database/sql"
	"fmt"
	"runtime"
	"strings"

	_ "github.com/go-sql-driver/mysql" //Libreria para mysql
)

//Funcion para mostrar cuentas de un user
func mostrarCuentas(id string) []cuenta {
	var cuentas []cuenta
	var cuenta cuenta

	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/gestorpass")
	chk(err)
	defer db.Close() //Para que se cierre la bd al finalizar

	dato, err2 := db.Query("select * from cuentas where id_user = " + id)
	chk(err2)
	defer dato.Close() //Para que se cierre la query al finalizar

	for dato.Next() {
		var user, password, url, notes, credit string
		var idCuenta, idUser int
		err3 := dato.Scan(&idCuenta, &idUser, &user, &password, &url, &notes, &credit)
		chk(err3)

		cuenta.User = user
		cuenta.Pass = password
		cuenta.URL = url
		cuenta.Notes = notes
		cuenta.Credit = credit

		cuentas = append(cuentas, cuenta)

		//cuentas += fmt.Sprintln(". User: " + user + " || Pass: " + password + " || Sitio: " + url)
	}

	return cuentas
}

func devolverUser(user string) (bool, int, string, string) {
	existe := false
	var id_user int
	var usr, password, salt string

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

		for dato.Next() {
			err3 := dato.Scan(&id_user, &usr, &password, &salt)
			if err3 != nil {
				fmt.Print("Error escaneando fila: ")
				fmt.Println(err3)
			}
		}
		existe = true
	}

	return existe, id_user, password, salt
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
			var user, password, salt string
			for dato.Next() {
				err3 := dato.Scan(&id_user, &user, &password, &salt)
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
func insertUser(user string, pass string, salt string) bool {
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
			insert, err2 := db.Query("INSERT INTO usuario (user, password, salt) VALUES ('" + user + "','" + pass + "','" + salt + "')")
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
func insertCuenta(id string, user string, pass string, url string, notes string, credit string) bool {
	var ok bool
	ok = false

	//Comprobamos que user y pass no esten vacios ni null
	if user != "" && pass != "" && url != "" {
		//Abrimos la base de datos
		db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/gestorpass")
		chk(err)
		defer db.Close() //Para que se cierre la bd al finalizar

		//Query para insertar cuenta
		insert, err2 := db.Query("INSERT INTO cuentas (id_user, user, password, sitio_web, notas, tarjeta) VALUES (" + id + ",'" + user + "','" + pass + "','" + url + "','" + notes + "','" + credit + "')")
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
