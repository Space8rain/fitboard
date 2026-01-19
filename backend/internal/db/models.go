package db

type User struct {
    ID         int64  `json:"id"`
    TelegramID int64  `json:"telegram_id"`
    Role       string `json:"role"`
    Username   string `json:"username"`
    FirstName  string `json:"first_name"`
    LastName   string `json:"last_name"`
    Email      string `json:"email"`
    Phone      string `json:"phone"`
    CreatedAt  string `json:"created_at"`
}
