package definitions

// TypeCategory is the category to which a type belongs.
// The different categories are defined as constants.
type TypeCategory int

const (
	// Builtin types (int, int8, string, bool, etc)
	Builtin = TypeCategory(iota)

	// Simple types. eg: `type myInt int`
	Simple

	// Slice types. eg: `type mySlice []int`
	Slice

	// Map types. eg: `type myMap map[keyType]valueType`
	Map

	// Struct types. eg: `type myStruct struct{}`
	Struct
)

// Type is the type of a parameter
// It's value should be a valid go TypeName (http://golang.org/ref/spec#TypeName)
type Type struct {
	// Name is the identifier of the type
	Name string

	// Category indicates the Type's category (builtin, simple, slice, map, struct)
	Category TypeCategory

	// Source tells where the given type was declared
	Source Source

	// SliceElementType defines the type for the slice elements
	SliceElementType *Type

	// MapKeyType and MapValueType define the key and value types for the map
	MapKeyType   *Type
	MapValueType *Type

	// StructFields holds the struct field definitions, only used when Category is Struct.
	StructFields []StructField
}

// StructField defines a single field in a struct
type StructField struct {
	Name string
	Type *Type
}

// Builtin types
var (
	TypeInt    = &Type{Name: "int"}
	TypeInt8   = &Type{Name: "int8"}
	TypeInt16  = &Type{Name: "int16"}
	TypeInt32  = &Type{Name: "int32"}
	TypeInt64  = &Type{Name: "int64"}
	TypeUint   = &Type{Name: "uint"}
	TypeUint8  = &Type{Name: "uint8"}
	TypeUint16 = &Type{Name: "uint16"}
	TypeUint32 = &Type{Name: "uint32"}
	TypeUint64 = &Type{Name: "uint64"}
	TypeString = &Type{Name: "string"}
	TypeBool   = &Type{Name: "bool"}
)
