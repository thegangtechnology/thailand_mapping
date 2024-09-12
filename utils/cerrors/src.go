package cerrors

import (
	"go-template/database"
	"sync"
)

type RoutineManager interface {
	FanInErrors() RoutineError
	GetProcesses() []RoutineError
}

type RoutineRunner interface {
	Execute(functionToRun func() error)
	WaitForError() error
	WaitForErrors() []error
}

type RoutineError struct {
	Error chan error
	Done  chan bool
}

type RoutineManagerImpl struct {
	processes []RoutineError
}

type RoutineRunnerImpl struct {
	wg      *sync.WaitGroup
	ConnCh  chan bool
	ConnErr chan error
	Done    chan bool
}

func NewRoutineManager(size int) RoutineManager {
	routines := make([]RoutineError, size)
	for i := 0; i < len(routines); i++ {
		routines[i].Error = make(chan error)
		routines[i].Done = make(chan bool)
	}

	return RoutineManagerImpl{processes: routines}
}

func NewRoutineRunner(size int) RoutineRunner {
	connCh := make(chan bool, database.MaxIdleConns)
	connErr := make(chan error)
	done := make(chan bool)
	wg := sync.WaitGroup{}

	wg.Add(size)

	return RoutineRunnerImpl{wg: &wg, ConnCh: connCh, ConnErr: connErr, Done: done}
}

func (g RoutineRunnerImpl) Execute(functionToRun func() error) {
	run := func() {
		defer g.wg.Done()

		g.ConnCh <- true

		err := functionToRun()

		if err != nil {
			g.ConnErr <- err
		}

		<-g.ConnCh
	}

	catch := func() {
		g.wg.Wait()

		g.Done <- true
	}

	go run()
	go catch()
}

func (g RoutineRunnerImpl) WaitForError() error {
	select {
	case err := <-g.ConnErr:
		if err != nil {
			return err
		}
	case <-g.Done:
	}

	return nil
}

func (g RoutineRunnerImpl) WaitForErrors() []error {
	var errs []error

	out := false

	for !out {
		select {
		case err := <-g.ConnErr:
			if err != nil {
				errs = append(errs, err)
			}
		case out = <-g.Done:
		}
	}

	return errs
}

func (r RoutineManagerImpl) GetProcesses() []RoutineError {
	return r.processes
}

func (r RoutineManagerImpl) FanInErrors() RoutineError {
	var wg sync.WaitGroup

	check := RoutineError{make(chan error), make(chan bool)}

	// Start an output goroutine for each input channel in errs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	wg.Add(len(r.processes))

	for _, process := range r.processes {
		go r.fan(&wg, &check, process.Error, process.Done)
	}

	// Start a goroutine to close out once all the output goroutines are done.
	// This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(check.Error)

		check.Done <- true
		close(check.Done)
	}()

	return check
}

func (r RoutineManagerImpl) fan(wg *sync.WaitGroup, check *RoutineError, c <-chan error, done <-chan bool) {
	select {
	case b := <-c:
		check.Error <- b
		r.fan(wg, check, c, done)
	case <-done:
		wg.Done()

		return
	}
}
