package plan

import (
	"context"
	"slices"

	"core_service/internal/domain"
	"core_service/internal/usecase/ai"

	"github.com/google/uuid"
)

func (s *planService) GenerateSkillsIfCompleted(ctx context.Context, planID uuid.UUID) error {
	planEntity, err := s.planRepo.GetByID(ctx, planID)
	if err != nil {
		return err
	}

	if planEntity.Progress < 100 {
		return nil
	}

	tasks, err := s.taskRepo.ListByPlanID(ctx, planID)
	if err != nil {
		return err
	}

	taskTitles := make([]string, 0, len(tasks))
	for _, task := range tasks {
		taskTitles = append(taskTitles, task.Title)
	}
	existingSkills, err := s.skillRepo.ListByUserID(ctx, planEntity.EmployeeID)
	if err != nil {
		return err
	}

	existingSkillNames := make([]string, 0, len(existingSkills))

	for _, skill := range existingSkills {
		existingSkillNames = append(existingSkillNames, skill.Name)
	}
	description := ""

	if planEntity.Description != nil {
		description = *planEntity.Description
	}

	generatedSkills, err := s.aiService.ExtractSkills(ctx, ai.ExtractSkillsInput{
		PlanTitle:       planEntity.Title,
		PlanDescription: description,
		Tasks:           taskTitles,
		ExistingSkills:  existingSkillNames,
	},
	)
	if err != nil {
		return err
	}
	for _, generated := range generatedSkills {
		exists := slices.Contains(existingSkillNames, generated.Name)
		if exists {
			continue
		}
		existing, err := s.skillRepo.GetByName(ctx, generated.Name)
		var skillID uuid.UUID
		if err != nil {
			description := generated.Description
			entity := domain.NewSkill(generated.Name, generated.Category, &description)
			skillID, err =
				s.skillRepo.Create(ctx, entity)
			if err != nil {
				return err
			}
		} else {
			skillID = existing.ID
		}
		err = s.skillRepo.AttachToUser(ctx, planEntity.EmployeeID, skillID, planID)
		if err != nil {
			return err
		}
	}
	return nil
}
