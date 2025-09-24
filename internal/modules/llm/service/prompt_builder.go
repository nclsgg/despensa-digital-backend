package service

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/nclsgg/despensa-digital/backend/internal/modules/llm/model"
)

// PromptBuilderImpl implementa a interface PromptBuilder
type PromptBuilderImpl struct {
	variablePattern *regexp.Regexp
}

// NewPromptBuilder cria uma nova instância do construtor de prompts
func NewPromptBuilder() *PromptBuilderImpl {
	return &PromptBuilderImpl{
		variablePattern: regexp.MustCompile(`\{\{(\w+)\}\}`),
	}
}

// BuildSystemPrompt constrói o prompt do sistema
func (pb *PromptBuilderImpl) BuildSystemPrompt(template string, variables map[string]string) (string, error) {
	return pb.replaceVariables(template, variables)
}

// BuildUserPrompt constrói o prompt do usuário
func (pb *PromptBuilderImpl) BuildUserPrompt(template string, variables map[string]string) (string, error) {
	return pb.replaceVariables(template, variables)
}

// BuildMessages constrói uma lista de mensagens
func (pb *PromptBuilderImpl) BuildMessages(systemPrompt, userPrompt string) []model.Message {
	messages := []model.Message{}

	if systemPrompt != "" {
		messages = append(messages, model.Message{
			Role:    model.RoleSystem,
			Content: systemPrompt,
		})
	}

	if userPrompt != "" {
		messages = append(messages, model.Message{
			Role:    model.RoleUser,
			Content: userPrompt,
		})
	}

	return messages
}

// AddContext adiciona contexto ao prompt
func (pb *PromptBuilderImpl) AddContext(prompt string, context map[string]string) string {
	if len(context) == 0 {
		return prompt
	}

	var contextBuilder strings.Builder
	contextBuilder.WriteString("Contexto adicional:\n")

	for key, value := range context {
		if value != "" {
			contextBuilder.WriteString(fmt.Sprintf("- %s: %s\n", key, value))
		}
	}

	contextBuilder.WriteString("\n")
	contextBuilder.WriteString(prompt)

	return contextBuilder.String()
}

// ValidateTemplate valida um template de prompt
func (pb *PromptBuilderImpl) ValidateTemplate(template string, requiredVariables []string) error {
	// Encontra todas as variáveis no template
	matches := pb.variablePattern.FindAllStringSubmatch(template, -1)
	templateVariables := make(map[string]bool)

	for _, match := range matches {
		if len(match) > 1 {
			templateVariables[match[1]] = true
		}
	}

	// Verifica se todas as variáveis obrigatórias estão presentes
	for _, required := range requiredVariables {
		if !templateVariables[required] {
			return fmt.Errorf("variável obrigatória '%s' não encontrada no template", required)
		}
	}

	return nil
}

// replaceVariables substitui variáveis no template
func (pb *PromptBuilderImpl) replaceVariables(template string, variables map[string]string) (string, error) {
	result := template
	missingVariables := []string{}

	// Primeiro, remove linhas condicionais que não têm valores
	result = pb.removeEmptyConditionals(result, variables)

	// Substitui todas as variáveis encontradas
	result = pb.variablePattern.ReplaceAllStringFunc(result, func(match string) string {
		// Extrai o nome da variável
		varName := strings.Trim(strings.Trim(match, "{"), "}")

		if value, exists := variables[varName]; exists && value != "" {
			return value
		}

		// Para variáveis opcionais conhecidas, retorna string vazia em vez de erro
		optionalVars := map[string]bool{
			"dietary_restrictions": true,
			"purpose":              true,
			"additional_notes":     true,
			"cuisine":              true,
		}

		if optionalVars[varName] {
			return "" // Remove a variável opcional se não tiver valor
		}

		missingVariables = append(missingVariables, varName)
		return match // mantém a variável se não encontrar valor
	})

	if len(missingVariables) > 0 {
		return "", fmt.Errorf("variáveis obrigatórias não fornecidas: %v", missingVariables)
	}

	return result, nil
}

// removeEmptyConditionals remove linhas com condicionais que não têm valores
func (pb *PromptBuilderImpl) removeEmptyConditionals(template string, variables map[string]string) string {
	lines := strings.Split(template, "\n")
	result := []string{}

	for _, line := range lines {
		// Verifica se a linha contém uma condicional opcional
		if strings.Contains(line, "{{#") {
			// Extrai o nome da variável condicional (ex: {{#cuisine}} -> cuisine)
			start := strings.Index(line, "{{#") + 3
			end := strings.Index(line[start:], "}}") + start
			if end > start {
				varName := line[start:end]
				// Se a variável não existe ou está vazia, pula a linha
				if value, exists := variables[varName]; !exists || value == "" {
					continue
				}
				// Remove a sintaxe condicional da linha
				line = strings.Replace(line, "{{#"+varName+"}}", "", 1)
				line = strings.Replace(line, "{{/"+varName+"}}", "", 1)
			}
		}
		result = append(result, line)
	}

	return strings.Join(result, "\n")
}

// RecipePromptTemplates contém templates pré-definidos para receitas
type RecipePromptTemplates struct {
	SystemPrompt string
	UserPrompt   string
}

