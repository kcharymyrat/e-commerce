package main

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/kcharymyrat/e-commerce/internal/data"
	"github.com/kcharymyrat/e-commerce/internal/validator"
)

func (app *application) createCategoryHandler(w http.ResponseWriter, r *http.Request) {
	var categoryInput struct {
		Name        string    `json:"name"`
		Parent      uuid.UUID `json:"parent,omitempty"`
		Slug        string    `json:"slug"`
		Description string    `json:"description,omitempty"`
		ImageUrl    string    `json:"image_url"`
		CreatedByID uuid.UUID `json:"created_by_id"`
		UpdatedByID uuid.UUID `json:"updated_by_id"`
	}

	err := app.readJSON(w, r, &categoryInput)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	category := &data.Category{
		Name:        categoryInput.Name,
		Parent:      categoryInput.Parent,
		Slug:        categoryInput.Slug,
		Description: categoryInput.Description,
		ImageUrl:    categoryInput.ImageUrl,
		CreatedByID: categoryInput.CreatedByID,
		UpdatedByID: categoryInput.UpdatedByID,
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

	// category := data.Category{
	// 	ID:        id,
	// 	NameTk:    "Telewizor",
	// 	NameEn:    "TV",
	// 	NameRu:    "Телевизор",
	// 	CreatedAt: time.Now(),
	// 	UpdatedAt: time.Now(),
	// 	Parent:    nil,
	// 	Children:  nil,
	// }

	type category struct {
		ID   uuid.UUID
		Name string
	}

	cat := category{
		ID:   id,
		Name: "Example category",
	}

	envelope := envelope{
		"category": cat,
	}

	fmt.Println("envelope =", envelope)

	err = app.writeJson(w, http.StatusOK, envelope, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
