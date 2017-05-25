package workflow

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode"
)

func cleanLine(line string) string {
	return removeWS(chopComments(line))
}

func chopComments(line string) string {
	splits := strings.Split(line, "#")
	return splits[0]
}

func removeWS(line string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, line)
}

func parseCommand(line string, lineNo int, rawLine string, parallel bool) (*Step, error) {
	splits := strings.SplitN(line, ":", 2)
	if len(splits) != 2 || splits[1] == "" {
		return nil, fmt.Errorf("parsing error at line %d, '%s' malformed command", lineNo, rawLine)
	}
	step := &Step{Command: splits[1], LineNo: lineNo, Arguments: make(map[string]string), Parallel: parallel}
	return step, nil
}

func parseArgument(line string, lineNo int, rawLine string) (string, string, error) {
	splits := strings.SplitN(line, ":", 2)
	if len(splits) != 2 {
		return "", "", fmt.Errorf("parsing error for at line %d, '%s' malformed argument", lineNo, rawLine)
	}
	return splits[0], splits[1], nil
}

// parseSteps splits the workflow DSL file into steps for processing
func ParseSteps(reader io.Reader) (*Steps, error) {
	steps := &Steps{}
	scanner := bufio.NewScanner(reader)
	lineNo := 0
	for scanner.Scan() {
		lineNo += 1
		line := scanner.Text()
		rawLine := line
		line = cleanLine(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "command:") {
			currentStep, err := parseCommand(line, lineNo, rawLine, false)
			if err != nil {
				return nil, err
			}
			steps.List = append(steps.List, *currentStep)
			continue
		}
		if strings.HasPrefix(line, "parallel:") {
			currentStep, err := parseCommand(line, lineNo, rawLine, true)
			if err != nil {
				return nil, err
			}
			steps.List = append(steps.List, *currentStep)
			continue
		}
		key, value, err := parseArgument(line, lineNo, rawLine)
		if err != nil {
			return nil, err
		}
		steps.List[len(steps.List)-1].Arguments[key] = value

	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	steps.b = make(chan int, len(steps.List))
	steps.e = make(chan string, len(steps.List))
	return steps, nil
}
