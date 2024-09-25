package data

import (
	"errors"
	"time"
	"unicode"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kcharymyrat/e-commerce/internal/validator"
	"github.com/shopspring/decimal"
)

type User struct {
	ID                     uuid.UUID       `json:"id,omitempty" db:"id"`
	Phone                  string          `json:"phone" db:"phone"`
	PasswordHash           string          `json:"-" db:"password_hash"`                 // No need to serialize password hash
	FirstName              *string         `json:"first_name,omitempty" db:"first_name"` // Can be NULL
	LastName               *string         `json:"last_name,omitempty" db:"last_name"`   // Can be NULL
	Patronomic             *string         `json:"patronomic,omitempty" db:"patronomic"` // Can be NULL
	DOB                    *time.Time      `json:"dob,omitempty" db:"dob"`               // Can be NULL
	Email                  *string         `json:"email,omitempty" db:"email"`           // Can be NULL
	IsActive               bool            `json:"is_active" db:"is_active"`
	IsBanned               bool            `json:"is_banned" db:"is_banned"`
	IsTrusted              bool            `json:"is_trusted" db:"is_trusted"`
	InvitedById            *uuid.UUID      `json:"invited_by_id,omitempty" db:"invited_by_id"`     // Can be NULL
	InvRefId               *uuid.UUID      `json:"inv_ref_id,omitempty" db:"inv_ref_id"`           // Can be NULL
	InvProdRefId           *uuid.UUID      `json:"inv_prod_ref_id,omitempty" db:"inv_prod_ref_id"` // Can be NULL
	RefSignups             int             `json:"ref_signups" db:"ref_signups"`
	ProdRefSignups         int             `json:"prod_ref_signups" db:"prod_ref_signups"`
	ProdRefBought          int             `json:"prod_ref_bought" db:"prod_ref_bought"`
	TotalRefferals         int             `json:"total_referrals" db:"total_referrals"`
	DynamicDiscountPercent decimal.Decimal `json:"dynamic_discount_percent" db:"_dynamic_discount_percent"`
	DynDiscPercent         decimal.Decimal `json:"dyn_disc_percent" db:"dyn_disc_percent"`
	BonusPoints            decimal.Decimal `json:"bonus_points" db:"bonus_points"`
	IsStaff                bool            `json:"is_staff" db:"is_staff"`
	IsAdmin                bool            `json:"is_admin" db:"is_admin"`
	IsSuperuser            bool            `json:"is_superuser" db:"is_superuser"`
	CreatedAt              time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time       `json:"updated_at" db:"updated_at"`
	CreatedById            *uuid.UUID      `json:"created_by_id,omitempty" db:"created_by_id"` // Can be NULL
	UpdatedById            *uuid.UUID      `json:"updated_by_id,omitempty" db:"updated_by_id"` // Can be NULL
}

var (
	ErrPhoneNumber = errors.New("can not use this phone number")
)

func IsValidPassword(passwordPlaintext string) (bool, []string) {
	var hasUpper, hasLower, hasNumber, isAscii bool

	isValid := true
	reasons := []string{}

	if len([]byte(passwordPlaintext)) <= 8 {
		isValid = false
		reasons = append(reasons, "must be at least 8 bytes long")
	}

	if len([]byte(passwordPlaintext)) >= 72 {
		isValid = false
		reasons = append(reasons, "password", "must not be more than 500 bytes long")
	}

	for _, char := range passwordPlaintext {
		if char > unicode.MaxASCII {
			isAscii = false
		}

		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		}
	}

	if !isAscii {
		isValid = false
		reasons = append(reasons, "must not contain non ASCII character")
	}

	if !hasUpper {
		isValid = false
		reasons = append(reasons, "must contain upper case character")
	}

	if !hasLower {
		isValid = false
		reasons = append(reasons, "must contain lower case character")
	}

	if !hasNumber {
		isValid = false
		reasons = append(reasons, "must contain number")
	}

	return isValid, reasons
}

func ValidatePhone(v *validator.Validator, phone string) {
	v.Check(phone != "", "phone", "must be provided")
	v.Check(validator.Matches(phone, validator.PhoneNumberRX), "phone", "must be a valid phone number")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	isValid, reasons := IsValidPassword(password)
	if !isValid {
		for _, reason := range reasons {
			v.Check(false, "password", reason)
		}
	}
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidateUser(v *validator.Validator, user *User, passwordPlaintext string) {
	ValidatePhone(v, user.Phone)
	ValidatePasswordPlaintext(v, passwordPlaintext)
	if user.Email != nil {
		ValidateEmail(v, *user.Email)
	}
}

type UserModel struct {
	DBPOOL *pgxpool.Pool
}

// func (u UserModel) Insert(user *User) error {
// 	query := `INSERT INTO users (phone, password_hash, first_name, last_name, patronymic, dob, email)
// 		VALUES ($1, $2, $3, $4, $5, %6, $7)
// 		RETURNING id, phone, first_name, last_name;
// 	`
// }
