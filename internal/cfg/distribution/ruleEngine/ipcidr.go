package ruleEngine

import (
	"fmt"
	"github.com/igoogolx/itun2socks/internal/constants"
	"net/netip"
)

type IpCidr struct {
	RuleType constants.RuleType `json:"ruleType"`
	Payload  string             `json:"payload"`
	prefix   netip.Prefix
	Policy   constants.Policy `json:"policy"`
}

func (i IpCidr) GetPolicy() constants.Policy {
	return i.Policy
}

func (i IpCidr) Valid() bool {
	_, err := netip.ParsePrefix(i.Payload)
	if err != nil {
		return false
	}
	return true
}

func (i IpCidr) Type() constants.RuleType {
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

func NewIpCidrRule(payload string, policy constants.Policy) (*IpCidr, error) {
	prefix, err := netip.ParsePrefix(payload)
	if err != nil {
		return nil, fmt.Errorf("invalid ipcidr payload: %s, err:%s", payload, err)
	}

	return &IpCidr{constants.RuleIpCidr, payload, prefix, policy}, nil
}
