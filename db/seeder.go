package db

import (
	"fmt"
	"workHub/internal/entity"
	"workHub/logger"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SeedBooks tạo dữ liệu mẫu cho bảng books
func SeedBooks(db *gorm.DB) error {
	logger.Info("seeder", "SeedBooks", "Starting to seed books...")

	books := []entity.Book{
		{
			ID:              uuid.New().String(),
			Title:           "Những Người Khốn Khổ",
			Author:          "Victor Hugo",
			ISBN:            "978-604-1-12345-6",
			Category:        "Novel",
			Description:     getStringPtr("Tác phẩm kinh điển của văn học Pháp, kể về cuộc đời của Jean Valjean trong xã hội Pháp thế kỷ 19."),
			TotalCopies:     5,
			AvailableCopies: 5,
			PublishedYear:   getIntPtr(1862),
			Publisher:       getStringPtr("Nhà xuất bản Văn học"),
		},
		{
			ID:              uuid.New().String(),
			Title:           "Sapiens: Lược Sử Loài Người",
			Author:          "Yuval Noah Harari",
			ISBN:            "978-604-1-23456-7",
			Category:        "History",
			Description:     getStringPtr("Cuốn sách khám phá lịch sử và quá trình phát triển của loài người từ thời kỳ đồ đá đến hiện tại."),
			TotalCopies:     8,
			AvailableCopies: 8,
			PublishedYear:   getIntPtr(2011),
			Publisher:       getStringPtr("Nhà xuất bản Tri thức"),
		},
		{
			ID:              uuid.New().String(),
			Title:           "Clean Code: A Handbook of Agile Software Craftsmanship",
			Author:          "Robert C. Martin",
			ISBN:            "978-604-1-34567-8",
			Category:        "Technology",
			Description:     getStringPtr("Hướng dẫn viết code sạch, dễ đọc và dễ bảo trì cho các lập trình viên."),
			TotalCopies:     10,
			AvailableCopies: 10,
			PublishedYear:   getIntPtr(2008),
			Publisher:       getStringPtr("Prentice Hall"),
		},
		{
			ID:              uuid.New().String(),
			Title:           "Flutter Complete Reference",
			Author:          "Alberto Miola",
			ISBN:            "978-604-1-45678-9",
			Category:        "Technology",
			Description:     getStringPtr("Hướng dẫn toàn diện về Flutter framework để phát triển ứng dụng mobile đa nền tảng."),
			TotalCopies:     6,
			AvailableCopies: 6,
			PublishedYear:   getIntPtr(2020),
			Publisher:       getStringPtr("Packt Publishing"),
		},
		{
			ID:              uuid.New().String(),
			Title:           "Đất Rừng Phương Nam",
			Author:          "Đoàn Giỏi",
			ISBN:            "978-604-1-56789-0",
			Category:        "Novel",
			Description:     getStringPtr("Tiểu thuyết nổi tiếng về cuộc sống và con người vùng đất Nam Bộ trong thời kỳ kháng chiến."),
			TotalCopies:     7,
			AvailableCopies: 7,
			PublishedYear:   getIntPtr(1957),
			Publisher:       getStringPtr("Nhà xuất bản Kim Đồng"),
		},
		{
			ID:              uuid.New().String(),
			Title:           "A Brief History of Time",
			Author:          "Stephen Hawking",
			ISBN:            "978-604-1-67890-1",
			Category:        "Science",
			Description:     getStringPtr("Cuốn sách phổ biến về vũ trụ học, thuyết tương đối, hố đen và nguồn gốc của vũ trụ."),
			TotalCopies:     5,
			AvailableCopies: 5,
			PublishedYear:   getIntPtr(1988),
			Publisher:       getStringPtr("Bantam Books"),
		},
		{
			ID:              uuid.New().String(),
			Title:           "The Design of Everyday Things",
			Author:          "Don Norman",
			ISBN:            "978-604-1-78901-2",
			Category:        "Technology",
			Description:     getStringPtr("Cuốn sách về thiết kế UX/UI và tâm lý học nhận thức trong thiết kế sản phẩm."),
			TotalCopies:     4,
			AvailableCopies: 4,
			PublishedYear:   getIntPtr(1988),
			Publisher:       getStringPtr("Basic Books"),
		},
		{
			ID:              uuid.New().String(),
			Title:           "Lịch Sử Việt Nam Từ Nguồn Gốc Đến Giữa Thế Kỷ XX",
			Author:          "Lê Thành Khôi",
			ISBN:            "978-604-1-89012-3",
			Category:        "History",
			Description:     getStringPtr("Công trình nghiên cứu toàn diện về lịch sử Việt Nam từ thời tiền sử đến thế kỷ 20."),
			TotalCopies:     6,
			AvailableCopies: 6,
			PublishedYear:   getIntPtr(1955),
			Publisher:       getStringPtr("Nhà xuất bản Thế giới"),
		},
		{
			ID:              uuid.New().String(),
			Title:           "Tôi Tài Giỏi, Bạn Cũng Thế",
			Author:          "Adam Khoo",
			ISBN:            "978-604-1-90123-4",
			Category:        "Science",
			Description:     getStringPtr("Phương pháp học tập hiệu quả và kỹ thuật ghi nhớ để đạt được thành công trong học tập."),
			TotalCopies:     12,
			AvailableCopies: 12,
			PublishedYear:   getIntPtr(2009),
			Publisher:       getStringPtr("Nhà xuất bản Trẻ"),
		},
		{
			ID:              uuid.New().String(),
			Title:           "Dế Mèn Phiêu Lưu Ký",
			Author:          "Tô Hoài",
			ISBN:            "978-604-1-01234-5",
			Category:        "Novel",
			Description:     getStringPtr("Truyện đồng thoại nổi tiếng về chú dế mèn và những cuộc phiêu lưu kỳ thú."),
			TotalCopies:     15,
			AvailableCopies: 15,
			PublishedYear:   getIntPtr(1941),
			Publisher:       getStringPtr("Nhà xuất bản Kim Đồng"),
		},
		{
			ID:              uuid.New().String(),
			Title:           "The Pragmatic Programmer",
			Author:          "Andrew Hunt & David Thomas",
			ISBN:            "978-604-1-23456-8",
			Category:        "Technology",
			Description:     getStringPtr("Cuốn sách về tư duy và cách tiếp cận lập trình thực tế, hiệu quả cho các developer."),
			TotalCopies:     8,
			AvailableCopies: 8,
			PublishedYear:   getIntPtr(1999),
			Publisher:       getStringPtr("Addison-Wesley"),
		},
		{
			ID:              uuid.New().String(),
			Title:           "Súng, Vi Trùng và Thép",
			Author:          "Jared Diamond",
			ISBN:            "978-604-1-34567-9",
			Category:        "History",
			Description:     getStringPtr("Phân tích về sự phát triển khác biệt của các nền văn minh trên thế giới."),
			TotalCopies:     5,
			AvailableCopies: 5,
			PublishedYear:   getIntPtr(1997),
			Publisher:       getStringPtr("W. W. Norton & Company"),
		},
	}

	// Kiểm tra xem đã có dữ liệu chưa
	var count int64
	db.Model(&entity.Book{}).Count(&count)
	if count > 0 {
		logger.Info("seeder", "SeedBooks", fmt.Sprintf("Books already exist (%d records), skipping seed", count))
		return nil
	}

	// Tạo dữ liệu
	for _, book := range books {
		if err := db.Create(&book).Error; err != nil {
			logger.Error("seeder", "SeedBooks", fmt.Sprintf("Failed to create book %s: %v", book.Title, err))
			return err
		}
	}

	logger.Info("seeder", "SeedBooks", fmt.Sprintf("Successfully seeded %d books", len(books)))
	return nil
}

// Helper functions
func getStringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func getIntPtr(i int) *int {
	return &i
}

