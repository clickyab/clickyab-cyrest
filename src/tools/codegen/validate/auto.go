package validate

import (
	"fmt"
	"html/template"
	"tools/codegen/annotate"
	"tools/codegen/plugins"

	"bytes"
	"io/ioutil"

	"path/filepath"

	"github.com/goraz/humanize"
	"golang.org/x/tools/imports"
)

type validatePlugin struct {
}

type validate struct {
	pkg  humanize.Package
	file humanize.File
	ann  annotate.Annotate
	typ  humanize.TypeName

	Rec  string
	Type string
}

type context []validate

var (
	validateFunc = `
package {{ .PackageName }}
// AUTO GENERATED CODE. DO NOT EDIT!

import (
	"gopkg.in/labstack/echo.v3"
	"gopkg.in/go-playground/validator.v9"
)

	{{ range $m := .Data }}
	func ({{ $m.Rec }} {{ $m.Type }}) Validate(ctx echo.Context ) error {
		return validator.New().Struct({{ $m.Rec }})
	}
	{{ end }}
	`

	tpl = template.Must(template.New("validate").Parse(validateFunc))
)

// GetType return all types that this plugin can operate on
// for example if the result contain Route then all @Route sections are
// passed to this plugin
func (e validatePlugin) GetType() []string {
	return []string{"Validate"}
}

// Finalize is called after all the functions are done. the context is the one from the
// process
func (e validatePlugin) Finalize(c interface{}, p humanize.Package) error {
	var ctx context
	if c != nil {
		var ok bool
		ctx, ok = c.(context)
		if !ok {
			return fmt.Errorf("invalid context, need %T , got %T", ctx, c)
		}
	}

	buf := &bytes.Buffer{}
	err := tpl.Execute(buf, struct {
		Data        context
		PackageName string
	}{
		Data:        ctx,
		PackageName: p.Name,
	})
	if err != nil {
		return err
	}
	f := filepath.Dir(p.Files[0].FileName)
	f = filepath.Join(f, "validators.gen.go")
	res, err := imports.Process("", buf.Bytes(), nil)
	if err != nil {
		fmt.Println(buf.String())
		return err
	}

	err = ioutil.WriteFile(f, res, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (r *validatePlugin) ProcessStructure(
	c interface{},
	pkg humanize.Package,
	p humanize.File,
	f humanize.TypeName,
	a annotate.Annotate,
) (interface{}, error) {
	var ctx context
	if c != nil {
		var ok bool
		ctx, ok = c.(context)
		if !ok {
			return nil, fmt.Errorf("invalid context, need %T , got %T", ctx, c)
		}
	}

	dt := validate{
		pkg:  pkg,
		file: p,
		ann:  a,
		typ:  f,

		Type: f.Name,
		Rec:  "pl",
	}

bigLoop:
	for i := range pkg.Files {
		for _, fn := range pkg.Files[i].Functions {
			if fn.Reciever != nil {
				rec := fn.Reciever.Type
				if s, ok := rec.(*humanize.StarType); ok {
					rec = s.Target
				}
				if f.Name == rec.GetDefinition() {
					dt.Rec = fn.Reciever.Name
					break bigLoop
				}
			}
		}
	}

	ctx = append(ctx, dt)
	return ctx, nil
}

func (r *validatePlugin) StructureIsSupported(file humanize.File, fn humanize.TypeName) bool {
	return true
}
func init() {
	plugins.Register(&validatePlugin{})
}
