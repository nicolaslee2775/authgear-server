// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package oauth

import (
	"github.com/skygeario/skygear-server/pkg/auth"
	auth2 "github.com/skygeario/skygear-server/pkg/auth/dependency/auth"
	redis3 "github.com/skygeario/skygear-server/pkg/auth/dependency/auth/redis"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/authenticator/bearertoken"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/authenticator/oob"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/authenticator/password"
	provider2 "github.com/skygeario/skygear-server/pkg/auth/dependency/authenticator/provider"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/authenticator/recoverycode"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/authenticator/totp"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/challenge"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/hook"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/identity/anonymous"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/identity/loginid"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/identity/oauth"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/identity/provider"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/interaction"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/interaction/flows"
	redis2 "github.com/skygeario/skygear-server/pkg/auth/dependency/interaction/redis"
	oauth2 "github.com/skygeario/skygear-server/pkg/auth/dependency/oauth"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/oauth/handler"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/oauth/pq"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/oauth/redis"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/oidc"
	handler2 "github.com/skygeario/skygear-server/pkg/auth/dependency/oidc/handler"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/session"
	redis4 "github.com/skygeario/skygear-server/pkg/auth/dependency/session/redis"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/urlprefix"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/user"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/webapp"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/welcomemessage"
	"github.com/skygeario/skygear-server/pkg/core/async"
	"github.com/skygeario/skygear-server/pkg/core/config"
	"github.com/skygeario/skygear-server/pkg/core/db"
	"github.com/skygeario/skygear-server/pkg/core/logging"
	"github.com/skygeario/skygear-server/pkg/clock"
	"net/http"
)

// Injectors from wire.go:

