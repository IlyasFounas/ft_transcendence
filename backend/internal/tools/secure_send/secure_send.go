package securesend

import (
	"log"
	"context"
)
// SecureSend sends data on c, unless ctx is cancelled first.
//
// The hub channels have a single reader: the Run()'s loops. Once Run()'s returns on
// shutdown, a plain send on those channels would block forever, leaking the
// caller's goroutine and hanging the WaitGroup. Aborting on ctx.Done() means
// the message is dropped, but the goroutine can exit.
func SecureSend[T any](
	ctx context.Context, c chan<- *T, data *T, errMsg string) bool {
	select {
	case c <- data:
		return true
	case <-ctx.Done():
		log.Println(errMsg)
		return false
	}
}
