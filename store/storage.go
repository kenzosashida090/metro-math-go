package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrInternalServer    = errors.New("something went wrong")
	ErrNotFound          = errors.New("Record not found")
	QueryTimeoutDuration = time.Second * 5
	ErrDuplicatedName    = errors.New("Duplicated name")
)

func ErrorFactoryDB(err error) error {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return ErrNotFound
	case err.Error() == `pq:duplicated key value violates unique constraint "songs_name"`:
		return ErrDuplicatedName
	default:
		return ErrInternalServer
	}
}

type Storage struct {
	Songs interface {
		CreateSong(context.Context, *Song) error
		UpdateSong(context.Context, int64, *Song) error
		DeleteSong(context.Context, int64) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Songs: &SongStore{db},
	}
}
