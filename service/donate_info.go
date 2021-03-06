package service

import (
	"errors"
	"github.com/AH-dark/random-donate/model"
	"github.com/AH-dark/random-donate/pkg/utils"
	"gorm.io/gorm"
	"math/rand"
)

func DonateInfoIsExist(info *model.DonateInfo) (bool, error) {
	count, err := Count(&model.DonateInfo{}, info)
	if err != nil {
		return false, err
	}

	utils.Log().Debug("Count Table DonateInfo, has %d rows", count)

	return count != 0, nil
}

func DonateInfoSave(info *model.DonateInfo) error {
	err := model.DB.Save(&info).Error
	return err
}

func DonateInfoFind(info *model.DonateInfo) (model.DonateInfo, error) {
	var dbData model.DonateInfo
	err := model.DB.Where(&info).First(&dbData).Error
	return dbData, err
}

func DonateInfoRandomGet(prevId uint) (model.DonateInfo, error) {
	var data model.DonateInfo
	not := model.DonateInfo{
		Model: gorm.Model{
			ID: prevId,
		},
	}

	var count int64

	err := model.DB.Model(&model.DonateInfo{}).Not(&not).Count(&count).Error
	if err != nil {
		return data, err
	}

	if count < 1 {
		return data, errors.New("no data exist")
	}

	err = model.DB.Not(&not).Offset(rand.Intn(int(count))).First(&data).Error
	if err != nil {
		return data, err
	}

	return data, nil
}
