package datatable

import (
	"encoding/json"
	"fmt"
	"regexp"
	"tools/codegen/annotate"
	"tools/codegen/plugins"

	"strings"

	"common/utils"

	"github.com/goraz/humanize"
)

type dataTablePlugin struct {
}

type dataTable struct {
	pkg  humanize.Package
	file humanize.File
	Ann  annotate.Annotate
	typ  humanize.TypeName

	Format []string
}

type context []dataTable

type ColumnDef struct {
	Data       string `json:"data"`
	Name       string `json:"name"`
	Searchable bool   `json:"searchable"`
	Sortable   bool   `json:"sortable"`
	Visible    bool   `json:"visible"`
	Filter     bool   `json:"filter"`
	Title      string `json:"title"`
	Format     bool   `json:"format"`
}

var (
	formater = regexp.MustCompile("Format([a-zA-Z]+)")
	prefix   = regexp.MustCompile("_([a-zA-Z]+)")
)

const (
	filterFunc = `
	
	func ({ .Rec } { .Type }) Filter(u )
	
	`
)

func isExported(s string) bool {
	if len(s) == 0 {
		panic("empty?")
	}

	return strings.ToUpper(s[:1]) == s[:1]
}

// GetType return all types that this plugin can operate on
// for example if the result contain Route then all @Route sections are
// passed to this plugin
func (e dataTablePlugin) GetType() []string {
	return []string{"DataTable"}
}

// Finalize is called after all the functions are done. the context is the one from the
// process
func (e dataTablePlugin) Finalize(c interface{}, p humanize.Package) error {
	var ctx context
	if c != nil {
		var ok bool
		ctx, ok = c.(context)
		if !ok {
			return fmt.Errorf("invalid context, need %T , got %T", ctx, c)
		}
	}

	for i := range ctx {
		res := make(map[string]interface{})
		for key := range ctx[i].Ann.Items {
			if prefix.MatchString(key) {
				res[key[1:]] = ctx[i].Ann.Items[key]
			}
		}
		columns := make([]ColumnDef, 0)
		st := ctx[i].typ.Type.(*humanize.StructType)
		for _, f := range st.Fields {
			if isExported(f.Name) && f.Tags.Get("json") != "-" {
				clm := ColumnDef{}
				tag := f.Tags.Get("json")
				if tag == "" {
					tag = f.Name
				}
				clm.Data = tag
				clm.Name = tag
				clm.Searchable = strings.ToLower(f.Tags.Get("search")) == "true"
				clm.Sortable = strings.ToLower(f.Tags.Get("sort")) == "true"
				clm.Filter = strings.ToLower(f.Tags.Get("filter")) == "true"
				if clm.Filter && clm.Searchable {
					return fmt.Errorf("both filter and search can not set on one field : %s", st.GetDefinition())
				}
				// Every thing is visible except when we note that
				clm.Visible = strings.ToLower(f.Tags.Get("visible")) != "false"
				clm.Title = f.Tags.Get("title")
				if clm.Title == "" {
					clm.Title = f.Name
				}
				clm.Format = utils.StringInArray(f.Name, ctx[i].Format...)
				columns = append(columns, clm)
			}
		}
		res["columns"] = columns

		j, _ := json.MarshalIndent(res, "\t", "\t")
		fmt.Println(string(j))
	}

	j, _ := json.MarshalIndent(ctx, "\t", "\t")
	fmt.Println(string(j))
	return nil
}

func (r *dataTablePlugin) ProcessStructure(
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

	dt := dataTable{
		pkg:  pkg,
		file: p,
		Ann:  a,
		typ:  f,
	}

	for i := range pkg.Files {
		for _, fn := range pkg.Files[i].Functions {
			if fn.Reciever != nil {
				rec := fn.Reciever.Type
				if s, ok := rec.(*humanize.StarType); ok {
					rec = s.Target
				}
				if f.Name == rec.GetDefinition() {
					// found a function
					res := formater.FindStringSubmatch(fn.Name)
					if len(res) == 2 {
						dt.Format = append(dt.Format, res[1])
					}
				}
			}
		}
	}

	ctx = append(ctx, dt)
	return ctx, nil
}

func (r *dataTablePlugin) StructureIsSupported(file humanize.File, fn humanize.TypeName) bool {
	return true
}

func init() {
	plugins.Register(&dataTablePlugin{})
}
