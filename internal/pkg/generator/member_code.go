package generator

import (
	"hotel-loyalty/internal/infrastructure/logger"
	"strings"

	"github.com/google/uuid"
)

func GenerateMemberCode(role string) string {
	prefix := getPrefix(role)

	code := strings.ToUpper(uuid.New().String()[:6])

	memberCode := prefix + "-" + code

	logger.InfoLogger.Println(
		"[GENERATOR][MemberCode] generated",
		"role=", role,
		"member_code=", memberCode,
	)

	return memberCode
}

func getPrefix(role string) string {
	switch role {
	case "manager":
		return "MNG"
	case "staff":
		return "STF"
	case "admin":
		return "ADM"
	default:
		logger.ErrorLogger.Println(
			"[GENERATOR][MemberCode] unknown role, fallback to MEM",
			"role=", role,
		)
		return "MEM"
	}
}
