package gospider

import (
	"fmt"
	"sync/atomic"
)

/**************************************************************
* type:	MID : module ID
**************************************************************/
// MID stands for module ID
type MID string

// variable for generating MID
var (
	midTemplate = "%s%d"
	midSerial   = NewSerial()
)

// GenMID will generate MID according to module't type
func GenMID(mtype ModuleType) (mid MID) {
	return MID(fmt.Sprintf(midTemplate, ModuleType2Letter[mtype], midSerial.Next()))
}

/**************************************************************
* type:	ModuleType : module type
**************************************************************/
// ModuleType type
type ModuleType string

// Allowed module type
const (
	ModuleTypeGenerator  ModuleType = "generator"
	ModuleTypeDownloader ModuleType = "downloader"
	ModuleTypeAnalyzer   ModuleType = "analyzer"
	ModuleTypePipeline   ModuleType = "pipeline"
)

// ModuleType2Letter contains mapping from Type string to Heading letter
var ModuleType2Letter = map[ModuleType]string{
	ModuleTypeGenerator:  "G",
	ModuleTypeDownloader: "D",
	ModuleTypeAnalyzer:   "A",
	ModuleTypePipeline:   "P",
}

// Letter2ModuleType contains mapping from Heading letter to type string
var Letter2ModuleType = map[string]ModuleType{
	"G": ModuleTypeGenerator,
	"D": ModuleTypeDownloader,
	"A": ModuleTypeAnalyzer,
	"P": ModuleTypePipeline,
}

// CheckModuleType will judge whether given module fits corresponding type
func CheckModuleType(moduleModuleType ModuleType, module Module) bool {
	if moduleModuleType == "" || module == nil {
		return false
	}
	switch moduleModuleType {
	case ModuleTypeGenerator:
		if _, ok := module.(Generator); ok {
			return true
		}
	case ModuleTypeDownloader:
		if _, ok := module.(Downloader); ok {
			return true
		}
	case ModuleTypeAnalyzer:
		if _, ok := module.(Analyzer); ok {
			return true
		}
	case ModuleTypePipeline:
		if _, ok := module.(Pipeline); ok {
			return true
		}
	}
	return false
}

/**************************************************************
* type:	Counters
**************************************************************/

// Counters holds module counter's values.
// Counters have following constrains: Total = Doing + Done + Fail
// * Valid + Invalid = Invoke
type Counters struct {
	CallCount  uint64
	DoingCount uint64
	DoneCount  uint64
	FailCount  uint64
}

/**************************************************************
* type:	Summary
**************************************************************/
// Summary holds module's summary info
type Summary struct {
	ID         MID         `json:"id"`
	CallCount  uint64      `json:"called_count"`
	DoingCount uint64      `json:"doing_count"`
	DoneCount  uint64      `json:"done_count"`
	FailCount  uint64      `json:"fail_count"`
	Other      interface{} `json:"other,omitempty"`
}

/**************************************************************
* interface: Module
**************************************************************/
// Module is interface for Downloader, Generator, Pipeline, Analyzer,etc...
// Module's implementation must be thread-safe
type Module interface {
	// ID returns module's ID
	ID() MID
	// DoCount indicates total times being invoked
	CallCount() uint64
	// DoingCount indicates working job numbers
	DoingCount() uint64
	// DoneCount indicates done job numbers
	DoneCount() uint64
	// FailCount indicates failed job numbers
	FailCount() uint64
	// Counters returns all counters simultaneously
	Counters() Counters
	// Summary returns module's summary report
	Summary() Summary
}

/**************************************************************
* interface: ModuleInternal
**************************************************************/

// ModuleInternal : basic interface for implementing a module
type ModuleInternal interface {
	Module
	// Count means received a job. callCounter++
	Call()
	// Doing means put hands on this job. doingCounter++
	Doing()
	// Done means job have successfully done. doingCounter-- & doneCounter++
	Done()
	// Fail means job have failed. doingCounter-- & failCounter++
	Fail()
	// Clear resets all counters
	Clear()
}

/**************************************************************
* struct: defaultModule
**************************************************************/

// defaultModule implements Module & ModuleInternal simultaneously
type defaultModule struct {
	mid        MID
	callCount  uint64
	doingCount uint64
	doneCount  uint64
	failCount  uint64
}

// Module methods implementation

// defaultModule_ID
func (m *defaultModule) ID() MID {
	return m.mid
}

// defaultModule_CallCount
func (m *defaultModule) CallCount() uint64 {
	return atomic.LoadUint64(&m.callCount)
}

// defaultModule_DoingCount
func (m *defaultModule) DoingCount() uint64 {
	return atomic.LoadUint64(&m.doingCount)
}

// defaultModule_DoneCount
func (m *defaultModule) DoneCount() uint64 {
	count := atomic.LoadUint64(&m.doneCount)
	return count
}

// defaultModule_FailCount
func (m *defaultModule) FailCount() uint64 {
	return atomic.LoadUint64(&m.failCount)
}

// defaultModule_Counters
func (m *defaultModule) Counters() Counters {
	return Counters{
		CallCount:  atomic.LoadUint64(&m.callCount),
		DoingCount: atomic.LoadUint64(&m.doingCount),
		DoneCount:  atomic.LoadUint64(&m.doneCount),
		FailCount:  atomic.LoadUint64(&m.failCount),
	}
}

// defaultModule_Summary
func (m *defaultModule) Summary() Summary {
	return Summary{
		ID:         m.ID(),
		CallCount:  atomic.LoadUint64(&m.callCount),
		DoingCount: atomic.LoadUint64(&m.doingCount),
		DoneCount:  atomic.LoadUint64(&m.doneCount),
		FailCount:  atomic.LoadUint64(&m.failCount),
		Other:      nil,
	}
}

// ModuleInternal methods implementation

// defaultModule_Call
func (m *defaultModule) Call() {
	atomic.AddUint64(&m.doneCount, 1)
}

// defaultModule_Doing
func (m *defaultModule) Doing() {
	atomic.AddUint64(&m.doingCount, 1)
}

// defaultModule_Done
func (m *defaultModule) Done() {
	atomic.AddUint64(&m.doingCount, ^uint64(0))
	atomic.AddUint64(&m.doneCount, 1)
}

// defaultModule_Fail
func (m *defaultModule) Fail() {
	atomic.AddUint64(&m.doingCount, ^uint64(0))
	atomic.AddUint64(&m.failCount, 1)
}

// defaultModule_Clear
func (m *defaultModule) Clear() {
	atomic.StoreUint64(&m.callCount, 0)
	atomic.StoreUint64(&m.doingCount, 0)
	atomic.StoreUint64(&m.doneCount, 0)
	atomic.StoreUint64(&m.failCount, 0)
}

// NewModuleInternal will create a new module internal instance
func NewModuleInternal(mid MID) (ModuleInternal) {
	return &defaultModule{mid: mid}
}

// NewModuleInternalFromType
func NewModuleInternalFromType(mtype ModuleType) (ModuleInternal) {
	return &defaultModule{mid: GenMID(mtype)}
}
