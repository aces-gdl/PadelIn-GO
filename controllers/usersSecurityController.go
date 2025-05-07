package controllers

import (
	"PadelIn/initializers"
	"PadelIn/models"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func getInstitutionByURL(myURL string) uint {
	u, err := url.Parse(myURL)
	if err != nil {
		panic(err)
	}
	u.Path = ""
	u.RawQuery = ""
	u.Fragment = ""

	var Origin models.Origin
	initializers.DB.Where("url =?", u.String()).First(&Origin)
	if Origin.ID == 0 {
		fmt.Println("Origin not found", u.String())
		return 0
	}
	fmt.Println(u)
	return Origin.InstitutionID
}
func Signup(c *gin.Context) {
	type bodyType struct {
		models.EndUser
		RequestURL string
	}

	var body bodyType

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Fallo al leer body...",
		})
		return
	}

	user := models.EndUser{
		Phone:         body.Phone,
		Name:          body.Name,
		InstitutionID: body.InstitutionID,
		Email:         body.Email,
		User:          body.Email,
		DateOfBirth:   time.Now(),
		Type:          "Player",
		Active:        "1",
	}
	user.InstitutionID = getInstitutionByURL(body.RequestURL)
	if user.InstitutionID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "URL de la institución no válida... ",
		})
		return
	}
	if body.Phone == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Faltan datos requeridos (telefono o correo)... ",
		})
		return
	}

	if body.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Fallo al convertir password a hash...",
			})
			return
		}
		user.Password = string(hash)
	}

	result := initializers.DB.Debug().Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Fallo al crear usuario... ",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Usuario creado correctamente", "user": user})
}

func Login(c *gin.Context) {
	type bodyType struct {
		models.EndUser
		RequestURL string
	}

	var body bodyType
	error := c.ShouldBindJSON(&body)
	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Fallo al leer body...",
		})
		return
	}

	var user models.EndUser
	user.InstitutionID = getInstitutionByURL(body.RequestURL)
	if user.InstitutionID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Institución no válida... ",
		})
		return
	}

	results := initializers.DB.Table("end_users").Debug().First(&user, "phone=? and type ='Player' and institution_id=?", body.Phone, user.InstitutionID)
	if results.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Telefono o clave invalido ...",
		})
		return
	}

	var institution models.Institution
	results = initializers.DB.Debug().First(&institution, user.InstitutionID)
	if results.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Institución invalida ...",
		})
		return
	}

	if results.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Telefono o clave invalido ...",
		})
	}
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Telefono o clave invalido ...",
		})
		return
	}

	if body.Password != "" {
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Correo o clave invalido ...",
			})
			return
		}

	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":       user.ID,
		"name":      user.Name,
		"user_type": user.Type,
		"inst":      user.InstitutionID,
		"exp":       time.Now().Add(time.Hour * 8).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Fallo al crear token ...",
		})
		return
	}
	user.Password = ""
	institution.Origin = ""

	c.JSON(http.StatusOK, gin.H{"message": "OK", "data": user, "token": tokenString, "institution": institution})
}

func Logout(c *gin.Context) {
	strHost := strings.Split(c.Request.Host, ":")[0]
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("club-jwt", "", -1, "/", strHost, true, true)
	c.JSON(http.StatusOK, gin.H{
		"message": "Logout exitoso",
	})
}
func Validate(c *gin.Context) {
	user, _ := c.Get("user")

	c.JSON(http.StatusOK, gin.H{
		"messaje": user,
	})
}
func ChangePassword(c *gin.Context) {
	var body struct {
		ID          uint
		OldPassword string
		NewPassword string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Fallo al leer body...",
		})
		return
	}

	// Fetch the user from the database
	var user models.EndUser
	result := initializers.DB.Table("user").First(&user, body.ID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Usuario no encontrado",
		})
		return
	}

	// Verify the old password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.OldPassword))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Contraseña actual incorrecta",
		})
		return
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error al generar nueva contraseña",
		})
		return
	}

	// Update the user's password in the database
	user.Password = string(hashedPassword)
	user.Active = "1"
	result = initializers.DB.Debug().Table("user").Where("id_user= ?", user.ID).Save(&user)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error al actualizar la contraseña",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Contraseña actualizada exitosamente",
	})
}
