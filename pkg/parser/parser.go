package parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type MetaInfo struct {
	Severity    string
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

	splitByColon := strings.SplitN(line, ":", 2)
	if len(splitByColon) < 2 {
		return MetaInfo{}, errors.New("unexpected vulnerability format")
	}

	severity := strings.TrimSpace(splitByColon[0])
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

type Location struct {
	Path string
	Line int
}

func ParseLocation(raw string) (Location, error) {

	splitByColon := strings.SplitN(raw, ":", 2)
	if len(splitByColon) < 2 {
		return Location{}, errors.New("unexpected location format")
	}

	pathAndLine := strings.TrimSpace(splitByColon[1])
	lastColon, err := lastIndex(pathAndLine, ':')
	if err != nil {
		return Location{}, err
	}

	path := pathAndLine[:lastColon]

	lineRaw := pathAndLine[lastColon+1:]
	line, err := strconv.Atoi(lineRaw)
	if err != nil {
		return Location{}, err
	}

	return Location{
		Path: path,
		Line: line,
	}, nil
}

func ParseSnippet(raw string) (string, error) {
	splitByColon := strings.SplitN(raw, " ", 2)
	if len(splitByColon) < 2 {
		return "", errors.New("unexpected snippet format")
	}

	return strings.TrimSpace(splitByColon[1]), nil
}
