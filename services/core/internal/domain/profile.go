package domain

type EmployeeProfile struct {
	User   *User
	Skills []*Skill
	Plans  []*Plan
}
