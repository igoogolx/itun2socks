package ruleEngine

import (
	"errors"
	"github.com/igoogolx/itun2socks/internal/constants"
	"testing"
)

var BYPASS_IP = "119.29.29.29"
var BYPASS_DOMAIN = "bing.cn"

var PROXY_IP = "8.8.8.8"
var PROXY_DOMAIN = "google.com"

func TestProxyAll(t *testing.T) {
	engine, err := New("proxy_all", []string{})
	if err != nil {
		t.Fatal(err)
	}
	//not found
	_, err = engine.Match(BYPASS_IP, constants.IpRuleTypes)
	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			t.Fatal(err)
		}
	}

	//not found
	_, err = engine.Match(BYPASS_DOMAIN, constants.DomainRuleTypes)
	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			t.Fatal(err)
		}
	}

	//not found
	_, err = engine.Match(PROXY_IP, constants.IpRuleTypes)
	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			t.Fatal(err)
		}
	}

	//not found
	_, err = engine.Match(PROXY_DOMAIN, constants.DomainRuleTypes)
	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			t.Fatal(err)
		}
	}

}

func TestBypassAll(t *testing.T) {
	engine, err := New("bypass_all", []string{})
	if err != nil {
		t.Fatal(err)
	}

	//bypass
	bypassIpRes, err := engine.Match(BYPASS_IP, constants.IpRuleTypes)
	if err != nil {
		t.Fatal(err)
	}

	if bypassIpRes.GetPolicy() != constants.PolicyDirect {
		t.Fatal("bypassIp in bypass_all failed")
	}

	//bypass
	bypassDomainRes, err := engine.Match(BYPASS_DOMAIN, constants.DomainRuleTypes)
	if err != nil {
		t.Fatal(err)
	}
	if bypassDomainRes.GetPolicy() != constants.PolicyDirect {
		t.Fatal("bypassDomain in bypass_all failed")
	}

	//bypass
	proxyIpRes, err := engine.Match(PROXY_IP, constants.IpRuleTypes)
	if err != nil {
		t.Fatal(err)
	}
	if proxyIpRes.GetPolicy() != constants.PolicyDirect {
		t.Fatal("proxyIp in bypass_all failed")
	}

	//bypass
	proxyDomainRes, err := engine.Match(PROXY_DOMAIN, constants.DomainRuleTypes)
	if err != nil {
		t.Fatal(err)
	}
	if proxyDomainRes.GetPolicy() != constants.PolicyDirect {
		t.Fatal("proxyDomain in bypass_all failed")
	}
}

func TestBypassCn(t *testing.T) {
	engine, err := New("bypass_cn", []string{})
	if err != nil {
		t.Fatal(err)
	}

	//bypass
	bypassIpRes, err := engine.Match(BYPASS_IP, constants.IpRuleTypes)
	if err != nil {
		t.Fatal(err)
	}

	if bypassIpRes.GetPolicy() != constants.PolicyDirect {
		t.Fatal("bypassIp in bypass_cn failed")
	}

	//bypass
	bypassDomainRes, err := engine.Match(BYPASS_DOMAIN, constants.DomainRuleTypes)
	if err != nil {
		t.Fatal(err)
	}
	if bypassDomainRes.GetPolicy() != constants.PolicyDirect {
		t.Fatal("bypassDomain in bypass_cn failed")
	}

	//not found
	_, err = engine.Match(PROXY_IP, constants.IpRuleTypes)
	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			t.Fatal(err)
		}
	}

	//not found
	_, err = engine.Match(PROXY_DOMAIN, constants.DomainRuleTypes)
	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			t.Fatal(err)
		}
	}
}

func TestProxyGfw(t *testing.T) {
	engine, err := New("proxy_gfw", []string{})
	if err != nil {
		t.Fatal(err)
	}

	//bypass
	bypassIpRes, err := engine.Match(BYPASS_IP, constants.IpRuleTypes)
	if err != nil {
		t.Fatal(err)
	}
	if bypassIpRes.GetPolicy() != constants.PolicyDirect {
		t.Fatal("bypassIp in proxy_gfw failed")
	}

	//bypass
	bypassDomainRes, err := engine.Match(BYPASS_DOMAIN, constants.DomainRuleTypes)
	if err != nil {
		t.Fatal(err)
	}
	if bypassDomainRes.GetPolicy() != constants.PolicyDirect {
		t.Fatal("bypassDomain in proxy_gfw failed")
	}

	//bypass
	proxyIpRes, err := engine.Match(PROXY_IP, constants.IpRuleTypes)
	if err != nil {
		t.Fatal(err)
	}
	if proxyIpRes.GetPolicy() != constants.PolicyDirect {
		t.Fatal("proxyIpRes in proxy_gfw failed")
	}

	//proxy
	proxyDomainRes, err := engine.Match(PROXY_DOMAIN, constants.DomainRuleTypes)
	if err != nil {
		t.Fatal(err)
	}
	if proxyDomainRes.GetPolicy() != constants.PolicyProxy {
		t.Fatal("proxyDomain in proxy_gfw failed")
	}
}
