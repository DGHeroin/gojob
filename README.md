#GoJob

```go
import "github.com/DGHeroin/jobs"
func Foo() {
    q := make(chan Job)
    dispatcher := NewDispatcher(10, q)
    dispatcher.Start()

    q <- &testJob{args: 1, t: t}
}
```

```go
import "github.com/DGHeroin/jobs"
func Bar() {
    dispatcher := NewDispatcher(10, nil)
    dispatcher.Start()

    dispatcher.Run(func() {
        // your code
    })
}
```