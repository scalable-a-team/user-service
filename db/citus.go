package db

import (
	"gorm.io/gorm"
	"strings"
)

func InitDistributedTable(db *gorm.DB) error {
	queries := [...]string{
		"SELECT create_distributed_table('sellers', 'id')",
		"SELECT create_distributed_table('seller_wallets', 'seller_id')",
		"SELECT create_distributed_table('seller_profiles', 'seller_id')",
		"SELECT create_distributed_table('buyers', 'id')",
		"SELECT create_distributed_table('buyer_wallets', 'buyer_id')",
		"SELECT create_distributed_table('buyer_profiles', 'buyer_id')",
	}
	for _, query := range queries {
		if err := db.Exec(query).Error; err != nil {
			// If error is not "table .... is already distributed" then it's true error
			if !strings.Contains(err.Error(), "is already distributed") {
				return err
			}
		}
	}
	return nil
}
