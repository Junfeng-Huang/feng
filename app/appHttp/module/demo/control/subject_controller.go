package demo

import (
	"feng/framework"
	"fmt"
)

func SubjectAddController(c *framework.Context) {
	c.SetOkStatus().Json("ok, SubjectAddController")

}

func SubjectListController(c *framework.Context) {
	c.SetOkStatus().Json("ok, SubjectListController")

}

func SubjectDelController(c *framework.Context) {
	c.SetOkStatus().Json("ok, SubjectDelController")

}

func SubjectUpdateController(c *framework.Context) {
	c.SetOkStatus().Json("ok, SubjectUpdateController")

}

func SubjectGetController(c *framework.Context) {
	subjectId, _ := c.ParamInt("id", 0)
	c.SetOkStatus().Json("ok, SubjectGetController:" + fmt.Sprint(subjectId))

}

func SubjectNameController(c *framework.Context) {
	c.SetOkStatus().Json("ok, SubjectNameController")

}
