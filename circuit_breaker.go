/*
Шаблон Circuit Breaker (Размыкатель цепи) автоматически отключает сервисные функции в ответ на вероятную неисправность,
чтобы предотвратить более крупные или каскадные отказы, устранить повторяющиеся ошибки
и обеспечить разумную реакцию на ошибки.

Этот шаблон включает следующие компоненты:
Circuit - функция, взаимодействующая со службой.
Breaker - функция-замыкание с той же сигнатурой, что и функция Circuit.
*/

package go_cloud_patterns

import (
	"context"
	"errors"
	"sync"
	"time"
)

type Circuit func(ctx context.Context) (string, error)

func Breaker(circuit Circuit, failureThreshold uint) Circuit {
	var consecutiveFailures int = 0 // Количество последовательных сбоев
	var lastAttempt = time.Now()
	var m sync.RWMutex

	return func(ctx context.Context) (string, error) {
		m.RLock() // Устанавливаем блокировку чтения

		d := consecutiveFailures - int(failureThreshold)

		if d >= 0 {
			shouldRetryAt := lastAttempt.Add(time.Second * 2 << d)
			if !time.Now().After(shouldRetryAt) {
				m.RUnlock()
				return "", errors.New("service unreachable")
			}
		}

		m.RUnlock() // Снимаем блокировку чтения

		response, err := circuit(ctx) // Посылаем запрос, как обычно

		m.Lock() // Блокировка общих ресурсов
		defer m.Unlock()

		lastAttempt = time.Now() // Фиксируем время попытки

		if err != nil {
			consecutiveFailures++ // Увеличиваем счетчик сбоев
			return response, err  // И возвращаем ошибку
		}

		consecutiveFailures = 0 // Сброс счетчика сбоев

		return response, nil
	}
}
