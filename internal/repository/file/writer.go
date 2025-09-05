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
	NoteText     string
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

	res, err := s.db.ExecContext(
		ctx,
		`INSERT INTO notas (hour, note, reminder, plusreminder) VALUES (?, ?, ?, ?)`,
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

func (s SqliteHandler) QueryNote(firstId int, lastId int, ctx context.Context) (map[int]Note, error) {
	rows, err := s.db.QueryContext(
		ctx,
		`SELECT * FROM notas WHERE id BETWEEN (?) AND (?)`,
		firstId, lastId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var queryMap = map[int]Note{}
	for rows.Next() {
		var note NoteRow
		err := rows.Scan(&note.ID, &note.Hour, &note.NoteText, &note.Reminder, &note.PlusReminder)
		if err != nil {
			return nil, err
		}
		queryMap[note.ID] = note.Note

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
	row, err := s.db.QueryContext(
		ctx,
		`SELECT COUNT(*) FROM (?)`, s.TableName,
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
	return count, nil
}
