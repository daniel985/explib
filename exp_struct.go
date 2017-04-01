package explib

/**
 * exp information
 */
type ExpInfo struct {
	exp_name   string
	base_id    string
	base_start uint64
	base_end   uint64
	exp_id     string
	exp_start  uint64
	exp_end    uint64
	exp_params map[string]string
}

/**
 * exp layer information
 */
type LayerInfo struct {
	layer_type string
	layer_name string
	layer_sign uint64
	exp_list   []ExpInfo
}

/**
 * exp request information
 */
type ExpRequestInfo struct {
	Androidid string
}

/**
 * exp assign result
 */
type ExpAssignInfo struct {
	Exp_ids     []string
	Param_infos map[string]string
}
