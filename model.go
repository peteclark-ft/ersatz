package main

import (
	"encoding/json"
	"errors"

	"github.com/peteclark-ft/ersatz/v1"
	"github.com/peteclark-ft/ersatz/v2"
)

var ErrUnsupportedVersion = errors.New("unsupported ersatz version, please confirm the ersatz-fixtures.yml version number")

type ersatz struct {
	Version  string   `json:"version"`
	Fixtures fixtures `json:"fixtures"`
}

type fixtures interface {
	Version() int
}

func (e *ersatz) UnmarshalJSON(data []byte) error {
	v := struct {
		Version string `json:"version"`
	}{}

	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	e.Version = v.Version

	f := struct {
		Fixtures fixtures `json:"fixtures"`
	}{}

	switch e.Version {
	case "1.0.0-rc1":
	case "1.0.0":
		f.Fixtures = &v1.Fixtures{}
	case "2.0.0-rc1":
		f.Fixtures = &v2.Fixtures{}
	default:
		return ErrUnsupportedVersion
	}

	err = json.Unmarshal(data, &f)
	if err != nil {
		return err
	}

	e.Fixtures = f.Fixtures
	return nil
}
