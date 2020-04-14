package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)


func ObtenerBase () (db *sql.DB, e error){
	usuario := "root"
	pass := ""
	host := "tcp(127.0.0.1:3306)"
	BaseDeDatos := "sala"

	db, err := sql.Open("mysql",fmt.Sprintf("%s:%s@%s/%s", usuario, pass, host, BaseDeDatos))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return db, nil
}

func insertarJugador(c Jugador) (e error) {
	db, err := ObtenerBase()
	if err != nil {
		return err
	}
	defer db.Close()

	// Preparamos para prevenir inyecciones SQL
	sentenciaPreparada, err := db.Prepare("INSERT INTO Jugador (idJugador, Puntaje, Dibujado) VALUES(?, ?, ?)")
	if err != nil {
		return err
	}
	defer sentenciaPreparada.Close()
	// Ejecutar sentencia, un valor por cada '?'
	_, err = sentenciaPreparada.Exec(c.ID, c.puntaje, c.dibujado)
	if err != nil {
		return err
	}
	return nil
}

func insertarJugadores(){
	c := Jugador{
		ID:       5,
		puntaje:  15,
		dibujado: 2,
	}
	err:= insertarJugador(c)
	if err != nil{
		fmt.Printf("Error insertado: %V", err)
	}else{
		fmt.Printf("insertado correctamente")
	}
}

func ObtenerJugador() ([]Jugador, error){
	Jugadores := []Jugador{}
	db, err := ObtenerBase()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	filas, err := db.Query("SELECT idjugador, puntaje, dibujado FROM Jugador")

	if err != nil {
		return nil, err
	}
	// Si llegamos aquí, significa que no ocurrió ningún error
	defer filas.Close()

	// Aquí vamos a "mapear" lo que traiga la consulta en el while de más abajo
	var c Jugador

	// Recorrer todas las filas, en un "while"
	for filas.Next() {
		err = filas.Scan(&c.ID, &c.puntaje, &c.dibujado)
		// Al escanear puede haber un error
		if err != nil {
			return nil, err
		}
		// Y si no, entonces agregamos lo leído al arreglo
		Jugadores = append(Jugadores, c)
	}
	// Vacío o no, regresamos el arreglo de jugadores
	return Jugadores, nil
}

func eliminarJugador(c Jugador)error{
	db, err := ObtenerBase()
	if err != nil {
		return err
	}
	defer db.Close()

	sentenciaPreparada, err := db.Prepare("DELETE FROM jugador WHERE idJugador = ?")
	if err != nil {
		return err
	}
	defer sentenciaPreparada.Close()
	_, err = sentenciaPreparada.Exec(c.ID)
	if err != nil {
		return err
	}
	return nil
}

func agregarPuntaje(c Jugador) error {
	db, err := ObtenerBase()
	if err != nil {
		return err
	}
	defer db.Close()

	sentenciaPreparada, err := db.Prepare("UPDATE jugador SET Puntaje = ? WHERE idJugador = ?")
	if err != nil {
		return err
	}
	defer sentenciaPreparada.Close()
	// Pasar argumentos en el mismo orden que la consulta
	_, err = sentenciaPreparada.Exec(c.puntaje, c.ID)
	return err // Ya sea nil o sea un error, lo manejaremos desde donde hacemos la llamada
}

