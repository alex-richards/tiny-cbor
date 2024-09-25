package cbor

import "errors"

type MajorType byte

const (
	MajorTypeUInt        MajorType = 0 << 5
	MajorTypeNInt                  = 1 << 5
	MajorTypeBstr                  = 2 << 5
	MajorTypeTstr                  = 3 << 5
	MajorTypeArray                 = 4 << 5
	MajorTypeMap                   = 5 << 5
	MajorTypeTagged                = 6 << 5
	MajorTypeSimpleFloat           = 7 << 5
)

const (
	majorTypeMask = 0b111_00000
)

type Arg byte

const (
	Arg8  Arg = 24
	Arg16     = 25
	Arg32     = 26
	Arg64     = 27
	// 28..30 reserved
	ArgIndefinite = 31
)

const (
	argMask = 0b000_11111
)

const (
	SimpleFalse     byte = 20
	SimpleTrue           = 21
	SimpleNull           = 22
	SimpleUndefined      = 23
	SimpleUint8          = Arg8  // 24
	SimpleFloat16        = Arg16 // 25
	SimpleFloat32        = Arg32 // 26
	SimpleFloat64        = Arg64 // 27
	// 28..30 reserved
	SimpleBreak = 31
)

var (
	ErrUnsupportedMajorType = errors.New("cbor: unsupported major type")
	ErrUnsupportedValue     = errors.New("cbor: unsupported value")
	ErrNotWellFormed        = errors.New("cbor: not well formed")
	ErrOverflow             = errors.New("cbor: overflow")
	ErrNestedIndefinite     = errors.New("cbor: nested indefinite")
)

const (
	valueBreak = MajorTypeSimpleFloat | SimpleBreak
)

const lenSharedBuffer = 64

var sharedBuffer [lenSharedBuffer]byte
