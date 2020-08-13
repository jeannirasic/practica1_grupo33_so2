package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

//STRUCTS DE LA PARTE LOGICA
type DatosRam struct {
	total      float64
	consumida  float64
	porcentaje float64
}

type Procesos struct {
	pid           int
	nombre        string
	usuario       string
	estado        string
	porcentajeRam float64
	ppid          int
}

type Usuario struct {
	pid    int
	nombre string
}

//STRUCTS DEL WEB SERVICE
type structRam struct {
	Total      string `json:"total,omitempty"`
	Consumida  string `json:"consumida,omitempty"`
	Porcentaje string `json:"porcentaje,omitempty"`
}

type structProcesos struct {
	Pid           string `json:"pid,omitempty"`
	Nombre        string `json:"nombre,omitempty"`
	Usuario       string `json:"usuario,omitempty"`
	Estado        string `json:"estado,omitempty"`
	PorcentajeRam string `json:"porcentajeram,omitempty"`
	Ppid          string `json:"ppid,omitempty"`
}

type structKill struct {
	Pid string `json:"pid,omitempty"`
}

type structCpu struct {
	Porcentaje string `json:"porcentaje,omitempty"`
}

//VARIABLES
var (
	arregloUsuarios []Usuario
	monitorRAM      = DatosRam{0, 0, 0}
	listaProcesos   []Procesos
	monitorCPU      float64
)

