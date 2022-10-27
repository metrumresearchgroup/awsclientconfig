package awsclientconfig

import (
	"context"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/elisarver/zbt"
)

type Loginor interface {
	Login(ctx context.Context, sessionName string, optFns ...func(loadOptions *config.LoadOptions) error) (aws.Config, error)
}

func CreateGetClientFunc[T any](ctx context.Context, loginor Loginor, createClient func(ctx context.Context, mtx *sync.RWMutex, loginor Loginor) (T, error),
	refresh time.Duration) (func() (T, error), error) {
	// these are the closed-over vars we want to return to users.
	var (
		cli T
		err error
	)

	// limits read access when the client is being re-created.
	mtx := &sync.RWMutex{}

	// We have to wait for the first tick to occur.
	firstCreation := &sync.WaitGroup{}

	// the "work" is the first tick of this timer marking the WaitGroup
	// as done.
	firstCreation.Add(1)

	go func() {
		// ticker sets the cadence of refreshes. ZBT ticks immediately
		// upon creation.
		ticker := zbt.NewTicker(refresh)

		// once is to make sure we only ever mark done one time
		once := &sync.Once{}
		doneCh := ctx.Done()

		for {
			select {
			// global context is closed (server is shut down)
			case <-doneCh:
				break
			case <-ticker.C:
				cli, err = createClient(ctx, mtx, loginor)
				once.Do(firstCreation.Done)
				// we'll pass client and err on outside this loop
			}

			time.Sleep(1)
		}
	}()

	// This make sure the client's first creation ran.
	firstCreation.Wait()

	return func() (T, error) {
		mtx.RLock()
		defer mtx.RUnlock()

		return cli, err
	}, nil
}
