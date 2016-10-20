

# go_threadpool
##基于go lang的线程池库，使用了库 "github.com/vived/go_queue"
###使用方法:
引入包:<br/>
```
import (
    "fmt"
	pool "github.com/vived/go_threadpool"
)
```

书写处理函数<br/>

```
func myfun(data interface{}) {
    fmt.Println("处理:" + data.(string))	
}
```

使用:<br/>
```
p := pool.NewPool(3, myfun)
p.Start()
```

其它线程中使用

```
p.PushData("abc")
```
结束时使用:

```
p.Stop() //阻塞等待
```


