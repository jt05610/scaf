package couch

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/jt05610/scaf/core"
	"io"
	"log"
	"net/http"
)

type Couch[T core.Storable] struct {
	url string
}

type CreateResponse struct {
	Ok  bool   `json:"ok"`
	Id  string `json:"id"`
	Rev string `json:"rev"`
}

func (c *Couch[T]) idURL(t T) string {
	return c.url + "/" + t.Meta().ID
}

func (c *Couch[T]) Create(t T) error {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(t); err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPut, c.idURL(t), &buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	var cr CreateResponse
	if err := json.NewDecoder(resp.Body).Decode(&cr); err != nil {
		return err
	}
	if !cr.Ok {
		return errors.New("failed to create document")
	}
	t.Meta().Rev = cr.Rev
	return nil
}

func (c *Couch[T]) Update(t T) error {
	if t.Meta().Rev == "" {
		return errors.New("missing revision, use Create instead")
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(t); err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPut, c.idURL(t), &buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	var cr CreateResponse
	if err := json.NewDecoder(resp.Body).Decode(&cr); err != nil {
		return err
	}
	t.Meta().Rev = cr.Rev
	if !cr.Ok {
		return errors.New("failed to update document")
	}
	return nil
}

func (c *Couch[T]) Delete(t T) error {
	if t.Meta().Rev == "" {
		return errors.New("cannot delete without revision")
	}
	req, err := http.NewRequest(http.MethodDelete, c.idURL(t)+"?rev="+t.Meta().Rev, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	var cr CreateResponse
	if err := json.NewDecoder(resp.Body).Decode(&cr); err != nil {
		return err
	}
	t.Meta().Rev = cr.Rev
	if !cr.Ok {
		return errors.New("failed to delete document")
	}
	return nil
}

type DocQuery struct {
	Docs []struct {
		Id  string `json:"id"`
		Rev string `json:"rev"`
	} `json:"docs"`
}

type BulkGetResult[T any] struct {
	Results []struct {
		Id   string `json:"id"`
		Docs []struct {
			Ok    T `json:"ok"`
			Error struct {
				ID     string `json:"id"`
				Rev    string `json:"rev"`
				Error  string `json:"error"`
				Reason string `json:"reason"`
			} `json:"error"`
		} `json:"docs"`
	}
}

func (r *BulkGetResult[T]) Load(rdr io.Reader) ([]T, error) {
	if err := json.NewDecoder(rdr).Decode(r); err != nil {
		return nil, err
	}
	var ts []T
	for _, result := range r.Results {
		for _, doc := range result.Docs {
			if doc.Error.Error != "" {
				return nil, errors.New(doc.Error.Reason)
			}
			ts = append(ts, doc.Ok)
		}
	}
	return ts, nil
}

type AllDocResult struct {
	Offset int `json:"offset"`
	Rows   []struct {
		Id    string `json:"id"`
		Key   string `json:"key"`
		Value struct {
			Rev string `json:"rev"`
		}
	}
	TotalRows int `json:"total_rows"`
}

func (c *Couch[T]) allDocs() (*AllDocResult, error) {
	res, err := http.Get(c.url + "/_all_docs")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = res.Body.Close()
	}()
	var result AllDocResult
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func listQuery(docs *AllDocResult) (*DocQuery, error) {
	var d = &DocQuery{}
	for _, doc := range docs.Rows {
		d.Docs = append(d.Docs, struct {
			Id  string `json:"id"`
			Rev string `json:"rev"`
		}{
			Id:  doc.Id,
			Rev: doc.Value.Rev,
		})
	}
	return d, nil
}

func (c *Couch[T]) buildListQuery() (*DocQuery, error) {
	docs, err := c.allDocs()
	if err != nil {
		return nil, err
	}
	return listQuery(docs)
}

func (c *Couch[T]) List() ([]T, error) {
	query, err := c.buildListQuery()
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(query); err != nil {
		return nil, err
	}
	resp, err := http.Post(c.url+"/_bulk_get", "application/json", &buf)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	var result BulkGetResult[T]
	return result.Load(resp.Body)
}

func NewCouch[T core.Storable](url string) *Couch[T] {
	return &Couch[T]{
		url: url,
	}
}
