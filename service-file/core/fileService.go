package core

import (
	"context"
	"errors"
	"file/fileService"
	"file/pkg/consts"
	"file/pkg/statuscode"
	"file/pkg/utils/idmaker"
	iLogger "file/pkg/utils/logger"
	"file/repositry/dao"
	"file/repositry/model"
	"file/repositry/mq"
	"file/repositry/oss"
	"fmt"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
	"net/url"
	"strings"
	"time"
)

type FileService struct {
}

func (fs *FileService) GetAvatar(ctx context.Context, req *fileService.GetAvatarRequest, resp *fileService.AvatarResp) error {
	// 默认成功
	code := statuscode.Success
	resp.Code = uint32(code)

	// 获取操作用户的数据库对象
	userDao := dao.NewResourceDao(ctx)

	// 获取用户
	user, err := userDao.GetUserById(req.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果没有找到该记录
			resp.Code = statuscode.EmailNotRegister
			return nil
		} else {
			iLogger.LogrusObj.Error("根据userID查询用户时发生错误:", err)
			return err
		}
	}

	// 获取操作minio的客户端
	s3Client := oss.GetMinioClient()

	// 判断是否为默认头像
	objectName := ""
	if user.Avatar == "" {
		objectName = consts.DefaultAvatarName
	}

	// 搜索用户的头像
	prefix := fmt.Sprintf("%d", req.UserId)
	opts := minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	}
	objectCh := s3Client.ListObjects(ctx, consts.AvatarBucketName, opts)
	for obj := range objectCh {
		objectName = obj.Key // 根据前缀获取到该用户名
		break
	}

	// 设置请求参数
	reqParams := make(url.Values)

	// 获取用户的头像
	expires := time.Duration(consts.PresignedGetAvatarURLExpiresSecs) * time.Second
	presigenedURL, err := s3Client.PresignedGetObject(ctx, consts.AvatarBucketName, objectName, expires, reqParams)
	if err != nil {
		iLogger.LogrusObj.Error("生成头像预签名URL发生错误:", err)
		return err
	}

	// 绑定返回数据
	resp.AvatarUrl = presigenedURL.String()
	return nil
}

func (fs *FileService) UpdateAvatar(ctx context.Context, req *fileService.UpdateAvatarRequest, resp *fileService.AvatarResp) error {
	// 默认状态为成功
	code := statuscode.Success
	resp.Code = uint32(code)

	// 创建操作数据库的对象
	userDao := dao.NewResourceDao(ctx)

	// 根据用户id查询用户
	user, err := userDao.GetUserById(req.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果没有找到该记录
			resp.Code = statuscode.EmailNotRegister
			return nil
		} else {
			iLogger.LogrusObj.Error("通过UserID检索用户出错:", err)
			return err
		}

	}

	// 更新用户数据
	user.Avatar = req.AvatarName

	// 获取minio客户端
	s3Client := oss.GetMinioClient()

	// prefix
	prefix := fmt.Sprintf("%d/", user.UserId)

	// 检测是否已经存在
	// 这里不需要事务保护，由于 MinIO 对象存储是分布式的，单个删除操作是原子的，并且删除操作本身已经提供了一定程度的保护。
	// 这意味着，在删除对象时，MinIO 会确保要么对象被完全删除，要么不进行任何修改。
	// 这就排除了并发删除同一对象可能导致的问题。
	objectCh := s3Client.ListObjects(context.Background(), consts.AvatarBucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	// 循环遍历对象列表
	var deleteErrors []error // 用于存储删除对象时的错误
	for object := range objectCh {
		// 尝试删除对象
		err := s3Client.RemoveObject(context.Background(), consts.AvatarBucketName, object.Key, minio.RemoveObjectOptions{})
		if err != nil {
			// 如果删除失败，记录错误
			iLogger.LogrusObj.Error("更新头像时删除原有头像出错:", err)
			deleteErrors = append(deleteErrors, err)
		}
	}

	// 检查是否有删除失败的对象
	if len(deleteErrors) > 0 {
		// 返回删除对象时的错误
		return fmt.Errorf("删除原有头像时出错: %v", deleteErrors)
	}

	// 确定objectName
	objectName := fmt.Sprintf("%s%s", prefix, user.Avatar)

	// 确定上传的预签名
	expires := time.Duration(consts.PresignedPutAvatarURLExpiresSecs) * time.Second
	presignedURL, err := s3Client.PresignedPutObject(context.Background(), consts.AvatarBucketName, objectName, expires)

	// 绑定数据
	resp.AvatarUrl = presignedURL.String()

	// 更新用户入库
	if err := userDao.UpdateUser(user); err != nil {
		iLogger.LogrusObj.Error("用户更新头像文件时发生错误:", err)
		return err
	}

	return nil
}

