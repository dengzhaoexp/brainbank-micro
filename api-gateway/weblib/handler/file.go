package handler

import (
	"api-gateway/pkg/respmsg"
	"api-gateway/service/fileService"
	"api-gateway/types/filedtos"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAvatar() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取Get请求的参数
		userId := ctx.Param("userId")

		// 绑定参数
		var req fileService.GetAvatarRequest
		req.UserId = userId

		// 从gin中获取服务实例
		service := ctx.Keys["rpcFileService"].(fileService.FileService)
		rpcResp, err := service.GetAvatar(ctx, &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
			return
		}

		// 返回数据
		resp := filedtos.AvatarResponse{
			Code:   rpcResp.Code,
			Status: "请求成功",
			Info:   respmsg.GetMsg(rpcResp.Code),
			Avatar: rpcResp.AvatarUrl,
		}
		ctx.JSON(http.StatusOK, resp)
	}
}

func UpdateAvatar() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 绑定参数
		var req filedtos.UpdateAvatarRequest
		if err := bindAndValidate(ctx, &req); err != nil {
			returnError(ctx, http.StatusBadRequest, err)
			return
		}

		// 绑定参数
		rpcReq := &fileService.UpdateAvatarRequest{
			UserId:     ctx.Keys["userId"].(string),
			AvatarName: req.AvatarName,
		}

		// 从gin中获取服务实例
		service, ok := ctx.Keys["rpcFileService"].(fileService.FileService)
		if !ok {
			returnError(ctx, http.StatusInternalServerError, errors.New("rpcChatService not found in context"))
			return
		}

		// 通过实例调用具体服务
		rpcResp, err := service.UpdateAvatar(ctx, rpcReq)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
			return
		}

		// 返回数据
		resp := filedtos.AvatarResponse{
			Code:   rpcResp.Code,
			Status: "请求成功",
			Info:   respmsg.GetMsg(rpcResp.Code),
			Avatar: rpcResp.AvatarUrl,
		}
		ctx.JSON(http.StatusOK, resp)
	}

}

func CreateNeuron() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 绑定参数
		rpcReq := &fileService.CreateNeuronRequest{
			UserId: ctx.Keys["userId"].(string),
		}

		// 获取服务实例
		service := ctx.Keys["rpcFileService"].(fileService.FileService)
		rpcResp, err := service.CreateNeuron(ctx, rpcReq)
		if err != nil {
			returnError(ctx, http.StatusInternalServerError, err)
			return
		}

		// 绑定返回数据
		resp := filedtos.NeuronResponse{
			Code:   rpcResp.Code,
			Status: "请求成功",
			Info:   respmsg.GetMsg(rpcResp.Code),
			Data:   rpcResp.Data,
		}
		ctx.JSON(http.StatusOK, resp)
	}
}

func DeleteNeuron() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取参数
		id := ctx.Param("neuronId")

		// 绑定参数
		rpcReq := &fileService.DeleteNeuronRequest{
			NeuronId: id,
			UserId:   ctx.Keys["userId"].(string),
		}

		// 获取服务实例
		service := ctx.Keys["rpcFileService"].(fileService.FileService)
		rpcResp, err := service.DeleteNeuron(ctx, rpcReq)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
			return
		}

		// 绑定返回数据
		resp := filedtos.FileResponse{
			Code:   rpcResp.Code,
			Status: "请求成功",
			Info:   respmsg.GetMsg(rpcResp.Code),
		}
		ctx.JSON(http.StatusOK, resp)

	}
}

func RenameNeuron() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取参数
		var req filedtos.RenameNeuronRequest
		if err := bindAndValidate(ctx, &req); err != nil {
			returnError(ctx, http.StatusBadRequest, err)
			return
		}

		// 绑定参数
		userId := ctx.Keys["userId"].(string)
		rpcReq := &fileService.RenameNeuronRequest{
			UserId:     userId,
			NeuronId:   req.NeuronID,
			NeuronName: req.NeuronName,
		}

		// 获取服务实例
		service := ctx.Keys["rpcFileService"].(fileService.FileService)
		rpcResp, err := service.RenameNeuron(ctx, rpcReq)
		if err != nil {
			returnError(ctx, http.StatusInternalServerError, err)
			return
		}

		// 绑定返回数据
		resp := filedtos.NeuronResponse{
			Code:   rpcResp.Code,
			Status: "请求成功",
			Info:   respmsg.GetMsg(rpcResp.Code),
			Data:   rpcResp.Data,
		}
		ctx.JSON(http.StatusOK, resp)
	}
}

func ListNeuron() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 绑定请求参数
		rpcReq := &fileService.ListNeuronRequest{UserId: ctx.Keys["userId"].(string)}

		// 调用服务
		service := ctx.Keys["rpcFileService"].(fileService.FileService)
		rpcResp, err := service.ListNeuron(ctx, rpcReq)
		if err != nil {
			returnError(ctx, http.StatusInternalServerError, err)
			return
		}

		// 绑定返回数据
		resp := filedtos.ListNeuronResponse{
			Code:   rpcResp.Code,
			Status: "请求成功",
			Info:   respmsg.GetMsg(rpcResp.Code),
			Data:   rpcResp.Data,
		}
		ctx.JSON(http.StatusOK, resp)

	}
}

