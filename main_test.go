package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
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
	requests := []struct {
		count int
		want  int
	}{
		{count: 0, want: 0},
		{count: 1, want: 1},
		{count: 2, want: 2},
		{count: 100, want: min(100, len(cafeList["moscow"]))},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest(
			"GET", fmt.Sprintf("/cafe?count=%d&city=moscow", v.count), nil)
		handler.ServeHTTP(response, req)
		require.Equal(t, http.StatusOK, response.Code)

		result := response.Body.String()
		total := 0
		if result != "" {
			total = len(strings.Split(result, ","))
		}
		assert.Equal(t, v.want, total)
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)
	requests := []struct {
		search    string
		wantCount int
	}{
		{search: "фасоль", wantCount: 0},
		{search: "кофе", wantCount: 2},
		{search: "вилка", wantCount: 1},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest(
			"GET", fmt.Sprintf("/cafe?search=%s&city=moscow", v.search), nil)
		handler.ServeHTTP(response, req)
		require.Equal(t, http.StatusOK, response.Code)

		result := response.Body.String()
		total := 0
		if result != "" {
			cafes := strings.Split(result, ",")
			total = len(cafes)
			for _, cafe := range cafes {
				nameLower := strings.ToLower(cafe)
				searchLower := strings.ToLower(v.search)
				assert.True(t, strings.Contains(nameLower, searchLower))

			}
		}
		assert.Equal(t, v.wantCount, total)
	}
}
