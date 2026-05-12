package display

import (
	"fmt"
	"sort"
	"strings"
)

const labelWidth = 14

func Print(
	current map[string]float64,
	last map[string]float64,
	lastTime string,
	label string,
) {
	fmt.Printf("%-*s %s\n", labelWidth, "余额:", formatLine(current))

	if len(last) == 0 {
		return
	}

	diff := make(map[string]float64)
	for c, v := range last {
		if cv, ok := current[c]; ok {
			diff[c] = cv - v
		}
	}

	if len(diff) == 0 {
		return
	}

	timeStr := formatTime(lastTime)
	fmt.Printf("%-*s %s (%s)\n", labelWidth, label+":", formatLine(diff), timeStr)
}

func formatLine(items map[string]float64) string {
	currencies := make([]string, 0, len(items))
	for c := range items {
		currencies = append(currencies, c)
	}
	sort.Strings(currencies)

	parts := make([]string, len(currencies))
	for i, c := range currencies {
		parts[i] = fmt.Sprintf("%*.2f  %s", 8, items[c], c)
	}
	return strings.Join(parts, "  |  ")
}

func formatTime(iso string) string {
	if len(iso) >= 10 {
		return iso[5:10]
	}
	return iso
}
