package api

import (
	"errors"
	"net/http"

	"github.com/paych3ck/final_project/pkg/db"
)

const defaultTasksLimit = 50

type tasksResponse struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
		return
	}

	tasks, err := db.Tasks(defaultTasksLimit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	if tasks == nil {
		tasks = make([]*db.Task, 0)
	}

	writeJSON(w, http.StatusOK, tasksResponse{Tasks: tasks})
}
