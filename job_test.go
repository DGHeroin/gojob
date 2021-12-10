package jobs

import (
    "fmt"
    "sync/atomic"
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
    dispatcher.Run(func() {
        t.Log("testing 5")
    })
    time.Sleep(time.Second)
}

func TestQPS(t *testing.T) {
    var (
        count     int64
        last      int64
        isRunning = true
    )
    dispatcher := NewDispatcher(50, nil)
    dispatcher.Start()

    go func() {
        n := 0
        for n < 10 {
            time.Sleep(time.Second)
            now := atomic.LoadInt64(&count)
            fmt.Println(n, " qps:", now-last)
            last = now
            n++
        }
        isRunning = false
    }()
    fn := func() {
        atomic.AddInt64(&count, 1)
    }
    for isRunning {
        dispatcher.Run(fn)
    }

}
