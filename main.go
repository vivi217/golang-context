package main

import (
	"context"
	"fmt"
	"time"
)

type cancelervivi interface {
	cancel(removeFromParent bool, err error)
	Done() <-chan struct{}
}

type Pen interface {
	Done() <-chan struct{}
}

type cancelCtxvivi struct {
	Pen
	a string
	b string
}

type timerCtx struct {
	cancelCtxvivi
	timer    *time.Timer // Under cancelCtx.mu.
	deadline time.Time
	done     <-chan struct{}
}

/*func (t *timerCtx) Done() <-chan struct{} {
	return t.done
}*/

func (t *timerCtx) cancel(removeFromParent bool, err error) {
}

func newCancelCtx(p Pen) cancelCtxvivi {
	return cancelCtxvivi{Pen: p}
}

func addChild(c cancelervivi) {
	child := make(map[cancelervivi]struct{})
	child[c] = struct{}{}
}
func test_cancel() {
	//var c canceler
	//c.cancel()
	//c.Done()
	var parent Pen
	c := &timerCtx{
		cancelCtxvivi: newCancelCtx(parent),
		deadline:      time.Now(),
	}
	addChild(c)
}

//如果主父协程超时了，那么子协程都会收到信号，但是也可以不予理会，但是 这个信号 是怎么在通道里面传递的呢
func main() {
	/*ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second) //3秒钟结束，会收到ctx.Done的通道信号
	defer cancel()

	go handle(ctx, 3*time.Second)
	select {
	case <-ctx.Done():
		fmt.Println("main", ctx.Err())
	}*/
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	go cancel_handle(ctx, cancelFunc)
	time.Sleep(time.Duration(5 * time.Second))
	cancelFunc()
	fmt.Printf("i am going to do other thing\n")
	time.Sleep(time.Duration(2 * time.Second))
	fmt.Printf("going to exit\n")
}

func doSome(num int) error {
	time.Sleep(time.Duration(2 * time.Second))
	//return errors.New("some error")
	return nil
}

func cancel_handle(ctx context.Context, cancelFunc context.CancelFunc) {
	for {
		select {
		case <-ctx.Done(): //父协程取消后， 这里收到的消息，然后退出了
			fmt.Printf("ctx.Done\n")
			return
		default:
			fmt.Println("exec default func")
			err := doSome(3)
			if err != nil {
				fmt.Printf("cancelFunc()\n")
				cancelFunc() // 这里可以是结束了它自己
				//return
			}
		}
	}

}

func handle(ctx context.Context, duration time.Duration) {
	go test(ctx, 2*time.Second)
	select {
	case <-ctx.Done():
		fmt.Println("handle", ctx.Err())
	case <-time.After(duration):
		fmt.Println("process request with", duration)
	}

}
func test(ctx context.Context, duration time.Duration) {
	select {
	case <-ctx.Done():
		fmt.Println("test handle", ctx.Err())
	case <-time.After(duration):
		fmt.Println("test process request with", duration)
	}
	//fmt.Println("test handle")
}
