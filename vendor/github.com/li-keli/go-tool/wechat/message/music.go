package message

//Music 音乐消息
type Music struct {
	CommonToken

	Music struct {
		Title        string `json:"title" xml:"Title"        `
		Description  string `json:"description" xml:"Description"  `
		MusicURL     string `json:"music_url" xml:"MusicUrl"     `
		HQMusicURL   string `json:"hq_music_url" xml:"HQMusicUrl"   `
		ThumbMediaID string `json:"thumb_media_id" xml:"ThumbMediaId"`
	} `json:"music" xml:"Music"`
}

//NewMusic  回复音乐消息
func NewMusic(title, description, musicURL, hQMusicURL, thumbMediaID string) *Music {
	music := new(Music)
	music.Music.Title = title
	music.Music.Description = description
	music.Music.MusicURL = musicURL
	music.Music.ThumbMediaID = thumbMediaID
	return music
}
