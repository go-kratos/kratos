package ipdb

import (
	"reflect"
	"time"
	"os"
)

type DistrictInfo struct {
	CountryName	string	`json:"country_name"`
	RegionName string 	`json:"region_name"`
	CityName string 	`json:"city_name"`
	DistrictName string `json:"district_name"`
	ChinaAdminCode string `json:"china_admin_code"`
	CoveringRadius string `json:"covering_radius"`
	Latitude     string `json:"latitude"`
	Longitude    string `json:"longitude"`
}

type District struct {
	reader *reader
}

func NewDistrict(name string) (*District, error) {

	r, e := newReader(name, &DistrictInfo{})
	if e != nil {
		return nil, e
	}

	return &District{
		reader: r,
	}, nil
}

func (db *District) Reload(name string) error {

	_, err := os.Stat(name)
	if err != nil {
		return err
	}

	reader, err := newReader(name, &DistrictInfo{})
	if err != nil {
		return err
	}

	db.reader = reader

	return nil
}

func (db *District) Find(addr, language string) ([]string, error) {
	return db.reader.find1(addr, language)
}

func (db *District) FindMap(addr, language string) (map[string]string, error) {

	data, err := db.reader.find1(addr, language)
	if err != nil {
		return nil, err
	}
	info := make(map[string]string, len(db.reader.meta.Fields))
	for k, v := range data {
		info[db.reader.meta.Fields[k]] = v
	}

	return info, nil
}

func (db *District) FindInfo(addr, language string) (*DistrictInfo, error) {

	data, err := db.reader.FindMap(addr, language)
	if err != nil {
		return nil, err
	}

	info := &DistrictInfo{}

	for k, v := range data {
		sv := reflect.ValueOf(info).Elem()
		sfv := sv.FieldByName(db.reader.refType[k])

		if !sfv.IsValid() {
			continue
		}
		if !sfv.CanSet() {
			continue
		}

		sft := sfv.Type()
		fv := reflect.ValueOf(v)
		if sft == fv.Type() {
			sfv.Set(fv)
		}
	}

	return info, nil
}

func (db *District) IsIPv4() bool {
	return db.reader.IsIPv4Support()
}

func (db *District) IsIPv6() bool {
	return db.reader.IsIPv6Support()
}

func (db *District) Languages() []string {
	return db.reader.Languages()
}

func (db *District) Fields() []string {
	return db.reader.meta.Fields
}

func (db *District) BuildTime() time.Time {
	return db.reader.Build()
}