func (fs *FileService) CreateNeuron(ctx context.Context, req *fileService.CreateNeuronRequest, resp *fileService.CreateNeuronResponse) error {
	// 默认成功
	resp.Code = statuscode.Success

	// 生成id
	id := idmaker.GenerateNeuronID()

	// 获取操作Neuron数据对象
	neuronDao := dao.NewResourceDao(ctx)

	// 创建Neuron
	now := time.Now()
	n := &model.Neuron{
		NeuronID:      id,
		Description:   "",
		OwnedBy:       req.UserId,
		CreatedAt:     &now,
		UpdatedAt:     &now,
		DocumentCount: 0,
		SessionCount:  0,
		TokenUsage:    0,
	}

	// 检查名称
	_, err := neuronDao.GetNeuronByNameAndUserID(req.NeuronName, req.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 没有同名神经元
			n.Name = req.NeuronName
		} else {
			// 查询过程发生错误
			iLogger.LogrusObj.Error("查询神经元名称时发生错误:", err)
			return err
		}
	} else {
		// 存在重名的神经元，计算同名神经元的数量，并设置NeuronName为NeuronName+count
		count := 1
		for {
			// Todo:这里可以性能优化，先全部查出来，避免多次访问数据库
			newNeuronName := fmt.Sprintf("%s(%d)", req.NeuronName, count)
			_, err := neuronDao.GetNeuronByNameAndUserID(newNeuronName, req.UserId)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// 新名称不重复，设置NeuronName为新名称
					n.Name = newNeuronName
					break
				} else {
					// 查询过程发生错误
					iLogger.LogrusObj.Error("查询神经元名称时发生错误:", err)
					return err
				}
			}
			count++
		}
	}

	// 数据入库
	err = neuronDao.CreateNeuron(n)
	if err != nil {
		iLogger.LogrusObj.Error("创建Neuron时候入库失败:", err)
		return err
	}

	// 绑定返回数据
	d := &fileService.NeuronData{
		NeuronId:   n.NeuronID,
		NeuronName: n.Name,
	}
	resp.Data = d
	return nil
}

func (fs *FileService) DeleteNeuron(ctx context.Context, req *fileService.DeleteNeuronRequest, resp *fileService.FileServiceResponse) error {
	// 默认请求成功
	resp.Code = statuscode.Success

	// 获取操作Neuron的数据库对象
	nDao := dao.NewResourceDao(ctx)

	// 确保该用户拥有这个Neuron
	n, err := nDao.GetNeuronByID(req.NeuronId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果没有找到记录
			resp.Code = statuscode.NeuronNotFound
			return nil
		}
		return err
	}

	// 确保该神经元属于该用户
	if n.OwnedBy != req.UserId {
		resp.Code = statuscode.UserHasNoNeuron
		return nil
	}

	// 从数据库中删除该条记录
	if err = nDao.DeleteNeuron(n.NeuronID); err != nil {
		iLogger.LogrusObj.Error("删除神经元发生错误:", err)
		return err
	}
	return nil
}

