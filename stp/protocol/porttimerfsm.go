// 17.22 Port Timers state machine
package stp

import (
	//"fmt"
	"time"
	"utils/fsm"
)

const PtmMachineModuleStr = "Port Timer State Machine"

const (
	PtmStateNone = iota + 1
	PtmStateOneSecond
	PtmStateTick
)

var PtmStateStrMap map[fsm.State]string

func PtmMachineStrStateMapInit() {
	PtmStateStrMap = make(map[fsm.State]string)
	PtmStateStrMap[PtmStateNone] = "None"
	PtmStateStrMap[PtmStateOneSecond] = "OneSecond"
	PtmStateStrMap[PtmStateTick] = "Tick"
}

const (
	PtmEventBegin = iota + 1
	PtmEventTickEqualsTrue
	PtmEventUnconditionalFallthrough
)

// LacpRxMachine holds FSM and current State
// and event channels for State transitions
type PtmMachine struct {
	// for debugging
	PreviousState fsm.State

	Machine *fsm.Machine

	// State transition log
	log chan string

	// timer type
	TickTimer *time.Timer
	Tick      bool

	// Reference to StpPort
	p *StpPort

	// machine specific events
	PtmEvents chan MachineEvent
	// stop go routine
	PtmKillSignalEvent chan bool
	// enable logging
	PtmLogEnableEvent chan bool
}

func (ptm *PtmMachine) PrevState() fsm.State { return ptm.PreviousState }

// PrevStateSet will set the previous State
func (ptm *PtmMachine) PrevStateSet(s fsm.State) { ptm.PreviousState = s }

// NewLacpRxMachine will create a new instance of the LacpRxMachine
func NewStpPtmMachine(p *StpPort) *PtmMachine {
	ptm := &PtmMachine{
		p:                  p,
		PreviousState:      PtmStateNone,
		PtmEvents:          make(chan MachineEvent, 10),
		PtmKillSignalEvent: make(chan bool),
		PtmLogEnableEvent:  make(chan bool)}

	// start then stop
	ptm.TickTimerStart()
	ptm.TickTimerStop()

	p.PtmMachineFsm = ptm

	return ptm
}

// A helpful function that lets us apply arbitrary rulesets to this
// instances State machine without reallocating the machine.
func (ptm *PtmMachine) Apply(r *fsm.Ruleset) *fsm.Machine {
	if ptm.Machine == nil {
		ptm.Machine = &fsm.Machine{}
	}

	// Assign the ruleset to be used for this machine
	ptm.Machine.Rules = r
	ptm.Machine.Curr = &StpStateEvent{
		strStateMap: PtmStateStrMap,
		//logEna:      ptxm.p.logEna,
		logEna: false,
		logger: StpLoggerInfo,
		owner:  PtmMachineModuleStr,
	}

	return ptm.Machine
}

// Stop should clean up all resources
func (ptm *PtmMachine) Stop() {

	ptm.TickTimerDestroy()

	// stop the go routine
	ptm.PtmKillSignalEvent <- true

	close(ptm.PtmEvents)
	close(ptm.PtmLogEnableEvent)
	close(ptm.PtmKillSignalEvent)

}

// LacpPtxMachineNoPeriodic stops the periodic transmission of packets
func (ptm *PtmMachine) PtmMachineOneSecond(m fsm.Machine, data interface{}) fsm.State {
	ptm.Tick = false
	return PtmStateOneSecond
}

// LacpPtxMachineFastPeriodic sets the periodic transmission time to fast
// and starts the timer
func (ptm *PtmMachine) PtmMachineTick(m fsm.Machine, data interface{}) fsm.State {
	p := ptm.p
	p.DecrementTimerCounters()

	return PtmStateTick
}

func PtmMachineFSMBuild(p *StpPort) *PtmMachine {

	rules := fsm.Ruleset{}

	// Instantiate a new LacpPtxMachine
	// Initial State will be a psuedo State known as "begin" so that
	// we can transition to the NO PERIODIC State
	ptm := NewStpPtmMachine(p)

	//BEGIN -> ONE SECOND
	rules.AddRule(PtmStateNone, PtmEventBegin, ptm.PtmMachineOneSecond)
	rules.AddRule(PtmStateOneSecond, PtmEventBegin, ptm.PtmMachineOneSecond)
	rules.AddRule(PtmStateTick, PtmEventBegin, ptm.PtmMachineOneSecond)

	// TICK EQUALS TRUE	 -> TICK
	rules.AddRule(PtmStateOneSecond, PtmEventTickEqualsTrue, ptm.PtmMachineTick)

	// PORT DISABLED -> NO PERIODIC
	rules.AddRule(PtmStateTick, PtmEventUnconditionalFallthrough, ptm.PtmMachineOneSecond)

	// Create a new FSM and apply the rules
	ptm.Apply(&rules)

	return ptm
}

// LacpRxMachineMain:  802.1ax-2014 Table 6-18
// Creation of Rx State Machine State transitions and callbacks
// and create go routine to pend on events
func (p *StpPort) PtmMachineMain() {

	// Build the State machine for Lacp Receive Machine according to
	// 802.1ax Section 6.4.13 Periodic Transmission Machine
	ptm := PtmMachineFSMBuild(p)
	p.wg.Add(1)

	// set the inital State
	ptm.Machine.Start(ptm.PrevState())

	// lets create a go routing which will wait for the specific events
	// that the Port Timer State Machine should handle
	go func(m *PtmMachine) {
		StpMachineLogger("INFO", "PTM", "Machine Start")
		defer m.p.wg.Done()
		for {
			select {
			case <-m.PtmKillSignalEvent:
				StpMachineLogger("INFO", "PTM", "Machine End")
				return

			case <-m.TickTimer.C:

				m.Machine.ProcessEvent(PtmMachineModuleStr, PtmEventTickEqualsTrue, nil)

			case event := <-m.PtmEvents:
				m.Machine.ProcessEvent(event.src, event.e, nil)

				if event.responseChan != nil {
					SendResponse(PtmMachineModuleStr, event.responseChan)
				}

			case ena := <-m.PtmLogEnableEvent:
				m.Machine.Curr.EnableLogging(ena)
			}
		}
	}(ptm)
}
