package httpcontroller

// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

//go:generate minimock -i github.com/hablof/order-viewer/internal/httpcontroller.TemplateExecutor -o ./internal\httpcontroller\template_executor_mock_test.go -n TemplateExecutorMock

import (
	"io"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
)

// TemplateExecutorMock implements TemplateExecutor
type TemplateExecutorMock struct {
	t minimock.Tester

	funcExecuteTemplate          func(wr io.Writer, name string, data interface{}) (err error)
	inspectFuncExecuteTemplate   func(wr io.Writer, name string, data interface{})
	afterExecuteTemplateCounter  uint64
	beforeExecuteTemplateCounter uint64
	ExecuteTemplateMock          mTemplateExecutorMockExecuteTemplate
}

// NewTemplateExecutorMock returns a mock for TemplateExecutor
func NewTemplateExecutorMock(t minimock.Tester) *TemplateExecutorMock {
	m := &TemplateExecutorMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ExecuteTemplateMock = mTemplateExecutorMockExecuteTemplate{mock: m}
	m.ExecuteTemplateMock.callArgs = []*TemplateExecutorMockExecuteTemplateParams{}

	return m
}

type mTemplateExecutorMockExecuteTemplate struct {
	mock               *TemplateExecutorMock
	defaultExpectation *TemplateExecutorMockExecuteTemplateExpectation
	expectations       []*TemplateExecutorMockExecuteTemplateExpectation

	callArgs []*TemplateExecutorMockExecuteTemplateParams
	mutex    sync.RWMutex
}

// TemplateExecutorMockExecuteTemplateExpectation specifies expectation struct of the TemplateExecutor.ExecuteTemplate
type TemplateExecutorMockExecuteTemplateExpectation struct {
	mock    *TemplateExecutorMock
	params  *TemplateExecutorMockExecuteTemplateParams
	results *TemplateExecutorMockExecuteTemplateResults
	Counter uint64
}

// TemplateExecutorMockExecuteTemplateParams contains parameters of the TemplateExecutor.ExecuteTemplate
type TemplateExecutorMockExecuteTemplateParams struct {
	wr   io.Writer
	name string
	data interface{}
}

// TemplateExecutorMockExecuteTemplateResults contains results of the TemplateExecutor.ExecuteTemplate
type TemplateExecutorMockExecuteTemplateResults struct {
	err error
}

// Expect sets up expected params for TemplateExecutor.ExecuteTemplate
func (mmExecuteTemplate *mTemplateExecutorMockExecuteTemplate) Expect(wr io.Writer, name string, data interface{}) *mTemplateExecutorMockExecuteTemplate {
	if mmExecuteTemplate.mock.funcExecuteTemplate != nil {
		mmExecuteTemplate.mock.t.Fatalf("TemplateExecutorMock.ExecuteTemplate mock is already set by Set")
	}

	if mmExecuteTemplate.defaultExpectation == nil {
		mmExecuteTemplate.defaultExpectation = &TemplateExecutorMockExecuteTemplateExpectation{}
	}

	mmExecuteTemplate.defaultExpectation.params = &TemplateExecutorMockExecuteTemplateParams{wr, name, data}
	for _, e := range mmExecuteTemplate.expectations {
		if minimock.Equal(e.params, mmExecuteTemplate.defaultExpectation.params) {
			mmExecuteTemplate.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmExecuteTemplate.defaultExpectation.params)
		}
	}

	return mmExecuteTemplate
}

// Inspect accepts an inspector function that has same arguments as the TemplateExecutor.ExecuteTemplate
func (mmExecuteTemplate *mTemplateExecutorMockExecuteTemplate) Inspect(f func(wr io.Writer, name string, data interface{})) *mTemplateExecutorMockExecuteTemplate {
	if mmExecuteTemplate.mock.inspectFuncExecuteTemplate != nil {
		mmExecuteTemplate.mock.t.Fatalf("Inspect function is already set for TemplateExecutorMock.ExecuteTemplate")
	}

	mmExecuteTemplate.mock.inspectFuncExecuteTemplate = f

	return mmExecuteTemplate
}

