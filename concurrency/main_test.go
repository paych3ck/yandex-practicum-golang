package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRandomElements(t *testing.T) {
	tests := []struct {
		name string
		n    int
	}{
		{"zero", 0},
		{"negative", -5},
		{"small", 5},
		{"large", 10_000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateRandomElements(tt.n)
			if tt.n <= 0 {
				assert.Nil(t, got, "ожидаем nil для n<=0")
				return
			}
			assert.Len(t, got, tt.n)
		})
	}
}

func TestMaximum(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		want int
	}{
		{"empty", nil, 0},
		{"single", []int{42}, 42},
		{"increasing", []int{1, 2, 3, 4, 5}, 5},
		{"decreasing", []int{9, 7, 3, 1}, 9},
		{"mixed", []int{-10, 0, 5, -3, 4}, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, maximum(tt.in))
		})
	}
}

func TestMaxChunks_EqualsSequential(t *testing.T) {
	data := generateRandomElements(50_000)
	seq := maximum(data)
	par := maxChunks(data)
	assert.Equal(t, seq, par, "параллельный результат должен совпадать с последовательным")
}
