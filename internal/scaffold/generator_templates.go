package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const mainTemplate = `package main

import (
	"log"
	"net/http"

	"github.com/Tylores/egot/internal/{{.ServiceName}}/handler"
	"github.com/Tylores/egot/internal/{{.ServiceName}}/repository/memory"
	"github.com/Tylores/egot/internal/routes"
	"github.com/Tylores/egot/internal/tlsutil"
)

const MAX_ENTITIES memory.Entity = {{.MaxEntities}}

func main() {
	cfg, err := tlsutil.NewServerConfig("./ssl")
	if err != nil {
		log.Fatal(err)
	}
	server := http.Server{
		Addr:      routes.{{.ServiceConstant}},
		TLSConfig: cfg,
	}

	repo := memory.NewRepository(MAX_ENTITIES)
	repo.InitRepository("./ssl")

	h := handler.NewHandler(repo)
{{range .Routes}}	http.Handle("{{.HTTPMethod}} {{.Path}}", http.HandlerFunc(h.{{.MethodName}}))
{{end}}
	err = server.ListenAndServeTLS("./ssl/server.crt", "./ssl/server.key")
	if err != nil {
		log.Fatal(err)
	}
}
`

const handlerTemplate = `package handler

import (
	"crypto/sha256"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Tylores/egot/internal/{{.ServiceName}}/repository/memory"
	"github.com/Tylores/egot/sep"
)

type Handler struct {
	repo *memory.Repository
}

func NewHandler(repo *memory.Repository) *Handler {
	return &Handler{repo}
}

{{range .Methods}}
func (h *Handler) {{.MethodName}}(w http.ResponseWriter, req *http.Request) {
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]

	_, err := h.repo.GetEntity(lfdi)
	if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

{{if .PathParams}}{{range .PathParams}}	_ , err = strconv.Atoi(req.PathValue("{{.Name}}"))
	if err != nil {
		log.Printf("path {{.Name}} value error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
{{end}}
{{end}}	w.Header().Set("Content-Type", sep.ContentType)
	err = xml.NewEncoder(w).Encode(nil)
	if err != nil {
		log.Printf("Response encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
{{end}}
`

const repositoryTemplate = `package memory

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Tylores/egot/internal/{{.ServiceName}}/repository"
)

type Entity uint32

type Pool struct {
	// Add resource storage slices here
	// Example: items []sep.EndDevice
}

func NewPool(size Entity) *Pool {
	return &Pool{
		// Initialize slices here
		// Example: items: make([]sep.EndDevice, size),
	}
}

type Repository struct {
	sync.RWMutex
	tag_lookup      map[string]Entity
	active_entities []bool
	pool            Pool
}

func NewRepository(size Entity) *Repository {
	return &Repository{
		tag_lookup:      make(map[string]Entity),
		pool:            *NewPool(size),
		active_entities: make([]bool, size),
	}
}

func (r *Repository) NextFreeEntity() (*Entity, error) {
	r.Lock()
	defer r.Unlock()

	for i, e := range r.active_entities {
		if e {
			continue
		}
		r.active_entities[i] = true
		entity := Entity(i)
		return &entity, nil
	}

	return nil, repository.ErrPoolFull
}

func (r *Repository) InitRepository(dir string) {
	println("Initializing")
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() || filepath.Ext(path) != ".crt" {
			return nil
		}

		if !strings.Contains(path, "client") {
			return nil
		}

		cert_file, error := os.ReadFile(path)
		if error != nil {
			return error
		}

		block, _ := pem.Decode(cert_file)
		if block == nil {
			return nil
		}

		cert, error := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return error
		}

		lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[:40]
		e, error := r.NextFreeEntity()
		if error != nil {
			return error
		}

		fmt.Printf("\t%s : %d\n", lfdi, *e)
		error = r.TagEntity(lfdi, *e)
		if error != nil {
			return error
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (r *Repository) GetEntity(tag string) (Entity, error) {
	r.RLock()
	defer r.RUnlock()
	value, exists := r.tag_lookup[tag]
	if exists {
		return value, nil
	}
	return value, repository.ErrNotFound
}

func (r *Repository) TagEntity(tag string, e Entity) error {
	r.Lock()
	defer r.Unlock()
	_, exists := r.tag_lookup[tag]
	if exists {
		return repository.ErrTagExists
	}
	r.tag_lookup[tag] = e
	return nil
}
`

