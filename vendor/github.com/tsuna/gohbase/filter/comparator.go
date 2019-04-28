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

const comparatorPath = "org.apache.hadoop.hbase.filter."

// BitComparatorBitwiseOp is TODO
type BitComparatorBitwiseOp int32

func (o BitComparatorBitwiseOp) isValid() bool {
	return o >= 1 && o <= 3
}

// Constants are TODO
const (
	BitComparatorAND BitComparatorBitwiseOp = 1
	BitComparatorOR  BitComparatorBitwiseOp = 2
	BitComparatorXOR BitComparatorBitwiseOp = 3
)

// Ensure our types implement Comparator correctly.
var _ Comparator = (*BinaryComparator)(nil)
var _ Comparator = (*LongComparator)(nil)
var _ Comparator = (*BinaryPrefixComparator)(nil)
var _ Comparator = (*BitComparator)(nil)
var _ Comparator = (*NullComparator)(nil)
var _ Comparator = (*RegexStringComparator)(nil)
var _ Comparator = (*SubstringComparator)(nil)

// Comparator is TODO
type Comparator interface {
	// ConstructPBComparator creates and returns the comparator encoded in a
	// pb.Comparator type
	ConstructPBComparator() (*pb.Comparator, error)
}

// ByteArrayComparable is used across many Comparators.
type ByteArrayComparable pb.ByteArrayComparable

// NewByteArrayComparable is TODO
func NewByteArrayComparable(value []byte) *ByteArrayComparable {
	return &ByteArrayComparable{
		Value: value,
	}
}

func (b *ByteArrayComparable) toPB() *pb.ByteArrayComparable {
	return (*pb.ByteArrayComparable)(b)
}

// BinaryComparator is TODO
type BinaryComparator pb.BinaryComparator

// NewBinaryComparator is TODO
func NewBinaryComparator(comparable *ByteArrayComparable) *BinaryComparator {
	return &BinaryComparator{
		Comparable: comparable.toPB(),
	}
}

// ConstructPBComparator is TODO
func (c *BinaryComparator) ConstructPBComparator() (*pb.Comparator, error) {
	serializedComparator, err := proto.Marshal((*pb.BinaryComparator)(c))
	if err != nil {
		return nil, err
	}
	comparator := &pb.Comparator{
		Name:                 proto.String(comparatorPath + "BinaryComparator"),
		SerializedComparator: serializedComparator,
	}
	return comparator, nil
}

// LongComparator is TODO
type LongComparator pb.LongComparator

// NewLongComparator is TODO
func NewLongComparator(comparable *ByteArrayComparable) *LongComparator {
	return &LongComparator{
		Comparable: comparable.toPB(),
	}
}

// ConstructPBComparator is TODO
func (c *LongComparator) ConstructPBComparator() (*pb.Comparator, error) {
	serializedComparator, err := proto.Marshal((*pb.LongComparator)(c))
	if err != nil {
		return nil, err
	}
	comparator := &pb.Comparator{
		Name:                 proto.String(comparatorPath + "LongComparator"),
		SerializedComparator: serializedComparator,
	}
	return comparator, nil
}

// BinaryPrefixComparator is TODO
type BinaryPrefixComparator pb.BinaryPrefixComparator

// NewBinaryPrefixComparator is TODO
func NewBinaryPrefixComparator(comparable *ByteArrayComparable) *BinaryPrefixComparator {
	return &BinaryPrefixComparator{
		Comparable: comparable.toPB(),
	}
}

// ConstructPBComparator is TODO
func (c *BinaryPrefixComparator) ConstructPBComparator() (*pb.Comparator, error) {
	serializedComparator, err := proto.Marshal((*pb.BinaryPrefixComparator)(c))
	if err != nil {
		return nil, err
	}
	comparator := &pb.Comparator{
		Name:                 proto.String(comparatorPath + "BinaryPrefixComparator"),
		SerializedComparator: serializedComparator,
	}
	return comparator, nil
}

// BitComparator is TODO
type BitComparator pb.BitComparator

// NewBitComparator is TODO
func NewBitComparator(bitwiseOp BitComparatorBitwiseOp,
	comparable *ByteArrayComparable) *BitComparator {
	op := pb.BitComparator_BitwiseOp(bitwiseOp)
	return &BitComparator{
		Comparable: comparable.toPB(),
		BitwiseOp:  &op,
	}
}

// ConstructPBComparator is TODO
func (c *BitComparator) ConstructPBComparator() (*pb.Comparator, error) {
	if !BitComparatorBitwiseOp(*c.BitwiseOp).isValid() {
		return nil, errors.New("Invalid bitwise operator specified")
	}
	serializedComparator, err := proto.Marshal((*pb.BitComparator)(c))
	if err != nil {
		return nil, err
	}
	comparator := &pb.Comparator{
		Name:                 proto.String(comparatorPath + "BitComparator"),
		SerializedComparator: serializedComparator,
	}
	return comparator, nil
}

// NullComparator is TODO
type NullComparator struct{}

// NewNullComparator is TODO
func NewNullComparator() NullComparator {
	return NullComparator{}
}

// ConstructPBComparator is TODO
func (c NullComparator) ConstructPBComparator() (*pb.Comparator, error) {
	serializedComparator, err := proto.Marshal(&pb.NullComparator{})
	if err != nil {
		return nil, err
	}
	comparator := &pb.Comparator{
		Name:                 proto.String(comparatorPath + "NullComparator"),
		SerializedComparator: serializedComparator,
	}
	return comparator, nil
}

// RegexStringComparator is TODO
type RegexStringComparator pb.RegexStringComparator

// NewRegexStringComparator is TODO
func NewRegexStringComparator(pattern string, patternFlags int32,
	charset, engine string) *RegexStringComparator {
	return &RegexStringComparator{
		Pattern:      proto.String(pattern),
		PatternFlags: proto.Int32(patternFlags),
		Charset:      proto.String(charset),
		Engine:       proto.String(engine),
	}
}

// ConstructPBComparator is TODO
func (c *RegexStringComparator) ConstructPBComparator() (*pb.Comparator, error) {
	serializedComparator, err := proto.Marshal((*pb.RegexStringComparator)(c))
	if err != nil {
		return nil, err
	}
	comparator := &pb.Comparator{
		Name:                 proto.String(comparatorPath + "RegexStringComparator"),
		SerializedComparator: serializedComparator,
	}
	return comparator, nil
}

// SubstringComparator is TODO
type SubstringComparator pb.SubstringComparator

// NewSubstringComparator is TODO
func NewSubstringComparator(substr string) *SubstringComparator {
	return &SubstringComparator{
		Substr: proto.String(substr),
	}
}

// ConstructPBComparator is TODO
func (c *SubstringComparator) ConstructPBComparator() (*pb.Comparator, error) {
	serializedComparator, err := proto.Marshal((*pb.SubstringComparator)(c))
	if err != nil {
		return nil, err
	}
	comparator := &pb.Comparator{
		Name:                 proto.String(comparatorPath + "SubstringComparator"),
		SerializedComparator: serializedComparator,
	}
	return comparator, nil
}
