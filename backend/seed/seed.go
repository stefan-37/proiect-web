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
