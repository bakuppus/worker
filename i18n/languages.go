package i18n

type Language string

func (l Language) String() string {
	return string(l)
}

const (
	Arabic        Language = "ar"
	Bulgarian     Language = "bg"
	Czech         Language = "cz"
	Danish        Language = "dk"
	German        Language = "de"
	Greek         Language = "el"
	English       Language = "en"
	Spanish       Language = "es"
	Finnish       Language = "fi"
	French        Language = "fr"
	Croatian      Language = "hr"
	Hungarian     Language = "hu"
	Italian       Language = "it"
	Japanese      Language = "jp"
	Korean        Language = "kr"
	Lithuanian    Language = "lt"
	Dutch         Language = "nl"
	Norwegian     Language = "no"
	Polish        Language = "pl"
	PortugueseBR  Language = "br"
	Portuguese    Language = "pt"
	Romanian      Language = "ro"
	Russian       Language = "ru"
	Slovak        Language = "sk"
	Swedish       Language = "sv"
	Thai          Language = "th"
	Turkish       Language = "tr"
	Ukrainian     Language = "ua"
	Vietnamese    Language = "vn"
	Welsh         Language = "cy"
	Chinese       Language = "cn"
	ChineseTaiwan Language = "tw"
)

var Flags = map[Language]string{
	Arabic:        "đ¸đŚ",
	Bulgarian:     "đ§đŹ",
	Czech:         "đ¨đż",
	Danish:        "đŠđ°",
	German:        "đŠđŞ",
	Greek:         "đŹđˇ",
	English:       "đ´ó §ó ˘ó Ľó Žó §ó ż",
	Spanish:       "đŞđ¸",
	Finnish:       "đŤđŽ",
	French:        "đŤđˇ",
	Croatian:      "đ­đˇ",
	Hungarian:     "đ­đş",
	Italian:       "đŽđš",
	Japanese:      "đŻđľ",
	Korean:        "đ°đˇ",
	Lithuanian:    "đąđš",
	Dutch:         "đłđą",
	Norwegian:     "đłđ´",
	Polish:        "đľđą",
	PortugueseBR:  "đ§đˇ",
	Portuguese:    "đľđš",
	Romanian:      "đˇđ´",
	Russian:       "đˇđş",
	Slovak:        "đ¸đ°",
	Swedish:       "đ¸đŞ",
	Thai:          "đšđ­",
	Turkish:       "đšđˇ",
	Ukrainian:     "đşđŚ",
	Vietnamese:    "đťđł",
	Welsh:         "đ´ó §ó ˘ó ˇó Źó łó ż",
	Chinese:       "đ¨đł",
	ChineseTaiwan: "đšđź",
}

// https://discord.com/developers/docs/reference#locales
var DiscordLocales = map[string]Language{
	"en-US": English,
	"en-GB": English,
	"bg":    Bulgarian,
	"zh-CN": Chinese,
	"zh-TW": ChineseTaiwan,
	"hr":    Croatian,
	"cs":    Czech,
	"da":    Danish,
	"nl":    Dutch,
	"el":    Greek,
	"fi":    Finnish,
	"fr":    French,
	"de":    German,
	"hu":    Hungarian,
	"it":    Italian,
	"ja":    Japanese,
	"ko":    Korean,
	"lt":    Lithuanian,
	"no":    Norwegian,
	"pl":    Polish,
	"pt-BR": PortugueseBR,
	"ro":    Romanian,
	"ru":    Russian,
	"es-ES": Spanish,
	"sv-SE": Swedish,
	"th":    Thai,
	"tr":    Turkish,
	"uk":    Ukrainian,
}

// Used by CrowdIn
var FullLocales = map[string]Language{
	"ar-SA": Arabic,
	"bg-BG": Bulgarian,
	"cs-CZ": Czech,
	"cy-GB": Welsh,
	"da-DK": Danish,
	"de-DE": German,
	"el-GR": Greek,
	"en-GB": English,
	"es-ES": Spanish,
	"fi-FI": Finnish,
	"fr-FR": French,
	"hr-HR": Croatian,
	"hu-HU": Hungarian,
	"it-IT": Italian,
	"ja-JP": Japanese,
	"ko-KR": Korean,
	"lt-LT": Lithuanian,
	"nl-NL": Dutch,
	"no-NO": Norwegian,
	"pl-PL": Polish,
	"pt-BR": PortugueseBR,
	"pt-PT": Portuguese,
	"ro-RO": Romanian,
	"ru-RU": Russian,
	"sk-SK": Slovak,
	"sv-SE": Swedish,
	"th-TH": Thai,
	"tr-TR": Turkish,
	"uk-UA": Ukrainian,
	"vi-VN": Vietnamese,
	"zh-CN": Chinese,
	"zh-TW": ChineseTaiwan,
}

