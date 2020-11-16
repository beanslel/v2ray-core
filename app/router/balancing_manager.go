// +build !confonly

package router

import (
	"sync/atomic"

	"v2ray.com/core/common/dice"
)

type BalancerManager struct {
	Random   *RandomManager
	Fallback *FallbackManager
}

type RandomStrategy struct {
}

type RandomManager struct {
}

type FallbackStrategy struct {
	mode string
}

type FallbackManager struct {
	curIndex       int
	maxAttempts    int64
	failedAttempts int64
}

const defaultMaxAttempts = 10

// temp and will removed
var balancerManager BalancerManager

func newBalancerManager() {
	balancerManager = BalancerManager{
		Random:   NewRandomManager(),
		Fallback: NewFallbackManager(int64(defaultMaxAttempts)),
	}
}

func NewBalancerManager() BalancerManager {
	return balancerManager
}

func NewRandomManager() *RandomManager {
	return &RandomManager{}
}

// NewFallbackManager returns a new instance of FallbackManager
func NewFallbackManager(maxAttempts int64) *FallbackManager {
	return &FallbackManager{
		curIndex:       0,
		failedAttempts: int64(0),
		maxAttempts:    maxAttempts,
	}
}

// GetFailedAttempts implements outbound.FailedAttemptsRecorder
func (m *FallbackManager) getFailedAttempts() int64 {
	return atomic.LoadInt64(&m.failedAttempts)
}
func (m *FallbackManager) GetFailedAttempts() int64 {
	return balancerManager.Fallback.getFailedAttempts()
}

// ResetFailedAttempts implements outbound.FailedAttemptsRecorder
func (m *FallbackManager) resetFailedAttempts() int64 {
	return atomic.SwapInt64(&m.failedAttempts, int64(0))
}
func (m *FallbackManager) ResetFailedAttempts() int64 {
	return balancerManager.Fallback.resetFailedAttempts()
}

// AddFailedAttempts implements outbound.FailedAttemptsRecorder
func (m *FallbackManager) addFailedAttempts() int64 {
	return atomic.AddInt64(&m.failedAttempts, int64(1))
}
func (m *FallbackManager) AddFailedAttempts() int64 {
	return balancerManager.Fallback.addFailedAttempts()
}

// PickOutbound picks an outbound with fallback strategy
func (m *FallbackManager) pickOutbound(tags []string) string {
	if m.failedAttempts > m.maxAttempts {
		m.ResetFailedAttempts()
		m.curIndex = (m.curIndex + 1) % len(tags)
		newError("balancer: switched to fallback " + tags[m.curIndex]).AtInfo().WriteToLog()
	}
	return tags[m.curIndex]
}
func (m *FallbackManager) PickOutbound(tags []string) string {
	return balancerManager.Fallback.pickOutbound(tags)
}

func (m *RandomManager) PickOutbound(tags []string) string {
	n := len(tags)
	if n == 0 {
		panic("0 tags")
	}

	return tags[dice.Roll(n)]
}

func init() {
	newBalancerManager()
}
