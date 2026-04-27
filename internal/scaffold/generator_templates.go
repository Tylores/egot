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
	"path/filepath"

	"github.com/Tylores/egot/internal/{{.ServiceName}}/handler"
	"github.com/Tylores/egot/internal/store"
	"github.com/Tylores/egot/internal/routes"
	"github.com/Tylores/egot/internal/tlsutil"
)

func main() {
	cfg, err := tlsutil.NewServerConfig("./ssl")
	if err != nil {
		log.Fatal(err)
	}
	server := http.Server{
		Addr:      routes.{{.ServiceConstant}},
		TLSConfig: cfg,
	}

	repo := store.New(filepath.Join("data", "{{.ServiceName}}.db"))

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

	"github.com/Tylores/egot/internal/store"
	"github.com/Tylores/egot/sep"
)

type Handler struct {
	repo *store.Store
}

func NewHandler(repo *store.Store) *Handler {
	return &Handler{repo}
}

// getLFDI extracts and validates LFDI from certificate
func (h *Handler) getLFDI(req *http.Request) (string, error) {
	if req.TLS == nil || len(req.TLS.PeerCertificates) == 0 {
		return "0000000000000000000000000000000000000000", nil
	}
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]
	return lfdi, nil
}

// getSFDI extracts SFDI from LFDI
func (h *Handler) getSFDI(lfdi string) (sep.SFDIType, error) {
	return sep.ToSFDI(lfdi)
}

// buildStoreKey builds a store key from SFDI and optional MRID
func (h *Handler) buildStoreKey(sfdi sep.SFDIType, mrid ...string) string {
	key := fmt.Sprintf("%d", sfdi)
	if len(mrid) > 0 && mrid[0] != "" {
		key = fmt.Sprintf("%d:%s", sfdi, mrid[0])
	}
	return key
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
	lfdi, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	sfdi, err := h.getSFDI(lfdi)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_ = h.buildStoreKey(sfdi)
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
	w.Header().Set("location", "{{.Path}}/"+fmt.Sprintf("%d", sfdi))
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

// This file is deprecated. Services now use internal/store.Store directly.
// Kept for backward compatibility during migration.
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
	"path/filepath"

	"github.com/Tylores/egot/internal/{{.ServiceName}}/handler"
	"github.com/Tylores/egot/internal/store"
	"github.com/Tylores/egot/internal/routes"
	"github.com/Tylores/egot/internal/tlsutil"
)

func AddRoutes(h *handler.Handler) {
	// Add routes from WADL specification
	// Example: http.Handle("GET /resource", http.HandlerFunc(h.GetResource))
}

func ServeHTTPS() {
	cfg, err := tlsutil.NewServerConfig("./ssl")
	if err != nil {
		log.Fatal(err)
	}
	server := http.Server{
		Addr:      routes.{{.ServiceConstant}},
		TLSConfig: cfg,
	}

	repo := store.New(filepath.Join("data", "{{.ServiceName}}.db"))

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

		// Find the closing paren of the const block (which ends with "\n)")
		pattern := regexp.MustCompile(`(?m)^const \([^)]*\)`)
		match := pattern.FindString(content)
		if match == "" {
			return fmt.Errorf("invalid routes.go format: const block not found")
		}
		// Find where this match ends and insert before the closing paren
		constEnd := strings.Index(content, "const (")
		if constEnd == -1 {
			return fmt.Errorf("invalid routes.go format: const ( not found")
		}
		// Search forward from const ( to find the closing paren
		parenCount := 1
		i := constEnd + len("const (")
		for i < len(content) && parenCount > 0 {
			if content[i] == '(' {
				parenCount++
			} else if content[i] == ')' {
				parenCount--
			}
			i++
		}
		if parenCount != 0 {
			return fmt.Errorf("invalid routes.go format: unmatched parens")
		}
		lastParen := i - 1 // Position of closing paren
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
