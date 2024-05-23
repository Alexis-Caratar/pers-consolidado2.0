package models

//Survey modelo
type SurveyMaster struct {
	ID        int64
	IdArtefact int
	IdSurvey int
	Latitude  float64
	Longitude float64
	Altitude  float64
	FechaFinish string
	FechaCreate string
	Hora string
}
