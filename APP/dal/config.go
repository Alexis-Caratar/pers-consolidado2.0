package dal

import (
	"encoding/json"
	"github.com/promartingranja/SISNOVA-PERS-CONSOLIDADO/APP/models"
	"os"
)

func GetConfig() (config models.Configuracion, err error){
	js, err := os.Open("../config/config.json")
	if err != nil {
		return config, err
	}
	defer js.Close()
	deco := json.NewDecoder(js)
	err = deco.Decode(&config)
	return config, err
}