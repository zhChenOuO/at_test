package pkg

import gorm "gorm.io/gorm"

// IService 抽象，請在這定義要實作的方法
type IService interface {
	IdentityAccountService
	EamilService
	PhoneService
}

// IRepository 抽象，請在這定義要實作的方法
type IRepository interface {
	WriteDB() *gorm.DB
	IdentityAccountRepository
}
