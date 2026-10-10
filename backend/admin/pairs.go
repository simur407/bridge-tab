package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type pairLine struct {
	Table *int
	NS    int
	EW    int
}

func parsePairLines(text string) ([]pairLine, error) {
	var pairs []pairLine
	for rawIndex, raw := range strings.Split(text, "\n") {
		lineNo := rawIndex + 1
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}

		var table *int
		body := line
		if left, right, ok := strings.Cut(line, ":"); ok && !strings.Contains(left, ";") {
			n, err := strconv.Atoi(strings.TrimSpace(left))
			if err != nil || n < 1 {
				return nil, fmt.Errorf("linia %d: zły numer stołu", lineNo)
			}
			table = &n
			body = right
		}

		parts := strings.Split(body, ";")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		switch len(parts) {
		case 3:
			if table != nil {
				return nil, fmt.Errorf("linia %d: podaj stół raz", lineNo)
			}
			n, err := strconv.Atoi(parts[0])
			if err != nil || n < 1 {
				return nil, fmt.Errorf("linia %d: zły numer stołu", lineNo)
			}
			table = &n
			parts = parts[1:]
		case 2:
		default:
			return nil, fmt.Errorf("linia %d: oczekiwano NS;EW albo stół;NS;EW", lineNo)
		}

		ns, errNS := strconv.Atoi(parts[0])
		ew, errEW := strconv.Atoi(parts[1])
		if errNS != nil || errEW != nil || ns < 1 || ew < 1 {
			return nil, fmt.Errorf("linia %d: złe numery drużyn", lineNo)
		}
		if ns == ew {
			return nil, fmt.Errorf("linia %d: drużyna nie może grać sama ze sobą", lineNo)
		}
		pairs = append(pairs, pairLine{Table: table, NS: ns, EW: ew})
	}
	if len(pairs) == 0 {
		return nil, errors.New("podaj co najmniej jedną parę")
	}
	return pairs, nil
}

func parseBoardNos(text string) ([]int, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, errors.New("podaj numery rozdań")
	}
	var nums []int
	for _, chunk := range strings.Split(text, ",") {
		chunk = strings.TrimSpace(chunk)
		if chunk == "" {
			continue
		}
		n, err := strconv.Atoi(chunk)
		if err != nil || n < 1 {
			return nil, fmt.Errorf("zły numer rozdania %q", chunk)
		}
		nums = append(nums, n)
	}
	if len(nums) == 0 {
		return nil, errors.New("podaj numery rozdań")
	}
	return nums, nil
}

func normalizePostedContract(value string) string {
	value = strings.TrimSpace(value)
	if strings.EqualFold(value, "pass") {
		return "Pass"
	}
	return strings.ToUpper(value)
}
