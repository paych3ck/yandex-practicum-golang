package api

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/paych3ck/final_project/pkg/db"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	writeJSON(w, http.StatusOK, task)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := decodeJSON(r, &task); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if task.ID == 0 {
		writeError(w, http.StatusBadRequest, errors.New("task id is required"))
		return
	}

	normalizeTask(&task)
	if err := prepareTask(&task); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, errors.New("task not found"))
			return
		}

		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.FormValue("id"))
	if id == "" {
		writeError(w, http.StatusBadRequest, errors.New("task id is required"))
		return
	}

	if err := db.DeleteTask(id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, errors.New("task not found"))
			return
		}

		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{})
}
