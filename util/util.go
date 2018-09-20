package util

func SelectBPUrl (n int) string {
	var bpUrl string
	// https://api.eosdetroit.io:443 BP API 작동안함
	// https://eos.saltblock.io BP API 작동안함
	// https://api.eosuk.io:443 BP API 작동안함
	switch n % 15 {
	case 0:
		bpUrl = "https://api.eosnewyork.io:443"
		break
	case 1:
		bpUrl = "https://eos.greymass.com:443"
		break
	case 2:
		bpUrl = "http://api.hkeos.com:80"
		break
	case 3:
		bpUrl = "https://eosapi.blockmatrix.network:443"
		break
	case 4:
		bpUrl = "https://fn.eossweden.se:443"
		break
	case 5:
		bpUrl = "http://mainnet.eoscalgary.io:80"
		break
	case 6:
		bpUrl = "https://user-api.eoseoul.io:443"
		break
	case 7:
		bpUrl = "http://api1.eosdublin.io:80"
		break
	case 8:
		bpUrl = "http://api.cypherglass.com:8888"
		break
	case 9:
		bpUrl = "http://bp.cryptolions.io:8888"
		break
	case 10:
		bpUrl = "https://api.eosio.cr:443"
		break
	case 11:
		bpUrl = "https://api.eosn.io"
		break
	case 12:
		bpUrl = "https://eu1.eosdac.io:443"
		break
	case 13:
		bpUrl = "https://api.main.alohaeos.com:443"
		break
	case 14:
		bpUrl = "https://node1.eosphere.io"
		break
	}
	return bpUrl
}