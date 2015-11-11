package japi

import (
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSend(t *testing.T) {

	Convey("Send Tests", t, func() {

		response := httptest.NewRecorder()

		Convey("->SendObject()", func() {

		})

		Convey("->SendList()", func() {

		})
	})
}
