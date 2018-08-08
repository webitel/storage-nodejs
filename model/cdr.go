package model

import "encoding/json"

const (
	CDR_TYPE_NAME = "cdr"
)

type CdrData struct {
	Event StringInterface
}

type NoSQLRequest struct {
	Size      int         `json:"size"`
	Aggregate interface{} `json:"aggs"`
	Version   bool        `json:"version"`
	Query     string      `json:"query"`
	Filter    interface{} `json:"filter"`
	Includes  []string    `json:"includes"`
	Columns   []string    `json:"columns"`
	Sort      interface{} `json:"sort"`
}

type CdrCall struct {
	LegA  StringInterface   `json:"leg_a" db:"leg_a"`
	LegsB []StringInterface `json:"legs_b" db:"legs_b"`
}

func (self *CdrCall) ToJSON() string {
	b, _ := json.Marshal(self)
	return string(b)
}
func (self *CdrData) ToJSON() string {
	b, _ := json.Marshal(self.Event)
	return string(b)
}

type CdrTimes struct {
	CreatedTime        uint64 `json:"created_time,omitempty"`
	ProfileCreatedTime uint64 `json:"profile_created_time,omitempty"`
	ProgressTime       uint64 `json:"progress_time,omitempty"`
	ProgressMediaTime  uint64 `json:"progress_media_time,omitempty"`
	AnsweredTime       uint64 `json:"answered_time,omitempty"`
	BridgedTime        uint64 `json:"bridged_time,omitempty"`
	LastHoldTime       uint64 `json:"last_hold_time,omitempty"`
	HoldAccumTime      uint64 `json:"hold_accum_time,omitempty"`
	HangupTime         uint64 `json:"hangup_time,omitempty"`
	ResurrectTime      uint64 `json:"resurrect_time,omitempty"`
	TransferTime       uint64 `json:"transfer_time,omitempty"`
}

type CdrLocations struct {
	Geo         string `json:"geo,omitempty"`
	City        string `json:"city,omitempty"`
	Country     string `json:"country,omitempty"`
	CountryCode string `json:"country_code,omitempty"`
	Type        string `json:"type,omitempty"`
}

type CdrQueue struct {
	CC_Queue_Name          string `json:"name,omitempty"`
	Queue_CallDuration     uint32 `json:"duration,omitempty"`
	Queue_WaitingDuration  uint32 `json:"wait_duration,omitempty"`
	CC_CancelReason        string `json:"cancel_reason,omitempty"`
	CC_Cause               string `json:"cause,omitempty"`
	CC_Queue_AnsweredEpoch uint64 `json:"answered_time,omitempty"`
	CC_Queue_Hangup        uint64 `json:"exit_time,omitempty"`
	CC_Queue_JoinedEpoch   uint64 `json:"joined_time,omitempty"`
	CC_Side                string `json:"side,omitempty"`
}

type CdrCallFlow struct {
	CdrCallerProfile `json:"caller_profile,omitempty"`
	CdrTimes         `json:"times,omitempty"`
}

type CdrCallerProfile struct {
	Username          string `json:"username,omitempty"`
	CallerIdName      string `json:"caller_id_name,omitempty"`
	Ani               string `json:"ani,omitempty"`
	Aniii             string `json:"aniii,omitempty"`
	CallerIdNumber    string `json:"caller_id_number,omitempty"`
	NetworkAddr       string `json:"network_addr,omitempty"`
	Rdnis             string `json:"rdnis,omitempty"`
	DestinationNumber string `json:"destination_number,omitempty"`
	Uuid              string `json:"uuid,omitempty"`
	Source            string `json:"source,omitempty"`
}

type CdrSearch struct {
	Leg                  string `json:"leg,omitempty"`
	Parent_uuid          string `json:"parent_uuid,omitempty"`
	Uuid                 string `json:"uuid"`
	Direction            string `json:"direction,omitempty"`
	CallerIdName         string `json:"caller_id_name,omitempty"`
	CallerIdNumber       string `json:"caller_id_number,omitempty"`
	NetworkAddr          string `json:"network_addr,omitempty"`
	DestinationNumber    string `json:"destination_number,omitempty"`
	DomainName           string `json:"domain_name,omitempty"`
	Extension            string `json:"extension,omitempty"`
	PresenceId           string `json:"presence_id,omitempty"`
	Source               string `json:"source,omitempty"`
	Gateway              string `json:"gateway,omitempty"`
	Q850HangupCode       uint32 `json:"hangup_cause_q850"`
	HangupCause          string `json:"hangup_cause,omitempty"`
	HangupDisposition    string `json:"hangup_disposition,omitempty"`
	OriginateDisposition string `json:"originate_disposition,omitempty"`
	TransferDisposition  string `json:"transfer_disposition,omitempty"`
	CallCreatedTime      uint64 `json:"created_time,omitempty"`
	//times
	// BridgedTime     uint64 `json:"bridged_time,omitempty"`
	// CallAnswerTime  uint64 `json:"answered_time,omitempty"`
	// ProgressTime    uint64 `json:"progress_time,omitempty"`
	// CallHangupTime  uint64 `json:"hangup_time,omitempty"`
	//TransferTime    uint64 `json:"transfer_time,omitempty"`
	///////
	Duration              uint32 `json:"duration"`
	ConnectedCallDuration uint32 `json:"billsec"`
	ProgressSeconds       uint32 `json:"progresssec"`
	AnswerSeconds         uint32 `json:"answersec"`
	WaitSeconds           uint32 `json:"waitsec"`
	HoldAccumSeconds      uint32 `json:"holdsec"`
	HoldSecB              uint32 `json:"holdsec_b,omitempty"`
	TalkSec               uint32 `json:"talksec,omitempty"`
	///////
	QualityPercentageAudio uint32                 `json:"quality_percentage_audio,omitempty"`
	QualityPercentageVideo uint32                 `json:"quality_percentage_video,omitempty"`
	Variables              map[string]interface{} `json:"variables"`
	*CdrLocations          `json:"locations,omitempty"`
	*CdrQueue              `json:"queue,omitempty"`
	CallFlow               *[]CdrCallFlow `json:"callflow,omitempty"`
}
