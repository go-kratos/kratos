package typemap

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

import (
	"strings"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/pkg/errors"
)

// Registry is the place of descriptors resolving
type Registry struct {
	allFiles    []*descriptor.FileDescriptorProto
	filesByName map[string]*descriptor.FileDescriptorProto

	// Mapping of fully-qualified names to their definitions
	messagesByProtoName map[string]*MessageDefinition
}

// New Registry
func New(files []*descriptor.FileDescriptorProto) *Registry {
	r := &Registry{
		allFiles:            files,
		filesByName:         make(map[string]*descriptor.FileDescriptorProto),
		messagesByProtoName: make(map[string]*MessageDefinition),
	}

	// First, index the file descriptors by name. We need this so
	// messageDefsForFile can correctly scan imports.
	for _, f := range files {
		r.filesByName[f.GetName()] = f
	}

	// Next, index all the message definitions by their fully-qualified proto
	// names.
	for _, f := range files {
		defs := messageDefsForFile(f, r.filesByName)
		for name, def := range defs {
			r.messagesByProtoName[name] = def
		}
	}
	return r
}

// FileComments comment of file
func (r *Registry) FileComments(file *descriptor.FileDescriptorProto) (DefinitionComments, error) {
	return commentsAtPath([]int32{packagePath}, file), nil
}

// ServiceComments comments of service
func (r *Registry) ServiceComments(file *descriptor.FileDescriptorProto, svc *descriptor.ServiceDescriptorProto) (DefinitionComments, error) {
	for i, s := range file.Service {
		if s == svc {
			path := []int32{servicePath, int32(i)}
			return commentsAtPath(path, file), nil
		}
	}
	return DefinitionComments{}, errors.Errorf("service not found in file")
}

// FieldComments ...
func (r *Registry) FieldComments(message *MessageDefinition, field *descriptor.FieldDescriptorProto) (DefinitionComments, error) {
	file := message.File
	mpath := message.path
	for i, f := range message.Descriptor.Field {
		if f == field {
			path := append(mpath, messageFieldPath, int32(i))
			return commentsAtPath(path, file), nil
		}
	}
	return DefinitionComments{}, errors.Errorf("field not found in msg")
}

// MethodComments comment of method
func (r *Registry) MethodComments(file *descriptor.FileDescriptorProto, svc *descriptor.ServiceDescriptorProto, method *descriptor.MethodDescriptorProto) (DefinitionComments, error) {
	for i, s := range file.Service {
		if s == svc {
			path := []int32{servicePath, int32(i)}
			for j, m := range s.Method {
				if m == method {
					path = append(path, serviceMethodPath, int32(j))
					return commentsAtPath(path, file), nil
				}
			}
		}
	}
	return DefinitionComments{}, errors.Errorf("service not found in file")
}

// MethodInputDefinition returns MethodInputDefinition
func (r *Registry) MethodInputDefinition(method *descriptor.MethodDescriptorProto) *MessageDefinition {
	return r.messagesByProtoName[method.GetInputType()]
}

// MethodOutputDefinition returns MethodOutputDefinition
func (r *Registry) MethodOutputDefinition(method *descriptor.MethodDescriptorProto) *MessageDefinition {
	return r.messagesByProtoName[method.GetOutputType()]
}

// MessageDefinition by name
func (r *Registry) MessageDefinition(name string) *MessageDefinition {
	return r.messagesByProtoName[name]
}

// MessageDefinition msg info
type MessageDefinition struct {
	// Descriptor is is the DescriptorProto defining the message.
	Descriptor *descriptor.DescriptorProto
	// File is the File that the message was defined in. Or, if it has been
	// publicly imported, what File was that import performed in?
	File *descriptor.FileDescriptorProto
	// Parent is the parent message, if this was defined as a nested message. If
	// this was defiend at the top level, parent is nil.
	Parent *MessageDefinition
	// Comments describes the comments surrounding a message's definition. If it
	// was publicly imported, then these comments are from the actual source file,
	// not the file that the import was performed in.
	Comments DefinitionComments

	// path is the 'SourceCodeInfo' path. See the documentation for
	// github.com/golang/protobuf/protoc-gen-go/descriptor.SourceCodeInfo for an
	// explanation of its format.
	path []int32
}

