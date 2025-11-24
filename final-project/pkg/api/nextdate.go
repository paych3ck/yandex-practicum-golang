package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/paych3ck/final_project/pkg/nextdate"
)

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowParam := strings.TrimSpace(r.FormValue("now"))
	dateParam := strings.TrimSpace(r.FormValue("date"))
	repeatParam := strings.TrimSpace(r.FormValue("repeat"))

	now := normalizeNow(nowParam)
	if now.IsZero() {
		http.Error(w, fmt.Sprintf("invalid now date: %q", nowParam), http.StatusBadRequest)
		return
	}

	if dateParam == "" {
		http.Error(w, "date parameter is required", http.StatusBadRequest)
		return
	}

	if repeatParam == "" {
		http.Error(w, "repeat parameter is required", http.StatusBadRequest)
		return
	}

	next, err := nextdate.NextDate(now, dateParam, repeatParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = fmt.Fprint(w, next)
}

func normalizeNow(nowParam string) time.Time {
	if nowParam == "" {
		return truncateToDay(time.Now().UTC())
	}

	parsed, err := time.Parse(dateLayout, nowParam)
	if err != nil {
		return time.Time{}
	}

	return truncateToDay(parsed)
}

func truncateToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}
