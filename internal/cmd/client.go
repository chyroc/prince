package cmd

import (
	"github.com/chyroc/prince/internal/transferclient"
	"sync"
)

func RunClient(transferHost string) error {
	wg := new(sync.WaitGroup)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				transferclient.RunTransferClient(transferHost)
			}
		}()
	}
	wg.Wait()

	return nil
}
