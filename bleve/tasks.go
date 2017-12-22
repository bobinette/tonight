package bleve

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/blevesearch/bleve"
	_ "github.com/blevesearch/bleve/analysis/analyzer/keyword"
	_ "github.com/blevesearch/bleve/analysis/analyzer/simple"
	"github.com/blevesearch/bleve/mapping"
	"github.com/blevesearch/bleve/search/query"

	"github.com/bobinette/tonight"
)

type Index struct {
	index bleve.Index
}

func (s *Index) Open(path string) error {
	index, err := bleve.Open(path)
	if err != nil {
		if err != bleve.ErrorIndexPathDoesNotExist {
			return err
		}

		data, err := ioutil.ReadFile("bleve/mapping.json")
		if err != nil {
			return err
		}

		var m mapping.IndexMappingImpl
		err = json.Unmarshal(data, &m)
		if err != nil {
			return err
		}

		index, err = bleve.New(path, &m)
		if err != nil {
			return err
		}
	}

	s.index = index
	return nil
}

func (s *Index) Close() error {
	if s.index == nil {
		return nil
	}

	return s.index.Close()
}

func (s *Index) Index(ctx context.Context, task tonight.Task) error {
	data := map[string]interface{}{
		"id":          task.ID,
		"title":       task.Title,
		"description": task.Description,
		"tags":        task.Tags,
		"done":        strconv.Itoa(int(task.Done())),
		"rank":        task.Rank,
	}

	if doneAt := task.DoneAt(); doneAt != nil {
		data["done_at"] = *doneAt
	}

	return s.index.Index(fmt.Sprintf("%d", task.ID), data)
}

func (s *Index) Delete(ctx context.Context, id uint) error {
	return s.index.Delete(fmt.Sprintf("%d", id))
}

func (s *Index) Search(ctx context.Context, p tonight.TaskSearchParameters) ([]uint, error) {
	total := 100 // Default...
	if sr, err := s.index.Search(bleve.NewSearchRequest(query.NewMatchAllQuery())); err != nil {
		return nil, err
	} else {
		total = int(sr.Total)
	}

	query := andQ(
		query.NewMatchAllQuery(),
		s.searchQ(p.Q),
		searchDoneStatuses(p.Statuses),
		searchIDs(p.IDs),
	)

	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.Size = total
	searchRequest.SortBy([]string{"rank"})
	if len(p.Statuses) != 1 {
		searchRequest.SortBy([]string{"-done_at"})
	}

	// Activate for debugging
	// searchRequest.Highlight = bleve.NewHighlight()

	searchResults, err := s.index.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	ids := make([]uint, len(searchResults.Hits))
	for i, hit := range searchResults.Hits {
		id64, err := strconv.ParseUint(hit.ID, 10, 64)
		if err != nil {
			return nil, err
		}
		ids[i] = uint(id64)

		// Activate for debugging
		// fmt.Printf("%s: %f\n", hit.ID, hit.Score)
		// fmt.Println(hit.Fragments)
		// fmt.Println(hit.Expl)
	}

	return ids, nil
}

func andQ(qs ...query.Query) query.Query {
	ands := make([]query.Query, 0, len(qs))
	for _, q := range qs {
		if q != nil {
			ands = append(ands, q)
		}
	}

	if len(ands) == 0 {
		return nil
	}
	return query.NewConjunctionQuery(ands)
}

func orQ(qs ...query.Query) query.Query {
	ors := make([]query.Query, 0, len(qs))
	for _, q := range qs {
		if q != nil {
			ors = append(ors, q)
		}
	}

	if len(ors) == 0 {
		return nil
	}
	return query.NewDisjunctionQuery(ors)
}

func (s *Index) searchQ(queryString string) query.Query {
	words := strings.Fields(queryString)
	if len(words) == 0 {
		return nil
	}

	ands := make([]query.Query, 0, len(words))
	for _, word := range words {
		var q query.Query
		if strings.HasPrefix(word, "#") && len(word) > 1 {
			fmt.Println("tag", word)
			q = s.matches(word[1:], "tags")
		} else {
			fmt.Println("others", word)
			q = orQ(
				s.matches(word, "title"),
				s.matches(word, "description"),
			)
		}

		ands = append(ands, q)
	}

	return andQ(ands...)
}

func (s *Index) matches(queryString, field string) query.Query {
	analyzer := s.index.Mapping().AnalyzerNamed(s.index.Mapping().AnalyzerNameForPath(field))
	tokens := analyzer.Analyze([]byte(queryString))
	if len(tokens) == 0 {
		return nil
	}

	conjuncts := make([]query.Query, len(tokens))
	for i, token := range tokens {
		conjuncts[i] = &query.PrefixQuery{
			Prefix:   string(token.Term),
			FieldVal: field,
		}
	}

	return query.NewConjunctionQuery(conjuncts)
}

func searchDoneStatuses(statuses []tonight.DoneStatus) query.Query {
	if len(statuses) == 0 {
		return nil
	}

	qs := make([]query.Query, len(statuses))
	for i, s := range statuses {
		query := bleve.NewTermQuery(strconv.Itoa(int(s)))
		query.FieldVal = "done"
		qs[i] = query
	}
	return orQ(qs...)
}

func searchIDs(ids []uint) query.Query {
	docIDs := make([]string, len(ids))
	for i, id := range ids {
		docIDs[i] = fmt.Sprintf("%d", id)
	}
	return query.NewDocIDQuery(docIDs)
}
