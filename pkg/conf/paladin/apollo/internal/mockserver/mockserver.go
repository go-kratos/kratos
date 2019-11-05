package mockserver

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"
)

type notification struct {
	NamespaceName  string `json:"namespaceName,omitempty"`
	NotificationID int    `json:"notificationId,omitempty"`
}

type result struct {
	// AppID          string            `json:"appId"`
	// Cluster        string            `json:"cluster"`
	NamespaceName  string            `json:"namespaceName"`
	Configurations map[string]string `json:"configurations"`
	ReleaseKey     string            `json:"releaseKey"`
}

type mockServer struct {
	server http.Server

	lock          sync.Mutex
	notifications map[string]int
	config        map[string]map[string]string
}

func (s *mockServer) NotificationHandler(rw http.ResponseWriter, req *http.Request) {
	s.lock.Lock()
	defer s.lock.Unlock()
	req.ParseForm()
	var notifications []notification
	if err := json.Unmarshal([]byte(req.FormValue("notifications")), &notifications); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	var changes []notification
	for _, noti := range notifications {
		if currentID := s.notifications[noti.NamespaceName]; currentID != noti.NotificationID {
			changes = append(changes, notification{NamespaceName: noti.NamespaceName, NotificationID: currentID})
		}
	}

	if len(changes) == 0 {
		rw.WriteHeader(http.StatusNotModified)
		return
	}
	bts, err := json.Marshal(&changes)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.Write(bts)
}

func (s *mockServer) ConfigHandler(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	strs := strings.Split(req.RequestURI, "/")
	var namespace, releaseKey = strings.Split(strs[4], "?")[0], req.FormValue("releaseKey")
	config := s.Get(namespace)

	var result = result{NamespaceName: namespace, Configurations: config, ReleaseKey: releaseKey}
	bts, err := json.Marshal(&result)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.Write(bts)
}

var server *mockServer

func (s *mockServer) Set(namespace, key, value string) {
	server.lock.Lock()
	defer server.lock.Unlock()

	notificationID := s.notifications[namespace]
	notificationID++
	s.notifications[namespace] = notificationID

	if kv, ok := s.config[namespace]; ok {
		kv[key] = value
		return
	}
	kv := map[string]string{key: value}
	s.config[namespace] = kv
}

func (s *mockServer) Get(namespace string) map[string]string {
	server.lock.Lock()
	defer server.lock.Unlock()

	return s.config[namespace]
}

func (s *mockServer) Delete(namespace, key string) {
	server.lock.Lock()
	defer server.lock.Unlock()

	if kv, ok := s.config[namespace]; ok {
		delete(kv, key)
	}

	notificationID := s.notifications[namespace]
	notificationID++
	s.notifications[namespace] = notificationID
}

// Set namespace's key value
func Set(namespace, key, value string) {
	server.Set(namespace, key, value)
}

// Delete namespace's key
func Delete(namespace, key string) {
	server.Delete(namespace, key)
}

// Run mock server
func Run() error {
	initServer()
	return server.server.ListenAndServe()
}

func initServer() {
	server = &mockServer{
		notifications: map[string]int{},
		config:        map[string]map[string]string{},
	}
	mux := http.NewServeMux()
	mux.Handle("/notifications/", http.HandlerFunc(server.NotificationHandler))
	mux.Handle("/configs/", http.HandlerFunc(server.ConfigHandler))
	server.server.Handler = mux
	server.server.Addr = ":8010"
}

// Close mock server
func Close() error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second))
	defer cancel()

	return server.server.Shutdown(ctx)
}
