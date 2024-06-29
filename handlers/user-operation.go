package handlers

import (
	"fmt"
	"template/database"
	"template/function"
	"template/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func LoginView(c *fiber.Ctx) error {
	return c.Render("login", fiber.Map{
		"Title": "Login Form",
		"Alert": "",
	})
}

func LoginPost(c *fiber.Ctx) error {
	// Parses the request body
	getUserName := c.FormValue("username")
	getPassword := c.FormValue("password")

	//fmt.Println(getUserName, getPassword)
	Alerts := ""
	loginList := models.LoginList{}
	result := database.DB.Db.Table("client_master").Where("username = ?", getUserName).Find(&loginList)

	//fmt.Println(loginList.Status)

	if result.Error != nil {
		//fmt.Println("ERROR in QUERY")
		Alerts = "ERROR in QUERY"
	}

	if result.RowsAffected == 1 {
		//fmt.Println(loginList)
		//fmt.Println(Full_name)

		if loginList.Status != 1 {
			//fmt.Println("Account Not Activate / Deleted")
		} else if loginList.Password != "" {
			//fmt.Println(loginList.Password)
			err := bcrypt.CompareHashAndPassword([]byte(loginList.Password), []byte(getPassword))
			if err == nil {
				//fmt.Println("You have successfully logged in")

				s, _ := store.Get(c)

				//s.Set("name", "john")

				// Set key/value
				loginIp := c.Context().RemoteIP().String()
				s.Set("LoginMerchantName", loginList.Full_name)
				s.Set("LoginMerchantID", loginList.Client_id)
				s.Set("LoginMerchantEmail", getUserName)
				s.Set("LoginMerchantStatus", loginList.Status)
				s.Set("LoginVoltID", loginList.Volt_id)
				s.Set("LoginIP", c.Context().RemoteIP().String())
				s.Set("LoginTime", time.Unix(time.Now().Unix(), 0).UTC().String())
				s.Set("LoginAgent", string(c.Request().Header.UserAgent()))

				//Save sessions
				if err := s.Save(); err != nil {
					panic(err)
				}

				qry := models.LoginHistory{Client_id: loginList.Client_id, Login_ip: loginIp}
				result := database.DB.Db.Table("login_history").Select("client_id", "login_ip").Create(&qry)
				if result.Error != nil {
					fmt.Println(result.Error)
				}

				return c.Redirect("/")

			} else {
				//fmt.Println("Wrong Password")
				Alerts = "Wrong Password"
			}

		}

	} else {
		//fmt.Println("Account Not Found")
		Alerts = "Account Not Found"

	}

	return c.Render("login", fiber.Map{
		"Title": "Login Form",
		"Alert": Alerts,
		//"Facts":    facts,
	})
}

func RegistrationView(c *fiber.Ctx) error {
	//facts := []models.Fact{}
	//fmt.Println("===>", facts)
	return c.Render("registration", fiber.Map{
		"Title": "Registration Form",
		"Alert": "",
		//"Facts":    facts,
	})
}

func RegistrationPost(c *fiber.Ctx) error {
	// Parses the request body
	getName := c.FormValue("name")
	getEmail := c.FormValue("email")

	//fmt.Println(getName, getEmail)

	// Find Duplicate Email in DB

	Alerts := ""
	loginList := models.LoginList{}
	result := database.DB.Db.Table("client_master").Where("username = ?", getEmail).Find(&loginList)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	receivedId := loginList.Client_id
	fmt.Println("XXX ", receivedId)

	if receivedId == 0 {

		// END Find Duplicate Email in DB

		var password = function.PasswordGenerator(8)
		//fmt.Println(password)

		var hash []byte
		// func GenerateFromPassword(password []byte, cost int) ([]byte, error)
		hash, _ = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		//fmt.Println(hash)
		//loginList := models.Client_Master{}
		//result := database.DB.Db.Table("client_master").Create(&loginList)

		qry := models.Client_Master{Username: getEmail, Password: string(hash), Full_name: getName, Status: 1}
		//result := database.DB.Db.Table("client_master").Create(&[]models.Client_Master{{Username: getEmail, Password: string(hash), Full_name: getName, Status: 1}})
		result = database.DB.Db.Table("client_master").Select("username", "full_name", "password", "status").Create(&qry)
		//fmt.Println(result)

		if result.Error != nil {
			fmt.Println(result.Error)
		} else {
			fmt.Println(result.RowsAffected)
			fmt.Println(qry.Client_id)
			ClientID := qry.Client_id

			MyData := struct {
				Client_id uint `json:"name"`
			}{
				Client_id: ClientID,
			}
			result = database.DB.Db.Table("client_details").Select("client_id").Create(&MyData)
			if result.Error != nil {
				fmt.Println(result.Error)
			}
			//  Email///
			//var domName = "http://localhost:8080"
			var subject = "Test Message"
			//var HTMLbody = "Hi this is message"
			HTMLbody :=
				`<html>
			<p><strong>Hello , ` + getName + `</strong></p>
			<br/>
			<p>Welcome to Golang Bank! We are pleased to inform that your account has been created.</p>
			<br/>
			<strong>Login Details for Your Account:<br/>=====================<br/><strong>
			<p>Username :  ` + getEmail + `</p>
			<p>Password :  ` + password + `</p>
			
			<br/>
			Cheers,
			<br/>
            <strong>Golang Bank</strong>
		</html>`
			err := function.SendEmail(subject, HTMLbody)
			if err != nil {
				fmt.Println("issue sending verification email")
			} else {
				fmt.Println("Mail Going")
			}

			s, _ := store.Get(c)

			//s.Set("name", "john")

			// Set key/value
			s.Set("LoginMerchantName", getName)
			s.Set("LoginMerchantID", ClientID)
			s.Set("LoginMerchantEmail", getEmail)
			s.Set("LoginVoltID", "")
			s.Set("LoginMerchantStatus", 1)
			s.Set("LoginIP", c.Context().RemoteIP().String())
			s.Set("LoginTime", time.Unix(time.Now().Unix(), 0).UTC().String())
			s.Set("LoginAgent", string(c.Request().Header.UserAgent()))

			if err := s.Save(); err != nil {
				panic(err)
			}

			return c.Redirect("/")

		}
	} else {
		//fmt.Println("Duplicate = ", loginList.Client_id)
		Alerts = "Duplicate Email ID"

	}

	return c.Render("registration", fiber.Map{
		"Title": "Registration Form",
		"Alert": Alerts,
		//"Facts":    facts,
	})
}

