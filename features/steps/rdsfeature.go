package steps

import (
	"context"
	"github.com/ONSdigital/dp-areas-api/config"
	"github.com/ONSdigital/dp-areas-api/rds"
	componenttest "github.com/ONSdigital/dp-component-test"
)

type RDSFeature struct {
	componenttest.ErrorFeature
	client rds.RDS
	cfg    *config.Config
}

func NewRDSFeature(ef componenttest.ErrorFeature, cfg *config.Config) *RDSFeature {
	rf := &RDSFeature{
		ErrorFeature: ef,
	}

	if err := rf.client.Init(context.Background(), cfg); err != nil {
		return nil
	}

	return rf
}

func (rf *RDSFeature) Reset() {}

func (rf *RDSFeature) Close() { rf.client.Close() }

func (r RDSFeature) insertIntoDatabase() {

}
