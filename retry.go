/*
Шаблон Retry (Повтор) учитывает возможный временный характер ошибки
в распределенной системе и осуществляет повторные попытки выполнить
неудачную операцию.

Этот шаблон включает следующие компоненты:
Effector - функция, взаимодействующая со службой.
Retry - функция, принимающая Effector и возвращающая замыкание с той же сигнатурой, что и Effector.
*/

package go_cloud_patterns

import (
	"context"
	"log"
	"time"
)

type Effector func(context.Context) (string, error)

func Retry(effector Effector, retries int, delay time.Duration) Effector {
	// retries повторов через каждые delay интервалов времени
	return func(ctx context.Context) (string, error) {
		for r := 0; ; r++ {
			response, err := effector(ctx)
			if err == nil || r >= retries {
				return response, err
			}

			log.Printf("Attempt %d failed; retrying in %v", r+1, delay)

			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return "", ctx.Err()
			}
		}
	}
}

//var count int
//
//func EmulateTransientError(ctx context.Context) (string, error){
//	count++
//
//	if count <= 3 {
//		return "intentional fail", errors.New("error")
//	} else {
//		return "success", nil
//	}
//}
//
//func main() {
//	r := Retry(EmulateTransientError, 5, 2*time.Second)
//
//	res, err := r(context.Background())
//
//	fmt.Println(res, err)
//}
