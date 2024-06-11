package middleware

import (
	"deals_chatting_app_backend/internal/constant"
	"deals_chatting_app_backend/internal/data"
	"bytes"
	"fmt"
	"io"
	"net/http"
    "strings"
    "context"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
    
	"github.com/google/uuid"
	gocloak "github.com/Nerzal/gocloak/v13"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var tracer = otel.Tracer(viper.GetString("appName"))
const UserIDKey = "userID"

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func Middleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		childLogger := logger
		// zap.ReplaceGlobals(logger)

		// Set CORS headers and content type
		c.Header("Access-Control-Allow-Headers", "Content-Type, RequestId, RequestTime, Authorization")
		c.Header("Content-Type", "application/json")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")

		// Handle preflight requests
		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		spanName := c.FullPath()
		if c.Request.Method == http.MethodOptions {
			spanName = fmt.Sprintf("%s [%s] %s", viper.GetString("appName"), c.Request.Method, c.Request.URL.Path)
		}

		// Start a span
		ctx, span := tracer.Start(
			c.Request.Context(),
			spanName,
			trace.WithAttributes(semconv.HTTPMethodKey.String(c.Request.Method)),
			trace.WithAttributes(semconv.HTTPURLKey.String(c.Request.URL.Path)),
		)
		defer span.End()

		sc := span.SpanContext()

		// Set required values in Gin context
		c.Set(constant.TraceID, sc.TraceID())
		childLogger = childLogger.With(
			zap.String("traceID", sc.TraceID().String()),
			zap.String("spanID", sc.SpanID().String()),
		)
		zap.ReplaceGlobals(childLogger)

		// Update request context with the new context
		c.Request = c.Request.WithContext(ctx)

		// Read request body
		requestBody, err := io.ReadAll(c.Request.Body)
		if err != nil {
			logger.Sugar().Warnf("Error reading request body: %v", err)
			// log.Println("Error reading request body:", err)
		}
		// Restore request body
		c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))

		// Read response body
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// Proceed to the next middleware/handler
		c.Next()

		span.SetAttributes(semconv.HTTPStatusCodeKey.Int(c.Writer.Status()))
		// Log request and response
		childLogger.Sugar().Infof("TraceID: %s - %d", sc.TraceID().String(), c.Writer.Status())
		childLogger.Sugar().Infof("Request body: %s", requestBody)
		childLogger.Sugar().Infof("Response body: %s", blw.body.String())

		if len(c.Errors) == 0 {
			span.SetStatus(codes.Ok, "OK!")
			return
		}
		// Handle errors, if any
		childLogger.Sugar().Infof("codess error: %s", codes.Error)
		span.SetStatus(codes.Error, c.Errors.String())
		c.JSON(-1, data.BaseErrorResponse{
			BaseResponse: data.BaseResponse{
				TxnRef:        sc.TraceID().String(),
				ProcessStatus: constant.PROCESS_STATUS_FAILURE,
			},
			Error: c.Errors.String(),
		})
		for _, ginErr := range c.Errors {
			logger.Error(ginErr.Error())
		}
	}
}


func KeycloakAuthMiddleware(client *gocloak.GoCloak, clientId, clientSecret, realm string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Authorization token is required"))
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("Invalid authorization header format"))
			return
		}

		accessToken := tokenParts[1]
		ctx := c.Request.Context()
		rptResult, err := client.RetrospectToken(ctx, accessToken, clientId, clientSecret, realm)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		if !*rptResult.Active {
			c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Authorization token is invalid"))
			return
		}

        userInfo, err := client.GetUserInfo(c, accessToken, realm)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unable to get user info"})
			return
		}
        userSub := userInfo.Sub
		if userSub == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
			return
		}

		userID, err := uuid.Parse(*userSub)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
			return
		}

        fmt.Println("isisuususu", userID)
		ctx = context.WithValue(c.Request.Context(), UserIDKey, userID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func PaginationMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Parse pagination parameters from query string
        page := c.DefaultQuery("page", "1")
        pageSize := c.DefaultQuery("limit", "10")

        // Convert parameters to integers
        pageInt, _ := strconv.Atoi(page)
        pageSizeInt, _ := strconv.Atoi(pageSize)

        // Calculate offset based on page number and page size
        offset := (pageInt - 1) * pageSizeInt

        // Add pagination information to context
        c.Set("offset", offset)
        c.Set("limit", pageSizeInt)

        // Continue processing the request
        c.Next()
    }
}