func main() {
	//Leo los metodos para llenar las variables
	obtenerUsuarios()
	leerMeminfo()
	leerProcesos()
	leerStat()
	//Inicio el codigo del servidor
	router := mux.NewRouter()

	router.HandleFunc("/", inicio)
	router.HandleFunc("/principal", informacionPrincipal).Methods("GET", "OPTIONS")
	router.HandleFunc("/ram", informacionRAM).Methods("GET", "OPTIONS")
	router.HandleFunc("/cpu", informacionCPU).Methods("GET", "OPTIONS")
	router.HandleFunc("/kill/{id}", matarProceso).Methods("POST", "OPTIONS")

	fmt.Println("El servidor se ha iniciado en el puerto 3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}

//FUNCIONES DEL SERVIDOR-------------------------------------------------------------------------------------------------------------------------------
//Esta funcion solo se agrego para que se vea bonito el servidor cuando inicia :)
func inicio(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	fmt.Fprintf(w, "Proyecto 1 Sistemas Operativos 1 \nFernando Vidal Ruiz Piox - 201503984 \nJeannira Del Rosario Sic Menéndez - 201602434")
}

//Esta funcion va a devolver el json con la info de la pagina principal pero solo los procesos
func informacionPrincipal(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	leerProcesos()
	temporalProcesos := []structProcesos{}
	for i := 0; i < (len(listaProcesos)); i++ {
		datoPid := strconv.Itoa(listaProcesos[i].pid)
		datoNombre := listaProcesos[i].nombre
		datoUsuario := listaProcesos[i].usuario
		datoEstado := listaProcesos[i].estado
		datoPorcentaje := fmt.Sprintf("%f", listaProcesos[i].porcentajeRam)
		datoPpid := strconv.Itoa(listaProcesos[i].ppid)
		proc := structProcesos{Pid: datoPid, Nombre: datoNombre, Usuario: datoUsuario, Estado: datoEstado, PorcentajeRam: datoPorcentaje,
			Ppid: datoPpid}
		temporalProcesos = append(temporalProcesos, proc)
	}
	//fmt.Printf("Total: %v , Nuevo: %v, Temporal: %v \n", len(listaProcesos), estadisticasGenerales.total, len(temporalProcesos))
	json.NewEncoder(w).Encode(temporalProcesos)
}

//Esta funcion va a devolver el json con la informacion de la pagina de la RAM
func informacionRAM(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	leerMeminfo()
	informacion := structRam{fmt.Sprintf("%f", monitorRAM.total), fmt.Sprintf("%f", monitorRAM.consumida), fmt.Sprintf("%f", monitorRAM.porcentaje)}
	json.NewEncoder(w).Encode(informacion)
}

//Esta funcion va a devolver el json con la informacion del CPU
func informacionCPU(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	leerStat()
	informacion := structCpu{fmt.Sprintf("%f", monitorCPU)}
	json.NewEncoder(w).Encode(informacion)
}

//Esta funcion va a matar el proceso especificado
func matarProceso(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	params := mux.Vars(req)
	var valor structKill
	valor.Pid = params["id"]
	_, err := exec.Command("sh", "-c", "sudo -S pkill -SIGINT "+valor.Pid).Output()
	if err != nil {
		fmt.Printf("Error matando el proceso: %v", err)
	}
	json.NewEncoder(w).Encode(valor)
}

//FUNCIONES PRINCIPALES--------------------------------------------------------------------------------------------------------------------------------

//Esta funcion es utilizada para obtener la informacion de la ram utilizada y el total de la RAM
func leerMeminfo() {
	bytesLeidos, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		fmt.Printf("Error leyendo archivo: %v", err)
		return
	}

	contenidoArchivo := string(bytesLeidos)
	archivoCortado := strings.Split(contenidoArchivo, "\n")
	monitorRAM.total = (limpiarNumeroFloat(archivoCortado[0]) / 1000)
	libre := limpiarNumeroFloat(archivoCortado[2]) / 1000
	monitorRAM.consumida = (monitorRAM.total - libre)
	monitorRAM.porcentaje = (monitorRAM.consumida / monitorRAM.total) * 100
	//fmt.Printf("Total: %g Mb, Consumida: %g Mb, Porcentaje: %g \n", monitorRAM.total, monitorRAM.consumida, monitorRAM.porcentaje)
}

/*Esta funcion es utilizada para obtener toda la informacion de los procesos y estadisticas generales.
Hace uso del dato del total de la RAM calculado anteriormente*/
func leerProcesos() {
	archivos, err := ioutil.ReadDir("/proc/")
	if err != nil {
		log.Fatal(err)
	}

	listaProcesos = []Procesos{}
	for _, archivo := range archivos {
		if archivo.IsDir() {
			nombreCarpeta := archivo.Name()
			if buscarNumero(nombreCarpeta) {
				bytesLeidos, err := ioutil.ReadFile("/proc/" + nombreCarpeta + "/status")
				if err != nil {
					fmt.Printf("Error leyendo archivo: %v", err)
					return
				}
				contenidoArchivo := string(bytesLeidos)
				archivoCortado := strings.Split(contenidoArchivo, "\n")
				pid := limpiarNumeroInt(archivoCortado[5])
				nombre := limpiarProc(archivoCortado[0])
				usuario := limpiarUsuario(archivoCortado[8])
				estado := limpiarEstado(archivoCortado[2])
				porcentaje := (buscarPorcentajeRam(archivoCortado[17]))
				ppid := limpiarNumeroInt(archivoCortado[6])
				proceso := Procesos{pid, nombre, usuario, estado, porcentaje, ppid}
				listaProcesos = append(listaProcesos, proceso)

				leerProcesosHijos("/proc/"+nombreCarpeta+"/task", nombreCarpeta)
			}
		}
	}
}

func leerProcesosHijos(ruta string, nombre string) {
	archivos, err := ioutil.ReadDir(ruta)
	if err != nil {
		log.Fatal(err)
	}

	for _, archivo := range archivos {
		if archivo.IsDir() {
			nombreCarpeta := archivo.Name()
			if buscarNumero(nombreCarpeta) {
				if nombreCarpeta != nombre {
					bytesLeidos, err := ioutil.ReadFile(ruta + "/" + nombreCarpeta + "/status")
					if err != nil {
						fmt.Printf("Error leyendo archivo: %v", err)
						return
					}
					contenidoArchivo := string(bytesLeidos)
					archivoCortado := strings.Split(contenidoArchivo, "\n")
					pid := limpiarNumeroInt(archivoCortado[5])
					nombre := limpiarProc(archivoCortado[0])
					usuario := limpiarUsuario(archivoCortado[8])
					estado := limpiarEstado(archivoCortado[2])
					porcentaje := (buscarPorcentajeRam(archivoCortado[17]))
					ppid := limpiarNumeroInt(archivoCortado[6])
					proceso := Procesos{pid, nombre, usuario, estado, porcentaje, ppid}
					listaProcesos = append(listaProcesos, proceso)
				}
			}
		}
	}
}

//Esta funcion es utilizada para obtener la informacion del cpu
func leerStat() {
	bytesLeidos, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		fmt.Printf("Error leyendo archivo: %v", err)
		return
	}

	contenidoArchivo := string(bytesLeidos)
	archivoCortado := strings.Split(contenidoArchivo, "\n") //split por salto de linea
	//fmt.Printf("%s\n", archivoCortado[0])
	primeraLinea := strings.Split(archivoCortado[0], " ") // cpu  360012 2703 109712 948527 101360 0 10511 0 0 0

	s2 := limpiarCPU(primeraLinea[2])
	s4 := limpiarCPU(primeraLinea[4])
	s5 := limpiarCPU(primeraLinea[5])
	//fmt.Printf("s2:%s s4:%s s5:%s\n", primeraLinea[2], primeraLinea[4], primeraLinea[5])

	monitorCPU = (s2 + s4) * 100 / (s2 + s4 + s5)
	//fmt.Printf("Porcentaje: %g %% \n", monitorCPU)
}

