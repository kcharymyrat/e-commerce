package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/kcharymyrat/e-commerce/internal/data"
)

func (app *application) createCategoryHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new category")
}

func (app *application) getCategoryHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readUUIDParam(r)
	if err != nil {
		http.NotFound(w, r)
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
		app.logger.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
