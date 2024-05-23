package dal

import (
	"github.com/promartingranja/SISNOVA-PERS-CONSOLIDADO/APP/models"
)

var qTest = models.Query{
	SQLGet:  "select ID, TITLE from Question",
}

//EmpresaDAL acceso a datos
type ConsolidadoDAL struct{}

func (ConsolidadoDAL) GetAllByID(id int64) (registro models.Consolidado, err error) {

	bd, err := GetSQLServer()
	if err != nil {
		return registro, err
	}
	defer bd.Close()

	stmt, err := bd.Prepare(qTest.SQLGet)
	if err != nil {
		return registro, err
	}
	defer stmt.Close()

	fila, err := stmt.Query()
	if err != nil {
		return registro, err
	}
	defer fila.Close()

	fila.Next()
	err = fila.Scan(&registro.ID, &registro.Nit)
	registro.Direccion = "sqlserver"

	return registro, err
}
