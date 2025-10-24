package migrations
	import (
		"gorm.io/gorm"
		"log"
		"github.com/codetheuri/poster-gen/internal/app/posters/models"
)
		// Createposterstable struct implements migration interface
		type Createposterstable struct {}

		func (m *Createposterstable) Version() string{
			return "20251011151351"
			}
		func (m *Createposterstable) Name() string {
			return "create_posters_table"
		}	
			//up migration method
		func (m *Createposterstable) Up(tx *gorm.DB) error {
		log.Printf("Running Up migration: %s", m.Name())
		// Add your migration logic here
		//example : 
		// 1. using SQL calls
		// if  err := tx.Exec("CREATE TABLE IF NOT EXISTS example (id INT PRIMARY KEY)").Error; err != nil {
		// 	return err
		// }

		// 2. using gorm methods
		//type NewModel struct {
		// gorm.Model
		// Field1 string
		//}
		// if err := tx.AutoMigrate(&NewModel{}); err != nil {
		// 	return err
		// }
		if err := tx.AutoMigrate(&models.Asset{}); err != nil {
			return err
		}
		// if err := tx.AutoMigrate(&models.Order{}); err != nil {
		// 	return err
		// }
		if err := tx.AutoMigrate(&models.Layout{}); err != nil {
			return err
		}
		if err := tx.AutoMigrate(&models.PosterTemplate{}); err != nil {
			return err
		}
		if err := tx.AutoMigrate(&models.Poster{}); err != nil {
			return err
		}
		log.Printf("Successfully applied Up migration: %s", m.Name())
		return nil
		}
		//down migration method
		func (m *Createposterstable) Down(tx *gorm.DB) error {
		log.Printf("Running Down migration: %s", m.Name())
		// Example:
		// 1. using SQL calls
		// if err := tx.Exec("DROP TABLE IF EXISTS example").Error; err != nil {
		// 	return err
		// }

		// 2. using gorm methods
		// if err := tx.Migrator().DropTable("new_models"); err != nil {
		// 	return err
		// }
		if err := tx.Migrator().DropTable("posters"); err != nil {
			return err
		}
		// if err := tx.Migrator().DropTable("orders"); err != nil {
		// 	return err
		// }
		if err := tx.Migrator().DropTable("poster_templates"); err != nil {
			return err
		}
		if err := tx.Migrator().DropTable("layouts"); err != nil {
			return err
		}
		if err := tx.Migrator().DropTable("assets"); err != nil {
			return err
		}
		log.Printf("Successfully applied Down migration: %s", m.Name())
		return nil
		}

		func init() {
		  // Register the migration
		  RegisteredMigrations = append(RegisteredMigrations, &Createposterstable{})
		}