func newAuthorizeHandler(r *http.Request, m auth.DependencyMap) http.Handler {
	context := auth.ProvideContext(r)
	tenantConfiguration := auth.ProvideTenantConfig(context, m)
	factory := logging.ProvideLoggerFactory(context, tenantConfiguration)
	pool := _wirePoolValue
	dbContext := db.ProvideContextOLD(context, pool, tenantConfiguration)
	sqlBuilder := db.ProvideSQLBuilderOLD(tenantConfiguration)
	sqlExecutor := db.ProvideSQLExecutor(dbContext)
	authorizationStore := &pq.AuthorizationStore{
		SQLBuilder:  sqlBuilder,
		SQLExecutor: sqlExecutor,
	}
	clock := _wireSystemClockValue
	grantStore := redis.ProvideGrantStore(context, factory, tenantConfiguration, sqlBuilder, sqlExecutor, clock)
	urlprefixProvider := urlprefix.NewProvider(r)
	endpointsProvider := &auth.EndpointsProvider{
		PrefixProvider: urlprefixProvider,
	}
	urlProvider := &handler.URLProvider{
		Endpoints: endpointsProvider,
	}
	isAnonymousIdentityEnabled := flows.ProvideIsAnonymousIdentityEnabled(tenantConfiguration)
	store := redis2.ProvideStore(context, tenantConfiguration, clock)
	reservedNameChecker := auth.ProvideReservedNameChecker(m)
	typeCheckerFactory := loginid.ProvideTypeCheckerFactory(tenantConfiguration, reservedNameChecker)
	checker := loginid.ProvideChecker(tenantConfiguration, typeCheckerFactory)
	normalizerFactory := loginid.ProvideNormalizerFactory(tenantConfiguration)
	loginidProvider := loginid.ProvideProvider(sqlBuilder, sqlExecutor, clock, tenantConfiguration, checker, normalizerFactory)
	oauthProvider := oauth.ProvideProvider(sqlBuilder, sqlExecutor, clock)
	anonymousProvider := anonymous.ProvideProvider(sqlBuilder, sqlExecutor)
	providerProvider := provider.ProvideProvider(tenantConfiguration, loginidProvider, oauthProvider, anonymousProvider)
	historyStoreImpl := password.ProvideHistoryStore(clock, sqlBuilder, sqlExecutor)
	passwordChecker := password.ProvideChecker(tenantConfiguration, historyStoreImpl)
	passwordProvider := password.ProvideProvider(sqlBuilder, sqlExecutor, clock, factory, historyStoreImpl, passwordChecker, tenantConfiguration)
	totpProvider := totp.ProvideProvider(sqlBuilder, sqlExecutor, clock, tenantConfiguration)
	engine := auth.ProvideTemplateEngine(tenantConfiguration, m)
	executor := auth.ProvideTaskExecutor(m)
	queue := async.ProvideTaskQueue(context, dbContext, tenantConfiguration, executor)
	oobProvider := oob.ProvideProvider(context, tenantConfiguration, sqlBuilder, sqlExecutor, clock, engine, urlprefixProvider, queue)
	bearertokenProvider := bearertoken.ProvideProvider(sqlBuilder, sqlExecutor, clock, tenantConfiguration)
	recoverycodeProvider := recoverycode.ProvideProvider(sqlBuilder, sqlExecutor, clock, tenantConfiguration)
	provider3 := &provider2.Provider{
		Password:     passwordProvider,
		TOTP:         totpProvider,
		OOBOTP:       oobProvider,
		BearerToken:  bearertokenProvider,
		RecoveryCode: recoverycodeProvider,
	}
	userStore := &user.Store{
		SQLBuilder:  sqlBuilder,
		SQLExecutor: sqlExecutor,
	}
	welcomemessageProvider := welcomemessage.ProvideProvider(context, tenantConfiguration, engine, queue)
	rawCommands := user.ProvideRawCommands(userStore, clock, urlprefixProvider, queue, tenantConfiguration, welcomemessageProvider)
	queries := &user.Queries{
		Store:      userStore,
		Identities: providerProvider,
	}
	hookUserProvider := &auth.HookUserProvider{
		Queries:     queries,
		RawCommands: rawCommands,
	}
	hookProvider := hook.ProvideHookProvider(context, sqlBuilder, sqlExecutor, tenantConfiguration, dbContext, clock, hookUserProvider, loginidProvider, factory)
	commands := &user.Commands{
		Raw:   rawCommands,
		Hooks: hookProvider,
	}
	userProvider := &user.Provider{
		Commands: commands,
		Queries:  queries,
	}
	interactionProvider := interaction.ProvideProvider(store, clock, factory, providerProvider, provider3, userProvider, oobProvider, tenantConfiguration, hookProvider)
	provider4 := challenge.ProvideProvider(context, clock, tenantConfiguration)
	anonymousFlow := &flows.AnonymousFlow{
		Enabled:      isAnonymousIdentityEnabled,
		Interactions: interactionProvider,
		Anonymous:    anonymousProvider,
		Challenges:   provider4,
	}
	stateStoreImpl := &webapp.StateStoreImpl{
		Context: context,
	}
	webappURLProvider := &webapp.URLProvider{
		Endpoints: endpointsProvider,
		Anonymous: anonymousFlow,
		States:    stateStoreImpl,
	}
	scopesValidator := _wireScopesValidatorValue
	tokenGenerator := _wireTokenGeneratorValue
	authorizationHandler := handler.ProvideAuthorizationHandler(context, tenantConfiguration, factory, authorizationStore, grantStore, urlProvider, webappURLProvider, scopesValidator, tokenGenerator, clock)
	httpHandler := provideAuthorizeHandler(factory, dbContext, authorizationHandler)
	return httpHandler
}

var (
	_wirePoolValue            = (*db.Pool)(nil)
	_wireSystemClockValue     = clock.NewSystemClock()
	_wireScopesValidatorValue = handler.ScopesValidator(oidc.ValidateScopes)
	_wireTokenGeneratorValue  = handler.TokenGenerator(oauth2.GenerateToken)
)

