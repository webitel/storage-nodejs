package store

import (
	"context"
	"github.com/olivere/elastic"
	"github.com/webitel/storage/model"
	"net/http"
)

type SearchEngine struct {
	DatabaseLayer LayeredStoreDatabaseLayer
	ElasticLayer  *ElasticSupplier
}

func (self *SearchEngine) Search(request *model.SearchEngineRequest) StoreChannel {

	return Do(func(result *StoreResult) {

		qs := elastic.NewQueryStringQuery(request.Query)
		qs.AnalyzeWildcard(true)

		filter := elastic.NewBoolQuery()
		filter.Must(request.Filter, qs)

		fetchContext := elastic.NewFetchSourceContext(true).
			Include(request.Includes...).
			Exclude(request.Excludes...)

		fetchContext.SetFetchSource(true)

		searchSource := elastic.NewSearchSource().
			StoredFields(request.Columns...).
			DocvalueFields(request.Columns...).
			FetchSource(true).
			FetchSourceContext(fetchContext)

		var err error
		var res *elastic.SearchResult

		if request.ScrollKeepAlive != nil && *request.ScrollKeepAlive != "" {
			res, err = self.ElasticLayer.client.
				Scroll(request.Index).
				KeepAlive(*request.ScrollKeepAlive).
				Type(request.Type).
				AllowNoIndices(true).
				IgnoreUnavailable(true).
				SearchSource(searchSource).
				Size(request.Size).
				SortBy(request.Sort).
				Query(filter).
				Do(context.TODO())

		} else {
			res, err = self.ElasticLayer.client.
				Search(request.Index).
				Type(request.Type).
				AllowNoIndices(true).
				IgnoreUnavailable(true).
				SearchSource(searchSource).
				Size(request.Size).
				From(request.Size * request.Page).
				SortBy(request.Sort).
				Query(filter).
				Do(context.TODO())
		}

		if err != nil {
			result.Err = model.NewAppError("CdrSupplier.Search", "store.no_sql.search.app_error",
				map[string]interface{}{"Type": request.Type}, err.Error(), http.StatusInternalServerError)
			return
		}

		result.Data = toResponse(res)
	})
}

func (self *SearchEngine) Scroll(request *model.SearchEngineScroll) StoreChannel {
	return Do(func(result *StoreResult) {
		res, err := self.ElasticLayer.client.
			Scroll().
			ScrollId(request.ScrollId).
			Scroll(request.ScrollKeepAlive).
			Do(context.TODO())

		if err != nil {
			result.Err = model.NewAppError("CdrSupplier.Scroll", "store.no_sql.scroll.app_error",
				map[string]interface{}{"ScrollId": request.ScrollId}, err.Error(), http.StatusInternalServerError)
			return
		}

		result.Data = toResponse(res)
	})
}

func toResponse(src *elastic.SearchResult) *model.SearchEngineResponse {
	res := new(model.SearchEngineResponse)
	res.Timeout = src.TimedOut
	res.Shards = map[string]int{
		"failed":     src.Shards.Failed,
		"successful": src.Shards.Successful,
		"total":      src.Shards.Total,
	}

	if src.ScrollId != "" {
		res.ScrollId = model.NewString(src.ScrollId)
	}

	res.Hits = model.BaseHitsResponseSearchEngine{
		Total: src.Hits.TotalHits,
		Hits:  make([]*model.SearchEngineHitsResponse, 0, len(src.Hits.Hits)),
	}

	for _, h := range src.Hits.Hits {
		res.Hits.Hits = append(res.Hits.Hits, &model.SearchEngineHitsResponse{
			Id:     h.Id,
			Index:  h.Index,
			Source: h.Source,
			Fields: h.Fields,
		})
	}

	return res
}
