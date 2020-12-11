package core

type Perm int

const (
	PermRead Perm = 1 << iota
	PermWrite
	PermSystem
)

func (p Perm) System() bool {
	return p&PermSystem > 0
}

func (p Perm) Read() bool {
	return p.System() || p&PermRead > 0
}

func (p Perm) Write() bool {
	return p.System() || p&PermWrite > 0
}