func newTokenHandler(r *http.Request, m auth.DependencyMap) http.Handler {
	context := auth.ProvideContext(r)
	tenantConfiguration := auth.ProvideTenantConfig(context, m)
	factory := logging.ProvideLoggerFactory(context, tenantConfiguration)
	pool := _wirePoolValue
	dbContext := db.ProvideContextOLD(context, pool, tenantConfiguration)
	sqlBuilder := db.ProvideSQLBuilderOLD(tenantConfiguration)
	sqlExecutor := db.ProvideSQLExecutor(dbContext)
	authorizationStore := &pq.AuthorizationStore{
		SQLBuilder:  sqlBuilder,
		SQLExecutor: sqlExecutor,
	}
	clock := _wireSystemClockValue
	grantStore := redis.ProvideGrantStore(context, factory, tenantConfiguration, sqlBuilder, sqlExecutor, clock)
	eventStore := redis3.ProvideEventStore(context, tenantConfiguration)
	accessEventProvider := auth2.AccessEventProvider{
		Store: eventStore,
	}
	store := redis4.ProvideStore(context, tenantConfiguration, clock, factory)
	authAccessEventProvider := &auth2.AccessEventProvider{
		Store: eventStore,
	}
	sessionProvider := session.ProvideSessionProvider(r, store, authAccessEventProvider, tenantConfiguration)
	isAnonymousIdentityEnabled := flows.ProvideIsAnonymousIdentityEnabled(tenantConfiguration)
	redisStore := redis2.ProvideStore(context, tenantConfiguration, clock)
	reservedNameChecker := auth.ProvideReservedNameChecker(m)
	typeCheckerFactory := loginid.ProvideTypeCheckerFactory(tenantConfiguration, reservedNameChecker)
	checker := loginid.ProvideChecker(tenantConfiguration, typeCheckerFactory)
	normalizerFactory := loginid.ProvideNormalizerFactory(tenantConfiguration)
	loginidProvider := loginid.ProvideProvider(sqlBuilder, sqlExecutor, clock, tenantConfiguration, checker, normalizerFactory)
	oauthProvider := oauth.ProvideProvider(sqlBuilder, sqlExecutor, clock)
	anonymousProvider := anonymous.ProvideProvider(sqlBuilder, sqlExecutor)
	providerProvider := provider.ProvideProvider(tenantConfiguration, loginidProvider, oauthProvider, anonymousProvider)
	historyStoreImpl := password.ProvideHistoryStore(clock, sqlBuilder, sqlExecutor)
	passwordChecker := password.ProvideChecker(tenantConfiguration, historyStoreImpl)
	passwordProvider := password.ProvideProvider(sqlBuilder, sqlExecutor, clock, factory, historyStoreImpl, passwordChecker, tenantConfiguration)
	totpProvider := totp.ProvideProvider(sqlBuilder, sqlExecutor, clock, tenantConfiguration)
	engine := auth.ProvideTemplateEngine(tenantConfiguration, m)
	urlprefixProvider := urlprefix.NewProvider(r)
	executor := auth.ProvideTaskExecutor(m)
	queue := async.ProvideTaskQueue(context, dbContext, tenantConfiguration, executor)
	oobProvider := oob.ProvideProvider(context, tenantConfiguration, sqlBuilder, sqlExecutor, clock, engine, urlprefixProvider, queue)
	bearertokenProvider := bearertoken.ProvideProvider(sqlBuilder, sqlExecutor, clock, tenantConfiguration)
	recoverycodeProvider := recoverycode.ProvideProvider(sqlBuilder, sqlExecutor, clock, tenantConfiguration)
	provider3 := &provider2.Provider{
		Password:     passwordProvider,
		TOTP:         totpProvider,
		OOBOTP:       oobProvider,
		BearerToken:  bearertokenProvider,
		RecoveryCode: recoverycodeProvider,
	}
	userStore := &user.Store{
		SQLBuilder:  sqlBuilder,
		SQLExecutor: sqlExecutor,
	}
	welcomemessageProvider := welcomemessage.ProvideProvider(context, tenantConfiguration, engine, queue)
	rawCommands := user.ProvideRawCommands(userStore, clock, urlprefixProvider, queue, tenantConfiguration, welcomemessageProvider)
	queries := &user.Queries{
		Store:      userStore,
		Identities: providerProvider,
	}
	hookUserProvider := &auth.HookUserProvider{
		Queries:     queries,
		RawCommands: rawCommands,
	}
	hookProvider := hook.ProvideHookProvider(context, sqlBuilder, sqlExecutor, tenantConfiguration, dbContext, clock, hookUserProvider, loginidProvider, factory)
	commands := &user.Commands{
		Raw:   rawCommands,
		Hooks: hookProvider,
	}
	userProvider := &user.Provider{
		Commands: commands,
		Queries:  queries,
	}
	interactionProvider := interaction.ProvideProvider(redisStore, clock, factory, providerProvider, provider3, userProvider, oobProvider, tenantConfiguration, hookProvider)
	provider4 := challenge.ProvideProvider(context, clock, tenantConfiguration)
	anonymousFlow := &flows.AnonymousFlow{
		Enabled:      isAnonymousIdentityEnabled,
		Interactions: interactionProvider,
		Anonymous:    anonymousProvider,
		Challenges:   provider4,
	}
	idTokenIssuer := oidc.ProvideIDTokenIssuer(tenantConfiguration, urlprefixProvider, queries, clock)
	tokenGenerator := _wireTokenGeneratorValue
	tokenHandler := handler.ProvideTokenHandler(r, tenantConfiguration, factory, authorizationStore, grantStore, grantStore, grantStore, accessEventProvider, sessionProvider, anonymousFlow, idTokenIssuer, tokenGenerator, clock)
	httpHandler := provideTokenHandler(factory, dbContext, tokenHandler)
	return httpHandler
}

