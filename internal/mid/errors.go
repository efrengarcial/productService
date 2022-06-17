package mid

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"net/http"
)

type iErrBadRequest interface {
	GetTitle() string
	GetEntityName() string
}

type iErrCustomParameterized interface {
	GetTitle() string
	GetParams() map[string]string
}
type iErrRecordNotFound interface {
	error
	GetMessageWithID() string
}

type iErrConcurrencyFailure interface {
	error
	GetMessage() string
}
type iErrDuplicateRecord interface {
	error
	GetMessageDuplicate() string
}

// Middleware Error Handler in server package
func Error(logger *logrus.Logger) gin.HandlerFunc {
	return jsonAppErrorReporterT(gin.ErrorTypeAny, logger)
}

func jsonAppErrorReporterT(errType gin.ErrorType, logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		detectedErrors := c.Errors.ByType(errType)

		if len(detectedErrors) > 0 {

			err := detectedErrors[0].Err
			var (
				status     int
				title      string
				entityName string
				params     map[string]string
			)

			switch err.(type) {
			case iErrBadRequest:
				status = http.StatusBadRequest
				iError, _ := err.(iErrBadRequest)
				title = iError.GetTitle()
				entityName = iError.GetEntityName()
				c.Header("X-app-error", err.Error())
				c.Header("x-app-params", entityName)
				c.AbortWithStatusJSON(status, gin.H{
					"message":    err.Error(),
					"status":     status,
					"title":      title,
					"params":     entityName,
					"entityName": entityName,
					"errorKey":   err.Error(),
				})
			case iErrCustomParameterized:
				status = http.StatusBadRequest
				iError, _ := err.(iErrCustomParameterized)
				title = iError.GetTitle()
				params = iError.GetParams()
				c.AbortWithStatusJSON(status, gin.H{
					"message": err.Error(),
					"status":  status,
					"title":   title,
					"params":  params,
				})
				//https://medium.com/@seb.nyberg/better-validation-errors-in-go-gin-88f983564a3d
			case validator.ValidationErrors:
				status = http.StatusBadRequest
				c.AbortWithStatusJSON(status, gin.H{
					"message": err.Error(),
					"status":  status,
				})
			case iErrRecordNotFound:
				status = http.StatusNotFound
				c.AbortWithStatusJSON(status, gin.H{
					"message": err.Error(),
					"status":  status,
				})
			case iErrDuplicateRecord:
				status = http.StatusConflict
				c.AbortWithStatusJSON(status, gin.H{
					"message": "error.validationDuplicateRecord",
					"status":  status,
					"title":   "El registro ya existe en el sistema.",
				})
			case iErrConcurrencyFailure:
				status = http.StatusConflict
				c.AbortWithStatusJSON(status, gin.H{
					"message": "error.concurrencyFailure",
					"status":  status,
					"title":   "Otro usuario ha modificado esta data al mismo tiempo que tú. Tus cambios fueron rechazados.",
				})
			default:
				//https://blog.logrocket.com/error-handling-in-golang/
				//fmt.Printf("Error: %+v", err)
				logger.Error(fmt.Sprintf("Error: %+v", err))
				//fmt.Printf("with stack trace => %+v \n\n", err)
				status = http.StatusInternalServerError
				c.AbortWithStatusJSON(status, gin.H{
					"message": "error.http.500",
					"status":  status,
					"title":   "Error Interno del Servidor",
				})
			}

			return
		}
	}
}
