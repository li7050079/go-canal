package web

import (
	"context"
	"fmt"
	"time"
)

func getBucket(capacityPs, maxCapacity int) (chan struct{}, *time.Ticker) {
	var bucketToken = make(chan struct{}, maxCapacity)
	ticker := time.NewTicker(time.Second / time.Duration(capacityPs))
	go func() {
		for {
			select {
			case <-ticker.C:
				bucketToken <- struct{}{}
			default:
			}
		}
	}()
	return bucketToken, ticker
}

func getToken(ctx context.Context, bucket *chan struct{}, wait bool) bool {
	if wait {
		select {
		case <-*bucket:
			return true
		case <-ctx.Done():
			fmt.Println("请求超时...")
			return false
		}
	} else { // 不阻塞
		select {
		case <-*bucket:
			return true
		default:
			return false
		}
	}
}

// 模拟qps 100， 超时等待时间为1s，模拟1000个请求进来
// 令牌桶最大每秒产生30个令牌，最大突发为50（桶最大容量）
func test1() {
	bucket, ticker := getBucket(30, 50)
	i := 0
	total := 0
	fail := 0
	go func() {
		for { // qps 100
			i++
			time.Sleep(time.Millisecond * 10)
			ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*100)
			ok := getToken(ctx, &bucket, false)
			if ok {
				total++
			} else {
				fail++
			}
			fmt.Println("当前获取token是否成功：", ok)
			if i >= 1000 {
				ticker.Stop()
				break
			}
		}
	}()
	time.Sleep(time.Second * 20)
	fmt.Println(i, total, fail)
}

// 模拟qps 100， 超时等待时间为1s，模拟1000个请求进来
// 令牌桶最大每秒产生30个令牌，最大突发为50（桶最大容量）
func test2() {
	bucket, ticker := getBucket(300, 500)
	i := 0
	total := 0
	fail := 0
	go func() {
		for { // qps 100
			i++
			time.Sleep(time.Millisecond * 10)
			ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*20)
			ok := getToken(ctx, &bucket, true)
			if ok {
				total++
			} else {
				fail++
			}
			fmt.Println("当前获取token是否成功：", ok)
			if i >= 1000 {
				ticker.Stop()
				break
			}
		}
	}()
	time.Sleep(time.Second * 20)
	fmt.Println(i, total, fail)
}

func main() {
	//test1()
	test2()
}
