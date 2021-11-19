package db

// schema.go provides data models in DB
import (
	"time"
)

// Task corresponds to a row in `tasks` table
type Task struct {
	ID        uint64    `db:"id"`
	Title     string    `db:"title"`
	CreatedAt time.Time `db:"created_at"`
	Deadline  time.Time `db:"deadline"`
	IsDone    bool      `db:"is_done"`
}

// User corresponds to a row in `users` table
type User struct {
	ID   uint64 `db:"id"`
	Name string `db:"name"`
	Pwd  string `db:"pwd"`
}

// Owner corresponds to a row in `task_owners` table
type Owner struct {
	Task uint64 `db:"task_id"`
	User uint64 `db:"user_id"`
}