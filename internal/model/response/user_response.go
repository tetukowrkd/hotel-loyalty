package response

type UserResponse struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Phone      string `json:"phone,omitempty"`
	Role       string `json:"role"`
	MemberCode string `json:"member_code"`
}
