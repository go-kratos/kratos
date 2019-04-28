// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"testing"
)

func TestFunctionScoreQuery(t *testing.T) {
	q := NewFunctionScoreQuery().
		Query(NewTermQuery("name.last", "banon")).
		Add(NewTermQuery("name.last", "banon"), NewWeightFactorFunction(1.5)).
		AddScoreFunc(NewWeightFactorFunction(3)).
		AddScoreFunc(NewRandomFunction()).
		Boost(3).
		MaxBoost(10).
		ScoreMode("avg")
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"function_score":{"boost":3,"functions":[{"filter":{"term":{"name.last":"banon"}},"weight":1.5},{"weight":3},{"random_score":{}}],"max_boost":10,"query":{"term":{"name.last":"banon"}},"score_mode":"avg"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestFunctionScoreQueryWithNilFilter(t *testing.T) {
	q := NewFunctionScoreQuery().
		Query(NewTermQuery("tag", "wow")).
		AddScoreFunc(NewRandomFunction()).
		Boost(2.0).
		MaxBoost(12.0).
		BoostMode("multiply").
		ScoreMode("max")
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"function_score":{"boost":2,"boost_mode":"multiply","functions":[{"random_score":{}}],"max_boost":12,"query":{"term":{"tag":"wow"}},"score_mode":"max"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestFieldValueFactor(t *testing.T) {
	q := NewFunctionScoreQuery().
		Query(NewTermQuery("name.last", "banon")).
		AddScoreFunc(NewFieldValueFactorFunction().Modifier("sqrt").Factor(2).Field("income")).
		Boost(2.0).
		MaxBoost(12.0).
		BoostMode("multiply").
		ScoreMode("max")
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"function_score":{"boost":2,"boost_mode":"multiply","functions":[{"field_value_factor":{"factor":2,"field":"income","modifier":"sqrt"}}],"max_boost":12,"query":{"term":{"name.last":"banon"}},"score_mode":"max"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestFieldValueFactorWithWeight(t *testing.T) {
	q := NewFunctionScoreQuery().
		Query(NewTermQuery("name.last", "banon")).
		AddScoreFunc(NewFieldValueFactorFunction().Modifier("sqrt").Factor(2).Field("income").Weight(2.5)).
		Boost(2.0).
		MaxBoost(12.0).
		BoostMode("multiply").
		ScoreMode("max")
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"function_score":{"boost":2,"boost_mode":"multiply","functions":[{"field_value_factor":{"factor":2,"field":"income","modifier":"sqrt"},"weight":2.5}],"max_boost":12,"query":{"term":{"name.last":"banon"}},"score_mode":"max"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestFieldValueFactorWithMultipleScoreFuncsAndWeights(t *testing.T) {
	q := NewFunctionScoreQuery().
		Query(NewTermQuery("name.last", "banon")).
		AddScoreFunc(NewFieldValueFactorFunction().Modifier("sqrt").Factor(2).Field("income").Weight(2.5)).
		AddScoreFunc(NewScriptFunction(NewScript("_score * doc['my_numeric_field'].value")).Weight(1.25)).
		AddScoreFunc(NewWeightFactorFunction(0.5)).
		Boost(2.0).
		MaxBoost(12.0).
		BoostMode("multiply").
		ScoreMode("max")
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"function_score":{"boost":2,"boost_mode":"multiply","functions":[{"field_value_factor":{"factor":2,"field":"income","modifier":"sqrt"},"weight":2.5},{"script_score":{"script":{"inline":"_score * doc['my_numeric_field'].value"}},"weight":1.25},{"weight":0.5}],"max_boost":12,"query":{"term":{"name.last":"banon"}},"score_mode":"max"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestFunctionScoreQueryWithGaussScoreFunc(t *testing.T) {
	q := NewFunctionScoreQuery().
		Query(NewTermQuery("name.last", "banon")).
		AddScoreFunc(NewGaussDecayFunction().FieldName("pin.location").Origin("11, 12").Scale("2km").Offset("0km").Decay(0.33))
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"function_score":{"functions":[{"gauss":{"pin.location":{"decay":0.33,"offset":"0km","origin":"11, 12","scale":"2km"}}}],"query":{"term":{"name.last":"banon"}}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestFunctionScoreQueryWithGaussScoreFuncAndMultiValueMode(t *testing.T) {
	q := NewFunctionScoreQuery().
		Query(NewTermQuery("name.last", "banon")).
		AddScoreFunc(NewGaussDecayFunction().FieldName("pin.location").Origin("11, 12").Scale("2km").Offset("0km").Decay(0.33).MultiValueMode("avg"))
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"function_score":{"functions":[{"gauss":{"multi_value_mode":"avg","pin.location":{"decay":0.33,"offset":"0km","origin":"11, 12","scale":"2km"}}}],"query":{"term":{"name.last":"banon"}}}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}

func TestFunctionScoreQueryBug660(t *testing.T) {
	q := NewFunctionScoreQuery().
		Query(NewBoolQuery().Filter(NewTermsQuery("prod", 8199))).
		AddScoreFunc(
			NewGaussDecayFunction().
				FieldName("update").
				Origin(1501747403).
				Scale(1209600).
				Offset("259200").
				Decay(0.6).
				Weight(1),
		).
		BoostMode("sum").
		ScoreMode("sum")
	src, err := q.Source()
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)
	expected := `{"function_score":{"boost_mode":"sum","functions":[{"gauss":{"update":{"decay":0.6,"offset":"259200","origin":1501747403,"scale":1209600}},"weight":1}],"query":{"bool":{"filter":{"terms":{"prod":[8199]}}}},"score_mode":"sum"}}`
	if got != expected {
		t.Errorf("expected\n%s\n,got:\n%s", expected, got)
	}
}
