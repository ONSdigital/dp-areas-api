package models_test

import (
	"testing"

	"github.com/ONSdigital/dp-topic-api/models"
	. "github.com/smartystreets/goconvey/convey"
)

func TestStateValidation(t *testing.T) {
	//!!! these need rationalising when final code is done
	Convey("Given a Created State, then only transitions to uploaded and deleted are allowed", t, func() {
		So(models.StateCreated.TransitionAllowed(models.StateCreated), ShouldBeFalse)
		So(models.StateCreated.TransitionAllowed(models.StateUploaded), ShouldBeTrue)
		So(models.StateCreated.TransitionAllowed(models.StateImporting), ShouldBeFalse)
		So(models.StateCreated.TransitionAllowed(models.StateImported), ShouldBeFalse)
		So(models.StateCreated.TransitionAllowed(models.StatePublished), ShouldBeFalse)
		So(models.StateCreated.TransitionAllowed(models.StateCompleted), ShouldBeFalse)
		So(models.StateCreated.TransitionAllowed(models.StateDeleted), ShouldBeTrue)
		So(models.StateCreated.TransitionAllowed(models.StateFailedImport), ShouldBeFalse)
		So(models.StateCreated.TransitionAllowed(models.StateFailedPublish), ShouldBeFalse)
	})

	Convey("Given a Uploaded State, then only transitions to importing and deleted are allowed", t, func() {
		So(models.StateUploaded.TransitionAllowed(models.StateCreated), ShouldBeFalse)
		So(models.StateUploaded.TransitionAllowed(models.StateUploaded), ShouldBeFalse)
		So(models.StateUploaded.TransitionAllowed(models.StateImporting), ShouldBeTrue)
		So(models.StateUploaded.TransitionAllowed(models.StateImported), ShouldBeFalse)
		So(models.StateUploaded.TransitionAllowed(models.StatePublished), ShouldBeFalse)
		So(models.StateUploaded.TransitionAllowed(models.StateCompleted), ShouldBeFalse)
		So(models.StateUploaded.TransitionAllowed(models.StateDeleted), ShouldBeTrue)
		So(models.StateUploaded.TransitionAllowed(models.StateFailedImport), ShouldBeFalse)
		So(models.StateUploaded.TransitionAllowed(models.StateFailedPublish), ShouldBeFalse)
	})

	Convey("Given an Importing State, then only transitions to imported, failedImport and deleted are allowed", t, func() {
		So(models.StateImporting.TransitionAllowed(models.StateCreated), ShouldBeFalse)
		So(models.StateImporting.TransitionAllowed(models.StateUploaded), ShouldBeFalse)
		So(models.StateImporting.TransitionAllowed(models.StateImporting), ShouldBeFalse)
		So(models.StateImporting.TransitionAllowed(models.StateImported), ShouldBeTrue)
		So(models.StateImporting.TransitionAllowed(models.StatePublished), ShouldBeFalse)
		So(models.StateImporting.TransitionAllowed(models.StateCompleted), ShouldBeFalse)
		So(models.StateImporting.TransitionAllowed(models.StateDeleted), ShouldBeTrue)
		So(models.StateImporting.TransitionAllowed(models.StateFailedImport), ShouldBeTrue)
		So(models.StateImporting.TransitionAllowed(models.StateFailedPublish), ShouldBeFalse)
	})

	Convey("Given an Imported State, then only transitions to published and deleted are allowed", t, func() {
		So(models.StateImported.TransitionAllowed(models.StateCreated), ShouldBeFalse)
		So(models.StateImported.TransitionAllowed(models.StateUploaded), ShouldBeFalse)
		So(models.StateImported.TransitionAllowed(models.StateImporting), ShouldBeFalse)
		So(models.StateImported.TransitionAllowed(models.StateImported), ShouldBeFalse)
		So(models.StateImported.TransitionAllowed(models.StatePublished), ShouldBeTrue)
		So(models.StateImported.TransitionAllowed(models.StateCompleted), ShouldBeFalse)
		So(models.StateImported.TransitionAllowed(models.StateDeleted), ShouldBeTrue)
		So(models.StateImported.TransitionAllowed(models.StateFailedImport), ShouldBeFalse)
		So(models.StateImported.TransitionAllowed(models.StateFailedPublish), ShouldBeFalse)
	})

	Convey("Given a Published State, then only transitions to failedPublish, completed and deleted are allowed", t, func() {
		So(models.StatePublished.TransitionAllowed(models.StateCreated), ShouldBeFalse)
		So(models.StatePublished.TransitionAllowed(models.StateUploaded), ShouldBeFalse)
		So(models.StatePublished.TransitionAllowed(models.StateImporting), ShouldBeFalse)
		So(models.StatePublished.TransitionAllowed(models.StateImported), ShouldBeFalse)
		So(models.StatePublished.TransitionAllowed(models.StatePublished), ShouldBeFalse)
		So(models.StatePublished.TransitionAllowed(models.StateCompleted), ShouldBeTrue)
		So(models.StatePublished.TransitionAllowed(models.StateDeleted), ShouldBeTrue)
		So(models.StatePublished.TransitionAllowed(models.StateFailedImport), ShouldBeFalse)
		So(models.StatePublished.TransitionAllowed(models.StateFailedPublish), ShouldBeTrue)
	})

	Convey("Given a Completed State, then only transitions to deleted are allowed", t, func() {
		So(models.StateCompleted.TransitionAllowed(models.StateCreated), ShouldBeFalse)
		So(models.StateCompleted.TransitionAllowed(models.StateUploaded), ShouldBeFalse)
		So(models.StateCompleted.TransitionAllowed(models.StateImporting), ShouldBeFalse)
		So(models.StateCompleted.TransitionAllowed(models.StateImported), ShouldBeFalse)
		So(models.StateCompleted.TransitionAllowed(models.StatePublished), ShouldBeFalse)
		So(models.StateCompleted.TransitionAllowed(models.StateCompleted), ShouldBeFalse)
		So(models.StateCompleted.TransitionAllowed(models.StateDeleted), ShouldBeTrue)
		So(models.StateCompleted.TransitionAllowed(models.StateFailedImport), ShouldBeFalse)
		So(models.StateCompleted.TransitionAllowed(models.StateFailedPublish), ShouldBeFalse)
	})

	Convey("Given a Deleted State, then no transitions are allowed", t, func() {
		So(models.StateDeleted.TransitionAllowed(models.StateCreated), ShouldBeFalse)
		So(models.StateDeleted.TransitionAllowed(models.StateUploaded), ShouldBeFalse)
		So(models.StateDeleted.TransitionAllowed(models.StateImporting), ShouldBeFalse)
		So(models.StateDeleted.TransitionAllowed(models.StateImported), ShouldBeFalse)
		So(models.StateDeleted.TransitionAllowed(models.StatePublished), ShouldBeFalse)
		So(models.StateDeleted.TransitionAllowed(models.StateCompleted), ShouldBeFalse)
		So(models.StateDeleted.TransitionAllowed(models.StateDeleted), ShouldBeFalse)
		So(models.StateDeleted.TransitionAllowed(models.StateFailedImport), ShouldBeFalse)
		So(models.StateDeleted.TransitionAllowed(models.StateFailedPublish), ShouldBeFalse)
	})

	Convey("Given a FailedImport State, then only transitions to deleted are allowed", t, func() {
		So(models.StateFailedImport.TransitionAllowed(models.StateCreated), ShouldBeFalse)
		So(models.StateFailedImport.TransitionAllowed(models.StateUploaded), ShouldBeFalse)
		So(models.StateFailedImport.TransitionAllowed(models.StateImporting), ShouldBeFalse)
		So(models.StateFailedImport.TransitionAllowed(models.StateImported), ShouldBeFalse)
		So(models.StateFailedImport.TransitionAllowed(models.StatePublished), ShouldBeFalse)
		So(models.StateFailedImport.TransitionAllowed(models.StateCompleted), ShouldBeFalse)
		So(models.StateFailedImport.TransitionAllowed(models.StateDeleted), ShouldBeTrue)
		So(models.StateFailedImport.TransitionAllowed(models.StateFailedImport), ShouldBeFalse)
		So(models.StateFailedImport.TransitionAllowed(models.StateFailedPublish), ShouldBeFalse)
	})

	Convey("Given a FailedPublish State, then only transitions to deleted are allowed", t, func() {
		So(models.StateFailedPublish.TransitionAllowed(models.StateCreated), ShouldBeFalse)
		So(models.StateFailedPublish.TransitionAllowed(models.StateUploaded), ShouldBeFalse)
		So(models.StateFailedPublish.TransitionAllowed(models.StateImporting), ShouldBeFalse)
		So(models.StateFailedPublish.TransitionAllowed(models.StateImported), ShouldBeFalse)
		So(models.StateFailedPublish.TransitionAllowed(models.StatePublished), ShouldBeFalse)
		So(models.StateFailedPublish.TransitionAllowed(models.StateCompleted), ShouldBeFalse)
		So(models.StateFailedPublish.TransitionAllowed(models.StateDeleted), ShouldBeTrue)
		So(models.StateFailedPublish.TransitionAllowed(models.StateFailedImport), ShouldBeFalse)
		So(models.StateFailedPublish.TransitionAllowed(models.StateFailedPublish), ShouldBeFalse)
	})
}
