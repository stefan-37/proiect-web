package repository

import(
	"backend/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateAdmin(admin models.Admin, database *gorm.DB) error{
	hash, hasherr := bcrypt.GenerateFromPassword([]byte(admin.Password),bcrypt.DefaultCost)
	if hasherr != nil{
		return hasherr
	}
	admin.Password = string(hash)
	err := database.Create(&admin)
	if err.Error != nil{
		return err.Error
	}
	return nil
}

func GetAdminByEmail(email string, database *gorm.DB) (models.Admin, error){
	var admin models.Admin
	err := database.Where("email = ?", email).First(&admin)
	if err.Error != nil{
		return models.Admin{}, err.Error
	}
	return admin, nil
}

func GetAdminByID(id uint, database *gorm.DB) (models.Admin, error){
	var admin models.Admin
	err := database.Where("id = ?", id).First(&admin)
	if err.Error != nil{
		return models.Admin{}, err.Error
	}	
	return admin, nil
}

func UpdateAdmin(admin models.Admin, database *gorm.DB) error{
	err := database.Save(&admin)
	if err.Error != nil{
		return err.Error	
	}
	return nil
}

func DeleteAdminByEmail(email string, database *gorm.DB) error{
	err := database.Where("email = ?", email).Delete(&models.Admin{})
	if err.Error != nil{
		return err.Error
	}
	return nil
}	

