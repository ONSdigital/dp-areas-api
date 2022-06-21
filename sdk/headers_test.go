package sdk

import (
	"net/http"
	"testing"

	dprequest "github.com/ONSdigital/dp-net/v2/request"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAdd(t *testing.T) {
	t.Parallel()

	req := &http.Request{
		Header: http.Header{},
	}

	Convey("Given the sdk Headers struct is nil", t, func() {
		var headers Headers

		Convey("When calling the Add method on Headers", func() {
			headers.Add(req)

			Convey("Then no headers is set on the request", func() {
				So(req.Header, ShouldBeEmpty)
			})
		})
	})

	Convey("Given the sdk Headers struct contains a value for Service Auth Token", t, func() {
		headers := &Headers{
			ServiceAuthToken: "4egqsf4378gfwqf44356h",
		}

		Convey("When calling the Add method on Headers", func() {
			headers.Add(req)

			expectedHeader := dprequest.BearerPrefix + headers.ServiceAuthToken
			Convey("Then an Authorization header is set on the request", func() {
				So(req.Header, ShouldContainKey, dprequest.AuthHeaderKey)
				So(req.Header[dprequest.AuthHeaderKey], ShouldHaveLength, 1)
				So(req.Header[dprequest.AuthHeaderKey][0], ShouldEqual, expectedHeader)
			})
		})
	})

	Convey("Given the sdk Headers struct contains a value for User Auth Token", t, func() {
		headers := &Headers{
			UserAuthToken: "3wOx0VLG6xnbyDNjfGP4",
		}

		Convey("When calling the Add method on Headers", func() {
			headers.Add(req)

			Convey("Then an X-Florence-Token header is set on the request", func() {
				So(req.Header, ShouldContainKey, dprequest.FlorenceHeaderKey)
				So(req.Header[dprequest.FlorenceHeaderKey], ShouldHaveLength, 1)
				So(req.Header[dprequest.FlorenceHeaderKey][0], ShouldEqual, headers.UserAuthToken)
			})
		})
	})
}
