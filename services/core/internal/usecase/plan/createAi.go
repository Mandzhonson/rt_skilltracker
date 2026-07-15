package plan

import (
	"context"
	"errors"
	"strings"

	"core_service/internal/domain"
	"core_service/internal/repository/postgres"
	"core_service/internal/usecase/ai"
	"core_service/internal/usecase/user"

	"github.com/google/uuid"
)

func (s *planService) CreateAI(ctx context.Context, input CreateAIInput) (uuid.UUID, error) {

	if strings.TrimSpace(input.Topic) == "" {
		return uuid.Nil, ErrInvalidTitle
	}

	employee, err := s.userRepo.GetById(ctx, input.EmployeeID)
	if err != nil {
		if errors.Is(err, postgres.ErrUserNotFound) {
			return uuid.Nil, user.ErrUserNotFound
		}
		return uuid.Nil, err
	}

	if !employee.IsEmployee() {
		return uuid.Nil, ErrInvalidEmployee
	}

	manager, err := s.userRepo.GetById(ctx, input.CreatedBy)
	if err != nil {
		if errors.Is(err, postgres.ErrUserNotFound) {
			return uuid.Nil, user.ErrUserNotFound
		}
		return uuid.Nil, err
	}

	if !manager.IsManager() {
		return uuid.Nil, ErrInvalidCreator
	}

	if employee.ManagerID == nil || *employee.ManagerID != manager.ID {
		return uuid.Nil, ErrEmployeeNotAssigned
	}

	userSkills, err := s.skillRepo.ListByUserID(ctx, input.EmployeeID)
	if err != nil {
		return uuid.Nil, err
	}

	skills := make([]string, 0, len(userSkills))
	for _, skill := range userSkills {
		skills = append(skills, skill.Name)
	}

	generated, err := s.aiService.GeneratePlan(
		ctx,
		ai.GeneratePlanInput{
			Topic:          input.Topic,
			Description:    input.Description,
			ExistingSkills: skills,
		},
	)
	if err != nil {
		return uuid.Nil, err
	}

	var planDescription *string
	if generated.Description != "" {
		planDescription = &generated.Description
	}

	planEntity := domain.NewPlan(
		input.EmployeeID,
		input.CreatedBy,
		generated.Title,
		planDescription,
		domain.CreationAI,
	)

	taskEntities := make([]*domain.Task, 0, len(generated.Tasks))

	for i, t := range generated.Tasks {

		var desc *string
		if t.Description != "" {
			desc = &t.Description
		}
		taskEntities = append(taskEntities, domain.NewTask(uuid.Nil, t.Title, desc, i+1))
	}

	planID, err := s.planRepo.CreateWithTasks(ctx, planEntity, taskEntities)
	if err != nil {
		return uuid.Nil, err
	}

	planEntity.ID = planID

	err = s.attachTesting(ctx, planEntity)
	if err != nil {
		return uuid.Nil, err
	}

	return planID, nil
}
