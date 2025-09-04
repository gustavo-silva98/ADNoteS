package file

import (
	"context"
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

type Writer interface {
	InsertNote(noteData string)
	InitDB(DbPath string) error
}

type SqliteHandler struct {
	DbPath    string
	TableName string
	db        *sql.DB
}

type Note struct {
	Hour         int
	Note         string
	Reminder     int
	PlusReminder int
}

type NoteRow struct {
	ID int
	Note
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
			note TEXT NOT NULL,
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
		db:        db,
	}, nil
}

func (s SqliteHandler) InsertNote(n *Note, ctx context.Context) (int64, error) {
	defer s.db.Close()

	res, err := s.db.ExecContext(
		ctx,
		`INSERT INTO notas (hour, note, reminder, plusreminder) VALUES (?, ?, ?, ?)`,
		n.Hour, n.Note, n.Reminder, n.PlusReminder,
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
