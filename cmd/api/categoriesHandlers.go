package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/kcharymyrat/e-commerce/internal/data"
	"github.com/kcharymyrat/e-commerce/internal/validator"
)

func (app *application) createCategoryHandler(w http.ResponseWriter, r *http.Request) {
	var categoryInput struct {
		NameTk   string           `json:"name_tk"`
		NameEn   string           `json:"name_en"`
		NameRu   string           `json:"name_ru"`
		Parent   *data.Category   `json:"parent,omitempty"`
		Children []*data.Category `json:"children,omitempty"`
	}

	err := app.readJSON(w, r, &categoryInput)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	category := &data.Category{
		NameTk:   categoryInput.NameTk,
		NameEn:   categoryInput.NameEn,
		NameRu:   categoryInput.NameRu,
		Parent:   categoryInput.Parent,
		Children: categoryInput.Children,
	}

	v := validator.New()

	if data.ValidateCategory(v, category); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	fmt.Fprintf(w, "%v+\n", categoryInput)
}

func (app *application) getCategoryHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readUUIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	category := data.Category{
		ID:        id,
		NameTk:    "Telewizor",
		NameEn:    "TV",
		NameRu:    "Телевизор",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Parent:    nil,
		Children:  nil,
	}
	envelope := envelope{
		"category": category,
	}

	fmt.Println("envelope =", envelope)

	err = app.writeJson(w, http.StatusOK, envelope, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
