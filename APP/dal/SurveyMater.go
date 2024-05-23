package dal

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/promartingranja/SISNOVA-PERS-CONSOLIDADO/APP/models"
)

var qSurvey = models.Query{
	SQLGetAll: `SELECT Id, Id_Artefact, Id_Survey, Latitude, Longitude, Altitude, DateFinish, DateCreate FROM SurveyMaster WHERE Id_SurveyTemplate = @id`,
	SQLGetAllRango: `SELECT Id, Id_Artefact, Id_Survey, Latitude, Longitude, Altitude, DateFinish, DateCreate
					 FROM SurveyMaster 
					 WHERE Id_SurveyTemplate = @id and ( DateCreate  >= @fechainicio and  DateCreate <= @fechafin )`,
}

//SurveyDAL acceso a datos
type SurveyMaterDAL struct{}

func (SurveyMaterDAL) GetAllByIdSurveryTemplate(idSurveryTemplate int64) (mdl []models.SurveyMaster, err error) {

	ctx := context.Background()

	bd, err := GetSQLServer()
	if err != nil {
		fmt.Println(err)
		return mdl, err
	}
	defer bd.Close()

	stmt, err := bd.Prepare(qSurvey.SQLGetAll)
	if err != nil {
		return mdl, err
	}
	defer stmt.Close()

	filas, err := stmt.QueryContext(
		ctx,
		sql.Named("id", idSurveryTemplate),
	)
	if err != nil {
		return mdl, err
	}
	defer filas.Close()

	for filas.Next() {

		var fila models.SurveyMaster
		err = filas.Scan(&fila.ID, &fila.IdArtefact, &fila.IdSurvey, &fila.Latitude, &fila.Longitude, &fila.Altitude, &fila.FechaFinish, &fila.FechaCreate)
		if err != nil {
			fmt.Println("rp", err)
			return nil, err
		}

		s := strings.Split(fila.FechaFinish, " ")
		fila.FechaFinish = s[0]
		fila.Hora = s[1]

		//fila.Question, _ = QuestionDAL{}.GetBy(fila.Question.ID)
		mdl = append(mdl, fila)
	}

	return mdl, err
}

func (SurveyMaterDAL) GetAllByIdSurveryTemplateRango(idSurveryTemplate int64, fechaInicio string, fechafin string) (mdl []models.SurveyMaster, err error) {

	ctx := context.Background()

	bd, err := GetSQLServer()
	if err != nil {
		return mdl, err
	}
	defer bd.Close()

	stmt, err := bd.Prepare(qSurvey.SQLGetAll)
	if err != nil {
		fmt.Println(err)
		return mdl, err
	}
	defer stmt.Close()

	filas, err := stmt.QueryContext(
		ctx,
		sql.Named("id", idSurveryTemplate),
		sql.Named("fechainicio", fechaInicio),
		sql.Named("fechafin", fechafin),
	)
	if err != nil {
		fmt.Println(err)
		return mdl, err
	}
	defer filas.Close()

	for filas.Next() {

		var fila models.SurveyMaster
		err = filas.Scan(&fila.ID, &fila.IdArtefact, &fila.IdSurvey, &fila.Latitude, &fila.Longitude, &fila.Altitude, &fila.FechaFinish, &fila.FechaCreate)
		if err != nil {
			fmt.Println("rp", err)
			return nil, err
		}

		s := strings.Split(fila.FechaFinish, " ")
		fila.FechaFinish = s[0]
		fila.Hora = s[1]

		//fila.Question, _ = QuestionDAL{}.GetBy(fila.Question.ID)
		mdl = append(mdl, fila)
	}

	return mdl, err
}
