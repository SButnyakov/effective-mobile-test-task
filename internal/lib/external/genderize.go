package external

const GenderizeURL = "https://api.genderize.io/?name="

type GenderizeResponse struct {
	Gender string `json:"gender" validate:"oneof=male female"`
}

func (a APIs) GetGender(name string) (string, error) {
	requestURL := GenderizeURL + name

	var response GenderizeResponse
	err := a.get(requestURL, &response)
	if err != nil {
		return "", err
	}

	return response.Gender, nil
}
