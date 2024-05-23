package models

//Configuracion
type Configuracion struct {
	Port int `json:"port"`
	NumRows int `json:"num_rows"`
	PatchImg string `json:"patch_img"`
}
