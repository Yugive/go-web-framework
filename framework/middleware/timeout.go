package middleware

import (
	"context"
	"fmt"
	"github.com/ngyugive/go-web-framework/framework"
	"log"
	"time"
)

func Timeout(d time.Duration) framework.ControllerHandler {
	return func(c *framework.Context) error {
		finish := make(chan struct{}, 1)
		panicChan := make(chan interface{}, 1)

		durationCtx, cancel := context.WithTimeout(c.BaseContext(), d)
		defer cancel()

		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()
			c.Next()

			finish <- struct{}{}
		}()

		select {
		case p := <-panicChan:
			c.SetStatus(500).Json("inner error")
			log.Println(p)
		case <-finish:
			fmt.Println("finish")
		case <-durationCtx.Done():
			c.SetHasTimeout()
			c.SetStatus(500).Json("timed out")
		}
		return nil
	}
}
