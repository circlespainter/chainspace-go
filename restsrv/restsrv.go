package restsrv // import "chainspace.io/prototype/restsrv"

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"chainspace.io/prototype/config"
	"chainspace.io/prototype/crypto/signature"
	"chainspace.io/prototype/log"
	"chainspace.io/prototype/log/fld"
	"chainspace.io/prototype/network"
	"chainspace.io/prototype/transactor"
	"chainspace.io/prototype/transactor/client"
	transactorclient2 "chainspace.io/prototype/transactor/client2"
	"github.com/rs/cors"
)

type Config struct {
	Addr       string
	Key        signature.KeyPair
	Port       int
	Top        *network.Topology
	MaxPayload config.ByteSize
	SelfID     uint64
}

type Service struct {
	port       int
	srv        *http.Server
	top        *network.Topology
	maxPayload config.ByteSize
	client     transactorclient2.Client
}

type resp struct {
	Data   interface{} `json:"data"`
	Status string      `json:"status"`
}

func response(rw http.ResponseWriter, status int, resp resp) {
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(status)
	b, _ := json.Marshal(resp)
	rw.Write(b)
}

func fail(rw http.ResponseWriter, status int, data interface{}) {
	response(rw, status, resp{data, "fail"})
}

func errorr(rw http.ResponseWriter, status int, data interface{}) {
	response(rw, status, resp{data, "error"})
}

func success(rw http.ResponseWriter, status int, data interface{}) {
	response(rw, status, resp{data, "success"})
}

func (s *Service) object(rw http.ResponseWriter, r *http.Request) {
	if !strings.EqualFold(r.Header.Get("Content-Type"), "application/json") {
		fail(rw, http.StatusBadRequest, "unsupported content-type")
		return
	}
	switch r.Method {
	case http.MethodPost:
		s.createObject(rw, r)
		return
	case http.MethodDelete:
		s.deleteObject(rw, r)
		return
	case http.MethodGet:
		s.queryObject(rw, r)
		return
	default:
		fail(rw, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
}

func readdata(rw http.ResponseWriter, r *http.Request) ([]byte, bool) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fail(rw, http.StatusBadRequest, fmt.Sprintf("unable to read request: %v", err))
		return nil, false
	}
	req := struct {
		Data string `json:"data"`
	}{}
	if err := json.Unmarshal(body, &req); err != nil {
		fail(rw, http.StatusBadRequest, fmt.Sprintf("unable to unmarshal: %v", err))
		return nil, false
	}
	if len(req.Data) <= 0 {
		fail(rw, http.StatusBadRequest, "empty data")
		return nil, false
	}
	key, err := base64.StdEncoding.DecodeString(req.Data)
	if err != nil {
		fail(rw, http.StatusBadRequest, fmt.Sprintf("unable to b64decode: %v", err))
		return nil, false
	}
	return key, true
}

func BuildObjectResponse(objects []*transactor.Object) (Object, error) {
	if len(objects) <= 0 {
		return Object{}, errors.New("object already inactive")
	}
	data := []Object{}
	for _, v := range objects {
		o := Object{
			Key:    base64.StdEncoding.EncodeToString(v.Key),
			Value:  base64.StdEncoding.EncodeToString(v.Value),
			Status: v.Status.String(),
		}
		data = append(data, o)

	}
	for _, v := range data {
		if v != data[0] {
			return Object{}, errors.New("inconsistent data")
		}
	}
	return data[0], nil
}

func (s *Service) createObject(rw http.ResponseWriter, r *http.Request) {
	rawObject, ok := readdata(rw, r)
	if !ok {
		return
	}
	txclient := transactorclient.New(&transactorclient.Config{Top: s.top, MaxPayload: s.maxPayload})
	defer txclient.Close()
	ids, err := txclient.Create(rawObject)
	if err != nil {
		errorr(rw, http.StatusInternalServerError, err.Error())
		return
	}
	for _, v := range ids {
		if string(v) != string(ids[0]) {
			errorr(rw, http.StatusInternalServerError, "inconsistent data")
			return
		}
	}
	res := struct {
		ID string `json:"id"`
	}{
		ID: base64.StdEncoding.EncodeToString(ids[0]),
	}
	success(rw, http.StatusOK, res)
}

func (s *Service) queryObject(rw http.ResponseWriter, r *http.Request) {
	key, ok := readdata(rw, r)
	if !ok {
		return
	}
	txclient := transactorclient.New(&transactorclient.Config{Top: s.top, MaxPayload: s.maxPayload})
	defer txclient.Close()
	objs, err := txclient.Query(key)
	if err != nil {
		errorr(rw, http.StatusInternalServerError, err.Error())
		return
	}
	obj, err := BuildObjectResponse(objs)
	if err != nil {
		errorr(rw, http.StatusInternalServerError, err.Error())
		return
	}

	success(rw, http.StatusOK, obj)
}

func (s *Service) deleteObject(rw http.ResponseWriter, r *http.Request) {
	key, ok := readdata(rw, r)
	if !ok {
		return
	}
	txclient := transactorclient.New(&transactorclient.Config{Top: s.top, MaxPayload: s.maxPayload})
	defer txclient.Close()
	objs, err := txclient.Delete(key)
	if err != nil {
		errorr(rw, http.StatusInternalServerError, err.Error())
		return
	}
	obj, err := BuildObjectResponse(objs)
	if err != nil {
		errorr(rw, http.StatusInternalServerError, err.Error())
		return
	}

	success(rw, http.StatusOK, obj)
}

