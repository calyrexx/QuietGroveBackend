package repository

import "context"

type IReservations interface {
	BookAHouse(ctx context.Context, req string) error
}