// GetRecipePromptTemplates retorna os templates para geração de receitas
func GetRecipePromptTemplates() RecipePromptTemplates {
	return RecipePromptTemplates{
		SystemPrompt: `Você é um chef experiente e especialista em culinária brasileira e internacional. Sua missão é criar receitas deliciosas, práticas e personalizadas com base nos ingredientes disponíveis na despensa do usuário.

DIRETRIZES IMPORTANTES:
1. Sempre priorize ingredientes que o usuário JÁ POSSUI na despensa
2. Se precisar de ingredientes adicionais, sugira apenas itens básicos e comuns
3. Adapte a receita ao tempo de preparo solicitado
4. Considere as restrições alimentares informadas
5. Forneça instruções claras e detalhadas
6. Inclua dicas úteis quando apropriado
7. Seja criativo mas prático

REGRAS DE FORMATAÇÃO JSON:
- Use SEMPRE números decimais para quantidades (ex: 0.5 para meio, 1.0 para inteiro)
- NUNCA use frações matemáticas como 1/2, sempre use decimal: 0.5
- NUNCA use expressões matemáticas como 1/2, 1/4, 2/3 - converta para decimal
- Para quantidades "a gosto", use null no campo amount
- Certifique-se de que o JSON seja válido e bem formatado
- Exemplo correto: "amount": 0.5 (para meio)
- Exemplo INCORRETO: "amount": 1/2 (isto causará erro!)

FORMATO DA RESPOSTA:
Responda SEMPRE em JSON válido com a seguinte estrutura:
{
  "title": "Nome da receita",
  "description": "Descrição breve da receita",
  "ingredients": [
    {
      "name": "Nome do ingrediente",
      "amount": quantidade_numerica_decimal,
      "unit": "unidade de medida",
      "available": true/false,
      "alternative": "ingrediente alternativo se não disponível"
    }
  ],
  "instructions": [
    {
      "step": numero_do_passo,
      "description": "Descrição detalhada do passo",
      "time": tempo_em_minutos,
      "temperature": "temperatura se aplicável"
    }
  ],
  "cooking_time": tempo_total_em_minutos,
  "preparation_time": tempo_de_preparo_em_minutos,
  "total_time": tempo_total_em_minutos,
  "serving_size": numero_de_porcoes,
  "difficulty": "easy/medium/hard",
  "meal_type": "tipo da refeição",
  "cuisine": "tipo de culinária",
  "dietary_restrictions": ["restrições aplicáveis"],
  "tips": ["dicas úteis"],
  "nutrition_info": {
    "calories": calorias_aproximadas,
    "protein": proteinas_em_gramas,
    "carbohydrates": carboidratos_em_gramas,
    "fat": gorduras_em_gramas
  }
}`,

		UserPrompt: `Crie uma receita personalizada com as seguintes especificações:

INGREDIENTES DISPONÍVEIS NA DESPENSA:
{{available_ingredients}}

PREFERÊNCIAS:
- Tempo de preparo: {{cooking_time}} minutos
- Tipo de refeição: {{meal_type}}
- Dificuldade: {{difficulty}}
- Número de porções: {{serving_size}}
{{#cuisine}}- Tipo de culinária: {{cuisine}}{{/cuisine}}
{{#dietary_restrictions}}- Restrições alimentares: {{dietary_restrictions}}{{/dietary_restrictions}}
{{#purpose}}- Finalidade: {{purpose}}{{/purpose}}
{{#additional_notes}}- Observações: {{additional_notes}}{{/additional_notes}}

Por favor, crie uma receita que maximize o uso dos ingredientes disponíveis na despensa e atenda às preferências especificadas. Se precisar de ingredientes adicionais, liste apenas os essenciais e comuns.`,
	}
}

// SearchPromptTemplates contém templates para busca de receitas
func GetSearchPromptTemplates() RecipePromptTemplates {
	return RecipePromptTemplates{
		SystemPrompt: `Você é um especialista em culinária que ajuda a encontrar receitas baseadas em ingredientes específicos. Sua função é analisar receitas existentes e determinar quais são mais adequadas para os ingredientes disponíveis.

DIRETRIZES:
1. Analise a compatibilidade entre ingredientes disponíveis e receitas
2. Priorize receitas que usem mais ingredientes disponíveis
3. Considere substituições possíveis
4. Avalie a adequação às restrições e preferências
5. Forneça uma pontuação de compatibilidade (0-100)

FORMATO DA RESPOSTA:
Responda em JSON com uma lista de receitas ranqueadas:
{
  "recipes": [
    {
      "title": "Nome da receita",
      "compatibility_score": pontuação_0_a_100,
      "available_ingredients_count": numero_de_ingredientes_disponíveis,
      "total_ingredients_count": total_de_ingredientes,
      "missing_ingredients": ["ingredientes que faltam"],
      "possible_substitutions": {"ingrediente": "substituto"},
      "source_url": "URL da receita original se disponível"
    }
  ]
}`,

		UserPrompt: `Analise e ranqueie receitas baseadas nos seguintes critérios:

INGREDIENTES DISPONÍVEIS:
{{available_ingredients}}

RECEITAS PARA ANALISAR:
{{recipes_data}}

PREFERÊNCIAS:
{{#cooking_time}}- Tempo máximo: {{cooking_time}} minutos{{/cooking_time}}
{{#meal_type}}- Tipo de refeição: {{meal_type}}{{/meal_type}}
{{#difficulty}}- Dificuldade: {{difficulty}}{{/difficulty}}
{{#dietary_restrictions}}- Restrições: {{dietary_restrictions}}{{/dietary_restrictions}}

Ranqueie as receitas por compatibilidade com os ingredientes disponíveis e preferências especificadas.`,
	}
}
