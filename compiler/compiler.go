package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/banshanhanfu/ws-lang/parser"
)

// InstructionSet is the canonical YAML instruction set format
type InstructionSet struct {
	TaskID  string  `yaml:"task_id,omitempty" json:"task_id,omitempty"`
	Name    string  `yaml:"name" json:"name"`
	Status  string  `yaml:"status,omitempty" json:"status,omitempty"`
	Steps   []*Step `yaml:"steps" json:"steps"`
}

// Step represents a single execution step
type Step struct {
	ID        string     `yaml:"id" json:"id"`
	Name      string     `yaml:"name" json:"name"`
	Cap       string     `yaml:"cap" json:"cap"`
	Input     *Input     `yaml:"input,omitempty" json:"input,omitempty"`
	Args      map[string]string `yaml:"args,omitempty" json:"args,omitempty"`
	Outputs   []*Output  `yaml:"outputs,omitempty" json:"outputs,omitempty"`
	DependsOn []string   `yaml:"depends_on,omitempty" json:"depends_on,omitempty"`
	OnError   *ErrorHandling `yaml:"on_error,omitempty" json:"on_error,omitempty"`
}

// Input describes step input source
type Input struct {
	Type    string   `yaml:"type" json:"type"`
	Value   string   `yaml:"value,omitempty" json:"value,omitempty"`
	Path    string   `yaml:"path,omitempty" json:"path,omitempty"`
	From    string   `yaml:"from,omitempty" json:"from,omitempty"`
	Sources []*Input `yaml:"sources,omitempty" json:"sources,omitempty"`
}

// Output describes a named output port
type Output struct {
	Name string `yaml:"name" json:"name"`
	From string `yaml:"from" json:"from"`
}

// ErrorHandling describes error handling policy
type ErrorHandling struct {
	Action     string `yaml:"action" json:"action"`
	MaxRetries int    `yaml:"max_retries,omitempty" json:"max_retries,omitempty"`
	Interval   int    `yaml:"interval,omitempty" json:"interval,omitempty"`
}

// Translation tables
var translations = map[string]map[string]string{
	"zh": {
		"任务": "task", "步骤": "step", "能力": "cap",
		"来自": "from", "合并": "merge", "参数": "args",
		"输入": "input", "输出": "output", "出错时": "on_error",
		"重试": "retry", "忽略": "ignore", "停止": "stop",
		"名称": "name", "描述": "description", "状态": "status",
		"结果": "result", "错误": "error", "依赖": "depends_on",
		"来源": "sources",
	},
	"ja": {
		"タスク": "task", "ステップ": "step", "能力": "cap",
		"呼出": "cap", "引数": "args", "入力": "input",
		"出力": "output", "エラー時": "on_error", "リトライ": "retry",
		"無視": "ignore", "停止": "stop",
	},
}

func main() {
	lang := flag.String("lang", "en", "Input language (en/zh/ja)")
	outputJSON := flag.Bool("json", false, "Output in JSON format instead of YAML")
	outputFile := flag.String("o", "", "Output file path (default: stdout)")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: compiler [-lang en|zh|ja] [-json] [-o output.yaml] <input.ws>\n")
		os.Exit(1)
	}

	inputPath := args[0]

	// Read input file
	data, err := os.ReadFile(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input file: %v\n", err)
		os.Exit(1)
	}

	content := string(data)

	// If language is not English, apply translation
	if *lang != "en" {
		content = translateContent(content, *lang)
	}

	// Parse KV format
	result, err := parser.Parse(content)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Parse error: %v\n", err)
		os.Exit(1)
	}

	// Compile to instruction set
	instructionSet := compile(result, inputPath)

	// Output
	var output []byte
	if *outputJSON {
		output, err = json.MarshalIndent(instructionSet, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "JSON marshal error: %v\n", err)
			os.Exit(1)
		}
	} else {
		output = []byte(toYAML(instructionSet))
	}

	if *outputFile != "" {
		if err := os.WriteFile(*outputFile, output, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Compiled to %s\n", *outputFile)
	} else {
		fmt.Println(string(output))
	}
}

