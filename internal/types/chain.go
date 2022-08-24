package types

import "fmt"

type Chain struct {
	name    string
	chainID string
}

func NewChain(name string, chainID string) Chain {
	return Chain{
		name:    name,
		chainID: chainID,
	}
}

func (c Chain) ID() string {
	return fmt.Sprintf("%s.%s", c.name, c.chainID)
}
