package parser

import (
	"errors"
	"fmt"
	"strings"
)

const (
	Low = iota
	Medium
	High
	Warning
	Critical
)

type MetaInfo struct {
	Severity    int
	Description string
	Code        string
}

func lastIndex(str string, toFind byte) (int, error) {
	for i := len(str) - 1; i > 0; i-- {
		if str[i] == toFind {
			return i, nil
		}
	}
	return -1, fmt.Errorf("cannot find %s from from \"%s\"", string(toFind), str)
}

func parseSeverity(raw string) (int, error) {
	severityMap := map[string]int{
		"LOW":      Low,
		"MEDIUM":   Medium,
		"High":     High,
		"Warning":  Warning,
		"Critical": Critical,
	}

	if _, exist := severityMap[raw]; !exist {
		return -1, fmt.Errorf("severity level %s is not exist", raw)
	}

	return severityMap[raw], nil
}

func parseDescriptionAndCode(raw string) (string, string, error) {
	if raw[len(raw)-1] != ']' {
		return "", "", fmt.Errorf("\"%s\"does not end with ']'", raw)
	}

	bracketStart, err := lastIndex(raw, '[')
	if err != nil {
		return "", "", err
	}

	//exclude brackets

	description := strings.TrimSpace(raw[:bracketStart])
	code := raw[bracketStart+1 : len(raw)-1]

	return description, code, nil
}

func ParseMetaInfo(line string) (MetaInfo, error) {

	splitByColon := strings.Split(line, ":")
	if len(splitByColon) < 2 {
		return MetaInfo{}, errors.New("unexpected vulnerability format")
	}

	severity, err := parseSeverity(splitByColon[0])
	if err != nil {
		return MetaInfo{}, err
	}

	description, code, err := parseDescriptionAndCode(splitByColon[1])
	if err != nil {
		return MetaInfo{}, err
	}

	return MetaInfo{
		Severity:    severity,
		Description: description,
		Code:        code,
	}, nil
}
