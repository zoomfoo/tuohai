package im_api

import (
	"fmt"
	"time"

	"tuohai/im_api/models"
)

func (api *ImApi) friendLoop() {
	models.SyncFriends()
	ticker := time.Tick(300 * time.Second)
	for {
		select {
		case <-ticker:
			models.SyncFriends()
		case <-api.exitChan:
			goto exit
		}
	}

exit:
	fmt.Println("friendloop closing")
}
