package service

import (
	"testing"
)

func TestServiceConf(t *testing.T) {
	c := ServiceConf{
		Name: "foo",
		Mode: "dev",
	}
	c.MustSetUp()
}
