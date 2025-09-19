package testhelper

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	jwthelper "github.com/Manuel-Leleuly/kanban-flow-go/helpers/jwt"
	"github.com/Manuel-Leleuly/kanban-flow-go/models"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

var TEST_USER models.User = models.User{
	ID:        "6f213f7d399c482eaace1bcd5a35b9bd",
	FirstName: "Test",
	LastName:  "User",
	Email:     "testuser@example.com",
	Password:  "testing123",
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
}

func ConnectToTestDB(d *models.DBInstance) error {
	if err := godotenv.Load("../../.env"); err != nil {
		return err
	}

	if err := d.ConnectToDB(os.Getenv("DB_TEST_NAME")); err != nil {
		return err
	}

	if err := d.SyncDatabase(); err != nil {
		return err
	}

	return nil
}

// user
func GetTestToken(d *models.DBInstance) (*models.Token, error) {
	loginReqBody := models.Login{
		Email:    TEST_USER.Email,
		Password: TEST_USER.Password,
	}

	var user models.User
	result := d.DB.Where("email = ?", loginReqBody.Email).First(&user)
	if result.Error != nil || user.ID == "" {
		return nil, result.Error
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReqBody.Password)); err != nil {
		return nil, err
	}

	accessTokenString, err := jwthelper.CreateAccessToken(user)
	if err != nil {
		return nil, err
	}

	return &models.Token{
		Status:      "success",
		AccessToken: accessTokenString,
	}, nil
}

func CreateTestUser(d *models.DBInstance) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(TEST_USER.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	newUser := models.User{
		ID:        TEST_USER.ID,
		FirstName: TEST_USER.FirstName,
		LastName:  TEST_USER.LastName,
		Email:     TEST_USER.Email,
		Password:  string(hashedPassword),
	}
	if err := d.DB.Create(&newUser).Error; err != nil {
		return err
	}

	return nil
}

func DeleteAllTestUsers(d *models.DBInstance) error {
	var users []models.User
	if err := d.DB.Raw("DELETE from users").Scan(&users).Error; err != nil {
		return err
	}
	return nil
}

// ticket
var TEST_TICKET models.Ticket = models.Ticket{
	ID:          "7fa00bcc3bc94bada4992d321e94528a",
	Title:       "Test Ticket",
	Description: "Test Ticket Description",
	Assignees:   []string{"frontend", "backend"},
	Status:      "todo",
	CreatedAt:   time.Now(),
	UpdatedAt:   time.Now(),
	User:        TEST_USER,
}

func CreateTestTicket(d *models.DBInstance) error {
	newTicket := models.Ticket{
		ID:          TEST_TICKET.ID,
		Title:       TEST_TICKET.Title,
		Description: TEST_TICKET.Description,
		Assignees:   TEST_TICKET.Assignees,
		Status:      TEST_TICKET.Status,
		User:        TEST_USER,
	}

	if err := d.DB.Create(&newTicket).Error; err != nil {
		return err
	}

	return nil
}

func DeleteAllTestTickets(d *models.DBInstance) error {
	var tickets []models.Ticket
	if err := d.DB.Raw("TRUNCATE tickets").Scan(&tickets).Error; err != nil {
		return err
	}
	return nil
}

func GetHTTPRequest(method string, path string, body io.Reader, accessToken string) *http.Request {
	request := httptest.NewRequest(method, path, body)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")

	if accessToken != "" {
		request.Header.Add("Authorization", "Bearer "+accessToken)
	}

	return request
}
