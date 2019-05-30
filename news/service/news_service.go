package news

import (
	"encoding/json"
	"time"

	"github.com/gomodule/redigo/redis"

	"github.com/jinzhu/gorm"
	"github.com/yogaagungk/newsupdate/common"
	"github.com/yogaagungk/newsupdate/elasticsearch"
	"github.com/yogaagungk/newsupdate/model"
	"github.com/yogaagungk/newsupdate/news"
)

type repository struct {
	db          *gorm.DB
	newselastic elasticsearch.NewsElasticsearch
	redisConn   redis.Conn
}

//InitDependencyNewsService digunakan untuk menginject dependency yang dibutuhkan
//pada news service
func InitDependencyNewsService(db *gorm.DB,
	newselastic elasticsearch.NewsElasticsearch,
	redisConn redis.Conn) news.Service {

	return &repository{db, newselastic, redisConn}
}

//Save berisi operasi untuk menyimpan data kedalam database dan elasticsearch
func (repo *repository) Save(data *model.News) (model.News, error) {
	entity := model.News{Body: data.Body, Author: data.Author, Created: time.Now()}

	err := repo.db.Create(&entity)

	if err != nil {
		common.HandleError(err.Error, "Failed save data to database")
	}

	repo.newselastic.InsertNews(entity)

	return entity, nil
}

//FetchAll berisi operasi untuk mengambil data news (id, created) dari elastic
//kemudian dilakukan concurrency menggunakan goroutine untuk mengambil seluruh data news
//di dalam goroutine mengecek terlebih dahulu di cache, apabila ada datanya, maka nilai dari cache akan dikembalikan
//namun apabila tidak ada, maka akan mengambil data news dari database, dan kemudian menyimpannya kedalam cache
func (repo *repository) FetchAll(page int) ([]model.News, error) {
	var results []model.News

	newsChannel := make(chan model.News, 10)

	newsFromElastics := repo.newselastic.FindAll(page)

	go func() {
		for _, newsFromElastic := range newsFromElastics {
			var result model.News

			value, err := redis.String(repo.redisConn.Do("GET", newsFromElastic.ID))

			if err == redis.ErrNil {
				repo.db.First(&result, newsFromElastic.ID)

				json, err := json.Marshal(result)

				if err != nil {
					common.HandleError(err, "Error Encoding JSON")
				}

				repo.redisConn.Do("SET", result.ID, json)
			} else {
				rslt := model.News{}

				json.Unmarshal([]byte(value), &rslt)

				result = rslt
			}

			newsChannel <- result
		}

		close(newsChannel)
	}()

	for news := range newsChannel {
		results = append(results, news)
	}

	return results, nil
}
