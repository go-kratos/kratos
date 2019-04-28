// Copyright (C) 2015  The GoHBase Authors.  All rights reserved.
// This file is part of GoHBase.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package filter

import (
	"errors"

	"github.com/golang/protobuf/proto"
	"github.com/tsuna/gohbase/pb"
)

const filterPath = "org.apache.hadoop.hbase.filter."

// ListOperator is TODO
type ListOperator int32

func (o ListOperator) isValid() bool {
	return o >= 1 && o <= 2
}

func (o ListOperator) toPB() *pb.FilterList_Operator {
	op := pb.FilterList_Operator(o)
	return &op
}

// Constants is TODO
const (
	MustPassAll ListOperator = 1
	MustPassOne ListOperator = 2
)

// CompareType is TODO
type CompareType int32

func (c CompareType) isValid() bool {
	return c >= 0 && c <= 6
}

// Constants is TODO
const (
	Less           CompareType = 0
	LessOrEqual    CompareType = 1
	Equal          CompareType = 2
	NotEqual       CompareType = 3
	GreaterOrEqual CompareType = 4
	Greater        CompareType = 5
	NoOp           CompareType = 6
)

// Ensure our types implement Filter correctly.
var _ Filter = (*List)(nil)
var _ Filter = (*ColumnCountGetFilter)(nil)
var _ Filter = (*ColumnPaginationFilter)(nil)
var _ Filter = (*ColumnPrefixFilter)(nil)
var _ Filter = (*ColumnRangeFilter)(nil)
var _ Filter = (*CompareFilter)(nil)
var _ Filter = (*DependentColumnFilter)(nil)
var _ Filter = (*FamilyFilter)(nil)
var _ Filter = (*Wrapper)(nil)
var _ Filter = (*FirstKeyOnlyFilter)(nil)
var _ Filter = (*FirstKeyValueMatchingQualifiersFilter)(nil)
var _ Filter = (*FuzzyRowFilter)(nil)
var _ Filter = (*InclusiveStopFilter)(nil)
var _ Filter = (*KeyOnlyFilter)(nil)
var _ Filter = (*MultipleColumnPrefixFilter)(nil)
var _ Filter = (*PageFilter)(nil)
var _ Filter = (*PrefixFilter)(nil)
var _ Filter = (*QualifierFilter)(nil)
var _ Filter = (*RandomRowFilter)(nil)
var _ Filter = (*RowFilter)(nil)
var _ Filter = (*SingleColumnValueFilter)(nil)
var _ Filter = (*SingleColumnValueExcludeFilter)(nil)
var _ Filter = (*SkipFilter)(nil)
var _ Filter = (*TimestampsFilter)(nil)
var _ Filter = (*ValueFilter)(nil)
var _ Filter = (*WhileMatchFilter)(nil)
var _ Filter = (*AllFilter)(nil)
var _ Filter = (*RowRange)(nil)
var _ Filter = (*MultiRowRangeFilter)(nil)

// Filter is TODO
type Filter interface {
	// ConstructPBFilter creates and returns the filter encoded in a pb.Filter type
	//	- For most filters this just involves creating the special filter object,
	//	  serializing it, and then creating a standard Filter object with the name and
	//	  serialization inside.
	//	- For FilterLists this requires creating the protobuf FilterList which contains
	//	  an array []*pb.Filter (meaning we have to create, serialize, create all objects
	//	  in that array), serialize the newly created pb.FilterList and then create a
	//	  pb.Filter object containing that new serialization.
	ConstructPBFilter() (*pb.Filter, error)
}

// BytesBytesPair is a type used in FuzzyRowFilter. Want to avoid users having
// to interact directly with the protobuf generated file so exposing here.
type BytesBytesPair pb.BytesBytesPair

// NewBytesBytesPair is TODO
func NewBytesBytesPair(first []byte, second []byte) *BytesBytesPair {
	return &BytesBytesPair{
		First:  first,
		Second: second,
	}
}

