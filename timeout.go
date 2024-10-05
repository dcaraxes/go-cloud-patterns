/*
Шаблон Timeout (Тайм-аут) позволяет процессу прекратить ожидание ответа,
когда станет очевидно, что его можно вообще не получить.

Client - клиент, обращающийся к SlowFunction.
SlowFunction - функция, возвращающая ответ, необходимый клиенту, и выполняющаяся
очень долго.
Timeout - функция-обертка вокруг SlowFunction, которая реализует логику тайм-аута.
*/

package go_cloud_patterns

import (
	"context"
)

type SlowFunction func(string) (string, error)

type WithContext func(context.Context, string) (string, error)

func Timeout(f SlowFunction) WithContext {
	return func(ctx context.Context, arg string) (string, error) {
		chres := make(chan string)
		cherr := make(chan error)

		go func() {
			res, err := f(arg)
			chres <- res
			cherr <- err
		}()

		select {
		case res := <-chres:
			return res, <-cherr
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}
}

//func main() {
//	ctx := context.Background()
//	ctxt, cancel := context.WithTimeout(ctx, 1*time.Second)
//	defer cancel
//
//	timeout := Timeout(Slow)
//	res, err := timeout(ctxt, "some input")
//	fmt.Println(res, err)
//}
