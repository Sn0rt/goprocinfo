package linux

import (
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
)

type IpMib struct {
	Forwarding      uint64 `json:"forwarding"`
	DefaultTTL      uint64 `json:"default_ttl"`
	InReceives      uint64 `json:"in_receives"`
	InHdrErrors     uint64 `json:"in_hdr_errors"`
	InAddrErrors    uint64 `json:"in_addr_errors"`
	ForwDatagrams   uint64 `json:"forw_datagrams"`
	InUnknownProtos uint64 `json:"in_unknown_protos"`
	InDiscards      uint64 `json:"in_discards"`
	InDelivers      uint64 `json:"in_delivers"`
	OutRequests     uint64 `json:"out_requests"`
	OutDiscards     uint64 `json:"out_discards"`
	OutNoRoutes     uint64 `json:"out_no_routes"`
	ReasmTimeout    uint64 `json:"reasm_timeout"`
	ReasmReqds      uint64 `json:"reasm_reqds"`
	ReasmOKs        uint64 `json:"reasm_oKs"`
	ReasmFails      uint64 `json:"reasm_fails"`
	FragOKs         uint64 `json:"frag_oKs"`
	FragFails       uint64 `json:"frag_fails"`
	FragCreates     uint64 `json:"frag_creates"`
}

type IcmpMib struct {
	InMsgsInErrors   uint64 `json:"in_msgs_in_errors"`
	InCsumErrors     uint64 `json:"in_csum_errors"`
	InDestUnreachs   uint64 `json:"in_dest_unreachs"`
	InTimeExcds      uint64 `json:"in_time_excds"`
	InParmProbs      uint64 `json:"in_parm_probs"`
	InSrcQuenchs     uint64 `json:"in_src_quenchs"`
	InRedirects      uint64 `json:"in_redirects"`
	InEchos          uint64 `json:"in_echos"`
	InEchoReps       uint64 `json:"in_echo_reps"`
	InTimestamps     uint64 `json:"in_timestamps"`
	InTimestampReps  uint64 `json:"in_timestamp_reps"`
	InAddrMasks      uint64 `json:"in_addr_masks"`
	InAddrMaskReps   uint64 `json:"in_addr_mask_reps"`
	OutMsgs          uint64 `json:"out_msgs"`
	OutErrors        uint64 `json:"out_errors"`
	OutDestUnreachs  uint64 `json:"out_dest_unreachs"`
	OutTimeExcds     uint64 `json:"out_time_excds"`
	OutParmProbs     uint64 `json:"out_parm_probs"`
	OutSrcQuenchs    uint64 `json:"out_src_quenchs"`
	OutRedirects     uint64 `json:"out_redirects"`
	OutEchos         uint64 `json:"out_echos"`
	OutEchoReps      uint64 `json:"out_echo_reps"`
	OutTimestamps    uint64 `json:"out_timestamps"`
	OutTimestampReps uint64 `json:"out_timestamp_reps"`
	OutAddrMasks     uint64 `json:"out_addr_masks"`
	OutAddrMaskReps  uint64 `json:"out_addr_mask_reps"`
}

type IcmpMsgMib struct {
	InTypes  map[string]uint64 `json:"in_types"`
	OutTypes map[string]uint64 `json:"out_types"`
}

type TcpMib struct {
	RtoAlgorithmRtoMin uint64 `json:"rto_algorithm_rto_min"`
	RtoMax             uint64 `json:"rto_max"`
	MaxConn            uint64 `json:"max_conn"`
	ActiveOpens        uint64 `json:"active_opens"`
	PassiveOpens       uint64 `json:"passive_opens"`
	AttemptFails       uint64 `json:"attempt_fails"`
	EstabResets        uint64 `json:"estab_resets"`
	CurrEstab          uint64 `json:"curr_estab"`
	InSegs             uint64 `json:"in_segs"`
	OutSegs            uint64 `json:"out_segs"`
	RetransSegs        uint64 `json:"retrans_segs"`
	InErrs             uint64 `json:"in_errs"`
	OutRsts            uint64 `json:"out_rsts"`
	InCsumErrors       uint64 `json:"in_csum_errors"`
}

type UdpMib struct {
	InDatagramsNoPorts uint64 `json:"in_datagrams_no_ports"`
	InErrors           uint64 `json:"in_errors"`
	OutDatagrams       uint64 `json:"out_datagrams"`
	RcvbufErrors       uint64 `json:"rcvbuf_errors"`
	SndbufErrors       uint64 `json:"sndbuf_errors"`
	InCsumErrors       uint64 `json:"in_csum_errors"`
}

type UdpLiteMib struct {
	InDatagramsNoPorts uint64 `json:"in_datagrams_no_ports"`
	InErrors           uint64 `json:"in_errors"`
	OutDatagrams       uint64 `json:"out_datagrams"`
	RcvbufErrors       uint64 `json:"rcvbuf_errors"`
	SndbufErrors       uint64 `json:"sndbuf_errors"`
	InCsumErrors       uint64 `json:"in_csum_errors"`
}

type Snmp struct {
	Ip      IpMib
	Icmp    IcmpMib
	IcmpMsg IcmpMsgMib
	Tcp     TcpMib
	Udp     UdpMib
	UdpLite UdpLiteMib
}

func ReadSnmp(path string) (*Snmp, error) {
	data, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")

	var snmp Snmp = Snmp{}

	elem := reflect.ValueOf(&snmp).Elem()
	typeOfElem := elem.Type()

	for i := 1; i < len(lines); i = i + 2 {
		headers := strings.Fields(lines[i-1][strings.Index(lines[i-1], ":")+1:])
		values := strings.Fields(lines[i][strings.Index(lines[i], ":")+1:])
		fieldName := lines[i-1][:strings.Index(lines[i-1], ":")]

		if fieldName == "IcmpMsg" {
			inTypes := map[string]uint64{}
			outTypes := map[string]uint64{}
			for j, header := range headers {
				parsedValue, _ := strconv.ParseUint(values[j], 10, 64)
				if strings.Contains(header, "InType") {
					inTypes[header[6:]] = parsedValue
				} else {
					outTypes[header[7:]] = parsedValue
				}
			}
			snmp.IcmpMsg = IcmpMsgMib{InTypes: inTypes, OutTypes: outTypes}
			continue
		}

		if structType, ok := typeOfElem.FieldByName(fieldName); ok {
			fieldType := structType.Type
			fieldElem := reflect.New(fieldType).Elem()
			for j, header := range headers {
				if _, ok := fieldType.FieldByName(header); ok {
					parsedValue, _ := strconv.ParseUint(values[j], 10, 64)
					fieldElem.FieldByName(header).SetUint(parsedValue)
				}
			}
			elem.FieldByName(fieldName).Set(fieldElem)
		}
	}
	return &snmp, nil
}