var LocalesInverse = map[Language]string{
	Arabic:        "ar-SA",
	Bulgarian:     "bg-BG",
	Czech:         "cs-CZ",
	Danish:        "da-DK",
	German:        "de-DE",
	Greek:         "el-GR",
	English:       "en-GB",
	Spanish:       "es-ES",
	Finnish:       "fi-FI",
	French:        "fr-FR",
	Croatian:      "hr-HR",
	Hungarian:     "hu-HU",
	Italian:       "it-IT",
	Japanese:      "ja-JP",
	Korean:        "ko-KR",
	Lithuanian:    "lt-LT",
	Dutch:         "nl-NL",
	Norwegian:     "no-NO",
	Polish:        "pl-PL",
	PortugueseBR:  "pt-BR",
	Portuguese:    "pt-PT",
	Romanian:      "ro-RO",
	Russian:       "ru-RU",
	Slovak:        "sk-SK",
	Swedish:       "sv-SE",
	Thai:          "th-TH",
	Turkish:       "tr-TR",
	Ukrainian:     "uk-UA",
	Vietnamese:    "vi-VN",
	Welsh:         "cy-GB",
	Chinese:       "zh-CN",
	ChineseTaiwan: "zh-TW",
}

var LanguagesAlphabetical = [...]Language{
	Arabic,
	Bulgarian,
	Czech,
	Danish,
	German,
	Greek,
	English,
	Spanish,
	Finnish,
	French,
	Croatian,
	Hungarian,
	Italian,
	Japanese,
	Korean,
	Lithuanian,
	Dutch,
	Norwegian,
	Polish,
	PortugueseBR,
	Portuguese,
	Romanian,
	Russian,
	Slovak,
	Swedish,
	Thai,
	Turkish,
	Ukrainian,
	Vietnamese,
	Welsh,
	Chinese,
	ChineseTaiwan,
}

var FullNames = map[Language]string{
	Arabic:        "Ř§ŮŮŮŘšŮŘąŮŘ¨ŮŮŮŮŘŠŮ",
	Bulgarian:     "ĐąŃĐťĐłĐ°ŃŃĐşĐ¸",
	Czech:         "ÄeĹĄtina",
	Danish:        "Dansk",
	German:        "Deutsch",
	Greek:         "ÎÎťÎťÎˇÎ˝ÎšÎşÎŹ",
	English:       "English",
	Spanish:       "EspaĂąol",
	Finnish:       "Suomi",
	French:        "FranĂ§ais",
	Croatian:      "Hrvatski",
	Hungarian:     "Magyar",
	Italian:       "Italiano",
	Japanese:      "ćĽćŹčŞ",
	Korean:        "íęľ­ě´",
	Lithuanian:    "LietuviĹĄkai",
	Dutch:         "Nederlands",
	Norwegian:     "Norsk",
	Polish:        "Polski",
	PortugueseBR:  "PortuguĂŞs do Brasil",
	Portuguese:    "PortuguĂŞs",
	Romanian:      "RomĂ˘nÄ",
	Russian:       "PŃŃŃĐşĐ¸Đš",
	Slovak:        "SlovenskĂ˝",
	Swedish:       "Svenska",
	Thai:          "ŕšŕ¸ŕ¸˘",
	Turkish:       "TĂźrkĂ§e",
	Ukrainian:     "ĐŁĐşŃĐ°ŃĐ˝ŃŃĐşĐ°",
	Vietnamese:    "Tiáşżng Viáťt",
	Welsh:         "Cymraeg",
	Chinese:       "ä¸­ć",
	ChineseTaiwan: "çšéŤä¸­ć",
}

var FullNamesEnglish = map[Language]string{
	Arabic:        "Arabic",
	Bulgarian:     "Bulgarian",
	Czech:         "Czech",
	Danish:        "Danish",
	German:        "German",
	Greek:         "Greek",
	English:       "English",
	Spanish:       "Spanish",
	Finnish:       "Finnish",
	French:        "French",
	Croatian:      "Croatian",
	Hungarian:     "Hungarian",
	Italian:       "Italian",
	Japanese:      "Japanese",
	Korean:        "Korean",
	Lithuanian:    "Lithuanian",
	Dutch:         "Dutch",
	Norwegian:     "Norwegian",
	Polish:        "Polish",
	PortugueseBR:  "Portuguese (Brazilian)",
	Portuguese:    "Portuguese",
	Romanian:      "Romanian",
	Russian:       "Russian",
	Slovak:        "Slovak",
	Swedish:       "Swedish",
	Thai:          "Thai",
	Turkish:       "Turkish",
	Ukrainian:     "Ukrainian",
	Vietnamese:    "Vietnamese",
	Welsh:         "Welsh",
	Chinese:       "Chinese",
	ChineseTaiwan: "Chinese (Taiwan)",
}
