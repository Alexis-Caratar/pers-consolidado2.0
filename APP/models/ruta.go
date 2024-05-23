package models

//Rutas modelo
type infRutas struct {
	Layout string `json:"layout,omitempty"`
	Shared string `json:"shared,omitempty"`
}

//Rutas modelo para la Rutas Publicas
var Rutas infRutas

func init() {
	//variables para rutas definidas
	Rutas.Layout = "../tpl/layout/"
	Rutas.Shared = "../tpl/shared/"
}
