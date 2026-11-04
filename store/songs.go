package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type Section struct {
	ID     string `json:"id"`
	Groups []int  `json:"groups"`
	Bars   int    `json:"bars"`
}
type SongStore struct {
	db *sql.DB
}
type Song struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	BPM      int       `json:"bpm"`
	Sections []Section `json:"sections"`
}

func (s *SongStore) CreateSong(ctx context.Context, song *Song) error {
	query := `
	 			INSERT INTO songs (name, bpm, sections) VALUES ($1,$2,$3)	
				RETURNING id, name, bpm, sections
	 `
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		song.Name,
		song.BPM,
		song.Sections,
	).Scan(
		&song.ID,
		&song.Name,
		&song.BPM,
		&song.Sections,
	)
	if err != nil {
		return ErrorFactoryDB(err)
	}
	return nil
}
func (s *SongStore) UpdateSong(ctx context.Context, songId int64, song *Song) error {
	conditions := []string{}
	args := []any{}
	argsIdx := 1
	if song.BPM > 0 {
		conditions = append(conditions, fmt.Sprintf("bpm = $%d", argsIdx))
		args = append(args, song.BPM)
		argsIdx++
	}
	if len(song.Sections) > 0 {
		conditions = append(conditions, fmt.Sprintf("sections = $%id", argsIdx))
		args = append(args, song.Sections)
		argsIdx++
	}
	query := fmt.Sprintf(
		`
				UPDATE songs
				SET %s	
				WHERE id = $%d	
				RETURNING id, name, bpm, sections
			`, strings.Join(conditions, ", "), argsIdx,
	)
	args = append(args, songId)
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	err := s.db.QueryRowContext(ctx, query, args...).Scan(
		&song.ID,
		&song.Name,
		&song.BPM,
		&song.Sections,
	)
	if err != nil {
		return ErrorFactoryDB(err)
	}
	return nil
}

func (s *SongStore) DeleteSong(ctx context.Context, songID int64) error {
	query := `
			DELETE  from songs
			WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	_, err := s.db.ExecContext(ctx, query, songID)
	if err != nil {
		return ErrorFactoryDB(err)
	}
	return nil
}
