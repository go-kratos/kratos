package types

// ServiceSpec service spec
type ServiceSpec struct {
	Name string
	// origin service package name
	Package string
	// service import path
	ImportPath string
	Receiver   string
	Methods    []*Method
}

// Method service method
type Method struct {
	Name       string // method name
	Comments   []string
	Options    []string
	Parameters []*Field
	Results    []*Field
}

// Field a pair of parameter name and type
type Field struct {
	Name string
	Type Typer
}

func (f *Field) String() string {
	if f.Name == "" {
		return f.Type.String()
	}
	return f.Name + " " + f.Type.String()
}

// Typer type interface
type Typer interface {
	String() string
	IsReference() bool
	SetReference() Typer
}

// BasicType go buildin type
type BasicType struct {
	Name      string
	Reference bool
}

// IsReference return is reference
func (t *BasicType) IsReference() bool {
	return t.Reference
}

// SetReference return is reference
func (t *BasicType) SetReference() Typer {
	t.Reference = true
	return t
}

func (t *BasicType) String() string {
	return t.Name
}

// ArrayType array type
type ArrayType struct {
	EltType   Typer
	Reference bool
}

func (t *ArrayType) String() string {
	return "[]" + t.EltType.String()
}

// IsReference return is reference
func (t *ArrayType) IsReference() bool {
	return t.Reference
}

// SetReference return is reference
func (t *ArrayType) SetReference() Typer {
	t.Reference = true
	return t
}

// MapType map
type MapType struct {
	KeyType   Typer
	ValueType Typer
	Reference bool
}

// IsReference return is reference
func (t *MapType) IsReference() bool {
	return t.Reference
}

func (t *MapType) String() string {
	return "[" + t.KeyType.String() + "]" + t.ValueType.String()
}

// SetReference return is reference
func (t *MapType) SetReference() Typer {
	t.Reference = true
	return t
}

// StructType struct type
type StructType struct {
	Package    string
	ImportPath string
	IdentName  string
	Reference  bool
	Fields     []*Field
	ProtoFile  string
}

func (t *StructType) String() string {
	return t.Package + "." + t.IdentName
}

// IsReference return is reference
func (t *StructType) IsReference() bool {
	return t.Reference
}

// SetReference return is reference
func (t *StructType) SetReference() Typer {
	t.Reference = true
	return t
}

// InterfaceType struct type
type InterfaceType struct {
	Package    string
	ImportPath string
	IdentName  string
	Reference  bool
}

func (t *InterfaceType) String() string {
	return t.Package + "." + t.IdentName
}

// IsReference return is reference
func (t *InterfaceType) IsReference() bool {
	return t.Reference
}

// SetReference return is reference
func (t *InterfaceType) SetReference() Typer {
	t.Reference = true
	return t
}
