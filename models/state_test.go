package models_test

import (
	"testing"

	"github.com/ONSdigital/dp-topic-api/models"
	. "github.com/smartystreets/goconvey/convey"
)

func TestStateValidation(t *testing.T) {
	// TODO these need rationalising when final code is done
	Convey("Given a Created State, then only transitions to completed allowed", t, func() {
		So(models.StateCreated.TransitionAllowed(models.StateCreated), ShouldBeFalse)
		So(models.StateCreated.TransitionAllowed(models.StatePublished), ShouldBeFalse)
		So(models.StateCreated.TransitionAllowed(models.StateCompleted), ShouldBeTrue)
	})

	Convey("Given a Published State, then only transitions to created allowed", t, func() {
		So(models.StatePublished.TransitionAllowed(models.StateCreated), ShouldBeTrue)
		So(models.StatePublished.TransitionAllowed(models.StatePublished), ShouldBeFalse)
		So(models.StatePublished.TransitionAllowed(models.StateCompleted), ShouldBeFalse)
	})

	Convey("Given a Completed State, then only transitions to published allowed", t, func() {
		So(models.StateCompleted.TransitionAllowed(models.StateCreated), ShouldBeFalse)
		So(models.StateCompleted.TransitionAllowed(models.StatePublished), ShouldBeTrue)
		So(models.StateCompleted.TransitionAllowed(models.StateCompleted), ShouldBeFalse)
	})
}
