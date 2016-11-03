package http

import (
	"encoding/json"
	"proxy_service/internal/options"
	"proxy_service/internal/pb/MsgSrv"
	"proxy_service/internal/util"
	"testing"
)

func init() {
	util.ConfPath = "../../conf/web.config"
	options.NewOptions()
}

func TestPost(t *testing.T) {
	// t.Log(post("120.27.45.244:12138/change_member_notify", map[string]interface{}{
	// 	"user_id":          216,
	// 	"change_type":      1,
	// 	"group_id":         1076,
	// 	"cur_user_id_list": []int{1, 2, 3, 4},
	// 	"chg_user_id_list": []int{5, 6},
	// }))
}

func TestChangeMemberNotify(t *testing.T) {
	t.Log(ChangeMemberNotify(map[string]interface{}{
		"user_id":          216,
		"change_type":      1,
		"group_id":         1076,
		"cur_user_id_list": []int{1, 2, 3, 4},
		"chg_user_id_list": []int{5, 6},
	}))
}

type IMGroupChangeMemberNotify struct {
	// cmd id:       0x040b
	UserId        string   `protobuf:"bytes,1,opt,name=user_id,json=userId" json:"user_id,omitempty"`
	ChangeType    int      `protobuf:"varint,2,opt,name=change_type,json=changeType,enum=IM.BaseDefine.GroupModifyType" json:"change_type,omitempty"`
	GroupId       string   `protobuf:"bytes,3,opt,name=group_id,json=groupId" json:"group_id,omitempty"`
	CurUserIdList []string `protobuf:"bytes,4,rep,name=cur_user_id_list,json=curUserIdList" json:"cur_user_id_list,omitempty"`
	ChgUserIdList []string `protobuf:"bytes,5,rep,name=chg_user_id_list,json=chgUserIdList" json:"chg_user_id_list,omitempty"`
	RcvId         string   `protobuf:"bytes,6,opt,name=rcv_id,json=rcvId" json:"rcv_id,omitempty"`
}

func TestDispatchNotify(t *testing.T) {
	p := &IMGroupChangeMemberNotify{
		UserId:        "10",
		ChangeType:    1,
		GroupId:       group_id,
		CurUserIdList: []string{"10", "11", "12", "13"},
		ChgUserIdList: []string{"15", "16", "17"},
	}

	js, err := json.Marshal(p)
	if err != nil {
		t.Error(err)
		return
	}

	dispatchNotify("127.0.0.1:5004", &MsgSrv.SendNotifyReq{
		NotifyType: MsgSrv.NotifyType_NOTIFY_TYPE_GROUP_NOTIFY,
		Content:    string(js),
	})
}
