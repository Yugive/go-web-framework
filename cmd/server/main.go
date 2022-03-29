package main

import (
	"context"
	"fmt"
	"github.com/ngyugive/go-web-framework/framework"
	"github.com/ngyugive/go-web-framework/framework/middleware"
	_ "github.com/ngyugive/goat_kit"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	core := framework.NewCore()
	core.Use(middleware.Recovery(), middleware.Cost())
	registerRouter(core)
	server := &http.Server{
		Handler: core,
		Addr:    ":8084",
	}

	go func() {
		server.ListenAndServe()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	timeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(timeCtx); err != nil {
		log.Fatal(err)
	}
}

func registerRouter(core *framework.Core) {
	//core.Get("foo", FooControllerHandler)
	core.Get("/", middleware.Timeout(1*time.Second), UserLoginController)

	subjectApi := core.Group("/subject")
	{
		subjectApi.Use(middleware.Timeout(500 * time.Millisecond))
		subjectApi.Delete("/:id", SubjectDelController)
		subjectApi.Get("/:id", SubjectGetController)
		subjectApi.Get("/list/all", SubjectListController)
		subjectApi.Put("/:id", SubjectUpdateController)
		subjectInnerApi := subjectApi.Group("/info")
		{
			subjectInnerApi.Get("/name", SubjectNameController)
		}
	}
}

func UserLoginController(c *framework.Context) error {
	c.Json("ok, userController 444")
	return nil
}

func SubjectAddController(c *framework.Context) error {
	c.SetOkStatus().Json("ok, subject SubjectAddController")
	return nil
}

func SubjectListController(c *framework.Context) error {
	c.SetOkStatus().Json("ok, SubjectListController")
	return nil
}

func SubjectDelController(c *framework.Context) error {
	c.SetOkStatus().Json("ok, SubjectDelController")
	return nil
}

func SubjectUpdateController(c *framework.Context) error {
	c.SetOkStatus().Json("ok, SubjectUpdateController")
	return nil
}

func SubjectGetController(c *framework.Context) error {
	subjectId, _ := c.ParamInt("id", 0)
	c.SetOkStatus().Json("ok, SubjectGetController:" + fmt.Sprint(subjectId))
	return nil
}

func SubjectNameController(c *framework.Context) error {
	c.SetOkStatus().Json("ok, SubjectNameController")
	return nil
}
