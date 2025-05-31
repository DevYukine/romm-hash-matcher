package playmatch

type GameMatchType string

const (
	GameMatchTypeSHA256          GameMatchType = "SHA256"
	GameMatchTypeSHA1            GameMatchType = "SHA1"
	GameMatchTypeMD5             GameMatchType = "MD5"
	GameMatchTypeFileNameAndSize GameMatchType = "FileNameAndSize"
	GameMatchTypeNoMatch         GameMatchType = "NoMatch"
)

type MetadataMatchType string

const (
	MetadataMatchTypeAutomatic MetadataMatchType = "Automatic"
	MetadataMatchTypeManual    MetadataMatchType = "Manual"
	MetadataMatchTypeFailed    MetadataMatchType = "Failed"
	MetadataMatchTypeNone      MetadataMatchType = "None"
)

type MetadataFailedMatchReason string

const (
	MetadataFailedMatchReasonNoDirectMatch  MetadataFailedMatchReason = "NoDirectMatch"
	MetadataFailedMatchReasonTooManyMatches MetadataFailedMatchReason = "TooManyMatches"
)

type MetadataManualMatchType string

const (
	MetadataManualMatchTypeAdmin     MetadataManualMatchType = "Admin"
	MetadataManualMatchTypeTrusted   MetadataManualMatchType = "Trusted"
	MetadataManualMatchTypeCommunity MetadataManualMatchType = "Community"
)

type MetadataAutomaticMatchReason string

const (
	MetadataAutomaticMatchReasonAlternativeName MetadataAutomaticMatchReason = "AlternativeName"
	MetadataAutomaticMatchReasonDirectName      MetadataAutomaticMatchReason = "DirectName"
	MetadataAutomaticMatchReasonViaChild        MetadataAutomaticMatchReason = "ViaChild"
	MetadataAutomaticMatchReasonViaParent       MetadataAutomaticMatchReason = "ViaParent"
)

type MatchResponse struct {
	ExternalMetadata []struct {
		AutomaticMatchReason MetadataAutomaticMatchReason `json:"automaticMatchReason"`
		Comment              string                       `json:"comment"`
		FailedMatchReason    MetadataFailedMatchReason    `json:"failedMatchReason"`
		ManualMatchType      MetadataManualMatchType      `json:"manualMatchType"`
		MatchType            MetadataMatchType            `json:"matchType"`
		ProviderId           string                       `json:"providerId"`
		ProviderName         string                       `json:"providerName"`
	} `json:"externalMetadata"`
	GameMatchType GameMatchType `json:"gameMatchType"`
	Id            string        `json:"id"`
}

type SlimmedGameResponse struct {
	Id int64 `json:"id"`
}
