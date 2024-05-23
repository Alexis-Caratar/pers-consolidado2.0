package dal

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/promartingranja/SISNOVA-PERS-CONSOLIDADO/APP/models"
)

var qQuestion = models.Query{

	SQLGet: `SELECT Id, Title, Id_SurveyTemplate, OrderQuestion, isnull(EnumLabel,''), Id_ControlGroup 
			 FROM Question 
			 WHERE Id = @id`,

	SQLGetAll: `SELECT Id, Title, Id_SurveyTemplate, OrderQuestion, isnull(EnumLabel,''), Id_ControlGroup 
                FROM Question 
                WHERE Id_SurveyTemplate = @id order by OrderQuestion ASC `,

	SQLGetConsolid: `SELECT Q.ID, Q.Title, Q.OrderQuestion, isnull(Q.EnumLabel,'') as EnumLabel, Q.Id_ControlGroup, R.id, isnull(R.Title,''), R.OrderResponse, R.Index_Row, R.Index_Column, isnull(R.Value,''), R.Id_Response_Type, isnull(R.EnumLabel,''), isnull(R.AditionalResponse,''), isnull(Q.Is_Images, '') as is_image
					 FROM   Response R inner join
						    Question Q on R.Id_Question = q.Id
					 WHERE  Q.Id_SurveyTemplate = @id
					 ORDER BY Q.OrderQuestion, Q.ID, OrderResponse ASC`,
}

//EmpresaDAL acceso a datos
type QuestionDAL struct{}

func (QuestionDAL) GetBy(id int64) (mdl models.Question, err error) {
	ctx := context.Background()

	bd, err := GetSQLServer()
	if err != nil {
		return mdl, err
	}
	defer bd.Close()

	stmt, err := bd.Prepare(qQuestion.SQLGet)
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

	var fila models.Question
	err = filas.Scan(&fila.ID, &fila.Title, &fila.IdSurveyTemplate, &fila.OrderQuestion, &fila.EnumLabel, &fila.IdControlGroup)
	if(err != nil){
		return mdl, err
	}

	return mdl, err
}

func (QuestionDAL) GetAllByIDSurveyTemplate(id int64) (mdl []models.Question, err error) {
	ctx := context.Background()

	bd, err := GetSQLServer()
	if err != nil {
		return mdl, err
	}
	defer bd.Close()

	stmt, err := bd.Prepare(qQuestion.SQLGetAll)
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

		var fila models.Question
		err = filas.Scan(&fila.ID, &fila.Title, &fila.IdSurveyTemplate, &fila.OrderQuestion, &fila.EnumLabel, &fila.IdControlGroup)
		if(err != nil){
			fmt.Println(err)
			return nil, err
		}
		mdl = append(mdl, fila)
	}

	return mdl, err
}

func (QuestionDAL) GetAllByIDSurveyTemplate2(id int64) (mdl []models.ResponseQuestion, err error) {
	ctx := context.Background()

	bd, err := GetSQLServer()
	if err != nil {
		return mdl, err
	}
	defer bd.Close()

	stmt, err := bd.Prepare(qQuestion.SQLGetConsolid)
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

		var fila models.ResponseQuestion
		err = filas.Scan(&fila.QID, &fila.QTitle, &fila.OrderQuestion, &fila.QEnumLabel, &fila.IdControlGroup, &fila.RID, &fila.RTitle, &fila.OrderResponse, &fila.IndexRow, &fila.IndexColumn, &fila.Value, &fila.IdResponseType, &fila.REnumLabel, &fila.AditionalResponse, &fila.IsImages)
		if(err != nil){
			fmt.Println(err)
			return nil, err
		}
		mdl = append(mdl, fila)
	}

	return mdl, err
}
