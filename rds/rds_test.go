package rds

import (
	"errors"
	"github.com/ONSdigital/dp-areas-api/models"
	pgxMock "github.com/ONSdigital/dp-areas-api/pgx/mock"
	"github.com/jackc/pgx/v4"
	"testing"
)

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRDS_GetArea(t *testing.T) {
	Convey("Given an valid area code", t, func() {

		rowMock := &pgxMock.PGXRowMock{
			ScanFunc: func(dest ...interface{}) error {
				id := dest[0].(*int64)
				code := dest[1].(*string)
				active := dest[2].(*bool)

				*id = 1
				*code = "Wales"
				*active = true
				return nil
			},
		}

		rds := RDS{
			conn: &pgxMock.PGXPoolMock{
				QueryRowFunc: func(ctx context.Context, sql string, args ...interface{}) pgx.Row {
					return rowMock
				},
			}}
		area, err := rds.GetArea("W92000004")

		Convey("When GetArea is invoked", func() {

			Convey("Then area details are returned", func() {
				So(err, ShouldBeNil)
				So(area.Code, ShouldEqual, "Wales")
				So(area.Id, ShouldEqual, 1)
				So(area.Active, ShouldEqual, true)
			})
		})
	})

	Convey("Given an invalid area code", t, func() {

		rowMock := &pgxMock.PGXRowMock{
			ScanFunc: func(dest ...interface{}) error {
				return errors.New("no rows in result set")
			},
		}

		rds := RDS{
			conn: &pgxMock.PGXPoolMock{
				QueryRowFunc: func(ctx context.Context, sql string, args ...interface{}) pgx.Row {
					return rowMock
				},
			}}
		area, err := rds.GetArea("123")

		Convey("When GetArea is invoked", func() {

			Convey("Then an error is returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "no rows in result set")
				So(area, ShouldBeNil)
			})
		})
	})
}

func TestRDS_ValidateArea(t *testing.T) {
	Convey("Given valid area code", t, func() {

		rowMock := &pgxMock.PGXRowMock{
			ScanFunc: func(dest ...interface{}) error {
				code := dest[0].(*string)

				*code = "Wales"
				return nil
			},
		}

		rds := RDS{
			conn: &pgxMock.PGXPoolMock{
				QueryRowFunc: func(ctx context.Context, sql string, args ...interface{}) pgx.Row {
					return rowMock
				},
			}}
		err := rds.ValidateArea("W92000004")

		Convey("When area code is validated", func() {

			Convey("Then nil is returned", func() {
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given an invalid area code", t, func() {

		rowMock := &pgxMock.PGXRowMock{
			ScanFunc: func(dest ...interface{}) error {
				return errors.New("no rows in result set")
			},
		}

		rds := RDS{
			conn: &pgxMock.PGXPoolMock{
				QueryRowFunc: func(ctx context.Context, sql string, args ...interface{}) pgx.Row {
					return rowMock
				},
			}}
		err := rds.ValidateArea("invalid")

		Convey("When invalid area  code is validated", func() {

			Convey("Then error is returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "no rows in result set")
			})
		})
	})
}

func TestRDS_GetRelationships(t *testing.T) {
	Convey("Given a valid area code with relationships", t, func() {
		callCount := 0

		relationships := []*models.AreaBasicData{
			{"E12000001", "North East"},
			{"E12000002", "North West"},
			{"E12000003", "Yorkshire and The Humbe"},
		}

		rowMock := &pgxMock.PGXRowsMock{
			CloseFunc: func() {
			},
			NextFunc: func() bool {
				response := callCount < len(relationships)
				return response
			},
			ScanFunc: func(dest ...interface{}) error {
				code := dest[0].(*string)
				name := dest[1].(*string)

				*code = relationships[callCount].Code
				*name = relationships[callCount].Name

				callCount = callCount + 1
				return nil
			},
		}

		rds := RDS{
			conn: &pgxMock.PGXPoolMock{
				QueryFunc: func(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
					return rowMock, nil
				},
			}}
		actualRelationships, err := rds.GetRelationships("E92000001")

		Convey("When relationships are fetched", func() {

			Convey("Then all relationships available for the area code is returned", func() {
				So(err, ShouldBeNil)
				So(actualRelationships, ShouldResemble, relationships)
			})
		})
	})

	Convey("Given an valid area code", t, func() {
		errorMsg := "error while connecting to DB"
		rds := RDS{
			conn: &pgxMock.PGXPoolMock{
				QueryFunc: func(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
					return nil, errors.New(errorMsg)
				},
			}}
		actualRelationships, err := rds.GetRelationships("E92000001")

		Convey("When failed to connect to DB", func() {

			Convey("Then error is returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, errorMsg)
				So(actualRelationships, ShouldBeNil)
			})
		})
	})

	Convey("Given an invalid area code", t, func() {
		rowMock := &pgxMock.PGXRowsMock{
			CloseFunc: func() {
			},
			NextFunc: func() bool {
				return false
			},
		}

		rds := RDS{
			conn: &pgxMock.PGXPoolMock{
				QueryFunc: func(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
					return rowMock, nil
				},
			}}
		actualRelationships, err := rds.GetRelationships("E92000001")

		Convey("When relationships are fetched", func() {

			Convey("Then empty relationships are returned", func() {
				So(err, ShouldBeNil)
				So(actualRelationships, ShouldBeEmpty)
			})
		})
	})
}
