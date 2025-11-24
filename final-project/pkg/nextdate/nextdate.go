package nextdate

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const dateLayout = "20060102"

var (
	ErrEmptyRepeat      = errors.New("repeat rule is empty")
	ErrUnsupportedRule  = errors.New("unsupported repeat rule")
	ErrInvalidRepeat    = errors.New("invalid repeat rule")
	ErrInvalidStartDate = errors.New("invalid start date")
)

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	repeat = strings.TrimSpace(repeat)
	if repeat == "" {
		return "", ErrEmptyRepeat
	}

	start, err := parseDate(dstart)
	if err != nil {
		return "", err
	}

	nowDay := normalizeDate(now)
	ruleParts := strings.Fields(repeat)
	if len(ruleParts) == 0 {
		return "", ErrInvalidRepeat
	}

	next, err := applyRule(start, nowDay, ruleParts)
	if err != nil {
		return "", err
	}

	return next.Format(dateLayout), nil
}

func parseDate(raw string) (time.Time, error) {
	date := strings.TrimSpace(raw)
	t, err := time.Parse(dateLayout, date)
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: %v", ErrInvalidStartDate, err)
	}

	return normalizeDate(t), nil
}

func normalizeDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

func applyRule(start, now time.Time, rule []string) (time.Time, error) {
	switch rule[0] {
	case "y":
		if len(rule) != 1 {
			return time.Time{}, ErrInvalidRepeat
		}
		return advanceUntil(start, now, func(t time.Time) time.Time {
			next := t.AddDate(1, 0, 0)
			if t.Month() == time.February && t.Day() == 29 &&
				next.Month() == time.February && next.Day() == 28 {
				next = next.AddDate(0, 0, 1)
			}
			return next
		}), nil
	case "d":
		if len(rule) != 2 {
			return time.Time{}, ErrInvalidRepeat
		}
		interval, err := strconv.Atoi(rule[1])
		if err != nil || interval <= 0 || interval > 400 {
			return time.Time{}, ErrInvalidRepeat
		}
		return advanceUntil(start, now, func(t time.Time) time.Time {
			return t.AddDate(0, 0, interval)
		}), nil
	default:
		return time.Time{}, fmt.Errorf("%w: %s", ErrUnsupportedRule, rule[0])
	}
}

func advanceUntil(start, now time.Time, step func(time.Time) time.Time) time.Time {
	cur := start
	for {
		cur = step(cur)
		if cur.After(now) {
			return cur
		}
	}
}
