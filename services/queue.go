package services

import (
	"database/sql"
	"trullio-kyc/models"
)

func enqueueTask(db *sql.DB, taskData string) error {
	query := `INSERT INTO task_queue (task_data, status) VALUES($1, $2)`
	_, err := db.Exec(query, taskData, "pending")
	return err
}

func DequeueTask(db *sql.DB) (*models.Task, error) {
	var task models.Task

	query := `
		UPDATE task_queue
	          SET status = 'processing', updated_at = NOW()
	          WHERE id = (
	              SELECT id FROM task_queue
	              WHERE status = 'pending'
	              ORDER BY created_at ASC
	              LIMIT 1
	          )
	          RETURNING id, task_data, status, created_at, updated_at
	`

	row := db.QueryRow(query)
	err := row.Scan(&task.ID, &task.TaskData, &task.Status, &task.CreatedAt, &task.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &task, err
}

func updateTaskStatus(db *sql.DB, taskId int, status string) error {
	query := `UPDATE task_queue SET status = $1, update_at = NOW() WHERE id = $2`
	_, err := db.Exec(query, status, taskId)
	return err
}
