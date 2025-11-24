package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/paych3ck/final_project/pkg/db"
	"github.com/paych3ck/final_project/pkg/nextdate"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := decodeJSON(r, &task); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	normalizeTask(&task)
	if err := prepareTask(&task); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"id": strconv.FormatInt(id, 10),
	})
}

func decodeJSON(r *http.Request, v any) error {
	if r.Body == nil {
		return errors.New("empty request body")
	}

	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	return dec.Decode(v)
}

func normalizeTask(task *db.Task) {
	task.Title = strings.TrimSpace(task.Title)
	task.Comment = strings.TrimSpace(task.Comment)
	task.Repeat = strings.TrimSpace(task.Repeat)
	task.Date = strings.TrimSpace(task.Date)
}

func prepareTask(task *db.Task) error {
	if task.Title == "" {
		return errors.New("title is required")
	}

	now := truncateToDay(time.Now().UTC())
	if task.Date == "" {
		task.Date = now.Format(dateLayout)
	}

	parsed, err := time.Parse(dateLayout, task.Date)
	if err != nil {
		return fmt.Errorf("invalid date: %w", err)
	}

	parsed = truncateToDay(parsed)

	var next string
	if task.Repeat != "" {
		next, err = nextdate.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	if afterNow(now, parsed) {
		if task.Repeat == "" {
			task.Date = now.Format(dateLayout)
		} else {
			task.Date = next
		}
	}

	return nil
}

func afterNow(a, b time.Time) bool {
	return a.After(b)
}
