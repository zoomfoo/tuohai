package im_api

import (
	"fmt"
	"time"

	"tuohai/im_api/models"
)

func (api *ImApi) friendLoop() {
	err := models.SyncFriends()
	fmt.Println(err)
	ticker := time.Tick(20 * time.Second)
	for {
		select {
		case <-ticker:
			fmt.Println(models.SyncFriends())
		case <-api.exitChan:
			goto exit
		}
	}

exit:
	fmt.Println("friendloop closing")
}
