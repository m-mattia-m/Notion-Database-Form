package model

// Map the form-request-data manually to the form because mongodb will otherwise create a subobject and it can no longer be mapped by database ID.

type Form struct {
	UrlId      string     `json:"url_id"`
	IamUserId  string     `json:"iam_user_id"`
	DatabaseId string     `json:"database_id"`
	Password   string     `json:"password"`
	Overrides  []Override `json:"overrides"`
	Storage
}

type Override struct {
	ColumnId        string `json:"column_id"`
	NewPropertyType string `json:"new_property_type"`
}

type FormRequest struct {
	DatabaseId string     `json:"database_id"`
	Password   string     `json:"password"`
	Overrides  []Override `json:"overrides"`
}

type FormResponse struct {
	Url        string     `json:"url"`
	IamUserId  string     `json:"iam_user_id"`
	DatabaseId string     `json:"database_id"`
	Password   string     `json:"password"`
	Overrides  []Override `json:"overrides"`
	Storage
}
