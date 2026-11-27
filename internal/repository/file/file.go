package file

import (
	"context"
	"database/sql"
	"fmt"
	"gustav-tgbot/internal/model"
	"log/slog"
)

type File struct {
	db  *sql.DB
	log *slog.Logger
}

func (f *File) GetFilesByUserID(ctx context.Context, userID int64) ([]model.File, error) {
	query := `
		SELECT 
			id,
			filename,
			fileext,
			user_id,
			short_url_original,
			short_url_upgrade,
			created_at
		FROM files
		WHERE user_id = ?
		ORDER BY created_at ASC
	`

	rows, err := f.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("get files by user id: %w", err)
	}
	defer rows.Close()

	var files []model.File

	for rows.Next() {
		var file model.File

		err := rows.Scan(
			&file.ID,
			&file.FileName,
			&file.FileExt,
			&file.UserID,
			&file.ShortUrlOriginal,
			&file.ShortUrlUpgrade,
			&file.CreatedAt,
			//&file.StateFile,
		)
		if err != nil {
			return nil, fmt.Errorf("scan file: %w", err)
		}
		files = append(files, file)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate files: %w", err)
	}

	return files, nil
}

func (f *File) CreateFile(ctx context.Context, file *model.File) (int64, error) {
	query := `
        INSERT INTO files (filename, fileext, user_id, short_url_original, short_url_upgrade, created_at)
        VALUES (?, ?, ?, ?, ?, ?)
    `

	result, err := f.db.ExecContext(ctx, query, file.FileName, file.FileExt, file.UserID, file.ShortUrlOriginal, file.ShortUrlUpgrade, file.CreatedAt)
	if err != nil {
		return 0, fmt.Errorf("create file: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get last inserted file ID: %w", err)
	}

	return id, nil
}

func NewFileRepository(db *sql.DB, log *slog.Logger) *File {
	return &File{
		db:  db,
		log: log,
	}
}
