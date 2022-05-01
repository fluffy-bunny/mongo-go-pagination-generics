package mongopagination

import (
	"context"

	base "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

// PagingQuery is an interface that provides list of function
// you can perform on pagingQuery
type PagingQuery[T any] interface {
	// Find set the filter for query results.
	Find() (records []T, paginatedData *base.PaginatedData, err error)

	Aggregate(criteria ...interface{}) (records []T, paginatedData *base.PaginatedData, err error)

	// Select used to enable fields which should be retrieved.
	Select(selector interface{}) PagingQuery[T]
	Filter(selector interface{}) PagingQuery[T]
	Limit(limit int64) PagingQuery[T]
	Page(page int64) PagingQuery[T]
	Sort(sortField string, sortValue interface{}) PagingQuery[T]
	Context(ctx context.Context) PagingQuery[T]
	SetCollation(ctx *options.Collation) PagingQuery[T]
}

// PagingQuery struct for holding mongo
// connection, filter needed to apply
// filter data with page, limit, sort key
// and sort value
type pagingQuery[T any] struct {
	internal base.PagingQuery
}

// New is to construct PagingQuery object with mongo.Database and collection name
func New[T any](collection *mongo.Collection) PagingQuery[T] {
	return &pagingQuery[T]{
		internal: base.New(collection),
	}
}
func (paging *pagingQuery[T]) SetCollation(ctx *options.Collation) PagingQuery[T] {
	paging.internal.SetCollation(ctx)
	return paging
}
func (paging *pagingQuery[T]) _decode(decode *[]T) PagingQuery[T] {
	paging.internal.Decode(decode)
	return paging
}
func (paging *pagingQuery[T]) Context(ctx context.Context) PagingQuery[T] {
	paging.internal.Context(ctx)
	return paging
}
func (paging *pagingQuery[T]) Sort(sortField string, sortValue interface{}) PagingQuery[T] {
	paging.internal.Sort(sortField, sortValue)
	return paging
}
func (paging *pagingQuery[T]) Page(page int64) PagingQuery[T] {
	paging.internal.Page(page)
	return paging
}
func (paging *pagingQuery[T]) Limit(limit int64) PagingQuery[T] {
	paging.internal.Limit(limit)
	return paging
}
func (paging *pagingQuery[T]) Filter(selector interface{}) PagingQuery[T] {
	paging.internal.Filter(selector)
	return paging
}
func (paging *pagingQuery[T]) Select(selector interface{}) PagingQuery[T] {
	paging.internal.Select(selector)
	return paging
}
func (paging *pagingQuery[T]) Aggregate(criteria ...interface{}) (records []T, paginatedData *base.PaginatedData, err error) {
	paginatedData, err = paging.internal.Aggregate(criteria...)
	if err != nil {
		return
	}
	var aggList []T
	for _, raw := range paginatedData.Data {
		var record *T
		if marshallErr := bson.Unmarshal(raw, &record); marshallErr == nil {
			aggList = append(aggList, *record)
		}
	}
	return
}

func (paging *pagingQuery[T]) Find() (records []T, paginatedData *base.PaginatedData, err error) {
	paging._decode(&records)
	paginatedData, err = paging.internal.Find()
	return
}