func translateContent(content string, lang string) string {
	table, ok := translations[lang]
	if !ok {
		return content
	}

	result := content
	// Sort keys by length (longest first) to avoid partial replacements
	type kv struct{ k, v string }
	var pairs []kv
	for k, v := range table {
		pairs = append(pairs, kv{k, v})
	}
	// Sort by key length descending
	for i := 0; i < len(pairs); i++ {
		for j := i + 1; j < len(pairs); j++ {
			if len(pairs[i].k) < len(pairs[j].k) {
				pairs[i], pairs[j] = pairs[j], pairs[i]
			}
		}
	}

	for _, p := range pairs {
		result = strings.ReplaceAll(result, p.k, p.v)
	}
	return result
}

func compile(result *parser.ParseResult, inputPath string) *InstructionSet {
	baseName := filepath.Base(inputPath)
	baseName = strings.TrimSuffix(baseName, filepath.Ext(baseName))

	is := &InstructionSet{
		Name:   result.TaskName,
		Status: "pending",
	}

	stepIndex := 0

	// Build a name→id mapping for dependency resolution
	stepNameToID := make(map[string]string)

	// First pass: register all step names
	for _, node := range result.Nodes {
		stepIndex = registerStepNames(node, &stepIndex, stepNameToID)
	}

	// Second pass: compile steps with dependency resolution
	stepIndex = 0
	for _, node := range result.Nodes {
		stepIndex = compileNode(node, is, &stepIndex, stepNameToID)
	}

	if is.Name == "" {
		is.Name = baseName
	}

	return is
}

func registerStepNames(node *parser.Node, stepIndex *int, nameToID map[string]string) int {
	// Top-level step nodes
	if node.Key == "step" || node.Key == "steps" || node.Key == "步骤" || node.Key == "ステップ" {
		*stepIndex++
		id := fmt.Sprintf("s-%d", *stepIndex)
		if node.Value != "" {
			nameToID[node.Value] = id
		}
		return *stepIndex
	}

	// Task wrapper children
	for _, child := range node.Children {
		if child.Key == "step" || child.Key == "steps" || child.Key == "步骤" || child.Key == "ステップ" {
			*stepIndex++
			id := fmt.Sprintf("s-%d", *stepIndex)
			if child.Value != "" {
				nameToID[child.Value] = id
			}
		}
	}
	return *stepIndex
}

func compileNode(node *parser.Node, is *InstructionSet, stepIndex *int, nameToID map[string]string) int {
	// Top-level step nodes
	if node.Key == "step" || node.Key == "steps" || node.Key == "步骤" || node.Key == "ステップ" {
		*stepIndex++
		step := buildStep(node, *stepIndex, nameToID)
		is.Steps = append(is.Steps, step)
		return *stepIndex
	}

	// Task wrapper: process children
	for _, child := range node.Children {
		if child.Key == "step" || child.Key == "steps" || child.Key == "步骤" || child.Key == "ステップ" {
			*stepIndex++
			step := buildStep(child, *stepIndex, nameToID)
			is.Steps = append(is.Steps, step)
		} else {
			// Nested structure, recurse
			*stepIndex = compileNode(child, is, stepIndex, nameToID)
		}
	}
	return *stepIndex
}

