package cbor

import "errors"

type MajorType byte

const (
	MajorTypeUInt        MajorType = 0 << 5
	MajorTypeNInt        MajorType = 1 << 5
	MajorTypeBstr        MajorType = 2 << 5
	MajorTypeTstr        MajorType = 3 << 5
	MajorTypeArray       MajorType = 4 << 5
	MajorTypeMap         MajorType = 5 << 5
	MajorTypeTagged      MajorType = 6 << 5
	MajorTypeSimpleFloat MajorType = 7 << 5
)

const (
	MajorTypeMask byte = 0b111_00000
)

type Arg byte

const (
	Arg8  Arg = 24
	Arg16 Arg = 25
	Arg32 Arg = 26
	Arg64 Arg = 27
	// 28..30 reserved
	ArgIndefinite Arg = 31
)

const (
	argMask byte = 0b000_11111
)

const (
	SimpleFalse     Arg = 20
	SimpleTrue      Arg = 21
	SimpleNull      Arg = 22
	SimpleUndefined Arg = 23
	SimpleUint8         = Arg8  // 24
	SimpleFloat16       = Arg16 // 25
	SimpleFloat32       = Arg32 // 26
	SimpleFloat64       = Arg64 // 27
	// 28..30 reserved
	SimpleBreak Arg = 31
)

var (
	ErrUnsupportedMajorType = errors.New("cbor: unsupported major type")
	ErrUnsupportedValue     = errors.New("cbor: unsupported value")
	ErrNotWellFormed        = errors.New("cbor: not well formed")
	ErrOverflow             = errors.New("cbor: overflow")
	ErrNestedIndefinite     = errors.New("cbor: nested indefinite")
)

const (
	valueBreak = byte(MajorTypeSimpleFloat) | byte(SimpleBreak)
)

const lenSharedBuffer = 64

var sharedBuffer [lenSharedBuffer]byte
