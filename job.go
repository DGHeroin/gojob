package jobs

import "log"

type (
    Worker struct {
        WorkerPool chan chan Job
        JobChannel chan Job
        quit       chan bool
    }
    Job interface {
        PreRun() bool
        Run()
    }
    jobWrap struct {
        fn func()
    }
)

func (j jobWrap) PreRun() bool {
    return true
}

func (j jobWrap) Run() {
    j.fn()
}

func NewWorker(workerPool chan chan Job) Worker {
    return Worker{
        WorkerPool: workerPool,
        JobChannel: make(chan Job, 10),
        quit:       make(chan bool)}
}

func (w Worker) Start() {
    go func() {
        for {
            w.WorkerPool <- w.JobChannel
            select {
            case job := <-w.JobChannel:
                if job.PreRun() {
                    job.Run()
                }
            case <-w.quit:
                return
            }
        }
    }()
}

func (w Worker) Stop() {
    go func() {
        w.quit <- true
    }()
}

type Dispatcher struct {
    jobQueue   chan Job
    workerPool chan chan Job
    maxWorkers int
    closeChan  chan struct{}
}

func NewDispatcher(maxWorkers int, jobQueue chan Job) *Dispatcher {
    if jobQueue == nil {
        jobQueue = make(chan Job)
    }
    pool := make(chan chan Job, maxWorkers)
    return &Dispatcher{
        workerPool: pool,
        maxWorkers: maxWorkers,
        jobQueue:   jobQueue,
        closeChan:  make(chan struct{}),
    }
}

func (d *Dispatcher) Start() {
    for i := 0; i < d.maxWorkers; i++ {
        worker := NewWorker(d.workerPool)
        worker.Start()
    }

    go d.dispatch()
}

func (d *Dispatcher) dispatch() {
    for {
        select {
        case job := <-d.jobQueue:
            go func(job Job) {
                defer func() {
                    if e := recover(); e != nil {
                        log.Println(e)
                    }
                }()
                jobChannel := <-d.workerPool
                jobChannel <- job
            }(job)
        case <-d.closeChan:
            return
        }
    }
}

func (d *Dispatcher) Close() {
    close(d.closeChan)
}

func (d *Dispatcher) Run(f func()) {
    d.jobQueue <- &jobWrap{
        fn: f,
    }
}
