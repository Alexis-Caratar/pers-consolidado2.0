package handlers

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/husobee/vestigo"
	"github.com/promartingranja/SISNOVA-PERS-CONSOLIDADO/APP/dal"
	"github.com/promartingranja/SISNOVA-PERS-CONSOLIDADO/APP/models"
	"net/http"
	"strconv"
)

//EmpresaHandler Handler para el modulo Empresa
type ConsolidadoHandler struct{}

//var wg sync.WaitGroup
//var wg2 sync.WaitGroup

//calculando encabezados

func GetConstructorSurvery(IdSurveyTemplate int64) (consolidadoEncabezados []models.ConsolidadoEncabezados) {

	responseQuestions, _ := dal.QuestionDAL{}.GetAllByIDSurveyTemplate2(IdSurveyTemplate)

	var Qid int64
	for _, q := range responseQuestions {

		if Qid != q.QID {

			var auxQ models.ConsolidadoEncabezados
			auxQ.Question = models.Question{
				ID:               q.QID,
				Title:            q.QTitle,
				OrderQuestion:    q.OrderQuestion,
				EnumLabel:        q.QEnumLabel,
				IdControlGroup:   q.IdControlGroup,
				IdSurveyTemplate: 0,
				IsImages:         q.IsImages,
			}

			for _, r := range responseQuestions {

				if auxQ.Question.ID == r.QID {
					auxQ.ResponseCon = append(auxQ.ResponseCon, models.Response{
						ID:                r.RID,
						Title:             r.RTitle,
						OrdenResponse:     r.OrderResponse,
						IndexRow:          r.IndexRow,
						IndexColumn:       r.IndexColumn,
						Value:             r.Value,
						IdResponseType:    r.IdResponseType,
						EnumLabel:         r.REnumLabel,
						AditionalResponse: r.AditionalResponse,
					})
				}

			}

			consolidadoEncabezados = append(consolidadoEncabezados, auxQ)
			Qid = q.QID
		}
	}

	return consolidadoEncabezados
}

//trayendo encuentas

func GetConstructorSurveryResponse(IdSurveyTemplate int64, fechaInicio string, fechaFin string) (surveysResponse []models.ConsolidadoSurveyResponse) {

	var auxSurveryResponses []models.SurveryResponse
	if fechaInicio == "" || fechaFin == "" {
		auxSurveryResponses, _ = dal.SurveyResponseDAL{}.GetAllBySurveyResponse(IdSurveyTemplate)
	} else {
		auxSurveryResponses, _ = dal.SurveyResponseDAL{}.GetAllBySurveyResponseRango(IdSurveyTemplate, fechaInicio, fechaFin)
	}

	var Sid int64

	for _, s := range auxSurveryResponses {

		if Sid != s.IdSurvey {
			surveysResponse = append(surveysResponse, models.ConsolidadoSurveyResponse{
				IDSurvey:          s.IdSurvey,
				SurveyResponseQue: GetQuetionResponse(s.IdSurvey, auxSurveryResponses),
			})
			Sid = s.IdSurvey
		}
	}

	return surveysResponse
}

func GetQuetionResponse(IdSurvey int64, sr []models.SurveryResponse) (resConQSR []models.ConsolidadoQuesSurveyResponse) {

	var Qid int64
	for _, q := range sr {

		if IdSurvey == q.IdSurvey {

			if Qid != q.IdQuestion {

				resConQSR = append(resConQSR, models.ConsolidadoQuesSurveyResponse{
					IDQuestion:  q.IdQuestion,
					ResponseCon: GetSurveryResponse(IdSurvey, q.IdQuestion, sr),
				})

				Qid = q.IdQuestion
			}
		}
	}

	return resConQSR
}

func GetSurveryResponse(IdSurvey int64, IdQuertion int64, sr []models.SurveryResponse) (resSR []models.SurveryResponse) {

	for _, q := range sr {

		if IdSurvey == q.IdSurvey && IdQuertion == q.IdQuestion {
			resSR = append(resSR, q)
		}
	}

	return resSR
}

// generar consolidado

