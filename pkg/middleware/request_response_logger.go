package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"golang-rest-api-template/pkg/auth"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type LogRequestModel struct {
	T                             string            `bson:"_t" json:"_t"`
	RequestDateTime               time.Time         `bson:"RequestDateTime" json:"RequestDateTime"`
	RequestDateTimeUtc            time.Time         `bson:"RequestDateTimeUtc" json:"RequestDateTimeUtc"`
	RequestDateTimeUtcActionLevel time.Time         `bson:"RequestDateTimeUtcActionLevel" json:"RequestDateTimeUtcActionLevel"`
	RequestMethod                 string            `bson:"RequestMethod" json:"RequestMethod"`
	RequestPath                   string            `bson:"RequestPath" json:"RequestPath"`
	RequestHeaders                map[string]string `bson:"RequestHeaders" json:"RequestHeaders"`
	RequestBody                   interface{}       `bson:"RequestBody" json:"RequestBody"`
	RequestContentType            string            `bson:"RequestContentType" json:"RequestContentType"`
	RequestQueryParams            map[string]string `bson:"RequestQueryParams" json:"RequestQueryParams"`
}

type LogResponseModel struct {
	T                              string            `bson:"_t" json:"_t"`
	ResponseDateTime               time.Time         `bson:"ResponseDateTime" json:"ResponseDateTime"`
	ResponseDateTimeUtc            time.Time         `bson:"ResponseDateTimeUtc" json:"ResponseDateTimeUtc"`
	ResponseDateTimeUtcActionLevel time.Time         `bson:"ResponseDateTimeUtcActionLevel" json:"ResponseDateTimeUtcActionLevel"`
	ResponseStatus                 string            `bson:"ResponseStatus" json:"ResponseStatus"`
	ResponseHeaders                map[string]string `bson:"ResponseHeaders" json:"ResponseHeaders"`
	ResponseBody                   interface{}       `bson:"ResponseBody" json:"ResponseBody"`
	ResponseContentType            string            `bson:"ResponseContentType" json:"ResponseContentType"`
}

type RequestResponseLogModel struct {
	T                           string            `bson:"_t" json:"_t"`
	Type                        string            `bson:"Type" json:"Type"`
	ModelType                   string            `bson:"ModelType" json:"ModelType"`
	LogId                       string            `bson:"LogId" json:"LogId"`
	Node                        string            `bson:"Node" json:"Node"`
	LocalIpAddress              string            `bson:"LocalIpAddress" json:"LocalIpAddress"`
	RemoteIpAddress             string            `bson:"RemoteIpAddress" json:"RemoteIpAddress"`
	TraceId                     string            `bson:"TraceId" json:"TraceId"`
	CurrentUserId               string            `bson:"CurrentUserId" json:"CurrentUserId"`
	CurrentUsername             string            `bson:"CurrentUsername" json:"CurrentUsername"`
	Environment                 string            `bson:"Environment" json:"Environment"`
	OsVersion                   string            `bson:"OsVersion" json:"OsVersion"`
	BeginDateTimeUtcActionLevel time.Time         `bson:"BeginDateTimeUtcActionLevel" json:"BeginDateTimeUtcActionLevel"`
	EndDateTimeUtcActionLevel   time.Time         `bson:"EndDateTimeUtcActionLevel" json:"EndDateTimeUtcActionLevel"`
	RequestModel                *LogRequestModel  `bson:"RequestModel,omitempty" json:"RequestModel,omitempty"`
	ReferenceModel              *LogResponseModel `bson:"ReferenceModel,omitempty" json:"ReferenceModel,omitempty"`
	IsExceptionActionLevel      *bool             `bson:"IsExceptionActionLevel" json:"IsExceptionActionLevel"`
}

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func extractUserFromToken(c *gin.Context) (userId, username string) {
	const BearerSchema = "Bearer "
	header := c.GetHeader("Authorization")
	if header == "" || !strings.HasPrefix(header, BearerSchema) {
		return "", ""
	}

	tokenStr := header[len(BearerSchema):]
	claims := &auth.CustomClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return auth.JwtKey, nil
	})

	if err != nil || !token.Valid {
		return "", ""
	}

	return claims.UserId, claims.Username
}

func getLocalIP() string {
	return "::ffff:127.0.0.1"
}

func getEnvironment() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}

