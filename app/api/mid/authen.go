package mid

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zhangpetergo/gin-service/app/api/auth"
	"github.com/zhangpetergo/gin-service/app/api/authclient"
	"github.com/zhangpetergo/gin-service/app/api/errs"
	"github.com/zhangpetergo/gin-service/foundation/logger"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// Authenticate 通过 auth 服务验证身份验证
func Authenticate(log *logger.Logger, client *authclient.Client) gin.HandlerFunc {

	return func(c *gin.Context) {

		if len(c.Errors) > 0 {
			// 如果前面的中间件发生错误，还需要身份验证吗，当然不
			// 直接跳过当前中间件，下一个到 处理错误中间件
			return
		}

		ctx := c.Request.Context()

		authorization := c.Request.Header.Get("Authorization")

		resp, err := client.Authenticate(ctx, authorization)
		if err != nil {
			c.Error(errs.New(errs.Unauthenticated, err))
			return
		}

		ctx = setUserID(ctx, resp.UserID)
		ctx = setClaims(ctx, resp.Claims)

		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}

}

// AuthenticateLocal validates a JWT from the `Authorization` header.
// AuthenticateLocal 通过 “Authorization” 标头验证 JWT。
// AuthenticateLocal 本地处理验证
func AuthenticateLocal(auth *auth.Auth) gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx := c.Request.Context()
		authorization := c.Request.Header.Get("Authorization")

		var err error
		parts := strings.Split(authorization, " ")

		switch parts[0] {
		case "Bearer":
			ctx, err = processJWT(ctx, auth, authorization)

		case "Basic":
			ctx, err = processBasic(ctx)
		default:
			// 没有 Authorization 标头
			c.Error(errs.New(errs.Unauthenticated, errors.New("expected authorization header format: Bearer <token>")))
			return
		}

		if err != nil {
			c.Error(err)
			return
		}

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// Bearer processes JWT authentication logic.
func Bearer(ath *auth.Auth) gin.HandlerFunc {

	return func(c *gin.Context) {
		ctx := c.Request.Context()
		authorization := c.Request.Header.Get("Authorization")

		claims, err := ath.Authenticate(ctx, authorization)
		if err != nil {
			c.Error(errs.New(errs.Unauthenticated, err))
			return
		}

		if claims.Subject == "" {
			c.Error(errs.Newf(errs.Unauthenticated, "authorize: you are not authorized for that action, no claims"))
			return
		}

		subjectID, err := uuid.Parse(claims.Subject)
		if err != nil {
			c.Error(errs.New(errs.Unauthenticated, fmt.Errorf("parsing subject: %w", err)))
			return
		}

		ctx = setUserID(ctx, subjectID)
		ctx = setClaims(ctx, claims)
	}
}

// Basic processes basic authentication logic.
func Basic() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		claims := auth.Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				Subject:   "38dc9d84-018b-4a15-b958-0b78af11c301",
				Issuer:    "service project",
				ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(8760 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			},
			Roles: []string{"ADMIN"},
		}

		subjectID, err := uuid.Parse(claims.Subject)
		if err != nil {
			c.Error(errs.Newf(errs.Unauthenticated, "parsing subject: %s", err))
			return
		}

		ctx = setUserID(ctx, subjectID)
		ctx = setClaims(ctx, claims)
	}
}

func processJWT(ctx context.Context, auth *auth.Auth, token string) (context.Context, error) {
	claims, err := auth.Authenticate(ctx, token)
	if err != nil {
		return ctx, errs.New(errs.Unauthenticated, err)
	}

	if claims.Subject == "" {
		return ctx, errs.Newf(errs.Unauthenticated, "authorize: you are not authorized for that action, no claims")
	}

	subjectID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return ctx, errs.New(errs.Unauthenticated, fmt.Errorf("parsing subject: %w", err))
	}

	ctx = setUserID(ctx, subjectID)
	ctx = setClaims(ctx, claims)

	return ctx, nil
}

func processBasic(ctx context.Context) (context.Context, error) {
	claims := auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "38dc9d84-018b-4a15-b958-0b78af11c301",
			Issuer:    "service project",
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(8760 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Roles: []string{"ADMIN"},
	}

	subjectID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return ctx, errs.Newf(errs.Unauthenticated, "parsing subject: %s", err)
	}

	ctx = setUserID(ctx, subjectID)
	ctx = setClaims(ctx, claims)

	return ctx, nil
}

func parseBasicAuth(auth string) (string, string, bool) {
	parts := strings.Split(auth, " ")
	if len(parts) != 2 || parts[0] != "Basic" {
		return "", "", false
	}

	c, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", "", false
	}

	username, password, ok := strings.Cut(string(c), ":")
	if !ok {
		return "", "", false
	}

	return username, password, true
}
