package thsdk

// Opt is an optional request field: a pointer whose nil value means "not sent".
// An alias is required, because a defined pointer type cannot carry methods and
// would need explicit conversions at every call site.
//
// The constructors Num, Text and Flag take literals directly:
//
//	thsdk.ListNewsFlashItemsRequest{Traversal: thsdk.FlashPage, Page: thsdk.Num(20)}
type Opt = *any

// Num stores a number or any other value in an optional field.
func Num[T any](value T) Opt {
	boxed := any(value)
	return &boxed
}

// Text stores a string in an optional field.
func Text(value string) Opt { return Num(value) }

// Flag stores a bool in an optional field.
func Flag(value bool) Opt { return Num(value) }
