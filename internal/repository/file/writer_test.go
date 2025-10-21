package file_test

import (
	"context"
	"database/sql"
	"log"
	"testing"

	"github.com/gustavo-silva98/adnotes/internal/repository/file"
)

func setupTestDB(t *testing.T) *file.SqliteHandler {
	t.Helper()
	db_path := ":memory:"
	ctx := context.Background()

	handler, err := file.InitDB(db_path, ctx)
	if err != nil {
		t.Fatalf("Erro ao inicializar banco de teste - %v", err)
	}

	_, err = handler.DB.ExecContext(ctx, "DELETE FROM notas")
	if err != nil {
		t.Fatalf("Erro ao limpar o banco de teste - %v", err)
	}

	return handler
}

func TestInsertNote(t *testing.T) {
	handler := setupTestDB(t)
	ctx := context.Background()

	note := &file.Note{
		Hour:         int64(1),
		NoteText:     "Teste de insert",
		Reminder:     1,
		PlusReminder: 2,
	}

	id, err := handler.InsertNote(note, ctx)
	if err != nil {
		t.Fatalf("Falha ao inserir nota teste - %v", err)
	}
	if id <= 0 {
		t.Errorf("ID de insert Note retornado inválido - %d", id)
	}
}

func TestQueryNote(t *testing.T) {
	handler := setupTestDB(t)
	ctx := context.Background()

	note := &file.Note{
		Hour:         int64(1),
		NoteText:     "Teste de insert",
		Reminder:     1,
		PlusReminder: 2,
	}

	_, _ = handler.InsertNote(note, ctx)

	result, err := handler.QueryNote(10, 0, ctx)
	if err != nil {
		t.Fatalf("Erro ao consultar notas - %v", err)
	}

	if len(result) == 0 {
		t.Errorf("Nenhuma nota retornada")
	}
}

func TestUpdateEditNote(t *testing.T) {
	handler := setupTestDB(t)
	ctx := context.Background()

	note := &file.Note{
		Hour:         int64(1),
		NoteText:     "Teste de insert",
		Reminder:     1,
		PlusReminder: 2,
	}

	id, _ := handler.InsertNote(note, ctx)
	note.ID = int(id)
	note.NoteText = "Nota atualizada"

	rowsAffected, err := handler.UpdateEditNoteRepository(ctx, *note)
	if err != nil {
		t.Fatalf("Erro ao atualizar nota - %v", err)
	}
	if rowsAffected != 1 {
		t.Errorf("Esperado 1 linha afetada. Resultado divergente")
	}

}

func TestDeleteNoteRepository(t *testing.T) {
	ctx := context.Background()
	handler, _ := file.InitDB(":memory:", ctx)
	note := &file.Note{
		Hour:         int64(1),
		NoteText:     "Teste de insert",
		Reminder:     1,
		PlusReminder: 2,
	}
	id, _ := handler.InsertNote(note, ctx)

	rowsAffected, err := handler.DeleteNoteRepository(ctx, int(id))
	if err != nil {
		t.Fatalf("Falha ao deletar nota no teste - %v", err)
	}
	if rowsAffected != 1 {
		t.Errorf("Era esperado somente 1 linha. Valor divergente.")
	}
}

func BenchmarkInsertNote(b *testing.B) {
	ctx := context.Background()
	handler, _ := file.InitDB(":memory:", ctx)

	for i := 0; i < b.N; i++ {
		note := &file.Note{
			Hour:         int64(i),
			NoteText:     "Benchmark",
			Reminder:     1,
			PlusReminder: 2,
		}
		_, _ = handler.InsertNote(note, ctx)
	}

}

func FTSTableExists(db *file.SqliteHandler) (bool, error) {
	query := `SELECT name FROM sqlite_master WHERE type='table' AND name='notes_fts'`

	row := db.DB.QueryRow(query)
	var name string
	err := row.Scan(&name)

	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func TestCreateFTSTable(t *testing.T) {
	db, _ := sql.Open("sqlite", "teste_db.db")
	handler := &file.SqliteHandler{
		DbPath:    "teste_db.db",
		TableName: "teste_notas",
		DB:        db,
	}

	err := handler.CreateFTSTable()
	if err != nil {
		log.Fatalf("Erro ao criar tabela FTS: %v", err)
	}
	exists, err := FTSTableExists(handler)
	if err != nil {
		log.Fatalf("Erro ao verificar tabela FTS: %v", err)
	}

	if exists {
		log.Println("Tabela FTS Criada com sucesso")
	} else {
		log.Println("Tabela FTS não foi criada")
	}
}
