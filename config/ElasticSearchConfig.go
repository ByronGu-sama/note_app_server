package config

import (
	"context"
	"github.com/elastic/go-elasticsearch/v8"
	"note_app_server/global"
)

func InitElasticSearchConfig() {
	url := "http://" + AC.Es.Host + ":" + AC.Es.Port
	client, err := elasticsearch.NewTypedClient(elasticsearch.Config{
		Addresses: []string{url},
	})
	if err != nil {
		panic(err)
	}
	ctx := context.TODO()
	_, err = client.Ping().IsSuccess(ctx)
	if err != nil {
		panic(err)
	}
	global.ESClient = client
}
