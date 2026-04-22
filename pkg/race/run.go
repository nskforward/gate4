package race

import (
	"context"
	"sync"
)

func Run(ctx context.Context, funcs ...func(context.Context) error) error {
	if len(funcs) == 0 {
		return nil
	}

	// Создаём дочерний контекст с возможностью отмены.
	ctx, cancel := context.WithCancel(ctx)
	defer cancel() // Гарантированно освобождаем ресурсы контекста.

	var wg sync.WaitGroup
	wg.Add(len(funcs))

	// Канал для получения результата первой завершившейся функции.
	// Буфер размера 1, чтобы отправитель не блокировался.
	resultCh := make(chan error, 1)

	// once гарантирует, что только первый результат будет отправлен в канал.
	var once sync.Once

	// Запускаем каждую функцию в отдельной горутине.
	for _, f := range funcs {
		go func(fn func(context.Context) error) {
			defer wg.Done()
			err := fn(ctx)
			once.Do(func() {
				cancel()
				resultCh <- err
			})
		}(f)
	}

	// Ожидаем первый результат (или завершение всех горутин, если контекст отменён извне).
	firstErr := <-resultCh

	// Ждём, пока все горутины корректно завершатся после отмены контекста.
	wg.Wait()

	return firstErr
}
