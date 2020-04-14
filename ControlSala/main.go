package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)


var palabras[]string
var palabra int
var tiempo = time.Now().Unix()

type Jugador struct{
	ID string `json:"id,omitempty"`
	puntaje string `json:"puntaje,omitempty"`
	dibujado string `json:"dibujado,omitempty"`
}

var Jugadores[]Jugador

func leerPalabras( nombre string){
	file, err := os.Open(nombre)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		palabras =append(palabras, scanner.Text())
	}
}

func GetMostrarPalabra (w http.ResponseWriter, req *http.Request){
	palabra = rand.Intn(len(palabras))
	tiempo = time.Now().Unix()
	json.NewEncoder(w).Encode(palabras[palabra])
}

func GetPuntaje (w http.ResponseWriter, req * http.Request){
	Jugadores, err := ObtenerJugadores()
	if err != nil {
		fmt.Printf("Error obteniendo contactos: %v", err)
		return
	}
	for _, jugador := range Jugadores {
		json.NewEncoder(w).Encode(jugador.puntaje)
	}
}

func GetPuntajeJugador (w http.ResponseWriter, req * http.Request){
	params := mux.Vars(req)
	Jugadores, err := ObtenerJugadores()
	if err != nil {
		fmt.Printf("Error obteniendo jugadores: %v", err)
		return
	}
	for _, jugador := range Jugadores {
		if jugador.ID == params["id"]{
			json.NewEncoder(w).Encode(jugador.puntaje)
			return
		}
	}
	json.NewEncoder(w).Encode("No Encontrado")
}

func puntajeactual(id string) string{
	Jugadores, err := ObtenerJugadores()
	if err != nil {
		fmt.Printf("Error obteniendo jugadores: %v", err)
	}
	for _, jugador := range Jugadores {
		if jugador.ID == id{
			puntaje := jugador.puntaje
			return puntaje
		}
}
	return "0"
}

func palabraAcertada (w http.ResponseWriter, req *http.Request){
	params := mux.Vars(req)
	var jugador Jugador
	jugador.ID = params["id"]
	puntaje := puntajeactual(jugador.ID)
	i1, err := strconv.ParseInt(puntaje, 10,64)
	if err == nil {
		i1 = i1 + 50
	}
	puntaje = strconv.FormatInt(int64(i1),10)
	jugador.puntaje = puntaje

	errr := agregarPuntaje(jugador)
	if errr != nil {
		fmt.Printf("Error actualizando: %v", errr)
	} else {
		fmt.Println("Actualizando correctamente")
	}
}

func AgregarJugador (w http.ResponseWriter, req * http.Request){
	params := mux.Vars(req)
	var jugador Jugador
	_ = json.NewDecoder(req.Body).Decode(&jugador)
	jugador.ID = params["id"]
	fmt.Printf("%v\n", params)
	err:= insertarJugador(jugador)
	if err != nil{
		fmt.Printf("Error insertado: %V", err)
	}else{
		fmt.Printf("insertado correctamente")
	}
}

func BorrarJugador (w http.ResponseWriter, req * http.Request){
	params := mux.Vars(req)
	var jugador Jugador
	jugador.ID = params["id"]
	err := eliminarJugador(jugador)
	if err != nil {
		fmt.Printf("Error al eliminar: %v", err)
	} else {
		fmt.Println("Eliminado correctamente")
	}
}

func main() {
 	leerPalabras("Palabras")
 	router := mux.NewRouter()

	router.HandleFunc("/palabra", GetMostrarPalabra).Methods("GET")
	router.HandleFunc("/puntaje", GetPuntaje).Methods("GET")
	router.HandleFunc("/puntajeJugador/{id}", GetPuntajeJugador).Methods("GET")
	router.HandleFunc("/palabraAcertada/{id}", palabraAcertada).Methods("POST")
	router.HandleFunc("/agregarjugador/{id}", AgregarJugador).Methods("POST")
	router.HandleFunc("/quitarjugador/{id}", BorrarJugador).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":3000", router))
}






