package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/server"
	"postgres_server/config"
	"postgres_server/security"
	"postgres_server/tools"
)

func main() {
	config.Init()
	startedAt := time.Now().UTC()

	authMode := security.MCPAuthModeFromEnv()
	adminEnabled := security.AdminAPITokenFromEnv() != ""
	exchangeCfg := security.TencentCloudTokenExchangeConfigFromEnv(authMode)
	needsTokenStore := authMode == security.AuthModeIssuedToken || adminEnabled || exchangeCfg.Enabled

	var (
		sqliteStore        *security.SQLiteTokenStore
		tokenStore         security.TokenStore
		tokenIssuer        *security.TokenIssuer
		issuerCfg          security.TokenIssuerConfig
		credentialCipher   *security.CredentialCipher
		credentialProvider security.CredentialProvider
		tokenExchange      *security.TencentCloudTokenExchangeService
		staticProvider     *security.StaticCredentialProvider
		err                error
	)

	if needsTokenStore {
		issuerCfg = security.TokenIssuerConfigFromEnv()
		sqliteStore, err = security.NewSQLiteTokenStoreFromEnv()
		if err != nil {
			log.Fatalf("init token store failed: %v", err)
		}
		tokenStore = sqliteStore
		tokenIssuer = security.NewTokenIssuer(sqliteStore, issuerCfg)
	}

	switch authMode {
	case security.AuthModeIssuedToken:
		credentialCipher, err = security.NewCredentialCipherFromEnv()
		if err != nil {
			log.Fatal(err)
		}
		credentialProvider, err = security.NewTokenBoundCredentialProvider(sqliteStore, credentialCipher)
		if err != nil {
			log.Fatal(err)
		}
	default:
		staticProvider, err = security.NewStaticCredentialProviderFromEnv()
		if err != nil {
			log.Fatal(err)
		}
		credentialProvider = staticProvider
		if staticProvider.UsingLegacyEnv() {
			log.Printf("warning: using legacy TENCENTCLOUD_SECRET_* envs; please migrate to MCP_SECRET_*")
		}
	}

	if exchangeCfg.Enabled {
		if credentialCipher == nil {
			credentialCipher, err = security.NewCredentialCipherFromEnv()
			if err != nil {
				log.Fatal(err)
			}
		}
		tokenExchange, err = security.NewTencentCloudTokenExchangeService(tokenIssuer, sqliteStore, credentialCipher, exchangeCfg)
		if err != nil {
			log.Fatal(err)
		}
	}

	guard := security.NewGuard()
	mcpServerName := "mcp-server-postgres"
	mcpsvr := server.NewMCPServer(
		"腾讯云 Postgres MCP",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)

	toolCount := 0
	if config.IsEnabled("instance") {
		tools.RegisterInstanceTools(mcpsvr, credentialProvider, guard)
		toolCount += 15
	}
	if config.IsEnabled("account") {
		tools.RegisterAccountTools(mcpsvr, credentialProvider, guard)
		toolCount += 6
	}
	if config.IsEnabled("database") {
		tools.RegisterDatabaseTools(mcpsvr, credentialProvider, guard)
		toolCount += 4
	}
	if config.IsEnabled("parameter") {
		tools.RegisterParameterTools(mcpsvr, credentialProvider, guard)
		toolCount += 5
	}
	if config.IsEnabled("backup") {
		tools.RegisterBackupTools(mcpsvr, credentialProvider, guard)
		toolCount += 8
	}
	if config.IsEnabled("monitoring") {
		tools.RegisterMonitoringTools(mcpsvr, credentialProvider, guard)
		toolCount += 3
	}
	if config.IsEnabled("network") {
		tools.RegisterNetworkTools(mcpsvr, credentialProvider, guard)
		toolCount += 4
	}
	if config.IsEnabled("readonly") {
		tools.RegisterReadonlyTools(mcpsvr, credentialProvider, guard)
		toolCount += 2
	}
	tools.RegisterSSLTools(mcpsvr, credentialProvider, guard)
	toolCount += 1
	log.Printf("Total tools registered: %d", toolCount)

	authenticator, err := security.NewMCPAuthenticator(authMode, tokenStore, issuerCfg.Pepper)
	if err != nil {
		log.Fatalf("init auth failed: %v", err)
	}

	sseEndpoint := getEnv("MCP_SERVER_SSE_ENDPOINT", "/sse")
	messageEndpoint := getEnv("MCP_SERVER_MESSAGE_ENDPOINT", "/message")
	ssePort := getEnv("MCP_SERVER_PORT", getEnv("MCP_SERVER_SSE_PORT", "9000"))
	bindHost := getEnv("MCP_SERVER_BIND_HOST", "127.0.0.1")
	listenAddr := fmt.Sprintf("%s:%s", bindHost, ssePort)
	serverURL := buildServerURL(bindHost, ssePort, sseEndpoint)

	sseServer := server.NewSSEServer(mcpsvr,
		server.WithSSEEndpoint(sseEndpoint),
		server.WithMessageEndpoint(messageEndpoint),
		server.WithAppendQueryToMessageEndpoint())

	mux := http.NewServeMux()
	security.RegisterHealthRoutes(mux, security.HealthStatus{
		Service:   mcpServerName,
		Version:   "1.0.0",
		StartedAt: startedAt,
	})
	if adminEnabled {
		security.RegisterAdminRoutes(mux, tokenIssuer, tokenStore)
	}
	if tokenExchange != nil {
		security.RegisterTokenExchangeRoutes(mux, tokenExchange, security.TokenExchangeBootstrapConfig{
			MCPServerName: mcpServerName,
			ServerURL:     serverURL,
		})
	}
	mux.Handle("/", security.WrapMCPAuth(authenticator, sseServer))

	httpServer := &http.Server{
		Addr:              listenAddr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       2 * time.Minute,
	}

	log.Printf("SSE server listening on %s", listenAddr)
	switch authMode {
	case security.AuthModeSharedToken:
		log.Printf("MCP auth mode: shared-token (%s)", security.APIAuthFallbackHeader)
	case security.AuthModeIssuedToken:
		log.Printf("MCP auth mode: issued-token")
		log.Printf("Token store: %s", security.TokenStorePathFromEnv())
		log.Printf("Credential source: token-bound dynamic credentials")
	default:
		log.Printf("MCP auth mode: none")
	}
	if adminEnabled {
		log.Printf("Admin token API enabled on /admin/tokens")
	} else if authMode == security.AuthModeIssuedToken {
		log.Printf("warning: issued-token mode is enabled but %s is empty; admin token API is disabled", security.AdminAPITokenEnv)
	}
	if tokenExchange != nil {
		log.Printf("TencentCloud token exchange enabled on /auth/token-exchange/tencentcloud (%s)", exchangeCfg.Mode)
	}
	log.Printf("Health check endpoints enabled on /healthz and /readyz")

	outputMCPServerConfig(mcpServerName, serverURL, authMode)
	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Server error: %v", err)
	}
}

