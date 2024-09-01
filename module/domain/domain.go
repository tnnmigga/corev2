package domain

import "github.com/tnnmigga/corev2/iface"

type Root interface {
	PutCase(index int, useCase any)
	GetCase(index int) any
}

type root struct {
	iface.IModule
	useCases []any
}

func New(m iface.IModule, maxCaseIndex int) Root {
	return &root{
		IModule:  m,
		useCases: make([]any, maxCaseIndex),
	}
}

func (p *root) PutCase(caseIndex int, useCase any) {
	p.useCases[caseIndex] = useCase
}

func (p *root) GetCase(caseIndex int) any {
	return p.useCases[caseIndex]
}