/*
    Each filter below has three primary methods/declarations, each of which can be summarized
    as follows -

	1. Type declaration. Create a new type for each filter. A 'Name' field is required but
	   you can create as many other fields as you like. These are purely local and will be
	   transcribed into a pb.Filter type by ConstructPBFilter()
	2. Constructor. Given a few parameters create the above type and return it to the callee.
	3. ConstructPBFilter. Take our local representation of a filter object and create the
	   appropriate pb.Filter object. Return the pb.Filter object.

	You may define any additional methods you like (see FilterList) but be aware that as soon
	as the returned object is type casted to a Filter (e.g. appending it to an array of Filters)
	it loses the ability to call those additional functions.
*/

// List is TODO
type List pb.FilterList

// NewList is TODO
func NewList(operator ListOperator, filters ...Filter) *List {
	f := &List{
		Operator: operator.toPB(),
	}
	f.AddFilters(filters...)
	return f
}

// AddFilters is TODO
func (f *List) AddFilters(filters ...Filter) {
	for _, filter := range filters {
		fpb, err := filter.ConstructPBFilter()
		if err != nil {
			panic(err)
		}
		f.Filters = append(f.Filters, fpb)
	}
}

// ConstructPBFilter is TODO
func (f *List) ConstructPBFilter() (*pb.Filter, error) {
	if !ListOperator(*f.Operator).isValid() {
		return nil, errors.New("invalid operator specified")
	}

	serializedFilter, err := proto.Marshal((*pb.FilterList)(f))
	if err != nil {
		return nil, err
	}
	filter := &pb.Filter{
		Name:             proto.String(filterPath + "FilterList"),
		SerializedFilter: serializedFilter,
	}
	return filter, nil
}

// ColumnCountGetFilter is TODO
type ColumnCountGetFilter pb.ColumnCountGetFilter

// NewColumnCountGetFilter is TODO
func NewColumnCountGetFilter(limit int32) *ColumnCountGetFilter {
	return &ColumnCountGetFilter{
		Limit: proto.Int32(limit),
	}
}

// ConstructPBFilter is TODO
func (f *ColumnCountGetFilter) ConstructPBFilter() (*pb.Filter, error) {
	serializedFilter, err := proto.Marshal((*pb.ColumnCountGetFilter)(f))
	if err != nil {
		return nil, err
	}
	filter := &pb.Filter{
		Name:             proto.String(filterPath + "ColumnCountGetFilter"),
		SerializedFilter: serializedFilter,
	}
	return filter, nil
}

// ColumnPaginationFilter is TODO
type ColumnPaginationFilter pb.ColumnPaginationFilter

// NewColumnPaginationFilter is TODO
func NewColumnPaginationFilter(limit, offset int32, columnOffset []byte) *ColumnPaginationFilter {
	return &ColumnPaginationFilter{
		Limit:        proto.Int32(limit),
		Offset:       proto.Int32(offset),
		ColumnOffset: columnOffset,
	}
}

// ConstructPBFilter is TODO
func (f *ColumnPaginationFilter) ConstructPBFilter() (*pb.Filter, error) {
	serializedFilter, err := proto.Marshal((*pb.ColumnPaginationFilter)(f))
	if err != nil {
		return nil, err
	}
	filter := &pb.Filter{
		Name:             proto.String(filterPath + "ColumnPaginationFilter"),
		SerializedFilter: serializedFilter,
	}
	return filter, nil
}

// ColumnPrefixFilter is TODO
type ColumnPrefixFilter pb.ColumnPrefixFilter

// NewColumnPrefixFilter is TODO
func NewColumnPrefixFilter(prefix []byte) *ColumnPrefixFilter {
	return &ColumnPrefixFilter{
		Prefix: prefix,
	}
}

// ConstructPBFilter is TODO
func (f *ColumnPrefixFilter) ConstructPBFilter() (*pb.Filter, error) {
	serializedFilter, err := proto.Marshal((*pb.ColumnPrefixFilter)(f))
	if err != nil {
		return nil, err
	}
	filter := &pb.Filter{
		Name:             proto.String(filterPath + "ColumnPrefixFilter"),
		SerializedFilter: serializedFilter,
	}
	return filter, nil
}

// ColumnRangeFilter is TODO
type ColumnRangeFilter pb.ColumnRangeFilter