func (fs *FileService) RenameNeuron(ctx context.Context, req *fileService.RenameNeuronRequest, resp *fileService.RenameNeuronResponse) error {
	// 默认请求成功
	resp.Code = statuscode.Success

	// 获取操作Neuron的数据库对象
	nDao := dao.NewResourceDao(ctx)
	n, err := nDao.GetNeuronByID(req.NeuronId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果没有找到记录
			resp.Code = statuscode.NeuronNotFound
			return nil
		}
		return err
	}

	// 确保该神经元属于该用户
	if n.OwnedBy != req.UserId {
		resp.Code = statuscode.UserHasNoNeuron
		return nil
	}

	// 检查更新的名称
	fmt.Println(1)
	_, err = nDao.GetNeuronByNameAndUserID(req.NeuronName, req.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 命名未重复，继续检查新名称是否与原名称相同
			if n.Name != req.NeuronName {
				// 新名称不与原名称相同，更新神经元信息
				n.Name = req.NeuronName
			} else {
				// 新名称与原名称相同，不进行更新，直接返回
				resp.Code = statuscode.NameSameAsOriginal
				return nil
			}
		} else {
			// 查询过程发生错误
			iLogger.LogrusObj.Error("更新神经元名称时通过用户ID和新名称查询神经元时发生错误:", err)
			return err
		}
	} else {
		// 命名重复
		resp.Code = statuscode.NameExisted
		return nil
	}

	// 对象入库
	fmt.Println(2)
	if err = nDao.UpdateNeuron(n); err != nil {
		iLogger.LogrusObj.Error("重命名Neuron时发生错误:", err)
		return nil
	}

	// 绑定返回数据
	fmt.Println(3)
	d := &fileService.NeuronData{
		NeuronId:   n.NeuronID,
		NeuronName: n.Name,
	}
	resp.Data = d
	return nil
}

func (fs *FileService) ListNeuron(ctx context.Context, req *fileService.ListNeuronRequest, resp *fileService.ListNeuronResponse) error {
	// 默认请求成功
	resp.Code = statuscode.Success

	// 获取操作Neuron的数据库对象
	nDao := dao.NewResourceDao(ctx)

	// 获取Neurons
	ns, err := nDao.ListNeuron(req.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果没有找到该记录
			resp.Code = statuscode.NeuronNotFound
			return nil
		} else {
			iLogger.LogrusObj.Error("遍历Neurons时出错:", err)
			return err
		}
	}

	// 检查结果
	if len(ns) == 0 {
		resp.Data = []*fileService.NeuronData{}
		return nil
	}

	// 返回数据
	data := make([]*fileService.NeuronData, len(ns))
	for i, neuron := range ns {
		data[i] = &fileService.NeuronData{
			NeuronId:   neuron.NeuronID,
			NeuronName: neuron.Name,
		}
	}

	// 绑定数据
	resp.Data = data
	return nil
}

func (fs *FileService) UploadDocument(ctx context.Context, req *fileService.UploadDocumentRequest, resp *fileService.UploadDocumentResponse) error {
	// 默认请求成功
	resp.Code = statuscode.Success

	// 创建操作document数据库对象
	dDao := dao.NewResourceDao(ctx)

	// 生成文档id
	id := idmaker.GenerateDocumentID()

	// 创建一个document
	d := &model.Document{
		ID:          id,
		Type:        req.Type,
		Extension:   req.Extension,
		Size:        0,
		ContentHash: "",
		Metadata:    model.Metadata{},
		UploadedAt:  req.NeuronId,
		UploadedBy:  req.UserId,
		Status:      consts.FileStatusActivate, // 文件状态,如 Active, Deleted, Released 等
	}

	// 检查名称
	_, err := dDao.GetDocumentByNameUserIDAndNeuronID(req.Name, req.NeuronId, req.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 没有同名文件
			d.Name = req.Name
		} else {
			// 查询过程发生错误
			iLogger.LogrusObj.Error("查询文件名称时发生错误:", err)
			return err
		}
	} else {
		// 存在重命名的文件，计算重命名文件数量，设置新的名称
		count := 1
		nameOnly := strings.Trim(req.Name, req.Extension)
		for {
			newDocumentName := fmt.Sprintf("%s(%d).%s", nameOnly, count, req.Type)
			_, err = dDao.GetDocumentByNameUserIDAndNeuronID(newDocumentName, req.NeuronId, req.UserId)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// 新名称不重复，设置新名称
					d.Name = newDocumentName
					break
				} else {
					// 查询过程发生错误
					iLogger.LogrusObj.Error("查询文件名时发生错误:", err)
					return err
				}
			}
			count++
		}
	}

	// 对象入库
	if err = dDao.CreateDocument(d); err != nil {
		iLogger.LogrusObj.Error("文件对象入库时发生错误:", err)
		return err
	}

	// 生成oss的objectName
	objectName := fmt.Sprintf("%d/%s/%s", req.UserId, req.NeuronId, d.ID)

	// 生成上传文件的链接
	minioClient := oss.GetMinioClient()
	expires := time.Duration(consts.PresignedPutAvatarURLExpiresSecs) * time.Second
	presignedURL, err := minioClient.PresignedPutObject(ctx, consts.FileBucketName, objectName, expires)
	if err != nil {
		return err
	}

	// 返回数据
	data := &fileService.DocumentData{
		DocumentId: d.ID,
		Name:       d.Name,
		Type:       d.Type,
	}

	// 绑定返回数据
	resp.Data = data
	resp.FileUrl = presignedURL.String()
	return nil
}

