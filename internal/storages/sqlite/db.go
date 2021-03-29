package sqllite

import (
	"context"
	"database/sql"

	"github.com/manabie-com/togo/internal/storages"
)

// LiteDB for working with sqllite
type LiteDB struct {
	DB *sql.DB
}

// RetrieveTasks returns tasks if match userID AND createDate.
func (l *LiteDB) RetrieveTasks(ctx context.Context, userID, createdDate sql.NullString) ([]*storages.Task, error) {
	stmt := `SELECT id, content, user_id, created_date FROM tasks WHERE user_id = ? AND created_date = ?`
	rows, err := l.DB.QueryContext(ctx, stmt, userID, createdDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*storages.Task
	for rows.Next() {
		t := &storages.Task{}
		err := rows.Scan(&t.ID, &t.Content, &t.UserID, &t.CreatedDate)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// AddTask adds a new task to DB
func (l *LiteDB) AddTask(ctx context.Context, t *storages.Task) error {
	stmt := `INSERT INTO tasks (id, content, user_id, created_date) VALUES (?, ?, ?, ?)`
	_, err := l.DB.ExecContext(ctx, stmt, &t.ID, &t.Content, &t.UserID, &t.CreatedDate)
	if err != nil {
		return err
	}

	return nil
}

// ValidateUser returns tasks if match userID AND password
func (l *LiteDB) ValidateUser(ctx context.Context, userID, pwd sql.NullString) bool {
	stmt := `SELECT id FROM users WHERE id = ? AND password = ?`
	row := l.DB.QueryRowContext(ctx, stmt, userID, pwd)
	u := &storages.User{}
	err := row.Scan(&u.ID)
	if err != nil {
		return false
	}

	return true
}

// Add by VienNV start
// CheckIsLimitReachedTask returns tasks created is reached or not
func (l *LiteDB) CheckIsLimitReachedTask(ctx context.Context, userID, createdDate string) bool {
	stmt := `SELECT COUNT(tasks.id) AS TasksCount, users.max_todo AS MaxTodo
			 FROM tasks join users ON tasks.user_id = users.id 
		 	 WHERE users.id = ? AND tasks.created_date = ?`
	row := l.DB.QueryRowContext(ctx, stmt, userID, createdDate)
	reachedTask := struct {
		TasksCount int
		MaxTodo int
	}{}
	err := row.Scan(&reachedTask.TasksCount, &reachedTask.MaxTodo)
	if err != nil {
		return false
	}

	if reachedTask.TasksCount >= reachedTask.MaxTodo {
		return true
	} else {
		return false
	}
}
// Add by VienNV end