package consts

type Mode string

const (
	Production  Mode = "production"
	Development Mode = "development"
)

const (
	KeyId   = "keyId"
	KeyRole = "keyRole"
)

const (
	MaxSolutionFileSize = 1 * 1024 * 1024 * 1024
)

const EmptyProjectJson = "{\"targets\":[{\"isStage\":true,\"name\":\"Stage\",\"variables\":{\"`jEk@4|i[#Fk?(8x)AV." +
	"-my variable\":[\"my variable\",0]},\"lists\":{},\"broadcasts\":{},\"blocks\":{},\"currentCostume\":0,\"costum" +
	"es\":[{\"assetId\":\"cd21514d0531fdffb22204e0ec5ed84a\",\"name\":\"backdrop1\",\"md5ext\":\"cd21514d0531fdffb2" +
	"2204e0ec5ed84a.svg\",\"dataFormat\":\"svg\",\"rotationCenterX\":240,\"rotationCenterY\":180}],\"sounds\":[{\"a" +
	"ssetId\":\"83a9787d4cb6f3b7632b4ddfebf74367\",\"name\":\"pop\",\"dataFormat\":\"wav\",\"format\":\"\",\"rate\"" +
	":11025,\"sampleCount\":258,\"md5ext\":\"83a9787d4cb6f3b7632b4ddfebf74367.wav\"}],\"volume\":100},{\"isStage\":" +
	"false,\"name\":\"Sprite1\",\"variables\":{},\"lists\":{},\"broadcasts\":{},\"blocks\":{},\"currentCostume\":0," +
	"\"costumes\":[{\"assetId\":\"f8e72b8244738d0b448e46b38c5db6c2\",\"name\":\"costume1\",\"md5ext\":\"f8e72b82447" +
	"38d0b448e46b38c5db6c2.svg\",\"dataFormat\":\"svg\",\"bitmapResolution\":1,\"rotationCenterX\":67,\"rotationCen" +
	"terY\":95},{\"assetId\":\"bb3a866c4db08353f6faf43c54990f10\",\"name\":\"costume2\",\"md5ext\":\"bb3a866c4db083" +
	"53f6faf43c54990f10.svg\",\"dataFormat\":\"svg\",\"bitmapResolution\":1,\"rotationCenterX\":67,\"rotationCenter" +
	"Y\":95},{\"assetId\":\"bb6b82c9fa7c432c552ca2f251ae2078\",\"name\":\"costume3\",\"md5ext\":\"bb6b82c9fa7c432c5" +
	"52ca2f251ae2078.svg\",\"dataFormat\":\"svg\",\"bitmapResolution\":1,\"rotationCenterX\":67,\"rotationCenterY\"" +
	":95},{\"assetId\":\"be2345a4417ff516f9c1a5ece86a8c64\",\"name\":\"costume4\",\"md5ext\":\"be2345a4417ff516f9c1" +
	"a5ece86a8c64.svg\",\"dataFormat\":\"svg\",\"bitmapResolution\":1,\"rotationCenterX\":67,\"rotationCenterY\":95" +
	"}],\"sounds\":[{\"assetId\":\"28c76b6bebd04be1383fe9ba4933d263\",\"name\":\"Beep\",\"dataFormat\":\"wav\",\"fo" +
	"rmat\":\"\",\"rate\":11025,\"sampleCount\":9536,\"md5ext\":\"28c76b6bebd04be1383fe9ba4933d263.wav\"}],\"volume" +
	"\":100,\"visible\":true,\"x\":0,\"y\":0,\"size\":100,\"direction\":90,\"draggable\":false,\"rotationStyle\":\"" +
	"all around\"}],\"meta\":{\"semver\":\"3.0.0\",\"vm\":\"0.1.0\",\"agent\":\"Mozilla/5.0 (Macintosh; Intel Mac O" +
	"S X 10_13_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.181 Safari/537.36\"}}"
