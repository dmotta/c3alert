package c3alert

type Alert struct {
	Attributes struct {
		class       string `json:"CLASS"`
		McObjectURI string `json:"mc_object_uri"`
		Severity    string `json:"severity"`
		Msg         string `json:"msg"`
		McHost      string `json:"mc_host"`
		McSmcAlias  string `json:"mc_smc_alias"`
		McSmcID     string `json:"mc_smc_id"`
		McOwner     string `json:"mc_owner"`
		McPriority  string `json:"mc_priority"`
	} `json:"attributes"`
}