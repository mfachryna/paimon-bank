package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/mfachryna/paimon-bank/internal/common/response"
	"github.com/mfachryna/paimon-bank/internal/common/utils/validation"
	dto "github.com/mfachryna/paimon-bank/internal/domain/dto/balance"
	"github.com/mfachryna/paimon-bank/internal/entity"
	"go.uber.org/zap"
)

func (bh *BalanceHandler) Store(w http.ResponseWriter, r *http.Request) {
	var (
		data   dto.BalanceAdd
		userId string
	)

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		bh.log.Info("required fields are missing or invalid", zap.Error(err))
		(&response.Response{
			HttpStatus: http.StatusBadRequest,
			Message:    "required fields are missing or invalid",
		}).GenerateResponse(w)
		return
	}

	if err := bh.val.Struct(data); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, e := range validationErrors {
			bh.log.Info(validation.CustomError(e), zap.Error(err))
			(&response.Response{
				HttpStatus: http.StatusBadRequest,
				Message:    validation.CustomError(e),
			}).GenerateResponse(w)
			return
		}
	}

	if err := validation.UrlValidation(data.Receipt); err != nil {
		bh.log.Info("failed to validate url", zap.Error(err))
		(&response.Response{
			HttpStatus: http.StatusBadRequest,
			Message:    "URL malformed",
		}).GenerateResponse(w)

		return
	}

	ctx := r.Context()
	userId = ctx.Value("user_id").(string)
	dataId := uuid.NewString()
	balanceHistory := entity.BalanceHistory{
		ID:       dataId,
		Balance:  data.Balance,
		Currency: data.Currency,
		Receipt:  data.Receipt,
		UserId:   userId,
		Source: entity.BalanceSource{
			BankName:   data.BankName,
			BankNumber: data.BankNumber,
		},
	}

	if err := bh.br.Insert(ctx, balanceHistory); err != nil {
		bh.log.Info("failed to insert data", zap.Error(err))
		(&response.Response{
			HttpStatus: http.StatusInternalServerError,
			Message:    err.Error(),
		}).GenerateResponse(w)
		return
	}

	(&response.Response{
		HttpStatus: http.StatusOK,
		Message:    "Add balance success",
	}).GenerateResponse(w)

}
