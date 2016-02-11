package jsh

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRelationship(t *testing.T) {
	Convey("ResourceLinkage Tests", t, func() {
		Convey("->UnmarshalJSON()", func() {

			Convey("should handle a linkage object", func() {
				jObj := `{"type": "testRelationship", "id": "ID456"}`

				rl := ResourceLinkage{}
				err := rl.UnmarshalJSON([]byte(jObj))
				So(err, ShouldBeNil)
				So(len(rl), ShouldEqual, 1)
			})

			Convey("should handle a linkage list", func() {
				jList := `[
					{"type": "testRelationship", "id": "ID456"},
					{"type": "testRelationship", "id": "ID789"}
				]`

				rl := ResourceLinkage{}
				err := rl.UnmarshalJSON([]byte(jList))
				So(err, ShouldBeNil)
				So(len(rl), ShouldEqual, 2)
			})
		})
	})
}
