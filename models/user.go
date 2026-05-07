package models

import (
	"errors"
	"time"

	"hotbox-adm-backend/cli"
	"hotbox-adm-backend/pkg/constant"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SysUser struct {
	UserId        int64     `gorm:"column:user_id;primary_key;AUTO_INCREMENT"`
	LoginName     string    `gorm:"column:login_name"`
	Code          string    `gorm:"column:code"`
	UserPicUrl    string    `gorm:"column:user_pic_url"`
	Mobile        string    `gorm:"column:mobile"`
	UserName      string    `gorm:"column:username"`
	IdCardNo      string    `gorm:"column:id_card_no"`
	RealName      string    `gorm:"column:real_name" json:"-"`
	PayPwd        string    `gorm:"column:pay_pwd" json:"-"`
	Age           string    `gorm:"column:age"`
	Constell      string    `gorm:"column:constell"`
	Gender        string    `gorm:"column:gender"`
	Explain       string    `gorm:"column:explain"`
	CreateTime    string    `gorm:"column:create_time"`
	Certification uint32    `gorm:"column:certification" json:"certification"`
	FaceCertify   uint32    `gorm:"column:face_certify" json:"face_certify"`
	LoginTime     time.Time `gorm:"column:login_time" json:"login_time"`
	IsDelete      uint32    `gorm:"column:is_delete" json:"is_delete"`
	Privacy       int32     `gorm:"column:privacy" json:"privacy"`
}

type AdmUser struct {
	UserId int    `gorm:"column:user_id" json:"user_id"`
	OrgId  int    `gorm:"column:org_id" json:"org_id"`
	Name   string `gorm:"name" json:"name"`
}

type User struct {
	Ctx *gin.Context
}

func (u User) FindSysUserById(userId uint64) (sysUser *SysUser, err error) {
	if err = cli.HotDogGormDB.Model(&sysUser).
		Select("user_id,username,login_name,`code`,user_pic_url,mobile,id_card_no,real_name,create_time,certification,face_certify").
		Where("user_id = ? ", userId).First(&sysUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("查无此用户")
		}
		return nil, err
	}
	return
}

func (u User) FindAdmUserById(userId uint64) (adu *AdmUser, err error) {
	if err = cli.HotDogGormDB.Where("user_id = ? ", userId).Limit(1).Find(&adu).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("查无此用户")
		}
		return nil, err
	}
	return
}

func (u User) FindSysUserByMobile(mobile string) (sysUser *SysUser, err error) {
	if err = cli.HotDogGormDB.Model(&sysUser).
		Select("user_id,username,login_name,`code`,user_pic_url,mobile,id_card_no,real_name,create_time,certification,face_certify").
		Where("mobile = ? AND app_type = ? AND is_delete = ?", mobile, constant.APP_TYPE, 0).First(&sysUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("查无此用户")
		}
		return nil, err
	}
	return
}

func (u User) FindSysUserByParams(where map[string]any) (list []SysUser, err error) {
	if err = cli.HotDogGormDB.Model(&SysUser{}).
		Select("user_id,username,login_name,`code`,user_pic_url,mobile,id_card_no,real_name,create_time,certification,face_certify").
		Where(where).Scan(&list).Error; err != nil {
		return list, err
	}
	return
}

func (u User) GetSysUserByUserId(userId uint64) (sysUser *SysUser, err error) {
	if err = cli.HotDogGormDB.Model(&sysUser).
		Select("user_id,username,login_name,`code`,user_pic_url,mobile,id_card_no,real_name,create_time,certification,face_certify").
		Where("user_id", userId).First(&sysUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("查无此用户")
		}
		return nil, err
	}
	return
}

func (u User) DeleteTestUserByIds(userId []int64) error {
	return cli.HotDogGormDB.Model(&SysUser{}).Where("user_id in ? ", userId).UpdateColumn("is_delete", 1).Error
}
