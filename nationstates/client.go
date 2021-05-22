package nationstates

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type Client struct {
	Password  string
	Autologin string
	Pin       string

	client *http.Client
}

func (c *Client) do(options map[string]interface{}) (Nation, error) {
	req, err := http.NewRequest(http.MethodGet, "https://www.nationstates.net/cgi-bin/api.cgi", nil)
	if err != nil {
		return Nation{}, err
	}
	req.Header.Set("User-Agent", "NationStates Go client")
	if password := c.Password; password != "" {
		req.Header.Set("X-Password", password)
	}
	if autologin := c.Autologin; autologin != "" {
		req.Header.Set("X-Autologin", autologin)
	}
	if pin := c.Pin; pin != "" {
		req.Header.Set("X-Pin", pin)
	}
	var params []string
	for k, v := range options {
		params = append(params, fmt.Sprintf("%v=%v", k, v))
	}
	req.URL.RawQuery = strings.Join(params, "&")
	client := http.DefaultClient
	if c.client != nil {
		client = c.client
	}
	res, err := client.Do(req)
	if err != nil {
		return Nation{}, err
	}
	if pin := res.Header.Get("X-Pin"); pin != "" {
		c.Pin = pin
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return Nation{}, err
	}
	var n Nation
	err = xml.Unmarshal(body, &n)
	if err != nil {
		return Nation{}, errors.New(string(body))
	}
	return n, nil
}

// GetNation is a generic method for querying the Nation API.
func (c *Client) GetNation(nation string, shards []string, options map[string]interface{}) (Nation, error) {
	opts := make(map[string]interface{})
	for k, v := range options {
		opts[k] = v
	}
	opts["nation"] = nation
	opts["q"] = strings.Join(shards, "+")
	return c.do(opts)
}

// GetIssues is a convenience method for getting issues for a nation.
func (c *Client) GetIssues(nation string) ([]Issue, error) {
	n, err := c.GetNation(nation, []string{"issues"}, nil)
	if err != nil {
		return nil, err
	}
	return n.Issues, nil
}

// GetNotices is a convenience method for getting notices for a nation.
func (c *Client) GetNotices(nation string) ([]Notice, error) {
	n, err := c.GetNation(nation, []string{"notices"}, nil)
	if err != nil {
		return nil, err
	}
	return n.Notices, nil
}

// GetNoticesSince is a convenience method for getting notices for a nation since a given offset.
func (c *Client) GetNoticesSince(nation string, from int) ([]Notice, error) {
	n, err := c.GetNation(nation, []string{"notices"}, map[string]interface{}{"from": from})
	if err != nil {
		return nil, err
	}
	return n.Notices, nil
}

func (c *Client) AnswerIssue(nation string, issue, option int) (Consequences, error) {
	n, err := c.do(map[string]interface{}{
		"nation": nation,
		"c":      "issue",
		"issue":  issue,
		"option": option,
	})
	if err != nil {
		return Consequences{}, err
	}
	return n.Consequences, nil
}

// CreateCensusShard creates a shard to query the census.
// scales takes:
// - a list of scale IDs to select which scales to query.
// - "all" to retrieve all scales.
// - nil to retreive the daily World Census scale.
// modes takes:
// - a list of modes to select which stats to return.
// - a map of from and to times in int64 UNIX format to retrieve a history of points.
// - nil to retrieve the score, rank and regional rank.
func (c *Client) CreateCensusShard(scales interface{}, modes interface{}) string {
	shard := []string{"census"}

	switch v := scales.(type) {
	case []int:
		var scales []string
		for _, scale := range v {
			scales = append(scales, strconv.Itoa(scale))
		}
		shard = append(shard, fmt.Sprintf("scale=%s", strings.Join(scales, "+")))
	case string:
		shard = append(shard, fmt.Sprintf("scale=%s", v))
	}

	switch v := modes.(type) {
	case []string:
		shard = append(shard, fmt.Sprintf("mode=%s", strings.Join(v, "+")))
	case map[string]int64:
		shard = append(shard, "mode=history", fmt.Sprintf("from=%v", v["from"]), fmt.Sprintf("to=%v", v["to"]))
	}

	return strings.Join(shard, ";")
}