// ProtoName returns the dot-delimited, fully-qualified protobuf name of the
// message.
func (m *MessageDefinition) ProtoName() string {
	prefix := "."
	if pkg := m.File.GetPackage(); pkg != "" {
		prefix += pkg + "."
	}

	if lineage := m.Lineage(); len(lineage) > 0 {
		for _, parent := range lineage {
			prefix += parent.Descriptor.GetName() + "."
		}
	}

	return prefix + m.Descriptor.GetName()
}

// Lineage returns m's parental chain all the way back up to a top-level message
// definition. The first element of the returned slice is the highest-level
// parent.
func (m *MessageDefinition) Lineage() []*MessageDefinition {
	var parents []*MessageDefinition
	for p := m.Parent; p != nil; p = p.Parent {
		parents = append([]*MessageDefinition{p}, parents...)
	}
	return parents
}

// descendants returns all the submessages defined within m, and all the
// descendants of those, recursively.
func (m *MessageDefinition) descendants() []*MessageDefinition {
	descendants := make([]*MessageDefinition, 0)
	for i, child := range m.Descriptor.NestedType {
		path := append(m.path, []int32{messageMessagePath, int32(i)}...)
		childDef := &MessageDefinition{
			Descriptor: child,
			File:       m.File,
			Parent:     m,
			Comments:   commentsAtPath(path, m.File),
			path:       path,
		}
		descendants = append(descendants, childDef)
		descendants = append(descendants, childDef.descendants()...)
	}
	return descendants
}

// messageDefsForFile gathers a mapping of fully-qualified protobuf names to
// their definitions. It scans a singles file at a time. It requires a mapping
// of .proto file names to their definitions in order to correctly handle
// 'import public' declarations; this mapping should include all files
// transitively imported by f.
func messageDefsForFile(f *descriptor.FileDescriptorProto, filesByName map[string]*descriptor.FileDescriptorProto) map[string]*MessageDefinition {
	byProtoName := make(map[string]*MessageDefinition)
	// First, gather all the messages defined at the top level.
	for i, d := range f.MessageType {
		path := []int32{messagePath, int32(i)}
		def := &MessageDefinition{
			Descriptor: d,
			File:       f,
			Parent:     nil,
			Comments:   commentsAtPath(path, f),
			path:       path,
		}

		byProtoName[def.ProtoName()] = def
		// Next, all nested message definitions.
		for _, child := range def.descendants() {
			byProtoName[child.ProtoName()] = child
		}
	}

	// Finally, all messages imported publicly.
	for _, depIdx := range f.PublicDependency {
		depFileName := f.Dependency[depIdx]
		depFile := filesByName[depFileName]
		depDefs := messageDefsForFile(depFile, filesByName)
		for _, def := range depDefs {
			imported := &MessageDefinition{
				Descriptor: def.Descriptor,
				File:       f,
				Parent:     def.Parent,
				Comments:   commentsAtPath(def.path, depFile),
				path:       def.path,
			}
			byProtoName[imported.ProtoName()] = imported
		}
	}

	return byProtoName
}

//   // ignored detached comments.
type DefinitionComments struct {
	Leading         string
	Trailing        string
	LeadingDetached []string
}

func commentsAtPath(path []int32, sourceFile *descriptor.FileDescriptorProto) DefinitionComments {
	if sourceFile.SourceCodeInfo == nil {
		// The compiler didn't provide us with comments.
		return DefinitionComments{}
	}

	for _, loc := range sourceFile.SourceCodeInfo.Location {
		if pathEqual(path, loc.Path) {
			return DefinitionComments{
				Leading:         strings.TrimSuffix(loc.GetLeadingComments(), "\n"),
				LeadingDetached: loc.GetLeadingDetachedComments(),
				Trailing:        loc.GetTrailingComments(),
			}
		}
	}
	return DefinitionComments{}
}

func pathEqual(path1, path2 []int32) bool {
	if len(path1) != len(path2) {
		return false
	}
	for i, v := range path1 {
		if path2[i] != v {
			return false
		}
	}
	return true
}

const (
	// tag numbers in FileDescriptorProto
	packagePath = 2 // package
	messagePath = 4 // message_type
	servicePath = 6 // service
	// tag numbers in DescriptorProto
	messageFieldPath   = 2 // field
	messageMessagePath = 3 // nested_type
	// tag numbers in ServiceDescriptorProto
	serviceMethodPath = 2 // method
)
