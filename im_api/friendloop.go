package im_api

import (
	"fmt"
	"time"
)

func (api *ImApi) friendLoop() {
	ticker := time.Tick(2 * time.Second)
	for {
		select {
		case <-ticker:
			fmt.Println("同步好友关系ing")
		case <-api.exitChan:
			goto exit
		}
	}

exit:
	fmt.Println("friendloop closing")
}
