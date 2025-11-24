package api

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/paych3ck/final_project/pkg/db"
	"github.com/paych3ck/final_project/pkg/nextdate"
)

func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
		return
	}

	id := strings.TrimSpace(r.FormValue("id"))
	if id == "" {
		writeError(w, http.StatusBadRequest, errors.New("task id is required"))
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, errors.New("task not found"))
			return
		}

		writeError(w, http.StatusBadRequest, err)
		return
	}

	if strings.TrimSpace(task.Repeat) == "" {
		if err := db.DeleteTask(id); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				writeError(w, http.StatusNotFound, errors.New("task not found"))
				return
			}

			writeError(w, http.StatusInternalServerError, err)
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{})
		return
	}

	now := truncateToDay(time.Now().UTC())
	next, err := nextdate.NextDate(now, task.Date, task.Repeat)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if err := db.UpdateTaskDate(task.ID, next); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, errors.New("task not found"))
			return
		}

		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{})
}
