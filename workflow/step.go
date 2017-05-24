package workflow

import (
	"fmt"
	"log"
	"sync"
)

type Step struct {
	Command   string
	Arguments map[string]string
	Parallel  bool
	lineNo    int
	err       string
}
type Steps struct {
	List []Step
}

func (s Step) Error() string {
	return fmt.Sprintf("Error running step %s line %d: %s", s.Command, s.lineNo, s.err)
}

func (s Step) Run() (int, error) {
	log.Println("Running", s.Command)
	f, ok := commandList[s.Command]
	if !ok {
		s.err = fmt.Sprintf("unknown command '%s'", s.Command)
		return 0, fmt.Errorf(s.Error())
	}
	return f(s.Arguments)
}

func (steps Steps) waitCount() int {
	count := 0
	for _, s := range steps.List {
		if s.Parallel {
			count += 1
		}
	}
	return count
}

func (steps Steps) Run() (int, error) {
	var wg sync.WaitGroup
	wg.Add(steps.waitCount())
	errs := make(chan string, steps.waitCount())
	totalBytes := 0
	var errorState string
	by := make(chan int, steps.waitCount())
	go func() {
		totalBytes += <-by
	}()
	go func() {
		errorState = <-errs
	}()
	for _, step := range steps.List {
		if errorState != "" {
			return 0, fmt.Errorf(errorState)
		}
		if step.Parallel {
			go func(s Step) {
				defer wg.Done()
				b, err := s.Run()
				if err != nil {
					errs <- err.Error()
				}
				by <- b
			}(step)
		} else {
			b, err := step.Run()
			if err != nil {
				return 0, err
			}
			totalBytes += b
		}
	}
	wg.Wait()
	return totalBytes, nil
}
