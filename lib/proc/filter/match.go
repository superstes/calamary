package filter

import (
	"fmt"
	"net"

	"github.com/superstes/calamary/cnf"
	"github.com/superstes/calamary/proc/meta"
	"github.com/superstes/calamary/proc/parse"
)

func matchProtoL3(pkt parse.ParsedPacket, rule cnf.Rule, rid int) meta.Match {
	// protocols layer 3
	if rule.Match.ProtoL3N != nil && len(rule.Match.ProtoL3N) > 0 {
		ruleDebug(pkt, rid, fmt.Sprintf("!Proto L3: %v vs %v", rule.Match.ProtoL3N, pkt.L3.Proto))
		return ruleMatch(!anyProtoMatch(rule.Match.ProtoL3N, pkt.L3.Proto))
	}
	if rule.Match.ProtoL3 != nil && len(rule.Match.ProtoL3) > 0 {
		ruleDebug(pkt, rid, fmt.Sprintf("Proto L3: %v vs %v", rule.Match.ProtoL3, pkt.L3.Proto))
		return ruleMatch(anyProtoMatch(rule.Match.ProtoL3, pkt.L3.Proto))
	}
	return meta.MatchNeutral
}

func matchProtoL4(pkt parse.ParsedPacket, rule cnf.Rule, rid int) meta.Match {
	// protocols layer 4
	if rule.Match.ProtoL4 != nil && len(rule.Match.ProtoL4) > 0 {
		ruleDebug(pkt, rid, fmt.Sprintf("Proto L4: %v vs %v", rule.Match.ProtoL4, pkt.L4.Proto))
		return ruleMatch(anyProtoMatch(rule.Match.ProtoL4, pkt.L4.Proto))
	}
	if rule.Match.ProtoL4N != nil && len(rule.Match.ProtoL4N) > 0 {
		ruleDebug(pkt, rid, fmt.Sprintf("!Proto L4: %v vs %v", rule.Match.ProtoL4N, pkt.L4.Proto))
		return ruleMatch(!anyProtoMatch(rule.Match.ProtoL4N, pkt.L4.Proto))
	}
	return meta.MatchNeutral
}

func matchProtoL5(pkt parse.ParsedPacket, rule cnf.Rule, rid int) meta.Match {
	result := meta.MatchNeutral
	// protocols layer 5
	if rule.Match.ProtoL5 != nil && pkt.L5.Proto != meta.ProtoNone {
		ruleDebug(pkt, rid, fmt.Sprintf("Proto L5: %v vs %v", rule.Match.ProtoL5, pkt.L5.Proto))
		if !anyProtoMatch(rule.Match.ProtoL5, pkt.L5.Proto) {
			return meta.MatchNegative
		} else {
			result = meta.MatchPositive
		}
	}
	if rule.Match.ProtoL5N != nil && pkt.L5.Proto != meta.ProtoNone {
		ruleDebug(pkt, rid, fmt.Sprintf("!Proto L5: %v vs %v", rule.Match.ProtoL5N, pkt.L5.Proto))
		if anyProtoMatch(rule.Match.ProtoL5N, pkt.L5.Proto) {
			return meta.MatchNegative
		} else {
			result = meta.MatchPositive
		}
	}
	// todo: make sure 'none' is OK for 'pkt.L5.Encrypted'
	if rule.Match.Encrypted != meta.OptBoolNone && pkt.L5.Encrypted != meta.OptBoolNone {
		ruleDebug(pkt, rid, fmt.Sprintf("Encrypted: %v vs %v", rule.Match.Encrypted, pkt.L5.Encrypted))
		if rule.Match.Encrypted != pkt.L5.Encrypted {
			return meta.MatchNegative
		} else {
			result = meta.MatchPositive
		}
	}
	return result
}

