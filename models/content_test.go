package models_test

import (
	"fmt"
	"testing"

	"github.com/ONSdigital/dp-topic-api/apierrors"
	"github.com/ONSdigital/dp-topic-api/models"
	. "github.com/smartystreets/goconvey/convey"
)

func TestContentValidation(t *testing.T) {
	Convey("Given an empty content, it is successfully validated", t, func() {
		content := models.Content{
			State: models.StateCreated.String(),
		}
		err := content.Validate()
		So(err, ShouldBeNil)
	})

	Convey("Given an content with no state supplied, it fails to validate  with the expected error", t, func() {
		content := models.Content{}
		err := content.Validate()
		So(err, ShouldResemble, apierrors.ErrTopicInvalidState)
	})

	Convey("Given an content with a state that does not correspond to any expected state, it fails to validate with the expected error", t, func() {
		content := models.Content{
			State: "wrong",
		}
		err := content.Validate()
		So(err, ShouldResemble, apierrors.ErrTopicInvalidState)
	})

	Convey("Given a fully populated valid content with a valid download variant, it is successfully validated", t, func() {
		content := models.Content{
			State: models.StatePublished.String(),
		}
		err := content.Validate()
		So(err, ShouldBeNil)
	})
}

func TestContentValidateTransitionFrom(t *testing.T) {
	Convey("Given an existing content in an created state", t, func() {
		existing := &models.Content{
			State: models.StateCreated.String(),
		}

		Convey("When we try to transition to an completed state", func() {
			content := &models.Content{
				State: models.StateCompleted.String(),
			}
			err := content.ValidateTransitionFrom(existing)
			So(err, ShouldBeNil)
		})

		Convey("When we try to transition to a published state", func() {
			content := &models.Content{
				State: models.StatePublished.String(),
			}
			err := content.ValidateTransitionFrom(existing)
			So(err, ShouldResemble, apierrors.ErrTopicStateTransitionNotAllowed)
		})
	})

	Convey("Given an existing content in a Completed state", t, func() {
		existing := &models.Content{
			State: models.StateCompleted.String(),
		}

		Convey("When we try to transition to a published state", func() {
			content := &models.Content{
				State: models.StatePublished.String(),
			}
			err := content.ValidateTransitionFrom(existing)
			So(err, ShouldBeNil)
		})

		Convey("When we try to transition to an created state", func() {
			content := &models.Content{
				State: models.StateCreated.String(),
			}
			err := content.ValidateTransitionFrom(existing)
			So(err, ShouldResemble, apierrors.ErrTopicStateTransitionNotAllowed)
		})
	})
}

func TestContentStateTransitionAllowed(t *testing.T) {
	Convey("Given an content in created state", t, func() {
		content := models.Content{
			State: models.StateCreated.String(),
		}
		contentValidateTransitionsToCreated(content)
	})

	Convey("Given an content with a wrong state value, then no transition is allowed", t, func() {
		content := models.Content{State: "wrong"}
		contentValidateTransitionsToCreated(content)
	})

	Convey("Given an content without state, then created state is assumed when checking for transitions", t, func() {
		content := models.Content{}
		contentValidateTransitionsToCreated(content)
	})
}

// validateTransitionsToCreated validates that the provided content can transition to created state,
// and not to any forbidden of invalid state
func contentValidateTransitionsToCreated(content models.Content) {
	Convey("Then an allowed transition is successfully checked", func() {
		So(content.StateTransitionAllowed(models.StateCompleted.String()), ShouldBeTrue)
	})
	Convey("Then a forbidden transition to a valid state is not allowed", func() {
		So(content.StateTransitionAllowed(models.StatePublished.String()), ShouldBeFalse)
	})
	Convey("Then a transition to an invalid state is not allowed", func() {
		So(content.StateTransitionAllowed("wrong"), ShouldBeFalse)
	})
}

