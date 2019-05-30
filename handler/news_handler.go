package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yogaagungk/newsupdate/message"
	"github.com/yogaagungk/newsupdate/model"
	"github.com/yogaagungk/newsupdate/news"
)

//handler pointer
type handler struct {
	messaging message.Messaging
	service   news.Service
}

//InitDependencyNewsHandler digunakan untuk menginject dependency yang dibutuhkan
//pada handler endpoint
//return handler
func InitDependencyNewsHandler(messaging message.Messaging, service news.Service) *handler {
	return &handler{messaging, service}
}

// FindAllNews, sebagai handler untuk endpoint GET data News
func (handler *handler) FindAllNews(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))

	news, err := handler.service.FetchAll(page)

	if err != nil {
		c.JSON(http.StatusNotFound, news)
	} else {
		c.JSON(http.StatusOK, news)
	}
}

// SaveNews, sebagai handler untuk endpoint POST data News
func (handler *handler) SaveNews(c *gin.Context) {
	var data model.News

	c.BindJSON(&data)

	handler.messaging.Publish(data)
}
