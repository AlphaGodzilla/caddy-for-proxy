package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"strconv"
	"strings"
	"text/template"
)

const (
	PORT            = "PORT"
	DEFAULT_URL     = "DEFAULT_URL"
	TRUSTED_PROXIES = "TRUSTED_PROXIES"
	METRICS         = "METRICS"
	METRICS_PORT    = "METRICS_PORT"
	DEBUG           = "DEBUG"
	LOG_LEVEL       = "LOG_LEVEL"
	PATH_I          = "PATH"
	URL_I           = "URL"
	DOMAIN_I        = "DOMAIN"
	DOMAIN_TO_I     = "DOMAIN_TO"
)

// 对外服务端口
var port uint64

// 默认反代上游
var defaultProxyTo = getEnvOrDefault(DEFAULT_URL, "")

// 信任的代理列表
var trustedProxies string

var metrics bool
var metricsPort uint64

var debug bool

var logLevel = getEnvOrDefault(LOG_LEVEL, "INFO")

type Proxy struct {
	PathPrefix string
	ProxyTo    string
}

var proxies = make([]Proxy, 0, 10)

type Domain struct {
	Domain  string
	ProxyTo string
}

var domains = make([]Domain, 0, 10)

type RenderConfig struct {
	Port           uint64
	DefaultProxyTo string
	TrustedProxies string
	LogLevel       string
	Metrics        bool
	MetricsPort    uint64
	Debug          bool
	Proxies        *[]Proxy
	Domains        *[]Domain
}

func initialize() {
	// port
	portStr := getEnvOrDefault(PORT, "80")
	if port1, err := strconv.ParseUint(portStr, 10, 64); err == nil {
		port = port1
	} else {
		log.Fatalln("PORT parse error "+portStr, err)
	}
	// trustedProxies
	trustedProxies = getEnvOrDefault(TRUSTED_PROXIES, "")
	trustedProxies = strings.ReplaceAll(trustedProxies, " ", "")
	// 如果CIDR的IP是用逗号分隔的，将逗号替换为空格
	trustedProxies = strings.ReplaceAll(trustedProxies, ",", " ")
	// metrics
	metricsStr := getEnvOrDefault(METRICS, "true")
	if _metrics, err := strconv.ParseBool(metricsStr); err == nil {
		metrics = _metrics
	} else {
		log.Fatalf("METRICS parse error "+metricsStr, err)
	}
	metricsPortStr := getEnvOrDefault(METRICS_PORT, "2023")
	if metricsPort1, err := strconv.ParseUint(metricsPortStr, 10, 64); err == nil {
		metricsPort = metricsPort1
	} else {
		log.Fatalln("METRICS_PORT parse error "+metricsPortStr, err)
	}
	debugStr := getEnvOrDefault(DEBUG, "false")
	if _debug, err := strconv.ParseBool(debugStr); err == nil {
		debug = _debug
	} else {
		log.Fatalf("DEBUG parse error "+debugStr, err)
	}
	// path and url
	i := 0
	for {
		pathKey := fmt.Sprintf("%v_%v", PATH_I, i)
		path := os.Getenv(pathKey)
		if path == "" {
			break
		}
		urlKey := fmt.Sprintf("%v_%v", URL_I, i)
		url := os.Getenv(urlKey)
		if url == "" {
			break
		}
		i += 1
		// path url ready
		proxies = append(proxies, Proxy{path, url})
	}
	// multi domain
	j := 0
	for {
		domainKey := fmt.Sprintf("%v_%v", DOMAIN_I, j)
		domain := os.Getenv(domainKey)
		if domain == "" {
			break
		}
		domainToKey := fmt.Sprintf("%v_%v", DOMAIN_TO_I, j)
		domainTo := os.Getenv(domainToKey)
		if domainTo == "" {
			break
		}
		j += 1
		domains = append(domains, Domain{domain, domainTo})
	}
}

func getEnvOrDefault(envKey string, defaultValue string) string {
	value := os.Getenv(envKey)
	if value == "" {
		return defaultValue
	}
	return value
}

func buildCaddyFile(config *RenderConfig) {
	tmpl, err := template.ParseFiles("CaddyfileTemplate")
	if err != nil {
		log.Fatalln("Parse template error <CaddyfileTemplate>", err)
	}

	saveFile, err := os.OpenFile("Caddyfile", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fs.ModePerm)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(saveFile, config)
	if err != nil {
		log.Fatalln("Render template error", err)
	}
	data, err := os.ReadFile("Caddyfile")
	if err != nil {
		panic(err)
	}
	println(string(data))

}

func main() {
	// 初始化
	initialize()
	config := &RenderConfig{
		Port:           port,
		DefaultProxyTo: defaultProxyTo,
		TrustedProxies: trustedProxies,
		Metrics:        metrics,
		MetricsPort:    metricsPort,
		Debug:          debug,
		LogLevel:       logLevel,
		Proxies:        &proxies,
		Domains:        &domains,
	}
	fmt.Printf("%+v\n", config)
	buildCaddyFile(config)
}
