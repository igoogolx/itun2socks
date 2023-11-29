package ruleEngine

import (
	"github.com/igoogolx/itun2socks/internal/constants"
	"net/netip"
)

type IpCidr struct {
	RuleType constants.RuleConfig `json:"ruleType"`
	Payload  string               `json:"payload"`
	prefix   netip.Prefix
	Policy   string `json:"policy"`
}

func (d IpCidr) GetPolicy() constants.RuleType {
	return constants.RuleType(d.Policy)
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
