package models_test

import (
	"testing"
	"time"

	"github.com/ONSdigital/dp-topic-api/apierrors"
	"github.com/ONSdigital/dp-topic-api/models"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTopicValidation(t *testing.T) {
	t.Parallel()

	Convey("Given an empty topic, it is successfully validated", t, func() {
		topic := models.Topic{
			State: models.StateCreated.String(),
		}
		err := topic.Validate()
		So(err, ShouldBeNil)
	})

	Convey("Given an topic with no state supplied, it fails to validate  with the expected error", t, func() {
		topic := models.Topic{}
		err := topic.Validate()
		So(err, ShouldResemble, apierrors.ErrTopicInvalidState)
	})

	Convey("Given an topic with a state that does not correspond to any expected state, it fails to validate with the expected error", t, func() {
		topic := models.Topic{
			State: "wrong",
		}
		err := topic.Validate()
		So(err, ShouldResemble, apierrors.ErrTopicInvalidState)
	})

	Convey("Given a fully populated valid topic with a valid download variant, it is successfully validated", t, func() {
		topic := models.Topic{
			ID:    "123",
			State: models.StatePublished.String(),
		}
		err := topic.Validate()
		So(err, ShouldBeNil)
	})
}

func TestTopicValidateTransitionFrom(t *testing.T) {
	t.Parallel()

	Convey("Given an existing topic in an created state", t, func() {
		existing := &models.Topic{
			State: models.StateCreated.String(),
		}

		Convey("When we try to transition to an completed state", func() {
			topic := &models.Topic{
				State: models.StateCompleted.String(),
			}
			err := topic.ValidateTransitionFrom(existing)
			So(err, ShouldBeNil)
		})

		Convey("When we try to transition to a published state", func() {
			topic := &models.Topic{
				State: models.StatePublished.String(),
			}
			err := topic.ValidateTransitionFrom(existing)
			So(err, ShouldResemble, apierrors.ErrTopicStateTransitionNotAllowed)
		})
	})

	Convey("Given an existing topic in a Completed state", t, func() {
		existing := &models.Topic{
			State: models.StateCompleted.String(),
		}

		Convey("When we try to transition to a published state", func() {
			topic := &models.Topic{
				State: models.StatePublished.String(),
			}
			err := topic.ValidateTransitionFrom(existing)
			So(err, ShouldBeNil)
		})

		Convey("When we try to transition to an created state", func() {
			topic := &models.Topic{
				State: models.StateCreated.String(),
			}
			err := topic.ValidateTransitionFrom(existing)
			So(err, ShouldResemble, apierrors.ErrTopicStateTransitionNotAllowed)
		})
	})
}

func TestTopicStateTransitionAllowed(t *testing.T) {
	t.Parallel()

	Convey("Given an topic in created state", t, func() {
		topic := models.Topic{
			State: models.StateCreated.String(),
		}
		validateTransitionsToCreated(topic)
	})

	Convey("Given an topic with a wrong state value, then no transition is allowed", t, func() {
		topic := models.Topic{State: "wrong"}
		validateTransitionsToCreated(topic)
	})

	Convey("Given an topic without state, then created state is assumed when checking for transitions", t, func() {
		topic := models.Topic{}
		validateTransitionsToCreated(topic)
	})
}

func TestTopicReleaseDateValidation(t *testing.T) {
	t.Parallel()

	Convey("Given a valid topic release object", t, func() {
		topicRelease := models.TopicRelease{
			ReleaseDate: "2022-10-14T11:30:00Z",
		}
		releaseDate, err := topicRelease.Validate()
		So(err, ShouldBeNil)
		So(*releaseDate, ShouldHaveSameTypeAs, time.Now())
	})

	Convey("Given topic release object is empty", t, func() {
		topicRelease := models.TopicRelease{}
		releaseDate, err := topicRelease.Validate()
		So(err, ShouldEqual, apierrors.ErrInvalidReleaseDate)
		So(releaseDate, ShouldBeNil)
	})

	Convey("Given topic release object has empty release date value", t, func() {
		topicRelease := models.TopicRelease{
			ReleaseDate: "",
		}
		releaseDate, err := topicRelease.Validate()
		So(err, ShouldEqual, apierrors.ErrInvalidReleaseDate)
		So(releaseDate, ShouldBeNil)
	})

	Convey("Given topic release object has a non RFC3339 format date value", t, func() {
		Convey("Where the release date is missing timezone location notation", func() {
			topicRelease := models.TopicRelease{
				ReleaseDate: "2022-10-14T11:30:00",
			}
			releaseDate, err := topicRelease.Validate()
			So(err, ShouldEqual, apierrors.ErrInvalidReleaseDate)
			So(releaseDate, ShouldBeNil)
		})

		Convey("Where the release date is not in the correct structure", func() {
			topicRelease := models.TopicRelease{
				ReleaseDate: "10-10-2022T14.09.10Z",
			}
			releaseDate, err := topicRelease.Validate()
			So(err, ShouldEqual, apierrors.ErrInvalidReleaseDate)
			So(releaseDate, ShouldBeNil)
		})
	})
}

// validateTransitionsToCreated validates that the provided topic can transition to created state,
// and not to any forbidden of invalid state
func validateTransitionsToCreated(topic models.Topic) {
	Convey("Then an allowed transition is successfully checked", func() {
		So(topic.StateTransitionAllowed(models.StateCompleted.String()), ShouldBeTrue)
	})
	Convey("Then a forbidden transition to a valid state is not allowed", func() {
		So(topic.StateTransitionAllowed(models.StatePublished.String()), ShouldBeFalse)
	})
	Convey("Then a transition to an invalid state is not allowed", func() {
		So(topic.StateTransitionAllowed("wrong"), ShouldBeFalse)
	})
}
