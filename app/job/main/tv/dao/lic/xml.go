package lic

import (
	"encoding/xml"
	"net/url"

	model "go-common/app/job/main/tv/model/pgc"
)

// PrepareXML combine the xml data to sync
func PrepareXML(v *model.License) (body string) {
	xmlRes, _ := xml.MarshalIndent(v.XMLData.Service, " ", " ")
	params := url.Values{}
	params.Set("tId", v.TId)
	params.Set("inputTime", v.InputTime)
	params.Set("sign", v.Sign)
	body = params.Encode()
	body = body + "&xmlData=<?xml version=\"1.0\" encoding=\"UTF-8\"?> " + string(xmlRes)
	return body
}