func newRevokeHandler(r *http.Request, m auth.DependencyMap) http.Handler {
	context := auth.ProvideContext(r)
	tenantConfiguration := auth.ProvideTenantConfig(context, m)
	factory := logging.ProvideLoggerFactory(context, tenantConfiguration)
	pool := _wirePoolValue
	dbContext := db.ProvideContextOLD(context, pool, tenantConfiguration)
	sqlBuilder := db.ProvideSQLBuilderOLD(tenantConfiguration)
	sqlExecutor := db.ProvideSQLExecutor(dbContext)
	clock := _wireSystemClockValue
	grantStore := redis.ProvideGrantStore(context, factory, tenantConfiguration, sqlBuilder, sqlExecutor, clock)
	revokeHandler := &handler.RevokeHandler{
		OfflineGrants: grantStore,
		AccessGrants:  grantStore,
	}
	httpHandler := provideRevokeHandler(factory, dbContext, revokeHandler)
	return httpHandler
}

func newMetadataHandler(r *http.Request, m auth.DependencyMap) http.Handler {
	urlprefixProvider := urlprefix.NewProvider(r)
	endpointsProvider := &auth.EndpointsProvider{
		PrefixProvider: urlprefixProvider,
	}
	metadataProvider := &oauth2.MetadataProvider{
		AuthorizeEndpoint: endpointsProvider,
		TokenEndpoint:     endpointsProvider,
		RevokeEndpoint:    endpointsProvider,
	}
	oidcMetadataProvider := &oidc.MetadataProvider{
		URLPrefix:          urlprefixProvider,
		JWKSEndpoint:       endpointsProvider,
		UserInfoEndpoint:   endpointsProvider,
		EndSessionEndpoint: endpointsProvider,
	}
	httpHandler := provideMetadataHandler(metadataProvider, oidcMetadataProvider)
	return httpHandler
}

func newJWKSHandler(r *http.Request, m auth.DependencyMap) http.Handler {
	context := auth.ProvideContext(r)
	tenantConfiguration := auth.ProvideTenantConfig(context, m)
	httpHandler := provideJWKSHandler(tenantConfiguration)
	return httpHandler
}

func newUserInfoHandler(r *http.Request, m auth.DependencyMap) http.Handler {
	context := auth.ProvideContext(r)
	tenantConfiguration := auth.ProvideTenantConfig(context, m)
	factory := logging.ProvideLoggerFactory(context, tenantConfiguration)
	pool := _wirePoolValue
	dbContext := db.ProvideContextOLD(context, pool, tenantConfiguration)
	urlprefixProvider := urlprefix.NewProvider(r)
	sqlBuilder := db.ProvideSQLBuilderOLD(tenantConfiguration)
	sqlExecutor := db.ProvideSQLExecutor(dbContext)
	store := &user.Store{
		SQLBuilder:  sqlBuilder,
		SQLExecutor: sqlExecutor,
	}
	clock := _wireSystemClockValue
	reservedNameChecker := auth.ProvideReservedNameChecker(m)
	typeCheckerFactory := loginid.ProvideTypeCheckerFactory(tenantConfiguration, reservedNameChecker)
	checker := loginid.ProvideChecker(tenantConfiguration, typeCheckerFactory)
	normalizerFactory := loginid.ProvideNormalizerFactory(tenantConfiguration)
	loginidProvider := loginid.ProvideProvider(sqlBuilder, sqlExecutor, clock, tenantConfiguration, checker, normalizerFactory)
	oauthProvider := oauth.ProvideProvider(sqlBuilder, sqlExecutor, clock)
	anonymousProvider := anonymous.ProvideProvider(sqlBuilder, sqlExecutor)
	providerProvider := provider.ProvideProvider(tenantConfiguration, loginidProvider, oauthProvider, anonymousProvider)
	queries := &user.Queries{
		Store:      store,
		Identities: providerProvider,
	}
	idTokenIssuer := oidc.ProvideIDTokenIssuer(tenantConfiguration, urlprefixProvider, queries, clock)
	httpHandler := provideUserInfoHandler(factory, dbContext, idTokenIssuer)
	return httpHandler
}

