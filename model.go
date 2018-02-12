package main

import (
	"encoding/json"
	"errors"

	"github.com/peteclark-ft/ersatz/v1"
)

var ErrUnsupportedVersion = errors.New("unsupported ersatz version, please confirm the fixtures.yml version")

type ersatz struct {
	Version  int      `json:"version"`
	Fixtures fixtures `json:"fixtures"`
}

type fixtures interface {
	Version() int
}

func (e *ersatz) UnmarshalJSON(data []byte) error {
	v := struct {
		Version int `json:"version"`
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
	case 1:
		f.Fixtures = &v1.Fixtures{}
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
