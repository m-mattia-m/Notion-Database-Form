package service

import (
	"Notion-Forms/internal/model"
	"github.com/google/uuid"
)

func (svc Clients) CreateFormToDatabase(oidcUser model.OidcUser, formRequestBody model.FormRequest) (*model.Form, error) {
	var form model.Form
	form.IamUserId = oidcUser.Sub
	form.UrlId = uuid.New().String()
	form.DatabaseId = formRequestBody.DatabaseId
	form.Password = formRequestBody.Password
	form.Overrides = formRequestBody.Overrides

	err := svc.db.dao.CreateForm(form)
	if err != nil {
		return nil, err
	}

	return &form, err
}

func (svc Clients) ListForms(oidcUser model.OidcUser) ([]model.Form, error) {
	forms, err := svc.db.dao.ListForms(oidcUser.Sub)
	if err != nil {
		return nil, err
	}

	return forms, err
}

func (svc Clients) GetFormById(databaseId string, oidcUser model.OidcUser) (*model.Form, error) {
	form, err := svc.db.dao.GetFormById(databaseId, oidcUser.Sub)
	if err != nil {
		return nil, err
	}

	return form, err
}

func (svc Clients) UpdateFormById(formRequestBody model.FormRequest, oidcUser model.OidcUser) (*model.Form, error) {
	form, err := svc.GetFormById(formRequestBody.DatabaseId, oidcUser)
	if err != nil {
		return nil, err
	}

	form.Password = formRequestBody.Password
	form.Overrides = formRequestBody.Overrides

	updatedForm, err := svc.db.dao.UpdateFormById(*form)
	if err != nil {
		return nil, err
	}

	return updatedForm, err
}

func (svc Clients) DeleteFormById(databaseId string, oidcUser model.OidcUser) error {
	err := svc.db.dao.DeleteFormById(databaseId, oidcUser.Sub)
	if err != nil {
		return err
	}

	return err
}
