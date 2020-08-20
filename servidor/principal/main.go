package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"

	"github.com/gorilla/mux"
)

//STRUCTS DEL WEB SERVICE
type structRam struct {
	Memoria_Total     int `json:"Memoria_Total,omitempty"`
	Memoria_en_uso    int `json:"Memoria_en_uso,omitempty"`
	Porcentaje_en_uso int `json:"Porcentaje_en_uso,omitempty"`
}

type StructListaRam struct {
	StructListaRam []structRam `json:"struct_lista_ram"`
}

type structProcesos struct {
	Pid           int     `json:"PID,omitempty"`
	Nombre        string  `json:"Nombre,omitempty"`
	Usuario       string  `json:"Usuario,omitempty"`
	Estado        string  `json:"Estado,omitempty"`
	PorcentajeRam float64 `json:"PorcentajeRam,omitempty"`
	Ppid          int     `json:"PPID,omitempty"`
}

type StructListaProcesos struct {
	StructListaProcesos []structProcesos `json:"struct_lista_procesos"`
}

type structKill struct {
	Pid string `json:"pid,omitempty"`
}

//VARIABLES
var (
	tamanio float64 = 0
)

func main() {
	//Obtenemos el tamanio
	leerInicio()
	//Inicio el codigo del servidor
	router := mux.NewRouter()

	router.HandleFunc("/", inicio)
	router.HandleFunc("/procesos", enviarProcesos).Methods("GET", "OPTIONS")
	router.HandleFunc("/ram", informacionRAM).Methods("GET", "OPTIONS")
	router.HandleFunc("/kill/{id}", matarProceso).Methods("POST", "OPTIONS")

	fmt.Println("El servidor se ha iniciado en el puerto 3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}

//FUNCIONES DEL SERVIDOR-------------------------------------------------------------------------------------------------------------------------------
//Esta funcion solo se agrego para que se vea bonito el servidor cuando inicia :)
func inicio(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	fmt.Fprintf(w, "Practica 1 Sistemas Operativos 2 \nFernando Vidal Ruiz Piox - 201503984 \nJeannira Del Rosario Sic Menéndez - 201602434")
}

//Esta funcion va a devolver el json con la info de los procesos
func enviarProcesos(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	file, _ := ioutil.ReadFile("/proc/cpu_grupo33")

	data := StructListaProcesos{}

	_ = json.Unmarshal([]byte(file), &data)

	//fmt.Println("Tamanio: ", len(data.StructListaProcesos))
	for i := 0; i < len(data.StructListaProcesos); i++ {
		var temporal float64 = data.StructListaProcesos[i].PorcentajeRam
		data.StructListaProcesos[i].PorcentajeRam = (temporal / tamanio) * 100
		//fmt.Println("Valor: ", data.StructListaProcesos[i])
	}

	json.NewEncoder(w).Encode(data.StructListaProcesos)
}

//Esta funcion va a devolver el json con la informacion de la pagina de la RAM
func informacionRAM(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	file, _ := ioutil.ReadFile("/proc/mem_grupo33")

	data := StructListaRam{}

	_ = json.Unmarshal([]byte(file), &data)

	/*fmt.Println("Tamanio: ", len(data.StructListaRam))
	for i := 0; i < len(data.StructListaRam); i++ {
		fmt.Println("Valor: ", data.StructListaRam[i])
	}*/

	json.NewEncoder(w).Encode(data.StructListaRam[0])
}

//Esta funcion va a matar el proceso especificado
func matarProceso(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	params := mux.Vars(req)
	var valor structKill
	valor.Pid = params["id"]
	
        cmd := exec.Command("sudo -9 kill "+valor.Pid)
	stdout, err := cmd.Output()
        //_, err := exec.Command("sh", "-c", "sudo -9 kill "+valor.Pid).Output()
	if err != nil {
		fmt.Printf("Error matando el proceso: %v", err)
	}
	json.NewEncoder(w).Encode(valor)
}

//FUNCIONES ADICIONALES-------------------------------------------------------------------------------------------------------------------------------
func leerInicio() {
	file, _ := ioutil.ReadFile("/proc/mem_grupo33")

	data := StructListaRam{}

	_ = json.Unmarshal([]byte(file), &data)

	tamanio = float64(data.StructListaRam[0].Memoria_Total)
	/*fmt.Println("Tamanio: ", tamanio)
	for i := 0; i < len(data.StructListaRam); i++ {
		fmt.Println("Valor: ", data.StructListaRam[i])
	}*/
}
