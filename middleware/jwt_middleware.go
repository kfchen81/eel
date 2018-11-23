package middleware

import (
	go_context "context"
	"github.com/kfchen81/eel/config"
	"github.com/kfchen81/eel/handler"
	"github.com/kfchen81/eel/log"
	"strings"
	"fmt"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"github.com/bitly/go-simplejson"
	"github.com/opentracing/opentracing-go"
)

var SALT string = "030e2cf548cf9da683e340371d1a74ee"
var SKIP_JWT_CHECK_URLS []string = make([]string, 0)

type JWTMiddleware struct {
	handler.Middleware
}

func (this *JWTMiddleware) ProcessRequest(ctx *handler.Context) {
	uri := ctx.Request.HttpRequest.RequestURI
	for _, skipUrl := range SKIP_JWT_CHECK_URLS {
		if strings.Contains(uri, skipUrl) {
			log.Logger.Debug("[jwt_middleware] skip jwt check", "url", skipUrl)
			return
		}
	}
	
	//get jwt token
	jwtToken := ctx.Request.Header("AUTHORIZATION");
	if jwtToken == "" {
		//for dev
		jwtToken = ctx.Request.Query("_jwt")
	}
	
	if jwtToken != "" {
		items := strings.Split(jwtToken, ".")
		if len(items) != 3 {
			//jwt token 格式不对
			response := handler.MakeErrorResponse(500, "jwt:invalid_jwt_token", "无效的jwt token 1")
			ctx.Response.JSON(response)
			return
		}
		
		headerB64Code, payloadB64Code, expectedSignature := items[0], items[1], items[2]
		message := fmt.Sprintf("%s.%s", headerB64Code, payloadB64Code)
		
		h := hmac.New(sha256.New, []byte(SALT))
		h.Write([]byte(message))
		actualSignature := base64.StdEncoding.EncodeToString(h.Sum(nil));
		
		if expectedSignature != actualSignature {
			//jwt token的signature不匹配
			response := handler.MakeErrorResponse(500, "jwt:invalid_jwt_token", "无效的jwt token 2")
			ctx.Response.JSON(response)
			return
		}
		
		decodeBytes, err := base64.StdEncoding.DecodeString(payloadB64Code)
		if err != nil {
			log.Logger.Fatal(err)
		}
		js, err := simplejson.NewJson([]byte(decodeBytes))
		
		if err != nil {
			response := handler.MakeErrorResponse(500, "jwt:invalid_jwt_token", "无效的jwt token 3")
			ctx.Response.JSON(response)
			return
		}
		
		userId, err := js.Get("uid").Int()
		if err != nil {
			log.Logger.Fatal(err)
			response := handler.MakeErrorResponse(500, "jwt:invalid_jwt_token", "无效的jwt token 4")
			ctx.Response.JSON(response)
			return
		}
		
		var bCtx go_context.Context
		if config.Runtime.NewBusinessContext != nil {
			bCtx = config.Runtime.NewBusinessContext(go_context.Background(), ctx.Request.HttpRequest, userId, jwtToken, js) //bCtx is for "business context"
		}
		ctx.SetBusinessContext(bCtx)
		ctx.Set("span", opentracing.SpanFromContext(bCtx))
	} else {
		response := handler.MakeErrorResponse(500, "jwt:invalid_jwt_token", "无效的jwt token 5")
		ctx.Response.JSON(response)
		return
	}
	log.Logger.Info("i am in jwt middleware process request")
}

func (this *JWTMiddleware) ProcessResponse(ctx *handler.Context) {
	log.Logger.Info("i am in jwt middleware process response")
}

func init() {
	skipUrls := config.ServiceConfig.String("SKIP_JWT_CHECK_URLS")
	if skipUrls == "" {
		log.Logger.Info("SKIP_JWT_CHECK_URLS is empty")
	} else {
		SKIP_JWT_CHECK_URLS = strings.Split(skipUrls, ";")
	}
	
	log.Logger.Infow("[jwt_middleware]", "SKIP_JWT_CHECK_URLS", SKIP_JWT_CHECK_URLS)
}

