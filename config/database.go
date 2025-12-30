package config

import (
	"fmt"
	"os"
	"workHub/internal/entity"
	"workHub/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectDatabase trả về *gorm.DB và error
func ConnectDatabase() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_MASTER_HOST"),
		os.Getenv("DB_MASTER_USER"),
		os.Getenv("DB_MASTER_PASSWORD"),
		os.Getenv("DB_MASTER_NAME"),
		os.Getenv("DB_MASTER_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect DB: %w", err)
	}

	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	logger.Info("config", "AutoMigrate", "Starting database migration...")
	
	// Xóa các bảng junction trước nếu tồn tại (để tránh conflict foreign key)
	db.Exec("DROP TABLE IF EXISTS role_permissions CASCADE")
	db.Exec("DROP TABLE IF EXISTS user_roles CASCADE")
	db.Exec("DROP TABLE IF EXISTS workspace_members CASCADE")
	
	// Xóa và tạo lại các bảng cha nếu có vấn đề với primary key
	// (Chỉ xóa nếu cần thiết - có thể comment nếu muốn giữ data)
	// db.Exec("DROP TABLE IF EXISTS roles CASCADE")
	// db.Exec("DROP TABLE IF EXISTS permissions CASCADE")
	
	// Xử lý migration cho bảng users - thêm các cột mới nếu chưa có và set giá trị mặc định
	db.Exec(`
		DO $$ 
		BEGIN
			-- Thêm cột full_name nếu chưa có
			IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='users' AND column_name='full_name') THEN
				ALTER TABLE users ADD COLUMN full_name TEXT;
				-- Set giá trị mặc định cho các record cũ
				UPDATE users SET full_name = COALESCE(username, email, 'User') WHERE full_name IS NULL;
				-- Sau đó mới thêm NOT NULL constraint
				ALTER TABLE users ALTER COLUMN full_name SET NOT NULL;
			ELSE
				-- Nếu cột đã tồn tại nhưng có giá trị NULL, update chúng
				UPDATE users SET full_name = COALESCE(username, email, 'User') WHERE full_name IS NULL;
				-- Thêm NOT NULL constraint nếu chưa có
				ALTER TABLE users ALTER COLUMN full_name SET NOT NULL;
			END IF;
			
			-- Thêm cột phone_number nếu chưa có
			IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='users' AND column_name='phone_number') THEN
				ALTER TABLE users ADD COLUMN phone_number TEXT;
			END IF;
			
			-- Thêm cột address nếu chưa có
			IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='users' AND column_name='address') THEN
				ALTER TABLE users ADD COLUMN address TEXT;
			END IF;
			
			-- Thêm cột role nếu chưa có và set giá trị mặc định
			IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='users' AND column_name='role') THEN
				ALTER TABLE users ADD COLUMN role TEXT DEFAULT 'student';
				UPDATE users SET role = 'student' WHERE role IS NULL;
				ALTER TABLE users ALTER COLUMN role SET NOT NULL;
			END IF;
		END $$;
	`)
	
	// Migrate các bảng cha trước (không có foreign key)
	err := db.AutoMigrate(
		&entity.User{},
		&entity.Role{},
		&entity.Permission{},
		&entity.Session{},
		&entity.Workspace{},
		&entity.Book{},
	)
	if err != nil {
		logger.Error("config", "AutoMigrate", fmt.Sprintf("Migration failed (base tables): %v", err))
		return fmt.Errorf("failed to auto migrate base tables: %w", err)
	}
	
	// Đảm bảo primary key tồn tại trên các bảng cha
	// Kiểm tra và tạo primary key nếu chưa có (sử dụng DO block để tránh lỗi nếu đã tồn tại)
	db.Exec(`
		DO $$ 
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'roles_pkey') THEN
				ALTER TABLE roles ADD CONSTRAINT roles_pkey PRIMARY KEY (id);
			END IF;
		END $$;
	`)
	
	db.Exec(`
		DO $$ 
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'permissions_pkey') THEN
				ALTER TABLE permissions ADD CONSTRAINT permissions_pkey PRIMARY KEY (id);
			END IF;
		END $$;
	`)
	
	db.Exec(`
		DO $$ 
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'users_pkey') THEN
				ALTER TABLE users ADD CONSTRAINT users_pkey PRIMARY KEY (id);
			END IF;
		END $$;
	`)
	
	// Sau đó migrate các bảng có foreign key
	err = db.AutoMigrate(
		&entity.RolePermission{},
		&entity.UserRole{},
		&entity.WorkspaceMember{},
		&entity.Borrow{},
	)
	if err != nil {
		logger.Error("config", "AutoMigrate", fmt.Sprintf("Migration failed (junction tables): %v", err))
		return fmt.Errorf("failed to auto migrate junction tables: %w", err)
	}
	
	logger.Info("config", "AutoMigrate", "Database migration completed: roles, users, sessions, permissions, workspaces, books, borrows, role_permissions, user_roles, workspace_members tables created/updated")
	return nil
}
