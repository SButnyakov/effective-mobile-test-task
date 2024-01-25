package external

import "fmt"

const NationalizeURL = "https://api.nationalize.io/?name="

type NationalizeResponse struct {
	Country []CountryResponseObject `json:"country"`
}

type CountryResponseObject struct {
	CountryID   string  `json:"country_id"`
	Probability float32 `json:"probability"`
}

func (a APIs) GetMostProbableCountry(name string) (string, error) {
	requestURL := NationalizeURL + name

	var response NationalizeResponse

	err := a.get(requestURL, &response)
	if err != nil {
		return "", err
	}

	if len(response.Country) == 0 {
		return "", fmt.Errorf("no country found")
	}

	mostProbable := response.Country[0]

	for i := 1; i < len(response.Country); i++ {
		if response.Country[i].Probability > mostProbable.Probability {
			mostProbable = response.Country[i]
		}
	}

	return mostProbable.CountryID, nil
}
