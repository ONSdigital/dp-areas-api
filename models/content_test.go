package models_test

import (
	"testing"

	"github.com/ONSdigital/dp-topic-api/apierrors"
	"github.com/ONSdigital/dp-topic-api/models"
	. "github.com/smartystreets/goconvey/convey"
)

func TestContentValidation(t *testing.T) {
	Convey("Given an empty content, it is successfully validated", t, func() {
		content := models.Content{
			State: models.StateCreated.String(), //!!! remove topic out of the enum name
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

//!!! as the content functionality grows, add more tests to cover added state transitions.

///!!! these transitions should equally apply to content, therefore make all topic stat & transition stuff generic for topic & content.