// NewColumnRangeFilter is TODO
func NewColumnRangeFilter(minColumn, maxColumn []byte,
	minColumnInclusive, maxColumnInclusive bool) *ColumnRangeFilter {
	return &ColumnRangeFilter{
		MinColumn:          minColumn,
		MaxColumn:          maxColumn,
		MinColumnInclusive: proto.Bool(minColumnInclusive),
		MaxColumnInclusive: proto.Bool(maxColumnInclusive),
	}
}

// ConstructPBFilter is TODO
func (f *ColumnRangeFilter) ConstructPBFilter() (*pb.Filter, error) {
	serializedFilter, err := proto.Marshal((*pb.ColumnRangeFilter)(f))
	if err != nil {
		return nil, err
	}
	filter := &pb.Filter{
		Name:             proto.String(filterPath + "ColumnRangeFilter"),
		SerializedFilter: serializedFilter,
	}
	return filter, nil
}

// CompareFilter is TODO
type CompareFilter pb.CompareFilter

// NewCompareFilter is TODO
func NewCompareFilter(compareOp CompareType, comparatorObj Comparator) *CompareFilter {
	op := pb.CompareType(compareOp)
	obj, err := comparatorObj.ConstructPBComparator()
	if err != nil {
		panic(err)
	}
	return &CompareFilter{
		CompareOp:  &op,
		Comparator: obj,
	}
}

// ConstructPBFilter is TODO
func (f *CompareFilter) ConstructPBFilter() (*pb.Filter, error) {
	serializedFilter, err := proto.Marshal((*pb.CompareFilter)(f))
	if err != nil {
		return nil, err
	}
	filter := &pb.Filter{
		Name:             proto.String(filterPath + "CompareFilter"),
		SerializedFilter: serializedFilter,
	}
	return filter, nil
}

// DependentColumnFilter is TODO
type DependentColumnFilter pb.DependentColumnFilter

// NewDependentColumnFilter is TODO
func NewDependentColumnFilter(compareFilter *CompareFilter, columnFamily, columnQualifier []byte,
	dropDependentColumn bool) *DependentColumnFilter {
	return &DependentColumnFilter{
		CompareFilter:       (*pb.CompareFilter)(compareFilter),
		ColumnFamily:        columnFamily,
		ColumnQualifier:     columnQualifier,
		DropDependentColumn: proto.Bool(dropDependentColumn),
	}
}

// ConstructPBFilter is TODO
func (f *DependentColumnFilter) ConstructPBFilter() (*pb.Filter, error) {
	serializedFilter, err := proto.Marshal((*pb.DependentColumnFilter)(f))
	if err != nil {
		return nil, err
	}
	filter := &pb.Filter{
		Name:             proto.String(filterPath + "DependentColumnFilter"),
		SerializedFilter: serializedFilter,
	}
	return filter, nil
}

// FamilyFilter is TODO
type FamilyFilter pb.FamilyFilter

// NewFamilyFilter is TODO
func NewFamilyFilter(compareFilter *CompareFilter) *FamilyFilter {
	return &FamilyFilter{
		CompareFilter: (*pb.CompareFilter)(compareFilter),
	}
}

// ConstructPBFilter is TODO
func (f *FamilyFilter) ConstructPBFilter() (*pb.Filter, error) {
	serializedFilter, err := proto.Marshal((*pb.FamilyFilter)(f))
	if err != nil {
		return nil, err
	}
	filter := &pb.Filter{
		Name:             proto.String(filterPath + "FamilyFilter"),
		SerializedFilter: serializedFilter,
	}
	return filter, nil
}

// Wrapper is TODO
type Wrapper pb.FilterWrapper

// NewWrapper is TODO
func NewWrapper(wrappedFilter Filter) *Wrapper {
	f, err := wrappedFilter.ConstructPBFilter()
	if err != nil {
		panic(err)
	}
	return &Wrapper{
		Filter: f,
	}
}

// ConstructPBFilter is TODO
func (f *Wrapper) ConstructPBFilter() (*pb.Filter, error) {
	serializedFilter, err := proto.Marshal((*pb.FilterWrapper)(f))
	if err != nil {
		return nil, err
	}
	filter := &pb.Filter{
		Name:             proto.String(filterPath + "FilterWrapper"),
		SerializedFilter: serializedFilter,
	}
	return filter, nil
}