// package middleware

// import (
// 	"bytes"
// 	"io/ioutil"
// 	"log"
// 	"time"
// 	"fmt"

// 	"github.com/google/uuid"
// 	"github.com/gin-gonic/gin"
// 	"deals_chatting_app_backend/internal/data"
// 	"deals_chatting_app_backend/internal/constant"
// )

// // Middleware to generate and set trace_id and span_id
// func TraceAndSpanMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		// Generate a new UUID for the trace_id if not provided
// 		traceID := c.GetHeader("X-Trace-ID")
// 		if traceID == "" {
// 			traceID = uuid.New().String()
// 		}
// 		c.Set("trace_id", traceID)
// 		c.Writer.Header().Set("X-Trace-ID", traceID)

// 		// Generate a new UUID for the span_id
// 		spanID := uuid.New().String()
// 		c.Set("span_id", spanID)
// 		c.Writer.Header().Set("X-Span-ID", spanID)

// 		// Capture the current time for request_time
// 		requestTime := time.Now()
// 		c.Set("request_time", requestTime)
// 		c.Writer.Header().Set("X-Request-Time", requestTime.Format(time.RFC3339))

// 		// Proceed to the next middleware or handler
// 		c.Next()
// 	}
// }

// // BodyLogWriter is a custom response writer to capture the response body
// type BodyLogWriter struct {
// 	gin.ResponseWriter
// 	body *bytes.Buffer
// }

// func (w BodyLogWriter) Write(b []byte) (int, error) {
// 	w.body.Write(b)
// 	return w.ResponseWriter.Write(b)
// }

// // LoggerMiddleware logs the requests and responses along with other details
// func LoggerMiddleware() gin.HandlerFunc {
// 	fmt.Println("LoggerMiddleware initialized")
// 	return func(c *gin.Context) {
// 		fmt.Println("LoggerMiddleware started")
// 		start := time.Now()
//         fmt.Println("11111111111111111")
// 		// Capture request body
// 		var requestBody []byte
// 		if c.Request.Body != nil {
// 			requestBody, _ = ioutil.ReadAll(c.Request.Body)
// 		}
        
//         fmt.Println("222222222222222222")
// 		// Restore the request body for further use
// 		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))

// 		// Use the custom response writer
// 		blw := &BodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
// 		c.Writer = blw

// 		// Process the request
// 		c.Next()
        
//         fmt.Println("3333333333333333333333")
// 		// Capture response body
// 		responseBody := blw.body.String()
// 		latency := time.Since(start)
// 		status := c.Writer.Status()

// 		// Retrieve trace_id and span_id from context
// 		traceID, _ := c.Get("trace_id")
// 		spanID, _ := c.Get("span_id")

//         fmt.Println("traceIDtraceID", traceID)
        
//         fmt.Println("spanIDspanID", spanID)

// 		// Log the details
// 		log.Printf("status=%d method=%s path=%s trace_id=%v span_id=%v latency=%v request=%s response=%s",
// 			status, c.Request.Method, c.Request.URL.Path, traceID, spanID, latency, string(requestBody), responseBody)

// 		// Check for errors in Gin context
//         fmt.Println("errorororor", c.Errors)
// 		if len(c.Errors) > 0 {
// 			// Log errors
// 			for _, err := range c.Errors {
// 				log.Printf("Error: %s", err.Error())
// 			}

// 			// Set response status
// 			c.JSON(-1, data.BaseErrorResponse{
// 				BaseResponse: data.BaseResponse{
// 					TxnRef:        traceID.(string),
// 					ProcessStatus: constant.PROCESS_STATUS_FAILURE,
// 				},
// 				Error: c.Errors.String(),
// 			})
// 		}
// 	}
// }
