package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/mfachryna/paimon-bank/internal/common/response"
	"github.com/mfachryna/paimon-bank/internal/common/utils/validation"
	dto "github.com/mfachryna/paimon-bank/internal/domain/dto/user"
	"github.com/mfachryna/paimon-bank/internal/entity"
	"github.com/mfachryna/paimon-bank/pkg/jwt"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func (uh *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var (
		data    dto.UserCreate
		resData interface{}
	)

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		uh.log.Info("required fields are missing or invalid", zap.Error(err))
		(&response.Response{
			HttpStatus: http.StatusBadRequest,
			Message:    "required fields are missing or invalid",
		}).GenerateResponse(w)
		return
	}

	if err := uh.val.Struct(data); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, e := range validationErrors {
			uh.log.Info(validation.CustomError(e), zap.Error(err))
			(&response.Response{
				HttpStatus: http.StatusBadRequest,
				Message:    validation.CustomError(e),
			}).GenerateResponse(w)
			return
		}
	}

	ctx := r.Context()
	uuid := uuid.NewString()
	user := entity.User{
		ID:   uuid,
		Name: data.Name,
	}

	tokenString, err := jwt.SignedToken(jwt.Claim{UserId: uuid})
	if err != nil {
		uh.log.Info("failed to sign token", zap.Error(err))
		(&response.Response{
			HttpStatus: http.StatusInternalServerError,
			Message:    err.Error(),
		}).GenerateResponse(w)
		return
	}

	if err := validation.EmailValidation(data.Email); err != nil {
		uh.log.Info("failed to validate email credential", zap.Error(err))
		(&response.Response{
			HttpStatus: http.StatusBadRequest,
			Message:    err.Error(),
		}).GenerateResponse(w)
		return
	}

	count, err := uh.ur.EmailCheck(ctx, data.Email)
	if err != nil {
		uh.log.Info("failed to get email", zap.Error(err))
		(&response.Response{
			HttpStatus: http.StatusInternalServerError,
			Message:    err.Error(),
		}).GenerateResponse(w)
		return
	}

	if count > 0 {
		uh.log.Info("email is already used")
		(&response.Response{
			HttpStatus: http.StatusConflict,
			Message:    "email is already used",
		}).GenerateResponse(w)
		return
	}

	user.Email = data.Email
	resData = dto.UserCredential{
		Email:       user.Email,
		Name:        user.Name,
		AccessToken: tokenString,
	}

	salt, err := strconv.Atoi(uh.cfg.App.BcryptSalt)
	if err != nil {
		uh.log.Info("failed to convert string salt", zap.Error(err))
		(&response.Response{
			HttpStatus: http.StatusInternalServerError,
			Message:    err.Error(),
		}).GenerateResponse(w)
		return
	}

	hashedPasswordChan := make(chan string)
	go func() {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), salt)
		if err != nil {
			uh.log.Info("failed to hash password", zap.Error(err))
			(&response.Response{
				HttpStatus: http.StatusConflict,
				Message:    "email is already used",
			}).GenerateResponse(w)
			return
		}
		hashedPasswordChan <- string(hashedPassword)
	}()

	user.Password = <-hashedPasswordChan

	if err := uh.ur.Insert(ctx, user); err != nil {
		uh.log.Info("failed to insert data", zap.Error(err))
		(&response.Response{
			HttpStatus: http.StatusInternalServerError,
			Message:    err.Error(),
		}).GenerateResponse(w)
		return
	}

	(&response.Response{
		HttpStatus: http.StatusCreated,
		Message:    "User registered successfully",
		Data:       resData,
	}).GenerateResponse(w)
}
