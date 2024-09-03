package gateways

import (
	"github.com/skinnykaen/rpa_clone/internal/db"
	"github.com/skinnykaen/rpa_clone/internal/models"
	"gorm.io/gorm/clause"
)

//TODO maybe return []student or []parent а не rels

type ParentRel interface {
	CreateRel(parentId, childId uint) (models.ParentRelCore, error)
	DeleteRel(id uint) (err error)
	GetRelsByParentId(parentId uint) (rels []models.ParentRelCore, err error)
	GetRelsByChildId(childId uint) (rels []models.ParentRelCore, err error)
}

type ParentRelGatewayImpl struct {
	postgresClient db.PostgresClient
}

func (p ParentRelGatewayImpl) CreateRel(parentId, childId uint) (models.ParentRelCore, error) {
	rel := models.ParentRelCore{
		ParentID: parentId,
		ChildID:  childId,
	}
	result := p.postgresClient.Db.Create(&rel).Clauses(clause.Returning{})
	return rel, result.Error
}

func (p ParentRelGatewayImpl) DeleteRel(id uint) (err error) {
	return p.postgresClient.Db.Delete(&models.ParentRelCore{}, id).Error
}

func (p ParentRelGatewayImpl) GetRelsByParentId(parentId uint) (rels []models.ParentRelCore, err error) {
	result := p.postgresClient.Db.Where(models.ParentRelCore{ParentID: parentId}).Find(&rels)
	return rels, result.Error
}

func (p ParentRelGatewayImpl) GetRelsByChildId(childId uint) ([]models.ParentRelCore, error) {
	//TODO implement me
	panic("implement me")
}
