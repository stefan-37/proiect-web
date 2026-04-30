package seed

import (
	"backend/models"	
	"encoding/json"
	"os"

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
    		Attrs(models.Subscription{Price: plan.Price, AdminID: plan.AdminID}).
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
		err := database.Where(map[string]any{"email": admin.Email}).
    		Attrs(models.Admin{Name: admin.Name, Password: admin.Password}).
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
		err := database.Where(map[string]any{"email": trainer.Email}).
			Attrs(models.Trainer{Name: trainer.Name, Password: trainer.Password, Description: trainer.Description, AdminID: trainer.AdminID}).
			FirstOrCreate(&trainer)
		if err.Error != nil {
			return err.Error
		}	
	}
	return nil
}	