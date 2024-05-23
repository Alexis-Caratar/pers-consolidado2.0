package models

//Question modelo
type Response struct {
	ID                int64
	Title             string
	OrdenResponse     int64
	IndexRow          int64
	IndexColumn       int64
	Value             string
	IdResponseType    int64
	EnumLabel         string
	AditionalResponse string
	IdQuestion        int64
}
