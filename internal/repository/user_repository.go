package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kcharymyrat/e-commerce/api/requests"
	"github.com/kcharymyrat/e-commerce/internal/auth"
	"github.com/kcharymyrat/e-commerce/internal/common"
	"github.com/kcharymyrat/e-commerce/internal/data"
	"github.com/kcharymyrat/e-commerce/internal/filters"
	"github.com/kcharymyrat/e-commerce/internal/types"
)

type UserRepository struct {
	DBPOOL *pgxpool.Pool
}

func (r UserRepository) Create(user *data.User) error {
	passwordHashBytes, err := auth.GeneratePasswordHash(user.Password)
	if err != nil {
		return err
	}

	args := []interface{}{
		user.Phone,
		passwordHashBytes,
		user.FirstName,
		user.LastName,
		user.Patronomic,
		user.Email,
		user.IsActive,
	}

	query := `
		INSERT INTO users (
			phone, 
			password_hash, 
			first_name,
			last_name,
			patronmic,
			email,
			is_active 
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, phone, first_name, last_name, patronomic, email, is_active, created_by_id
	`

	if user.CreatedByID != nil {
		query = `
		INSERT INTO users (
			phone, 
			password_hash, 
			first_name,
			last_name,
			patronmic,
			email,
			is_active, 
			created_by_id
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, phone, first_name, last_name, patronomic, email, is_active, created_by_id
		`
		args = append(args, *user.CreatedByID)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = r.DBPOOL.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.Phone,
		&user.FirstName,
		&user.LastName,
		&user.Patronomic,
		&user.Email,
		&user.IsActive,
		&user.CreatedByID,
	)

	return err
}

