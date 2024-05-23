package models

type QuestionImage struct {
	ID int64 `json:"id"`
	IDQuestion int64 `json:"id_question"`
	IdSurvey int64 `json:"id_survey"`
	Patch string `json:"patch"`
	DateCreate string `json:"date_create"`
	IdMobileArtefact int `json:"id_mobile_artefact"`
}