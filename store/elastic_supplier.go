package store

import (
	"context"
	"fmt"
	"github.com/olivere/elastic"
	"github.com/webitel/storage/model"
	"github.com/webitel/wlog"
	"log"
	"os"
	"time"
)

const (
	DB_PING_ATTEMPTS     = 18
	DB_PING_TIMEOUT_SECS = 10
)

const (
	EXIT_DB_OPEN = 101
	EXIT_PING    = 102
)

type ElasticSupplier struct {
	client *elastic.Client
}

func (e ElasticSupplier) Name() string {
	return "Elastic"
}

func NewElasticSupplier(settings model.NoSqlSettings) *ElasticSupplier {
	var err error

	if settings.Host == nil {
		wlog.Critical("Failed to open NoSQL connection to err: bad settings host")
		time.Sleep(time.Second)
		os.Exit(EXIT_DB_OPEN)
	}

	supplier := &ElasticSupplier{}
	options := []elastic.ClientOptionFunc{
		elastic.SetURL(*settings.Host),
		elastic.SetSniff(false),
		elastic.SetHealthcheckTimeout(time.Second * DB_PING_TIMEOUT_SECS),
	}

	if settings.Trace {
		options = append(options, elastic.SetTraceLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)))
	}

	supplier.client, err = elastic.NewClient(options...)

	if err != nil {
		return nil
		wlog.Critical(fmt.Sprintf("Failed to open NoSQL connection to err:%v", err.Error()))
		time.Sleep(time.Second)
		os.Exit(EXIT_DB_OPEN)
	}

	for i := 0; i < DB_PING_ATTEMPTS; i++ {
		wlog.Info(fmt.Sprintf("Pinging NoSQL %v database", supplier.Name()))
		ctx, cancel := context.WithTimeout(context.Background(), DB_PING_TIMEOUT_SECS*time.Second)
		defer cancel()

		_, _, err := supplier.client.Ping(*settings.Host).Do(ctx)

		if err == nil {
			break
		} else {
			if i == DB_PING_ATTEMPTS-1 {
				wlog.Critical(fmt.Sprintf("Failed to ping DB, server will exit err=%v", err))
				time.Sleep(time.Second)
				os.Exit(EXIT_PING)
			} else {
				wlog.Error(fmt.Sprintf("Failed to ping DB retrying in %v seconds err=%v", DB_PING_TIMEOUT_SECS, err))
				time.Sleep(DB_PING_TIMEOUT_SECS * time.Second)
			}
		}
	}

	return supplier
}
