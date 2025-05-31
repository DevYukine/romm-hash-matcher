package hasheous

type LookupRequest struct {
	MD5  string `json:"mD5"`
	ShA1 string `json:"shA1"`
	Crc  string `json:"crc"`
}

type LookupResponse struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Platform struct {
		Name     string `json:"name"`
		Metadata []struct {
			Id                 string `json:"id"`
			ImmutableId        string `json:"immutableId"`
			Status             string `json:"status"`
			MatchMethod        string `json:"matchMethod"`
			Source             string `json:"source"`
			Link               string `json:"link"`
			LastSearch         string `json:"lastSearch"`
			NextSearch         string `json:"nextSearch"`
			WinningVoteCount   int    `json:"winningVoteCount"`
			TotalVoteCount     int    `json:"totalVoteCount"`
			WinningVotePercent int    `json:"winningVotePercent"`
		} `json:"metadata"`
	} `json:"platform"`
	Publisher struct {
		Name     string `json:"name"`
		Metadata []struct {
			Id                 string `json:"id"`
			ImmutableId        string `json:"immutableId"`
			Status             string `json:"status"`
			MatchMethod        string `json:"matchMethod"`
			Source             string `json:"source"`
			Link               string `json:"link"`
			LastSearch         string `json:"lastSearch"`
			NextSearch         string `json:"nextSearch"`
			WinningVoteCount   int    `json:"winningVoteCount"`
			TotalVoteCount     int    `json:"totalVoteCount"`
			WinningVotePercent int    `json:"winningVotePercent"`
		} `json:"metadata"`
	} `json:"publisher"`
	Signature struct {
		Game struct {
			Countries struct {
				JP string `json:"JP"`
				US string `json:"US"`
			} `json:"countries"`
			Languages struct {
			} `json:"languages"`
			SystemId       int    `json:"systemId"`
			PublisherId    int    `json:"publisherId"`
			MetadataSource int    `json:"metadataSource"`
			Score          int    `json:"score"`
			Id             string `json:"id"`
			Category       string `json:"category"`
			Name           string `json:"name"`
			Description    string `json:"description"`
			Year           string `json:"year"`
			Publisher      string `json:"publisher"`
			Demo           string `json:"demo"`
			System         string `json:"system"`
			SystemVariant  string `json:"systemVariant"`
			Video          string `json:"video"`
			Country        struct {
			} `json:"country"`
			Language struct {
			} `json:"language"`
			Copyright string        `json:"copyright"`
			Roms      []interface{} `json:"roms"`
			Flags     struct {
			} `json:"flags"`
			RomCount int `json:"romCount"`
		} `json:"game"`
		Rom struct {
			Score   int    `json:"score"`
			Id      string `json:"id"`
			Name    string `json:"name"`
			Size    int    `json:"size"`
			Crc     string `json:"crc"`
			Md5     string `json:"md5"`
			Sha1    string `json:"sha1"`
			Country struct {
				JP string `json:"JP"`
			} `json:"country"`
			Language struct {
			} `json:"language"`
			DevelopmentStatus string `json:"developmentStatus"`
			RomType           string `json:"romType"`
			RomTypeMedia      string `json:"romTypeMedia"`
			MediaDetail       struct {
			} `json:"mediaDetail"`
			MediaLabel      string `json:"mediaLabel"`
			SignatureSource string `json:"signatureSource"`
		} `json:"rom"`
	} `json:"signature"`
	Metadata   []Metadata `json:"metadata"`
	Attributes []struct {
		Link                  string `json:"link,omitempty"`
		Id                    int    `json:"id,omitempty"`
		AttributeType         string `json:"attributeType"`
		AttributeName         string `json:"attributeName"`
		AttributeRelationType string `json:"attributeRelationType"`
		Value                 string `json:"value"`
	} `json:"attributes"`
}

type Metadata struct {
	Id                 string `json:"id"`
	ImmutableId        string `json:"immutableId"`
	Status             string `json:"status"`
	MatchMethod        string `json:"matchMethod"`
	Source             string `json:"source"`
	Link               string `json:"link"`
	LastSearch         string `json:"lastSearch"`
	NextSearch         string `json:"nextSearch"`
	WinningVoteCount   int    `json:"winningVoteCount"`
	TotalVoteCount     int    `json:"totalVoteCount"`
	WinningVotePercent int    `json:"winningVotePercent"`
}
