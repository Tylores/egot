package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"
)

// NginxYAML represents the nginx configuration file
type NginxYAML struct {
	Server struct {
		Listen     string `yaml:"listen"`
		Port       int    `yaml:"port"`
		ServerName string `yaml:"server_name"`
	} `yaml:"server"`
	TLS struct {
		Cert   string `yaml:"cert"`
		Key    string `yaml:"key"`
		CACert string `yaml:"ca_cert"`
	} `yaml:"tls"`
	Logging struct {
		AccessLog      string `yaml:"access_log"`
		ErrorLog       string `yaml:"error_log"`
		ErrorLogLevel  string `yaml:"error_log_level"`
	} `yaml:"logging"`
	Proxy struct {
		ConnectTimeout int    `yaml:"connect_timeout"`
		ReadTimeout    int    `yaml:"read_timeout"`
		SendTimeout    int    `yaml:"send_timeout"`
		BufferSize     string `yaml:"buffer_size"`
	} `yaml:"proxy"`
}

type Service struct {
	Name string
	Host string
	Port string
}

type PathMapping struct {
	Path    string
	Service string
}

func main() {
	inputYAML := flag.String("yaml", "", "Path to nginx.yaml configuration")
	outputConf := flag.String("output", "", "Output path for nginx.conf")
	flag.Parse()

	if *inputYAML == "" || *outputConf == "" {
		log.Fatal("Usage: -yaml <input.yaml> -output <output.conf>")
	}

	// Parse YAML configuration
	cfg, err := parseNginxYAML(*inputYAML)
	if err != nil {
		log.Fatalf("Failed to parse YAML: %v", err)
	}

	// Extract services from routes.go
	repoRoot := filepath.Dir(filepath.Dir(*inputYAML))
	routesFile := filepath.Join(repoRoot, "internal", "routes", "routes.go")
	
	services, err := extractServices(routesFile)
	if err != nil {
		log.Fatalf("Failed to extract services: %v", err)
	}

	pathMappings, err := extractPathMappings(routesFile)
	if err != nil {
		log.Fatalf("Failed to extract path mappings: %v", err)
	}

	// Generate upstream blocks
	upstreamBlocks := generateUpstreamBlocks(services)

	// Generate location blocks
	locationBlocks := generateLocationBlocks(pathMappings)

	// Generate nginx.conf from template
	confContent := generateNginxConf(cfg, upstreamBlocks, locationBlocks)

	// Write output file
	if err := os.WriteFile(*outputConf, []byte(confContent), 0644); err != nil {
		log.Fatalf("Failed to write nginx.conf: %v", err)
	}

	fmt.Printf("✅ Generated nginx.conf: %s\n", *outputConf)
	fmt.Printf("📍 Configuration: %s\n", *inputYAML)
	fmt.Printf("📋 Services mapped: %d\n", len(services))
	fmt.Printf("📋 Paths mapped: %d\n", len(pathMappings))
}

func parseNginxYAML(filepath string) (*NginxYAML, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var cfg NginxYAML
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// Set defaults
	if cfg.Server.Listen == "" {
		cfg.Server.Listen = "127.0.0.1"
	}
	if cfg.Server.Port == 0 {
		cfg.Server.Port = 8443
	}
	if cfg.Server.ServerName == "" {
		cfg.Server.ServerName = "localhost"
	}
	if cfg.TLS.Cert == "" {
		cfg.TLS.Cert = "./ssl/server.crt"
	}
	if cfg.TLS.Key == "" {
		cfg.TLS.Key = "./ssl/server.key"
	}
	if cfg.Logging.AccessLog == "" {
		cfg.Logging.AccessLog = "./logs/nginx_access.log"
	}
	if cfg.Logging.ErrorLog == "" {
		cfg.Logging.ErrorLog = "./logs/nginx_error.log"
	}
	if cfg.Logging.ErrorLogLevel == "" {
		cfg.Logging.ErrorLogLevel = "warn"
	}
	if cfg.Proxy.ConnectTimeout == 0 {
		cfg.Proxy.ConnectTimeout = 10
	}
	if cfg.Proxy.ReadTimeout == 0 {
		cfg.Proxy.ReadTimeout = 30
	}
	if cfg.Proxy.SendTimeout == 0 {
		cfg.Proxy.SendTimeout = 30
	}
	if cfg.Proxy.BufferSize == "" {
		cfg.Proxy.BufferSize = "4k"
	}

	return &cfg, nil
}

