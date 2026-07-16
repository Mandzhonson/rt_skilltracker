package dto

type EmployeeProfileResponse struct {
	User   UserResponse    `json:"user"`
	Skills []SkillResponse `json:"skills"`
	Plans  []PlanResponse  `json:"plans"`
}