//FUNCIONES AUXILIARES----------------------------------------------------------------------------------------------------------------------------

//Esta funcion devuelve el porcentaje de ram utilizado por cada procesos
func buscarPorcentajeRam(cadena string) float64 {
	if strings.HasPrefix(cadena, "VmSize") {
		tamanioConvertido := limpiarNumeroFloat(cadena) / 1000
		porcentaje := (tamanioConvertido / monitorRAM.total) * 100
		return porcentaje
	}
	return 0
}

//Esta funcion es utilizada para quitar la informacion innecesaria de una linea del archivo proc
func limpiarProc(cadena string) string {
	cortado := strings.Split(cadena, ":")
	sinEspacios := strings.TrimSpace(cortado[1])
	return sinEspacios
}

//Esta funcion es utilizada para quitar informacion innecesaria en una linea que contiene un numero y lo devuelve convertido en float
func limpiarNumeroFloat(cadena string) float64 {
	cortado := strings.Split(cadena, ":")
	sinEspacios := strings.TrimSpace(cortado[1])
	final := strings.Fields(sinEspacios)
	retorno, error := strconv.Atoi(final[0])
	if error != nil {
		fmt.Println("Error al convertir: ", error)
		return 0
	}
	return float64(retorno)
}

//Esta funcion es utilizada para quitar informacion innecesaria en una linea que contiene un numero y lo devuelve convertido en entero
func limpiarNumeroInt(cadena string) int {
	cortado := strings.Split(cadena, ":")
	sinEspacios := strings.TrimSpace(cortado[1])
	final := strings.Fields(sinEspacios)
	retorno, error := strconv.Atoi(final[0])
	if error != nil {
		fmt.Println("Error al convertir: ", error)
		return 0
	}
	return retorno
}

//Esta funcion busca si el nombre de una carpeta contiene un numero
func buscarNumero(cadena string) bool {
	if strings.ContainsAny(cadena, "0") {
		return true
	} else if strings.ContainsAny(cadena, "1") {
		return true
	} else if strings.ContainsAny(cadena, "2") {
		return true
	} else if strings.ContainsAny(cadena, "3") {
		return true
	} else if strings.ContainsAny(cadena, "4") {
		return true
	} else if strings.ContainsAny(cadena, "5") {
		return true
	} else if strings.ContainsAny(cadena, "6") {
		return true
	} else if strings.ContainsAny(cadena, "7") {
		return true
	} else if strings.ContainsAny(cadena, "8") {
		return true
	} else if strings.ContainsAny(cadena, "9") {
		return true
	} else {
		return false
	}
}

//Esta funcion devuelve el usuario de un proceso
func limpiarUsuario(cadena string) string {
	cortado := strings.Split(cadena, ":")
	sinEspacios := strings.TrimSpace(cortado[1])
	final := strings.Fields(sinEspacios)
	id, error := strconv.Atoi(final[0])
	if error != nil {
		fmt.Println("Error al convertir: ", error)
		return ""
	}
	for i := 0; i < len(arregloUsuarios)-1; i++ {
		if id == arregloUsuarios[i].pid {
			return arregloUsuarios[i].nombre
		}
	}
	return ""
}

//Esta funcion devuelve solo la letra del estado en el que esta un proceso
func limpiarEstado(cadena string) string {
	cortado := strings.Split(cadena, ":")
	sinEspacios := strings.TrimSpace(cortado[1])
	final := strings.Fields(sinEspacios)
	return final[0]
}

//Esta funcion obtiene los usuarios del sistema y su id respectivo
func obtenerUsuarios() {
	nombreArchivo := "/etc/passwd"
	bytesLeidos, err := ioutil.ReadFile(nombreArchivo)
	if err != nil {
		fmt.Printf("Error leyendo archivo: %v", err)
	}

	contenido := string(bytesLeidos)
	arreglo := strings.Split(contenido, "\n")
	for i := 0; i < len(arreglo)-1; i++ {
		cortada := strings.Split(arreglo[i], ":")
		pid, error := strconv.Atoi(cortada[2])
		if error != nil {
			fmt.Println("Error al convertir: ", error)
		}
		nombre := cortada[0]
		usuario := Usuario{pid, nombre}
		arregloUsuarios = append(arregloUsuarios, usuario)

	}
}

//Esta funcion es utilizada para limpiar el dato del cpu
func limpiarCPU(cadena string) float64 {
	retorno, error := strconv.Atoi(cadena)
	if error != nil {
		fmt.Println("Error al convertir: ", error)
		return 0
	}
	return float64(retorno)
}
