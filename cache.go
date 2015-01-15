package main

import (
	"github.com/hashicorp/golang-lru"
	fb "github.com/huandu/facebook"
)

type FacebookInfo struct {
	Id   string
	Name string
}

var (
	NameCache, _ = lru.New(256)
	GlobalApp    = fb.New(ClientId, ApiSecret)
)

func GetId(accessToken string) (string, error) {
	if val, ok := NameCache.Get(accessToken); ok {
		return val.(FacebookInfo).Id, nil
	}

	return exchangeAccessToken(accessToken)

}

func exchangeAccessToken(accessToken string) (string, error) {

	session := GlobalApp.Session(accessToken)
	err := session.Validate()
	if err != nil {
		return "", err
	}

	res, _ := session.Get("me?fields=id,name", nil)

	var ret string = res.Get("id").(string)
	var name string = res.Get("name").(string)

	NameCache.Add(accessToken, FacebookInfo{ret, name})
	return ret, nil
}

func SetId(accessToken, id string, name string) {
	NameCache.Add(accessToken, FacebookInfo{id, name})
}
