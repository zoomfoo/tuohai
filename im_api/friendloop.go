package im_api

import (
	"fmt"
	"time"

	"tuohai/im_api/models"
)

func (api *ImApi) friendLoop() {
	ticker := time.Tick(1 * time.Minute)
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
