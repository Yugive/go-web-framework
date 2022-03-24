package main

import (
	"go-web-framework/framework"
	"go-web-framework/framework/middleware"
	"time"
)

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