func buildStep(node *parser.Node, index int, nameToID map[string]string) *Step {
	step := &Step{
		ID:   fmt.Sprintf("s-%d", index),
		Name: node.Value,
	}

	for _, child := range node.Children {
		switch child.Key {
		case "cap", "capability", "能力":
			step.Cap = child.Value
		case "input", "输入":
			step.Input = compileInput(child, nameToID)
		case "on_error", "出错时", "エラー時":
			step.OnError = compileErrorHandling(child.Value)
		case "output", "输出", "outputs":
			// Handled by nodeOutputs, skip args
		default:
			// Skip children that are output port definitions
			if child.IsOutput {
				continue
			}
			if step.Args == nil {
				step.Args = make(map[string]string)
			}
			step.Args[child.Key] = child.Value
		}
	}

	// Set default output
	if step.Cap != "" && len(nodeOutputs(node)) == 0 {
		step.Outputs = []*Output{
			{Name: "default", From: "stdout.text"},
		}
	} else {
		step.Outputs = nodeOutputs(node)
	}

	return step
}

func nodeOutputs(node *parser.Node) []*Output {
	var outputs []*Output
	for _, child := range node.Children {
		if child.IsOutput || child.Key == "output" || child.Key == "输出" {
			outputs = append(outputs, &Output{
				Name: child.Value,
				From: "stdout.text",
			})
		}
	}
	return outputs
}

func compileInput(node *parser.Node, nameToID map[string]string) *Input {
	if node.IsArray {
		// Merge input: [$ref1, $ref2]
		input := &Input{Type: "merge"}
		inner := strings.TrimSpace(node.Value[1 : len(node.Value)-1])
		refs := strings.Split(inner, ",")
		for _, ref := range refs {
			ref = strings.TrimSpace(ref)
			refName := strings.TrimPrefix(ref, "$")
			fromID := resolveStepRef(refName, nameToID)
			input.Sources = append(input.Sources, &Input{
				Type: "ref",
				From: fromID,
			})
		}
		return input
	}

	if node.IsRef {
		refName := strings.TrimPrefix(node.Value, "$")
		fromID := resolveStepRef(refName, nameToID)
		return &Input{Type: "ref", From: fromID}
	}

	// Check if value looks like a file path
	val := node.Value
	if strings.HasPrefix(val, "/") || strings.HasPrefix(val, "./") || strings.HasPrefix(val, "../") ||
		strings.HasPrefix(val, "s3://") || strings.HasPrefix(val, "http://") || strings.HasPrefix(val, "https://") {
		return &Input{Type: "file", Path: val}
	}

	return &Input{Type: "literal", Value: val}
}

func resolveStepRef(refName string, nameToID map[string]string) string {
	// Check if it's a known step name
	if id, ok := nameToID[refName]; ok {
		return id
	}
	// Also check without underscores/dots variations
	for name, id := range nameToID {
		if strings.ReplaceAll(name, " ", "") == strings.ReplaceAll(refName, " ", "") {
			return id
		}
	}
	// Return original name if not found (it might be a step ID directly)
	return refName
}

func compileErrorHandling(value string) *ErrorHandling {
	eh := &ErrorHandling{Action: "stop"}

	// Parse patterns like: retry(3, 5s), ignore, stop
	if strings.HasPrefix(value, "retry") || strings.HasPrefix(value, "重试") {
		eh.Action = "retry"
		eh.MaxRetries = 3
		eh.Interval = 5

		// Extract (n, ms)
		if idx := strings.Index(value, "("); idx > 0 {
			end := strings.Index(value, ")")
			if end < 0 {
				end = len(value)
			}
			params := value[idx+1 : end]
			parts := strings.Split(params, ",")
			if len(parts) > 0 {
				fmt.Sscanf(strings.TrimSpace(parts[0]), "%d", &eh.MaxRetries)
			}
			if len(parts) > 1 {
				var interval int
				intervalStr := strings.TrimSpace(parts[1])
				intervalStr = strings.TrimSuffix(intervalStr, "s")
				fmt.Sscanf(intervalStr, "%d", &interval)
				if interval > 0 {
					eh.Interval = interval
				}
			}
		}
	} else if value == "ignore" || value == "忽略" {
		eh.Action = "ignore"
	} else if value == "stop" || value == "停止" {
		eh.Action = "stop"
	}

	return eh
}

