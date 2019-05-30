package news

import (
	"github.com/yogaagungk/newsupdate/model"
)

//Service digunakan sebagai contract
type Service interface {
	Save(data *model.News) (model.News, error)
	FetchAll(page int) ([]model.News, error)
}
