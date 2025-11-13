package db

type User struct {
    ID         int64
    TelegramID int64
    Role       string
    Username   string
    FirstName  string
    LastName   string
    Email      string
    Phone      string
    CreatedAt  string
}
