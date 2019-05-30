package elasticsearch

import (
	"context"
	"encoding/json"
	"time"

	"github.com/yogaagungk/newsupdate/common"

	"github.com/yogaagungk/newsupdate/model"

	"gopkg.in/olivere/elastic.v7"
)

//NewsElasticsearch contarct
type NewsElasticsearch interface {
	InsertNews(data model.News)
	FindAll(page int) []News
}

//News struct digunakan untuk menyimpan data di elasticsearch
type News struct {
	ID      uint      `json:"id"`
	Created time.Time `json:"created"`
}

type newselasticsearch struct {
	context       context.Context
	elasticClient *elastic.Client
}

//InitDependencyElasticsearch digunakan untuk menginject dependency yang dibutuhkan
func InitDependencyElasticsearch(context context.Context, elasticClient *elastic.Client) NewsElasticsearch {
	return &newselasticsearch{context, elasticClient}
}

//InsertNews berisi operasi untuk insert news ke dalam elasticsearch
//data news yang di save adalah id dan created
func (newselasticsearch *newselasticsearch) InsertNews(data model.News) {
	news := News{data.ID, data.Created}

	_, err := newselasticsearch.elasticClient.Index().Index("newsupdate").Type("news").BodyJson(news).Do(newselasticsearch.context)

	if err != nil {
		common.HandleError(err, "Could not insert data to elasticsearch")
	}

	newselasticsearch.elasticClient.Flush().Index("newsupdate").Do(newselasticsearch.context)
}

//FindAll berisi operasi untuk mengambil data news dari elasticsearch
//Data di sort berdasarkan created DESC dan masing-masing result berisi 10 data
func (newselasticsearch *newselasticsearch) FindAll(page int) []News {
	query := elastic.MatchAllQuery{}

	from := page * 10

	result, err := newselasticsearch.elasticClient.Search().
		Index("newsupdate").
		Query(query).
		SortBy(elastic.NewFieldSort("created").Desc()).
		From(from - 10).
		Size(10).
		Do(newselasticsearch.context)

	if err != nil {
		common.HandleError(err, "Error when fetching result from query")
	}

	return convertResultToNews(result)
}

//Convert from searchResult elastic to slices News struct
func convertResultToNews(result *elastic.SearchResult) []News {
	var rslt []News

	for _, hit := range result.Hits.Hits {
		var newsObject News

		err := json.Unmarshal(hit.Source, &newsObject)

		if err != nil {
			common.HandleError(err, "Error while deserialize JSON")

			continue
		}

		rslt = append(rslt, newsObject)
	}

	return rslt
}
