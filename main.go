package main

import (
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"
)

var CLI struct {
	Gen GenCmd `cmd:"" help:"Generate a new React CRUD entity"`
}

type ComponentType string

const (
	TypesComponent     ComponentType = "types"
	ConstantsComponent ComponentType = "constants"
	ApiComponent       ComponentType = "api"
	FormComponent      ComponentType = "form"
	TableComponent     ComponentType = "table"
	ListComponent      ComponentType = "list"
	DetailComponent    ComponentType = "detail"
	AllComponents      ComponentType = "all"
)

type GenCmd struct {
	Name      string        `help:"Name of the entity" required:""`
	Fields    string        `help:"Fields in format name:type:required,... (e.g., name:string:true,age:number:false)" required:""`
	Component ComponentType `help:"Component type to generate" enum:"types,constants,api,form,table,list,detail,all" default:"all"`
	Out       bool          `help:"Create components in a subdirectory named after the entity" default:"false"`
}

type Field struct {
	Name     string
	Type     string
	Required bool
}

type Entity struct {
	Name   string
	Fields []Field
}

func (g *GenCmd) Run() error {
	// parse fields
	var fields []Field
	for _, field := range strings.Split(g.Fields, ",") {
		fieldParts := strings.Split(field, ":")
		if len(fieldParts) != 3 {
			return fmt.Errorf("invalid field format: %s", field)
		}
		fields = append(fields, Field{
			Name:     fieldParts[0],
			Type:     fieldParts[1],
			Required: fieldParts[2] == "true",
		})
	}

	entity := Entity{
		Name:   g.Name,
		Fields: fields,
	}

	// Determine output directory
	var outDir string
	if g.Out {
		outDir = strings.ToLower(entity.Name)
	} else {
		outDir = "."
	}

	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	log.Info().
		Str("entity", entity.Name).
		Int("fields", len(entity.Fields)).
		Str("dir", outDir).
		Str("component", string(g.Component)).
		Bool("subfolder", g.Out).
		Msg("generating entity")

	components := map[ComponentType]struct {
		template string
		output   string
	}{
		TypesComponent:     {"templates/types.ts.tmpl", entity.Name + ".types.ts"},
		ConstantsComponent: {"templates/constants.ts.tmpl", entity.Name + ".constants.ts"},
		ApiComponent:       {"templates/api.ts.tmpl", entity.Name + ".api.ts"},
		FormComponent:      {"templates/form.tsx.tmpl", entity.Name + ".form.tsx"},
		TableComponent:     {"templates/table.tsx.tmpl", entity.Name + ".table.tsx"},
		ListComponent:      {"templates/list.tsx.tmpl", entity.Name + ".list.tsx"},
		DetailComponent:    {"templates/detail.tsx.tmpl", entity.Name + ".detail.tsx"},
	}

	// Generate specified component or all components
	if g.Component == AllComponents {
		// Generate all components
		for _, comp := range components {
			if err := generateFile(entity, comp.template, filepath.Join(outDir, comp.output)); err != nil {
				return fmt.Errorf("failed to generate file: %v", err)
			}
		}
	} else {
		// Generate specific component
		if comp, ok := components[g.Component]; ok {
			if err := generateFile(entity, comp.template, filepath.Join(outDir, comp.output)); err != nil {
				return fmt.Errorf("failed to generate file: %v", err)
			}
		} else {
			return fmt.Errorf("invalid component type: %s", g.Component)
		}
	}

	return nil
}

func generateFile(entity Entity, tmplPath string, outFile string) error {
	// custom template functions
	funcMap := template.FuncMap{
		"lower": strings.ToLower,
		"upper": strings.ToUpper,
		"title": strings.Title,
		"eq":    reflect.DeepEqual, // Add this for the type comparison
		// Add function to check if field is required
		"required": func(field Field) bool {
			return field.Required
		},
	}

	// Read template content
	content, err := os.ReadFile(tmplPath)
	if err != nil {
		return fmt.Errorf("failed to read template: %w", err)
	}

	// read and parse the template with custom functions
	tmpl, err := template.New(filepath.Base(tmplPath)).Funcs(funcMap).Parse(string(content))
	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}

	// create output file
	file, err := os.Create(outFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Error().Err(err).Msg("failed to close file")
		}
	}(file)

	// execute template
	if err := tmpl.Execute(file, entity); err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}

	log.Info().Str("file", outFile).Msg("generated file")
	return nil
}

func main() {
	// initialize zerolog with pretty console output for development
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).
		With().
		Timestamp().
		Logger()

	ctx := kong.Parse(&CLI,
		kong.Name("gen"),
		kong.Description("Generate a new React CRUD entity"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}))

	err := ctx.Run()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to run command")
	}
}
