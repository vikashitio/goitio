package handlers

import (
	"fmt"
	"template/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var store = session.New()

func IndexView(c *fiber.Ctx) error {

	s, _ := store.Get(c)

	// For check session
	//keys := s.Keys()
	//fmt.Println("Keys = > ", keys)

	// Get value
	LoginMerchantID := s.Get("LoginMerchantID")
	if LoginMerchantID == nil {
		return c.Redirect("/login")
	}

	userProfileData, _ := GetUserSessionData(c)

	return c.Render("index", fiber.Map{
		"Title":          "Dashboard",
		"Subtitle":       "Home",
		"CurrentSession": userProfileData.LoginMerchantName,
		"Sessions":       userProfileData.Sessions,
	})
}

func GetUserSessionData(c *fiber.Ctx) (*models.UserSession, error) {
	// Get current session

	s, err := store.Get(c)

	if err != nil {
		//U := &models.UserSession{Session: "Error"}
		//return U, nil
		//function.CheckSession()

		fmt.Print("Session Store")
	}

	// Get value
	LoginMerchantName := s.Get("LoginMerchantName").(string)
	LoginMerchantID := s.Get("LoginMerchantID").(uint)
	LoginMerchantEmail := s.Get("LoginMerchantEmail").(string)
	LoginMerchantStatus := s.Get("LoginMerchantStatus").(int)
	LoginIP := s.Get("LoginIP").(string)
	LoginTime := s.Get("LoginTime").(string)
	LoginAgent := s.Get("LoginAgent").(string)
	// If there is a valid session
	if len(s.Keys()) > 0 {

		// Get profile info
		U := &models.UserSession{
			LoginMerchantName:   LoginMerchantName,
			LoginMerchantID:     LoginMerchantID,
			LoginMerchantEmail:  LoginMerchantEmail,
			LoginMerchantStatus: LoginMerchantStatus,
			Session:             "Test",
		}

		// Append session
		U.Sessions = append(
			U.Sessions,
			models.UserSessionOther{
				LoginIP:    LoginIP,
				LoginTime:  LoginTime,
				LoginAgent: LoginAgent,
			},
		)
		//fmt.Println(U)
		return U, nil
	}

	return nil, nil
}
