package gist

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/ellofae/authentication-deanery/internal/models"
)

var GistText = &models.GistText{}

func processFetchedText(text []byte) error {
	return json.Unmarshal(text, GistText)
}

func FetchGistData(gistUrl string) error {
	resp, err := http.Get(gistUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := processFetchedText(body); err != nil {
		return err
	}

	return nil
}
