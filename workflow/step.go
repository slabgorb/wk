package workflow

import (
	"fmt"
	"log"
)

type Step struct {
	Command   string
	Arguments map[string]string
	Parallel  bool
	LineNo    int
	err       string
}
type Steps struct {
	List []Step
	b    chan int
	e    chan string
}

func (s *Step) SetError(errMessage string) {
	s.err = errMessage
}

func (s *Step) Error() string {
	return fmt.Sprintf("error running step %s line %d: %s", s.Command, s.LineNo, s.err)
}

func (s *Step) Run() (int, error) {
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
	totalBytes := 0
	done := make(chan struct{})
	for _, step := range steps.List {
		if step.Parallel {
			go func(s Step) {
				b, err := s.Run()
				if err != nil {
					steps.e <- err.Error()
				}
				steps.b <- b
			}(step)
		} else {
			b, err := step.Run()
			if err != nil {
				steps.e <- err.Error()
			}
			steps.b <- b
		}
	}
	var errMessage string
	for i := 0; i < len(steps.List); i++ {
		select {
		case b := <-steps.b:
			totalBytes += b
		case errMessage = <-steps.e:
			return 0, fmt.Errorf(errMessage)
		case <-done:
			return totalBytes, nil
		}
	}
	return totalBytes, nil
}