// For Login History

func Loginhistory(c *fiber.Ctx) error {

	// check session
	s, _ := store.Get(c)
	// Get value
	LoginMerchantID := s.Get("LoginMerchantID")
	if LoginMerchantID == nil {
		return c.Redirect("/login")
	}

	loginHistory := []models.LoginHistory{}
	database.DB.Db.Table("login_history").Order("token_id desc").Where("client_id = ?", LoginMerchantID).Find(&loginHistory)
	//.Select("login_time")
	userProfileData, err := GetUserSessionData(c)
	if err != nil {
		panic(err)
	}

	//fmt.Println(loginHistory)
	return c.Render("login-history", fiber.Map{
		"Title":          "Login History",
		"Subtitle":       "Login History",
		"LoginHistory":   loginHistory,
		"CurrentSession": userProfileData.LoginMerchantName,
		"Sessions":       userProfileData.Sessions,
	})
}
func CreateVaultWalletView99(c *fiber.Ctx) error {
	VID := c.Params("VID")

	// check session
	s, _ := store.Get(c)
	// Get value
	LoginMerchantID := s.Get("LoginMerchantID")
	Alerts := s.Get("Alerts")
	s.Delete("Alerts")
	if err := s.Save(); err != nil {
		panic(err)
	}
	if LoginMerchantID == nil {
		return c.Redirect("/login")
	}

	userProfileData, err := GetUserSessionData(c)
	if err != nil {
		panic(err)
	}

	//fmt.Println(profile)
	//fmt.Println(userProfileData)
	return c.Render("create-wallet", fiber.Map{
		"Title":          "Create Wallet",
		"Subtitle":       "Create Wallet",
		"Alert":          Alerts,
		"VoltID":         VID,
		"CurrentSession": userProfileData.LoginMerchantName,
		"Sessions":       userProfileData.Sessions,
	})
}

func ProfileView(c *fiber.Ctx) error {

	// check session
	s, _ := store.Get(c)
	// Get value
	LoginMerchantID := s.Get("LoginMerchantID")
	Alerts := s.Get("Alerts")
	s.Delete("Alerts")
	if err := s.Save(); err != nil {
		panic(err)
	}
	if LoginMerchantID == nil {
		return c.Redirect("/login")
	}

	profile := []models.Profile{}
	database.DB.Db.Table("client_details").Where("client_id = ?", LoginMerchantID).Find(&profile)
	//.Select("login_time")
	userProfileData, err := GetUserSessionData(c)
	if err != nil {
		panic(err)
	}

	//fmt.Println(profile)
	//fmt.Println(userProfileData)
	return c.Render("profile", fiber.Map{
		"Title":          "Profile",
		"Subtitle":       "Profile",
		"Alert":          Alerts,
		"Profile":        profile,
		"CurrentSession": userProfileData.LoginMerchantName,
		"Sessions":       userProfileData.Sessions,
	})
}

func ProfilePost(c *fiber.Ctx) error {
	// Parses the request body
	getGender := c.FormValue("gender")
	getBirthDate := c.FormValue("birth_date")
	getCountryCode := c.FormValue("country_code")
	getMobile := c.FormValue("mobile")
	getAddressLine1 := c.FormValue("address_line1")
	getAddressLine2 := c.FormValue("address_line2")

	s, _ := store.Get(c)
	LoginMerchantID := s.Get("LoginMerchantID").(uint)
	//fmt.Println(getGender, getBirthDate, getCountryCode, getMobile, getAddressLine1, getAddressLine2)

	result := database.DB.Db.Table("client_details").Save(&models.Profile{Client_id: LoginMerchantID, Gender: getGender, BirthDate: getBirthDate, CountryCode: getCountryCode, Mobile: getMobile, AddressLine1: getAddressLine1, AddressLine2: getAddressLine2})

	//fmt.Println(loginList.Status)
	Alerts := "Profile Updated successfully"
	if result.Error != nil {
		//fmt.Println("ERROR in QUERY")
		Alerts = "Profile Not Updated"
	}

	// check session

	s.Set("Alerts", Alerts)
	if err := s.Save(); err != nil {
		panic(err)
	}

	return c.Redirect("/profile")

}

// Create new Fact View handler
func LogOut(c *fiber.Ctx) error {
	s, err := store.Get(c)
	if err != nil {
		panic(err)
	}

	s.Delete("LoginMerchantID")
	s.Delete("LoginVoltID")

	// Destroy session
	if err := s.Destroy(); err != nil {
		panic(err)
	}

	return c.Redirect("/login")
}