// FirstKeyOnlyFilter is TODO
type FirstKeyOnlyFilter struct{}

// NewFirstKeyOnlyFilter is TODO
func NewFirstKeyOnlyFilter() FirstKeyOnlyFilter {
	return FirstKeyOnlyFilter{}
}

// ConstructPBFilter is TODO
func (f FirstKeyOnlyFilter) ConstructPBFilter() (*pb.Filter, error) {
	return &pb.Filter{
		Name:             proto.String(filterPath + "FirstKeyOnlyFilter"),
		SerializedFilter: pb.MustMarshal(&pb.FirstKeyOnlyFilter{}),
	}, nil
}

// FirstKeyValueMatchingQualifiersFilter is TODO
type FirstKeyValueMatchingQualifiersFilter pb.FirstKeyValueMatchingQualifiersFilter

// NewFirstKeyValueMatchingQualifiersFilter is TODO
func NewFirstKeyValueMatchingQualifiersFilter(
	qualifiers [][]byte) *FirstKeyValueMatchingQualifiersFilter {
	return &FirstKeyValueMatchingQualifiersFilter{
		Qualifiers: qualifiers,
	}
}

// ConstructPBFilter is TODO
func (f *FirstKeyValueMatchingQualifiersFilter) ConstructPBFilter() (*pb.Filter, error) {
	serializedFilter, err := proto.Marshal((*pb.FirstKeyValueMatchingQualifiersFilter)(f))
	if err != nil {
		return nil, err
	}
	filter := &pb.Filter{
		Name:             proto.String(filterPath + "FirstKeyValueMatchingQualifiersFilter"),
		SerializedFilter: serializedFilter,
	}
	return filter, nil
}

// FuzzyRowFilter is TODO
type FuzzyRowFilter pb.FuzzyRowFilter

// NewFuzzyRowFilter is TODO
func NewFuzzyRowFilter(pairs []*BytesBytesPair) *FuzzyRowFilter {
	p := make([]*pb.BytesBytesPair, len(pairs))
	for i, pair := range pairs {
		p[i] = (*pb.BytesBytesPair)(pair)
	}
	return &FuzzyRowFilter{
		FuzzyKeysData: p,
	}
}

// ConstructPBFilter is TODO
func (f *FuzzyRowFilter) ConstructPBFilter() (*pb.Filter, error) {
	serializedFilter, err := proto.Marshal((*pb.FuzzyRowFilter)(f))
	if err != nil {
		return nil, err
	}
	filter := &pb.Filter{
		Name:             proto.String(filterPath + "FuzzyRowFilter"),
		SerializedFilter: serializedFilter,
	}
	return filter, nil
}

// InclusiveStopFilter is TODO
type InclusiveStopFilter pb.InclusiveStopFilter

// NewInclusiveStopFilter is TODO
func NewInclusiveStopFilter(stopRowKey []byte) *InclusiveStopFilter {
	return &InclusiveStopFilter{
		StopRowKey: stopRowKey,
	}
}

// ConstructPBFilter is TODO
func (f *InclusiveStopFilter) ConstructPBFilter() (*pb.Filter, error) {
	serializedFilter, err := proto.Marshal((*pb.InclusiveStopFilter)(f))
	if err != nil {
		return nil, err
	}
	filter := &pb.Filter{
		Name:             proto.String(filterPath + "InclusiveStopFilter"),
		SerializedFilter: serializedFilter,
	}
	return filter, nil
}

// KeyOnlyFilter is TODO
type KeyOnlyFilter pb.KeyOnlyFilter

// NewKeyOnlyFilter is TODO
func NewKeyOnlyFilter(lenAsVal bool) *KeyOnlyFilter {
	return &KeyOnlyFilter{
		LenAsVal: proto.Bool(lenAsVal),
	}
}

// ConstructPBFilter is TODO
func (f *KeyOnlyFilter) ConstructPBFilter() (*pb.Filter, error) {
	serializedFilter, err := proto.Marshal((*pb.KeyOnlyFilter)(f))
	if err != nil {
		return nil, err
	}
	filter := &pb.Filter{
		Name:             proto.String(filterPath + "KeyOnlyFilter"),
		SerializedFilter: serializedFilter,
	}
	return filter, nil
}

