package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Task struct {
	ID      int64  `db:"id" json:"id,string"`
	Date    string `db:"date" json:"date"`
	Title   string `db:"title" json:"title"`
	Comment string `db:"comment" json:"comment"`
	Repeat  string `db:"repeat" json:"repeat"`
}

const (
	addTaskQuery = `INSERT INTO scheduler(date, title, comment, repeat)
					VALUES(?, ?, ?, ?);`
	getTaskQuery     = `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	getAllTasksQuery = `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC, id ASC LIMIT ?`
	deleteTaskQuery  = `DELETE FROM scheduler WHERE id = ?`
	updateTaskQuery  = `UPDATE scheduler
					   SET date = ?, title = ?, comment = ?, repeat = ?
					   WHERE id = ?;`
	updateTaskDateQuery = `UPDATE scheduler SET date = ? WHERE id = ?`
)

func AddTask(task *Task) (int64, error) {
	if task == nil {
		return 0, errors.New("task is nil")
	}

	if conn == nil {
		return 0, errors.New("database is not initialized")
	}

	res, err := conn.Exec(addTaskQuery, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func Tasks(limit int) ([]*Task, error) {
	if conn == nil {
		return nil, errors.New("database is not initialized")
	}

	rows, err := conn.Query(getAllTasksQuery, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, err
		}

		result = append(result, &t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if result == nil {
		return []*Task{}, nil
	}

	return result, nil
}

func GetTask(id string) (*Task, error) {
	if conn == nil {
		return nil, errors.New("database is not initialized")
	}
	taskID, err := parseTaskID(id)
	if err != nil {
		return nil, err
	}

	row := conn.QueryRow(getTaskQuery, taskID)
	var t Task
	if err := row.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
		return nil, err
	}
	return &t, nil
}

func UpdateTask(task *Task) error {
	if task == nil {
		return errors.New("task is nil")
	}

	if conn == nil {
		return errors.New("database is not initialized")
	}

	if task.ID == 0 {
		return errors.New("task id is required")
	}

	res, err := conn.Exec(updateTaskQuery, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func DeleteTask(id string) error {
	if conn == nil {
		return errors.New("database is not initialized")
	}

	taskID, err := parseTaskID(id)
	if err != nil {
		return err
	}

	res, err := conn.Exec(deleteTaskQuery, taskID)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func UpdateTaskDate(id int64, nextDate string) error {
	if conn == nil {
		return errors.New("database is not initialized")
	}

	if id <= 0 {
		return errors.New("task id is required")
	}

	if strings.TrimSpace(nextDate) == "" {
		return errors.New("next date is required")
	}

	res, err := conn.Exec(updateTaskDateQuery, nextDate, id)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func parseTaskID(id string) (int64, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return 0, errors.New("task id is required")
	}

	taskID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid task id: %w", err)
	}

	return taskID, nil
}
