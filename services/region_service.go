package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Province struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type City struct {
	ID         string `json:"id"`
	ProvinceID string `json:"province_id"`
	Name       string `json:"name"`
}

const (
	BASE_URL = "https://www.emsifa.com/api-wilayah-indonesia/api"
)

// GetProvinces fetches all provinces from the Indonesian region API
func GetProvinces() ([]Province, error) {
	url := fmt.Sprintf("%s/provinces.json", BASE_URL)
	
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch provinces: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var provinces []Province
	if err := json.Unmarshal(body, &provinces); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	return provinces, nil
}

// GetCitiesByProvinceID fetches cities for a specific province
func GetCitiesByProvinceID(provinceID string) ([]City, error) {
	url := fmt.Sprintf("%s/regencies/%s.json", BASE_URL, provinceID)
	
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch cities: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var cities []City
	if err := json.Unmarshal(body, &cities); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	return cities, nil
}

// GetAllCities fetches all cities from all provinces
func GetAllCities() ([]City, error) {
	provinces, err := GetProvinces()
	if err != nil {
		return nil, err
	}

	var allCities []City
	for _, province := range provinces {
		cities, err := GetCitiesByProvinceID(province.ID)
		if err != nil {
			// Log error but continue with other provinces
			fmt.Printf("Error fetching cities for province %s: %v\n", province.Name, err)
			continue
		}
		
		// Add province_id to each city
		for i := range cities {
			cities[i].ProvinceID = province.ID
		}
		
		allCities = append(allCities, cities...)
	}

	return allCities, nil
}

// GetProvinceByID fetches a specific province by ID
func GetProvinceByID(provinceID string) (*Province, error) {
	provinces, err := GetProvinces()
	if err != nil {
		return nil, err
	}

	for _, province := range provinces {
		if province.ID == provinceID {
			return &province, nil
		}
	}

	return nil, fmt.Errorf("province with ID %s not found", provinceID)
}

// GetCityByID fetches a specific city by ID
func GetCityByID(cityID string) (*City, error) {
	allCities, err := GetAllCities()
	if err != nil {
		return nil, err
	}

	for _, city := range allCities {
		if city.ID == cityID {
			return &city, nil
		}
	}

	return nil, fmt.Errorf("city with ID %s not found", cityID)
}