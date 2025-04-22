package repo

import (
	"api/internal/entities"
	"api/internal/infrastructure/database"
	"context"
	"time"

	"gorm.io/gorm"
)

type organizationGORM struct {
	ID        uint32    `gorm:"primary_key"`
	Name      string    `gorm:"column:name"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

type OrganizationDB struct {
	dbManager database.DBManager
	tableName string
}

func NewOrganizationDB(dbManager database.DBManager) OrganizationDB {
	return OrganizationDB{
		dbManager: dbManager,
		tableName: "organizations",
	}
}

func (u OrganizationDB) Create(ctx context.Context, organization entities.Organization) (entities.Organization, error) {
	var organizationGORM = organizationGORM{
		Name:      organization.Name(),
		CreatedAt: organization.CreatedAt(),
	}

	if err := u.dbManager.With(ctx).Table(u.tableName).Create(&organizationGORM).Error; err != nil {
		return entities.Organization{}, err
	}

	organization.SetID(organizationGORM.ID)
	return organization, nil
}

func (u OrganizationDB) FindByID(ctx context.Context, id uint32) (entities.Organization, error) {
	var organizationGORM organizationGORM

	if err := u.dbManager.With(ctx).Table(u.tableName).Where("id = ?", id).First(&organizationGORM).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return entities.Organization{}, entities.ErrOrganizationNotFound
		}
		return entities.Organization{}, err
	}

	return entities.NewOrganization(
		organizationGORM.ID,
		organizationGORM.Name,
		organizationGORM.CreatedAt,
	), nil
}

func (u OrganizationDB) FindAll(ctx context.Context) ([]entities.Organization, int64, error) {
	var (
		organizationGORMs []organizationGORM
		total             int64
	)

	if err := u.dbManager.With(ctx).Table(u.tableName).Count(&total).Find(&organizationGORMs).Error; err != nil {
		return nil, 0, err
	}

	organizations := make([]entities.Organization, 0, len(organizationGORMs))
	for _, organizationGORM := range organizationGORMs {
		organizations = append(organizations, entities.NewOrganization(
			organizationGORM.ID,
			organizationGORM.Name,
			organizationGORM.CreatedAt,
		))
	}
	return organizations, total, nil
}
