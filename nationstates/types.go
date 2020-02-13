package nationstates

import (
	"encoding/xml"
)

const (
	NoticeTelegram          = "TG"
	NoticeIssue             = "I"
	NoticeEndorsementGained = "END"
	NoticeEndorsementLost   = "UNEND"
	NoticeBanner            = "U"
	NoticeRank              = "T"
	NoticePolicy            = "P"
	NoticeTradingCards      = "C"
	NoticeRMBMention        = "RMB"
	NoticeRMBQuote          = "RMBQ"
	NoticeRMBLike           = "RMBL"
	NoticeDispatchMention   = "D"
	NoticeDispatchPin       = "DP"
	NoticeDispatchQuote     = "DQ"
	NoticeEmbassy           = "EMB"
	NoticeLoomingApocalypse = "X"
)

const (
	CensusCivilRights                     = 0
	CensusEconomy                         = 1
	CensusPoliticalFreedom                = 2
	CensusPopulation                      = 3
	CensusAuthoritarianism                = 53
	CensusAverageDisposableIncome         = 85
	CensusAverageIncome                   = 72
	CensusAverageIncomeofPoor             = 73
	CensusAverageIncomeofRich             = 74
	CensusAverageness                     = 67
	CensusBlackMarket                     = 79
	CensusBusinessSubsidization           = 31
	CensusCharmlessness                   = 64
	CensusCheerfulness                    = 40
	CensusCompassion                      = 6
	CensusCompliance                      = 42
	CensusCorruption                      = 51
	CensusCrime                           = 77
	CensusCulture                         = 55
	CensusDeathRate                       = 5
	CensusDefenseForces                   = 46
	CensusEcoFriendliness                 = 7
	CensusEconomicFreedom                 = 48
	CensusEconomicOutput                  = 76
	CensusEmployment                      = 56
	CensusEnvironmentalBeauty             = 63
	CensusForeignAid                      = 78
	CensusFreedomFromTaxation             = 50
	CensusGovernmentSize                  = 27
	CensusHealth                          = 39
	CensusHumanDevelopmentIndex           = 68
	CensusIdeologicalRadicality           = 45
	CensusIgnorance                       = 37
	CensusInclusiveness                   = 71
	CensusIncomeEquality                  = 33
	CensusIndustryArmsManufacturing       = 16
	CensusIndustryAutomobileManufacturing = 10
	CensusIndustryBasketWeaving           = 12
	CensusIndustryBeverageSales           = 18
	CensusIndustryBookPublishing          = 24
	CensusIndustryCheeseExports           = 11
	CensusIndustryFurnitureRestoration    = 22
	CensusIndustryGambling                = 25
	CensusIndustryInformationTechnology   = 13
	CensusIndustryInsurance               = 21
	CensusIndustryMining                  = 20
	CensusIndustryPizzaDelivery           = 14
	CensusIndustryRetail                  = 23
	CensusIndustryTimberWoodchipping      = 19
	CensusIndustryTroutFishing            = 15
	CensusInfluence                       = 65
	CensusIntegrity                       = 52
	CensusIntelligence                    = 36
	CensusInternationalArtwork            = 86
	CensusLawEnforcement                  = 30
	CensusLifespan                        = 44
	CensusNiceness                        = 34
	CensusNudity                          = 9
	CensusObesity                         = 61
	CensusPacifism                        = 47
	CensusPoliticalApathy                 = 38
	CensusPrimitiveness                   = 69
	CensusPublicEducation                 = 75
	CensusPublicHealthcare                = 29
	CensusPublicTransport                 = 57
	CensusRecreationalDrugUse             = 60
	CensusReligiousness                   = 32
	CensusResidency                       = 80
	CensusRudeness                        = 35
	CensusSafety                          = 43
	CensusScientificAdvancement           = 70
	CensusSectorAgriculture               = 17
	CensusSectorManufacturing             = 26
	CensusSecularism                      = 62
	CensusSocialConservatism              = 8
	CensusTaxation                        = 49
	CensusTourism                         = 58
	CensusWealthGaps                      = 4
	CensusWeaponization                   = 59
	CensusWeather                         = 41
	CensusWelfare                         = 28
	CensusWorldAssemblyEndorsements       = 66
	CensusYouthRebelliousness             = 54
	CensusAlphabetical                    = 254
)

