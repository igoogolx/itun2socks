package ruleEngine

import (
	"github.com/igoogolx/itun2socks/internal/constants"
	"github.com/stretchr/testify/assert"
	"testing"
)

var BYPASS_IP = "119.29.29.29"
var BYPASS_DOMAIN = "bing.cn"

var PROXY_IP = "8.8.8.8"
var PROXY_DOMAIN = "google.com"

func TestProxyAll(t *testing.T) {
	engine, err := New("proxy_all", []string{})
	assert.NoError(t, err)

	//not found
	_, err = engine.Match(BYPASS_IP, constants.IpRuleTypes)
	assert.Equal(t, ErrNotFound, err)

	//not found
	_, err = engine.Match(BYPASS_DOMAIN, constants.DomainRuleTypes)
	assert.Equal(t, ErrNotFound, err)

	//not found
	_, err = engine.Match(PROXY_IP, constants.IpRuleTypes)
	assert.Equal(t, ErrNotFound, err)

	//not found
	_, err = engine.Match(PROXY_DOMAIN, constants.DomainRuleTypes)
	assert.Equal(t, ErrNotFound, err)

}

func TestBypassAll(t *testing.T) {
	engine, err := New("bypass_all", []string{})
	assert.NoError(t, err)

	//bypass
	bypassIpRes, err := engine.Match(BYPASS_IP, constants.IpRuleTypes)
	assert.NoError(t, err)
	assert.Equal(t, bypassIpRes.GetPolicy(), constants.PolicyDirect)

	//bypass
	bypassDomainRes, err := engine.Match(BYPASS_DOMAIN, constants.DomainRuleTypes)
	assert.NoError(t, err)
	assert.Equal(t, bypassDomainRes.GetPolicy(), constants.PolicyDirect)

	//bypass
	proxyIpRes, err := engine.Match(PROXY_IP, constants.IpRuleTypes)
	assert.NoError(t, err)
	assert.Equal(t, proxyIpRes.GetPolicy(), constants.PolicyDirect)

	//bypass
	proxyDomainRes, err := engine.Match(PROXY_DOMAIN, constants.DomainRuleTypes)
	assert.NoError(t, err)
	assert.Equal(t, proxyDomainRes.GetPolicy(), constants.PolicyDirect)
}

func TestBypassCn(t *testing.T) {
	engine, err := New("bypass_cn", []string{})
	assert.NoError(t, err)

	//bypass
	bypassIpRes, err := engine.Match(BYPASS_IP, constants.IpRuleTypes)
	assert.NoError(t, err)
	assert.Equal(t, bypassIpRes.GetPolicy(), constants.PolicyDirect)

	//bypass
	bypassDomainRes, err := engine.Match(BYPASS_DOMAIN, constants.DomainRuleTypes)
	assert.NoError(t, err)
	assert.Equal(t, bypassDomainRes.GetPolicy(), constants.PolicyDirect)

	//not found
	_, err = engine.Match(PROXY_IP, constants.IpRuleTypes)
	assert.Equal(t, ErrNotFound, err)

	//not found
	_, err = engine.Match(PROXY_DOMAIN, constants.DomainRuleTypes)
	assert.Equal(t, ErrNotFound, err)

}

func TestProxyGfw(t *testing.T) {
	engine, err := New("proxy_gfw", []string{})
	assert.NoError(t, err)

	//bypass
	bypassIpRes, err := engine.Match(BYPASS_IP, constants.IpRuleTypes)
	assert.NoError(t, err)
	assert.Equal(t, bypassIpRes.GetPolicy(), constants.PolicyDirect)

	//bypass
	bypassDomainRes, err := engine.Match(BYPASS_DOMAIN, constants.DomainRuleTypes)
	assert.NoError(t, err)
	assert.Equal(t, bypassDomainRes.GetPolicy(), constants.PolicyDirect)

	//bypass
	proxyIpRes, err := engine.Match(PROXY_IP, constants.IpRuleTypes)
	assert.NoError(t, err)
	assert.Equal(t, proxyIpRes.GetPolicy(), constants.PolicyDirect)

	//proxy
	proxyDomainRes, err := engine.Match(PROXY_DOMAIN, constants.DomainRuleTypes)
	assert.NoError(t, err)
	assert.Equal(t, proxyDomainRes.GetPolicy(), constants.PolicyProxy)
}
