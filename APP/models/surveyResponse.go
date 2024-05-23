package models

//SurveryResponse modelo
type SurveryResponse struct {
	ID         int64
	IdSurvey   int64
	IdQuestion int64
	IdResponse int64
	Value      string
	AditionalResponse string
}
