package dal

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/promartingranja/SISNOVA-PERS-CONSOLIDADO/APP/models"
)

var qQuestionImage = models.Query{

	SQLGetAll: `SELECT * FROM QuestionImage 
                WHERE 
                Id_Question in(SELECT id 
                               FROM Question 
                               WHERE Id_SurveyTemplate = @id AND Is_Images = 'true') 
				ORDER BY Id_Survey, Id_Question,  DateCreate ASC`,
}

//EmpresaDAL acceso a datos
type QuestionImageDAL struct{}

func (QuestionImageDAL) GetAllByIDSurveyTemplate(id int64) (mdl []models.QuestionImage, err error) {
	ctx := context.Background()

	bd, err := GetSQLServer()
	if err != nil {
		return mdl, err
	}
	defer bd.Close()

	stmt, err := bd.Prepare(qQuestionImage.SQLGetAll)
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

		var fila models.QuestionImage
		err = filas.Scan(&fila.ID, &fila.IDQuestion, &fila.IdSurvey, &fila.Patch, &fila.DateCreate, &fila.IdMobileArtefact)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		mdl = append(mdl, fila)
	}

	return mdl, err
}

