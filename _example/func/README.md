

## TryCatch

- 简单使用

```go
package main
import (
    "errors"
    "fmt"

    "github.com/guojia99/octopus/ofunc"
)

func main(){
  ofunc.Try(func() {
    // todo code
    var err error
    ofunc.Throw(err)
  }).Catch(func(err ofunc.StackTraceErr) {
    fmt.Println(err)
  })
}
```

- 多个Try的时候 使用`CatchAll`

```go
package main
import (
    "errors"
    "fmt"

    "github.com/guojia99/octopus/ofunc"
)

func main(){
    ofunc.Try(func() {
        var err error
        ofunc.Throw(err)
    }).Try(func() {
        // todo code
        panic(errors.New("test error by panic"))
    }).Try(func() {
        ofunc.Throw("test error by throw")
        fmt.Println("this code will not be executed ")
    }).Try(func() {
        // todo code
        a := 0
        b := 0
        fmt.Println(a / b)
    }).CatchAll(func(errs []ofunc.StackTraceErr) {
        for i := 0; i < len(errs); i++ {
            fmt.Println(errs[i])
        }
    })
}
```

- 指定某个错误的时候使用`CatchWithErr` 进行捕获

```go
var DataError = errors.New("test data error")
ofunc.Try(func() {
    ofunc.Throw(DataError)
}).Try(func() {
    ofunc.Throw("todo error")
}).CatchWithErr(DataError, func(err ofunc.StackTraceErr) {
    fmt.Println(err)
})
```

