package seed

import (
	"backend/models"
	"encoding/json"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func LoadPlans(path string, database *gorm.DB) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	var plans []models.Subscription
	err = json.NewDecoder(file).Decode(&plans)
	if err != nil {
		return err
	}

	for _, plan := range plans {
		err := database.Where(map[string]any{"type": plan.Type}).
    		Attrs(models.Subscription{Price: plan.Price, Description: plan.Description, FreeClasses: plan.FreeClasses, AdminID: plan.AdminID}).
    		FirstOrCreate(&plan)
			if err.Error != nil {
    			return err.Error
		}

	}
	return nil
}

func LoadAdmins(path string, database *gorm.DB) error {
	file, err := os.Open(path)
	if err != nil {
		return err	
	}
	defer file.Close()

	var admins []models.Admin
	err = json.NewDecoder(file).Decode(&admins)	
	if err != nil {
		return err
	}

	for _, admin := range admins {
		hash, herr := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
		if herr != nil {
			return herr
		}
		err := database.Where(map[string]any{"email": admin.Email}).
			Attrs(models.Admin{Name: admin.Name, Password: string(hash)}).
			FirstOrCreate(&admin)
		if err.Error != nil {
			return err.Error
		}
	}
	return nil
}

func LoadTrainers(path string, database *gorm.DB) error {
	file, err := os.Open(path)
	if err != nil {
		return err	
	}
	defer file.Close()

	var trainers []models.Trainer
	err = json.NewDecoder(file).Decode(&trainers)
	if err != nil {
		return err
	}	
	for _, trainer := range trainers {
		hash, herr := bcrypt.GenerateFromPassword([]byte(trainer.Password), bcrypt.DefaultCost)
		if herr != nil {
			return herr
		}
		err := database.Where(map[string]any{"email": trainer.Email}).
			Attrs(models.Trainer{Name: trainer.Name, Password: string(hash), Description: trainer.Description, Class: trainer.Class, AdminID: trainer.AdminID}).
			FirstOrCreate(&trainer)
		if err.Error != nil {
			return err.Error
		}	
	}
	return nil
}

func LoadClasses(path string, database *gorm.DB) error {
	file, err := os.Open(path)
	if err != nil {
		return err	
	}
	defer file.Close()

	var classes []models.Class
	err = json.NewDecoder(file).Decode(&classes)
	if err != nil {
		return err
	}

	for _, class := range classes {
		err := database.Where(map[string]any{"name": class.Name}).
			Attrs(models.Class{ScheduledAt: class.ScheduledAt, Description: class.Description, TrainerID: class.TrainerID, Capacity: class.Capacity, Users: class.Users, AdminID: class.AdminID}).
			FirstOrCreate(&class)
		if err.Error != nil {
			return err.Error
		}
	}
	return nil
}
