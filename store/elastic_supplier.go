package store

import "github.com/olivere/elastic"

type ElasticSupplier struct {
	client *elastic.Client
}

func NewElasticSupplier() *ElasticSupplier {
	return &ElasticSupplier{}
}
