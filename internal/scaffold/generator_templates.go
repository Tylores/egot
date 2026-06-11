package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/Tylores/egot/internal/routes"
)

const mainTemplate = `package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/Tylores/egot/internal/{{.ServiceName}}/handler"
	"github.com/Tylores/egot/internal/store"
	"github.com/Tylores/egot/internal/registry"
	"github.com/Tylores/egot/internal/routes"
	"github.com/Tylores/egot/internal/tlsutil"
)

func main() {
	cfg, err := tlsutil.NewServerConfig("./ssl")
	if err != nil {
		log.Fatal(err)
	}
	server := http.Server{
		Addr:              routes.{{.ServiceConstant}},
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		TLSConfig:         cfg,
		Handler:           tlsutil.CertHeaderMiddleware(http.DefaultServeMux),
	}

	reg := registry.New(filepath.Join("data", "registry.db"))
	if err := reg.Load(); err != nil {
		log.Fatal(err)
	}
	defer reg.Close()
	// Auto-populate registry from known client certs
	_ = reg.PopulateFromCertDir("./ssl")

	repo := store.New(filepath.Join("data", "{{.ServiceName}}.db"))
	if err := repo.Load(); err != nil {
		log.Fatal(err)
	}
	defer repo.Close()

	h := handler.NewHandler(repo, reg)
{{range .Routes}}	http.Handle("{{.HTTPMethod}} {{.Path}}", http.HandlerFunc(h.{{.FuncName}}))
{{end}}
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Starting {{.ServiceName}} on %s", routes.{{.ServiceConstant}})
		err = server.ListenAndServeTLS("./ssl/server.crt", "./ssl/server.key")
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	<-sigChan
	log.Println("Shutting down {{.ServiceName}} server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
	log.Println("Database connections closed.")
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
	"github.com/Tylores/egot/internal/registry"
	"github.com/Tylores/egot/sep"
)

type Handler struct {
	repo *store.Store
	reg  *registry.Registry
}

func NewHandler(repo *store.Store, reg *registry.Registry) *Handler {
	return &Handler{repo, reg}
}

// getLFDI extracts and validates LFDI from certificate against the registry
func (h *Handler) getLFDI(req *http.Request) (string, error) {
	if req.TLS == nil || len(req.TLS.PeerCertificates) == 0 {
		return "", fmt.Errorf("mTLS certificate required")
	}
	cert := req.TLS.PeerCertificates[0]
	lfdi := fmt.Sprintf("%X", sha256.Sum256(cert.Raw))[0:40]
	
	if !h.reg.IsAuthorized(lfdi) {
		return "", fmt.Errorf("device %s not authorized", lfdi)
	}
	
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
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", sep.ContentType)
	w.WriteHeader(http.StatusMethodNotAllowed)
{{- else}}
	lfdi, err := h.getLFDI(req)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
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

const serverTemplate = `package server

import (
	"time"
	"log"
	"net/http"
	"path/filepath"

	"github.com/Tylores/egot/internal/{{.ServiceName}}/handler"
	"github.com/Tylores/egot/internal/store"
	"github.com/Tylores/egot/internal/registry"
	"github.com/Tylores/egot/internal/routes"
	"github.com/Tylores/egot/internal/tlsutil"
)

func AddRoutes(h *handler.Handler) {
	// Add routes from WADL specification
}

func ServeHTTPS() {
	cfg, err := tlsutil.NewServerConfig("./ssl")
	if err != nil {
		log.Fatal(err)
	}
	server := http.Server{
		Addr:              routes.{{.ServiceConstant}},
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		TLSConfig:         cfg,
	}

	reg := registry.New(filepath.Join("data", "registry.db"))
	if err := reg.Load(); err != nil {
		log.Fatal(err)
	}
	defer reg.Close()

	repo := store.New(filepath.Join("data", "{{.ServiceName}}.db"))
	if err := repo.Load(); err != nil {
		log.Fatal(err)
	}
	defer repo.Close()

	h := handler.NewHandler(repo, reg)
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
	path := filepath.Join(outputDir, fmt.Sprintf("internal/%s/handler/handler.go", serviceName))
	if _, err := os.Stat(path); err == nil {
		fmt.Printf("  - Skipping handler generation: %s already exists\n", path)
		return nil
	}

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

	return renderTemplate("handler", handlerTemplate, path, data)
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
	constant := toCamelCase(serviceName)

	var paths []string
	for _, res := range g.spec.Resources {
		if res.Path != "" {
			paths = append(paths, res.Path)
		}
	}

	return routes.UpdateServiceRegistration(routesPath, serviceName, constant, paths)
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