func UploadDocument() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取数据
		var req filedtos.UploadDocumentRequest
		if err := bindAndValidate(ctx, &req); err != nil {
			returnError(ctx, http.StatusBadRequest, err)
			return
		}

		// 绑定数据
		rpcReq := &fileService.UploadDocumentRequest{
			Name:      req.Name,
			Type:      req.Type,
			Extension: req.Extension,
			UserId:    ctx.Keys["userId"].(string),
			NeuronId:  req.NeuronID,
		}

		// 从gin中获取服务
		service := ctx.Keys["rpcFileService"].(fileService.FileService)
		rpcResp, err := service.UploadDocument(ctx, rpcReq)
		if err != nil {
			returnError(ctx, http.StatusInternalServerError, err)
			return
		}

		// 绑定返回数据
		resp := &filedtos.UploadDocumentResponse{
			Code:   rpcResp.Code,
			Status: "请求成功",
			Info:   respmsg.GetMsg(rpcResp.Code),
			Data:   rpcResp.Data,
			Url:    rpcResp.FileUrl,
		}
		ctx.JSON(http.StatusOK, resp)
	}
}

func DeleteDocument() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取请求数据
		var req filedtos.DeleteDocumentRequest
		if err := bindAndValidate(ctx, &req); err != nil {
			returnError(ctx, http.StatusBadRequest, err)
			return
		}

		// 绑定数据
		rpcReq := &fileService.DeleteDocumentRequest{
			DocumentId: req.DocumentID,
			NeuronId:   req.NeuronID,
			UserId:     ctx.Keys["userId"].(string),
		}

		// 获取服务
		service := ctx.Keys["rpcFileService"].(fileService.FileService)
		rpcResp, err := service.DeleteDocument(ctx, rpcReq)
		if err != nil {
			returnError(ctx, http.StatusInternalServerError, err)
			return
		}

		// 返回数据
		resp := filedtos.DocumentResponse{
			Code:   rpcResp.Code,
			Status: "请求成功",
			Info:   respmsg.GetMsg(rpcResp.Code),
		}
		ctx.JSON(http.StatusOK, resp)
	}
}

func ListDocuments() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取请求数据
		var req filedtos.ListDocumentRequest
		if err := bindAndValidate(ctx, &req); err != nil {
			returnError(ctx, http.StatusBadRequest, err)
			return
		}

		// 绑定数据
		rpcReq := &fileService.ListDocumentsRequest{
			UserId:   ctx.Keys["userId"].(string),
			NeuronId: req.NeuronID,
		}

		// 获取请求
		service := ctx.Keys["rpcFileService"].(fileService.FileService)
		rpcResp, err := service.ListDocuments(ctx, rpcReq)
		if err != nil {
			returnError(ctx, http.StatusInternalServerError, err)
			return
		}

		// 处理返回
		resp := filedtos.ListDocumentsResponse{
			Code:   rpcResp.Code,
			Status: "请求成功",
			Info:   respmsg.GetMsg(rpcResp.Code),
			Data:   rpcResp.Data,
		}
		ctx.JSON(http.StatusOK, resp)
	}
}

func RenameDocument() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取请求参数
		var req filedtos.RenameDocumentRequest
		if err := bindAndValidate(ctx, &req); err != nil {
			returnError(ctx, http.StatusBadRequest, err)
			return
		}

		// 绑定请求参数
		rpcReq := &fileService.RenameDocumentRequest{
			DocumentId: req.DocumentID,
			NeuronId:   req.NeuronID,
			NewName:    req.NewName,
		}

		// 获取请求
		service := ctx.Keys["rpcFileService"].(fileService.FileService)
		rpcResp, err := service.RenameDocument(ctx, rpcReq)
		if err != nil {
			returnError(ctx, http.StatusInternalServerError, err)
			return
		}

		// 绑定返回数据
		resp := &filedtos.RenameDocumentResponse{
			Code:   rpcResp.Code,
			Status: "请求成功",
			Info:   respmsg.GetMsg(rpcResp.Code),
			Data:   rpcResp.Data,
		}
		ctx.JSON(http.StatusOK, resp)
	}
}

func ListDocumentInBin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 绑定请求参数
		rpcReq := &fileService.ListDocumentsInBinRequest{UserId: ctx.Keys["userId"].(string)}

		// 获取服务
		service := ctx.Keys["rpcFileService"].(fileService.FileService)
		rpcResp, err := service.ListDocumentsInBin(ctx, rpcReq)
		if err != nil {
			returnError(ctx, http.StatusInternalServerError, err)
			return
		}

		// 绑定返回数据
		resp := &filedtos.ListDocumentsResponse{
			Code:   rpcResp.Code,
			Status: "请求成功",
			Info:   respmsg.GetMsg(rpcResp.Code),
			Data:   rpcResp.Data,
		}
		ctx.JSON(http.StatusOK, resp)
	}
}

func RecoveryDocument() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取请求参数
		var req filedtos.RecoveryDocumentRequest
		if err := bindAndValidate(ctx, &req); err != nil {
			returnError(ctx, http.StatusBadRequest, err)
			return
		}

		// 绑定参数
		rpcReq := &fileService.RecoveryDocumentRequest{
			UserId:     ctx.Keys["userId"].(string),
			DocumentId: req.DocumentID,
		}

		// 获取服务
		service := ctx.Keys["rpcFileService"].(fileService.FileService)
		rpcResp, err := service.RecoveryDocument(ctx, rpcReq)
		if err != nil {
			returnError(ctx, http.StatusInternalServerError, err)
			return
		}

		// 绑定返回数据
		resp := &filedtos.RecoveryDocumentResponse{
			Code:   rpcResp.Code,
			Status: "请求成功",
			Info:   respmsg.GetMsg(rpcResp.Code),
			Data:   rpcResp.Data,
		}
		ctx.JSON(http.StatusOK, resp)
	}
}
