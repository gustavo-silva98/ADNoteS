package file

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

type Writer interface {
	InsertNote(n *Note, ctx context.Context) (int64, error)
	QueryNote(limit int, offset int, ctx context.Context) (map[int]Note, error)
	UpdateEditNoteRepository(ctx context.Context, note Note) (int64, error)
	DeleteNoteRepository(ctx context.Context, noteId int) (int64, error)
}

type SqliteHandler struct {
	DbPath    string
	TableName string
	DB        *sql.DB
}

type Note struct {
	ID           int
	Hour         int64
	NoteText     string
	Reminder     int
	PlusReminder int
}

func InitDB(pathString string, ctx context.Context) (*SqliteHandler, error) {
	db, err := sql.Open("sqlite", pathString)
	if err != nil {
		return nil, err
	}
	_, err = db.ExecContext(
		ctx,
		`CREATE TABLE IF NOT EXISTS notas (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			hour INTEGER NOT NULL,
			note_text TEXT NOT NULL,
			reminder INTEGER,
			plusreminder INTEGER
		)`,
	)
	if err != nil {
		return nil, err
	}
	return &SqliteHandler{
		DbPath:    pathString,
		TableName: "notas",
		DB:        db,
	}, nil
}

func (s SqliteHandler) InsertNote(n *Note, ctx context.Context) (int64, error) {

	res, err := s.DB.ExecContext(
		ctx,
		`INSERT INTO notas (hour, note_text, reminder, plusreminder) VALUES (?, ?, ?, ?)`,
		n.Hour, n.NoteText, n.Reminder, n.PlusReminder,
	)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s SqliteHandler) QueryNote(limit int, offset int, ctx context.Context) (map[int]Note, error) {
	rows, err := s.DB.QueryContext(
		ctx,
		`SELECT * FROM notas ORDER BY id DESC LIMIT ? OFFSET ?`,
		limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var queryMap = map[int]Note{}
	for rows.Next() {
		var note Note
		err := rows.Scan(&note.ID, &note.Hour, &note.NoteText, &note.Reminder, &note.PlusReminder)
		if err != nil {
			return nil, err
		}
		queryMap[note.ID] = note
	}

	return queryMap, nil

}

func WriteTxt(msg string) {
	f, err := os.OpenFile("notes.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Erro ao abrir o arquivo: %v", err)
	}

	defer f.Close()

	_, err = f.WriteString(msg + "\n")
	if err != nil {
		log.Fatalf("Erro ao escrever no arquivo: %v", err)
	}

	log.Println("Arquivo escrito com sucesso")
}

func (s SqliteHandler) GetFirsIndexPage(ctx context.Context) (int, error) {
	row, err := s.DB.QueryContext(
		ctx,
		fmt.Sprintf(`SELECT COUNT(*) FROM %v`, s.TableName),
	)

	if err != nil {
		return 0, err
	}
	defer row.Close()
	var count int
	for row.Next() {
		if err := row.Scan(&count); err != nil {
			return 0, nil
		}
	}
	if count < 10 {
		return 10, nil
	}
	return count, nil
}

func (s SqliteHandler) UpdateEditNoteRepository(ctx context.Context, note Note) (int64, error) {
	row, err := s.DB.ExecContext(
		ctx,
		`UPDATE notas
		SET hour = ?, note_text = ?, reminder = ?, plusreminder = ?
		WHERE id = ?`,
		note.Hour, note.NoteText, note.Reminder, note.PlusReminder, note.ID)
	if err != nil {
		return 0, err
	}
	ra, err := row.RowsAffected()
	if err != nil {
		return 0, err
	}

	return ra, nil

}

func (s SqliteHandler) DeleteNoteRepository(ctx context.Context, noteId int) (int64, error) {

	row, err := s.DB.ExecContext(ctx, `DELETE FROM notas WHERE id = ?`, noteId)
	if err != nil {
		return 0, err
	}
	ra, err := row.RowsAffected()
	if err != nil {
		return 0, err
	}

	return ra, nil
}
