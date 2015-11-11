package japi

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRequest(t *testing.T) {

	Convey("Request Tests", t, func() {

		Convey("->ParseObject()", func() {
			jsonStr := `{"data": {"type": "user", "id": "sweetID123", "attributes": {"ID":"123"}}}`

			closer := createIOCloser([]byte(jsonStr))

			object, err := ParseObject(closer)
			So(err, ShouldBeNil)
			So(object, ShouldNotBeEmpty)
			So(err, ShouldBeNil)
			So(object.Type, ShouldEqual, "user")
			So(object.ID, ShouldEqual, "sweetID123")
			So(object.Attributes, ShouldResemble, map[string]interface{}{"ID": "123"})
		})
	})
}

func createIOCloser(data []byte) io.ReadCloser {
	reader := bytes.NewReader(data)
	return ioutil.NopCloser(reader)
}