// Return sets up results that will be returned by TemplateExecutor.ExecuteTemplate
func (mmExecuteTemplate *mTemplateExecutorMockExecuteTemplate) Return(err error) *TemplateExecutorMock {
	if mmExecuteTemplate.mock.funcExecuteTemplate != nil {
		mmExecuteTemplate.mock.t.Fatalf("TemplateExecutorMock.ExecuteTemplate mock is already set by Set")
	}

	if mmExecuteTemplate.defaultExpectation == nil {
		mmExecuteTemplate.defaultExpectation = &TemplateExecutorMockExecuteTemplateExpectation{mock: mmExecuteTemplate.mock}
	}
	mmExecuteTemplate.defaultExpectation.results = &TemplateExecutorMockExecuteTemplateResults{err}
	return mmExecuteTemplate.mock
}

// Set uses given function f to mock the TemplateExecutor.ExecuteTemplate method
func (mmExecuteTemplate *mTemplateExecutorMockExecuteTemplate) Set(f func(wr io.Writer, name string, data interface{}) (err error)) *TemplateExecutorMock {
	if mmExecuteTemplate.defaultExpectation != nil {
		mmExecuteTemplate.mock.t.Fatalf("Default expectation is already set for the TemplateExecutor.ExecuteTemplate method")
	}

	if len(mmExecuteTemplate.expectations) > 0 {
		mmExecuteTemplate.mock.t.Fatalf("Some expectations are already set for the TemplateExecutor.ExecuteTemplate method")
	}

	mmExecuteTemplate.mock.funcExecuteTemplate = f
	return mmExecuteTemplate.mock
}

// When sets expectation for the TemplateExecutor.ExecuteTemplate which will trigger the result defined by the following
// Then helper
func (mmExecuteTemplate *mTemplateExecutorMockExecuteTemplate) When(wr io.Writer, name string, data interface{}) *TemplateExecutorMockExecuteTemplateExpectation {
	if mmExecuteTemplate.mock.funcExecuteTemplate != nil {
		mmExecuteTemplate.mock.t.Fatalf("TemplateExecutorMock.ExecuteTemplate mock is already set by Set")
	}

	expectation := &TemplateExecutorMockExecuteTemplateExpectation{
		mock:   mmExecuteTemplate.mock,
		params: &TemplateExecutorMockExecuteTemplateParams{wr, name, data},
	}
	mmExecuteTemplate.expectations = append(mmExecuteTemplate.expectations, expectation)
	return expectation
}

// Then sets up TemplateExecutor.ExecuteTemplate return parameters for the expectation previously defined by the When method
func (e *TemplateExecutorMockExecuteTemplateExpectation) Then(err error) *TemplateExecutorMock {
	e.results = &TemplateExecutorMockExecuteTemplateResults{err}
	return e.mock
}

