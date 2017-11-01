package vm

import (
	"os"
	"testing"
)

// Failures

func TestRWLockNewMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`
		require 'concurrent/rw_lock'
		Concurrent::RWLock.new(5)
		`, "ArgumentError: Expected 0 arguments, got 1", 3, 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestRWLockWithReadLockMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`
		require 'concurrent/rw_lock'
		Concurrent::RWLock.new.with_read_lock
		`, "InternalError: Can't yield without a block", 3, 1},
		{`
		require 'concurrent/rw_lock'
		Concurrent::RWLock.new.with_read_lock(5) do end
		`, "ArgumentError: Expected 0 arguments, got 1", 3, 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

func TestRWLockWithWriteLockMethodFail(t *testing.T) {
	testsFail := []errorTestCase{
		{`
		require 'concurrent/rw_lock'
		Concurrent::RWLock.new.with_write_lock
		`, "InternalError: Can't yield without a block", 3, 1},
		{`
		require 'concurrent/rw_lock'
		Concurrent::RWLock.new.with_write_lock(5) do end
		`, "ArgumentError: Expected 0 arguments, got 1", 3, 1},
	}

	for i, tt := range testsFail {
		v := initTestVM()
		evaluated := v.testEval(t, tt.input, getFilename())
		checkError(t, i, evaluated, tt.expected, getFilename(), tt.errorLine)
		v.checkCFP(t, i, tt.expectedCFP)
		v.checkSP(t, i, 1)
	}
}

// Isolated lock types

func TestRWLockWithReadLockMethod(t *testing.T) {
	code := `
	require 'concurrent/rw_lock'

	lock = Concurrent::RWLock.new
	message = nil

	lock.with_read_lock do
		message = "completed"
	end

	message
	`

	expected := "completed"

	v := initTestVM()
	evaluated := v.testEval(t, code, getFilename())
	testStringObject(t, i, evaluated, expected)
	v.checkCFP(t, i, 0)
	v.checkSP(t, i, 1)
}

func TestRWLockWithWriteLockMethod(t *testing.T) {
	code := `
	require 'concurrent/rw_lock'

	lock = Concurrent::RWLock.new
	message = nil

	lock.with_write_lock do
		message = "completed"
	end

	message
	`

	expected := "completed"

	v := initTestVM()
	evaluated := v.testEval(t, code, getFilename())
	testStringObject(t, i, evaluated, expected)
	v.checkCFP(t, i, 0)
	v.checkSP(t, i, 1)
}

// Mixed locks (functional tests)

func TestRWLockWithReadLockReadBlocksWriteNoRaceDetection(t *testing.T) {
	skipRWLockTestIfRaceDetectionEnabled(t)

	code := `
	require 'concurrent/rw_lock'

	lock = Concurrent::RWLock.new
	message = nil

	thread do
	  lock.with_read_lock do
	    sleep 2
	    message ||= "thread 1"
	  end
	end

	thread do
	  sleep 1
	  lock.with_write_lock do
	    message ||= "thread 2"
	  end
	end

	sleep 3
	lock.with_read_lock do
		message
	end
	`

	expected := "thread 1"

	v := initTestVM()
	evaluated := v.testEval(t, code, getFilename())
	testStringObject(t, i, evaluated, expected)
	v.checkCFP(t, i, 0)
	v.checkSP(t, i, 1)
}

func skipRWLockTestIfRaceDetectionEnabled(t *testing.T) {
	if os.Getenv("NO_RACE_DETECTION") == "" {
		t.Skip("skipping RW lock related tests")
	}
}
