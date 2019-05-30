package config

import (
	"github.com/yogaagungk/newsupdate/common"
	"gopkg.in/olivere/elastic.v7"
)

// InitialClientElastic , konfigurasi dan open connection ke elasticseaarch
func InitialClientElastic() *elastic.Client {
	client, err := elastic.NewClient(elastic.SetURL("http://localhost:9200"))

	if err != nil {
		common.HandleError(err, "Error create new client for elastic")

		return nil
	}

	return client
}
