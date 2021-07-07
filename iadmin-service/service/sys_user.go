package service

import (
	"errors"
	"ginProject/global"
	"ginProject/model"
	"ginProject/model/request"
	"ginProject/utils"
	"gorm.io/gorm"
)

func Login(u *model.SysUser) (err error, userInter *model.SysUser) {
	var user model.SysUser
	u.Password = utils.MD5V([]byte(u.Password))
	err = global.GvaDb.Where("phone = ? AND password = ?", u.Phone, u.Password).Preload("Role").First(&user).Error
	return err, &user
}

func Register(u model.SysUser) (err error, userInter model.SysUser) {
	var user model.SysUser
	//err = global.GvaDb.Where("phone = ?",u.Phone).First(&user).Error
	if err = global.GvaDb.Where("phone = ?", u.Phone).First(&user).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("该手机号已注册！"), userInter
	}
	u.Password = utils.MD5V([]byte(u.Password))
	err = global.GvaDb.Create(&u).Error
	return err, u
}

func FindSysUserById(id uint) (err error, user *model.SysUser) {
	var u model.SysUser
	if err = global.GvaDb.Where("id = ?", id).First(&u).Error; err != nil {
		return errors.New("用户不存在!"), &u
	}
	return nil, &u
}

func GetSysUserList(query *request.SysUserListSearch) (err error, total int64, list []model.SysUser) {
	db := global.GvaDb.Model(&model.SysUser{}).Where("phone like ? and real_name like ?", "%"+query.Phone+"%", "%"+query.RealName+"%")
	if len(query.CreatedAt) > 0 {
		db.Where("created_at between ? and ?", query.CreatedAt[0], query.CreatedAt[1])
	}
	err = db.Count(&total).Error
	if err != nil {
		return err, 0, nil
	}
	err = db.Scopes(utils.Paginate(query.Pagination.Current, query.Pagination.PageSize)).Order("created_at desc").Preload("Role").Find(&list).Error
	return err, total, list
}

func UpdateSysUserInfo(reqUser model.SysUser) (err error, user model.SysUser) {
	err = global.GvaDb.Updates(&reqUser).Error
	return err, reqUser
}

func AddSysUserInfo(reqUser model.SysUser) (err error, userInter model.SysUser) {
	var user model.SysUser
	if err = global.GvaDb.Where("phone = ?", reqUser.Phone).First(&user).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("该手机号已注册！"), userInter
	}
	reqUser.Password = utils.MD5V([]byte(utils.DefaultPassword))
	err = global.GvaDb.Create(&reqUser).Error
	return err, reqUser
}

func DeleteBatchSysUser(ids request.IdsReq) (err error) {
	err = global.GvaDb.Delete(&[]model.SysUser{}, "id in (?)", ids.Ids).Error
	return err
}
