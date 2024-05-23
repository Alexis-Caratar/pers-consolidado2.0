package models

//Empresa modelo
type Consolidado struct {
	ID                int64  `json:"id,string"`
	Nit               string `json:"nit"`
	Nombre            string `json:"nombre"`
	Estado            string `json:"estado"`
	Telefono          string `json:"telefono"`
	Direccion         string `json:"direccion"`
	Correo            string `json:"correo"`
	Logo              string `json:"logo"`
	MapaProceso       string `json:"mapa_procesos"`
	RepresentateLegal string `json:"representate_legal"`
	Mision            string `json:"mision"`
	Vision            string `json:"vision"`
	Politicas         string `json:"politicas"`
	Objetivos         string `json:"objetivos"`
	Estrategias       string `json:"estrategias"`
	Organigrama       string `json:"organigrama"`
}
