package domainops

import "github.com/tnnmigga/corev2/iface"

type IRoot interface {
	PutCase(index int, useCase any)
	GetCase(index int) any
}

func Register[I any](root IRoot, index int, useCase I) {
	root.PutCase(index, useCase)
}

type root struct {
	iface.IModule
	cases []any
}

func New(m iface.IModule, maxCaseIndex int) IRoot {
	return &root{
		IModule: m,
		cases:   make([]any, maxCaseIndex),
	}
}

func (p *root) PutCase(caseIndex int, useCase any) {
	p.cases[caseIndex] = useCase
}

func (p *root) GetCase(caseIndex int) any {
	return p.cases[caseIndex]
}
