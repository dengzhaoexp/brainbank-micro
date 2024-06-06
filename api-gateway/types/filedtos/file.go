package filedtos

import "api-gateway/service/fileService"

type FileResponse struct {
	Code   uint32 `form:"code" json:"code"`
	Status string `form:"status" json:"status"`
	Info   string `form:"info" json:"info"`
}

type AvatarResponse struct {
	Code   uint32 `form:"code" json:"code"`
	Status string `form:"status" json:"status"`
	Info   string `form:"info" json:"info"`
	Avatar string `form:"avatar" json:"avatar"`
}

type UpdateAvatarRequest struct {
	AvatarName string `form:"avatar_name" json:"avatar_name" binding:"required"`
}

type NeuronResponse struct {
	Code   uint32                  `form:"code" json:"code"`
	Status string                  `form:"status" json:"status"`
	Info   string                  `form:"info" json:"info"`
	Data   *fileService.NeuronData `form:"data" json:"data"`
}

type RenameNeuronRequest struct {
	NeuronID   string `form:"neuron_id" json:"neuron_id" binding:"required"`
	NeuronName string `form:"neuron_name" json:"neuron_name" binding:"required"`
}

type ListNeuronResponse struct {
	Code   uint32                    `form:"code" json:"code"`
	Status string                    `form:"status" json:"status"`
	Info   string                    `form:"info" json:"info"`
	Data   []*fileService.NeuronData `form:"data" json:"data"`
}

type UploadDocumentRequest struct {
	Name      string `form:"name" json:"name" binding:"required"`
	Type      string `form:"types" json:"types" binding:"required"`
	Extension string `form:"extension" json:"extension" binding:"required"`
	NeuronID  string `form:"neuron_id" json:"neuron_id" binding:"required"`
}

type UploadDocumentResponse struct {
	Code   uint32                    `form:"code" json:"code"`
	Status string                    `form:"status" json:"status"`
	Info   string                    `form:"info" json:"info"`
	Data   *fileService.DocumentData `form:"data" json:"data"`
	Url    string                    `form:"url" json:"url"`
}

type DeleteDocumentRequest struct {
	DocumentID string `form:"document_id" json:"document_Id" binding:"required"`
	NeuronID   string `form:"neuron_id" json:"neuron_id" binding:"required"`
}

type DocumentResponse struct {
	Code   uint32 `form:"code" json:"code"`
	Status string `form:"status" json:"status"`
	Info   string `form:"info" json:"info"`
}

type ListDocumentRequest struct {
	NeuronID string `form:"neuron_id" json:"neuron_id" binding:"required"`
}

type ListDocumentsResponse struct {
	Code   uint32                      `form:"code" json:"code"`
	Status string                      `form:"status" json:"status"`
	Info   string                      `form:"info" json:"info"`
	Data   []*fileService.DocumentData `form:"data" json:"data"`
}

type RenameDocumentRequest struct {
	DocumentID string `form:"document_id" json:"document_Id" binding:"required"`
	NeuronID   string `form:"neuron_id" json:"neuron_id" binding:"required"`
	NewName    string `form:"new_name" json:"new_name" binding:"required"`
}

type RenameDocumentResponse struct {
	Code   uint32                    `form:"code" json:"code"`
	Status string                    `form:"status" json:"status"`
	Info   string                    `form:"info" json:"info"`
	Data   *fileService.DocumentData `form:"data" json:"data"`
}

type RecoveryDocumentRequest struct {
	DocumentID string `form:"document_id" json:"document_Id" binding:"required"`
}

type RecoveryDocumentResponse struct {
	Code   uint32                    `form:"code" json:"code"`
	Status string                    `form:"status" json:"status"`
	Info   string                    `form:"info" json:"info"`
	Data   *fileService.DocumentData `form:"data" json:"data"`
}
