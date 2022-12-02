package api

import (
	"github.com/ONSdigital/dp-topic-api/config"
	"github.com/ONSdigital/dp-topic-api/mocks"
	storetest "github.com/ONSdigital/dp-topic-api/store/mock"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNavigationGetNavigationHandler(t *testing.T) {
	Convey("Given a topic API in publishing mode (private endpoints enabled)", t, func() {
		cfg, err := config.Get()
		So(err, ShouldBeNil)
		mockedDataStore := &storetest.StorerMock{}
		permissions := mocks.NewAuthHandlerMock()

		api := GetWebAPIWithMocks(testContext, cfg, mockedDataStore, permissions)
		w := httptest.NewRecorder()
		So(w.Header(), ShouldBeEmpty)

		request, err := createRequestWithAuth("GET", "http://localhost:25300/navigation", nil)

		api.getNavigationHandler(w, request)

		So(w.Code, ShouldEqual, http.StatusOK)
		So(w.Header(), ShouldNotBeEmpty)
		So(w.Header().Get("Cache-Control"), ShouldEqual, "public, max-age=1800.000000")

	})

}
