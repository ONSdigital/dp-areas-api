package sdk

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestOptionsLanguage(t *testing.T) {
	t.Parallel()

	Convey("Given the sdk Options struct contains a Lang variable for english", t, func() {
		options := Options{
			Lang: English,
		}

		Convey("When calling String method", func() {
			language, err := options.Lang.String()

			Convey("Then \"en\" value is returned", func() {
				So(language, ShouldEqual, "en")

				Convey("And no error is returned", func() {
					So(err, ShouldBeNil)
				})
			})
		})
	})

	Convey("Given the sdk Options struct contains a Lang variable for welsh", t, func() {
		options := Options{
			Lang: Welsh,
		}

		Convey("When calling String method", func() {
			language, err := options.Lang.String()

			Convey("Then \"cy\" value is returned", func() {
				So(language, ShouldEqual, "cy")

				Convey("And no error is returned", func() {
					So(err, ShouldBeNil)
				})
			})
		})
	})

	Convey("Given the sdk Options struct contains an unknown Lang variable", t, func() {
		unknownLanguage := Language("alien language")
		options := Options{
			Lang: unknownLanguage,
		}

		Convey("When calling String method", func() {
			language, err := options.Lang.String()

			Convey("Then empty string is returned", func() {
				So(language, ShouldEqual, "")

				Convey("And error is returned", func() {
					So(err, ShouldResemble, ErrUnrecognisedLanguage(unknownLanguage))
				})
			})
		})
	})
}
