package dns

import "strconv"

var queryType = map[uint16]string{
	1:  "A",
	2:  "NS",
	5:  "CNAME",
	6:  "SOA",
	12: "PTR",
	15: "MX",
	16: "TXT",
	28: "AAAA",
	33: "SRV",
	64: "SVCB",
	65: "HTTPS",
}

func getTypeString(t uint16) string {
	q, ok := queryType[t]
	if !ok {
		str := strconv.Itoa(int(t))
		return str
	}
	return q
}
