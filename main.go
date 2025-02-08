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

type GenCmd struct {
	Name   string `help:"Name of the entity" required:""`
	Fields string `help:"Fields in format name:type:required,... (e.g., name:string:true,age:number:false)" required:""`
	Out    string `help:"Output directory" default:"."`
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

	// create output directory
	outDir := filepath.Join(g.Out, strings.ToLower(entity.Name))
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	log.Info().
		Str("entity", entity.Name).
		Int("fields", len(entity.Fields)).
		Str("dir", outDir).
		Msg("generating entity")

	// TODO: add template generation logic here
	// generate types file
	tmplPath := "templates/types.ts.tmpl"
	outFile := filepath.Join(outDir, entity.Name+".types.ts")
	if err := generateFile(entity, tmplPath, outFile); err != nil {
		return fmt.Errorf("failed to generate file: %v", err)
	}

	// generate constants file
	// Add after your types file generation
	tmplPath = "templates/constants.ts.tmpl"
	outFile = filepath.Join(outDir, entity.Name+".constants.ts")
	if err := generateFile(entity, tmplPath, outFile); err != nil {
		return fmt.Errorf("failed to generate file: %v", err)
	}

	// Api functions and react-query hooks
	tmplPath = "templates/api.ts.tmpl"
	outFile = filepath.Join(outDir, entity.Name+".api.ts")
	if err := generateFile(entity, tmplPath, outFile); err != nil {
		return fmt.Errorf("failed to generate file: %v", err)
	}

	// Form component
	tmplPath = "templates/form.tsx.tmpl"
	outFile = filepath.Join(outDir, entity.Name+".form.tsx")
	if err := generateFile(entity, tmplPath, outFile); err != nil {
		return fmt.Errorf("failed to generate file: %v", err)
	}

	// Table component
	tmplPath = "templates/table.tsx.tmpl"
	outFile = filepath.Join(outDir, entity.Name+".table.tsx")
	if err := generateFile(entity, tmplPath, outFile); err != nil {
		return fmt.Errorf("failed to generate file: %v", err)
	}

	// List component
	tmplPath = "templates/list.tsx.tmpl"
	outFile = filepath.Join(outDir, entity.Name+".list.tsx")
	if err := generateFile(entity, tmplPath, outFile); err != nil {
		return fmt.Errorf("failed to generate file: %v", err)
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
