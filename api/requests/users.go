package requests

import (
	"github.com/google/uuid"
)

// TODO: UpdateUserInput - For Superusers and Admin - IsAdmin, IsStaff, IsSuperuser fields should be required

type CreateUserInput struct {
	Phone       string     `json:"phone" validate:"required,e164"`
	Password    string     `json:"password" validate:"required,min=8,max=72,password"`
	FirstName   string     `json:"first_name" validate:"omitempty,max=50,alpha"`
	LastName    string     `json:"last_name" validate:"omitempty,max=50,alpha"`
	Patronomic  string     `json:"patronomic" validate:"omitempty,max=50,alpha"`
	Email       string     `json:"email" validate:"omitempty,email"`
	IsActive    bool       `json:"is_active" validate:"required"`
	CreatedByID *uuid.UUID `json:"created_by_id,omitempty" validate:"omitempty,uuid"`
}
