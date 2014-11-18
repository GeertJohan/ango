package definitions

import (
	"strings"
)

// TypeCategory is the category to which a type belongs.
// The different categories are defined as constants.
type TypeCategory int

const (
	// Builtin types (int, int8, string, bool, etc)
	Builtin = TypeCategory(iota + 1)

	// Simple types. eg: `type myInt int`
	Simple

	// Slice types. eg: `type mySlice []int`
	Slice

	// Map types. eg: `type myMap map[keyType]valueType`
	Map

	// Struct types. eg: `type myStruct struct{}`
	Struct
)

// StructField defines a single field in a struct
type StructField struct {
	Name string
	Type *Type
}

// Type is the type of a parameter
// It's value should be a valid go TypeName (http://golang.org/ref/spec#TypeName)
type Type struct {
	// Name is the identifier of the type
	Name string

	// Category indicates the Type's category (builtin, simple, slice, map, struct)
	Category TypeCategory

	// SimpleType defines the type for the type that this type maps to
	SimpleType *Type

	// SliceElementType defines the type for the slice elements
	SliceElementType *Type

	// MapKeyType and MapValueType define the key and value types for the map
	MapKeyType   *Type
	MapValueType *Type

	// StructFields holds the struct field definitions, only used when Category is Struct.
	StructFields []StructField
}

func (t *Type) CapitalizedName() string {
	return strings.ToUpper(t.Name[:1]) + t.Name[1:]
}

func (t *Type) GoType() string {
	switch t.Category {
	case Builtin:
		return t.Name
	case Simple:
		return t.SimpleType.CapitalizedName()
	case Slice:
		return `[]` + t.SliceElementType.CapitalizedName()
	case Map:
		return `map[` + t.MapKeyType.CapitalizedName() + `]` + t.MapValueType.CapitalizedName()
	case Struct:
		s := "struct {\n"
		for _, f := range t.StructFields {
			s += strings.ToUpper(f.Name[:1]) + f.Name[1:] + ` ` + f.Type.GoType() + "\n"
		}
		s += `}`
		return s
	default:
		panic("unknown type")
	}
}

func (t *Type) GoTypeName() string {
	switch t.Category {
	case Builtin:
		return t.Name
	case Simple, Slice, Map:
		return t.CapitalizedName()
	case Struct:
		return `*` + t.CapitalizedName()
	default:
		panic("unknown type")
	}
}

// Builtin types
var (
	TypeInt      = &Type{Name: "int", Category: Builtin}
	TypeInt8     = &Type{Name: "int8", Category: Builtin}
	TypeInt16    = &Type{Name: "int16", Category: Builtin}
	TypeInt32    = &Type{Name: "int32", Category: Builtin}
	TypeInt64    = &Type{Name: "int64", Category: Builtin}
	TypeUint     = &Type{Name: "uint", Category: Builtin}
	TypeUint8    = &Type{Name: "uint8", Category: Builtin}
	TypeUint16   = &Type{Name: "uint16", Category: Builtin}
	TypeUint32   = &Type{Name: "uint32", Category: Builtin}
	TypeUint64   = &Type{Name: "uint64", Category: Builtin}
	TypeString   = &Type{Name: "string", Category: Builtin}
	TypeBool     = &Type{Name: "bool", Category: Builtin}
	BuiltinTypes = map[string]*Type{
		TypeInt.Name:    TypeInt,
		TypeInt8.Name:   TypeInt8,
		TypeInt16.Name:  TypeInt16,
		TypeInt32.Name:  TypeInt32,
		TypeInt64.Name:  TypeInt64,
		TypeUint.Name:   TypeUint,
		TypeUint8.Name:  TypeUint8,
		TypeUint16.Name: TypeUint16,
		TypeUint32.Name: TypeUint32,
		TypeUint64.Name: TypeUint64,
		TypeString.Name: TypeString,
		TypeBool.Name:   TypeBool,
	}
)
