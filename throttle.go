/*
Шаблон Throttle (Дроссельная заслонка) ограничивает частоту вызовов функции некоторым предельным числом
вызовов в единицу времени.

Этот шаблон включает следующие компоненты:
Effector - функция, частоту вызовов которой нужно ограничить.
Throttle - функция, принимающая Effector и возвращающая замыкание с той же сигнатурой, что и Effector.
*/

package go_cloud_patterns

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Throttle Реализация алгоритма "token bucket(корзина жетонов)"
// с самой простой стратегией "возврат ошибки"
func Throttle(e Effector, max uint, refill uint, d time.Duration) Effector {
	var tokens = max
	var once sync.Once

	return func(ctx context.Context) (string, error) {
		if ctx.Err() != nil {
			return "", ctx.Err()
		}

		once.Do(func() {
			ticker := time.NewTicker(d)

			go func() {
				defer ticker.Stop()

				for {
					select {
					case <-ctx.Done():
						return
					case <-ticker.C:
						t := tokens + refill
						if t > max {
							t = max
						}
						tokens = t
					}
				}
			}()
		})

		if tokens <= 0 {
			return "", fmt.Errorf("too many calls")
		}

		tokens--

		return e(ctx)
	}
}
