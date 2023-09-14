package helper

import (
	"Notion-Forms/internal/model"
	"fmt"
)

func ConvertFormsListToFormResponse(forms []model.Form, baseUrl string) []model.FormResponse {
	var formsResponse []model.FormResponse
	for _, form := range forms {
		formsResponse = append(formsResponse, ConvertFormToFormResponse(form, baseUrl))
	}
	return formsResponse
}

func ConvertFormToFormResponse(form model.Form, baseUrl string) model.FormResponse {
	return model.FormResponse{
		Url:        fmt.Sprintf("%s/form/%s", baseUrl, form.UrlId),
		IamUserId:  form.IamUserId,
		DatabaseId: form.DatabaseId,
		Password:   form.Password,
		Overrides:  form.Overrides,
		Storage:    form.Storage,
	}
}
