package main

import (
	"context"
	"net/url"
	"testing"
)

func TestOwns(t *testing.T) {
	db := DatabaseNew()
	c := context.Background()
	IRI, err := url.Parse("https://lol.net/user/lolbro")
	if err != nil {
		t.Error(err)
	}
	owns, err := db.Owns(c, IRI)
	if owns == true {
		t.Error(owns, err)
	}
	IRI, err = url.Parse("https://jade.moe/user/oct2pus")
	owns, err = db.Owns(c, IRI)
	if owns != true {
		t.Error(owns, err)
	}
}
