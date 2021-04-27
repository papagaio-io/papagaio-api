package test

import (
	"testing"

	"github.com/golang/mock/gomock"
)

func TestCreateOrganizationOK(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
}
