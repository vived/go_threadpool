package go_threadpool

import (
	"fmt"

	queue "github.com/vived/go_queue"
)

type Process func(data interface{})

type ThreadPool struct {
	count       uint16       //线程数
	threadqueue *queue.Queue //数据缓冲
	process     Process
	startedch   []chan bool //线程已启动
	exitCh      []chan bool //退出线程
	exitedCh    []chan bool //已退出线程
}

func handleFun(stated chan bool, exit chan bool,
	exited chan bool, process Process, data chan interface{}, id uint16) {

	stated <- true
	defer func() {
		exited <- true
	}()

	stop := false
	for !stop {
		select {
		case v := <-data:
			process(v)
		case <-exit:
			stop = true
			close(exit)
			break
		}
		if stop {
			break
		}
	}
}

func NewPool(count uint16, process Process) *ThreadPool {
	q := queue.NewFifoQueue(10)
	return &ThreadPool{count, q, process, make([]chan bool, count),
		make([]chan bool, count), make([]chan bool, count)}
}

func (pool *ThreadPool) PushData(data interface{}) {
	pool.threadqueue.Put(data)
}

func (pool *ThreadPool) Log(format string, a ...interface{}) {
	fmt.Printf(format, a)
}

func (pool *ThreadPool) Start() {
	for i := uint16(0); i < pool.count; i++ {
		pool.exitCh[i] = make(chan bool)
		pool.startedch[i] = make(chan bool)
		pool.exitedCh[i] = make(chan bool)
		go handleFun(pool.startedch[i], pool.exitCh[i], pool.exitedCh[i],
			pool.process, pool.threadqueue.GetChannel(), i)

	}

	for i := uint16(0); i < pool.count; i++ {
		<-pool.startedch[i]
		close(pool.startedch[i])
	}

}

func (pool *ThreadPool) Stop() {

	for i := uint16(0); i < pool.count; i++ {
		pool.exitCh[i] <- true
	}
	for i := uint16(0); i < pool.count; i++ {
		<-pool.exitedCh[i]
		close(pool.exitedCh[i])
	}

}
