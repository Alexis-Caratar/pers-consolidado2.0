package dal

import (
	"context"
	"fmt"
)

//SurveyDAL acceso a datos
type SubiendoValidadoDAL struct{}

func (SubiendoValidadoDAL) GetMaxPERSC() (max string, err error) {

	ctx := context.Background()

	bd, err := GetSQLServerUpload()
	if err != nil {
		return max, err
	}
	defer bd.Close()

	stmt, err := bd.Prepare("SELECT MAX(CAST(p0 AS INTEGER)) as max FROM PERSC")
	if err != nil {
		return max, err
	}
	defer stmt.Close()

	filas, err := stmt.QueryContext(
		ctx,
	)
	if err != nil {
		return max, err
	}
	defer filas.Close()

	for filas.Next() {

		err = filas.Scan(&max)
		if err != nil {
			fmt.Println("rp", err)
			return max, err
		}

	}

	return max, err
}

func (SubiendoValidadoDAL) CreatePERS(cadena string) (err error) {

	//ctx := context.Background()

	bd, err := GetSQLServerUpload()
	defer bd.Close()
	if err != nil {
		return err
	}

	stmt, err := bd.Prepare(cadena)
	defer stmt.Close()
	if err != nil {
		return err
	}

	filas, err := stmt.Query()
	defer filas.Close()
	if err != nil {
		fmt.Println("Error CREANDO SQL ----->>", cadena)
		return err
	}

	return err
}
