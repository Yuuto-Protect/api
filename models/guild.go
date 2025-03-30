package models

type Guild struct {
	GuildId          string `gorm:"primaryKey"`
	Owners           []string
	Prefix           string
	AntiSpam         bool
	AntiSpamTime     int
	AntiSpamMessages int
	AntiLink         bool
	AntiRaid         bool
	Captcha          bool
	CaptchaChannel   *string
	AntiBadWords     bool
	BadWords         []string
	AutoModRules     map[string]string
	XpEnabled        bool
	XpPerMessage     int
}
