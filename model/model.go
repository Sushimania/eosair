package model

type GetInfo struct {
	ServerVersion            string `json:"server_version"`
	ChainID                  string `json:"chain_id"`
	HeadBlockNum             int    `json:"head_block_num"`
	LastIrreversibleBlockNum int    `json:"last_irreversible_block_num"`
	LastIrreversibleBlockID  string `json:"last_irreversible_block_id"`
	HeadBlockID              string `json:"head_block_id"`
	HeadBlockTime            string `json:"head_block_time"`
	HeadBlockProducer        string `json:"head_block_producer"`
	VirtualBlockCPULimit     int    `json:"virtual_block_cpu_limit"`
	VirtualBlockNetLimit     int    `json:"virtual_block_net_limit"`
	BlockCPULimit            int    `json:"block_cpu_limit"`
	BlockNetLimit            int    `json:"block_net_limit"`
}

type BlockDetail struct {
	Timestamp         string        `json:"timestamp"`
	Producer          string        `json:"producer"`
	Confirmed         int           `json:"confirmed"`
	Previous          string        `json:"previous"`
	TransactionMroot  string        `json:"transaction_mroot"`
	ActionMroot       string        `json:"action_mroot"`
	ScheduleVersion   int           `json:"schedule_version"`
	NewProducers      interface{}   `json:"new_producers"`
	HeaderExtensions  []interface{} `json:"header_extensions"`
	ProducerSignature string        `json:"producer_signature"`
	Transactions      []struct {
		Status        string `json:"status"`
		CPUUsageUs    int    `json:"cpu_usage_us"`
		NetUsageWords int    `json:"net_usage_words"`
		Trx           struct {
			ID                    string        `json:"id"`
			Signatures            []string      `json:"signatures"`
			Compression           string        `json:"compression"`
			PackedContextFreeData string        `json:"packed_context_free_data"`
			ContextFreeData       []interface{} `json:"context_free_data"`
			PackedTrx             string        `json:"packed_trx"`
			Transaction           struct {
				Expiration         string        `json:"expiration"`
				RefBlockNum        int           `json:"ref_block_num"`
				RefBlockPrefix     int           `json:"ref_block_prefix"`
				MaxNetUsageWords   int           `json:"max_net_usage_words"`
				MaxCPUUsageMs      int           `json:"max_cpu_usage_ms"`
				DelaySec           int           `json:"delay_sec"`
				ContextFreeActions []interface{} `json:"context_free_actions"`
				Actions            []struct {
					Account       string `json:"account"`
					Name          string `json:"name"`
					Authorization []struct {
						Actor      string `json:"actor"`
						Permission string `json:"permission"`
					} `json:"authorization"`
					Data struct {
						From     string `json:"from"`
						To       string `json:"to"`
						Quantity string `json:"quantity"`
						Memo     string `json:"memo"`
					} `json:"data"`
					HexData string `json:"hex_data"`
				} `json:"actions"`
				TransactionExtensions []interface{} `json:"transaction_extensions"`
			} `json:"transaction"`
		} `json:"trx"`
	} `json:"transactions"`
	BlockExtensions []interface{} `json:"block_extensions"`
	ID              string        `json:"id"`
	BlockNum        int           `json:"block_num"`
	RefBlockPrefix  int64         `json:"ref_block_prefix"`
}

type Block struct {
	Timestamp         string        `json:"timestamp"`
	Producer          string        `json:"producer"`
	Confirmed         int           `json:"confirmed"`
	Previous          string        `json:"previous"`
	TransactionMroot  string        `json:"transaction_mroot"`
	ActionMroot       string        `json:"action_mroot"`
	ScheduleVersion   int           `json:"schedule_version"`
	NewProducers      interface{}   `json:"new_producers"`
	HeaderExtensions  []interface{} `json:"header_extensions"`
	ProducerSignature string        `json:"producer_signature"`
	Transactions      []struct {
		Status        string `json:"status"`
		CPUUsageUs    int    `json:"cpu_usage_us"`
		NetUsageWords int    `json:"net_usage_words"`
		Trx           string `json:"trx"`
	} `json:"transactions"`
	BlockExtensions []interface{} `json:"block_extensions"`
	ID              string        `json:"id"`
	BlockNum        int           `json:"block_num"`
	RefBlockPrefix  int           `json:"ref_block_prefix"`
}

