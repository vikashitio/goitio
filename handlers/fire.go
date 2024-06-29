package handlers

import (
	"encoding/json"
	"fmt"
	"template/database"
	"template/fireblocks"
	"template/models"

	"github.com/gofiber/fiber/v2"
)

var privateKeyPath = "./fireblocks_secret.key"
var apiKey = "053c5036-525b-41fd-af32-2cf9776be07c" //API user: vikash API

// type ApiTokenProvider struct {
// 	privateKey *rsa.PrivateKey
// 	apiKey     string
// }

func UsersView(c *fiber.Ctx) error {

	// check session
	s, _ := store.Get(c)
	// Get value
	LoginMerchantID := s.Get("LoginMerchantID")
	if LoginMerchantID == nil {
		return c.Redirect("/login")
	}

	//privateKeyPath := "./fireblocks_secret.key"
	//apiKey := "053c5036-525b-41fd-af32-2cf9776be07c" //API user: vikash API
	tokenProvider, err := fireblocks.NewApiTokenProvider(privateKeyPath, apiKey)
	if err != nil {
		fmt.Printf("Error initializing token provider: %v\n", err)
		//return
	}

	// Example API calls
	// Ensure functions getAccountsPaged and createAccount are correctly defined to use this main function.

	path := "/v1/users"
	respBody, err := fireblocks.MakeAPIRequest("GET", path, nil, tokenProvider)
	if err != nil {
		return fmt.Errorf("error making GET request to accounts_paged: %w", err)
	}

	//var DataList = string(respBody)
	//fmt.Println("======", DataList)
	//////////////////////////

	// Parse the JSON data into the struct
	var fireblocksData []models.FireblocksUsers
	if err := json.Unmarshal([]byte(respBody), &fireblocksData); err != nil {
		fmt.Println(err)
	}

	///////////////////////////////////
	userProfileData, err := GetUserSessionData(c)
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Println(fireblocksData)
	return c.Render("fireblocks-users", fiber.Map{
		"Title":          "Fire Blocked User List",
		"Subtitle":       "User List",
		"Data":           fireblocksData,
		"CurrentSession": userProfileData.LoginMerchantName,
		"Sessions":       userProfileData.Sessions,
	})
}

func CreateVaultWallet(c *fiber.Ctx) error {

	VID := c.FormValue("VID")
	WID := c.FormValue("WID")

	//fmt.Println(VID, WID)

	// check session
	s, err := store.Get(c)
	if err != nil {
		panic(err)
	}
	// Get value

	//privateKeyPath := "./fireblocks_secret.key"
	//apiKey := "053c5036-525b-41fd-af32-2cf9776be07c"
	tokenProvider, err := fireblocks.NewApiTokenProvider(privateKeyPath, apiKey)
	if err != nil {
		fmt.Println("Error")
	}
	path := "/v1/vault/accounts/" + VID + "/" + WID
	//fmt.Println("======", path)
	respBody, err := fireblocks.MakeAPIRequest("POST", path, nil, tokenProvider)
	if err != nil {
		fmt.Println(err)
	}

	//var DataList = string(respBody)
	//fmt.Println("======", DataList)

	//FireblocksWallet
	var fireblocksData models.FireblocksWallet
	if err := json.Unmarshal([]byte(respBody), &fireblocksData); err != nil {
		fmt.Println(err)
	}

	//fmt.Println("======", fireblocksData.Message)
	if fireblocksData.Message == "" {
		s.Set("Alerts", "Wallet Creates Successfully with ID : "+fireblocksData.Id)
	} else {
		s.Set("Alerts", fireblocksData.Message)
	}

	if err := s.Save(); err != nil {
		panic(err)
	}
	//Alerts := s.Get("Alerts")
	//fmt.Println("==>Message :: ", Alerts)

	//////////////////////////
	return c.Redirect("/vault")

}

