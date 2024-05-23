package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/husobee/vestigo"
	"io/ioutil"
	"net/http"
)

//EmpresaHandler Handler para el modulo Empresa
type ArchivoConValidado struct{}

func (ArchivoConValidado) PostArchvioConValidado(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	clave := vestigo.Param(r, "id")
	//ID_ST, _ := strconv.ParseInt(clave, 10, 32)

	file, _, err := r.FormFile("file")
	if err != nil {
		fmt.Println("1->",err)
		w.Write([]byte(fmt.Sprintf("%d", err)))
		return
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("2->",err)
		w.Write([]byte(fmt.Sprintf("%d", err)))
	}

	err = ioutil.WriteFile("../../PUBLIC/docs/" + clave + ".xlsx", data, 0666)
	if err != nil {
		fmt.Println("3->",err)
		w.Write([]byte(fmt.Sprintf("%d", err)))
	}
	//fmt.Println(header)

	js, err := json.Marshal(clave)
	w.Write(js)

}

func (ArchivoConValidado) GetConValidados(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	strCadenaval := "../../PUBLIC/docs/"

	archivos, err := ioutil.ReadDir(strCadenaval)
	if err != nil{
		//fmt.Println("1->",err)
		w.Write([]byte(fmt.Sprintf("%d", err)))
		return
	}

	var data []string

	for _, archivo := range archivos {
		//fmt.Println("Nombre:", archivo.Name())
		data = append(data, archivo.Name())
	}

	js, err := json.Marshal(data)
	w.Write(js)

}
