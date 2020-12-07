package models_test

import (
	"testing"

	"github.com/ONSdigital/dp-topic-api/models"
	. "github.com/smartystreets/goconvey/convey"
)

func TestStateValidation(t *testing.T) {
	//!!! these need rationalising when final code is done
	Convey("Given a Created State, then only transitions to completed allowed", t, func() {
		So(models.StateTopicCreated.TransitionAllowed(models.StateTopicCreated), ShouldBeFalse)
		So(models.StateTopicCreated.TransitionAllowed(models.StateTopicPublished), ShouldBeFalse)
		So(models.StateTopicCreated.TransitionAllowed(models.StateTopicCompleted), ShouldBeTrue)
	})

	Convey("Given a Published State, then only transitions to created allowed", t, func() {
		So(models.StateTopicPublished.TransitionAllowed(models.StateTopicCreated), ShouldBeTrue)
		So(models.StateTopicPublished.TransitionAllowed(models.StateTopicPublished), ShouldBeFalse)
		So(models.StateTopicPublished.TransitionAllowed(models.StateTopicCompleted), ShouldBeFalse)
	})

	Convey("Given a Completed State, then only transitions to published allowed", t, func() {
		So(models.StateTopicCompleted.TransitionAllowed(models.StateTopicCreated), ShouldBeFalse)
		So(models.StateTopicCompleted.TransitionAllowed(models.StateTopicPublished), ShouldBeTrue)
		So(models.StateTopicCompleted.TransitionAllowed(models.StateTopicCompleted), ShouldBeFalse)
	})
}
