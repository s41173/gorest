package utils

import (
	"bytes"
	"fmt"
	"go-rest/config"
	"go-rest/models"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

// Fungsi sendRequest mengirim permintaan HTTP dan mengembalikan hasilnya
func sendRequest(param []byte) (interface{}, error) {
	// Ambil data properti dari database
	var property models.Property
	data, err := property.GetProperty(config.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to get property data: %v", err)
	}

	// Mendapatkan apiKey dan url dari data properti, pastikan tipe string
	apiKeyPtr, apiKeyOk := data["notif_token"].(*string)
	urlPtr, urlOk := data["notif_url"].(*string)

	if !apiKeyOk || !urlOk || apiKeyPtr == nil || urlPtr == nil {
		return nil, fmt.Errorf("notif_token or notif_url is not a string or is nil")
	}

	apiKey := *apiKeyPtr
	url := *urlPtr + "notif/add_notif"

	// Mencetak body request yang akan dikirim
	// fmt.Printf("Request Body: %s\n", param)

	// Membuat client HTTP dengan batas waktu
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Membuat request POST
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(param))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Menetapkan header untuk request
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-Auth-Token", apiKey)

	// Melakukan request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Membaca response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	// Mencetak body dari response
	// fmt.Printf("Response Body: %s\n", body)

	// Mengembalikan sesuai tipe yang diminta
	result := map[string]interface{}{
		"body":      string(body),
		"http_code": resp.StatusCode,
	}
	fmt.Printf("HTTP Code: %d\n", resp.StatusCode)

	return result, nil
}

// Fungsi postNotif
func postNotif(typeVar, sentto, cust, custname, subject, content, modul string) bool {
	// Menyiapkan data POST dalam format URL-encoded
	postData := url.Values{}
	postData.Set("type", typeVar)
	postData.Set("customer", cust)
	postData.Set("custname", custname)
	postData.Set("subject", subject)
	postData.Set("sentto", sentto)
	postData.Set("content", content)
	postData.Set("modul", modul)

	// Mengonversi map ke format x-www-form-urlencoded
	postString := postData.Encode()

	// Mengirim permintaan ke URL notifikasi dengan mengubah postString ke []byte
	response, err := sendRequest([]byte(postString))
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}

	// Mengecek respons jika responseType bernilai true
	respData, ok := response.(map[string]interface{})
	if ok && respData["http_code"].(int) == 200 {
		return true
	}
	return false
}

// Fungsi SendNotif
func SendNotif(c *gin.Context, notifyType, custID, subject, content, module string) bool {
	customer := models.Customer{}
	if err := config.DB.First(&customer, custID).Error; err == nil {
		var res, res1, res2 bool
		switch notifyType {
		case "0":
			res = postNotif(notifyType, customer.Email, custID, customer.FirstName, subject, content, module)
		case "1":
			res = postNotif(notifyType, customer.Phone1, custID, customer.FirstName, subject, content, module)
		case "2":
			res1 = postNotif(notifyType, customer.Email, custID, customer.FirstName, subject, content, module)
			res2 = postNotif(notifyType, customer.Phone1, custID, customer.FirstName, subject, content, module)
			res = res1 && res2
		case "7":
			res = postNotif(notifyType, customer.Phone1, custID, customer.FirstName, subject, content, module)
		case "8":
			res1 = postNotif(notifyType, customer.Email, custID, customer.FirstName, subject, content, module)
			res2 = postNotif(notifyType, customer.Phone1, custID, customer.FirstName, subject, content, module)
			res = res1 && res2
		}
		return res
	}
	return false
}