// MultipleColumnPrefixFilter is TODO
type MultipleColumnPrefixFilter pb.MultipleColumnPrefixFilter

// NewMultipleColumnPrefixFilter is TODO
func NewMultipleColumnPrefixFilter(sortedPrefixes [][]byte) *MultipleColumnPrefixFilter {
	return &MultipleColumnPrefixFilter{
		SortedPrefixes: sortedPrefixes,
	}
}

// ConstructPBFilter is TODO
func (f *MultipleColumnPrefixFilter) ConstructPBFilter() (*pb.Filter, error) {
	serializedFilter, err := proto.Marshal((*pb.MultipleColumnPrefixFilter)(f))
	if err != nil {
		return nil, err
	}
	filter := &pb.Filter{
		Name:             proto.String(filterPath + "MultipleColumnPrefixFilter"),
		SerializedFilter: serializedFilter,
	}
	return filter, nil
}

// PageFilter is TODO
type PageFilter pb.PageFilter

// NewPageFilter is TODO
func NewPageFilter(pageSize int64) *PageFilter {
	return &PageFilter{
		PageSize: proto.Int64(pageSize),
	}
}

// ConstructPBFilter is TODO
func (f *PageFilter) ConstructPBFilter() (*pb.Filter, error) {
	serializedFilter, err := proto.Marshal((*pb.PageFilter)(f))
	if err != nil {
		return nil, err
	}
	filter := &pb.Filter{
		Name:             proto.String(filterPath + "PageFilter"),
		SerializedFilter: serializedFilter,
	}
	return filter, nil
}

// PrefixFilter is TODO
type PrefixFilter pb.PrefixFilter

// NewPrefixFilter is TODO
func NewPrefixFilter(prefix []byte) *PrefixFilter {
	return &PrefixFilter{
		Prefix: prefix,
	}
}

// ConstructPBFilter is TODO
func (f *PrefixFilter) ConstructPBFilter() (*pb.Filter, error) {
	serializedFilter, err := proto.Marshal((*pb.PrefixFilter)(f))
	if err != nil {
		return nil, err
	}
	filter := &pb.Filter{
		Name:             proto.String(filterPath + "PrefixFilter"),
		SerializedFilter: serializedFilter,
	}
	return filter, nil
}

// QualifierFilter is TODO
type QualifierFilter pb.QualifierFilter

// NewQualifierFilter is TODO
func NewQualifierFilter(compareFilter *CompareFilter) *QualifierFilter {
	return &QualifierFilter{
		CompareFilter: (*pb.CompareFilter)(compareFilter),
	}
}

// ConstructPBFilter is TODO
func (f *QualifierFilter) ConstructPBFilter() (*pb.Filter, error) {
	serializedFilter, err := proto.Marshal((*pb.QualifierFilter)(f))
	if err != nil {
		return nil, err
	}
	filter := &pb.Filter{
		Name:             proto.String(filterPath + "QualifierFilter"),
		SerializedFilter: serializedFilter,
	}
	return filter, nil
}

// RandomRowFilter is TODO
type RandomRowFilter pb.RandomRowFilter

// NewRandomRowFilter is TODO
func NewRandomRowFilter(chance float32) *RandomRowFilter {
	return &RandomRowFilter{
		Chance: proto.Float32(chance),
	}
}

// ConstructPBFilter is TODO
func (f *RandomRowFilter) ConstructPBFilter() (*pb.Filter, error) {
	serializedFilter, err := proto.Marshal((*pb.RandomRowFilter)(f))
	if err != nil {
		return nil, err
	}
	filter := &pb.Filter{
		Name:             proto.String(filterPath + "RandomRowFilter"),
		SerializedFilter: serializedFilter,
	}
	return filter, nil
}

// RowFilter is TODO
type RowFilter pb.RowFilter

// NewRowFilter is TODO
func NewRowFilter(compareFilter *CompareFilter) *RowFilter {
	return &RowFilter{
		CompareFilter: (*pb.CompareFilter)(compareFilter),
	}
}

