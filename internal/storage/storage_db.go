package storage

import (
	"context"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StorageDB struct {
	pool *pgxpool.Pool
}

func (db *StorageDB) SaveLongURL(ctx context.Context, longURL, userID string) (shortURL string, err error) {
	shortURL, err = GenerateShortLink(longURL)
	if err != nil {
		return "", err
	}

	if userID == "" {
		err = saveLink(ctx, db, longURL, shortURL)
	} else {
		err = saveLinkWithUser(ctx, db, longURL, shortURL, userID)
	}
	if err != nil {
		return shortURL, err
	}
	return shortURL, nil
}

func (db *StorageDB) GetLongURL(ctx context.Context, shortURL string) (longURL string, err error) {
	sql := `SELECT (long_url) FROM links WHERE short_url = $1`
	err = db.pool.QueryRow(context.Background(), sql, shortURL).Scan(&longURL)

	switch err {
	case nil:
		return longURL, nil
	case pgx.ErrNoRows:
		return "", err
	default:
		log.Print("failed to fetch long url", err)
		return "", err

	}
}

func (db *StorageDB) CreateUser(ctx context.Context) string {
	uuid := uuid.New().String()
	sql := `INSERT INTO users(uuid) VALUES ($1) ON CONFLICT(uuid) DO NOTHING`
	_, err := db.pool.Exec(context.Background(), sql, uuid)
	if err != nil {
		log.Fatal("failed to create user", err)
	}
	return uuid
}

func (db *StorageDB) GetUserLinks(ctx context.Context, userID string) (links []string, err error) {
	sql := "SELECT short_url FROM links WHERE id in (SELECT link_id FROM user_links WHERE user_id = $1);"
	rows, _ := db.pool.Query(context.Background(), sql, userID)

	for rows.Next() {
		var shortURL string
		err := rows.Scan(&shortURL)
		if err != nil {
			return links, err
		}
		links = append(links, shortURL)
	}

	return links, nil
}

func (db *StorageDB) SaveBatch(ctx context.Context, records []BatchInput) ([]BatchOutput, error) {
	var output []BatchOutput

	batch := &pgx.Batch{}
	for _, record := range records {
		shortURL, err := GenerateShortLink(record.OriginalURL)
		if err != nil {
			return output, nil
		}

		batch.Queue(`INSERT INTO links(long_url, short_url) VALUES ($1, $2) ON CONFLICT(long_url) DO NOTHING`, record.OriginalURL, shortURL)

		batchOutput := BatchOutput{
			CorrelationID: record.CorrelationID,
			ShortURL:      shortURL,
		}
		output = append(output, batchOutput)
	}

	br := db.pool.SendBatch(context.Background(), batch)

	_, err := br.Exec()
	if err != nil {
		return output, err
	}

	return output, nil
}

func (db *StorageDB) Ping(ctx context.Context) error {
	return db.pool.Ping(context.Background())
}

func (db *StorageDB) DeleteLink(ctx context.Context, longLink string) error {
	var linkID string

	sql := `SELECT (id) FROM links WHERE long_url = $1`
	_ = db.pool.QueryRow(context.Background(), sql, longLink).Scan(&linkID)

	sql = `DELETE FROM user_links WHERE link_id = $1`
	db.pool.Exec(context.Background(), sql, linkID)

	sql = `DELETE FROM links WHERE id = $1`
	_, err := db.pool.Exec(context.Background(), sql, linkID)

	if err != nil {
		return err
	}
	return nil
}

func saveLink(ctx context.Context, db *StorageDB, longLink, shortLink string) error {
	sql := `INSERT INTO links(long_url, short_url) VALUES ($1, $2) ON CONFLICT(long_url) DO NOTHING`
	code, err := db.pool.Exec(context.Background(), sql, longLink, shortLink)
	conflict := strings.Split(code.String(), " ")[2] == "0"
	if conflict {
		return NewLinkExistError(shortLink)
	}
	if err != nil {
		return err
	}
	return nil
}

func saveLinkWithUser(ctx context.Context, db *StorageDB, longLink, shortLink, userID string) error {
	var linkID string
	sql := `SELECT (id) FROM links WHERE short_url = $1`
	err := db.pool.QueryRow(context.Background(), sql, shortLink).Scan(&linkID)
	conflict := false
	switch err {
	case nil:
		conflict = true
	case pgx.ErrNoRows:
		sql = `INSERT INTO links(long_url, short_url) VALUES ($1, $2) RETURNING id`
		err = db.pool.QueryRow(context.Background(), sql, longLink, shortLink).Scan(&linkID)
		if err != nil {
			return err
		}
	default:
		return err
	}

	sql = `INSERT INTO user_links(user_id, link_id) VALUES ($1, $2) ON CONFLICT(user_id, link_id) DO NOTHING`
	_, err = db.pool.Exec(context.Background(), sql, userID, linkID)
	if err != nil {
		return err
	}

	if conflict {
		return NewLinkExistError(shortLink)
	}
	return nil
}
