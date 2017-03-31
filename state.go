package main

const (
	baseType = iota
	pathType
	programType
)

type State struct {
	Type    int
	Base    bool
	Path    string
	Program int
}

func NewState() *State {
	return &State{
		Type:    baseType,
		Base:    true,
		Path:    "",
		Program: -1,
	}
}

func (s *State) IsBase() bool {
	return s.Type == baseType
}

func (s *State) IsPath() bool {
	return s.Type == pathType
}

func (s *State) IsProgram() bool {
	return s.Type == programType
}

func (s *State) SetBase() {
	s.Type = baseType
	s.Base = true
	s.Path = ""
	s.Program = -1
}

func (s *State) SetPath(path string) {
	s.Type = pathType
	s.Base = false
	s.Path = path
}

func (s *State) UnsetPath() {
	s.Type = baseType
	s.Base = true
}

func (s *State) SetProgram(program int) {
	s.Type = programType
	s.Base = false
	s.Program = program
}

func (s *State) UnsetProgram() {
	s.Type = pathType
	s.Program = -1
}