// save result to handle if excluded subnet is inside included subnet
func matchSourceNetwork(pkt parse.ParsedPacket, rule cnf.Rule, rid int) meta.Match {
	result := meta.MatchNeutral
	// source network
	if rule.Match.SrcNet != nil && len(rule.Match.SrcNet) > 0 {
		ruleDebug(pkt, rid, fmt.Sprintf("SNet: %v vs %v", rule.Match.SrcNet, pkt.L3.SrcIP))
		if !anyNetMatch(rule.Match.SrcNet, pkt.L3.SrcIP) {
			result = meta.MatchNegative
		} else {
			result = meta.MatchPositive
		}
	}
	if rule.Match.SrcNetN != nil && len(rule.Match.SrcNetN) > 0 {
		ruleDebug(pkt, rid, fmt.Sprintf("!SNet: %v vs %v", rule.Match.SrcNetN, pkt.L3.SrcIP))
		return ruleMatch(!anyNetMatch(rule.Match.SrcNetN, pkt.L3.SrcIP))
	}
	return result
}

func matchDestinationNetwork(pkt parse.ParsedPacket, rule cnf.Rule, rid int) meta.Match {
	result := meta.MatchNeutral
	// destination network
	if rule.Match.DestNet != nil && len(rule.Match.DestNet) > 0 {
		ruleDebug(pkt, rid, fmt.Sprintf("DNet: %v vs %v", rule.Match.DestNet, pkt.L3.DestIP))
		if !anyNetMatch(rule.Match.DestNet, pkt.L3.DestIP) {
			result = meta.MatchNegative
		} else {
			result = meta.MatchPositive
		}
	}
	if rule.Match.DestNetN != nil && len(rule.Match.DestNetN) > 0 {
		ruleDebug(pkt, rid, fmt.Sprintf("!DNet: %v vs %v", rule.Match.DestNetN, pkt.L3.DestIP))
		return ruleMatch(!anyNetMatch(rule.Match.DestNetN, pkt.L3.DestIP))
	}
	return result
}

func matchSourcePort(pkt parse.ParsedPacket, rule cnf.Rule, rid int) meta.Match {
	// source port
	if rule.Match.SrcPort != nil && len(rule.Match.SrcPort) > 0 {
		ruleDebug(pkt, rid, fmt.Sprintf("SPort: %v vs %v", rule.Match.SrcPort, pkt.L4.SrcPort))
		return ruleMatch(anyPortMatch(rule.Match.SrcPort, pkt.L4.SrcPort))
	}
	if rule.Match.SrcPortN != nil && len(rule.Match.SrcPortN) > 0 {
		ruleDebug(pkt, rid, fmt.Sprintf("!SPort: %v vs %v", rule.Match.SrcPortN, pkt.L4.SrcPort))
		return ruleMatch(!anyPortMatch(rule.Match.SrcPortN, pkt.L4.SrcPort))
	}
	return meta.MatchNeutral
}

func matchDestinationPort(pkt parse.ParsedPacket, rule cnf.Rule, rid int) meta.Match {
	// destination port
	if rule.Match.DestPort != nil && len(rule.Match.DestPort) > 0 {
		ruleDebug(pkt, rid, fmt.Sprintf("DPort: %v vs %v", rule.Match.DestPort, pkt.L4.DestPort))
		return ruleMatch(anyPortMatch(rule.Match.DestPort, pkt.L4.DestPort))
	}
	if rule.Match.DestPortN != nil && len(rule.Match.DestPortN) > 0 {
		ruleDebug(pkt, rid, fmt.Sprintf("!DPort: %v vs %v", rule.Match.DestPortN, pkt.L4.DestPort))
		return ruleMatch(!anyPortMatch(rule.Match.DestPortN, pkt.L4.DestPort))
	}
	return meta.MatchNeutral
}

func anyProtoMatch(list []meta.Proto, single meta.Proto) bool {
	for i := range list {
		if list[i] == single {
			return true
		}
	}
	return false
}

func anyPortMatch(list []uint16, single uint16) bool {
	for i := range list {
		if list[i] == single {
			return true
		}
	}
	return false
}

func anyNetMatch(nets []*net.IPNet, ip net.IP) bool {
	for i := range nets {
		if nets[i].Contains(ip) {
			return true
		}
	}
	return false
}

func ruleMatch(match bool) meta.Match {
	if match {
		return meta.MatchPositive
	} else {
		return meta.MatchNegative
	}
}