func TestAppendLinkInfo(t *testing.T) {
	Convey("Given a nil ptr for Links in published state", t, func() {
		var result models.ContentResponseAPI

		result.AppendLinkInfo("spotlight", nil, "9", "published")
		fmt.Printf("%+v", result)
		Convey("Then the expected result of zero items", func() {
			So(result.Count, ShouldEqual, 0)
			So(result.Offset, ShouldEqual, 0)
			So(result.Limit, ShouldEqual, 0)
			So(result.TotalCount, ShouldEqual, 0)

			So(result.Items, ShouldBeNil)
		})
	})

	Convey("Given an empty list of article Links in published state", t, func() {
		var result models.ContentResponseAPI
		var articleObjects []models.TypeLinkObject

		result.AppendLinkInfo("spotlight", &articleObjects, "9", "published")
		fmt.Printf("%+v", result)
		Convey("Then the expected result is of zero items", func() {
			So(result.Count, ShouldEqual, 0)
			So(result.Offset, ShouldEqual, 0)
			So(result.Limit, ShouldEqual, 0)
			So(result.TotalCount, ShouldEqual, 0)

			So(result.Items, ShouldBeNil)
		})
	})

	Convey("Given a list containing one spotlight Link in published state", t, func() {
		var result models.ContentResponseAPI
		var spotlightObjects = []models.TypeLinkObject{
			{
				HRef:  "/a 1st",
				Title: "first",
			},
		}

		result.AppendLinkInfo("spotlight", &spotlightObjects, "9", "published")
		fmt.Printf("%+v", result)
		Convey("Then the expected result contains one Item", func() {
			So(result.Count, ShouldEqual, 0)
			So(result.Offset, ShouldEqual, 0)
			So(result.Limit, ShouldEqual, 0)
			So(result.TotalCount, ShouldEqual, 1)

			So((*result.Items)[0].Title, ShouldEqual, "first")
			So((*result.Items)[0].Type, ShouldEqual, "spotlight")
			So((*result.Items)[0].State, ShouldEqual, "published")
			So((*result.Items)[0].Links.Self.HRef, ShouldEqual, "/a 1st")
			So((*result.Items)[0].Links.Topic.HRef, ShouldEqual, "/topic/9")
			So((*result.Items)[0].Links.Topic.ID, ShouldEqual, "9")
		})
	})

	Convey("Given a list containing two spotlight Links in published state arranged by Href in non alphabetical order ", t, func() {
		var result models.ContentResponseAPI
		var spotlightObjects = []models.TypeLinkObject{
			{
				HRef:  "/b 2nd",
				Title: "second",
			},
			{
				HRef:  "/a 1st",
				Title: "first",
			},
		}

		result.AppendLinkInfo("spotlight", &spotlightObjects, "9", "published")
		fmt.Printf("%+v", result)
		Convey("Then the expected result is sorted by Href", func() {
			So(result.Count, ShouldEqual, 0)
			So(result.Offset, ShouldEqual, 0)
			So(result.Limit, ShouldEqual, 0)
			So(result.TotalCount, ShouldEqual, 2)

			So((*result.Items)[0].Title, ShouldEqual, "first")
			So((*result.Items)[0].Type, ShouldEqual, "spotlight")
			So((*result.Items)[0].State, ShouldEqual, "published")
			So((*result.Items)[0].Links.Self.HRef, ShouldEqual, "/a 1st")
			So((*result.Items)[0].Links.Topic.HRef, ShouldEqual, "/topic/9")
			So((*result.Items)[0].Links.Topic.ID, ShouldEqual, "9")

			So((*result.Items)[1].Title, ShouldEqual, "second")
			So((*result.Items)[1].Type, ShouldEqual, "spotlight")
			So((*result.Items)[1].State, ShouldEqual, "published")
			So((*result.Items)[1].Links.Self.HRef, ShouldEqual, "/b 2nd")
			So((*result.Items)[1].Links.Topic.HRef, ShouldEqual, "/topic/9")
			So((*result.Items)[1].Links.Topic.ID, ShouldEqual, "9")
		})
	})
}