func newEndSessionHandler(r *http.Request, m auth.DependencyMap) http.Handler {
	context := auth.ProvideContext(r)
	tenantConfiguration := auth.ProvideTenantConfig(context, m)
	factory := logging.ProvideLoggerFactory(context, tenantConfiguration)
	pool := _wirePoolValue
	dbContext := db.ProvideContextOLD(context, pool, tenantConfiguration)
	urlprefixProvider := urlprefix.NewProvider(r)
	endpointsProvider := &auth.EndpointsProvider{
		PrefixProvider: urlprefixProvider,
	}
	isAnonymousIdentityEnabled := flows.ProvideIsAnonymousIdentityEnabled(tenantConfiguration)
	clock := _wireSystemClockValue
	store := redis2.ProvideStore(context, tenantConfiguration, clock)
	sqlBuilder := db.ProvideSQLBuilderOLD(tenantConfiguration)
	sqlExecutor := db.ProvideSQLExecutor(dbContext)
	reservedNameChecker := auth.ProvideReservedNameChecker(m)
	typeCheckerFactory := loginid.ProvideTypeCheckerFactory(tenantConfiguration, reservedNameChecker)
	checker := loginid.ProvideChecker(tenantConfiguration, typeCheckerFactory)
	normalizerFactory := loginid.ProvideNormalizerFactory(tenantConfiguration)
	loginidProvider := loginid.ProvideProvider(sqlBuilder, sqlExecutor, clock, tenantConfiguration, checker, normalizerFactory)
	oauthProvider := oauth.ProvideProvider(sqlBuilder, sqlExecutor, clock)
	anonymousProvider := anonymous.ProvideProvider(sqlBuilder, sqlExecutor)
	providerProvider := provider.ProvideProvider(tenantConfiguration, loginidProvider, oauthProvider, anonymousProvider)
	historyStoreImpl := password.ProvideHistoryStore(clock, sqlBuilder, sqlExecutor)
	passwordChecker := password.ProvideChecker(tenantConfiguration, historyStoreImpl)
	passwordProvider := password.ProvideProvider(sqlBuilder, sqlExecutor, clock, factory, historyStoreImpl, passwordChecker, tenantConfiguration)
	totpProvider := totp.ProvideProvider(sqlBuilder, sqlExecutor, clock, tenantConfiguration)
	engine := auth.ProvideTemplateEngine(tenantConfiguration, m)
	executor := auth.ProvideTaskExecutor(m)
	queue := async.ProvideTaskQueue(context, dbContext, tenantConfiguration, executor)
	oobProvider := oob.ProvideProvider(context, tenantConfiguration, sqlBuilder, sqlExecutor, clock, engine, urlprefixProvider, queue)
	bearertokenProvider := bearertoken.ProvideProvider(sqlBuilder, sqlExecutor, clock, tenantConfiguration)
	recoverycodeProvider := recoverycode.ProvideProvider(sqlBuilder, sqlExecutor, clock, tenantConfiguration)
	provider3 := &provider2.Provider{
		Password:     passwordProvider,
		TOTP:         totpProvider,
		OOBOTP:       oobProvider,
		BearerToken:  bearertokenProvider,
		RecoveryCode: recoverycodeProvider,
	}
	userStore := &user.Store{
		SQLBuilder:  sqlBuilder,
		SQLExecutor: sqlExecutor,
	}
	welcomemessageProvider := welcomemessage.ProvideProvider(context, tenantConfiguration, engine, queue)
	rawCommands := user.ProvideRawCommands(userStore, clock, urlprefixProvider, queue, tenantConfiguration, welcomemessageProvider)
	queries := &user.Queries{
		Store:      userStore,
		Identities: providerProvider,
	}
	hookUserProvider := &auth.HookUserProvider{
		Queries:     queries,
		RawCommands: rawCommands,
	}
	hookProvider := hook.ProvideHookProvider(context, sqlBuilder, sqlExecutor, tenantConfiguration, dbContext, clock, hookUserProvider, loginidProvider, factory)
	commands := &user.Commands{
		Raw:   rawCommands,
		Hooks: hookProvider,
	}
	userProvider := &user.Provider{
		Commands: commands,
		Queries:  queries,
	}
	interactionProvider := interaction.ProvideProvider(store, clock, factory, providerProvider, provider3, userProvider, oobProvider, tenantConfiguration, hookProvider)
	provider4 := challenge.ProvideProvider(context, clock, tenantConfiguration)
	anonymousFlow := &flows.AnonymousFlow{
		Enabled:      isAnonymousIdentityEnabled,
		Interactions: interactionProvider,
		Anonymous:    anonymousProvider,
		Challenges:   provider4,
	}
	stateStoreImpl := &webapp.StateStoreImpl{
		Context: context,
	}
	urlProvider := &webapp.URLProvider{
		Endpoints: endpointsProvider,
		Anonymous: anonymousFlow,
		States:    stateStoreImpl,
	}
	endSessionHandler := handler2.ProvideEndSessionHandler(tenantConfiguration, endpointsProvider, urlProvider, urlProvider)
	httpHandler := provideEndSessionHandler(factory, dbContext, endSessionHandler)
	return httpHandler
}

