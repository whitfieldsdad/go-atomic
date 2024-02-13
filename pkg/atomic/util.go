package atomic

import (
	"regexp"
)

var (
	AttackTechniqueIdRegex = regexp.MustCompile("[Tt]\\d{4}(?:\\.\\d{3})?")
)

// isAttackTechniqueId returns true if the provided value looks like an ATT&CK technique ID (e.g. T1003)
func isAttackTechniqueId(value string) bool {
	return AttackTechniqueIdRegex.MatchString(value)
}
