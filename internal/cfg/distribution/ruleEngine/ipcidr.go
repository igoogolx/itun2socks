package ruleEngine

import (
	"github.com/igoogolx/itun2socks/internal/constants"
	"net/netip"
)

type IpCidr struct {
	RuleType constants.RuleConfig
	Payload  string
	prefix   netip.Prefix
	policy   string
}

func (d IpCidr) Policy() constants.RuleType {
	return constants.RuleType(d.policy)
}

func (i IpCidr) Type() constants.RuleConfig {
	return constants.RuleIpCidr
}

func (i IpCidr) Match(value string) bool {
	ip, err := netip.ParseAddr(value)
	if err != nil {
		return false
	}
	return i.prefix.Contains(ip)
}

func (i IpCidr) Value() string {
	return i.prefix.String()
}

func NewIpCidrRule(payload string, policy string) (*IpCidr, error) {
	prefix, err := netip.ParsePrefix(payload)
	if err != nil {
		return nil, err
	}

	return &IpCidr{constants.RuleIpCidr, payload, prefix, policy}, nil
}
