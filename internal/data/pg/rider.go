package pg

import (
	"context"

	"github.com/dead-letter/dead-letter-data/internal/data"
	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RiderService struct {
	pool        *pgxpool.Pool
	userService UserService
}

func NewRiderService(pool *pgxpool.Pool) RiderService {
	return RiderService{
		pool:        pool,
		userService: NewUserService(pool),
	}
}

func (s RiderService) Create(userID uuid.UUID) (*data.Rider, error) {
	var r data.Rider

	sql := `
		INSERT INTO rider_ (id_)
		VALUES($1);`

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	_, err := s.pool.Exec(ctx, sql, userID)
	if err != nil {
		return nil, err
	}

	r.User, err = s.userService.Read(userID)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (s RiderService) Read(userID uuid.UUID) (*data.Rider, error) {
	var r data.Rider
	var err error

	r.User, err = s.userService.Read(userID)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (s RiderService) Update(Rider *data.Rider) error {
	return nil
}

func (s RiderService) Delete(userID uuid.UUID) error {
	sql := `
		DELETE FROM rider_
		WHERE id_ = $1;`

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	res, err := s.pool.Exec(ctx, sql, userID)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return data.ErrRecordNotFound
	}

	return nil
}
