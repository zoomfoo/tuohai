package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/Sirupsen/logrus"
	"gopkg.in/gin-gonic/gin.v1"
)

func SignIn(url string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ""
		url = fmt.Sprintf("%s?session_token=%s", url, token)
		if ok, err := ValidationToken(url, "", token); err != nil || !ok {
			ctx.Abort()
			ctx.JSON(http.StatusUnauthorized, gin.H{"result": "no"})
		} else {
			ctx.Next()
		}
	}
}

func ValidationToken(url, user_id, token string) (bool, error) {
	uid, err := strconv.ParseFloat(user_id, 64)
	if err != nil {
		return false, fmt.Errorf("ERROR: %s", "Invalid user_id type")
	}
	if url == "" {
		return false, fmt.Errorf("%v", "url is empty")
	}

	body, err := getTokenInfo(url)
	if err != nil {
		return false, err
	}

	m, err := serialized(body)
	if err != nil {
		return false, err
	}

	if m["error_code"].(float64) != 0 {
		return false, fmt.Errorf("ERROR: %s", "Invalid token")
	}

	profile, ok := m["profile"].(map[string]interface{})
	if !ok {
		return false, fmt.Errorf("ERROR: %s", "Return to the profile of the json format is not correct")
	}

	imid, ok := profile["imid"].(float64)
	if !ok {
		return false, fmt.Errorf("ERROR: %s", "The types of imid not float64 returns")
	}

	return imid == uid, nil
}

func getTokenInfo(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func serialized(data []byte) (m map[string]interface{}, err error) {
	err = json.Unmarshal(data, &m)
	return
}
