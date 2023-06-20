package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"projects/Repairment_service/Repairment_order_service/genproto/order_service"
	"projects/Repairment_service/Repairment_order_service/models"
	"projects/Repairment_service/Repairment_order_service/pkg/helper"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type orderRepo struct {
	db *pgxpool.Pool
}

func NewOrderRepo(db *pgxpool.Pool) *orderRepo {
	return &orderRepo{
		db: db,
	}
}

func (u *orderRepo) Create(ctx context.Context, req *order_service.CreateOrderRequest) (resp *order_service.Order, err error) {
	id := uuid.New().String()

	query := `
		INSERT INTO "order" (
			id,
			user_id,
			updated_at
		)
		VALUES ($1, $2, NOW())
	`

	_, err = u.db.Exec(ctx, query, id, req.UserId)
	if err != nil {
		return nil, err
	}

	return &order_service.Order{
		Id:     id,
		UserId: req.UserId,
	}, nil
}

func (u *orderRepo) GetById(ctx context.Context, req *order_service.OrderPrimaryKey) (resp *order_service.Order, err error) {
	query := `
		Select 
			id,
			user_id,
			is_completed,
			created_at,
			updated_at
		from "order"
		where id = $1
	`

	var (
		id           sql.NullString
		user_id      sql.NullString
		is_completed sql.NullBool
		created_at   sql.NullString
		updated_at   sql.NullString
	)

	err = u.db.QueryRow(ctx, query, req.Id).Scan(
		&id,
		&user_id,
		&is_completed,
		&created_at,
		&updated_at,
	)

	if err != nil {
		return resp, err
	}

	resp = &order_service.Order{
		Id:          id.String,
		UserId:      user_id.String,
		IsCompleted: is_completed.Bool,
		CreatedAt:   created_at.String,
		UpdatedAt:   updated_at.String,
	}
	return
}

func (u *orderRepo) GetList(ctx context.Context, req *order_service.GetListOrderRequest) (resp *order_service.GetListOrderResponse, err error) {
	resp = &order_service.GetListOrderResponse{}

	var (
		query  string
		limit  = ""
		offset = " OFFSET 0"
		params = make(map[string]interface{})
		filter = " WHERE TRUE"
		sort   = " ORDER BY created_at DESC"
	)

	query = `
		SELECT
			COUNT(*) OVER(),
			id,
			user_id,
			is_completed,
			to_char(created_at, 'YYYY-MM-DD HH24:MI:SS'),
			to_char(updated_at, 'YYYY-MM-DD HH24:MI:SS')
		FROM "order"
		`

	if len(req.GetSearch()) > 0 {
		filter += " AND (user_id || ' ' is_completed) ILIKE '%' || '" + req.Search + " || '%' "
	}
	if req.GetLimit() > 0 {
		limit = " LIMIT :limit"
		params["limit"] = req.Limit
	}
	if req.GetOffset() > 0 {
		offset = " OFFSET :offset"
		params["offset"] = req.Offset
	}

	query += filter + sort + offset + limit

	query, args := helper.ReplaceQueryParams(query, params)
	rows, err := u.db.Query(ctx, query, args...)
	if err != nil {
		return resp, err
	}

	defer rows.Close()

	for rows.Next() {
		var (
			id           sql.NullString
			user_id      sql.NullString
			is_completed sql.NullBool
			created_at   sql.NullString
			updated_at   sql.NullString
		)

		err := rows.Scan(
			&resp.Total,
			&id,
			&user_id,
			&is_completed,
			&created_at,
			&updated_at,
		)
		if err != nil {
			return resp, err
		}

		resp.Orders = append(resp.Orders, &order_service.Order{
			Id:          id.String,
			UserId:      user_id.String,
			IsCompleted: is_completed.Bool,
			CreatedAt:   created_at.String,
			UpdatedAt:   updated_at.String,
		})
	}

	return
}

func (u *orderRepo) Update(ctx context.Context, req *order_service.UpdateOrderRequest) (rowsAffected int64, err error) {
	var (
		query  string
		params = make(map[string]interface{})
	)

	query = `
		UPDATE "order"
		SET
			is_completed = :is_completed,
			updated_at = NOW()
		WHERE id = :id
	`
	params = map[string]interface{}{
		"is_completed": req.GetIsCompleted(),
	}

	query, args := helper.ReplaceQueryParams(query, params)

	result, err := u.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}

func (u *orderRepo) UpdatePatch(ctx context.Context, req *models.UpdatePatchRequest) (rowsAffected int64, err error) {
	var (
		set   = " SET "
		ind   = 0
		query string
	)

	if len(req.Fields) == 0 {
		err = errors.New("no updates provided")
		return
	}

	req.Fields["id"] = req.Id

	for key := range req.Fields {
		set += fmt.Sprintf(" %s = :%s ", key, key)
		if ind != len(req.Fields)-1 {
			set += ", "
		}
		ind++
	}

	query = `
		UPDATE "order"
		    ` + set + ` , updated_at = NOW()
		WHERE id = :id
	`

	query, args := helper.ReplaceQueryParams(query, req.Fields)

	result, err := u.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}

func (u *orderRepo) Delete(ctx context.Context, req *order_service.OrderPrimaryKey) (*empty.Empty, error) {
	query := `
		DELETE FROM "order"
		WHERE id = $1
	`

	_, err := u.db.Exec(ctx, query, req.Id)
	if err != nil {
		return &empty.Empty{}, err
	}

	return &empty.Empty{}, nil
}
