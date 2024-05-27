package response

import (
	"fmt"
	"math/rand"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Payload struct {
	Meta interface{} `json:"meta"`
	Data interface{} `json:"data"`
}

type ResponseWithData struct {
	Success bool    `json:"success"`
	Message string  `json:"message"`
	Payload Payload `json:"payload"`
}

type ErrorMessages struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ResponseWithError struct {
	Success bool            `json:"success"`
	Errors  []ErrorMessages `json:"errors"`
}

func ResponseWithSuccessMessage(ctx *fiber.Ctx, code int, message string, data interface{}, meta interface{}) error {
	payload := Payload{
		Data: data,
		Meta: meta,
	}
	response := ResponseWithData{
		Success: true,
		Message: message,
		Payload: payload,
	}

	return ctx.Status(code).JSON(response)
}

func ResponseToBrokerWithSuccessMessage(message string, data interface{}, meta interface{}) ResponseWithData {
	payload := Payload{
		Data: data,
		Meta: meta,
	}
	response := ResponseWithData{
		Success: true,
		Message: message,
		Payload: payload,
	}

	return response
}
func ResponseToBrokerWithErrorMessage(errs []ErrorMessages) ResponseWithError {

	response := ResponseWithError{
		Success: false,
		Errors:  errs,
	}

	return response
}

func ResponseWithErrorMessage(ctx *fiber.Ctx, err error) error {

	fiberErr, ok := err.(*fiber.Error)
	if ok {
		errorMessage := ErrorMessages{
			Code:    rand.Int(),
			Message: fiberErr.Error(),
		}
		response := ResponseWithError{
			Success: false,
			Errors:  []ErrorMessages{errorMessage},
		}
		return ctx.Status(fiberErr.Code).JSON(response)
	}
	validationErrs, ok := err.(validator.ValidationErrors)
	if ok {
		errorMessages := make([]ErrorMessages, 0)
		for _, v := range validationErrs {

			message := ErrorMessages{
				Code:    rand.Int(),
				Message: fmt.Sprintf("[%s]: '%v' | Needs to implement '%s'", v.Field(), v.Value(), v.Tag()),
			}
			errorMessages = append(errorMessages, message)
		}

		response := ResponseWithError{
			Success: false,
			Errors:  errorMessages,
		}

		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(response)
	}
	errorMessage := ErrorMessages{
		Code:    rand.Int(),
		Message: err.Error(),
	}

	response := ResponseWithError{
		Success: false,
		Errors:  []ErrorMessages{errorMessage},
	}

	return ctx.Status(fiber.StatusInternalServerError).JSON(response)
}

func ResponseWithValidationError(ctx *fiber.Ctx, errors []ErrorMessages) error {
	response := ResponseWithError{
		Success: false,
		Errors:  errors,
	}

	return ctx.Status(422).JSON(response)
}
