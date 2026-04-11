package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
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
{{range .Routes}}	http.Handle("{{.HTTPMethod}} {{.Path}}", http.HandlerFunc(h.{{.FuncName}}))
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
	"net/http"
{{- if .HasPathParams}}
	"strconv"
{{- end}}

	"github.com/Tylores/egot/internal/{{.ServiceName}}/repository/memory"
	"github.com/Tylores/egot/sep"
)

type Handler struct {
	repo *memory.Repository
}

func NewHandler(repo *memory.Repository) *Handler {
	return &Handler{repo}
}

// getLFDI extracts and validates LFDI from certificate
func (h *Handler) getLFDI(req *http.Request) (string, error) {
	if req.TLS == nil || len(req.TLS.PeerCertificates) == 0 {
		return "0000000000000000000000000000000000000000", nil
	}
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]
	if _, err := h.repo.GetEntity(lfdi); err != nil {
		return "0000000000000000000000000000000000000000", nil
	}
	return lfdi, nil
}
{{range .Resources}}
// {{.Name}} resource handlers
{{range .Methods}}
func (h *Handler) {{.FuncName}}(w http.ResponseWriter, req *http.Request) {
{{- if eq .Mode "E"}}
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
{{- else}}
	_, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
{{range .PathParams}}	if _, err := strconv.Atoi(req.PathValue("{{.}}")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
{{end -}}
{{- if eq .HTTPMethod "GET"}}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
{{- if .ResponseType}}
	xml.NewEncoder(w).Encode(&sep.{{.ResponseType}}{})
{{- end}}
{{- else if eq .HTTPMethod "HEAD"}}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
{{- else if eq .HTTPMethod "POST"}}
	w.Header().Set("Content-Type", sep.ContentType)
	w.Header().Set("location", "{{.Path}}/1")
	w.WriteHeader(http.StatusCreated)
{{- else if eq .HTTPMethod "DELETE"}}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
{{- if .ResponseType}}
	xml.NewEncoder(w).Encode(&sep.{{.ResponseType}}{})
{{- end}}
{{- else if eq .HTTPMethod "PUT"}}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusOK)
{{- end}}
{{- end}}
}
{{end}}{{end}}`

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

// TemplateData holds all data passed to file-generation templates.
type TemplateData struct {
	ServiceName     string
	ServiceConstant string
	MaxEntities     int
	Routes          []RouteInfo
	Resources       []ResourceTemplateData
	HasPathParams   bool
}

// RouteInfo is one http.Handle registration line.
type RouteInfo struct {
	HTTPMethod string
	Path       string
	FuncName   string
}

// ResourceTemplateData groups methods under a named resource.
type ResourceTemplateData struct {
	Name    string
	Methods []MethodInfo
}

// MethodInfo describes a single handler function.
type MethodInfo struct {
	FuncName     string
	HTTPMethod   string
	Mode         string
	ResponseType string
	PathParams   []string
	Path         string
}

func (g *Generator) generateMain(outputDir, serviceName string) error {
	var routes []RouteInfo
	for _, res := range g.spec.Resources {
		for _, m := range res.Methods {
			routes = append(routes, RouteInfo{
				HTTPMethod: m.HTTPMethod,
				Path:       res.Path,
				FuncName:   m.FuncName,
			})
		}
	}

	data := TemplateData{
		ServiceName:     serviceName,
		ServiceConstant: toCamelCase(serviceName),
		MaxEntities:     g.spec.MaxEntities,
		Routes:          routes,
	}

	return renderTemplate("main", mainTemplate, filepath.Join(outputDir, fmt.Sprintf("cmd/%s/main.go", serviceName)), data)
}

func (g *Generator) generateHandler(outputDir, serviceName string) error {
	hasPathParams := false
	var resources []ResourceTemplateData

	for _, res := range g.spec.Resources {
		rtd := ResourceTemplateData{Name: res.Name}
		for _, m := range res.Methods {
			if len(m.PathParams) > 0 {
				hasPathParams = true
			}
			rtd.Methods = append(rtd.Methods, MethodInfo{
				FuncName:     m.FuncName,
				HTTPMethod:   m.HTTPMethod,
				Mode:         m.Mode,
				ResponseType: m.ResponseType,
				PathParams:   m.PathParams,
				Path:         res.Path,
			})
		}
		resources = append(resources, rtd)
	}

	data := TemplateData{
		ServiceName:   serviceName,
		Resources:     resources,
		HasPathParams: hasPathParams,
	}

	return renderTemplate("handler", handlerTemplate, filepath.Join(outputDir, fmt.Sprintf("internal/%s/handler/handler.go", serviceName)), data)
}

func (g *Generator) generateRepository(outputDir, serviceName string) error {
	data := TemplateData{ServiceName: serviceName}
	return renderTemplate("repository", repositoryTemplate, filepath.Join(outputDir, fmt.Sprintf("internal/%s/repository/memory/memory.go", serviceName)), data)
}

func (g *Generator) generateRepositoryError(outputDir, serviceName string) error {
	errorPath := filepath.Join(outputDir, fmt.Sprintf("internal/%s/repository/error.go", serviceName))
	return os.WriteFile(errorPath, []byte(repositoryErrorTemplate), 0644)
}

func (g *Generator) generateServer(outputDir, serviceName string) error {
	data := TemplateData{
		ServiceName:     serviceName,
		ServiceConstant: toCamelCase(serviceName),
	}
	return renderTemplate("server", serverTemplate, filepath.Join(outputDir, fmt.Sprintf("internal/%s/server/server.go", serviceName)), data)
}

func (g *Generator) updateRoutes(outputDir, serviceName string) error {
	routesPath := filepath.Join(outputDir, "internal/routes/routes.go")

	data, err := os.ReadFile(routesPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read routes.go: %w", err)
	}

	content := string(data)
	constant := toCamelCase(serviceName)

	// Insert the host:port constant into the const block if not already present.
	if _, exists := findRouteAddress(content, constant); !exists {
		routeAddr := nextRouteAddress(content)
		newRoute := fmt.Sprintf("\t%s = \"%s\"\n", constant, routeAddr)

		lastParen := strings.LastIndex(content, ")")
		if lastParen == -1 {
			return fmt.Errorf("invalid routes.go format")
		}
		content = content[:lastParen] + newRoute + content[lastParen:]
	}

	// Insert each resource path into serviceMap if not already present.
	// The end of serviceMap is identified by the closing brace before LinkFor.
	const mapEnd = "}\n\n// LinkFor"
	mapEndIdx := strings.Index(content, mapEnd)
	if mapEndIdx == -1 {
		return fmt.Errorf("routes: could not find serviceMap closing brace in routes.go")
	}
	insertAt := mapEndIdx // insert new entries before the closing }

	for _, res := range g.spec.Resources {
		path := res.Path
		if path == "" {
			continue
		}
		entry := fmt.Sprintf("\t%q: %s,\n", path, constant)
		// Skip paths already registered (idempotent).
		if strings.Contains(content, fmt.Sprintf("%q:", path)) {
			continue
		}
		content = content[:insertAt] + entry + content[insertAt:]
		insertAt += len(entry)
	}

	return os.WriteFile(routesPath, []byte(content), 0644)
}

func findRouteAddress(content, constant string) (string, bool) {
	pattern := regexp.MustCompile(`(?m)^\s*` + regexp.QuoteMeta(constant) + `\s*=\s*"([^"]+)"`)
	matches := pattern.FindStringSubmatch(content)
	if len(matches) != 2 {
		return "", false
	}
	return matches[1], true
}

func nextRouteAddress(content string) string {
	pattern := regexp.MustCompile(`"egot\.internal\.com:(\d+)"`)
	matches := pattern.FindAllStringSubmatch(content, -1)

	maxPort := 7999
	for _, match := range matches {
		port, err := strconv.Atoi(match[1])
		if err != nil {
			continue
		}
		if port > maxPort {
			maxPort = port
		}
	}

	return fmt.Sprintf("egot.internal.com:%d", maxPort+1)
}

func renderTemplate(name, tmplStr, outputPath string, data any) error {
	tmpl, err := template.New(name).Parse(tmplStr)
	if err != nil {
		return fmt.Errorf("failed to parse %s template: %w", name, err)
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create %s: %w", outputPath, err)
	}
	defer f.Close()

	if err = tmpl.Execute(f, data); err != nil {
		return fmt.Errorf("failed to execute %s template: %w", name, err)
	}
	return nil
}