func WalletView(c *fiber.Ctx) error {

	VID := c.Params("VID")
	WID := c.Params("WID")

	//fmt.Println(VID, WID)

	s, _ := store.Get(c)
	// Get value
	LoginMerchantID := s.Get("LoginMerchantID")
	if LoginMerchantID == nil {
		return c.Redirect("/login")
	}

	Alerts := s.Get("Alerts")
	s.Delete("Alerts")

	//privateKeyPath := "./fireblocks_secret.key"
	//apiKey := "053c5036-525b-41fd-af32-2cf9776be07c"
	tokenProvider, err := fireblocks.NewApiTokenProvider(privateKeyPath, apiKey)
	if err != nil {
		fmt.Println("Error")
	}
	path := "/v1/vault/accounts/" + VID + "/" + WID + "/addresses_paginated"
	respBody, err := fireblocks.MakeAPIRequest("GET", path, nil, tokenProvider)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(string(respBody))
	// Parse the JSON data into the struct
	var fireblocksData models.FireblocksAddress
	if err := json.Unmarshal([]byte(respBody), &fireblocksData); err != nil {
		fmt.Println(err)
	}
	///////////////////////////////////
	userProfileData, err := GetUserSessionData(c)

	if err != nil {
		fmt.Println(err)
	}

	//fmt.Println(fireblocksData)
	//return c.Render("vault", fireblocksData)
	return c.Render("wallet", fiber.Map{
		"Title":          "Wallet Address",
		"Alert":          Alerts,
		"VoltID":         VID,
		"AssetID":        WID,
		"CurrentSession": userProfileData.LoginMerchantName,
		"Sessions":       userProfileData.Sessions,
		"Assets":         fireblocksData.Addresses,
	})
}

func VoltView(c *fiber.Ctx) error {

	s, _ := store.Get(c)

	// Get value
	LoginMerchantID := s.Get("LoginMerchantID")

	if LoginMerchantID == nil {
		return c.Redirect("/login")
	}

	Alerts := s.Get("Alerts")
	s.Delete("Alerts")

	//privateKeyPath := "./fireblocks_secret.key"
	//apiKey := "053c5036-525b-41fd-af32-2cf9776be07c"
	tokenProvider, err := fireblocks.NewApiTokenProvider(privateKeyPath, apiKey)
	if err != nil {
		fmt.Println("Error")
	}
	//#############################//
	var voltID = s.Get("LoginVoltID").(string)
	//fmt.Println("LoginVoltID -> ", voltID)
	var fireblocksData models.FireblocksResponse

	if voltID != "" {
		path := "/v1/vault/accounts/" + voltID //+ voltID
		//fmt.Println(path)
		respBody, err := fireblocks.MakeAPIRequest("GET", path, nil, tokenProvider)
		if err != nil {
			fmt.Println(err)
		}
		//fmt.Println(string(respBody))
		// Parse the JSON data into the struct
		//var fireblocksData models.FireblocksResponse
		if err := json.Unmarshal([]byte(respBody), &fireblocksData); err != nil {
			fmt.Println(err)
		}

	}
	///////////////////////////////////
	userProfileData, err := GetUserSessionData(c)

	if err != nil {
		fmt.Println(err)
	}

	//fmt.Println(fireblocksData)
	//return c.Render("vault", fireblocksData)
	return c.Render("vault", fiber.Map{
		"Title":          "Wallet List",
		"Subtitle":       "Wallet List",
		"Alert":          Alerts,
		"LoginVoltID":    voltID,
		"CurrentSession": userProfileData.LoginMerchantName,
		"Sessions":       userProfileData.Sessions,
		"ID":             fireblocksData.ID,
		"Name":           fireblocksData.Name,
		"HiddenOnUI":     fireblocksData.HiddenOnUI,
		"AutoFuel":       fireblocksData.AutoFuel,
		"Assets":         fireblocksData.Assets,
	})
}

