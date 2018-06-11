package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"path"
	"strings"
	"time"
)

type RestClient struct {
	Addr    string
	Timeout int
	Headers map[string]string
}

func NewRest(addr string, timeout int, headers map[string]string) *RestClient {
	if addr == "" {
		addr = "http://127.0.0.1:8088"
	}

	return &RestClient{
		Addr:    addr,
		Timeout: timeout,
		Headers: headers,
	}
}

func (r *RestClient) Get(table, rowkey, cf string) (CellSet, error) {
	var cellset CellSet
	if table == "" {
		return cellset, errors.New("need params table.")
	}
	if rowkey == "" {
		return cellset, errors.New("need params rowkey.")
	}

	rowkey = url.QueryEscape(rowkey)
	u, _ := url.Parse(r.Addr)
	u.Path = path.Join(u.Path, table, rowkey, cf)
	link := u.String()

	req := NewRequest(r.Timeout, r.Headers)
	resp, err := req.Get(link)
	if err != nil {
		return cellset, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return cellset, err
	}

	if resp.StatusCode == 404 {
		return cellset, nil

	}

	if resp.StatusCode != 200 {
		return cellset, errors.New(string(body))
	}

	if err := json.Unmarshal(body, &cellset); err != nil {
		return cellset, err
	}

	for i, r := range cellset.Row {
		cellset.Row[i].Key, err = Base64Decode(r.Key)
		if err != nil {
			return cellset, err
		}
		for n, v := range r.Cell {
			cellset.Row[i].Cell[n].Column, _ = Base64Decode(v.Column)
			cellset.Row[i].Cell[n].Value, _ = Base64Decode(v.Value)
		}
	}
	return cellset, nil
}

func (r *RestClient) Put(table, rowkey, cf, value string) error {
	if table == "" {
		return errors.New("need params table.")
	}
	if rowkey == "" {
		return errors.New("need params rowkey.")
	}
	if cf == "" {
		return errors.New("need params cf.")
	}
	if value == "" {
		return errors.New("need params value.")
	}

	var cellset CellSet
	var row Row
	var cell Cell

	cell.Column = Base64Encode(cf)
	cell.Timestamp = time.Now().UnixNano() / 1e6
	cell.Value = Base64Encode(value)

	row.Key = Base64Encode(rowkey)
	row.Cell = []Cell{cell}

	cellset.Row = []Row{row}

	cs, err := json.Marshal(cellset)
	if err != nil {
		return errors.New(fmt.Sprintf("json marshal error: %v", err))
	}
	data := strings.NewReader(string(cs))

	rowkey = url.QueryEscape(rowkey)
	u, _ := url.Parse(r.Addr)
	u.Path = path.Join(u.Path, table, rowkey)
	link := u.String()

	req := NewRequest(r.Timeout, r.Headers)
	resp, err := req.Post(link, data)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return errors.New(string(body))
}

func (r *RestClient) Puts(res []map[string]string) error {
	var cellset CellSet
	var table string
	for _, re := range res {
		table, _ = re["table"]
		rowkey, _ := re["rowkey"]
		family, _ := re["family"]
		column, _ := re["column"]
		value, _ := re["value"]

		if len(table) == 0 || len(rowkey) == 0 || len(family) == 0 || len(column) == 0 || len(value) == 0 {
			return errors.New("the one of params is empty")
		}
		var cell Cell
		var row Row

		cf := fmt.Sprintf("%s:%s", family, column)
		cell.Column = Base64Encode(cf)
		cell.Timestamp = time.Now().UnixNano() / 1e6
		cell.Value = Base64Encode(value)

		row.Key = Base64Encode(rowkey)
		row.Cell = []Cell{cell}
		cellset.Row = append(cellset.Row, row)
	}

	cs, err := json.Marshal(cellset)
	if err != nil {
		return errors.New(fmt.Sprintf("json marshal error: %v", err))
	}
	data := strings.NewReader(string(cs))

	u, _ := url.Parse(r.Addr)
	u.Path = path.Join(u.Path, table, "false-row-key")
	link := u.String()

	req := NewRequest(r.Timeout, r.Headers)
	resp, err := req.Post(link, data)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return errors.New(string(body))
}
