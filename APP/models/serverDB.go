package models

//ServerDB modelo para el servidor de la base de datos
type ServerDB struct {
	Servidor string `json:"servidor,omitempty"`
	Puerto   string `json:"puerto,omitempty"`
	BD       string `json:"bd,omitempty"`
	Usuario  string `json:"usuario,omitempty"`
	Clave    string `json:"clave,omitempty"`
}
