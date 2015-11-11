package jsh

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParsing(t *testing.T) {

	Convey("Request Tests", t, func() {

		Convey("->ParseObject()", func() {
			objectJSON := `{"data": {"type": "user", "id": "sweetID123", "attributes": {"ID":"123"}}}`

			closer := createIOCloser([]byte(objectJSON))

			object, err := ParseObject(closer)
			So(err, ShouldBeNil)
			So(object, ShouldNotBeEmpty)
			So(object.Type, ShouldEqual, "user")
			So(object.ID, ShouldEqual, "sweetID123")
			So(object.Attributes, ShouldResemble, map[string]interface{}{"ID": "123"})
		})

		Convey("->ParseList()", func() {
			listJSON :=
				`{"data": [
	{"type": "user", "id": "sweetID123", "attributes": {"ID": "123"}},
	{"type": "user", "id": "sweetID456", "attributes": {"ID": "456"}}
]}`

			closer := createIOCloser([]byte(listJSON))

			list, err := ParseList(closer)
			So(err, ShouldBeNil)
			So(len(list), ShouldEqual, 2)

			object := list[1]
			So(object.Type, ShouldEqual, "user")
			So(object.ID, ShouldEqual, "sweetID456")
			So(object.Attributes, ShouldResemble, map[string]interface{}{"ID": "456"})
		})
	})
}

func createIOCloser(data []byte) io.ReadCloser {
	reader := bytes.NewReader(data)
	return ioutil.NopCloser(reader)
}
