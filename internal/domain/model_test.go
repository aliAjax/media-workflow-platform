package domain

import "testing"

func TestFailedJobCanRequeue(t *testing.T) {
	j := Job{State: JobFailed}
	if err := j.Transition(JobQueued); err != nil {
		t.Fatal(err)
	}
}

func TestRetryStateKeepsSuccess(t *testing.T) {
	if RetryableState(JobSucceeded) != JobSucceeded {
		t.Fatal("succeeded state changed")
	}
}

func TestRunningJobCanComplete(t *testing.T) {
	if !CanComplete(JobRunning) {
		t.Fatal("running job cannot complete")
	}
}

func TestTerminalStatesIncludeSuccess(t *testing.T) {
	if !IsTerminalState(JobSucceeded) {
		t.Fatal("succeeded job is not terminal")
	}
}
