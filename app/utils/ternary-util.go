package utils

func If[T any](cond bool, vtrue, vfalse T) T {
	if cond {
		return vtrue
	}
	return vfalse
}

func IfNil[T any](v *T, d T) T {
	if v == nil {
		return d
	}
	return *v
}

func Default[T any](v *T, d T) T {
	if v == nil {
		return d 
	}
	return *v
}

func VPtr[T any](s T) *T {
	return &s
}
