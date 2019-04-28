package agent

import (
	"net/http"
	"strings"

	"go-common/app/service/main/bns/agent/backend"
	"go-common/library/log"
)

// NSTranslation query name from http api
func (s *HTTPServer) NSTranslation(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	name := strings.TrimPrefix(req.URL.Path, "/v1/naming/")

	inss, err := s.agent.Query(name)
	if err != nil {
		log.Error("call easyns server failed with naming translation, err: %s", err.Error())
		return nil, err
	}
	obj := struct {
		Instances []*backend.Instance `json:"instances"`
	}{inss}
	return obj, nil
}