func outputMCPServerConfig(mcpServerName, serverURL string, authMode security.AuthMode) {
	serverConfig := map[string]any{
		"url":  serverURL,
		"type": "sse",
	}
	if tokenPlaceholder, ok := clientTokenPlaceholder(authMode); ok {
		serverConfig["headers"] = map[string]string{
			"Authorization": "Bearer " + tokenPlaceholder,
		}
	}

	configPayload := map[string]any{
		"mcpServers": map[string]any{
			mcpServerName: serverConfig,
		},
	}

	jsonOutput, _ := json.MarshalIndent(configPayload, "", " ")
	fmt.Println("=== MCP Server Configuration ===")
	fmt.Println("Copy the following configuration to your MCP client:")
	fmt.Println()
	fmt.Println(string(jsonOutput))
	fmt.Println()
	switch authMode {
	case security.AuthModeSharedToken:
		fmt.Println("Authentication is enabled for this server.")
		fmt.Printf("Replace <MCP_API_TOKEN> with the shared token, or use header %s if your client does not support Bearer config.\n", security.APIAuthFallbackHeader)
		fmt.Println()
	case security.AuthModeIssuedToken:
		fmt.Println("Issued-token auth is enabled for this server.")
		fmt.Println("Use POST /auth/token-exchange/tencentcloud to exchange TencentCloud credentials for a local MCP access token, or let an admin create one via POST /admin/tokens.")
		fmt.Println()
	}
	fmt.Println("The server is now ready to accept connections.")
}

func clientTokenPlaceholder(authMode security.AuthMode) (string, bool) {
	switch authMode {
	case security.AuthModeSharedToken:
		return "<MCP_API_TOKEN>", true
	case security.AuthModeIssuedToken:
		return "<MCP_ACCESS_TOKEN>", true
	default:
		return "", false
	}
}

func buildServerURL(bindHost, port, endpoint string) string {
	if publicURL := strings.TrimSpace(os.Getenv("MCP_SERVER_PUBLIC_URL")); publicURL != "" {
		return publicURL
	}
	host := strings.TrimSpace(bindHost)
	switch host {
	case "", "0.0.0.0", "::":
		host = "127.0.0.1"
	}
	return fmt.Sprintf("http://%s:%s%s", host, port, endpoint)
}

func getEnv(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}