const repositoryErrorTemplate = `package repository

import "errors"

var (
	ErrNotFound = errors.New("entity not found")
	ErrTagExists = errors.New("tag already exists")
	ErrPoolFull = errors.New("entity pool is full")
)
`

const serverTemplate = `package server

import (
	"log"
	"net/http"

	"github.com/Tylores/egot/internal/{{.ServiceName}}/handler"
	"github.com/Tylores/egot/internal/{{.ServiceName}}/repository/memory"
	"github.com/Tylores/egot/internal/routes"
	"github.com/Tylores/egot/internal/tlsutil"
)

func AddRoutes(h *handler.Handler) {
	// Add routes from WADL specification
	// Example: http.Handle("GET /resource", http.HandlerFunc(h.GetResource))
}

func ServeHTTPS(entities memory.Entity) {
	cfg, err := tlsutil.NewServerConfig("./ssl")
	if err != nil {
		log.Fatal(err)
	}
	server := http.Server{
		Addr:      routes.{{.ServiceConstant}},
		TLSConfig: cfg,
	}

	repo := memory.NewRepository(entities)
	repo.InitRepository("./ssl")

	h := handler.NewHandler(repo)
	AddRoutes(h)

	err = server.ListenAndServeTLS("./ssl/server.crt", "./ssl/server.key")
	if err != nil {
		log.Fatal(err)
	}
}
`

type TemplateData struct {
	ServiceName      string
	ServiceConstant  string
	MaxEntities      int
	Routes           []RouteInfo
	Methods          []MethodInfo
}

type RouteInfo struct {
	HTTPMethod string
	Path       string
	MethodName string
}

type MethodInfo struct {
	MethodName string
	HTTPMethod string
	Path       string
	PathParams []PathParamInfo
}

type PathParamInfo struct {
	Name      string
	ParamVar  string
	Type      string
}

func (g *Generator) generateMain(outputDir, serviceName string) error {
	routes := []RouteInfo{}
	for _, res := range g.wadl.Resources {
		for _, method := range res.Methods {
			routes = append(routes, RouteInfo{
				HTTPMethod: method.HTTPMethod,
				Path:       res.Path,
				MethodName: method.Name,
			})
		}
	}

	data := TemplateData{
		ServiceName:     serviceName,
		ServiceConstant: toCamelCase(serviceName),
		MaxEntities:     g.wadl.MaxEntities,
		Routes:          routes,
	}

	tmpl, err := template.New("main").Parse(mainTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse main template: %w", err)
	}

	mainPath := filepath.Join(outputDir, fmt.Sprintf("cmd/%s/main.go", serviceName))
	f, err := os.Create(mainPath)
	if err != nil {
		return fmt.Errorf("failed to create main.go: %w", err)
	}
	defer f.Close()

	err = tmpl.Execute(f, data)
	if err != nil {
		return fmt.Errorf("failed to execute main template: %w", err)
	}

	return nil
}