// ConstructPBFilter is TODO
func (f *RowFilter) ConstructPBFilter() (*pb.Filter, error) {
	serializedFilter, err := proto.Marshal((*pb.RowFilter)(f))
	if err != nil {
		return nil, err
	}
	filter := &pb.Filter{
		Name:             proto.String(filterPath + "RowFilter"),
		SerializedFilter: serializedFilter,
	}
	return filter, nil
}

// SingleColumnValueFilter is TODO
type SingleColumnValueFilter pb.SingleColumnValueFilter

// NewSingleColumnValueFilter is TODO
func NewSingleColumnValueFilter(columnFamily, columnQualifier []byte, compareOp CompareType,
	comparatorObj Comparator, filterIfMissing, latestVersionOnly bool) *SingleColumnValueFilter {
	obj, err := comparatorObj.ConstructPBComparator()
	if err != nil {
		panic(err)
	}
	return &SingleColumnValueFilter{
		ColumnFamily:      columnFamily,
		ColumnQualifier:   columnQualifier,
		CompareOp:         (*pb.CompareType)(&compareOp),
		Comparator:        obj,
		FilterIfMissing:   proto.Bool(filterIfMissing),
		LatestVersionOnly: proto.Bool(latestVersionOnly),
	}
}

// ConstructPB is TODO
func (f *SingleColumnValueFilter) ConstructPB() (*pb.SingleColumnValueFilter, error) {
	if !CompareType(*f.CompareOp).isValid() {
		return nil, errors.New("invalid compare operation specified")
	}

	return (*pb.SingleColumnValueFilter)(f), nil
}

// ConstructPBFilter is TODO
func (f *SingleColumnValueFilter) ConstructPBFilter() (*pb.Filter, error) {
	serializedFilter, err := proto.Marshal((*pb.SingleColumnValueFilter)(f))
	if err != nil {
		return nil, err
	}
	filter := &pb.Filter{
		Name:             proto.String(filterPath + "SingleColumnValueFilter"),
		SerializedFilter: serializedFilter,
	}
	return filter, nil
}

// SingleColumnValueExcludeFilter is TODO
type SingleColumnValueExcludeFilter pb.SingleColumnValueExcludeFilter

// NewSingleColumnValueExcludeFilter is TODO
func NewSingleColumnValueExcludeFilter(
	filter *SingleColumnValueFilter) *SingleColumnValueExcludeFilter {
	return &SingleColumnValueExcludeFilter{
		SingleColumnValueFilter: (*pb.SingleColumnValueFilter)(filter),
	}
}

// ConstructPBFilter is TODO
func (f *SingleColumnValueExcludeFilter) ConstructPBFilter() (*pb.Filter, error) {
	serializedFilter, err := proto.Marshal((*pb.SingleColumnValueExcludeFilter)(f))
	if err != nil {
		return nil, err
	}
	filter := &pb.Filter{
		Name:             proto.String(filterPath + "SingleColumnValueExcludeFilter"),
		SerializedFilter: serializedFilter,
	}
	return filter, nil
}

// SkipFilter is TODO
type SkipFilter pb.SkipFilter

// NewSkipFilter is TODO
func NewSkipFilter(skippingFilter Filter) *SkipFilter {
	f, err := skippingFilter.ConstructPBFilter()
	if err != nil {
		panic(err)
	}
	return &SkipFilter{
		Filter: f,
	}
}

// ConstructPBFilter is TODO
func (f *SkipFilter) ConstructPBFilter() (*pb.Filter, error) {
	serializedFilter, err := proto.Marshal((*pb.SkipFilter)(f))
	if err != nil {
		return nil, err
	}
	filter := &pb.Filter{
		Name:             proto.String(filterPath + "SkipFilter"),
		SerializedFilter: serializedFilter,
	}
	return filter, nil
}

// TimestampsFilter is TODO
type TimestampsFilter pb.TimestampsFilter

// NewTimestampsFilter is TODO
func NewTimestampsFilter(timestamps []int64) *TimestampsFilter {
	return &TimestampsFilter{
		Timestamps: timestamps,
	}
}