func (r UserRepository) GetByID(id uuid.UUID) (*data.User, error) {
	query := `
	SELECT * 
	FROM users
	WHERE id = $1
	`
	var user data.User

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.DBPOOL.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Phone,
		&user.FirstName,
		&user.LastName,
		&user.Patronomic,
		&user.DOB,
		&user.Email,
		&user.IsActive,
		&user.IsBanned,
		&user.IsTrusted,
		&user.InvitedByID,
		&user.InvRefID,
		&user.InvProdRefID,
		&user.RefSignups,
		&user.ProdRefSignups,
		&user.ProdRefBought,
		&user.TotalRefferals,
		&user.WholeDynDiscPercent,
		&user.DynDiscPercent,
		&user.BonusPoints,
		&user.IsStaff,
		&user.IsAdmin,
		&user.IsSuperuser,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.CreatedByID,
		&user.UpdatedByID,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, common.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (r UserRepository) GetByPhone(phone string) (*data.User, error) {
	query := `
	SELECT * 
	FROM users
	WHERE phone = $1
	`
	var user data.User

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.DBPOOL.QueryRow(ctx, query, phone).Scan(
		&user.ID,
		&user.Phone,
		&user.FirstName,
		&user.LastName,
		&user.Patronomic,
		&user.PasswordHash,
		&user.DOB,
		&user.Email,
		&user.IsActive,
		&user.IsBanned,
		&user.IsTrusted,
		&user.InvitedByID,
		&user.InvRefID,
		&user.InvProdRefID,
		&user.RefSignups,
		&user.ProdRefSignups,
		&user.ProdRefBought,
		&user.TotalRefferals,
		&user.WholeDynDiscPercent,
		&user.DynDiscPercent,
		&user.BonusPoints,
		&user.IsStaff,
		&user.IsAdmin,
		&user.IsSuperuser,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.CreatedByID,
		&user.UpdatedByID,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, common.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (r UserRepository) List(f *requests.UsersAdminFilters) ([]*data.User, types.PaginationMetadata, error) {
	query := `SELECT *
	FROM users
	WHERE 1=1`

	args := []interface{}{}
	argCounter := 1

	addUserSpecificFiltersToSQL(f, &query, &argCounter, args)

	if f.Search != nil {
		query += fmt.Sprintf(
			` AND (
				id = $%d OR 
				to_tsvector('simple', phone) @@ plainto_tsquery('phone', $%d) OR 
				to_tsvector('simple', email) @@ plainto_tsquery('email', $%d) OR 
				to_tsvector('simple', first_name) @@ plainto_tsquery('first_name', $%d) OR 
				to_tsvector('simple', last_name) @@ plainto_tsquery('last_name', $%d) OR 
				to_tsvector('simple', first_name || ' ' || last_name) @@ plainto_tsquery('simple', $%d) OR
			)`, argCounter, argCounter, argCounter, argCounter, argCounter, argCounter,
		)
		args = append(args, *f.Search)
		argCounter++
	}

	filters.AddCreatedUpdateAtFilterToSQL(&f.CreatedUpdatedAtFilter, &query, &argCounter, args)
	filters.AddCreatedUpdateByFilterToSQL(&f.CreatedUpdatedByFilter, &query, &argCounter, args)
	filters.AddSortListFilterToSQL(&f.SortListFilter, &query)
	filters.AddPaginationFilterToSQL(&f.PaginationFilter, &query, &argCounter, args)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.DBPOOL.Query(ctx, query, args...)
	if err != nil {
		return nil, types.PaginationMetadata{}, err
	}

	users := []*data.User{}

	for rows.Next() {
		var user data.User
		err := rows.Scan(
			&user.ID,
			&user.Phone,
			&user.FirstName,
			&user.LastName,
			&user.Patronomic,
			&user.DOB,
			&user.Email,
			&user.IsActive,
			&user.IsBanned,
			&user.IsTrusted,
			&user.InvitedByID,
			&user.InvRefID,
			&user.InvProdRefID,
			&user.RefSignups,
			&user.ProdRefSignups,
			&user.ProdRefBought,
			&user.TotalRefferals,
			&user.WholeDynDiscPercent,
			&user.DynDiscPercent,
			&user.BonusPoints,
			&user.IsStaff,
			&user.IsAdmin,
			&user.IsSuperuser,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.CreatedByID,
			&user.UpdatedByID,
			&user.Version,
		)
		if err != nil {
			return nil, types.PaginationMetadata{}, err
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, types.PaginationMetadata{}, err
	}

	metadata := common.CalculateMetadata(len(users), *f.Page, *f.PageSize)

	return users, metadata, nil
}

func (r UserRepository) Update(user *data.User) error {
	query := `UPDATE users
		SET 
			phone = $1,
			first_name = $2,
			last_name = $3,
			patronomic = $4,
			email = $5
			is_active = $6,
			updated_by_id = $7,
			version = version + 1
		WHERE id = $8 AND version = $9
		`

	args := []interface{}{
		user.Phone,
		user.FirstName,
		user.LastName,
		user.Patronomic,
		user.Email,
		user.IsActive,
		user.UpdatedByID,
		user.ID,
		user.Version,
	}

	if user.Password != "" {
		user.Password = strings.TrimSpace(user.Password)
		passwordHashBytes, err := auth.GeneratePasswordHash(user.Password)
		if err != nil {
			return err
		}

		query = `UPDATE users
		SET 
			phone = $1,
			password = $2,
			first_name = $3,
			last_name = $4,
			patronomic = $5,
			email = $6
			is_active = $7,
			updated_by_id = $8,
			version = version + 1,
		WHERE id = $9 AND version = $10
		`

		args = []interface{}{
			user.Phone,
			passwordHashBytes,
			user.FirstName,
			user.LastName,
			user.Patronomic,
			user.Email,
			user.IsActive,
			user.UpdatedByID,
			user.ID,
			user.Version,
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.DBPOOL.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.Phone,
		&user.FirstName,
		&user.LastName,
		&user.Patronomic,
		&user.DOB,
		&user.Email,
		&user.IsActive,
		&user.IsBanned,
		&user.IsTrusted,
		&user.InvitedByID,
		&user.InvRefID,
		&user.InvProdRefID,
		&user.RefSignups,
		&user.ProdRefSignups,
		&user.ProdRefBought,
		&user.TotalRefferals,
		&user.WholeDynDiscPercent,
		&user.DynDiscPercent,
		&user.BonusPoints,
		&user.IsStaff,
		&user.IsAdmin,
		&user.IsSuperuser,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.CreatedByID,
		&user.UpdatedByID,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return common.ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (r UserRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := r.DBPOOL.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() < 1 {
		return common.ErrRecordNotFound
	}

	return nil
}

func addUserSpecificFiltersToSQL(
	f *requests.UsersAdminFilters, query *string, argCounter *int, args []interface{},
) {
	if f.ID != nil {
		*query += fmt.Sprintf(" AND id = $%d", *argCounter)
		args = append(args, *f.ID)
		*argCounter++
	}

	if f.Phone != nil {
		*query += fmt.Sprintf(" AND phone = $%d", *argCounter)
		args = append(args, *f.Phone)
		*argCounter++
	}

	if f.Email != nil {
		*query += fmt.Sprintf(" AND email = $%d", *argCounter)
		args = append(args, *f.Email)
		*argCounter++
	}

	if f.IsActice != nil {
		*query += fmt.Sprintf(" AND is_active = $%d", *argCounter)
		args = append(args, *f.IsActice)
		*argCounter++
	}

	if f.IsBanned != nil {
		*query += fmt.Sprintf(" AND is_banned = $%d", *argCounter)
		args = append(args, *f.IsBanned)
		*argCounter++
	}

	if f.IsBanned != nil {
		*query += fmt.Sprintf(" AND is_banned = $%d", *argCounter)
		args = append(args, *f.IsBanned)
		*argCounter++
	}

	if f.IsTrusted != nil {
		*query += fmt.Sprintf(" AND is_trusted = $%d", *argCounter)
		args = append(args, *f.IsTrusted)
		*argCounter++
	}

	if f.IsTrusted != nil {
		*query += fmt.Sprintf(" AND is_trusted = $%d", *argCounter)
		args = append(args, *f.IsTrusted)
		*argCounter++
	}

	if f.IsInvited != nil {
		if *f.IsInvited {
			*query += " AND (invited_by_id IS NOT NULL OR inv_ref_id is NOT NULL OR inv_prod_ref_id IS NOT NULL)"
		} else {
			*query += " AND invited_by_id IS NULL AND inv_ref_id IS NULL AND inv_prod_ref_id IS NULL"
		}
	}

	if f.RefSignupsFrom != nil {
		*query += fmt.Sprintf(" AND ref_signups >= $%d", *argCounter)
		args = append(args, *f.RefSignupsFrom)
		*argCounter++
	}

	if f.RefSignupsTo != nil {
		*query += fmt.Sprintf(" AND ref_signups <= $%d", *argCounter)
		args = append(args, *f.RefSignupsTo)
		*argCounter++
	}

	if f.ProdRefSignupsFrom != nil {
		*query += fmt.Sprintf(" AND prod_ref_signups >= $%d", *argCounter)
		args = append(args, *f.ProdRefSignupsFrom)
		*argCounter++
	}

	if f.ProdRefSignupsTo != nil {
		*query += fmt.Sprintf(" AND prod_ref_signups <= $%d", *argCounter)
		args = append(args, *f.ProdRefSignupsTo)
		*argCounter++
	}

	if f.ProdRefBoughtFrom != nil {
		*query += fmt.Sprintf(" AND prod_ref_bought >= $%d", *argCounter)
		args = append(args, *f.ProdRefBoughtFrom)
		*argCounter++
	}

	if f.ProdRefBoughtTo != nil {
		*query += fmt.Sprintf(" AND prod_ref_bought <= $%d", *argCounter)
		args = append(args, *f.ProdRefBoughtTo)
		*argCounter++
	}

	if f.WholeDynDiscPercentFrom != nil {
		*query += fmt.Sprintf(" AND _dynamic_discount_percent >= $%d", *argCounter)
		args = append(args, *f.WholeDynDiscPercentFrom)
		*argCounter++
	}

	if f.WholeDynDiscPercentTo != nil {
		*query += fmt.Sprintf(" AND _dynamic_discount_percent <= $%d", *argCounter)
		args = append(args, *f.WholeDynDiscPercentTo)
		*argCounter++
	}

	if f.DynDiscPercentFrom != nil {
		*query += fmt.Sprintf(" AND dyn_disc_percent >= $%d", *argCounter)
		args = append(args, *f.DynDiscPercentFrom)
		*argCounter++
	}

	if f.DynDiscPercentTo != nil {
		*query += fmt.Sprintf(" AND dyn_disc_percent <= $%d", *argCounter)
		args = append(args, *f.DynDiscPercentTo)
		*argCounter++
	}

	if f.BonusPointsFrom != nil {
		*query += fmt.Sprintf(" AND bonus_points >= $%d", *argCounter)
		args = append(args, *f.BonusPointsFrom)
		*argCounter++
	}

	if f.BonusPointsTo != nil {
		*query += fmt.Sprintf(" AND bonus_points <= $%d", *argCounter)
		args = append(args, *f.BonusPointsTo)
		*argCounter++
	}

	if f.IsStaff != nil {
		*query += fmt.Sprintf(" AND is_staff = $%d", *argCounter)
		args = append(args, *f.IsStaff)
		*argCounter++
	}

	if f.IsAdmin != nil {
		*query += fmt.Sprintf(" AND is_admin = $%d", *argCounter)
		args = append(args, *f.IsAdmin)
		*argCounter++
	}

	if f.IsSuperuser != nil {
		*query += fmt.Sprintf(" AND is_superuser = $%d", *argCounter)
		args = append(args, *f.IsSuperuser)
		*argCounter++
	}

	fmt.Println(args)
}
