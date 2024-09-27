package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	K1 = 0.035
	K2 = 0.029
)

var (
	Format     = "20060102 15:04:05"
	StepLength = 0.65
	Weight     = 75.0
	Height     = 1.75
	Speed      = 1.39
)

func parsePackage(data string) (t time.Time, steps int, ok bool) {
	ds := strings.Split(data, ",")
	var err error

	if len(ds) == 2 {
		t, err = time.Parse(Format, ds[0])
		if err != nil {
			return time.Time{}, 0, false
		}

		steps, err = strconv.Atoi(ds[1])
		if err != nil || steps < 0 {
			return time.Time{}, 0, false
		}

		ok = true
		return t, steps, ok
	}
	return time.Time{}, 0, false
}

func stepsDay(storage []string) int {
	stepsSum := 0
	for _, pkg := range storage {
		_, steps, _ := parsePackage(pkg)
		stepsSum += steps
	}
	return stepsSum
}

func calories(distance float64) float64 {
	return (K1*Weight + (Speed*Speed/Height)*K2*Weight) * ((distance / Speed) / 60) * 1000
}

func achievement(distance float64) string {
	if distance >= 6.5 {
		return "Отличный результат! Цель достигнута."
	} else if distance >= 3.9 {
		return "Неплохо! День был продуктивный."
	} else if distance >= 2.0 {
		return "Завтра наверстаем!"
	} else {
		return "Лежать тоже полезно. Главное — участие, а не победа!"
	}
}

func showMessage(s string) {
	fmt.Printf("%s\n\n", s)
}

func AcceptPackage(data string, storage []string) []string {
	t, steps, ok := parsePackage(data)
	if !ok {
		showMessage("ошибочный формат пакета")
		return storage
	}

	if steps == 0 {
		return storage
	}

	now := time.Now().UTC()

	if t.Day() != now.Day() {
		showMessage("неверный день")
		return storage
	}

	if t.After(now) {
		showMessage("некорректное значение времени")
		return storage
	}

	if len(storage) > 0 {
		lastIndex := len(storage) - 1

		if data[:len(Format)] <= storage[lastIndex][:len(Format)] {
			showMessage("некорректное значение времени")
			return storage
		}

		if data[:8] != storage[len(storage)-1][:8] {
			storage = storage[:0]
		}
	}
	storage = append(storage, data)
	totalSteps := stepsDay(storage)
	totalDist := float64(totalSteps) * StepLength / 1000
	totalCalories := calories(totalDist)
	text := achievement(totalDist)

	msg := fmt.Sprintf(`Время: %s.
Количество шагов за сегодня: %d.
Дистанция составила %.2f км.
Вы сожгли %.2f ккал.
%s`, t.Format("15:04:05"), totalSteps, totalDist, totalCalories, text)
	showMessage(msg)

	return storage
}

func main() {
	now := time.Now().UTC()
	today := now.Format("20060102")

	input := []string{
		"01:41:03,-100",
		",3456",
		"12:40:00, 3456 ",
		"something is wrong",
		"02:11:34,678",
		"02:11:34,792",
		"17:01:30,1078",
		"03:25:59,7830",
		"04:00:46,5325",
		"04:45:21,3123",
	}

	var storage []string
	storage = AcceptPackage("20230720 00:11:33,100", storage)
	for _, v := range input {
		storage = AcceptPackage(today+" "+v, storage)
	}
}
