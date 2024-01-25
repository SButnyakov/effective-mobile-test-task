package external

const AgifyURL = "https://api.agify.io/?name="

type AgifyResponse struct {
	Age uint8 `json:"age"`
}

func (a APIs) GetAge(name string) (uint8, error) {
	requestURL := AgifyURL + name

	var response AgifyResponse

	err := a.get(requestURL, &response)
	if err != nil {
		return 0, err
	}

	return response.Age, nil
}