func (fs *FileService) DeleteDocument(ctx context.Context, req *fileService.DeleteDocumentRequest, resp *fileService.FileServiceResponse) error {
	// 默认请求成功
	resp.Code = statuscode.Success

	// 获取操作文件数据库的对象
	dDao := dao.NewResourceDao(ctx)

	// 根据id获取到文件信息
	d, err := dDao.GetDocumentByID(req.DocumentId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果没有找到该记录
			resp.Code = statuscode.DocumentNotFound
			return nil
		} else {
			iLogger.LogrusObj.Error("通过文件ID查询文件时出错:", err)
			return err
		}
	}

	// 校验信息
	if d.UploadedBy != req.UserId || d.UploadedAt != req.NeuronId {
		resp.Code = statuscode.InvalidDocumentOwnership
		return nil
	}

	// 校验通过，设置文件状态，文件将被放入垃圾桶
	d.Status = consts.FileStatusDeleted

	// 使用消息队列来实现定时删除
	mqConn := mq.GetMQClient()

	// 获取channel
	ch, err := mqConn.Channel()
	if err != nil {
		iLogger.LogrusObj.Error("发布删除document消息时创建消息队列channel失败:", err)
		return err
	}

	// 创建队列
	if err = mq.InitUserQueueAndBind(ch, req.DocumentId); err != nil {
		iLogger.LogrusObj.Error("创建删除document队列时失败:", err)
		return err
	}

	// 发布消息
	if err = mq.PublishMessage(ch, req.UserId, req.DocumentId); err != nil {
		iLogger.LogrusObj.Error("发布删除document消息时发布删除消息失败:", err)
		return err
	}

	// 更新document的信息
	if err = dDao.UpdateDocument(d); err != nil {
		iLogger.LogrusObj.Error("删除document时更新文件信息失败:", err)
		return err
	}
	return nil
}

func (fs *FileService) ListDocuments(ctx context.Context, req *fileService.ListDocumentsRequest, resp *fileService.ListDocumentsResponse) error {
	// 默认请求成功
	resp.Code = statuscode.Success

	// 获取操作document数据库对象
	dDao := dao.NewResourceDao(ctx)

	// 获取资源
	ds, err := dDao.GetDocumentByNeuronID(req.NeuronId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果没有找到该记录
			resp.Code = statuscode.DocumentNotFound
			return nil
		} else {
			iLogger.LogrusObj.Error("获取神经元下面所有文件出错:", err)
			return err
		}
	}

	// 返回数据
	data := make([]*fileService.DocumentData, 0, len(ds)) // 预先分配足够的资源

	// 遍历资源
	for _, d := range ds {
		// 判断文件状态
		if d.Status != consts.FileStatusActivate {
			continue
		}

		// 判断是否数据当前用户
		if d.UploadedBy != req.UserId {
			continue
		}

		// 绑定数据
		temp := &fileService.DocumentData{
			DocumentId: d.ID,
			Name:       d.Name,
			Type:       d.Type,
		}

		// 添加数据
		data = append(data, temp)
	}

	// 绑定返回资源
	resp.Data = data
	return nil
}

