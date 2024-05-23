package dal

import (
	"context"
	"database/sql"
	"github.com/promartingranja/SISNOVA-PERS-CONSOLIDADO/APP/models"
)

var qPersColombia = models.Query{
	SQLGetAll: "select p1 from PERSS WHERE p1 != '' GROUP BY p1",
	SQLGetPERSC: "select * from PERSC where p1 = @region",
	SQLGetPERSS: "select * from PERSS where p1 = @region",
}

type PersColombiaDAL struct {}

func (PersColombiaDAL) GetListadoPERS()(mdl []models.Pers, err error){

	ctx := context.Background()

	bd, err := GetSQLServerUpload()
	if err != nil {
		return mdl, err
	}
	defer bd.Close()

	stmt, err := bd.Prepare(qPersColombia.SQLGetAll)
	if err != nil {
		return mdl, err
	}
	defer stmt.Close()

	filas, err := stmt.QueryContext(
		ctx,
	)
	if err != nil {
		return mdl, err
	}
	defer filas.Close()

	for filas.Next() {

		var fila models.Pers
		err = filas.Scan(&fila.Nombre)
		if err != nil{
			return nil, err
		}
		mdl = append(mdl, fila)
	}

	return mdl, err
}

func (PersColombiaDAL) GetListadoPERSCxRegion(region string)(mdl []map[string]interface{}, err error){

	ctx := context.Background()

	bd, err := GetSQLServerUpload()
	if err != nil {
		return mdl, err
	}
	defer bd.Close()

	stmt, err := bd.Prepare(qPersColombia.SQLGetPERSC)
	if err != nil {
		return mdl, err
	}

	defer stmt.Close()
	filas, err := stmt.QueryContext(
		ctx,
		sql.Named("region", region),
	)

	if err != nil {
		return mdl, err
	}
	defer filas.Close()

	cols, _ := filas.Columns()
	for filas.Next() {
		mdlAux := parseToMap(filas, cols)
		mdl = append(mdl, mdlAux)
	}

	return mdl, err
}

func (PersColombiaDAL) GetListadoPERSSxRegion(region string)(mdl []map[string]interface{}, err error){

	ctx := context.Background()

	bd, err := GetSQLServerUpload()
	if err != nil {
		return mdl, err
	}
	defer bd.Close()

	stmt, err := bd.Prepare(qPersColombia.SQLGetPERSS)
	if err != nil {
		return mdl, err
	}

	defer stmt.Close()
	filas, err := stmt.QueryContext(
		ctx,
		sql.Named("region", region),
	)

	if err != nil {
		return mdl, err
	}
	defer filas.Close()

	cols, _ := filas.Columns()
	for filas.Next() {
		mdlAux := parseToMap(filas, cols)
		mdl = append(mdl, mdlAux)
	}

	return mdl, err
}

func parseToMap(rows *sql.Rows, cols []string) map[string]interface{} {
	values := make([]interface{}, len(cols))
	pointers := make([]interface{}, len(cols))
	for i := range values {
		pointers[i] = &values[i]
	}

	if err := rows.Scan(pointers...); err != nil {
		panic(err)
	}

	m := make(map[string]interface{})
	for i, colName := range cols {
		if values[i] == nil {
			m[colName] = nil
		} else {
			m[colName] = values[i]
		}
	}
	return m
}
