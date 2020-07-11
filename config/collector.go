package config

type collectorConfig struct {
	ErrorTimes int
	Channels int
	TitleMinLength int
	ContentMinLength int
	TitleExclude []string
	TitleExcludePrefix []string
	TitleExcludeSuffix []string
	ContentExclude []string
	ContentExcludeLine []string
}