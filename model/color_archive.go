package model

// this is unitentryarchive struct for Table unit_entry_archive
type Color_archive struct {
	Color_id   int    `json:"color_id"`
	Color_name string `json:"color_name"`
	Changedate string `json:"changedate"`
	Changeflag string `json:"changeflag"`
	Trackid    int    `json:"trackid"`
	Changeuser string `json:"changeuser"`
}
