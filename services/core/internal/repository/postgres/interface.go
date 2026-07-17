package postgres

import (
	"context"
	"core_service/internal/domain"
	"core_service/internal/repository/model"

	"github.com/google/uuid"
)

//go:generate mockgen -source=interface.go -destination=mocks/mock_repository.go -package=mocks

type AuthRepository interface {
	SaveRefreshToken(ctx context.Context, e domain.RefreshToken) error
	GetRefreshToken(ctx context.Context, jti string) (*domain.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, jti string) error
}

type UserRepository interface {
	Create(ctx context.Context, u *domain.User) (uuid.UUID, error)
	GetById(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, update *domain.UpdateUserProfile) error
	UpdatePassword(ctx context.Context, userID uuid.UUID, hashPassword string) error
	UpdateAvatar(ctx context.Context, userID uuid.UUID, avatarKey *string) error
	ExistsAdmin(ctx context.Context) (bool, error)
	CountAdmins(ctx context.Context) (int, error)
	ListUsers(ctx context.Context, params model.ListUsersParams) ([]*domain.User, error)
	UpdateRole(ctx context.Context, userID uuid.UUID, role domain.Role) error
	AssignManager(ctx context.Context, userID uuid.UUID, managerID uuid.UUID) error
	RemoveManager(ctx context.Context, userID uuid.UUID) error
	ListEmployeesByManager(ctx context.Context, managerID uuid.UUID) ([]*domain.User, error)
	ClearManagerAssignments(ctx context.Context, managerID uuid.UUID) error
	UpdatePosition(ctx context.Context, userID uuid.UUID, position string) error
}

type PlanRepository interface {
	Create(ctx context.Context, plan *domain.Plan) (uuid.UUID, error)
	CreateWithTasks(ctx context.Context, planID uuid.UUID, tasks []*domain.Task) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Plan, error)
	ListByEmployeeID(ctx context.Context, employeeID uuid.UUID) ([]*domain.Plan, error)
	GetEmployeePlan(ctx context.Context, employeeID uuid.UUID, planID uuid.UUID) (*domain.Plan, error)
	RecalculateProgress(ctx context.Context, planID uuid.UUID) (int, error)
	ListByManager(ctx context.Context, managerID uuid.UUID) ([]*domain.Plan, error)
	GetByIDWithTasks(ctx context.Context, id uuid.UUID) (*domain.PlanWithTasks, error)
	Update(ctx context.Context, planID uuid.UUID, title string, description *string) error
	Delete(ctx context.Context, planID uuid.UUID) error
	UpdateGenerationStatus(ctx context.Context, planID uuid.UUID, status domain.GenerationStatus) error
	UpdateAIContent(ctx context.Context, planID uuid.UUID, title string, description *string) error
	Archive(ctx context.Context, planID uuid.UUID) error
	ListAllByEmployeeID(ctx context.Context, employeeID uuid.UUID) ([]*domain.Plan, error)
	ManagerOwnsPlan(ctx context.Context, managerID uuid.UUID, planID uuid.UUID) (bool, error)
}

type TaskRepository interface {
	Create(ctx context.Context, task *domain.Task) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Task, error)
	GetNextPosition(ctx context.Context, planID uuid.UUID) (int, error)
	Update(ctx context.Context, id uuid.UUID, title *string, description *string) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	ListByPlanID(ctx context.Context, planID uuid.UUID) ([]*domain.Task, error)
	CompleteTestingTask(ctx context.Context, planID uuid.UUID) error
}

type SkillRepository interface {
	GetByName(ctx context.Context, name string) (*domain.Skill, error)
	Create(ctx context.Context, skill *domain.Skill) (uuid.UUID, error)
	AttachToUser(ctx context.Context, userID uuid.UUID, skillID uuid.UUID, planID uuid.UUID) error
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Skill, error)
}

type TestRepository interface {
	CreateWithQuestions(ctx context.Context, test *domain.Test, questions []*domain.Question) (uuid.UUID, error)
	GetByPlanID(ctx context.Context, planID uuid.UUID) (*domain.Test, error)
	GetQuestions(ctx context.Context, testID uuid.UUID) ([]*domain.Question, error)
	CreateAttempt(ctx context.Context, attempt *domain.TestAttempt) (uuid.UUID, error)
	CreateAnswers(ctx context.Context, answers []*domain.TestAnswer) error
	GetAttempt(ctx context.Context, testID uuid.UUID, userID uuid.UUID) (*domain.TestAttempt, error)
}
