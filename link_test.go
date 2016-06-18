package jsh

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLink(t *testing.T) {

	Convey("Link Tests", t, func() {

		Convey("->MarshalJSON()", func() {

			Convey("should marshal as empty string", func() {
				l := NewLink("")

				jData, err := json.Marshal(l)
				So(err, ShouldBeNil)
				So(string(jData), ShouldEqual, `""`)
			})

			Convey("should marshal as string when no metadata is present", func() {
				l := NewLink("/stringlink")

				jData, err := json.Marshal(&l)
				So(err, ShouldBeNil)
				So(string(jData), ShouldEqual, `"/stringlink"`)
			})

			Convey("should marshal as an object when empty metadata is present", func() {
				l := NewMetaLink("/metalink", make(map[string]interface{}))

				jData, err := json.Marshal(&l)
				So(err, ShouldBeNil)
				So(string(jData), ShouldEqual, `{"href":"/metalink"}`)
			})

			Convey("should marshal as an object when metadata is present", func() {
				meta := make(map[string]interface{})
				meta["count"] = 10
				l := NewMetaLink("/metalink", meta)

				jData, err := json.Marshal(&l)
				So(err, ShouldBeNil)
				So(string(jData), ShouldEqual, `{"href":"/metalink","meta":{"count":10}}`)
			})
		})

		Convey("->UnmarshalJSON()", func() {

			Convey("should handle a string link", func() {
				jStr := `"/stringlink"`

				l := Link{}
				err := l.UnmarshalJSON([]byte(jStr))
				So(err, ShouldBeNil)
				So(l.HREF, ShouldEqual, "/stringlink")
			})

			Convey("should handle a meta link", func() {
				jMeta := `{"href": "/metalink"}`

				l := Link{}
				err := l.UnmarshalJSON([]byte(jMeta))
				So(err, ShouldBeNil)
				So(l.HREF, ShouldEqual, "/metalink")
			})

			Convey("should handle a meta link with metadata", func() {
				jMeta := `{"href": "/metalink", "meta": {"count": 10}}`

				l := Link{}
				err := l.UnmarshalJSON([]byte(jMeta))
				So(err, ShouldBeNil)
				So(l.HREF, ShouldEqual, "/metalink")
				So(l.Meta, ShouldNotBeEmpty)
				So(l.Meta["count"], ShouldEqual, 10)
			})
		})
	})
}