// ConstructPBFilter is TODO
func (f *TimestampsFilter) ConstructPBFilter() (*pb.Filter, error) {
	serializedFilter, err := proto.Marshal((*pb.TimestampsFilter)(f))
	if err != nil {
		return nil, err
	}
	filter := &pb.Filter{
		Name:             proto.String(filterPath + "TimestampsFilter"),
		SerializedFilter: serializedFilter,
	}
	return filter, nil
}

// ValueFilter is TODO
type ValueFilter pb.ValueFilter

// NewValueFilter is TODO
func NewValueFilter(compareFilter *CompareFilter) *ValueFilter {
	return &ValueFilter{
		CompareFilter: (*pb.CompareFilter)(compareFilter),
	}
}

// ConstructPBFilter is TODO
func (f *ValueFilter) ConstructPBFilter() (*pb.Filter, error) {
	serializedFilter, err := proto.Marshal((*pb.ValueFilter)(f))
	if err != nil {
		return nil, err
	}
	filter := &pb.Filter{
		Name:             proto.String(filterPath + "ValueFilter"),
		SerializedFilter: serializedFilter,
	}
	return filter, nil
}

// WhileMatchFilter is TODO
type WhileMatchFilter pb.WhileMatchFilter

// NewWhileMatchFilter is TODO
func NewWhileMatchFilter(matchingFilter Filter) *WhileMatchFilter {
	f, err := matchingFilter.ConstructPBFilter()
	if err != nil {
		panic(err)
	}
	return &WhileMatchFilter{
		Filter: f,
	}
}

// ConstructPBFilter is TODO
func (f *WhileMatchFilter) ConstructPBFilter() (*pb.Filter, error) {
	serializedFilter, err := proto.Marshal((*pb.WhileMatchFilter)(f))
	if err != nil {
		return nil, err
	}
	filter := &pb.Filter{
		Name:             proto.String(filterPath + "WhileMatchFilter"),
		SerializedFilter: serializedFilter,
	}
	return filter, nil
}

// AllFilter is TODO
type AllFilter struct{}

// NewAllFilter is TODO
func NewAllFilter() AllFilter {
	return AllFilter{}
}

// ConstructPBFilter is TODO
func (f *AllFilter) ConstructPBFilter() (*pb.Filter, error) {
	return &pb.Filter{
		Name:             proto.String(filterPath + "FilterAllFilter"),
		SerializedFilter: pb.MustMarshal(&pb.FilterAllFilter{}),
	}, nil
}

// RowRange is TODO
type RowRange pb.RowRange

// NewRowRange is TODO
func NewRowRange(startRow, stopRow []byte, startRowInclusive, stopRowInclusive bool) *RowRange {
	return &RowRange{
		StartRow:          startRow,
		StartRowInclusive: proto.Bool(startRowInclusive),
		StopRow:           stopRow,
		StopRowInclusive:  proto.Bool(stopRowInclusive),
	}
}

// ConstructPBFilter is TODO
func (f *RowRange) ConstructPBFilter() (*pb.Filter, error) {
	serializedFilter, err := proto.Marshal((*pb.RowRange)(f))
	if err != nil {
		return nil, err
	}
	filter := &pb.Filter{
		Name:             proto.String(filterPath + "RowRange"),
		SerializedFilter: serializedFilter,
	}
	return filter, nil
}

// MultiRowRangeFilter is TODO
type MultiRowRangeFilter pb.MultiRowRangeFilter

// NewMultiRowRangeFilter is TODO
func NewMultiRowRangeFilter(rowRangeList []*RowRange) *MultiRowRangeFilter {
	rangeList := make([]*pb.RowRange, len(rowRangeList))
	for i, rr := range rowRangeList {
		rangeList[i] = (*pb.RowRange)(rr)
	}
	return &MultiRowRangeFilter{
		RowRangeList: rangeList,
	}
}

// ConstructPBFilter is TODO
func (f *MultiRowRangeFilter) ConstructPBFilter() (*pb.Filter, error) {
	serializedFilter, err := proto.Marshal((*pb.MultiRowRangeFilter)(f))
	if err != nil {
		return nil, err
	}
	filter := &pb.Filter{
		Name:             proto.String(filterPath + "MultiRowRangeFilter"),
		SerializedFilter: serializedFilter,
	}
	return filter, nil
}
