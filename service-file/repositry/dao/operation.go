package dao

import (
	"context"
	"file/repositry/model"
	"gorm.io/gorm"
)

type ResourceDao struct {
	*gorm.DB
}

func NewResourceDao(ctx context.Context) *ResourceDao {
	return &ResourceDao{_db}
}

func NewUserDaoByDb(db *gorm.DB) *ResourceDao {
	return &ResourceDao{
		db,
	}
}

func (rd *ResourceDao) GetUserById(userId string) (user *model.User, err error) {
	user = &model.User{}
	err = rd.DB.Model(&model.User{}).
		Where("user_id = ?", userId).
		First(user).Error
	return
}

func (rd *ResourceDao) UpdateUser(user *model.User) error {
	err := rd.DB.Model(&model.User{}).
		Where("user_id = ?", user.UserId).
		Save(user).Error
	return err
}

func (rd *ResourceDao) CreateNeuron(n *model.Neuron) error {
	err := rd.DB.Model(&model.Neuron{}).
		Create(n).Error
	return err
}

func (rd *ResourceDao) GetNeuronByID(id string) (*model.Neuron, error) {
	var n model.Neuron
	err := rd.DB.Model(&model.Neuron{}).
		Where("neuron_id = ?", id).First(&n).Error
	return &n, err
}

func (rd *ResourceDao) GetNeuronByNameAndUserID(name string, userId string) (*model.Neuron, error) {
	var n model.Neuron
	err := rd.DB.Model(&model.Neuron{}).
		Where("name = ? AND owned_by = ?", name, userId).
		First(&n).Error
	return &n, err
}

func (rd *ResourceDao) DeleteNeuron(id string) error {
	// 根据指定的列和条件删除记录
	err := rd.DB.Where("neuron_id = ?", id).
		Delete(&model.Neuron{}).Error
	return err
}

func (rd *ResourceDao) UpdateNeuron(n *model.Neuron) error {
	err := rd.DB.Model(&model.Neuron{}).
		Where("neuron_id = ?", n.NeuronID).
		Save(&n).Error
	return err
}

func (rd *ResourceDao) ListNeuron(userId string) (ns []*model.Neuron, err error) {
	err = rd.DB.Model(&model.Neuron{}).
		Where("owned_by = ?", userId).
		Find(&ns).Error
	return
}

func (rd *ResourceDao) GetDocumentByNameUserIDAndNeuronID(name, neuronId string, userId string) (*model.Document, error) {
	var document model.Document
	if err := rd.DB.Model(&model.Document{}).
		Where("name = ? AND uploaded_by = ? AND uploaded_at = ?", name, userId, neuronId).
		First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (rd *ResourceDao) CreateDocument(d *model.Document) error {
	err := rd.DB.Model(&model.Document{}).
		Create(d).Error
	return err
}

func (rd *ResourceDao) GetDocumentByID(ID string) (*model.Document, error) {
	var d model.Document
	err := rd.DB.Model(&model.Document{}).
		Where("id = ?", ID).First(&d).
		Error
	return &d, err
}

func (rd *ResourceDao) DeleteDocumentByID(ID string) error {
	err := rd.DB.Where("id = ?", ID).
		Delete(&model.Document{}).Error
	return err
}

func (rd *ResourceDao) UpdateDocument(d *model.Document) error {
	err := rd.DB.Model(&model.Document{}).
		Where("id = ?", d.ID).
		Save(&d).Error
	return err
}

func (rd *ResourceDao) GetDocumentByNeuronID(nd string) ([]*model.Document, error) {
	var ds []*model.Document
	err := rd.DB.Model(&model.Document{}).
		Where("uploaded_at = ?", nd).
		Find(&ds).Error
	return ds, err
}

func (rd *ResourceDao) GetDocumentByNameAndNeronID(name, nd string) (*model.Document, error) {
	var d model.Document
	err := rd.DB.Model(&model.Document{}).
		Where("name = ? and uploaded_at = ?", name, nd).
		First(&d).Error
	return &d, err
}

func (rd *ResourceDao) GetDocumentByUserID(userId string) ([]*model.Document, error) {
	var ds []*model.Document
	err := rd.DB.Model(&model.Document{}).
		Where("uploaded_by = ?", userId).
		Find(&ds).Error
	return ds, err
}
