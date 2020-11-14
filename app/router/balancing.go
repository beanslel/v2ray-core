// +build !confonly

package router

import (
	"v2ray.com/core/features/outbound"
)

type BalancingStrategy interface {
	PickOutbound([]string) string
}

type Balancer struct {
	selectors []string
	strategy  BalancingStrategy
	ohm       outbound.Manager
}

func (b *Balancer) PickOutbound() (string, error) {
	hs, ok := b.ohm.(outbound.HandlerSelector)
	if !ok {
		return "", newError("outbound.Manager is not a HandlerSelector")
	}
	tags := hs.Select(b.selectors)
	if len(tags) == 0 {
		return "", newError("no available outbounds selected")
	}
	tag := b.strategy.PickOutbound(tags)
	if tag == "" {
		return "", newError("balancing strategy returns empty tag")
	}
	return tag, nil
}
