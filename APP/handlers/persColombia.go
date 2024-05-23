package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/husobee/vestigo"
	"github.com/promartingranja/SISNOVA-PERS-CONSOLIDADO/APP/dal"
	"net/http"
	"strconv"
)

type PersColombia struct{}

func (PersColombia) GetListadoPERS(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	data, err := dal.PersColombiaDAL{}.GetListadoPERS()
	if err != nil {
		w.Write([]byte(fmt.Sprintf("%d", err)))
		return
	}

	js, err := json.Marshal(data)
	w.Write(js)

}

func (PersColombia) GetListadoPERSRegionConsolidado(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	nombrePers := vestigo.Param(r, "region")

	f := excelize.NewFile()

	conf, err := dal.GetConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Creando encabezados")
	for i := 1; i <= conf.NumRows; i++ {
		columnaLetra, _ := excelize.ColumnNumberToName(i)
		_ = f.SetCellValue("Sheet1", columnaLetra+"1", "P"+strconv.Itoa(i))
	}

	fmt.Println("Insertando filas de PERS->C")
	dataPERSC, err := dal.PersColombiaDAL{}.GetListadoPERSCxRegion(nombrePers)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("%d", err)))
		return
	}

	for i, fila := range dataPERSC {

		for encavezado, columna := range fila {
			auxEncabezado := encavezado[1:len(encavezado)]
			numeroEncabezado, _ := strconv.Atoi(auxEncabezado)
			columnaLetra, _ := excelize.ColumnNumberToName(numeroEncabezado)
			numFila := i + 2
			_ = f.SetCellValue("Sheet1", columnaLetra+strconv.Itoa(numFila), columna)
		}
	}

	fmt.Println("Insertando filas de PERS->S")
	dataPERSS, err := dal.PersColombiaDAL{}.GetListadoPERSSxRegion(nombrePers)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("%d", err)))
		return
	}

	for i, fila := range dataPERSS {

		for encavezado, columna := range fila {
			auxEncabezado := encavezado[1:len(encavezado)]
			numeroEncabezado, _ := strconv.Atoi(auxEncabezado)
			columnaLetra, _ := excelize.ColumnNumberToName(numeroEncabezado)
			numFila := i + 2

			colExeption := []int{829, 841, 853, 865, 877, 889, 901, 913, 925, 937, 949, 961, 973, 985, 1056, 1057, 1058, 1059, 1060}

			noExiste := false
			if numeroEncabezado >= 829 {
				noExiste = existeEnArreglo(colExeption, numeroEncabezado)
			}

			if noExiste == false {
				_ = f.SetCellValue("Sheet1", columnaLetra+strconv.Itoa(numFila), columna)
			}
		}
	}

	fmt.Println("Creando archivo y finalizando")
	dirGuardar := "../../PUBLIC/doc_pers_colombia/" + nombrePers + ".xlsx"
	err = f.SaveAs(dirGuardar)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Add("Content-Disposition", "Attachment")
	w.Header().Add("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	http.ServeFile(w, r, dirGuardar)
}

func existeEnArreglo(arreglo []int, busqueda int) bool {
	for _, numero := range arreglo {
		if numero == busqueda {
			return true
		}
	}
	return false
}
