package response

type RegisterUserResponse struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	MemberCode string `json:"member_code"`
}
