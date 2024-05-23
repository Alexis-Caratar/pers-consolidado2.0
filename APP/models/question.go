package models

//Question modelo
type Question struct {
	ID               int64  `json:"id,string"`
	Title            string `json:"title"`
	IdSurveyTemplate int    `json:"id_survey_template,string"`
	OrderQuestion    int64  `json:"order_question,string"`
	EnumLabel        string `json:"enum_label,string"`
	IdControlGroup   int64  `json:"id_control_group,string"`
	IsImages         string `json:"is_images"`
}
