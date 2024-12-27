package calendar

import (
	"backend/internal/db"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type DateOnly struct {
	time.Time
}

type CalendarDay struct {
	ID       int      `json:"id"`
	UserID   int      `json:"user_id" binding:"required"`
	DayDate  DateOnly `json:"day_date" binding:"required"`
	Title    string   `json:"title" binding:"required"`
	ImageURL string   `json:"image_url" binding:"required,url"`
	Content  string   `json:"content" binding:"required"`
}

func (d *DateOnly) UnmarshalJSON(b []byte) error {
	s := string(b)
	t, err := time.Parse(`"2006-01-02"`, s)
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}

func (d *DateOnly) Scan(value interface{}) error {
	if value == nil {
		*d = DateOnly{Time: time.Time{}}
		return nil
	}

	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("неподдерживаемый тип %T для DateOnly", value)
	}

	t, err := time.Parse("2006-01-02", string(b))
	if err != nil {
		return err
	}

	*d = DateOnly{Time: t}
	return nil
}

func CreateDay(c *gin.Context) {
	var day CalendarDay
	if err := c.ShouldBindJSON(&day); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dayDateStr := day.DayDate.Format("2006-01-02")

	_, err := db.DB.Exec("INSERT INTO calendar_days (user_id, day_date, title, image_url, content) VALUES (?, ?, ?, ?, ?)", day.UserID, dayDateStr, day.Title, day.ImageURL, day.Content)
	if err != nil {
		log.Printf("Ошибка при создании дня: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании дня"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "День успешно создан"})
}

func GetNextDay(c *gin.Context) {
	var day CalendarDay
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}

	userIDInt, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обработке ID пользователя"})
		return
	}

	log.Printf("Извлеченный userID: %d", userIDInt)

	err := db.DB.QueryRow("SELECT id, user_id, day_date, title, image_url, content FROM calendar_days WHERE user_id = ? AND day_date = CURDATE() ORDER BY day_date ASC LIMIT 1",
		userIDInt).Scan(&day.ID, &day.UserID, &day.DayDate, &day.Title, &day.ImageURL, &day.Content)
	if err != nil {
		log.Printf("Ошибка при получении дня: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Нет доступных дней"})
		return
	}

	c.JSON(http.StatusOK, day)
}