var CensusLabels = map[int]string{
	CensusCivilRights:                     "Civil Rights",
	CensusEconomy:                         "Economy",
	CensusPoliticalFreedom:                "Political Freedom",
	CensusPopulation:                      "Population",
	CensusAuthoritarianism:                "Authoritarianism",
	CensusAverageDisposableIncome:         "Average Disposable Income",
	CensusAverageIncome:                   "Average Income",
	CensusAverageIncomeofPoor:             "Average Income of Poor",
	CensusAverageIncomeofRich:             "Average Income of Rich",
	CensusAverageness:                     "Averageness",
	CensusBlackMarket:                     "Black Market",
	CensusBusinessSubsidization:           "Business Subsidization",
	CensusCharmlessness:                   "Charmlessness",
	CensusCheerfulness:                    "Cheerfulness",
	CensusCompassion:                      "Compassion",
	CensusCompliance:                      "Compliance",
	CensusCorruption:                      "Corruption",
	CensusCrime:                           "Crime",
	CensusCulture:                         "Culture",
	CensusDeathRate:                       "Death Rate",
	CensusDefenseForces:                   "Defense Forces",
	CensusEcoFriendliness:                 "Eco-Friendliness",
	CensusEconomicFreedom:                 "Economic Freedom",
	CensusEconomicOutput:                  "Economic Output",
	CensusEmployment:                      "Employment",
	CensusEnvironmentalBeauty:             "Environmental Beauty",
	CensusForeignAid:                      "Foreign Aid",
	CensusFreedomFromTaxation:             "Freedom From Taxation",
	CensusGovernmentSize:                  "Government Size",
	CensusHealth:                          "Health",
	CensusHumanDevelopmentIndex:           "Human Development Index",
	CensusIdeologicalRadicality:           "Ideological Radicality",
	CensusIgnorance:                       "Ignorance",
	CensusInclusiveness:                   "Inclusiveness",
	CensusIncomeEquality:                  "Income Equality",
	CensusIndustryArmsManufacturing:       "Industry: Arms Manufacturing",
	CensusIndustryAutomobileManufacturing: "Industry: Automobile Manufacturing",
	CensusIndustryBasketWeaving:           "Industry: Basket Weaving",
	CensusIndustryBeverageSales:           "Industry: Beverage Sales",
	CensusIndustryBookPublishing:          "Industry: Book Publishing",
	CensusIndustryCheeseExports:           "Industry: Cheese Exports",
	CensusIndustryFurnitureRestoration:    "Industry: Furniture Restoration",
	CensusIndustryGambling:                "Industry: Gambling",
	CensusIndustryInformationTechnology:   "Industry: Information Technology",
	CensusIndustryInsurance:               "Industry: Insurance",
	CensusIndustryMining:                  "Industry: Mining",
	CensusIndustryPizzaDelivery:           "Industry: Pizza Delivery",
	CensusIndustryRetail:                  "Industry: Retail",
	CensusIndustryTimberWoodchipping:      "Industry: Timber Woodchipping",
	CensusIndustryTroutFishing:            "Industry: Trout Fishing",
	CensusInfluence:                       "Influence",
	CensusIntegrity:                       "Integrity",
	CensusIntelligence:                    "Intelligence",
	CensusInternationalArtwork:            "International Artwork",
	CensusLawEnforcement:                  "Law Enforcement",
	CensusLifespan:                        "Lifespan",
	CensusNiceness:                        "Niceness",
	CensusNudity:                          "Nudity",
	CensusObesity:                         "Obesity",
	CensusPacifism:                        "Pacifism",
	CensusPoliticalApathy:                 "Political Apathy",
	CensusPrimitiveness:                   "Primitiveness",
	CensusPublicEducation:                 "Public Education",
	CensusPublicHealthcare:                "Public Healthcare",
	CensusPublicTransport:                 "Public Transport",
	CensusRecreationalDrugUse:             "Recreational Drug Use",
	CensusReligiousness:                   "Religiousness",
	CensusResidency:                       "Residency",
	CensusRudeness:                        "Rudeness",
	CensusSafety:                          "Safety",
	CensusScientificAdvancement:           "Scientific Advancement",
	CensusSectorAgriculture:               "Sector: Agriculture",
	CensusSectorManufacturing:             "Next: Sector: Manufacturing",
	CensusSecularism:                      "Secularism",
	CensusSocialConservatism:              "Social Conservatism",
	CensusTaxation:                        "Taxation",
	CensusTourism:                         "Tourism",
	CensusWealthGaps:                      "Wealth Gaps",
	CensusWeaponization:                   "Weaponization",
	CensusWeather:                         "Weather",
	CensusWelfare:                         "Welfare",
	CensusWorldAssemblyEndorsements:       "World Assembly Endorsements",
	CensusYouthRebelliousness:             "Youth Rebelliousness",
	CensusAlphabetical:                    "Alphabetical",
}

type Nation struct {
	XMLName      xml.Name     `xml:"NATION"`
	ID           string       `xml:"id,attr"`
	Consequences Consequences `xml:"ISSUE"`
	Issues       []Issue      `xml:"ISSUES>ISSUE"`
	Notices      []Notice     `xml:"NOTICES>NOTICE"`
}

type Issue struct {
	ID      int      `xml:"id,attr"`
	Title   string   `xml:"TITLE"`
	Text    string   `xml:"TEXT"`
	Options []Option `xml:"OPTION"`
}

type Consequences struct {
	Desc      string   `xml:"DESC"`
	Rankings  []Rank   `xml:"RANKINGS>RANK"`
	Headlines []string `xml:"HEADLINES>HEADLINE"`

	Error string `xml:"ERROR"`
}

type Rank struct {
	Score   float32 `xml:"SCORE"`
	Change  float32 `xml:"CHANGE"`
	PChange float32 `xml:"PCHANGE"`
}

type Option struct {
	ID   int    `xml:"id,attr"`
	Text string `xml:",chardata"`
}

type Notice struct {
	Text      string `xml:"TEXT"`
	Timestamp int    `xml:"TIMESTAMP"`
	Title     string `xml:"TITLE"`
	Who       string `xml:"WHO"`
	URL       string `xml:"URL"`
	Type      string `xml:"TYPE"`
}
