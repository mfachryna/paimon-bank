package handler

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
	"github.com/mfachryna/paimon-bank/internal/common/response"
	"github.com/mfachryna/paimon-bank/internal/common/utils/validation"
	balancedto "github.com/mfachryna/paimon-bank/internal/domain/dto/balance"
	metadto "github.com/mfachryna/paimon-bank/internal/domain/dto/meta"
	"go.uber.org/zap"
)

func (bh *BalanceHandler) History(w http.ResponseWriter, r *http.Request) {
	var (
		filter balancedto.BalanceFilter
	)

	if err := r.ParseForm(); err != nil {
		bh.log.Info("failed to parse form", zap.Error(err))
		(&response.Response{
			HttpStatus: http.StatusInternalServerError,
			Message:    err.Error(),
		}).GenerateResponse(w)
		return
	}

	if err := validation.ValidateParams(r, filter); err != nil {
		bh.log.Info("failed to validate params", zap.Error(err))
		(&response.Response{
			HttpStatus: http.StatusBadRequest,
			Message:    err.Error(),
		}).GenerateResponse(w)
		return
	}

	if err := schema.NewDecoder().Decode(&filter, r.Form); err != nil {
		bh.log.Info("required fields are missing or invalid", zap.Error(err))
		(&response.Response{
			HttpStatus: http.StatusBadRequest,
			Message:    err.Error(),
		}).GenerateResponse(w)
		return
	}

	if err := bh.val.Struct(filter); err != nil {
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

	ctx := r.Context()
	userId := ctx.Value("user_id").(string)

	if filter.Limit == 0 {
		filter.Limit = 5
	}
	filter.Offset = filter.Limit * filter.Offset

	data, count, err := bh.br.GetBalanceHistory(ctx, filter, userId)
	if err != nil {
		bh.log.Info("failed to get balance history with filter", zap.Error(err))
		(&response.Response{
			HttpStatus: http.StatusInternalServerError,
			Message:    err.Error(),
		}).GenerateResponse(w)
		return
	}

	(&response.ResponseWithMeta{
		HttpStatus: http.StatusOK,
		Data:       data,
		Meta: metadto.Meta{
			Limit:  filter.Limit,
			Offset: filter.Offset,
			Total:  count,
		},
	}).GenerateResponseMeta(w)

}
