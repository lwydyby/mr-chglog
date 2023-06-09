package git

import (
	"errors"
	"strings"
	"time"
)

type Tag struct {
	Name     string
	Subject  string
	Date     time.Time
	Next     *Tag
	Previous *Tag
}

func Select(tags []*Tag, query string) ([]*Tag, *Tag, error) {
	if len(query) == 0 {
		return tags, nil, nil
	}
	tokens := strings.Split(query, "..")

	switch len(tokens) {
	case 1:
		return selectSingleTag(tags, tokens[0])
	case 2:
		old := tokens[0]
		new := tokens[1]
		switch {
		case old == "" && new == "":
			return tags, nil, nil
		case old == "":
			return selectBeforeTags(tags, new)
		case new == "":
			return selectAfterTags(tags, old)
		default:
			return selectRangeTags(tags, tokens[0], tokens[1])
		}
	}

	return nil, nil, errors.New("failed to parse the query")
}

func selectSingleTag(tags []*Tag, token string) ([]*Tag, *Tag, error) {
	for i, tag := range tags {
		if tag.Name == token {
			if i+1 < len(tags) {
				return []*Tag{tag}, tags[i+1], nil
			}
			return []*Tag{tag}, nil, nil
		}
	}

	return nil, nil, errors.New("could not find the tag")
}

func selectBeforeTags(tags []*Tag, token string) ([]*Tag, *Tag, error) {
	var (
		res    []*Tag
		enable bool
	)
	for i := range tags {
		if tags[i].Name == token {
			enable = true
		}

		if enable {
			res = append(res, tags[i])
		}
	}

	if len(res) == 0 {
		return res, nil, errors.New("could not find the tag")
	}

	return res, nil, nil
}

func selectAfterTags(tags []*Tag, token string) ([]*Tag, *Tag, error) {
	var (
		res []*Tag
	)
	var from *Tag
	for i := range tags {
		res = append(res, tags[i])
		if tags[i].Name == token {
			if i+1 < len(tags)-1 {
				from = tags[i+1]
			}
			break
		}
	}

	if len(res) == 0 {
		return res, nil, errors.New("could not find the tag")
	}

	return res, from, nil
}

func selectRangeTags(tags []*Tag, old string, new string) ([]*Tag, *Tag, error) {
	var (
		res    []*Tag
		enable bool
	)
	var from *Tag
	for i := range tags {
		if tags[i].Name == new {
			enable = true
		}

		if enable {
			res = append(res, tags[i])
		}

		if tags[i].Name == old {
			if i+1 < len(tags)-1 {
				from = tags[i+1]
			}
			break
		}
	}

	if len(res) == 0 {
		return res, nil, errors.New("could not find the tag")
	}

	return res, from, nil
}
