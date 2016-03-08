package jsh

import (
	"net/http"
	"testing"

	"encoding/json"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDocument(t *testing.T) {

	Convey("Document Tests", t, func() {

		doc := New()

		Convey("->New()", func() {
			So(doc.JSONAPI.Version, ShouldEqual, JSONAPIVersion)
		})

		Convey("->HasErrors()", func() {
			err := &Error{Status: 400}
			addErr := doc.AddError(err)
			So(addErr, ShouldBeNil)

			So(doc.HasErrors(), ShouldBeTrue)
		})

		Convey("->HasData()", func() {
			obj, err := NewObject("1", "user", nil)
			So(err, ShouldBeNil)

			doc.Data = append(doc.Data, obj)
			So(doc.HasData(), ShouldBeTrue)
		})

		Convey("->AddObject()", func() {

			obj, err := NewObject("1", "user", nil)
			So(err, ShouldBeNil)

			Convey("should successfully add an object", func() {
				err := doc.AddObject(obj)
				So(err, ShouldBeNil)
				So(len(doc.Data), ShouldEqual, 1)
			})

			Convey("should prevent multiple data objects in ObjectMode", func() {
				err := doc.AddObject(obj)
				So(err, ShouldBeNil)

				err = doc.AddObject(obj)
				So(err, ShouldNotBeNil)
			})
		})

		Convey("->AddError()", func() {
			testError := &Error{Status: 400}

			Convey("should successfully add a valid error", func() {
				err := doc.AddError(testError)
				So(err, ShouldBeNil)
				So(len(doc.Errors), ShouldEqual, 1)
				So(doc.Mode, ShouldEqual, ErrorMode)
			})

			Convey("should error if validation fails while adding an error", func() {
				badError := &Error{
					Title:  "Invalid",
					Detail: "So badly",
				}

				err := doc.AddError(badError)
				So(err, ShouldNotBeNil)
				So(doc.Errors, ShouldBeEmpty)
			})
		})

		Convey("->First()", func() {

			Convey("should not explode for nil data", func() {
				data := doc.First()
				So(data, ShouldBeNil)
			})
		})

		Convey("->Build()", func() {

			testObject := &Object{
				ID:     "1",
				Type:   "Test",
				Status: http.StatusAccepted,
			}

			Convey("should accept an object", func() {
				doc := Build(testObject)

				So(doc.Data, ShouldResemble, List{testObject})
				So(doc.Status, ShouldEqual, http.StatusAccepted)
				So(doc.Mode, ShouldEqual, ObjectMode)
			})

			Convey("should accept a list", func() {
				list := List{testObject}
				doc := Build(list)

				So(doc.Data, ShouldResemble, list)
				So(doc.Status, ShouldEqual, http.StatusOK)
				So(doc.Mode, ShouldEqual, ListMode)
			})

			Convey("should accept an error", func() {
				err := &Error{Status: http.StatusInternalServerError}
				doc := Build(err)

				So(doc.Errors, ShouldNotBeEmpty)
				So(doc.Status, ShouldEqual, err.Status)
				So(doc.Mode, ShouldEqual, ErrorMode)
			})
		})

		Convey("->Validate()", func() {

			testObject := &Object{
				ID:     "1",
				Type:   "Test",
				Status: http.StatusAccepted,
			}

			testObjectForInclusion := &Object{
				ID:   "1",
				Type: "Included",
			}

			req := &http.Request{Method: "GET"}

			Convey("should not accept an included object without objects in data", func() {
				doc := New()
				doc.Included = append(doc.Included, testObjectForInclusion)
				doc.Status = http.StatusOK

				validationErrors := doc.Validate(req, true)

				So(validationErrors, ShouldNotBeNil)
			})

			Convey("should accept an object in data and an included object", func() {
				doc := Build(testObject)
				doc.Included = append(doc.Included, testObjectForInclusion)

				validationErrors := doc.Validate(req, true)

				So(validationErrors, ShouldBeNil)
				So(doc.Data, ShouldResemble, List{testObject})
				So(doc.Included, ShouldNotBeEmpty)
				So(doc.Included[0], ShouldResemble, testObjectForInclusion)
				So(doc.Status, ShouldEqual, http.StatusAccepted)
			})

		})
	})
}

func TestDocumentMarshaling(t *testing.T) {

	Convey("Document Marshal Tests", t, func() {

		doc := New()

		Convey("->MarshalJSON()", func() {

			testObject := &Object{
				ID:     "1",
				Type:   "Test",
				Status: http.StatusAccepted,
			}

			Convey("ListMode", func() {

				Convey("should marshal a list with a single element as an array", func() {
					list := List{testObject}
					doc := Build(list)

					rawJSON, err := json.Marshal(doc)
					So(err, ShouldBeNil)

					m := map[string]json.RawMessage{}
					err = json.Unmarshal(rawJSON, &m)
					So(err, ShouldBeNil)

					data := string(m["data"])
					So(data, ShouldStartWith, "[")
					So(data, ShouldEndWith, "]")
				})

				Convey("should marshal an empty list", func() {
					list := List{}
					doc := Build(list)

					rawJSON, err := json.Marshal(doc)
					So(err, ShouldBeNil)

					m := map[string]json.RawMessage{}
					err = json.Unmarshal(rawJSON, &m)
					So(err, ShouldBeNil)

					data := string(m["data"])
					So(data, ShouldEqual, "[]")
				})
			})

			Convey("ObjectMode", func() {

				doc := New()
				doc.Mode = ObjectMode

				Convey("should marshal a single object as an object", func() {
					addErr := doc.AddObject(testObject)
					So(addErr, ShouldBeNil)

					rawJSON, err := json.Marshal(doc)
					So(err, ShouldBeNil)

					m := map[string]json.RawMessage{}
					err = json.Unmarshal(rawJSON, &m)
					So(err, ShouldBeNil)

					data := string(m["data"])
					So(data, ShouldStartWith, "{")
					So(data, ShouldEndWith, "}")
				})

				Convey("null cases", func() {

					Convey("should marshal nil to null", func() {
						doc.Data = nil
						rawJSON, err := json.Marshal(doc)
						So(err, ShouldBeNil)

						m := map[string]json.RawMessage{}
						err = json.Unmarshal(rawJSON, &m)
						So(err, ShouldBeNil)

						data := string(m["data"])
						So(data, ShouldEqual, "null")
					})

					Convey("should marshal an empty list to null", func() {
						doc.Data = List{}
						rawJSON, err := json.Marshal(doc)
						So(err, ShouldBeNil)

						m := map[string]json.RawMessage{}
						err = json.Unmarshal(rawJSON, &m)
						So(err, ShouldBeNil)

						data := string(m["data"])
						So(data, ShouldEqual, "null")
					})
				})
			})

			Convey("ErrorMode", func() {

				Convey("should not include 'data' field for error response", func() {
					doc.AddError(ISE("Test Error"))

					rawJSON, err := json.Marshal(doc)
					So(err, ShouldBeNil)

					jMap := map[string]json.RawMessage{}
					err = json.Unmarshal(rawJSON, &jMap)
					So(err, ShouldBeNil)

					_, exists := jMap["data"]
					So(exists, ShouldBeFalse)

					errors := string(jMap["errors"])
					So(errors, ShouldNotBeEmpty)
					So(errors, ShouldStartWith, "[")
					So(errors, ShouldEndWith, "]")
				})
			})
		})
	})
}
