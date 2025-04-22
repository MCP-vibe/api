package action

import (
	"encoding/json"
	"net/http"

	"api/internal/adapters/api/logging"
	"api/internal/adapters/api/response"
	"api/internal/adapters/logger"
	"api/internal/adapters/validator"
	"api/internal/usecase"
)

type CreateUserAction struct {
	uc        usecase.CreateUserUseCase
	log       logger.Logger
	validator validator.Validator
}

func NewCreateUserAction(uc usecase.CreateUserUseCase, log logger.Logger, v validator.Validator) CreateUserAction {
	return CreateUserAction{
		uc:        uc,
		log:       log,
		validator: v,
	}
}

func (a CreateUserAction) Execute(w http.ResponseWriter, r *http.Request) {
	const logKey = "create_user"
	var input usecase.CreateUserInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logging.NewError(
			a.log,
			err,
			logKey,
			http.StatusBadRequest,
		).Log("error when decoding json")

		response.NewError(err, http.StatusBadRequest).Send(w)
		return
	}
	defer r.Body.Close()

	if errs := a.validateInput(input); len(errs) > 0 {
		logging.NewError(
			a.log,
			response.ErrInvalidInput,
			logKey,
			http.StatusBadRequest,
		).Log("invalid input")

		response.NewErrorMessage(errs, http.StatusBadRequest).Send(w)
		return
	}

	err := a.uc.Execute(r.Context(), input)
	if err != nil {
		response.NewErrorWithErrorStatus(err, w, a.log, logKey, "error when create_user")
		return
	}

	logging.NewInfo(a.log, logKey, http.StatusOK).Log("success create_user")

	response.NewSuccess(nil, http.StatusOK).Send(w)
}

func (c CreateUserAction) validateInput(input usecase.CreateUserInput) []string {
	var msgs []string

	err := c.validator.Validate(input)
	if err != nil {
		msgs = append(msgs, c.validator.Messages()...)
	}

	return msgs
}
