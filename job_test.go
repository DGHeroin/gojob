package jobs

import (
    "testing"
    "time"
)

type (
    testJob struct {
        args int
        t    *testing.T
    }
)

func (job testJob) PreRun() bool {
    if job.args == 1 {
        return false
    }
    return true
}

func (job testJob) Run() {
    job.t.Log("run job args:", job.args)
}

func TestNewDispatcher(t *testing.T) {
    q := make(chan Job)
    dispatcher := NewDispatcher(10, q)
    dispatcher.Start()

    q <- &testJob{args: 1, t: t}
    q <- &testJob{args: 2, t: t}
    q <- &testJob{args: 3, t: t}

    dispatcher.Run(func() {
        t.Log("testing 4")
    })
    time.Sleep(time.Second)
}
