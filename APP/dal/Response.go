package dal

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/promartingranja/SISNOVA-PERS-CONSOLIDADO/APP/models"
)

var qResponse = models.Query{

	SQLGet: `SELECT id, isnull(Title,''), isnull(Value,0), Id_Question, Id_Response_Type, OrderResponse, isnull(EnumLabel, '') 
			 FROM Response
			 WHERE Id = @id`,

	SQLGetAll: `SELECT id, isnull(Title,''), isnull(Value,0), Id_Question, Id_Response_Type, OrderResponse, isnull(EnumLabel, '')  
                FROM Response
                WHERE Id_Question = @id order by OrderResponse ASC `,
}

//EmpresaDAL acceso a datos
type ResponseDAL struct{}

func (ResponseDAL) GetBy(id int64) (mdl models.Response, err error) {
	ctx := context.Background()

	bd, err := GetSQLServer()
	if err != nil {
		return mdl, err
	}
	defer bd.Close()

	stmt, err := bd.Prepare(qResponse.SQLGet)
	if err != nil {
		return mdl, err
	}
	defer stmt.Close()

	filas, err := stmt.QueryContext(
		ctx,
		sql.Named("id", id),
	)
	if err != nil {
		return mdl, err
	}
	defer filas.Close()

	filas.Next()

	err = filas.Scan(&mdl.ID, &mdl.Title, &mdl.Value, &mdl.IdQuestion, &mdl.IdResponseType, &mdl.OrdenResponse, &mdl.EnumLabel)
	if(err != nil){
		return mdl, err
	}

	//mdl.Question, _ = QuestionDAL{}.GetBy(mdl.Question.ID)
	return mdl, err
}

func (ResponseDAL) GetAllByIDQuestion(id int64) (mdl []models.Response, err error) {

	ctx := context.Background()

	bd, err := GetSQLServer()
	if err != nil {
		return mdl, err
	}
	defer bd.Close()

	stmt, err := bd.Prepare(qResponse.SQLGetAll)
	if err != nil {
		return mdl, err
	}
	defer stmt.Close()

	filas, err := stmt.QueryContext(
		ctx,
		sql.Named("id", id),
	)
	if err != nil {
		return mdl, err
	}
	defer filas.Close()

	for filas.Next() {

		var fila models.Response
		err = filas.Scan(&fila.ID, &fila.Title, &fila.Value, &fila.IdQuestion, &fila.IdResponseType, &fila.OrdenResponse, &fila.EnumLabel)
		if(err != nil){
			fmt.Println("rp",err)
			return nil, err
		}

		//fila.Question, _ = QuestionDAL{}.GetBy(fila.Question.ID)
		mdl = append(mdl, fila)
	}

	return mdl, err
}
