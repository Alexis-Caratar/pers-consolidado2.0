package models

//infoApp modelo
type infoApp struct {
	Nombre       string  `json:"nombre,omitempty"`
	Slogan       string  `json:"slogan,omitempty"`
	Version      string  `json:"version,omitempty"`
	CookieSesion string  `json:"cookie_sesion,omitempty"`
}

//App modelo para la informacion de la aplicacion
var App infoApp

func init() {
	App.Nombre = "CONSOLIDADO PERS"
	App.Slogan = "SISNOVA -PERS"
	App.Version = "1.0"
	App.CookieSesion = "_TSU_"
}