func (fs *FileService) RenameDocument(ctx context.Context, req *fileService.RenameDocumentRequest, resp *fileService.RenameDocumentResponse) error {
	// 默认请求成功
	resp.Code = statuscode.Success

	// 获取操作document数据库对象
	dDao := dao.NewResourceDao(ctx)

	// 获取指定对象
	d, err := dDao.GetDocumentByID(req.DocumentId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果没有找到该记录
			resp.Code = statuscode.DocumentNotFound
			return nil
		} else {
			iLogger.LogrusObj.Error("重命名文件时候根据文件ID查询文件未找到")
			return err
		}
	}

	// 检查新名称
	_, err = dDao.GetDocumentByNameAndNeronID(req.NewName, req.NeuronId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 命名重复，继续检查新名称与原名称相同
			if d.Name != req.NewName {
				// 新名称与原名称不相同，更新文件名称
				d.Name = req.NewName
			} else {
				// 新旧名称相同
				resp.Code = statuscode.NameSameAsOriginal
				return nil
			}
		} else {
			// 查询过程发生错误
			iLogger.LogrusObj.Error("更新文件名时通过神经元ID和新名称查询时发生错误:", err)
			return err
		}
	} else {
		// 命名重复
		count := 1
		for {
			newNeuronName := fmt.Sprintf("%s(%d)", req.NewName, count)
			if _, err = dDao.GetDocumentByNameAndNeronID(newNeuronName, req.NeuronId); err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// 新名称不重复
					d.Name = newNeuronName
					break
				} else {
					// 查询过程发生错误
					iLogger.LogrusObj.Error("查询神经元时发生错误:", err)
					return err
				}
			}
			count++
		}
	}

	// 更新用户数据入库
	if err = dDao.UpdateDocument(d); err != nil {
		iLogger.LogrusObj.Error("重命名文件时文件信息入库失败:", err)
		return err
	}

	// 返回数据
	data := &fileService.DocumentData{
		DocumentId: d.ID,
		Name:       d.Name,
		Type:       d.Type,
	}
	resp.Data = data
	return nil
}

func (fs *FileService) ListDocumentsInBin(ctx context.Context, req *fileService.ListDocumentsInBinRequest, resp *fileService.ListDocumentsInBinResponse) error {
	// 默认返回成功
	resp.Code = statuscode.Success

	// 操作document数据对象
	dDao := dao.NewResourceDao(ctx)

	// 查询对象
	ds, err := dDao.GetDocumentByUserID(req.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果没有找到该记录
			resp.Code = statuscode.DocumentNotFound
			return nil
		} else {
			iLogger.LogrusObj.Error("通过userID查询用户下面所有文件出错:", err)
			return err
		}
	}

	// 返回数据
	data := make([]*fileService.DocumentData, 0, len(ds))
	for _, d := range ds {
		if d.Status != consts.FileStatusDeleted {
			continue
		}
		t := &fileService.DocumentData{
			DocumentId: d.ID,
			Name:       d.Name,
			Type:       d.Type,
		}
		data = append(data, t)

	}

	// 绑定返回数据
	resp.Data = data
	return nil
}

func (fs *FileService) RecoveryDocument(ctx context.Context, req *fileService.RecoveryDocumentRequest, resp *fileService.RecoveryDocumentResponse) error {
	// 默认请求成功
	resp.Code = statuscode.Success

	// 获取操作document的数据库对象
	dDao := dao.NewResourceDao(ctx)

	// 更新文件信息
	d, err := dDao.GetDocumentByID(req.DocumentId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果没有找到该记录
			resp.Code = statuscode.DocumentNotFound
			return nil
		} else {
			iLogger.LogrusObj.Error("通过ID查询文件时出错:", err)
			return err
		}
	}

	// 使用消息队列来实现定时删除
	mqConn := mq.GetMQClient()

	// 获取channel
	ch, err := mqConn.Channel()
	if err != nil {
		iLogger.LogrusObj.Error("恢复删除document消息时创建消息队列channel失败:", err)
		return err
	}

	// 正常消费
	ok, err := mq.RecoveryConsume(ch, req.DocumentId)
	if err != nil {
		iLogger.LogrusObj.Error("恢复删除document消息时从队列消费消息失败:", err)
		return err
	}

	// 检查文件状态
	if ok {
		d.Status = consts.FileStatusActivate
	}

	// 更新用户入库
	if err = dDao.UpdateDocument(d); err != nil {
		iLogger.LogrusObj.Error("恢复删除document消息时更新状态入库时失败:", err)
		return err
	}

	// 绑定数据
	data := &fileService.DocumentData{
		DocumentId: d.ID,
		Name:       d.Name,
		Type:       d.Type,
	}
	resp.Data = data

	return nil
}
