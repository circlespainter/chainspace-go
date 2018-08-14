package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

const (
	entityKind = "nodeConfigList"
)

type auth struct {
	NetworkID string `json:"network_id"`
	Token     string `json:"token"`
}

type nodeConfig struct {
	NodeID string `json:"node_id"`
	Port   string `json:"port"`
	Addr   string `json:"addr"`
}

type Entity struct {
	Conf []nodeConfig `json:"conf"`
}

func main() {
	http.HandleFunc("/contacts.list", handleContactsList)
	http.HandleFunc("/contacts.set", handleContactsSet)
	appengine.Main()
}

func handleContactsSet(rw http.ResponseWriter, r *http.Request) {
	if !strings.EqualFold(r.Header.Get("Content-Type"), "application/json") {
		http.Error(rw, "missing content-type", http.StatusBadRequest)
		return
	}
	if !strings.EqualFold(r.Method, http.MethodPost) {
		http.Error(rw, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, fmt.Sprintf("unable to read request: %v", err), http.StatusBadRequest)
		return
	}
	req := struct {
		Auth   auth       `json:"auth"`
		Config nodeConfig `json:"config"`
	}{}
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(rw, fmt.Sprintf("unable to unmarshal: %v", err), http.StatusBadRequest)
		return
	}

	ctx := appengine.NewContext(r)
	key := datastore.NewKey(ctx, entityKind, req.Auth.NetworkID+req.Auth.Token, 0, nil)
	err = datastore.RunInTransaction(ctx, func(tctx context.Context) error {
		var e Entity
		err := datastore.Get(tctx, key, &e)
		if err != nil && err != datastore.ErrNoSuchEntity {
			return err
		}
		err = nil
		req.Config.Addr = r.RemoteAddr
		e.Conf = append(e.Conf, req.Config)
		_, err = datastore.Put(tctx, key, &e)
		return err
	}, nil)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusNoContent)
}

func handleContactsList(rw http.ResponseWriter, r *http.Request) {
	if !strings.EqualFold(r.Header.Get("Content-Type"), "application/json") {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(rw, "missing content-type")
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, fmt.Sprintf("unable to read request: %v", err), http.StatusBadRequest)
		return
	}
	req := struct {
		Auth auth `json:"auth"`
	}{}
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(rw, fmt.Sprintf("unable to unmarshal: %v", err), http.StatusBadRequest)
		return
	}
	ctx := appengine.NewContext(r)
	key := datastore.NewKey(ctx, entityKind, req.Auth.NetworkID+req.Auth.Token, 0, nil)
	var e Entity
	if err = datastore.Get(ctx, key, &e); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	b, err := json.Marshal(e.Conf)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(rw, "%v", string(b))
}