func toYAML(is *InstructionSet) string {
	var b strings.Builder

	// Header
	b.WriteString(fmt.Sprintf("# ws-lang instruction set\n"))
	b.WriteString(fmt.Sprintf("# generated by ws-lang compiler\n"))
	b.WriteString(fmt.Sprintf("---\n"))
	b.WriteString(fmt.Sprintf("name: %s\n", is.Name))
	b.WriteString(fmt.Sprintf("status: %s\n", is.Status))
	if is.TaskID != "" {
		b.WriteString(fmt.Sprintf("task_id: %s\n", is.TaskID))
	}
	b.WriteString(fmt.Sprintf("steps:\n"))

	for _, step := range is.Steps {
		b.WriteString(fmt.Sprintf("  - id: %s\n", step.ID))
		b.WriteString(fmt.Sprintf("    name: %s\n", step.Name))
		b.WriteString(fmt.Sprintf("    cap: %s\n", step.Cap))

		if step.Input != nil {
			writeInputYAML(&b, step.Input, 4)
		}

		if len(step.Args) > 0 {
			b.WriteString(fmt.Sprintf("    args:\n"))
			for k, v := range step.Args {
				b.WriteString(fmt.Sprintf("      %s: %s\n", k, v))
			}
		}

		if len(step.Outputs) > 0 {
			b.WriteString(fmt.Sprintf("    outputs:\n"))
			for _, o := range step.Outputs {
				b.WriteString(fmt.Sprintf("      - name: %s\n", o.Name))
				b.WriteString(fmt.Sprintf("        from: %s\n", o.From))
			}
		}

		if step.OnError != nil {
			b.WriteString(fmt.Sprintf("    on_error:\n"))
			b.WriteString(fmt.Sprintf("      action: %s\n", step.OnError.Action))
			if step.OnError.MaxRetries > 0 {
				b.WriteString(fmt.Sprintf("      max_retries: %d\n", step.OnError.MaxRetries))
			}
			if step.OnError.Interval > 0 {
				b.WriteString(fmt.Sprintf("      interval: %d\n", step.OnError.Interval))
			}
		}

		if len(step.DependsOn) > 0 {
			b.WriteString(fmt.Sprintf("    depends_on:\n"))
			for _, d := range step.DependsOn {
				b.WriteString(fmt.Sprintf("      - %s\n", d))
			}
		}

		b.WriteString("\n")
	}

	return b.String()
}

func writeInputYAML(b *strings.Builder, input *Input, indent int) {
	prefix := strings.Repeat(" ", indent)

	switch input.Type {
	case "literal":
		b.WriteString(fmt.Sprintf("%sinput:\n", prefix))
		b.WriteString(fmt.Sprintf("%s  type: literal\n", prefix))
		b.WriteString(fmt.Sprintf("%s  value: %s\n", prefix, input.Value))
	case "file":
		b.WriteString(fmt.Sprintf("%sinput:\n", prefix))
		b.WriteString(fmt.Sprintf("%s  type: file\n", prefix))
		b.WriteString(fmt.Sprintf("%s  path: %s\n", prefix, input.Path))
	case "ref":
		b.WriteString(fmt.Sprintf("%sinput:\n", prefix))
		b.WriteString(fmt.Sprintf("%s  type: ref\n", prefix))
		b.WriteString(fmt.Sprintf("%s  from: %s\n", prefix, input.From))
	case "merge":
		b.WriteString(fmt.Sprintf("%sinput:\n", prefix))
		b.WriteString(fmt.Sprintf("%s  type: merge\n", prefix))
		b.WriteString(fmt.Sprintf("%s  sources:\n", prefix))
		for _, src := range input.Sources {
			b.WriteString(fmt.Sprintf("%s    - type: ref\n", prefix))
			b.WriteString(fmt.Sprintf("%s      from: %s\n", prefix, src.From))
		}
	}
}
