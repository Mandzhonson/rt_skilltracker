package ai

import (
	"strings"
)

var skillAliases = map[string]string{
	"go":                        "Go",
	"golang":                    "Go",
	"postgres":                  "PostgreSQL",
	"postgresql":                "PostgreSQL",
	"postgre sql":               "PostgreSQL",
	"sql":                       "SQL",
	"js":                        "JavaScript",
	"javascript":                "JavaScript",
	"ts":                        "TypeScript",
	"typescript":                "TypeScript",
	"react":                     "React",
	"reactjs":                   "React",
	"eda":                       "EDA",
	"exploratory data analysis": "EDA",
	"numpy":                     "NumPy",
	"pandas":                    "Pandas",
	"sklearn":                   "Scikit-learn",
	"scikit learn":              "Scikit-learn",
	"scikit-learn":              "Scikit-learn",
	"model evaluation":          "Оценка моделей",
	"cross validation":          "Кросс-валидация",
	"cross-validation":          "Кросс-валидация",
	"кросс валидация":           "Кросс-валидация",
	"roc auc":                   "ROC-AUC",
}

var ignoredSkills = map[string]struct{}{
	"обучение":                        {},
	"разработка":                      {},
	"программирование":                {},
	"реализация проекта":              {},
	"реализация проекта data science": {},
	"создание проекта":                {},
	"работа с данными":                {},
}

func normalizeSkillName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, "-", " ")
	name = strings.Join(strings.Fields(name), " ")

	if _, ok := ignoredSkills[name]; ok {
		return ""
	}

	if alias, ok := skillAliases[name]; ok {
		return alias
	}

	return strings.TrimSpace(name)
}

func normalizeCategory(category string) string {

	category = strings.TrimSpace(category)

	if category == "" {
		return "Другое"
	}

	switch strings.ToLower(category) {

	case "database",
		"базы данных",
		"работа с базами данных":

		return "Базы данных"

	case "machine learning",
		"машинное обучение":

		return "Машинное обучение"

	case "programming languages",
		"языки программирования":

		return "Языки программирования"

	case "data analysis",
		"анализ данных":

		return "Анализ данных"

	default:
		return category
	}
}

func normalizeSkills(skills []SkillCandidate) []SkillCandidate {
	result := make(map[string]SkillCandidate)
	for _, skill := range skills {
		name := normalizeSkillName(skill.Name)
		if name == "" {
			continue
		}

		skill.Name = name
		skill.Category = normalizeCategory(skill.Category)
		skill.Description = strings.TrimSpace(skill.Description)

		key := strings.ToLower(skill.Name)
		if _, exists := result[key]; exists {
			continue
		}

		result[key] = skill
	}
	output := make([]SkillCandidate, 0, len(result))

	for _, skill := range result {
		output = append(output, skill)
	}

	return output
}