func (s *Service) states(rw http.ResponseWriter, r *http.Request) {
	if !strings.EqualFold(r.Header.Get("Content-Type"), "application/json") {
		fail(rw, http.StatusBadRequest, "unsupported content-type")
		return
	}
	if r.Method != http.MethodPost {
		fail(rw, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	req := struct {
		Id uint64 `json:"id"`
	}{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fail(rw, http.StatusBadRequest, fmt.Sprintf("unable to read request: %v", err))
		return
	}
	if err := json.Unmarshal(body, &req); err != nil {
		fail(rw, http.StatusBadRequest, fmt.Sprintf("unable to unmarshal: %v", err))
		return
	}
	txclient := transactorclient.New(&transactorclient.Config{Top: s.top, MaxPayload: s.maxPayload})
	defer txclient.Close()
	states, err := txclient.States(req.Id)
	if err != nil {
		errorr(rw, http.StatusInternalServerError, err.Error())
		return
	}

	sort.Slice(states.States, func(i, j int) bool { return states.States[i].HashID < states.States[j].HashID })
	success(rw, http.StatusOK, states)
}

func (s *Service) transaction(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fail(rw, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	if !strings.EqualFold(r.Header.Get("Content-Type"), "application/json") {
		fail(rw, http.StatusBadRequest, "unsupported content-type")
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fail(rw, http.StatusBadRequest, fmt.Sprintf("unable to read request: %v", err))
		return
	}
	req := Transaction{}
	if err := json.Unmarshal(body, &req); err != nil {
		fail(rw, http.StatusBadRequest, fmt.Sprintf("unable to unmarshal: %v", err))
		return
	}
	// require at least input object.
	for _, v := range req.Traces {
		if len(v.InputObjectsKeys) <= 0 {
			fail(rw, http.StatusBadRequest, "no input objects for a trace")
			return
		}
	}
	// txclient := transactorclient.New(&transactorclient.Config{Top: s.top, MaxPayload: s.maxPayload})
	// defer txclient.Close()
	objects, err := s.client.SendTransaction(req.ToTransactor())
	if err != nil {
		errorr(rw, http.StatusInternalServerError, err.Error())
		return
	}
	data := []Object{}
	for _, v := range objects {
		v := v
		o := Object{
			Key:    base64.StdEncoding.EncodeToString(v.Key),
			Value:  base64.StdEncoding.EncodeToString(v.Value),
			Status: v.Status.String(),
		}
		data = append(data, o)
	}
	success(rw, http.StatusOK, data)
}

func (s *Service) objectsReady(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fail(rw, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	if !strings.EqualFold(r.Header.Get("Content-Type"), "application/json") {
		fail(rw, http.StatusBadRequest, "unsupported content-type")
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fail(rw, http.StatusBadRequest, fmt.Sprintf("unable to read request: %v", err))
		return
	}
	req := struct {
		Data []string `json:"data"`
	}{}
	if err := json.Unmarshal(body, &req); err != nil {
		fail(rw, http.StatusBadRequest, fmt.Sprintf("unable to unmarshal: %v", err))
		return
	}

	// defer
	for _, v := range req.Data {
		key, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			fail(rw, http.StatusBadRequest, fmt.Sprintf("unable to b64decode: %v", err))
			return
		}
		txclient := transactorclient.New(&transactorclient.Config{Top: s.top, MaxPayload: s.maxPayload})
		txclient.Close()
		objs, err := txclient.Query(key)
		if err != nil {
			success(rw, http.StatusOK, false)
			// errorr(rw, http.StatusInternalServerError, err.Error())
			return
		}
		// fmt.Printf("tx: %v\n", objs)
		if uint64(len(objs)) != s.top.ShardSize() {
			success(rw, http.StatusOK, false)
			return
		}
		for _, v := range objs {
			if !bytes.Equal(v.Key, objs[0].Key) {
				fail(rw, http.StatusInternalServerError, "inconsistent data")
				return
			}
		}
	}
	success(rw, http.StatusOK, true)
}

func (s *Service) makeServ(addr string, port int) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/object", s.object)
	mux.HandleFunc("/object/ready", s.objectsReady)
	mux.HandleFunc("/states", s.states)
	mux.HandleFunc("/transaction", s.transaction)
	handler := cors.Default().Handler(mux)
	h := &http.Server{
		Addr:    fmt.Sprintf("%v:%v", addr, port),
		Handler: handler,
		//ReadTimeout:  10 * time.Second,
		//WriteTimeout: 10 * time.Second,
	}
	h.SetKeepAlivesEnabled(false)
	return h
}

func New(cfg *Config) *Service {
	clcfg := transactorclient2.Config{
		NodeID:     cfg.SelfID,
		Top:        cfg.Top,
		MaxPayload: cfg.MaxPayload,
		Key:        cfg.Key,
	}
	txclient := transactorclient2.New(&clcfg)
	s := &Service{
		port:       cfg.Port,
		top:        cfg.Top,
		maxPayload: cfg.MaxPayload,
		client:     txclient,
	}
	s.srv = s.makeServ(cfg.Addr, cfg.Port)
	go func() {
		log.Info("http server started", fld.Port(cfg.Port))
		log.Fatal("http server exited", fld.Err(s.srv.ListenAndServe()))
	}()
	return s
}