func extractServices(routesFile string) (map[string]Service, error) {
	services := make(map[string]Service)
	file, err := os.Open(routesFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Regex to match: ServiceName = "host:port"
	re := regexp.MustCompile(`^\s*([A-Z][a-zA-Z0-9]*)\s*=\s*"([^:]+):(\d+)"`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindStringSubmatch(line)
		if len(matches) == 4 {
			name := matches[1]
			host := matches[2]
			port := matches[3]
			services[name] = Service{Name: name, Host: host, Port: port}
		}
	}

	return services, scanner.Err()
}

func extractPathMappings(routesFile string) ([]PathMapping, error) {
	var mappings []PathMapping
	file, err := os.Open(routesFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Regex to match: "/path": ServiceName,
	re := regexp.MustCompile(`"(/[^"]+)"\s*:\s*([A-Z][a-zA-Z0-9]*)`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindStringSubmatch(line)
		if len(matches) == 3 {
			path := matches[1]
			service := matches[2]
			mappings = append(mappings, PathMapping{Path: path, Service: service})
		}
	}

	return mappings, scanner.Err()
}

func generateUpstreamBlocks(services map[string]Service) string {
	var buf strings.Builder
	for _, svc := range services {
		fmt.Fprintf(&buf, "    upstream %s_backend {\n", strings.ToLower(svc.Name))
		fmt.Fprintf(&buf, "        server %s:%s;\n", svc.Host, svc.Port)
		fmt.Fprintf(&buf, "    }\n\n")
	}
	return buf.String()
}

func generateLocationBlocks(pathMappings []PathMapping) string {
	var buf strings.Builder
	for _, pm := range pathMappings {
		// Convert {id} patterns to nginx regex
		nginxPath := pm.Path
		nginxPath = regexp.MustCompile(`\{id[0-9]*\}`).ReplaceAllString(nginxPath, "[^/]+")

		fmt.Fprintf(&buf, "        location ~ ^%s(/.*)?$ {\n", nginxPath)
		fmt.Fprintf(&buf, "            proxy_pass https://%s_backend;\n", strings.ToLower(pm.Service))
		fmt.Fprintf(&buf, "            proxy_set_header Host $host;\n")
		fmt.Fprintf(&buf, "            proxy_set_header X-Real-IP $remote_addr;\n")
		fmt.Fprintf(&buf, "            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;\n")
		fmt.Fprintf(&buf, "            proxy_set_header X-Forwarded-Proto $scheme;\n")
		fmt.Fprintf(&buf, "            proxy_ssl_verify off;\n")
		fmt.Fprintf(&buf, "        }\n\n")
	}
	return buf.String()
}

func generateNginxConf(cfg *NginxYAML, upstreamBlocks, locationBlocks string) string {
	tmpl, err := template.New("nginx").Parse(nginxTemplate)
	if err != nil {
		log.Fatal(err)
	}

	data := map[string]interface{}{
		"ListenAddr":            cfg.Server.Listen,
		"Port":                  cfg.Server.Port,
		"ServerName":            cfg.Server.ServerName,
		"TLSCert":               cfg.TLS.Cert,
		"TLSKey":                cfg.TLS.Key,
		"AccessLog":             cfg.Logging.AccessLog,
		"ErrorLog":              cfg.Logging.ErrorLog,
		"ErrorLogLevel":         cfg.Logging.ErrorLogLevel,
		"ProxyConnectTimeout":   cfg.Proxy.ConnectTimeout,
		"ProxyReadTimeout":      cfg.Proxy.ReadTimeout,
		"ProxySendTimeout":      cfg.Proxy.SendTimeout,
		"ProxyBufferSize":       cfg.Proxy.BufferSize,
		"UpstreamBlocks":        upstreamBlocks,
		"LocationBlocks":        locationBlocks,
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		log.Fatal(err)
	}

	return buf.String()
}

const nginxTemplate = `user www-data www-data;
worker_processes auto;
error_log {{.ErrorLog}} {{.ErrorLogLevel}};
pid /var/run/nginx.pid;

events {
    worker_connections 1024;
}

http {
    include /etc/nginx/mime.types;
    default_type application/octet-stream;

    log_format main '$remote_addr - $remote_user [$time_local] "$request" '
                    '$status $body_bytes_sent "$http_referer" '
                    '"$http_user_agent" "$http_x_forwarded_for"';

    access_log {{.AccessLog}} main;

    sendfile on;
    tcp_nopush on;
    tcp_nodelay on;
    keepalive_timeout 65;
    types_hash_max_size 2048;

    # Upstream service definitions
{{.UpstreamBlocks}}
    # Main server block for API gateway
    server {
        listen {{.ListenAddr}}:{{.Port}} ssl http2;
        server_name {{.ServerName}};

        # TLS configuration
        ssl_certificate {{.TLSCert}};
        ssl_certificate_key {{.TLSKey}};
        ssl_protocols TLSv1.2 TLSv1.3;
        ssl_ciphers HIGH:!aNULL:!MD5;

        # Request body size limit
        client_max_body_size 10m;

        # Proxy timeouts
        proxy_connect_timeout {{.ProxyConnectTimeout}}s;
        proxy_read_timeout {{.ProxyReadTimeout}}s;
        proxy_send_timeout {{.ProxySendTimeout}}s;
        proxy_buffer_size {{.ProxyBufferSize}};

        # Service routing location blocks
{{.LocationBlocks}}        # Default catch-all for unmapped paths
        location / {
            return 404 "No service mapping for this path";
        }
    }
}
`
