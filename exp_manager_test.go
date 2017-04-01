package explib
import (
	"fmt"
	"testing"
)

func TestExplib (t *testing.T) {
	var exp_manager ExpManager
	exp_manager.Load_default_from_path("/home/wangxiaodan1")
	exp_manager.Load_exp_from_path("/home/wangxiaodan1")
	var request ExpRequestInfo
	request.Androidid = "android_test"
	res := exp_manager.Get_assign(request)
	for _, exp_id := range res.Exp_ids {
		fmt.Println("EXP_ID:", exp_id)
	}
	for k, v := range res.Param_infos {
		fmt.Println(k,v)
	}
}
