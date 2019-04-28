// Copyright 2018 Twitch Interactive, Inc.  All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may not
// use this file except in compliance with the License. A copy of the License is
// located at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// or in the "license" file accompanying this file. This file is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.
//
//
// This file contains some code from https://github.com/golang/protobuf:
// Copyright 2010 The Go Authors.  All rights reserved.
// https://github.com/golang/protobuf
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//     * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package gen

import (
	"fmt"
	"path"
	"strconv"
	"strings"

	"go-common/app/tool/liverpc/protoc-gen-liverpc/gen/stringutils"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

// Each type we import as a protocol buffer (other than FileDescriptorProto) needs
// a pointer to the FileDescriptorProto that represents it.  These types achieve that
// wrapping by placing each Proto inside a struct with the pointer to its File. The
// structs have the same names as their contents, with "Proto" removed.
// FileDescriptor is used to store the things that it points to.

// WrapTypes walks the incoming data, wrapping DescriptorProtos, EnumDescriptorProtos
// and FileDescriptorProtos into file-referenced objects within the Generator.
// It also creates the list of files to generate and so should be called before GenerateAllFiles.
func WrapTypes(req *plugin.CodeGeneratorRequest) (genFiles, allFiles []*FileDescriptor, allFilesByName map[string]*FileDescriptor) {
	allFiles = make([]*FileDescriptor, 0, len(req.ProtoFile))
	allFilesByName = make(map[string]*FileDescriptor, len(allFiles))

	for _, f := range req.ProtoFile {
		// We must wrap the descriptors before we wrap the enums
		descs := wrapDescriptors(f)
		buildNestedDescriptors(descs)
		enums := wrapEnumDescriptors(f, descs)
		buildNestedEnums(descs, enums)
		exts := wrapExtensions(f)
		svcs := wrapServices(f)
		fd := &FileDescriptor{
			FileDescriptorProto: f,
			Services:            svcs,
			Descriptors:         descs,
			Enums:               enums,
			Extensions:          exts,
			proto3:              fileIsProto3(f),
		}
		extractComments(fd)

		allFiles = append(allFiles, fd)
		allFilesByName[f.GetName()] = fd
	}

	for _, fd := range allFiles {
		fd.Imported = wrapImported(fd.FileDescriptorProto, allFilesByName)
	}

	genFiles = make([]*FileDescriptor, 0, len(req.FileToGenerate))
	for _, fileName := range req.FileToGenerate {
		fd := allFilesByName[fileName]
		if fd == nil {
			Fail("could not find file named", fileName)
		}
		fd.Index = len(genFiles)
		genFiles = append(genFiles, fd)
	}

	return genFiles, allFiles, allFilesByName
}

// The file and package name method are common to messages and enums.
type common struct {
	file *descriptor.FileDescriptorProto // File this object comes from.
}

func (c *common) File() *descriptor.FileDescriptorProto { return c.file }

func fileIsProto3(file *descriptor.FileDescriptorProto) bool {
	return file.GetSyntax() == "proto3"
}

// Descriptor represents a protocol buffer message.
type Descriptor struct {
	common
	*descriptor.DescriptorProto
	Parent   *Descriptor            // The containing message, if any.
	nested   []*Descriptor          // Inner messages, if any.
	enums    []*EnumDescriptor      // Inner enums, if any.
	ext      []*ExtensionDescriptor // Extensions, if any.
	typename []string               // Cached typename vector.
	index    int                    // The index into the container, whether the file or another message.
	path     string                 // The SourceCodeInfo path as comma-separated integers.
	group    bool
}

func newDescriptor(desc *descriptor.DescriptorProto, parent *Descriptor, file *descriptor.FileDescriptorProto, index int) *Descriptor {
	d := &Descriptor{
		common:          common{file},
		DescriptorProto: desc,
		Parent:          parent,
		index:           index,
	}
	if parent == nil {
		d.path = fmt.Sprintf("%d,%d", messagePath, index)
	} else {
		d.path = fmt.Sprintf("%s,%d,%d", parent.path, messageMessagePath, index)
	}

	// The only way to distinguish a group from a message is whether
	// the containing message has a TYPE_GROUP field that matches.
	if parent != nil {
		parts := d.TypeName()
		if file.Package != nil {
			parts = append([]string{*file.Package}, parts...)
		}
		exp := "." + strings.Join(parts, ".")
		for _, field := range parent.Field {
			if field.GetType() == descriptor.FieldDescriptorProto_TYPE_GROUP && field.GetTypeName() == exp {
				d.group = true
				break
			}
		}
	}

	for _, field := range desc.Extension {
		d.ext = append(d.ext, &ExtensionDescriptor{common{file}, field, d})
	}

	return d
}

// Return a slice of all the Descriptors defined within this file
func wrapDescriptors(file *descriptor.FileDescriptorProto) []*Descriptor {
	sl := make([]*Descriptor, 0, len(file.MessageType)+10)
	for i, desc := range file.MessageType {
		sl = wrapThisDescriptor(sl, desc, nil, file, i)
	}
	return sl
}

// Wrap this Descriptor, recursively
func wrapThisDescriptor(sl []*Descriptor, desc *descriptor.DescriptorProto, parent *Descriptor, file *descriptor.FileDescriptorProto, index int) []*Descriptor {
	sl = append(sl, newDescriptor(desc, parent, file, index))
	me := sl[len(sl)-1]
	for i, nested := range desc.NestedType {
		sl = wrapThisDescriptor(sl, nested, me, file, i)
	}
	return sl
}

func buildNestedDescriptors(descs []*Descriptor) {
	for _, desc := range descs {
		if len(desc.NestedType) != 0 {
			for _, nest := range descs {
				if nest.Parent == desc {
					desc.nested = append(desc.nested, nest)
				}
			}
			if len(desc.nested) != len(desc.NestedType) {
				Fail("internal error: nesting failure for", desc.GetName())
			}
		}
	}
}

// TypeName returns the elements of the dotted type name.
// The package name is not part of this name.
func (d *Descriptor) TypeName() []string {
	if d.typename != nil {
		return d.typename
	}
	n := 0
	for parent := d; parent != nil; parent = parent.Parent {
		n++
	}
	s := make([]string, n)
	for parent := d; parent != nil; parent = parent.Parent {
		n--
		s[n] = parent.GetName()
	}
	d.typename = s
	return s
}

// EnumDescriptor describes an enum. If it's at top level, its parent will be nil.
// Otherwise it will be the descriptor of the message in which it is defined.
type EnumDescriptor struct {
	common
	*descriptor.EnumDescriptorProto
	parent   *Descriptor // The containing message, if any.
	typename []string    // Cached typename vector.
	index    int         // The index into the container, whether the file or a message.
	path     string      // The SourceCodeInfo path as comma-separated integers.
}

// Construct a new EnumDescriptor
func newEnumDescriptor(desc *descriptor.EnumDescriptorProto, parent *Descriptor, file *descriptor.FileDescriptorProto, index int) *EnumDescriptor {
	ed := &EnumDescriptor{
		common:              common{file},
		EnumDescriptorProto: desc,
		parent:              parent,
		index:               index,
	}
	if parent == nil {
		ed.path = fmt.Sprintf("%d,%d", enumPath, index)
	} else {
		ed.path = fmt.Sprintf("%s,%d,%d", parent.path, messageEnumPath, index)
	}
	return ed
}

// Return a slice of all the EnumDescriptors defined within this file
func wrapEnumDescriptors(file *descriptor.FileDescriptorProto, descs []*Descriptor) []*EnumDescriptor {
	sl := make([]*EnumDescriptor, 0, len(file.EnumType)+10)
	// Top-level enums.
	for i, enum := range file.EnumType {
		sl = append(sl, newEnumDescriptor(enum, nil, file, i))
	}
	// Enums within messages. Enums within embedded messages appear in the outer-most message.
	for _, nested := range descs {
		for i, enum := range nested.EnumType {
			sl = append(sl, newEnumDescriptor(enum, nested, file, i))
		}
	}
	return sl
}

func buildNestedEnums(descs []*Descriptor, enums []*EnumDescriptor) {
	for _, desc := range descs {
		if len(desc.EnumType) != 0 {
			for _, enum := range enums {
				if enum.parent == desc {
					desc.enums = append(desc.enums, enum)
				}
			}
			if len(desc.enums) != len(desc.EnumType) {
				Fail("internal error: enum nesting failure for", desc.GetName())
			}
		}
	}
}

// TypeName returns the elements of the dotted type name.
// The package name is not part of this name.
func (e *EnumDescriptor) TypeName() (s []string) {
	if e.typename != nil {
		return e.typename
	}
	name := e.GetName()
	if e.parent == nil {
		s = make([]string, 1)
	} else {
		pname := e.parent.TypeName()
		s = make([]string, len(pname)+1)
		copy(s, pname)
	}
	s[len(s)-1] = name
	e.typename = s
	return s
}

// ExtensionDescriptor describes an extension. If it's at top level, its parent will be nil.
// Otherwise it will be the descriptor of the message in which it is defined.
type ExtensionDescriptor struct {
	common
	*descriptor.FieldDescriptorProto
	parent *Descriptor // The containing message, if any.
}

// Return a slice of all the top-level ExtensionDescriptors defined within this file.
func wrapExtensions(file *descriptor.FileDescriptorProto) []*ExtensionDescriptor {
	var sl []*ExtensionDescriptor
	for _, field := range file.Extension {
		sl = append(sl, &ExtensionDescriptor{common{file}, field, nil})
	}
	return sl
}

// TypeName returns the elements of the dotted type name.
// The package name is not part of this name.
func (e *ExtensionDescriptor) TypeName() (s []string) {
	name := e.GetName()
	if e.parent == nil {
		// top-level extension
		s = make([]string, 1)
	} else {
		pname := e.parent.TypeName()
		s = make([]string, len(pname)+1)
		copy(s, pname)
	}
	s[len(s)-1] = name
	return s
}

// DescName returns the variable name used for the generated descriptor.
func (e *ExtensionDescriptor) DescName() string {
	// The full type name.
	typeName := e.TypeName()
	// Each scope of the extension is individually CamelCased, and all are joined with "_" with an "E_" prefix.
	for i, s := range typeName {
		typeName[i] = stringutils.CamelCase(s)
	}
	return "E_" + strings.Join(typeName, "_")
}

// ImportedDescriptor describes a type that has been publicly imported from another file.
type ImportedDescriptor struct {
	common
	Object Object
}

// Return a slice of all the types that are publicly imported into this file.
func wrapImported(file *descriptor.FileDescriptorProto, fileMap map[string]*FileDescriptor) (sl []*ImportedDescriptor) {
	for _, index := range file.PublicDependency {
		df := fileMap[file.Dependency[index]]
		for _, d := range df.Descriptors {
			if d.GetOptions().GetMapEntry() {
				continue
			}
			sl = append(sl, &ImportedDescriptor{common{file}, d})
		}
		for _, e := range df.Enums {
			sl = append(sl, &ImportedDescriptor{common{file}, e})
		}
		for _, ext := range df.Extensions {
			sl = append(sl, &ImportedDescriptor{common{file}, ext})
		}
	}
	return
}

// TypeName ...
func (id *ImportedDescriptor) TypeName() []string { return id.Object.TypeName() }

// ServiceDescriptor represents a protocol buffer service.
type ServiceDescriptor struct {
	common
	*descriptor.ServiceDescriptorProto
	Methods []*MethodDescriptor
	Index   int // index of the ServiceDescriptorProto in its parent FileDescriptorProto

	Path string // The SourceCodeInfo path as comma-separated integers.
}

// TypeName ...
func (sd *ServiceDescriptor) TypeName() []string {
	return []string{sd.GetName()}
}

func wrapServices(file *descriptor.FileDescriptorProto) (sl []*ServiceDescriptor) {
	for i, svc := range file.Service {
		sd := &ServiceDescriptor{
			common:                 common{file},
			ServiceDescriptorProto: svc,
			Index: i,
			Path:  fmt.Sprintf("%d,%d", servicePath, i),
		}
		for j, method := range svc.Method {
			md := &MethodDescriptor{
				common:                common{file},
				MethodDescriptorProto: method,
				service:               sd,
				Path:                  fmt.Sprintf("%d,%d,%d,%d", servicePath, i, serviceMethodPath, j),
			}
			sd.Methods = append(sd.Methods, md)
		}
		sl = append(sl, sd)
	}
	return sl
}

// MethodDescriptor represents an RPC method on a protocol buffer
// service.
type MethodDescriptor struct {
	common
	*descriptor.MethodDescriptorProto
	service *ServiceDescriptor

	Path string // The SourceCodeInfo path as comma-separated integers.
}

// TypeName ...
func (md *MethodDescriptor) TypeName() []string {
	return []string{md.service.GetName(), md.GetName()}
}

// FileDescriptor describes an protocol buffer descriptor file (.proto).
// It includes slices of all the messages and enums defined within it.
// Those slices are constructed by WrapTypes.
type FileDescriptor struct {
	*descriptor.FileDescriptorProto
	Descriptors []*Descriptor          // All the messages defined in this file.
	Enums       []*EnumDescriptor      // All the enums defined in this file.
	Extensions  []*ExtensionDescriptor // All the top-level extensions defined in this file.
	Imported    []*ImportedDescriptor  // All types defined in files publicly imported by this file.
	Services    []*ServiceDescriptor   // All the services defined in this file.

	// Comments, stored as a map of path (comma-separated integers) to the comment.
	Comments map[string]*descriptor.SourceCodeInfo_Location

	Index int // The index of this file in the list of files to generate code for

	proto3 bool // whether to generate proto3 code for this file
}

// VarName is the variable name used in generated code to refer to the
// compressed bytes of this descriptor. It is not exported, so it is only valid
// inside the generated package.
//
// protoc-gen-go writes its own version of this file, but so does
// protoc-gen-gogo - with a different name! Twirp aims to be compatible with
// both; the simplest way forward is to write the file descriptor again as
// another variable that we control.
func (d *FileDescriptor) VarName() string { return fmt.Sprintf("twirpFileDescriptor%d", d.Index) }

// PackageComments get pkg comments
func (d *FileDescriptor) PackageComments() string {
	if loc, ok := d.Comments[strconv.Itoa(packagePath)]; ok {
		text := strings.TrimSuffix(loc.GetLeadingComments(), "\n")
		return text
	}
	return ""
}

// BaseFileName name strip extension
func (d *FileDescriptor) BaseFileName() string {
	name := *d.Name
	if ext := path.Ext(name); ext == ".proto" || ext == ".protodevel" {
		name = name[:len(name)-len(ext)]
	}

	return name
}

func extractComments(file *FileDescriptor) {
	file.Comments = make(map[string]*descriptor.SourceCodeInfo_Location)
	for _, loc := range file.GetSourceCodeInfo().GetLocation() {
		if loc.LeadingComments == nil {
			continue
		}
		var p []string
		for _, n := range loc.Path {
			p = append(p, strconv.Itoa(int(n)))
		}
		file.Comments[strings.Join(p, ",")] = loc
	}
}

// Object is an interface abstracting the abilities shared by enums, messages, extensions and imported objects.
type Object interface {
	TypeName() []string
	File() *descriptor.FileDescriptorProto
}

// The SourceCodeInfo message describes the location of elements of a parsed
// .proto file by way of a "path", which is a sequence of integers that
// describe the route from a FileDescriptorProto to the relevant submessage.
// The path alternates between a field number of a repeated field, and an index
// into that repeated field. The constants below define the field numbers that
// are used.
//
// See descriptor.proto for more information about this.
const (
	// tag numbers in FileDescriptorProto
	packagePath = 2 // package
	messagePath = 4 // message_type
	enumPath    = 5 // enum_type
	servicePath = 6 // service
	// tag numbers in DescriptorProto
	//messageFieldPath   = 2 // field
	messageMessagePath = 3 // nested_type
	messageEnumPath    = 4 // enum_type
	//messageOneofPath   = 8 // oneof_decl
	// tag numbers in ServiceDescriptorProto
	//serviceNamePath    = 1 // name
	serviceMethodPath = 2 // method
	//serviceOptionsPath = 3 // options
	// tag numbers in MethodDescriptorProto
	//methodNamePath   = 1 // name
	//methodInputPath  = 2 // input_type
	//methodOutputPath = 3 // output_type
)
