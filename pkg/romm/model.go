package romm

import "time"

type Rom struct {
	Items  []Item `json:"items"`
	Total  int    `json:"total"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

type Item struct {
	Id                  int64         `json:"id"`
	IgdbId              interface{}   `json:"igdb_id"`
	SgdbId              interface{}   `json:"sgdb_id"`
	MobyId              interface{}   `json:"moby_id"`
	SsId                interface{}   `json:"ss_id"`
	PlatformId          int           `json:"platform_id"`
	PlatformSlug        string        `json:"platform_slug"`
	PlatformFsSlug      string        `json:"platform_fs_slug"`
	PlatformName        string        `json:"platform_name"`
	PlatformCustomName  string        `json:"platform_custom_name"`
	PlatformDisplayName string        `json:"platform_display_name"`
	FsName              string        `json:"fs_name"`
	FsNameNoTags        string        `json:"fs_name_no_tags"`
	FsNameNoExt         string        `json:"fs_name_no_ext"`
	FsExtension         string        `json:"fs_extension"`
	FsPath              string        `json:"fs_path"`
	FsSizeBytes         int64         `json:"fs_size_bytes"`
	Name                string        `json:"name"`
	Slug                interface{}   `json:"slug"`
	Summary             interface{}   `json:"summary"`
	AlternativeNames    []interface{} `json:"alternative_names"`
	YoutubeVideoId      interface{}   `json:"youtube_video_id"`
	Metadatum           struct {
		RomId            int           `json:"rom_id"`
		Genres           []interface{} `json:"genres"`
		Franchises       []interface{} `json:"franchises"`
		Collections      []interface{} `json:"collections"`
		Companies        []interface{} `json:"companies"`
		GameModes        []interface{} `json:"game_modes"`
		AgeRatings       []interface{} `json:"age_ratings"`
		FirstReleaseDate interface{}   `json:"first_release_date"`
		AverageRating    interface{}   `json:"average_rating"`
	} `json:"metadatum"`
	IgdbMetadata struct {
	} `json:"igdb_metadata"`
	MobyMetadata struct {
	} `json:"moby_metadata"`
	SsMetadata struct {
	} `json:"ss_metadata"`
	PathCoverSmall string   `json:"path_cover_small"`
	PathCoverLarge string   `json:"path_cover_large"`
	UrlCover       string   `json:"url_cover"`
	HasManual      bool     `json:"has_manual"`
	PathManual     string   `json:"path_manual"`
	UrlManual      string   `json:"url_manual"`
	IsUnidentified bool     `json:"is_unidentified"`
	Revision       string   `json:"revision"`
	Regions        []string `json:"regions"`
	Languages      []string `json:"languages"`
	Tags           []string `json:"tags"`
	CrcHash        string   `json:"crc_hash"`
	Md5Hash        *string  `json:"md5_hash"`
	Sha1Hash       *string  `json:"sha1_hash"`
	Multi          bool     `json:"multi"`
	Files          []struct {
		Id            int         `json:"id"`
		RomId         int         `json:"rom_id"`
		FileName      string      `json:"file_name"`
		FilePath      string      `json:"file_path"`
		FileSizeBytes int         `json:"file_size_bytes"`
		FullPath      string      `json:"full_path"`
		CreatedAt     time.Time   `json:"created_at"`
		UpdatedAt     time.Time   `json:"updated_at"`
		LastModified  time.Time   `json:"last_modified"`
		CrcHash       string      `json:"crc_hash"`
		Md5Hash       string      `json:"md5_hash"`
		Sha1Hash      string      `json:"sha1_hash"`
		Category      interface{} `json:"category"`
	} `json:"files"`
	FullPath  string        `json:"full_path"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	Siblings  []interface{} `json:"siblings"`
	RomUser   struct {
		Id              int         `json:"id"`
		UserId          int         `json:"user_id"`
		RomId           int         `json:"rom_id"`
		CreatedAt       time.Time   `json:"created_at"`
		UpdatedAt       time.Time   `json:"updated_at"`
		LastPlayed      interface{} `json:"last_played"`
		NoteRawMarkdown string      `json:"note_raw_markdown"`
		NoteIsPublic    bool        `json:"note_is_public"`
		IsMainSibling   bool        `json:"is_main_sibling"`
		Backlogged      bool        `json:"backlogged"`
		NowPlaying      bool        `json:"now_playing"`
		Hidden          bool        `json:"hidden"`
		Rating          int         `json:"rating"`
		Difficulty      int         `json:"difficulty"`
		Completion      int         `json:"completion"`
		Status          interface{} `json:"status"`
		UserUsername    string      `json:"user__username"`
	} `json:"rom_user"`
}

type SearchResponse struct {
	Id           interface{} `json:"id"`
	IgdbId       *int64      `json:"igdb_id"`
	MobyId       interface{} `json:"moby_id"`
	SsId         interface{} `json:"ss_id"`
	Slug         *string     `json:"slug"`
	Name         string      `json:"name"`
	Summary      string      `json:"summary"`
	IgdbUrlCover string      `json:"igdb_url_cover"`
	MobyUrlCover string      `json:"moby_url_cover"`
	SsUrlCover   string      `json:"ss_url_cover"`
	PlatformId   int         `json:"platform_id"`
}
