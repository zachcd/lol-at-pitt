package main

import (
	"github.com/hashicorp/golang-lru"
	fb "github.com/huandu/facebook"
)

var (
	NameCache, _ = lru.New(256)
	GlobalApp    = fb.New(ClientId, ApiSecret)
)

func GetId(accessToken string) (string, error) {
	if val, ok := NameCache.Get(accessToken); ok {
		return val.(string), nil
	}

	session := GlobalApp.Session(accessToken)
	err := session.Validate()
	if err != nil {
		return "", err
	}

	res, _ := session.Get("me?fields=id,name", nil)
	var ret string = res.Get("id").(string)
	NameCache.Add(accessToken, ret)
	return ret, nil
}

func SetId(accessToken, id string) {
	NameCache.Add(accessToken, id)
}
