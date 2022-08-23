package types

import "fmt"

type Chain struct {
	name    string
	cluster string
}

func NewChain(name string, cluster string) Chain {
	return Chain{
		name:    name,
		cluster: cluster,
	}
}

func (c *Chain) ID() string {
	return fmt.Sprintf("%s.%s", c.name, c.cluster)
}
