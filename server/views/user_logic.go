package views

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/shumipro/meetapp/server/models"
)

func readBodyUser(body io.Reader) (models.User, error) {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return models.User{}, err
	}

	var user models.User
	if err := json.Unmarshal(data, &user); err != nil {
		return models.User{}, err
	}

	return user, nil
}
