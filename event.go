/*
------------------------------------------------------------------------------------------------------------------------
####### scheduler ####### (c) 2020-2021 mls-361 #################################################### MIT License #######
------------------------------------------------------------------------------------------------------------------------
*/

package scheduler

import "github.com/robfig/cron/v3"

type (
	event struct {
		name   string
		after  cron.EntryID
		repeat cron.Schedule
		data   interface{}
		tools  *tools
	}
)

func (e *event) Run() {
	e.tools.ecb(e.name, e.data)

	if e.after != 0 {
		e.tools.cron.Remove(e.after)
		e.after = 0
	}

	if e.repeat != nil {
		e.tools.cron.Schedule(e.repeat, e)
		e.repeat = nil
	}
}

/*
######################################################################################################## @(°_°)@ #######
*/
