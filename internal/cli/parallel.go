package cli

import "sync"

func parallel(actions ...func()) {
	wg := &sync.WaitGroup{}
	for i := range actions {
		wg.Add(1)
		go func(wg *sync.WaitGroup, action func()) {
			action()
			wg.Done()
		}(wg, actions[i])
	}
	wg.Wait()
}
