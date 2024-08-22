package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/kcharymyrat/e-commerce/internal/data"
)

func (app *application) createCategoryHandler(w http.ResponseWriter, r *http.Request) {
	var categoryInput data.Category
	err := app.readJSON(w, r, &categoryInput)
	if err != nil {
		app.badRequestResponse(w, r, err)
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