func getOsVersion() string {
	return runtime.GOOS
}
func RequestResponseLogger(logger *zap.Logger, collection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now().UTC()
		logId := uuid.New().String()
		traceId := c.GetHeader("X-Trace-Id")
		if traceId == "" {
			traceId = uuid.New().String()
		}

		userId, username := extractUserFromToken(c)

		c.Set("logId", logId)
		c.Set("traceId", traceId)

		var requestBodyBytes []byte
		var requestBodyParsed interface{}
		if c.Request.Body != nil {
			requestBodyBytes, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBodyBytes))

			if len(requestBodyBytes) > 0 {
				if err := json.Unmarshal(requestBodyBytes, &requestBodyParsed); err != nil {
					requestBodyParsed = string(requestBodyBytes)
				}
			}
		}

		requestHeaders := make(map[string]string)
		for key, values := range c.Request.Header {
			if len(values) > 0 {
				if strings.ToLower(key) == "authorization" {
					requestHeaders[key] = "***"
				} else {
					requestHeaders[key] = values[0]
				}
			}
		}

		queryParams := make(map[string]string)
		for key, values := range c.Request.URL.Query() {
			if len(values) > 0 {
				queryParams[key] = values[0]
			}
		}

		requestLogEntry := RequestResponseLogModel{
			T:                           "RequestResponseLogModel",
			Type:                        "Api",
			ModelType:                   "Request",
			LogId:                       logId,
			Node:                        getEnvironment(),
			LocalIpAddress:              getLocalIP(),
			RemoteIpAddress:             "::ffff:" + c.ClientIP(),
			TraceId:                     traceId,
			CurrentUserId:               userId,
			CurrentUsername:             username,
			Environment:                 getEnvironment(),
			OsVersion:                   getOsVersion(),
			BeginDateTimeUtcActionLevel: startTime,
			EndDateTimeUtcActionLevel:   startTime,
			RequestModel: &LogRequestModel{
				T:                             "LogRequestModel",
				RequestDateTime:               startTime,
				RequestDateTimeUtc:            startTime,
				RequestDateTimeUtcActionLevel: startTime,
				RequestMethod:                 c.Request.Method,
				RequestPath:                   c.Request.URL.Path,
				RequestHeaders:                requestHeaders,
				RequestBody:                   requestBodyParsed,
				RequestContentType:            c.Request.Header.Get("Content-Type"),
				RequestQueryParams:            queryParams,
			},
			IsExceptionActionLevel: nil,
		}

		if _, err := collection.InsertOne(context.TODO(), requestLogEntry); err != nil {
			logger.Error("Failed to log request to MongoDB", zap.Error(err))
		}
		responseBody := &bytes.Buffer{}
		writer := responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           responseBody,
		}
		c.Writer = writer

		c.Next()

		endTime := time.Now().UTC()
		responseTime := time.Now().UTC()

		responseHeaders := make(map[string]string)
		for key, values := range c.Writer.Header() {
			if len(values) > 0 {
				responseHeaders[key] = values[0]
			}
		}

		var responseBodyParsed interface{}
		if err := json.Unmarshal(responseBody.Bytes(), &responseBodyParsed); err != nil {
			responseBodyParsed = responseBody.String()
		}

		responseLogEntry := RequestResponseLogModel{
			T:                           "RequestResponseLogModel",
			Type:                        "Api",
			ModelType:                   "Response",
			LogId:                       logId,
			Node:                        getEnvironment(),
			LocalIpAddress:              getLocalIP(),
			RemoteIpAddress:             "::ffff:" + c.ClientIP(),
			TraceId:                     traceId,
			CurrentUserId:               userId,
			CurrentUsername:             username,
			Environment:                 getEnvironment(),
			OsVersion:                   getOsVersion(),
			BeginDateTimeUtcActionLevel: startTime,
			EndDateTimeUtcActionLevel:   endTime,
			ReferenceModel: &LogResponseModel{
				T:                              "LogResponseModel",
				ResponseDateTime:               responseTime,
				ResponseDateTimeUtc:            responseTime,
				ResponseDateTimeUtcActionLevel: responseTime,
				ResponseStatus:                 strconv.Itoa(c.Writer.Status()),
				ResponseHeaders:                responseHeaders,
				ResponseBody:                   responseBodyParsed,
				ResponseContentType:            c.Writer.Header().Get("Content-Type"),
			},
			IsExceptionActionLevel: nil,
		}

		if _, err := collection.InsertOne(context.TODO(), responseLogEntry); err != nil {
			logger.Error("Failed to log response to MongoDB", zap.Error(err))
		}
	}
}
