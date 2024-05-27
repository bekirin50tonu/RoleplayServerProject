package integrations

import (
	"services/pkg/common/parser"
	"services/pkg/common/response"

	"github.com/rabbitmq/amqp091-go"
)

type UploadImageIntegration struct {
	Disk     string `json:"disk"`
	FileName string `json:"filename"`
	Data     []byte `json:"data"`
}

func (h *StorageIntegrations) UploadImageIntegrationEventHandler(d amqp091.Delivery) (any, error) {

	req, err := parser.ParseDataWithInterface[UploadImageIntegration](d.Body)
	if err != nil {
		return nil, err
	}

	item, err := h.service_user.SaveImageToLocalStorage(req.Disk, req.FileName, req.Data)
	if err != nil {
		return nil, err
	}

	resp := response.ResponseToBrokerWithSuccessMessage("Image Upload Success.", item, nil)

	return resp, nil
}
