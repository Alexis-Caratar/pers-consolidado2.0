package models

//ResponseQuestion modelo
type ResponseQuestion struct {
	QID               int64
	QTitle            string
	OrderQuestion     int64
	QEnumLabel        string
	IdControlGroup    int64
	RID               int64
	RTitle            string
	OrderResponse     int64
	IndexRow          int64
	IndexColumn       int64
	Value             string
	IdResponseType    int64
	REnumLabel        string
	AditionalResponse string
	IsImages          string
}