func newChallengeHandler(r *http.Request, m auth.DependencyMap) http.Handler {
	validator := auth.ProvideValidator(m)
	context := auth.ProvideContext(r)
	clock := _wireSystemClockValue
	tenantConfiguration := auth.ProvideTenantConfig(context, m)
	provider3 := challenge.ProvideProvider(context, clock, tenantConfiguration)
	challengeHandler := &ChallengeHandler{
		Validator:  validator,
		Challenges: provider3,
	}
	httpHandler := provideChallengeHandler(challengeHandler)
	return httpHandler
}

// wire.go:

func provideAuthorizeHandler(lf logging.Factory, tx db.TxContext, ah oauthAuthorizeHandler) http.Handler {
	h := &AuthorizeHandler{
		logger:       lf.NewLogger("oauth-authz-handler"),
		txContext:    tx,
		authzHandler: ah,
	}
	return h
}

func provideTokenHandler(lf logging.Factory, tx db.TxContext, th oauthTokenHandler) http.Handler {
	h := &TokenHandler{
		logger:       lf.NewLogger("oauth-token-handler"),
		txContext:    tx,
		tokenHandler: th,
	}
	return h
}

func provideRevokeHandler(lf logging.Factory, tx db.TxContext, rh oauthRevokeHandler) http.Handler {
	h := &RevokeHandler{
		logger:        lf.NewLogger("oauth-revoke-handler"),
		txContext:     tx,
		revokeHandler: rh,
	}
	return h
}

func provideMetadataHandler(oauth3 *oauth2.MetadataProvider, oidc2 *oidc.MetadataProvider) http.Handler {
	h := &MetadataHandler{
		metaProviders: []oauthMetadataProvider{oauth3, oidc2},
	}
	return h
}

func provideJWKSHandler(config2 *config.TenantConfiguration) http.Handler {
	h := &JWKSHandler{
		config: *config2.AppConfig.OIDC,
	}
	return h
}

func provideUserInfoHandler(lf logging.Factory, tx db.TxContext, uip oauthUserInfoProvider) http.Handler {
	h := &UserInfoHandler{
		logger:           lf.NewLogger("oauth-userinfo-handler"),
		txContext:        tx,
		userInfoProvider: uip,
	}
	return h
}

func provideEndSessionHandler(lf logging.Factory, tx db.TxContext, esh oidcEndSessionHandler) http.Handler {
	h := &EndSessionHandler{
		logger:            lf.NewLogger("oauth-end-session-handler"),
		txContext:         tx,
		endSessionHandler: esh,
	}
	return h
}

func provideChallengeHandler(h *ChallengeHandler) http.Handler {
	return h
}
