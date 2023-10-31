package services

import (
	"practice04/entity"

	"gorm.io/gorm"
)

var UserLogs UserLogService

type UserLogService struct {
	DB *gorm.DB
}

func CreateUserLogService(db *gorm.DB) UserLogService {
	return UserLogService{db}
}

func (sl *UserLogService) AddOne(user_id uint, channel_name, operation string) (*entity.UserLog, error) {
	elem := &entity.UserLog{UserID: user_id, ChannelName: channel_name, Operation: operation}
	result := sl.DB.Create(elem)
	if result.Error != nil {
		return nil, result.Error
	}
	return elem, nil
}

func (sl *UserLogService) GetAll() ([]*entity.UserLog, error) {
	var data []*entity.UserLog
	result := sl.DB.Find(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

func (sl *UserLogService) DeleteAll() error {
	result := sl.DB.Unscoped().Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&entity.UserLog{})
	return result.Error
}

// func (sl *UserLogService) GetPeriod(user_id uint, date1, date2 string) ([]*entity.UserLog, error) {
// 	var res []*entity.UserLog
// 	d1, err := GetTimeDate(date1)
// 	if err != nil {
// 		return res, err
// 	}
// 	d2, err := GetTimeDate(date2)
// 	if err != nil {
// 		return res, err
// 	}
// 	result := sl.DB.Model(&entity.UserLog{}).Where("user_id = ? AND created_at BETWEEN (?) AND (?)", user_id, d1, d2).Find(&res)
// 	if result.Error != nil {
// 		return nil, result.Error
// 	}
// 	return res, nil
// }

// func (sl *UserLogService) GenerateCSV(user_id uint, date1, date2 string) (string, error) {
// 	logs, err := sl.GetPeriod(user_id, date1, date2)
// 	if err != nil {
// 		return "", err
// 	}
// 	name := "User" + fmt.Sprint(user_id) + "_" + date1 + "_" + date2
// 	err = CreateCSVUserLogs(name, logs)
// 	if err != nil {
// 		return "", err
// 	}
// 	return name + ".csv", nil
// }
