package models

type ConsolidadoSurveyResponse struct {
	IDSurvey   int64
	SurveyResponseQue []ConsolidadoQuesSurveyResponse
}
