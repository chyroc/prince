package helper

import (
	"io"
	"sync"
)

func Transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}

func TransferWait(c1, c2 io.ReadWriteCloser) {
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()
		Transfer(c1, c2)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		Transfer(c2, c1)
	}()
	wg.Wait()
}
