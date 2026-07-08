package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
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
	transportCfg := transportConfigFromEnv()
	if err := transportCfg.Validate(authMode); err != nil {
		log.Fatal(err)
	}

	supportsHTTP := transportCfg.UsesHTTP()
	adminEnabled := supportsHTTP && authMode == security.AuthModeIssuedToken && security.AdminAPITokenFromEnv() != ""
	exchangeCfg := security.TencentCloudTokenExchangeConfigFromEnv(authMode)
	exchangeEnabled := supportsHTTP && exchangeCfg.Enabled
	needsTokenStore := authMode == security.AuthModeIssuedToken || adminEnabled || exchangeEnabled

	var (
		sqliteStore        *security.SQLiteTokenStore
		tokenStore         security.TokenStore
		tokenIssuer        *security.TokenIssuer
		issuerCfg          security.TokenIssuerConfig
		credentialCipher   *security.CredentialCipher
		credentialProvider security.CredentialProvider
		tokenExchange      *security.TencentCloudTokenExchangeService
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
	case security.AuthModeRequestCredential:
		credentialProvider = security.NewRequestHeaderCredentialProvider()
	default:
		credentialProvider, err = security.NewStaticCredentialProviderFromEnv()
		if err != nil {
			log.Fatal(err)
		}
	}

	if exchangeEnabled {
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
		"1.0.3",
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

	if !supportsHTTP {
		stdioCtxFunc, err := buildStdioContextFunc(authMode)
		if err != nil {
			log.Fatal(err)
		}
		logTransportStartup(transportCfg, authMode, false, false)
		if err := server.ServeStdio(mcpsvr, server.WithStdioContextFunc(stdioCtxFunc)); err != nil {
			log.Fatalf("stdio server error: %v", err)
		}
		return
	}

	authenticator, err := security.NewMCPAuthenticator(authMode, tokenStore, issuerCfg.Pepper)
	if err != nil {
		log.Fatalf("init auth failed: %v", err)
	}

	mux := http.NewServeMux()
	security.RegisterHealthRoutes(mux, security.HealthStatus{
		Service:   mcpServerName,
		Version:   "1.0.3",
		StartedAt: startedAt,
	})
	if adminEnabled {
		security.RegisterAdminRoutes(mux, tokenIssuer, tokenStore)
	}
	if tokenExchange != nil {
		security.RegisterTokenExchangeRoutes(mux, tokenExchange, security.TokenExchangeBootstrapConfig{
			MCPServerName: mcpServerName,
			ServerURL:     transportCfg.ServerURL(),
		})
	}

	switch transportCfg.Mode {
	case transportModeSSE:
		sseServer := server.NewSSEServer(mcpsvr,
			server.WithSSEEndpoint(transportCfg.SSEEndpoint),
			server.WithMessageEndpoint(transportCfg.MessageEndpoint),
			server.WithAppendQueryToMessageEndpoint(),
		)
		mux.Handle("/", security.WrapMCPAuth(authenticator, sseServer))
	default:
		transportServer := server.NewStreamableHTTPServer(mcpsvr,
			server.WithEndpointPath(transportCfg.HTTPEndpoint),
			server.WithStateLess(transportCfg.StatelessHTTP),
		)
		mux.Handle(transportCfg.HTTPEndpoint, security.WrapMCPAuth(authenticator, transportServer))
	}

	httpServer := &http.Server{
		Addr:              transportCfg.ListenAddr(),
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       2 * time.Minute,
	}

	logTransportStartup(transportCfg, authMode, adminEnabled, tokenExchange != nil)
	outputMCPServerConfig(mcpServerName, transportCfg, authMode)
	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Server error: %v", err)
	}
}

func outputMCPServerConfig(mcpServerName string, transportCfg transportConfig, authMode security.AuthMode) {
	if !transportCfg.UsesHTTP() {
		return
	}
	serverConfig := map[string]any{
		"url":  transportCfg.ServerURL(),
		"type": transportCfg.ClientType(),
	}
	if headers, ok := clientHeaderPlaceholders(authMode); ok {
		serverConfig["headers"] = headers
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
	case security.AuthModeRequestCredential:
		fmt.Println("Request-credential auth is enabled for this server.")
		fmt.Printf("Every MCP request must carry %s and %s headers. %s is optional for STS temporary credentials.\n", security.RequestSecretIDHeader, security.RequestSecretKeyHeader, security.RequestSessionTokenHeader)
		fmt.Println("Use HTTPS in production and never place secret material in URL query parameters.")
		fmt.Println()
	}
	fmt.Println("The server is now ready to accept connections.")
}

func clientHeaderPlaceholders(authMode security.AuthMode) (map[string]string, bool) {
	switch authMode {
	case security.AuthModeSharedToken:
		return map[string]string{"Authorization": "Bearer <MCP_API_TOKEN>"}, true
	case security.AuthModeIssuedToken:
		return map[string]string{"Authorization": "Bearer <MCP_ACCESS_TOKEN>"}, true
	case security.AuthModeRequestCredential:
		return security.RequestCredentialPlaceholderHeaders(), true
	default:
		return nil, false
	}
}
