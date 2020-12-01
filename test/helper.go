package test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"temperature/config"
	"temperature/models"
	"log"
	"net/http"
	"testing"
)

var theFixtures *fixtures

func Init() {
	err := config.Init("../test")
	if err != nil {
		log.Fatal(err)
	}
	//theFixtures, err = newFixtures()
	//if err != nil {
	//	log.Fatal(err)
	//}
	err = models.InitForTest()
	if err != nil {
		log.Fatal(err)
	}
	return
}

/* options:
fixtures, bool, default: true
*/
func Prepare(t *testing.T, options ...map[string]interface{}) *assert.Assertions {
	prepareFixtures := true
	if len(options) > 0 {
		if v, ok := options[0]["fixtures"]; ok {
			prepareFixtures = v.(bool)
		}
	}

	if prepareFixtures {
		theFixtures.Prepare(t)
	}
	return assert.New(t)
}

func Authorization(r *http.Request, openIDOption ...int) (err error) {
	devId := 1472
	if len(openIDOption) > 0 {
		devId = openIDOption[0]
	}
	info, err := new(models.Device).FindDevByID(devId)
	if err != nil {
		return err
	} else if info == nil {
		return errors.New("failed to find dev for devId: %d" + string(devId))
	}
	//session, err := models.NewSession(int(info.ID))
	//if err != nil {
	//	return
	//} else if session == nil {
	//	return errors.New("failed to create session for devId: " + string(devId))
	//}
	//r.Header.Add("Authorization", "Bearer "+session.Token.String())
	return
}
