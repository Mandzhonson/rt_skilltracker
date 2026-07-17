package ai

import (
	"fmt"
	"strings"
)

func buildSkillPrompt(input ExtractSkillsInput) string {

	existingSkills := "Нет полученных навыков"

	if len(input.ExistingSkills) > 0 {
		existingSkills = strings.Join(
			input.ExistingSkills,
			"\n",
		)
	}
	fmt.Println(input)
	return fmt.Sprintf(`
Ты Senior Software Architect.

Твоя задача — определить навыки, которые сотрудник должен получить
после завершения плана обучения.

ВАЖНО:

1. Не выделяй слишком общие навыки:
- "обучение"
- "разработка"
- "программирование"

2. Выделяй только конкретные:
- языки программирования
- технологии
- базы данных
- инструменты
- архитектурные подходы
- профессиональные практические навыки


3. Нормализуй названия:

Примеры:

Go → Go

Golang → Go

Postgres → PostgreSQL

PostgreSQL → PostgreSQL

Docker Compose → Docker


4. У сотрудника уже есть следующие навыки:

%s


НЕ возвращай навыки из этого списка.

Если все необходимые навыки уже есть у сотрудника,
верни пустой массив:

[]


Название плана:

%s


Описание:

%s


Задачи:

%s


Верни только JSON.

Формат:

[
  {
    "name":"",
    "category":"",
    "description":""
  }
]

`,
		existingSkills,
		input.PlanTitle,
		input.PlanDescription,
		strings.Join(input.Tasks, "\n"),
	)
}
