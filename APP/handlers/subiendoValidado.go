package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/promartingranja/SISNOVA-PERS-CONSOLIDADO/APP/dal"
)

//EmpresaHandler Handler para el modulo Empresa
type SubiendoValidado struct{}

var wg, wgLotes sync.WaitGroup

func (SubiendoValidado) ArchivoSubiendo(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	// obteniendo file
	file, _, err := r.FormFile("file")
	if err != nil {
		fmt.Println(err)
		w.Write([]byte(fmt.Sprintf("%d", err)))
		return
	}

	//leyendo file y pasandolos a una variable
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		w.Write([]byte(fmt.Sprintf("%d", err)))
		return
	}

	// guardando file
	err = ioutil.WriteFile("../../PUBLIC/docs/tmp/subirArchivo.xlsx", data, 0666)
	if err != nil {
		fmt.Println(err)
		w.Write([]byte(fmt.Sprintf("%d", err)))
		return
	}

	// ############## MIGRACION A LA BD

	// se trae el archivo .xlsx
	f, err := excelize.OpenFile("../../PUBLIC/docs/tmp/subirArchivo.xlsx")
	if err != nil {
		w.Write([]byte(fmt.Sprintf("%d", err)))
		return
	}

	// se lee la Data
	rows, err := f.GetRows("Data")
	if err != nil {
		w.Write([]byte(fmt.Sprintf("%d", err)))
		return
	}

	fmt.Println("## Inicio Subida de Datos ##")

	//fmt.Println("1.  Validar Encabezados")
	//if len(rows[0]) > 1062 {
	//
	//	res := ValidarCrearEncabezados(rows[0])
	//	if !res {
	//		js, _ := json.Marshal("Error al Validar Encabezados")
	//		w.Write(js)
	//		return
	//	}
	//}

	// obteniendo el valor maximo de la tabla PERSC
	auxMax, err := dal.SubiendoValidadoDAL{}.GetMaxPERSC()
	if err != nil {
		w.Write([]byte(fmt.Sprintf("%d", err)))
		return
	}

	maximoID_p0, _ := strconv.Atoi(auxMax)
	fmt.Println("maximo de la tabla ", maximoID_p0)

	// obteniendo el tamaño de filas del archivo
	max := len(rows)
	var fin int

	fmt.Println("1. Insertando Filas en Pers, total de filas ", len(rows))

	var rangoLote int = 50
	var rutinas int = max / rangoLote
	total := max - (rutinas * rangoLote)
	if total > 0 {
		rutinas++
	}

	wgLotes.Add(rutinas)

	for i := 0; i < max; i = i + rangoLote {

		if max <= i+rangoLote {
			fin = max
		} else {
			fin = i + rangoLote
		}

		go func(des int, fin int) {

			// variables para pruebas
			//des := i

			if des == 0 {
				des = 1
			}

			consecutivo := maximoID_p0 + des
			cadenaPERSC, cadenaPERSS := CadenaPERSC(rows, consecutivo, des, fin)

			auxCadena := "INSERT INTO PERSC VALUES " + cadenaPERSC + "INSERT INTO PERSS VALUES " + cadenaPERSS

			//Insertando en PERSC
			//e := dal.SubiendoValidadoDAL{}.CreatePERS("INSERT INTO PERSC VALUES " + cadenaPERSC)
			e := dal.SubiendoValidadoDAL{}.CreatePERS(auxCadena)
			if e != nil {
				fmt.Println("PERSC/PERSS error -> linea desde ", des, " - hasta ", fin, " msj -> ", e)
			} else {
				fmt.Println("PERSC/PERSS lote desde ", des, " hasta ", fin)
			}

			// Insertando en PERSS
			//e = dal.SubiendoValidadoDAL{}.CreatePERS("INSERT INTO PERSS VALUES " + cadenaPERSS)
			//if e != nil {
			//	fmt.Println("PERSS error -> linea desde ", des, " - hasta ", fin , " msj -> ", e)
			//} else {
			//	fmt.Println("PERSS lote desde ", des , " hasta " , fin )
			//}

			defer wgLotes.Done()

		}(i, fin)

	}

	wgLotes.Wait()

	fmt.Println("Fin Subida de Datos")

	// ############## FIN MIGRACION A LA BD

	js, err := json.Marshal(true)
	w.Write(js)
}

// ## funciones adicionales

func CadenaPERSC(rows [][]string, consecutivo int, desde int, hasta int) (string, string) {

	cadenaPERSC := ""
	cadenaPERSS := ""
	//wg.Add(hasta - desde)

	//for f, row := range rows {
	for f := desde; f < hasta; f++ {

		//fmt.Println("--> " , f)

		if f == 0 {
			continue
		}

		consecutivo++

		// #
		row := rows[f]
		con := consecutivo
		// # /

		//go func(row []string, f int, con int) {
		//fmt.Println("Generando cadena de la fila ", f)

		auxCadPERSC, auxCadePERSS := filaCadenaPERSC(row, con)

		cadenaPERSC += auxCadPERSC
		cadenaPERSS += auxCadePERSS
		//defer wg.Done()

		//}(rows[f], f, consecutivo)

	}

	//wg.Wait()

	if len(cadenaPERSC) > 0 {
		cadenaPERSC = cadenaPERSC[0:len(cadenaPERSC)-1] + ";"
		cadenaPERSS = cadenaPERSS[0:len(cadenaPERSS)-1] + ";"
	} else {
		fmt.Println("error --> ", cadenaPERSC)
		fmt.Println("error --> ", cadenaPERSS)
	}

	return cadenaPERSC, cadenaPERSS
}

func filaCadenaPERSC(row []string, con int) (cadenaPERSC string, cadenaPERSS string) {

	cadenaPERSC = "( " + strconv.Itoa(con) + ", "
	cadenaPERSS = "( " + strconv.Itoa(con) + ", "

	// PERSC

	for c, colCell := range row {

		if c > 827 {
			break
		}
		cadenaPERSC = cadenaPERSC + "'" + strings.Trim(colCell, " ") + "', "
	}

	if len(row) < 827 {

		for i := len(row); i <= 827; i++ {
			cadenaPERSC = cadenaPERSC + "'', "
		}
	}

	cadenaPERSC = cadenaPERSC[0 : len(cadenaPERSC)-2] // recorta la cadena las ultims dos posiciones
	cadenaPERSC = cadenaPERSC + "),"

	// PERSS

	for c, colCell := range row {

		if c < 7 || c > 827 {
			cadenaPERSS = cadenaPERSS + "'" + strings.Trim(colCell, " ") + "', "
		}
	}

	if len(row) < 1060 {

		inicio := 828
		if len(row) > 828 {
			inicio = len(row)
		}

		for i := inicio; i < 1061; i++ {
			cadenaPERSS = cadenaPERSS + "'', " //llena las cadenas basias faltantes
		}
	}

	cadenaPERSS = cadenaPERSS[0 : len(cadenaPERSS)-2] // recorta la cadena las ultims dos posiciones
	cadenaPERSS = cadenaPERSS + "),"

	return cadenaPERSC, cadenaPERSS
}

func ValidarCrearEncabezados(encRow []string) bool {

	res := true
	cont := 0

	for _, colCell := range encRow {

		if cont > 1061 {
			fmt.Println(strings.Trim(colCell, " "), "\t")
		}
		cont++
	}

	return res
}