func (g *Generator) generateHandler(outputDir, serviceName string) error {
	methods := []MethodInfo{}
	for _, res := range g.wadl.Resources {
		for _, method := range res.Methods {
			pathParams := []PathParamInfo{}
			for _, pp := range method.PathParams {
				pathParams = append(pathParams, PathParamInfo{
					Name:     pp.Name,
					ParamVar: pp.Name,
					Type:     pp.Type,
				})
			}
			methods = append(methods, MethodInfo{
				MethodName: method.Name,
				HTTPMethod: method.HTTPMethod,
				Path:       res.Path,
				PathParams: pathParams,
			})
		}
	}

	data := TemplateData{
		ServiceName: serviceName,
		Methods:     methods,
	}

	tmpl, err := template.New("handler").Parse(handlerTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse handler template: %w", err)
	}

	handlerPath := filepath.Join(outputDir, fmt.Sprintf("internal/%s/handler/handler.go", serviceName))
	f, err := os.Create(handlerPath)
	if err != nil {
		return fmt.Errorf("failed to create handler.go: %w", err)
	}
	defer f.Close()

	err = tmpl.Execute(f, data)
	if err != nil {
		return fmt.Errorf("failed to execute handler template: %w", err)
	}

	return nil
}

func (g *Generator) generateRepository(outputDir, serviceName string) error {
	data := TemplateData{
		ServiceName: serviceName,
	}

	tmpl, err := template.New("repository").Parse(repositoryTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse repository template: %w", err)
	}

	repoPath := filepath.Join(outputDir, fmt.Sprintf("internal/%s/repository/memory/memory.go", serviceName))
	f, err := os.Create(repoPath)
	if err != nil {
		return fmt.Errorf("failed to create memory.go: %w", err)
	}
	defer f.Close()

	err = tmpl.Execute(f, data)
	if err != nil {
		return fmt.Errorf("failed to execute repository template: %w", err)
	}

	return nil
}

func (g *Generator) generateRepositoryError(outputDir, serviceName string) error {
	errorPath := filepath.Join(outputDir, fmt.Sprintf("internal/%s/repository/error.go", serviceName))
	f, err := os.Create(errorPath)
	if err != nil {
		return fmt.Errorf("failed to create error.go: %w", err)
	}
	defer f.Close()

	_, err = f.WriteString(repositoryErrorTemplate)
	if err != nil {
		return fmt.Errorf("failed to write error.go: %w", err)
	}

	return nil
}

func (g *Generator) generateServer(outputDir, serviceName string) error {
	data := TemplateData{
		ServiceName:     serviceName,
		ServiceConstant: toCamelCase(serviceName),
	}

	tmpl, err := template.New("server").Parse(serverTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse server template: %w", err)
	}

	serverPath := filepath.Join(outputDir, fmt.Sprintf("internal/%s/server/server.go", serviceName))
	f, err := os.Create(serverPath)
	if err != nil {
		return fmt.Errorf("failed to create server.go: %w", err)
	}
	defer f.Close()

	err = tmpl.Execute(f, data)
	if err != nil {
		return fmt.Errorf("failed to execute server template: %w", err)
	}

	return nil
}

func (g *Generator) updateRoutes(outputDir, serviceName string) error {
	routesPath := filepath.Join(outputDir, "internal/routes/routes.go")
	
	// Check if routes file exists - if not, skip (will exist in main project)
	data, err := os.ReadFile(routesPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Routes file doesn't exist, skip update
		}
		return fmt.Errorf("failed to read routes.go: %w", err)
	}

	content := string(data)
	
	// Check if service already exists in routes
	constant := toCamelCase(serviceName)
	if strings.Contains(content, constant) {
		return nil // Already exists
	}

	// Add new route before closing parenthesis
	routeAddr := fmt.Sprintf("%s.internal.com:%d", serviceName, g.wadl.Port)
	newRoute := fmt.Sprintf("\t%s       = \"%s\"\n", constant, routeAddr)
	
	// Find the last route entry and add after it
	lastNewline := strings.LastIndex(content, ")")
	if lastNewline == -1 {
		return fmt.Errorf("invalid routes.go format")
	}

	// Insert new route before closing paren
	newContent := content[:lastNewline] + newRoute + content[lastNewline:]

	err = os.WriteFile(routesPath, []byte(newContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to update routes.go: %w", err)
	}

	return nil
}
