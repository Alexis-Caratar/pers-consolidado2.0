package dal

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/promartingranja/SISNOVA-PERS-CONSOLIDADO/APP/models"
)

var qSurveyResponse = models.Query{

	SQLGet: `SELECT Id, Id_Survey, ID_Response, Value 
			 FROM SurveyResponse 
			 WHERE Id = @id`,

	//SQLGetAll: `SELECT Id, Id_Survey, ID_Response, Value
    //            FROM SurveyResponse
	//		    WHERE ID_Response = @id`,

	SQLGetAll: `select sr.id, sr.id_survey, sr.id_quesion, sr.id_response, value, sr.aditionalresponse
                from SurveyResponse sr inner join Question qu ON sr.Id_Quesion = qu.Id
				where qu.Id_SurveyTemplate = @id
				order by sr.Id_Survey, qu.OrderQuestion, qu.Id`,

	SQLGetAllRango: `select sr.id, sr.id_survey, sr.id_quesion, sr.id_response, value, sr.aditionalresponse
					 from SurveyResponse sr inner join Question qu ON sr.Id_Quesion = qu.Id inner join SurveyMaster sm ON sr.Id_Survey = sm.id
					 where qu.Id_SurveyTemplate = @id and ( sm.DateCreate  >= @fechainicio and  sm.DateCreate <= @fechafin )
					 order by sr.Id_Survey, qu.OrderQuestion, qu.Id`,

	SQLGetConsolid: `select distinct Id_Survey from SurveyResponse
					 where ID_Response in(select Id from Response
										  where Id_Question in(select id from Question where Id_SurveyTemplate = @id)
 									     )
					 order by Id_Survey`,
}

//EmpresaDAL acceso a datos
type SurveyResponseDAL struct{}

//func (SurveyResponseDAL) GetBy(id int64) (mdl models.SurveryResponse, err error) {
//	ctx := context.Background()
//
//	bd, err := GetSQLServer()
//	if err != nil {
//		return mdl, err
//	}
//	defer bd.Close()
//
//	stmt, err := bd.Prepare(qSurveyResponse.SQLGet)
//	if err != nil {
//		return mdl, err
//	}
//	defer stmt.Close()
//
//	filas, err := stmt.QueryContext(
//		ctx,
//		sql.Named("id", id),
//	)
//	if err != nil {
//		return mdl, err
//	}
//	defer filas.Close()
//
//	filas.Next()
//
//	var fila models.SurveryResponse
//	err = filas.Scan(&fila.ID, &fila.IdSurvey, &fila.IdResponse, &fila.Value)
//	if(err != nil){
//		return mdl, err
//	}
//
//	return mdl, err
//}

func (SurveyResponseDAL) GetAllBySurveyResponse(idST int64) (mdl []models.SurveryResponse, err error) {

	ctx := context.Background()

	bd, err := GetSQLServer()
	if err != nil {
		return mdl, err
	}
	defer bd.Close()

	stmt, err := bd.Prepare(qSurveyResponse.SQLGetAll)
	if err != nil {
		return mdl, err
	}
	defer stmt.Close()

	filas, err := stmt.QueryContext(
		ctx,
		sql.Named("id", idST),
	)
	if err != nil {
		return mdl, err
	}
	defer filas.Close()

	for filas.Next() {

		var fila models.SurveryResponse
		err = filas.Scan(&fila.ID, &fila.IdSurvey, &fila.IdQuestion, &fila.IdResponse, &fila.Value, &fila.AditionalResponse)
		if(err != nil){
			return nil, err
		}
		mdl = append(mdl, fila)
	}

	return mdl, err
}

func (SurveyResponseDAL) GetAllBySurveyResponseRango(idST int64, fechaInicio string, fechaFin string) (mdl []models.SurveryResponse, err error) {

	ctx := context.Background()

	bd, err := GetSQLServer()
	if err != nil {
		fmt.Println(err)
		return mdl, err
	}
	defer bd.Close()

	stmt, err := bd.Prepare(qSurveyResponse.SQLGetAllRango)
	if err != nil {
		fmt.Println(err)
		return mdl, err
	}
	defer stmt.Close()

	filas, err := stmt.QueryContext(
		ctx,
		sql.Named("id", idST),
		sql.Named("fechainicio", fechaInicio),
		sql.Named("fechafin", fechaFin),
	)
	if err != nil {
		return mdl, err
	}
	defer filas.Close()

	for filas.Next() {

		var fila models.SurveryResponse
		err = filas.Scan(&fila.ID, &fila.IdSurvey, &fila.IdQuestion, &fila.IdResponse, &fila.Value, &fila.AditionalResponse)
		if(err != nil){
			return nil, err
		}
		mdl = append(mdl, fila)
	}

	return mdl, err
}

//func (SurveyResponseDAL) GetAllByConsolid(id int64) (mdl []models.SurveryResponse, err error) {
//
//	ctx := context.Background()
//
//	bd, err := GetSQLServer()
//	if err != nil {
//		return mdl, err
//	}
//	defer bd.Close()
//
//	stmt, err := bd.Prepare(qSurveyResponse.SQLGetConsolid)
//	if err != nil {
//		return mdl, err
//	}
//	defer stmt.Close()
//
//	filas, err := stmt.QueryContext(
//		ctx,
//		sql.Named("id", id),
//	)
//	if err != nil {
//		return mdl, err
//	}
//	defer filas.Close()
//
//	for filas.Next() {
//
//		var fila models.SurveryResponse
//		err = filas.Scan(&fila.IdSurvey)
//		if(err != nil){
//			fmt.Println("sur",err)
//			return nil, err
//		}
//		mdl = append(mdl, fila)
//	}
//
//	return mdl, err
//}
