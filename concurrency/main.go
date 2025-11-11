package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	SIZE   = 100_000_000
	CHUNKS = 8
)

func generateRandomElements(size int) []int {
	if size <= 0 {
		return nil
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	data := make([]int, size)
	for i := range data {
		data[i] = r.Int()
	}
	return data
}

func maximum(data []int) int {
	if len(data) == 0 {
		return 0
	}
	max := data[0]
	for i := 1; i < len(data); i++ {
		if data[i] > max {
			max = data[i]
		}
	}
	return max
}

func maxChunks(data []int) int {
	n := len(data)
	if n == 0 {
		return 0
	}

	chunks := CHUNKS
	if n < chunks {
		chunks = n
	}

	maxes := make([]int, chunks)
	var wg sync.WaitGroup
	wg.Add(chunks)

	base := n / chunks
	rem := n % chunks

	start := 0
	for i := 0; i < chunks; i++ {
		size := base
		if i < rem {
			size++
		}
		end := start + size

		chunk := data[start:end]
		go func(chunk []int, i int) {
			defer wg.Done()
			maxes[i] = maximum(chunk)
		}(chunk, i)

		start = end
	}

	wg.Wait()

	return maximum(maxes)
}

func main() {
	fmt.Printf("Генерируем %d целых чисел", SIZE)
	data := generateRandomElements(SIZE)

	fmt.Println("Ищем максимальное значение в один поток")
	start := time.Now()
	max := maximum(data)
	elapsed := time.Since(start).Milliseconds()

	fmt.Printf("Максимальное значение элемента: %d\nВремя поиска: %d ms\n", max, elapsed)

	fmt.Printf("Ищем максимальное значение в %d потоков", CHUNKS)
	start = time.Now()
	max = maxChunks(data)
	elapsed = time.Since(start).Milliseconds()

	fmt.Printf("Максимальное значение элемента: %d\nВремя поиска: %d ms\n", max, elapsed)
}
