package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"thirdtask/backend/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("secret_key")

var (
	users     = []models.User{}
	items     = []models.Item{}
	usersLock sync.Mutex
	itemsLock sync.Mutex
	userID    uint = 1
	itemID    uint = 1
)

// Создаём администратора при старте приложения
func init() {
	hash, _ := bcrypt.GenerateFromPassword([]byte("adminpass"), bcrypt.DefaultCost)
	admin := models.User{
		ID:       userID,
		Username: "admin",
		Password: string(hash),
		Role:     "admin",
	}
	userID++
	users = append(users, admin)
}

// ---------- AUTH ----------

func Register(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	if input.Username == "" || input.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password required"})
		return
	}

	usersLock.Lock()
	defer usersLock.Unlock()

	for _, u := range users {
		if u.Username == input.Username {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
			return
		}
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	newUser := models.User{
		ID:       userID,
		Username: input.Username,
		Password: string(hash),
		Role:     "user",
	}
	userID++
	users = append(users, newUser)

	c.JSON(http.StatusOK, gin.H{"message": "User registered"})
}

func Login(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var found *models.User
	usersLock.Lock()
	for i := range users {
		if users[i].Username == input.Username {
			found = &users[i]
			break
		}
	}
	usersLock.Unlock()

	if found == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(found.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	claims := jwt.MapClaims{
		"user_id": found.ID,
		"role":    found.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(jwtKey)

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// ---------- CRUD ----------

func GetItems(c *gin.Context) {
	itemsLock.Lock()
	defer itemsLock.Unlock()
	c.JSON(http.StatusOK, items)
}

// ---------- ADMIN-ONLY ----------
func CreateItem(c *gin.Context) {
	if !hasAdminRole(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	if input.Name == "" || input.Description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name and description required"})
		return
	}

	itemsLock.Lock()
	defer itemsLock.Unlock()

	newItem := models.Item{
		ID:          itemID,
		Name:        input.Name,
		Description: input.Description,
	}
	itemID++
	items = append(items, newItem)

	c.JSON(http.StatusOK, newItem)
}

func UpdateItem(c *gin.Context) {
	if !hasAdminRole(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	idParam := c.Param("id")
	id64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}
	id := uint(id64)

	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	itemsLock.Lock()
	defer itemsLock.Unlock()

	for i := range items {
		if items[i].ID == id {
			if input.Name != "" {
				items[i].Name = input.Name
			}
			if input.Description != "" {
				items[i].Description = input.Description
			}
			c.JSON(http.StatusOK, items[i])
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
}

func DeleteItem(c *gin.Context) {
	if !hasAdminRole(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	idParam := c.Param("id")
	id64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}
	id := uint(id64)

	itemsLock.Lock()
	defer itemsLock.Unlock()

	for i := range items {
		if items[i].ID == id {
			items = append(items[:i], items[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "Item deleted"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
}

// ---------- HELPERS ----------
func hasAdminRole(c *gin.Context) bool {
	roleVal, exists := c.Get("role")
	if !exists {
		return false
	}
	role := fmt.Sprintf("%v", roleVal)
	return role == "admin"
}

func JwtKey() []byte {
	return jwtKey
}
