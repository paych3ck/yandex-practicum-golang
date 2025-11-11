package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
	}
}

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	city := "moscow"
	total := len(cafeList[city])

	requests := []struct {
		count int
		want  int
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{100, func() int {
			if total < 100 {
				return total
			}
			return 100
		}()},
	}

	for _, v := range requests {
		u := "/cafe?city=" + city + "&count=" + strconv.Itoa(v.count)

		rr := httptest.NewRecorder()
		req := httptest.NewRequest("GET", u, nil)
		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code, "count=%d", v.count)

		body := strings.TrimSpace(rr.Body.String())
		var got int
		if body == "" {
			got = 0
		} else {
			parts := strings.Split(body, ",")
			for i := range parts {
				parts[i] = strings.TrimSpace(parts[i])
			}
			got = len(parts)
		}

		assert.Equal(t, v.want, got, "count=%d; body=%q", v.count, body)
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		search    string
		wantCount int
	}{
		{"фасоль", 0},
		{"кофе", 2},
		{"вилка", 1},
	}

	for _, v := range requests {
		u := "/cafe?city=moscow&search=" + url.QueryEscape(v.search)

		rr := httptest.NewRecorder()
		req := httptest.NewRequest("GET", u, nil)
		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code, "search=%q", v.search)

		body := strings.TrimSpace(rr.Body.String())

		var items []string
		if body == "" {
			items = []string{}
		} else {
			items = strings.Split(body, ",")
		}

		assert.Equal(t, v.wantCount, len(items), "search=%q; body=%q", v.search, body)

		needle := strings.ToLower(v.search)
		for _, name := range items {
			name = strings.TrimSpace(name)
			lowerName := strings.ToLower(name)
			assert.Contains(t, lowerName, needle, "search=%q; cafe=%q", v.search, name)
		}
	}
}
