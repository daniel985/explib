package explib

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	//"encoding/json"
	json "github.com/bitly/go-simplejson"
	farmhash "github.com/leemcloughlin/gofarmhash"
)

const (
	DEFAULT_SUFFIX = "default"
	EXP_SUFFIX     = "exp"
	SINGLE         = "single"
	MULTIPLE       = "multiple"
	MODE_BASE      = 100
)

/**
 * exp manager
 */
type ExpManager struct {
	default_params    map[string]string
	single_layer_list []LayerInfo
	multi_layer_list  []LayerInfo
}

/**
 * load default param from one json file
 */
func (exp_manager *ExpManager) load_single_default(filename string) {
	json_conf, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalln("Open default_param_conf failed!")
	}
	js, err := json.NewJson([]byte(json_conf))
	if err != nil {
		log.Fatalln("Parse default_param_conf failed!")
	}
	//exp_manager.default_params = make(map[string]string)
	json_params, err := js.Map()
	if err != nil {
		log.Fatalln("Parse param_map failed!")
	}
	for k, v := range json_params {
		exp_manager.default_params[k] = v.(string)
	}
}

/**
 * load default param from one path
 */
func (exp_manager *ExpManager) Load_default_from_path(pathname string) {
	exp_manager.default_params = make(map[string]string)
	dir, err := ioutil.ReadDir(pathname)
	if err != nil {
		log.Fatalln("Read default_param_path failed!")
	}
	PathSep := string(os.PathSeparator)
	for _, file := range dir {
		if file.IsDir() {
			continue
		}
		if strings.HasSuffix(file.Name(), DEFAULT_SUFFIX) {
			filename := pathname + PathSep + file.Name()
			log.Println("load default param from:", filename)
			exp_manager.load_single_default(filename)
		}
	}
}

/**
 * load exp param from one path
 */
func (exp_manager *ExpManager) Load_exp_from_path(pathname string) {
	dir, err := ioutil.ReadDir(pathname)
	if err != nil {
		log.Fatalln("Read exp_param_path failed!")
	}
	PathSep := string(os.PathSeparator)
	for _, file := range dir {
		if file.IsDir() {
			continue
		}
		if strings.HasSuffix(file.Name(), EXP_SUFFIX) {
			filename := pathname + PathSep + file.Name()
			log.Println("load exp param from:", filename)
			exp_manager.load_single_exp(filename)
		}
	}
}

/**
 * load exp param from json file
 */
func (exp_manager *ExpManager) load_single_exp(filename string) {
	json_conf, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalln("Open exp_param_conf failed!")
	}
	js, err := json.NewJson([]byte(json_conf))
	if err != nil {
		log.Fatalln("Parse exp_param_conf failed!")
	}
	var layer_info LayerInfo
	layer_info.layer_type, err = js.Get("layer_type").String()
	if err != nil {
		log.Fatalln("Parse layer_type failed!")
	}
	layer_info.layer_name, err = js.Get("layer_name").String()
	if err != nil {
		log.Fatalln("Parse layer_name failed!")
	}
	layer_info.layer_sign, err = js.Get("layer_sign").Uint64()
	if err != nil {
		log.Fatalln("Parse layer_sign failed!")
	}
	json_explist, err := js.Get("exp_list").Array()
	if err != nil {
		log.Fatalln("Parse exp_list failed!")
	}
	for i := 0; i < len(json_explist); i++ {
		var exp_info ExpInfo
		json_exp := js.Get("exp_list").GetIndex(i)
		exp_info.exp_name, err = json_exp.Get("exp_name").String()
		if err != nil {
			log.Fatalln("Parse exp_name failed!")
		}
		exp_info.base_id, err = json_exp.Get("base_id").String()
		if err != nil {
			log.Fatalln("Parse base_id failed!")
		}
		exp_info.base_start, err = json_exp.Get("base_start").Uint64()
		if err != nil {
			log.Fatalln("Parse base_start failed!")
		}
		exp_info.base_end, err = json_exp.Get("base_end").Uint64()
		if err != nil {
			log.Fatalln("Parse base_end failed!")
		}
		exp_info.exp_id, err = json_exp.Get("exp_id").String()
		if err != nil {
			log.Fatalln("Parse exp_id failed!")
		}
		exp_info.exp_start, err = json_exp.Get("exp_start").Uint64()
		if err != nil {
			log.Fatalln("Parse exp_start failed!")
		}
		exp_info.exp_end, err = json_exp.Get("exp_end").Uint64()
		if err != nil {
			log.Fatalln("Parse exp_end failed!")
		}
		exp_info.exp_params = make(map[string]string)
		json_params, err := json_exp.Get("exp_params").Map()
		if err != nil {
			log.Fatalln("Parse exp_params failed!")
		}
		for k, v := range json_params {
			exp_info.exp_params[k] = v.(string)
		}
		layer_info.exp_list = append(layer_info.exp_list, exp_info)
	}
	if layer_info.layer_type == SINGLE {
		exp_manager.single_layer_list = append(exp_manager.single_layer_list, layer_info)
	} else if layer_info.layer_type == MULTIPLE {
		exp_manager.multi_layer_list = append(exp_manager.multi_layer_list, layer_info)
	} else {
		log.Fatalln("ERROR layer_type!")
	}
}

/**
 * get exp assign result
 */
func (exp_manager *ExpManager) Get_assign(exp_request_info ExpRequestInfo) ExpAssignInfo {
	var exp_assign_info ExpAssignInfo
	// copy default params to exp_assign_info
	exp_assign_info.Param_infos = make(map[string]string)
	for k, v := range exp_manager.default_params {
		exp_assign_info.Param_infos[k] = v
	}
	// first assign single layer
	for _, single_layer := range exp_manager.single_layer_list {
		hash_code := farmhash.Hash64WithSeed([]byte(exp_request_info.Androidid), single_layer.layer_sign)
		num := hash_code % MODE_BASE
		log.Println("single_hash_code:", hash_code)
		log.Println("single_num:", num)
		for _, exp_info := range single_layer.exp_list {
			if num >= exp_info.base_start && num < exp_info.base_end {
				exp_assign_info.Exp_ids = append(exp_assign_info.Exp_ids, exp_info.base_id)
			} else if num >= exp_info.exp_start && num < exp_info.exp_end {
				exp_assign_info.Exp_ids = append(exp_assign_info.Exp_ids, exp_info.exp_id)
				// modify default param to exp param
				for k, v := range exp_info.exp_params {
					exp_assign_info.Param_infos[k] = v
				}
				break
			}
		}
		// one request only in one single layer
		if len(exp_assign_info.Exp_ids) != 0 {
			return exp_assign_info
		}
	}
	// then assign multiple layer
	for _, multi_layer := range exp_manager.multi_layer_list {
		hash_code := farmhash.Hash64WithSeed([]byte(exp_request_info.Androidid), multi_layer.layer_sign)
		num := hash_code % MODE_BASE
		log.Println("multi_hash_code:", hash_code)
		log.Println("multi_num:", num)
		for _, exp_info := range multi_layer.exp_list {
			if num >= exp_info.base_start && num < exp_info.base_end {
				exp_assign_info.Exp_ids = append(exp_assign_info.Exp_ids, exp_info.base_id)
			} else if num >= exp_info.exp_start && num < exp_info.exp_end {
				exp_assign_info.Exp_ids = append(exp_assign_info.Exp_ids, exp_info.exp_id)
				// modify default param to exp param
				for k, v := range exp_info.exp_params {
					exp_assign_info.Param_infos[k] = v
				}
				break
			}
		}
	}
	return exp_assign_info
}