type Trx struct {
	ID                    string        `json:"id"`
	Signatures            []string      `json:"signatures"`
	Compression           string        `json:"compression"`
	PackedContextFreeData string        `json:"packed_context_free_data"`
	ContextFreeData       []interface{} `json:"context_free_data"`
	PackedTrx             string        `json:"packed_trx"`
	Transaction           struct {
		Expiration         string        `json:"expiration"`
		RefBlockNum        int           `json:"ref_block_num"`
		RefBlockPrefix     int           `json:"ref_block_prefix"`
		MaxNetUsageWords   int           `json:"max_net_usage_words"`
		MaxCPUUsageMs      int           `json:"max_cpu_usage_ms"`
		DelaySec           int           `json:"delay_sec"`
		ContextFreeActions []interface{} `json:"context_free_actions"`
		Actions            []struct {
			Account       string `json:"account"`
			Name          string `json:"name"`
			Authorization []struct {
				Actor      string `json:"actor"`
				Permission string `json:"permission"`
			} `json:"authorization"`
			Data struct {
				From     string `json:"from"`
				To       string `json:"to"`
				Quantity string `json:"quantity"`
				Memo     string `json:"memo"`
			} `json:"data"`
			HexData string `json:"hex_data"`
		} `json:"actions"`
		TransactionExtensions []interface{} `json:"transaction_extensions"`
	} `json:"transaction"`
}

type AccountInfo struct {
	AccountName       string `json:"account_name"`
	HeadBlockNum      int    `json:"head_block_num"`
	HeadBlockTime     string `json:"head_block_time"`
	Privileged        bool   `json:"privileged"`
	LastCodeUpdate    string `json:"last_code_update"`
	Created           string `json:"created"`
	CoreLiquidBalance string `json:"core_liquid_balance"`
	RAMQuota          int    `json:"ram_quota"`
	NetWeight         int    `json:"net_weight"`
	CPUWeight         int    `json:"cpu_weight"`
	NetLimit          struct {
		Used      int `json:"used"`
		Available int `json:"available"`
		Max       int `json:"max"`
	} `json:"net_limit"`
	CPULimit struct {
		Used      int `json:"used"`
		Available int `json:"available"`
		Max       int `json:"max"`
	} `json:"cpu_limit"`
	RAMUsage    int `json:"ram_usage"`
	Permissions []struct {
		PermName     string `json:"perm_name"`
		Parent       string `json:"parent"`
		RequiredAuth struct {
			Threshold int `json:"threshold"`
			Keys      []struct {
				Key    string `json:"key"`
				Weight int    `json:"weight"`
			} `json:"keys"`
			Accounts []interface{} `json:"accounts"`
			Waits    []interface{} `json:"waits"`
		} `json:"required_auth"`
	} `json:"permissions"`
	TotalResources struct {
		Owner     string `json:"owner"`
		NetWeight string `json:"net_weight"`
		CPUWeight string `json:"cpu_weight"`
		RAMBytes  int    `json:"ram_bytes"`
	} `json:"total_resources"`
	SelfDelegatedBandwidth struct {
		From      string `json:"from"`
		To        string `json:"to"`
		NetWeight string `json:"net_weight"`
		CPUWeight string `json:"cpu_weight"`
	} `json:"self_delegated_bandwidth"`
	RefundRequest struct {
		Owner       string `json:"owner"`
		RequestTime string `json:"request_time"`
		NetAmount   string `json:"net_amount"`
		CPUAmount   string `json:"cpu_amount"`
	} `json:"refund_request"`
	VoterInfo struct {
		Owner             string   `json:"owner"`
		Proxy             string   `json:"proxy"`
		Producers         []string `json:"producers"`
		Staked            int      `json:"staked"`
		LastVoteWeight    string   `json:"last_vote_weight"`
		ProxiedVoteWeight string   `json:"proxied_vote_weight"`
		IsProxy           int      `json:"is_proxy"`
	} `json:"voter_info"`
}