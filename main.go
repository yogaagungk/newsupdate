package main

import (
	"context"

	"github.com/yogaagungk/newsupdate/elasticsearch"

	"github.com/yogaagungk/newsupdate/common"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/yogaagungk/newsupdate/config"
	"github.com/yogaagungk/newsupdate/handler"
	"github.com/yogaagungk/newsupdate/message"
	news "github.com/yogaagungk/newsupdate/news/service"
	"gopkg.in/olivere/elastic.v7"
)

func main() {
	db := config.OpenDB()
	connection, amqpChannel, _ := config.InitialChannel()
	elasticClient := config.InitialClientElastic()
	redisClient := config.InitialRedisConn()

	context := context.Background()

	createIndexIfNotExist(context, elasticClient, "newsupdate")

	defer func() {
		db.Close()
		connection.Close()
		amqpChannel.Close()
		redisClient.Close()
	}()

	newsElasticsearch := elasticsearch.InitDependencyElasticsearch(context, elasticClient)
	newsService := news.InitDependencyNewsService(db, newsElasticsearch, redisClient)
	consumer := message.InitDependencyMessaging(amqpChannel, newsService)
	newsHandler := handler.InitDependencyNewsHandler(consumer, newsService)

	stopChan := make(chan bool)

	consumer.Listen()

	r := gin.Default()
	r.GET("/news", newsHandler.FindAllNews)
	r.POST("/news", newsHandler.SaveNews)
	r.Run(":5045")

	<-stopChan
}

//Check index in elastic, if not exist then create one, if exist do nothing
func createIndexIfNotExist(context context.Context, client *elastic.Client, indexName string) error {
	isExists, err := client.IndexExists(indexName).Do(context)

	if err != nil {
		common.HandleError(err, "Error while checking index is exist or not")

		return err
	}

	if isExists {
		return nil
	}

	newIndex, err := client.CreateIndex(indexName).Do(context)

	if err != nil {
		common.HandleError(err, "Error while create new index")

		return err
	}

	if !newIndex.Acknowledged {
		common.HandleError(err, "CreateIndex was not acknowledged. Check that timeout value is correct")

		return err
	}

	return nil
}