func (ConsolidadoHandler) GetConsolidado(w http.ResponseWriter, r *http.Request) {

	clave := vestigo.Param(r, "id")
	fechaInicio := vestigo.Param(r, "fechainicio")
	fechaFin := vestigo.Param(r, "fechafin")

	ID_ST, _ := strconv.ParseInt(clave, 10, 32)

	Survery := GetConstructorSurvery(ID_ST)

	var questionFijas []models.QuestionFija
	var questionImagesHeders []models.Question

	questionImagesData, _ := dal.QuestionImageDAL{}.GetAllByIDSurveyTemplate(ID_ST) // trae de la base de datos de la tabla querestion image

	f := excelize.NewFile()

	// creando encavezados

	fmt.Println("************* CREANDO ENCABEZADOS ARCHIVO *****************")

	contCol := 1
	for _, question := range Survery {

		// agregar preguntas de tipo imangen
		if question.Question.IsImages == "true" {
			questionImagesHeders = append(questionImagesHeders, question.Question)
		}

		// LLY
		if question.Question.IdControlGroup == 1 {
			for _, qR := range question.ResponseCon {

				if qR.IdResponseType == 0 {
					continue
				}

				columnaLetra, _ := excelize.ColumnNumberToName(contCol)
				_ = f.SetCellValue("Sheet1", columnaLetra+"1", question.Question.EnumLabel+" "+question.Question.Title)
				_ = f.SetCellValue("Sheet1", columnaLetra+"2", qR.Title)
				contCol++

				//añadiendo a question fija
				auxquestionFijas := models.QuestionFija{
					RID:            qR.ID,
					QID:            question.Question.ID,
					IdControlGroup: 1,
					Columna:        columnaLetra,
				}

				if qR.AditionalResponse != "" {
					columnaLetra, _ := excelize.ColumnNumberToName(contCol)
					_ = f.SetCellValue("Sheet1", columnaLetra+"1", question.Question.EnumLabel+" "+question.Question.Title)
					_ = f.SetCellValue("Sheet1", columnaLetra+"2", qR.AditionalResponse)
					contCol++

					auxquestionFijas.IsOtter = true
					auxquestionFijas.ColumnaOtro = columnaLetra
				}

				questionFijas = append(questionFijas, auxquestionFijas)

			}
		}

		// RGR -> creo que es radiobutons
		if question.Question.IdControlGroup == 2 {

			columnaLetra, _ := excelize.ColumnNumberToName(contCol)
			auxColumnaLetra := columnaLetra
			_ = f.SetCellValue("Sheet1", columnaLetra+"1", question.Question.EnumLabel+" "+question.Question.Title)
			contCol++

			auxquestionFijas := models.QuestionFija{
				RID:            0,
				QID:            question.Question.ID,
				IdControlGroup: 2,
				Columna:        auxColumnaLetra,
			}

			for _, q := range question.ResponseCon {

				if q.AditionalResponse != "" {
					columnaLetra, _ := excelize.ColumnNumberToName(contCol)
					_ = f.SetCellValue("Sheet1", columnaLetra+"1", question.Question.EnumLabel+" "+question.Question.Title)
					_ = f.SetCellValue("Sheet1", columnaLetra+"2", q.AditionalResponse)
					contCol++

					auxquestionFijas.IsOtter = true
					auxquestionFijas.ColumnaOtro = columnaLetra
				}
			}

			//añadiendo a question fija
			questionFijas = append(questionFijas, auxquestionFijas)
		}

		// TLY -> tabla
		if question.Question.IdControlGroup == 3 {

			for _, qR := range question.ResponseCon {

				if qR.IdResponseType == 0 {
					continue
				}

				var title string = ""

				if qR.IdResponseType != 3 && question.Question.ID != 350 { // diferente a caja de chequeo

					for _, qt := range question.ResponseCon {

						if qt.IndexColumn == qR.IndexColumn && qt.IndexRow == 0 && qt.IdResponseType == 0 {
							title = qt.Title + " - "
						}

						if qt.IndexColumn == 0 && qt.IndexRow == qR.IndexRow && qt.IdResponseType == 0 {
							title += qt.Title
						}

					}

					if title != "" {
						qR.Title = title
					}
				}

				// pregunta personalizada pregunta 40
				if qR.IdResponseType != 3 && question.Question.ID == 350 { // diferente a caja de chequeo

					for _, qt := range question.ResponseCon {

						if qt.IndexColumn == qR.IndexColumn && qt.IndexRow == 1 && qt.IdResponseType == 0 {
							title += qt.Title + " "
						}

						if qt.IndexColumn == qR.IndexColumn && qt.IndexRow == 2 && qt.IdResponseType == 0 {
							title += qt.Title + " - "
						}

						if qt.IndexColumn == 0 && qt.IndexRow == qR.IndexRow && qt.IdResponseType == 0 {
							title += qt.Title + " "
						}

					}

					if title != "" {
						qR.Title = title
					}
				}

				columnaLetra, _ := excelize.ColumnNumberToName(contCol)
				_ = f.SetCellValue("Sheet1", columnaLetra+"1", question.Question.EnumLabel+" "+question.Question.Title)
				_ = f.SetCellValue("Sheet1", columnaLetra+"2", qR.Title)
				contCol++

				//añadiendo a question fija
				auxquestionFijas := models.QuestionFija{
					RID:            qR.ID,
					QID:            question.Question.ID,
					IdControlGroup: 1,
					Columna:        columnaLetra,
				}

				if qR.AditionalResponse != "" {
					columnaLetra, _ := excelize.ColumnNumberToName(contCol)
					_ = f.SetCellValue("Sheet1", columnaLetra+"1", question.Question.EnumLabel+" "+question.Question.Title)
					_ = f.SetCellValue("Sheet1", columnaLetra+"2", qR.AditionalResponse)
					contCol++

					auxquestionFijas.IsOtter = true
					auxquestionFijas.ColumnaOtro = columnaLetra
				}

				questionFijas = append(questionFijas, auxquestionFijas)
			}
		}
	}

	// creando encavezados - LONG, LAT, ALT

	columnaLetra, _ := excelize.ColumnNumberToName(contCol)
	_ = f.SetCellValue("Sheet1", columnaLetra+"1", "Longitude")
	contCol++

	columnaLetra, _ = excelize.ColumnNumberToName(contCol)
	_ = f.SetCellValue("Sheet1", columnaLetra+"1", "Latitude")
	contCol++

	columnaLetra, _ = excelize.ColumnNumberToName(contCol)
	_ = f.SetCellValue("Sheet1", columnaLetra+"1", "Altitude")
	contCol++

	//columnaLetra, _ = excelize.ColumnNumberToName(contCol)
	////TODO: para la uspme putumayo
	////_ = f.SetCellValue("Sheet1", columnaLetra+"1", "Altitude")
	//_ = f.SetCellValue("Sheet1", columnaLetra+"1", "Precisión")
	//contCol++

	columnaLetra, _ = excelize.ColumnNumberToName(contCol)
	_ = f.SetCellValue("Sheet1", columnaLetra+"1", "Fecha")
	contCol++

	columnaLetra, _ = excelize.ColumnNumberToName(contCol)
	_ = f.SetCellValue("Sheet1", columnaLetra+"1", "ID")
	contCol++

	columnaLetra, _ = excelize.ColumnNumberToName(contCol)
	_ = f.SetCellValue("Sheet1", columnaLetra+"1", "Hora")
	contCol++

	//TODO: AÑADIENDO FILAS DE IMAGNES

	//columnaLetra, _ = excelize.ColumnNumberToName(contCol)
	//_ = f.SetCellValue("Sheet1", columnaLetra+"1", "P.24 img - De acuerdo con el recibo diligencie lo siguiente:")
	//contCol++
	//
	//columnaLetra, _ = excelize.ColumnNumberToName(contCol)
	//_ = f.SetCellValue("Sheet1", columnaLetra+"1", "P.48 img - De acuerdo con el recibo diligencie lo siguiente:")
	//contCol++
	//
	//columnaLetra, _ = excelize.ColumnNumberToName(contCol)
	//_ = f.SetCellValue("Sheet1", columnaLetra+"1", "P.60 img - Qué tipo de estufa de leña tiene?")

	constColImg := contCol
	for _, questionImage := range questionImagesHeders {
		columnaLetra, _ := excelize.ColumnNumberToName(constColImg)
		_ = f.SetCellValue("Sheet1", columnaLetra+"2", questionImage.EnumLabel+" "+questionImage.Title)
		constColImg++
	}

	// TODO: realizando prueba de preguntas autoguardadas

	//contColFija := 1
	//for _, questionFija := range questionFijas {
	//
	//	columnaLetra, _ := excelize.ColumnNumberToName(contColFija)
	//	_ = f.SetCellValue("Sheet1", columnaLetra+"3", strconv.FormatInt(questionFija.QID, 10) + " - " + strconv.FormatInt(questionFija.RID,10))
	//	contColFija++
	//
	//	if questionFija.IsOtter {
	//		columnaLetra, _ := excelize.ColumnNumberToName(contColFija)
	//		_ = f.SetCellValue("Sheet1", columnaLetra+"3", strconv.FormatInt(questionFija.QID, 10) + " - otro")
	//		contColFija++
	//	}
	//
	//}

	// realizando prueba de encuestas

	fmt.Println("************* INSERTANDO DATOS EN ARCHIVO *****************")

	var filaIncial int = 3
	dataSurveyRes := GetConstructorSurveryResponse(ID_ST, fechaInicio, fechaFin)

	var dataSurveyMaster []models.SurveyMaster
	if fechaInicio == "" || fechaFin == "" {
		dataSurveyMaster, _ = dal.SurveyMaterDAL{}.GetAllByIdSurveryTemplate(ID_ST)
	} else {
		dataSurveyMaster, _ = dal.SurveyMaterDAL{}.GetAllByIdSurveryTemplateRango(ID_ST, fechaInicio, fechaFin)
	}

	fmt.Println("************* INSERTANDO DATOS EN ARCHIVO *****************")
	fmt.Println("NUM. FILAS = " + strconv.Itoa(len(dataSurveyRes)))

	for _, dataSR := range dataSurveyRes {

		//fmt.Println(strconv.Itoa(row + 1) + "/" +  strconv.Itoa(len(dataSurveyRes)))

		for _, questionFija := range questionFijas {

			for _, dataQ := range dataSR.SurveyResponseQue { // questiions

				if dataQ.IDQuestion == questionFija.QID {

					for _, dataR := range dataQ.ResponseCon { // respuestas

						if questionFija.IdControlGroup == 2 {

							_ = f.SetCellValue("Sheet1", questionFija.Columna+strconv.Itoa(filaIncial), dataR.Value)

							if questionFija.IsOtter && dataR.Value != "0" {
								_ = f.SetCellValue("Sheet1", questionFija.ColumnaOtro+strconv.Itoa(filaIncial), dataR.AditionalResponse)
							}
							break

						} else {

							if dataR.IdResponse == questionFija.RID {
								_ = f.SetCellValue("Sheet1", questionFija.Columna+strconv.Itoa(filaIncial), dataR.Value)

								if questionFija.IsOtter && dataR.Value != "0" {
									_ = f.SetCellValue("Sheet1", questionFija.ColumnaOtro+strconv.Itoa(filaIncial), dataR.AditionalResponse)
								}
								break
							}
						}
					}
					break
				}
			}
		}
		// buscar cordenadas

		conf, err := dal.GetConfig()
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, dsm := range dataSurveyMaster {

			if dsm.ID == dataSR.IDSurvey {

				contColAux := contCol - 6 //iniciarla para poner desde la longitud, lo que llega se le baja a 6
				columnaLetra, _ := excelize.ColumnNumberToName(contColAux)
				_ = f.SetCellValue("Sheet1", columnaLetra+strconv.Itoa(filaIncial), dsm.Longitude)
				contColAux++

				columnaLetra, _ = excelize.ColumnNumberToName(contColAux)
				_ = f.SetCellValue("Sheet1", columnaLetra+strconv.Itoa(filaIncial), dsm.Latitude)
				contColAux++

				// obtener altitude

				//s := strings.Split(dsm.FechaCreate, "&")
				//var altitude string = "0"
				//if len(s) > 1 {
				//	altitude = s[1]
				//}
				//columnaLetra, _ = excelize.ColumnNumberToName(contCol - 2)
				//_ = f.SetCellValue("Sheet1", columnaLetra+strconv.Itoa(filaIncial), altitude)

				// precision (solo que se gurda en el campo de altitude)
				columnaLetra, _ = excelize.ColumnNumberToName(contColAux)
				_ = f.SetCellValue("Sheet1", columnaLetra+strconv.Itoa(filaIncial), dsm.Altitude)
				contColAux++

				columnaLetra, _ = excelize.ColumnNumberToName(contColAux)
				_ = f.SetCellValue("Sheet1", columnaLetra+strconv.Itoa(filaIncial), dsm.FechaFinish)
				contColAux++

				columnaLetra, _ = excelize.ColumnNumberToName(contColAux)
				_ = f.SetCellValue("Sheet1", columnaLetra+strconv.Itoa(filaIncial), dsm.ID)
				contColAux++

				columnaLetra, _ = excelize.ColumnNumberToName(contColAux)
				_ = f.SetCellValue("Sheet1", columnaLetra+strconv.Itoa(filaIncial), dsm.Hora)
				contColAux++

				// TODO: Imagenes

				for _, questionImageH := range questionImagesHeders {

					for _, questionImageD := range questionImagesData {

						if questionImageH.ID == questionImageD.IDQuestion && questionImageD.IdSurvey == dsm.ID && questionImageD.IdMobileArtefact == dsm.IdArtefact {
							columnaLetra, _ := excelize.ColumnNumberToName(contColAux)
							//_ = f.SetCellValue("Sheet1",  columnaLetra+strconv.Itoa(filaIncial), questionImageD.Patch)
							_ = f.SetCellValue("Sheet1", columnaLetra+strconv.Itoa(filaIncial), "link")
							_ = f.SetCellHyperLink("Sheet1", columnaLetra+strconv.Itoa(filaIncial), conf.PatchImg+questionImageD.Patch, "External")
						}

					}
					contColAux++
				}

				//columnaLetra, _ = excelize.ColumnNumberToName(contColAux)
				//_ = f.SetCellValue("Sheet1", columnaLetra+strconv.Itoa(filaIncial), "link")
				//_ = f.SetCellHyperLink("Sheet1", columnaLetra+strconv.Itoa(filaIncial), "https://sig.upme.gov.co/EncuestaPERS/UploadedFiles/SurveyPictures/Artefact_10054/SurveyTemplate_4/Survey_651_10054_2020_11_11__05_53_48/Question_376/picture_2020_11_12__11_34_00.jpg", "External")
				//contColAux++
				//
				//columnaLetra, _ = excelize.ColumnNumberToName(contColAux)
				//_ = f.SetCellValue("Sheet1", columnaLetra+strconv.Itoa(filaIncial), "link")
				//_ = f.SetCellHyperLink("Sheet1", columnaLetra+strconv.Itoa(filaIncial), "https://sig.upme.gov.co/EncuestaPERS/UploadedFiles/SurveyPictures/Artefact_10054/SurveyTemplate_4/Survey_651_10054_2020_11_11__05_53_48/Question_376/picture_2020_11_12__11_34_00.jpg", "External")
				//contColAux++
				//
				//columnaLetra, _ = excelize.ColumnNumberToName(contColAux)
				//_ = f.SetCellValue("Sheet1", columnaLetra+strconv.Itoa(filaIncial), "link")
				//_ = f.SetCellHyperLink("Sheet1", columnaLetra+strconv.Itoa(filaIncial), "https://sig.upme.gov.co/EncuestaPERS/UploadedFiles/SurveyPictures/Artefact_10054/SurveyTemplate_4/Survey_651_10054_2020_11_11__05_53_48/Question_376/picture_2020_11_12__11_34_00.jpg", "External")
				//contColAux++

				break
			}

		}
		filaIncial++
	}

	fmt.Println("************* FIN CREACION DE ARCHIVO *****************")

	dirGuardar := "../../PUBLIC/docs_sin_validar/" + clave + ".xlsx"

	err := f.SaveAs(dirGuardar)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Add("Content-Disposition", "Attachment")
	w.Header().Add("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	http.ServeFile(w, r, dirGuardar)
}
