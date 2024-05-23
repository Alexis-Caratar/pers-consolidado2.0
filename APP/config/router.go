package config

import (
	"github.com/husobee/vestigo"
	"github.com/promartingranja/SISNOVA-PERS-CONSOLIDADO/APP/handlers"
)

func Load(router *vestigo.Router) {

	//CONSOLIDADO
	iniPath := "/consolidado"
	router.Get(iniPath+"/:id", handlers.ConsolidadoHandler{}.GetConsolidado)
	router.Get(iniPath+"/:id/:fechainicio/:fechafin", handlers.ConsolidadoHandler{}.GetConsolidado)

	iniPath = "/subir-con-validado"
	router.Post(iniPath+"/:id", handlers.ArchivoConValidado{}.PostArchvioConValidado)

	iniPath = "/traer-validados"
	router.Get(iniPath , handlers.ArchivoConValidado{}.GetConValidados)

	iniPath = "/subir-consolidado"
	router.Post(iniPath, handlers.SubiendoValidado{}.ArchivoSubiendo)

	iniPath = "/lsitado-pers"
	router.Get(iniPath , handlers.PersColombia{}.GetListadoPERS)
	router.Get(iniPath +"/region/:region" , handlers.PersColombia{}.GetListadoPERSRegionConsolidado)

}
