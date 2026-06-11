package routes

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
)

// UpdateServiceRegistration adds or updates a service constant and its path mappings in routes.go.
func UpdateServiceRegistration(routesPath, serviceName, serviceConstant string, paths []string) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, routesPath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse routes file: %w", err)
	}

	// 1. Ensure service constant exists in the const block
	updateConstants(f, serviceConstant)

	// 2. Ensure service mappings exist in the serviceMap variable
	updateServiceMap(f, serviceConstant, paths)

	// Write back the modified file
	file, err := os.Create(routesPath)
	if err != nil {
		return fmt.Errorf("failed to open routes file for writing: %w", err)
	}
	defer file.Close()

	if err := format.Node(file, fset, f); err != nil {
		return fmt.Errorf("failed to format and write routes file: %w", err)
	}

	return nil
}

func updateConstants(f *ast.File, serviceConstant string) {
	// Find the const block
	var constDecl *ast.GenDecl
	for _, decl := range f.Decls {
		if d, ok := decl.(*ast.GenDecl); ok && d.Tok == token.CONST {
			constDecl = d
			break
		}
	}

	if constDecl == nil {
		return // Should not happen in a valid routes.go
	}

	// Check if constant already exists
	exists := false
	maxPort := 8000
	for _, spec := range constDecl.Specs {
		vspec := spec.(*ast.ValueSpec)
		for i, name := range vspec.Names {
			if name.Name == serviceConstant {
				exists = true
				break
			}
			// Track max port for next assignment
			if len(vspec.Values) > i {
				if lit, ok := vspec.Values[i].(*ast.BasicLit); ok && lit.Kind == token.STRING {
					val := strings.Trim(lit.Value, "\"")
					idx := strings.LastIndex(val, ":")
					if idx != -1 {
						var port int
						if _, err := fmt.Sscanf(val[idx:], ":%d", &port); err == nil {
							if port > maxPort {
								maxPort = port
							}
						}
					}
				}
			}
		}
	}

	if !exists {
		// Add new constant with next port
		newPort := maxPort + 1
		newSpec := &ast.ValueSpec{
			Names: []*ast.Ident{ast.NewIdent(serviceConstant)},
			Values: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf("\"localhost:%d\"", newPort),
				},
			},
		}
		constDecl.Specs = append(constDecl.Specs, newSpec)
	}
}

func updateServiceMap(f *ast.File, serviceConstant string, paths []string) {
	// Find serviceMap variable declaration
	var compositeLit *ast.CompositeLit
	for _, decl := range f.Decls {
		if d, ok := decl.(*ast.GenDecl); ok && d.Tok == token.VAR {
			for _, spec := range d.Specs {
				vspec := spec.(*ast.ValueSpec)
				for _, name := range vspec.Names {
					if name.Name == "serviceMap" {
						if cl, ok := vspec.Values[0].(*ast.CompositeLit); ok {
							compositeLit = cl
							break
						}
					}
				}
			}
		}
	}

	if compositeLit == nil {
		return
	}

	// Track existing paths
	existingPaths := make(map[string]bool)
	for _, elt := range compositeLit.Elts {
		if kv, ok := elt.(*ast.KeyValueExpr); ok {
			if lit, ok := kv.Key.(*ast.BasicLit); ok && lit.Kind == token.STRING {
				path := strings.Trim(lit.Value, "\"")
				existingPaths[path] = true
			}
		}
	}

	// Add missing paths
	for _, path := range paths {
		if !existingPaths[path] {
			newKV := &ast.KeyValueExpr{
				Key: &ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf("%q", path),
				},
				Value: ast.NewIdent(serviceConstant),
			}
			compositeLit.Elts = append(compositeLit.Elts, newKV)
		}
	}

	// Optional: Sort entries by path for cleanliness? 
	// The current routes.go is grouped by service, which is better.
	// We'll just leave them appended for now.
}

// EnsurePackageImport adds an import if it's missing.
func EnsurePackageImport(f *ast.File, path string) {
	if !astutil.UsesImport(f, path) {
		astutil.AddImport(nil, f, path)
	}
}
