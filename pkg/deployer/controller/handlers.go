package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/crossedbot/common/golang/logger"
	"github.com/crossedbot/common/golang/server"

	"github.com/crossedbot/matryoshka/pkg/deployer/models"
	"github.com/crossedbot/matryoshka/pkg/runner"
)

// CreateContainer handles requests to deploy and run code for a given language.
func CreateContainer(w http.ResponseWriter, r *http.Request, p server.Parameters) {
	var payload runner.Payload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		logger.Error(err)
		server.JsonResponse(
			w,
			models.Error{
				Code: models.ErrFailedConversionCode,
				Message: fmt.Sprintf(
					"failed to parse request body; %s",
					err,
				),
			},
			http.StatusBadRequest,
		)
		return
	}
	result, err := V1().CreateDeployment(payload)
	if err != nil {
		logger.Error(err)
		server.JsonResponse(
			w,
			models.Error{
				Code: models.ErrProcessingRequestCode,
				Message: fmt.Sprintf(
					"failed to deploy container; %s",
					err,
				),
			},
			http.StatusInternalServerError,
		)
		return
	}
	server.JsonResponse(w, &result, http.StatusOK)
}
