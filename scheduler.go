/*
------------------------------------------------------------------------------------------------------------------------
####### scheduler ####### (c) 2020-2021 mls-361 #################################################### MIT License #######
------------------------------------------------------------------------------------------------------------------------
*/

package scheduler

import (
	"fmt"
	"time"

	"github.com/mls-361/failure"
	"github.com/robfig/cron/v3"
)

type (
	emitCallback func(name string, data interface{})

	tools struct {
		ecb  emitCallback
		cron *cron.Cron
	}

	// Event AFAIRE.
	Event struct {
		Name     string
		Disabled bool
		After    time.Duration
		Repeat   string
		Data     interface{}
	}

	// Scheduler AFAIRE.
	Scheduler struct {
		parser cron.Parser
		tools  *tools
	}
)

// New AFAIRE.
func New(ecb emitCallback) *Scheduler {
	parser := cron.NewParser(
		cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
	)

	tools := &tools{
		ecb:  ecb,
		cron: cron.New(cron.WithParser(parser)),
	}

	return &Scheduler{
		parser: parser,
		tools:  tools,
	}
}

// AddEvent AFAIRE.
func (s *Scheduler) AddEvent(e *Event) error {
	if e.Name == "" {
		return failure.New(nil).Msg("an event name cannot be empty") ///////////////////////////////////////////////////
	}

	// Pour le moment, on les élimine.
	// Plus tard, on pourra les prendre en compte pour les activer en cours d'exécution.
	if e.Disabled {
		return nil
	}

	event := &event{
		name:  e.Name,
		data:  e.Data,
		tools: s.tools,
	}

	if e.After != 0 {
		entryID, err := s.tools.cron.AddJob(fmt.Sprintf("@every %s", e.After.String()), event)
		if err != nil {
			return err
		}

		event.after = entryID
	}

	if e.Repeat == "" {
		if e.After == 0 {
			return failure.New(nil).
				Msg("at least one of the fields After or Repeat must be initialized") //////////////////////////////////
		}

		return nil
	}

	schedule, err := s.parser.Parse(e.Repeat)
	if err != nil {
		return err
	}

	if e.After == 0 {
		_ = s.tools.cron.Schedule(schedule, event)
	} else {
		event.repeat = schedule
	}

	return nil
}

// Start AFAIRE.
func (s *Scheduler) Start() {
	s.tools.cron.Start()
}

// Stop AFAIRE.
func (s *Scheduler) Stop() {
	s.tools.cron.Stop()
}

/*
######################################################################################################## @(°_°)@ #######
*/