func CreateVaultWalletAddress(c *fiber.Ctx) error {

	VID := c.Params("VID")
	WID := c.Params("WID")

	fmt.Println(VID, WID)

	s, _ := store.Get(c)
	// Get value
	LoginMerchantID := s.Get("LoginMerchantID")
	if LoginMerchantID == nil {
		return c.Redirect("/login")
	}

	//privateKeyPath := "./fireblocks_secret.key"
	//apiKey := "053c5036-525b-41fd-af32-2cf9776be07c"
	tokenProvider, err := fireblocks.NewApiTokenProvider(privateKeyPath, apiKey)
	if err != nil {
		fmt.Println("Error")
	}
	path := "/v1/vault/accounts/" + VID + "/" + WID + "/addresses"
	respBody, err := fireblocks.MakeAPIRequest("POST", path, nil, tokenProvider)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(respBody))
	// Parse the JSON data into the struct
	// var fireblocksData models.FireblocksAddress
	// if err := json.Unmarshal([]byte(respBody), &fireblocksData); err != nil {
	// 	fmt.Println(err)
	// }
	///////////////////////////////////
	Alerts := "Address Generated Sucessfully "
	s.Set("Alerts", Alerts)
	if err := s.Save(); err != nil {
		panic(err)
	}
	path = "/wallet/" + VID + "/" + WID
	return c.Redirect(path)
}
func CreateNewVault(c *fiber.Ctx) error {

	s, _ := store.Get(c)
	Alerts := "Account Generated Successfully"
	// Get value
	LoginMerchantID := s.Get("LoginMerchantID")
	LoginMerchantEmail := s.Get("LoginMerchantEmail").(string)
	if LoginMerchantID == nil {
		return c.Redirect("/login")
	}

	//privateKeyPath := "./fireblocks_secret.key"
	//apiKey := "053c5036-525b-41fd-af32-2cf9776be07c"
	tokenProvider, err := fireblocks.NewApiTokenProvider(privateKeyPath, apiKey)
	if err != nil {
		fmt.Println("Error")
	}

	path := "/v1/vault/accounts"
	Mydata := struct {
		Name string `json:"name"`
	}{
		Name: LoginMerchantEmail,
	}

	respBody, err := fireblocks.MakeAPIRequest("POST", path, Mydata, tokenProvider)
	if err != nil {
		fmt.Println(err)
		Alerts = "Account Not Generated"
	}
	fmt.Println(string(respBody))
	// Parse the JSON data into the struct
	var fireblocksData models.CreateVaultAccountResponse
	if err := json.Unmarshal([]byte(respBody), &fireblocksData); err != nil {
		fmt.Println(err)
	}

	s.Set("LoginVoltID", fireblocksData.ID)
	fmt.Println(fireblocksData.ID)
	///////////////////////////////////

	fmt.Println(fireblocksData.ID)
	if fireblocksData.ID != "" {

		Voltid := fireblocksData.ID
		LoginID := LoginMerchantID.(uint)
		result := database.DB.Db.Table("client_master").Save(&models.UpdateVolt{Client_id: LoginID, Volt_id: Voltid})

		if result.Error != nil {
			fmt.Println("ERROR in QUERY")
			Alerts = "Account Not Generated - 2"
		}

	}
	s.Set("Alerts", Alerts)
	if err := s.Save(); err != nil {
		//panic(err)
		fmt.Println("session not store on line no 560")
	}
	return c.Redirect("/vault")
}

func VaultAccountsView(c *fiber.Ctx) error {

	s, _ := store.Get(c)
	// Get value
	LoginMerchantID := s.Get("LoginMerchantID")
	if LoginMerchantID == nil {
		return c.Redirect("/login")
	}

	tokenProvider, err := fireblocks.NewApiTokenProvider(privateKeyPath, apiKey)
	if err != nil {
		fmt.Println("Error")
	}

	path := "/v1/vault/accounts_paged"

	respBody, err := fireblocks.MakeAPIRequest("GET", path, nil, tokenProvider)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(string(respBody))
	// Parse the JSON data into the struct
	var fireblocksData models.FireblocksVoltResponse
	if err := json.Unmarshal([]byte(respBody), &fireblocksData); err != nil {
		fmt.Println(err)
	}

	///////////////////////////////////
	userProfileData, err := GetUserSessionData(c)
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Println(fireblocksData)
	return c.Render("vault-accounts", fiber.Map{
		"Title":          "Volt List",
		"Subtitle":       "Volt List",
		"CurrentSession": userProfileData.LoginMerchantName,
		"Sessions":       userProfileData.Sessions,
		"Data":           fireblocksData.Accounts,
	})
}
func CreateVaultWalletView(c *fiber.Ctx) error {
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