// ExecuteTemplate implements TemplateExecutor
func (mmExecuteTemplate *TemplateExecutorMock) ExecuteTemplate(wr io.Writer, name string, data interface{}) (err error) {
	mm_atomic.AddUint64(&mmExecuteTemplate.beforeExecuteTemplateCounter, 1)
	defer mm_atomic.AddUint64(&mmExecuteTemplate.afterExecuteTemplateCounter, 1)

	if mmExecuteTemplate.inspectFuncExecuteTemplate != nil {
		mmExecuteTemplate.inspectFuncExecuteTemplate(wr, name, data)
	}

	mm_params := &TemplateExecutorMockExecuteTemplateParams{wr, name, data}

	// Record call args
	mmExecuteTemplate.ExecuteTemplateMock.mutex.Lock()
	mmExecuteTemplate.ExecuteTemplateMock.callArgs = append(mmExecuteTemplate.ExecuteTemplateMock.callArgs, mm_params)
	mmExecuteTemplate.ExecuteTemplateMock.mutex.Unlock()

	for _, e := range mmExecuteTemplate.ExecuteTemplateMock.expectations {
		if minimock.Equal(e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.err
		}
	}

	if mmExecuteTemplate.ExecuteTemplateMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmExecuteTemplate.ExecuteTemplateMock.defaultExpectation.Counter, 1)
		mm_want := mmExecuteTemplate.ExecuteTemplateMock.defaultExpectation.params
		mm_got := TemplateExecutorMockExecuteTemplateParams{wr, name, data}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmExecuteTemplate.t.Errorf("TemplateExecutorMock.ExecuteTemplate got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmExecuteTemplate.ExecuteTemplateMock.defaultExpectation.results
		if mm_results == nil {
			mmExecuteTemplate.t.Fatal("No results are set for the TemplateExecutorMock.ExecuteTemplate")
		}
		return (*mm_results).err
	}
	if mmExecuteTemplate.funcExecuteTemplate != nil {
		return mmExecuteTemplate.funcExecuteTemplate(wr, name, data)
	}
	mmExecuteTemplate.t.Fatalf("Unexpected call to TemplateExecutorMock.ExecuteTemplate. %v %v %v", wr, name, data)
	return
}

// ExecuteTemplateAfterCounter returns a count of finished TemplateExecutorMock.ExecuteTemplate invocations
func (mmExecuteTemplate *TemplateExecutorMock) ExecuteTemplateAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmExecuteTemplate.afterExecuteTemplateCounter)
}

// ExecuteTemplateBeforeCounter returns a count of TemplateExecutorMock.ExecuteTemplate invocations
func (mmExecuteTemplate *TemplateExecutorMock) ExecuteTemplateBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmExecuteTemplate.beforeExecuteTemplateCounter)
}

// Calls returns a list of arguments used in each call to TemplateExecutorMock.ExecuteTemplate.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmExecuteTemplate *mTemplateExecutorMockExecuteTemplate) Calls() []*TemplateExecutorMockExecuteTemplateParams {
	mmExecuteTemplate.mutex.RLock()

	argCopy := make([]*TemplateExecutorMockExecuteTemplateParams, len(mmExecuteTemplate.callArgs))
	copy(argCopy, mmExecuteTemplate.callArgs)

	mmExecuteTemplate.mutex.RUnlock()

	return argCopy
}

// MinimockExecuteTemplateDone returns true if the count of the ExecuteTemplate invocations corresponds
// the number of defined expectations
func (m *TemplateExecutorMock) MinimockExecuteTemplateDone() bool {
	for _, e := range m.ExecuteTemplateMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.ExecuteTemplateMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterExecuteTemplateCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcExecuteTemplate != nil && mm_atomic.LoadUint64(&m.afterExecuteTemplateCounter) < 1 {
		return false
	}
	return true
}

// MinimockExecuteTemplateInspect logs each unmet expectation
func (m *TemplateExecutorMock) MinimockExecuteTemplateInspect() {
	for _, e := range m.ExecuteTemplateMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to TemplateExecutorMock.ExecuteTemplate with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.ExecuteTemplateMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterExecuteTemplateCounter) < 1 {
		if m.ExecuteTemplateMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to TemplateExecutorMock.ExecuteTemplate")
		} else {
			m.t.Errorf("Expected call to TemplateExecutorMock.ExecuteTemplate with params: %#v", *m.ExecuteTemplateMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcExecuteTemplate != nil && mm_atomic.LoadUint64(&m.afterExecuteTemplateCounter) < 1 {
		m.t.Error("Expected call to TemplateExecutorMock.ExecuteTemplate")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *TemplateExecutorMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockExecuteTemplateInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *TemplateExecutorMock) MinimockWait(timeout mm_time.Duration) {
	timeoutCh := mm_time.After(timeout)
	for {
		if m.minimockDone() {
			return
		}
		select {
		case <-timeoutCh:
			m.MinimockFinish()
			return
		case <-mm_time.After(10 * mm_time.Millisecond):
		}
	}
}

func (m *TemplateExecutorMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockExecuteTemplateDone()
}
