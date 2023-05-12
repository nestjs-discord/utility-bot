package algolia

type credential struct {
	appId  string
	apiKey string
	index  string
}

// credentials represents the map of App to credential used for Algolia authentication.
var credentials = map[App]credential{
	Discord:        {appId: "7TYOYF10Z2", apiKey: "786517d17e19e9d306758dd276bc6574", index: "discord"},
	DiscordJSGuide: {appId: "8XSLZMKC5R", apiKey: "a2edfe9f29fe917013b23d5767ae569a", index: "discordjs"},
	Express:        {appId: "BH4D9OD16A", apiKey: "7164e33055faa6ecddefd9e08fc59f5d", index: "expressjs"},
	Fastify:        {appId: "DMPMC33PLU", apiKey: "12d46b3bfeee6511031cfe00778f3e45", index: "fastify"},
	Necord:         {appId: "U7YH0EPYI9", apiKey: "c41976c1ed280e75acc3e9efd4aaaf00", index: "necord"},
	NestCommander:  {appId: "9O0K4CXI15", apiKey: "9689faf6550ca3133e69be1d9861ea92", index: "nest-commander"},
	NestJS:         {appId: "SDCBYAN96J", apiKey: "6d1869890dab96592b191e63a8deb5b5", index: "nestjs"},
	Ogma:           {appId: "U5N45YQUS6", apiKey: "dad79a1521426f184d0fac2ce3575149", index: "ogma"},
	TypeScript:     {appId: "BGCDYOIYZ5", apiKey: "37ee06fa68db6aef451a490df6df7c60", index: "typescriptlang"},
}
