package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	config "github.com/duyhtq/incognito-data-sync/config"
	pg "github.com/duyhtq/incognito-data-sync/databases/postgresql"
	"github.com/duyhtq/incognito-data-sync/models"
)

type ExchangeRateResp struct {
	Data struct {
		Quote struct {
			USD struct {
				Price float64 `json:"price"`
			} `json:"USD"`
		} `json:"quote"`
	} `json:"data"`
}

type PToken struct {
	TokenID  string `json:"TokenID"`
	Symbol   string `json:"Symbol"`
	Name     string `json:"Name"`
	Decimals int    `json:"Decimals"`
}
type PTokenResp struct {
	Result []PToken `json:"Result"`
}
type CronUpdatePrice struct {
	running bool
}

// NewServer is to new server instance
func NewServer() (*CronUpdatePrice, error) {
	return &CronUpdatePrice{}, nil
}

func main() {

	// exchangeRateCron := &CronUpdatePrice{
	// 	running: false,
	// }
	conf := config.GetConfig()

	db, err := pg.Init(conf)
	if err != nil {
		return
	}

	pTokenStore, err := pg.NewPTokensStore(db)
	if err != nil {
		fmt.Println("err", err)
		return
	}

	// tokens, err := tokenStore.ListTokenIds()
	// fmt.Println("errssss", err)
	s, err := NewServer()
	if err != nil {

		fmt.Println("asad", err)
		return
	}

	var opts = []struct {
		TokenID  string
		Symbol   string
		Decimals int
	}{
		{"ecdb1ba76eaa62a08f655bcc6b253f0203201c34786db78e67be20a1e0cdb98d", "0xECH", 1.00E+08},
		{"889d11f1affba816b7d1e1819c3f1c3974ceed1bd49a7f43622e8f24c6412ea6", "0xESV", 1.00E+08},
		{"37da535a32efe21279a5b992a56066cc122de903c51667233ea694749ab9fe96", "0xETC", 1.00E+08},
		{"6e16e257b53014713e89cfc3bff34b60093e76be7783ea313c6b30d8fa68518f", "0xETD", 1.00E+08},
		{"7f6c7544854195919c91c5e1f1920c75abe8a97948b270196caf2664c2d21fa0", "3W", 1.00E+09},
		{"d2b1dc4fbd2f6d5a684010ecc32cb4a34b627ac2c6e93654b597febd4226a55e", "ADI-6BB", 1.00E+09},
		{"391617b4960fd9d24b318fccd9e47b25b3ac07d04228be837f421610cbb96b0e", "AE", 1.00E+09},
		{"0b3c4aecfcb47b25ed491f047abe510f8c97872215d6765ff83e436fa533f011", "AGI", 1.00E+08},
		{"2dda855fb4660225882d11136a64ad80effbddfa18a168f78924629b8664a6b3", "BAND", 1.00E+09},
		{"1fe75e9afa01b85126370a1583c7af9f1a5731625ef076ece396fcc6584c2b44", "BAT", 1.00E+09},
		{"db8b15f2b8d8921dd9ba1ecb7efd5662b9863a4cd39922149a56b827081b5289", "BET-844", 1.00E+09},
		{"5175e14c7b2284d69e71a3aa358e6d550b59c27154d7297331891717347261a2", "BETX-A0C", 1.00E+09},
		{"b416aeb914a50df0c234de231de51b8e2e87efe7e65b2881bc8edb58b6c295db", "BGBP-CF3", 1.00E+09},
		{"94445ea09e116a0f1ca4dbb13278d98f59a62a4fa68f0a1ad1cbc608361ad9b8", "BHFT-BBE", 1.00E+09},
		{"0240b50b180fb315e59a20db45cc5d86e6c1605dfedc6e24d009de029bd71e93", "BLC", 1.00E+08},
		{"b2655152784e8639fa19521a7035f331eea1f1e911b2f3200a507ebb4554387b", "BNB", 1.00E+09},
		{"df9f96a3f47980ef1d01eae445cad28892cf585f3ed0d32bf6d947edc733d513", "BNK", 1.00E+08},
		{"6dee66a817e8ddc321097bd9c58046e7ccee79ad2e3693ac9645caef16494740", "BNT", 1.00E+09},
		{"7880435eec031b5846be1d40200ff6e97fcc5d863ebba07c9d7f73d477a28125", "BPRO-5A6", 1.00E+09},
		{"b832e5d3b1f01a4f0623f7fe91d6673461e1f5d37d91fe78c5c2e6183ff39696", "BTC", 1.00E+09},
		{"0ffbd46aab0099a52c469225b84df7c7127fc1521db0b57ad23c92d52587e106", "BTM", 1.00E+08},
		{"904fed9074d6192a25829f09f7d284f3f93cb55555c2aa960eb16c42526e317d", "BTU", 1.00E+09},
		{"2c22eeb3811b4c5760d2fea6f26824de253aa18fce0a3c8cb440d655d4688ef4", "BTV", 1.00E+08},
		{"79e629030cfd3e27c740129b363c495f033a802c84b06e3e35da04de2e4b63f2", "BULK", 1.00E+08},
		{"5f9b9480d825e276d8413c52e728b6b5fe1dbbbe0e6ba0923a89e33e04d0d01c", "BULK", 1.00E+08},
		{"7b764ae9c1b190ef1af0d9e6f5a77eda261bd564f5479484f7cfe1b0f746218d", "BULK", 1.00E+09},
		{"9e1142557e63fd20dee7f3c9524ffe0aa41198c494aa8d36447d12e85f0ddce7", "BUSD", 1.00E+09},
		{"6a6737e6cc4e8fa15959cb760c95753c9bb6b02117ceff9d86545b427fec4a9b", "BUSD-BD1", 1.00E+09},
		{"aa90dbeb4eb1469f095d79fbdb065f556ac491dba06a7ce227b6477a72b83f4c", "BZNT-464", 1.00E+09},
		{"208f60d39ff90dc12ba7ead9a549ad1da220531d690d17ac3d5d1c800797fb3b", "CCCX-10D", 1.00E+09},
		{"58ab8b4f60056161fdc1cdbdbcf1bbc2d9808ca3337412f828f188f1c99081e4", "cDAI", 1.00E+08},
		{"f30a8d9814d1e4d54e20dcfded7cdda8ed183ce8268ad9edaec67d03327afa66", "CELR", 1.00E+09},
		{"44e7d4e6dadcebbdc3dc4897029392b67afdc3d9ca59cb22f423a73a8c85d245", "COCOS", 1.00E+09},
		{"00d9e8811e193f42d7d4a6d7b8a24ded68c04aed3550295f0795fe3af20d0c9f", "CON", 1.00E+09},
		{"23a0b532315df303603fdaeefd7ce2712a5a115d3de8f6e9db70b80828ba1328", "CPTL", 1.00E+08},
		{"ed1ddc6f38a56d815c93b1b18784d8d96722d3480b84ce3254f36e7f8d3793ec", "CRO", 1.00E+08},
		{"dfc4f4e9b6b9356f2b7221475456a443ea4958936ae965d13dc980ec5267c063", "CVC", 1.00E+08},
		{"3f89c75324b46f13c7b036871060e641d996a24c09b3065835cb1d38b799d6c1", "DAI", 1.00E+09},
		{"c9914d8edf0b460cc0fb875988c83d5c637176bf4a56e9ced6ece5788a42b24a", "DATA", 1.00E+09},
		{"a455983046eacf5c22235496230e0eac0042fd88fc97c0cf20ae59853cd51da3", "DGX", 1.00E+09},
		{"29d145e1e25de67fecb71a98c13f7f649fd616e80c23954bfe4d4fc01cea38e9", "EBK", 1.00E+09},
		{"fe4bf3890c8f3acf1f2e3e0cc7086955ef5b7697f8abcb6ee727bbb54bcaede5", "EDO", 1.00E+09},
		{"11af39fe6faac9241c093cb5e8700c2fb51cc9942a3254b3277160359b484e0f", "ENG", 1.00E+08},
		{"87b1e95cd5eefad03509c02f232660e7214ce10f8bb398e5cfa3a94ad5b1e3e9", "ENJ", 1.00E+09},
		{"4f08d200ff701d8b41d875ad799566bf4be68b84275d7625728c85d0a7d7bad9", "ENQ", 1.00E+09},
		{"ffd8d42dc40a8d166ea4848baf8b5f6e912ad79875f4373070b59392b1756c8f", "ETH", 1.00E+09},
		{"3c730bd6d20290e86da985dc99457f635717993960595095907f413360ab039e", "EVEN", 1.00E+09},
		{"be46b272e40f2bc594b146ee610020ce368cd7a8297bf45993cb25e05a25619b", "FET", 1.00E+09},
		{"962edee91bd0aa3ba9e15b03e218b0bd58d8dea6e2ffc8d1afaec59440a84e8f", "FLX", 1.00E+02},
		{"d09ad0af0a34ea3e13b772ef9918b71793a18c79b2b75aec42c53b69537029fe", "FTM", 1.00E+09},
		{"6bcf0666df8d2b88968b7718a61f91d754099c21087ade086ae585f4ffc45f10", "FTX Token", 1.00E+09},
		{"102c9349ed9765503053eb3b06bf4f118bab2f5845b5f3e3751976300f473323", "FXC", 1.00E+09},
		{"bfeb79c3c26ffaf5cb30a3659641238ef85b828d9731f46a41c8d47c618c90f7", "FXC", 1.00E+04},
		{"285a5eacb9ab9c816bb947ad4b65bfad0bc616b032cdd7f697715f11474df8a7", "FYN", 1000000000},
		{"445f70fe7a29bf71f3b5c28260c3497be27d80e0b09d0e9f351222721e26fc4b", "GIGS", 1000000000},
		{"8c2c35d207f2746a97b8220ca3ca32086fe0d661e9a2b0cfc86aa126e11fb1a0", "GKRW", 1000000000},
		{"e16389a9a1f58d50eaf6c4e8315566c22936174da8a0f98c37b5d7138eb07e8d", "GNT", 1000000000},
		{"465b0f709844be95d97e1f5c484e79c6c1ac51d28de2a68020e7313d34f644fe", "GUSD", 100},
		{"9ba01cbf9a2c9d9b49889312d47884603d010f427464f0dc2051e4da2be5a44e", "HOT", 1000000000},
		{"849404959d569b592e3a2622b1b9327cd3e3c616686f450fbef419a95b227821", "HOX", 1000000000},
		{"85408253fb6929ab526a05bad979e4910627d351ed891ef95a3c9c267736d585", "HUSD", 100000000},
		{"36d0456d98690a59beda78489efb997fba7e0710c252a9a657c94506ca7f1748", "HYN", 1000000000},
		{"97a03e552b52f79c840c6339bc64b4af1afcc69284a82d940073881d5778229b", "IDN", 100000000},
		{"638faeb96f33446caa505c70f1e8383a9d8304d7806d31bd84505673c7bddb18", "IDRTB-178", 1000000000},
		{"050156975a56e7b3a381e916c7fe0617659ce74f378757819fe7f47761489161", "INO", 1},
		{"e7daee0e9b6699cee24d8e5195f3cb5aaee7856419fcaa4b914f2db8a90122de", "IOST", 1000000000},
		{"fd3ca9ef63e25aae27bdcf27a56175d31f79682a61d1b45765550267912160d2", "IOTA", 1000000000},
		{"421e49911ff2bd47f4d0b0747dbb9e724965a8726c3040d2a74b74aca390fed9", "JSE", 1000000000},
		{"f1efcb6ebab30cd61415a45217b9b8ad3151a10078e5b0e9b83fed5b183e389f", "JWEL", 1000000000},
		{"1376f89b42184b7009afc40b4e51e42f762d14aca8e10dcdc0def0a9b8c69079", "KAT", 1000000000},
		{"5271886b40a5125f00cd5f2f0ea78731f42ba26fadf10c4aef1b7793aed193ce", "KAT-7BB", 1000000000},
		{"b6533b4f5fc8a61c36cd98263aad628cd2b14983d9ffd381b3656be87a0ac68e", "KAU", 100000000},
		{"ae59e1d68bb96de15936383fb0b7d4d1595ccf975a41c82dc93d933004d9f7ea", "KAVA-10C", 1000000000},
		{"513467653e06af73cd2b2874dd4af948f11f1c6f2689e994c055fd6934349e05", "KCS", 1000000},
		{"1a072deb34b87f0e828534ffeccb975c1de5c3a72da1cb01f11a9951d9f7a495", "KEY", 1000000000},
		{"2ea778c817e1fe8b2194fb6906dd0992c849dc807e784e46369f28c0d4b269ff", "KNC", 1000000000},
		{"0274c7705fadb265e32fadeac2cb93dc26ebfdf56048a18eb891a0114d0cd173", "LBIT", 1000000000},
		{"e0926da2436adc42e65ca174e590c7b17040cd0b7bdf35982f0dd7fc067f6bcf", "LINK", 1000000000},
		{"71ae5930b2fe30a4d9e850279ea231698f031d65c2862108b6c1101412d2e0b1", "LMY", 1000000000},
		{"298286913b236f4eadd6a5154e78882dbba8466b6a213ae35fd448404126bca3", "LRC", 1000000000},
		{"dae027b21d8d57114da11209dce8eeb587d01adf59d4fc356a8be5eedc146859", "MATIC", 1000000000},
		{"caaf286e889a8e0cee122f434d3770385a0fd92d27fcee737405b73c45b4f05f", "MCO", 100000000},
		{"d2b20be6e8ea55e6ae24f15d0614f9f3999ed880efdbfd8697291a7d389944f6", "MUA", 1000000000},
		{"9e2b1632e72dddc0d2bef791de93a589c2572940a8b3428719947bfcfd0a6e77", "NEXO-A84", 1000000000},
		{"c9e032e7c236b31e37eac238e61e555e0ae6745b63bc0c961159185dbf800176", "OCEAN", 1000000000},
		{"613c80968b8270808380be723a59107dd9745a8f9580e8cf268ce7310d07d867", "OKB", 1000000000},
		{"249ca174b4dce58ea6e1f8eda6e6f74ab6a3de4e4913c4f50c15101001bb467b", "OMG", 1000000000},
		{"0369ef074eee01fe42ce10bcddf4f411435598f451b31ad169f5aa8def9d940f", "ONE", 1000000000},
		{"4077654cf585a99b448564d1ecc915baf7b8ac58693d9f0a6af6c12b18143044", "ONE-5F9", 1000000000},
		{"6c3a3f4f32618118502f1e852ad39704eb1dd5c9388101617f3a5b8fc4cc6f31", "OVF", 1000000000},
		{"4a790f603aa2e7afe8b354e63758bb187a4724293d6057a46859c81b7bd0e9fb", "PAX", 1000000000},
		{"9e0058816c7c5696606dafd0b708785b685cd881229dd22139490eca179a9dfa", "PHM", 100},
		{"2aba7867047af73597c1952beda00e5dfb7d2dfe60f83c498e3731d0e8d19d29", "POWR", 1000000},
		{"27125ca4a6db9463ea67d745498b6e738ad1ffcce9d2678bc34ee6e10965b380", "RAVEN-F66", 1000000000},
		{"977c9d8b6109ce84a880a5d56d629fe4752cc6621f393fa1d01bca0c4ed7ed85", "REN", 1000000000},
		{"068905f717ef2e047f437ab9451de898558298aacd824d96e9b594a91be2d02e", "RHOC", 100000000},
		{"acb93edf5eedfa6bff48d53d10de5db1927c0ff01f46f7555e7ed5c6e8e025a9", "RSR", 1000000000},
		{"e23b22c0b6b95cd6785ccdd2cbd88e88083110bbac168eb6cd702279679de732", "RSV", 1000000000},
		{"d240c61c6066fed0535df9302f1be9f5c9728ef6d01ce88d525c4f6ff9d65a56", "SAI", 1000000000},
		{"7189e556a2f61340c51fe346da5876e8528b311e0af7fc2a773bf442c9147eda", "Seele", 1000000000},
		{"28fecd1e4b7aef27283fc3f7928e70359d9e4537a4823656f00174ac86a96a92", "SHIT", 1000000},
		{"7df7011bd3051b7713ca8dd227a0092f6f62d09f021df8c6e75a969920690d20", "SNX", 1000000000},
		{"2dc746750e91248c384d01e71665245b24dfcc2ebe605d2bb75dd0e8f1dc6bc1", "SWAP", 1},
		{"8c1f15200a4c7542e6c23da98f18930b7bd96e26763e344056924cd73172d1ac", "TAOL", 1000000000},
		{"af202a350df3e3c1bf0b424913769f80a185d6d6a60387f3cf121e0f6c9ba310", "TAU", 1000000000},
		{"ad7f3d0767f323fbc38fb810738666e0af6dd90883e87b4b1a949570b63d7c63", "TAUD", 1000000000},
		{"f99fcd95f552afd1330c8b4b954bfa0cf86809d88d9f18e44601902c966f466a", "TCAD", 1000000000},
		{"2097942c05ea056e56198c262dfe0d90a014e8bfd8f7f84eba4a6ee6938ffb1a", "TGBP", 1000000000},
		{"4acbf0c603f9880834c71cbdb7f74d5344ecd9a5c9b9c3d4a29b27adb864beb1", "THKD", 1000000000},
		{"c2a63f8484f6f8fd5740360781f4546c5bdedeb7347ff1f9df08ae9da2f32447", "TNX", 100000000},
		{"a0a22d131bbfdc892938542f0dbe1a7f2f48e16bc46bf1c5404319335dc1f0df", "TOMO", 1000000000},
		{"6f00ae1afa3eb1bd043fa8be348888e0a0da26933d82113c11fed9a0938d92f0", "TRB", 1000000000},
		{"8c3a61e77061265aaefa1e7160abfe343c2189278dd224bb7da6e7edc6a1d4db", "TUSD", 1000000000},
		{"c8a90c2cd895e2c9ffbb07812f397c35105472ccabf1d442407362e3cc27dce8", "UND-EBC", 1000000000},
		{"7efd6c3db16cd03398d56d6c6ae9cd20af5799f88d011ed994b7890a09415a96", "Uniswap:WBTC", 1000000000},
		{"b20810f4d2a1dde8046028819d9fa12549e04ce14fb299594da8cfca9be5d856", "USD", 1000000000},
		{"1ff2da446abfebea3ba30385e2ca99b0f0bbeda5c6371f4c23c939672b429a42", "USDC", 1000000},
		{"8d1645eb82ba1fa56d16a4b5350719decd61dbd343a700e01730492790be121d", "USDS", 1000000},
		{"716fd1009e2a1669caacc36891e707bfdf02590f96ebd897548e8963c95ebac0", "USDT", 1000000},
		{"622cf9a19bbe1ca9e29d30ace38c6639dbcb74a819a1f6a735fb150fc709f2c1", "VEN", 1000000000},
		{"247602bd5fcb0bac41274ef6edee16a7e6118828d3e8948960b19c5591dfcf5b", "VIX", 1000000000},
		{"ce98e782cf31f35807b7c3a73cf9b7271b30feb371232ad3a0178ad420e83673", "VNC", 100000000},
		{"89d88b3cc3a97ea2cf9530bbb3ea9b9203c1e398d4ac084a889958e855522e5f", "VYA", 100000000},
		{"cce982d217198e9a6f3e8a37ce59cbd47c6e5dd6f57bce09e2ebca1d7cde6008", "W0xETH", 100000000},
		{"e993fbade7d3df5c14c944c27b1cdac9d80548c2132fb5adc42490b1e8de8098", "WABI", 1000000000},
		{"d3b13aa763492461074224cef0e469b89f144d3993683c743a927ba057261b68", "WIC", 100000000},
		{"43389a9e391c0462b960cbd3b7ff4c5d8e8b19bd294caa1e2dce61d6972c634b", "WISH", 1000000000},
		{"e32f35d0db21b2538172d2979c4b0e8f2d19156a775550c62bd73368c77f74ab", "WISH-2D5", 1000000000},
		{"538ed5effbab40822f55bb0942485b59538561943ca08b5755dc7aaf777d8b06", "WRX-ED1", 1000000000},
		{"530cd74f506edcdd263e34e6dacdd15097f87677036cf412f8ebeb1c494e352d", "WTC", 1000000000},
		{"0566f05c549a8062e6e5ff5b84b2e8fb1b2d71088f14232b15c831b49ff66e59", "XAK", 1000000000},
		{"7d731a220b13c7cc2f3ffb96834adbbe8352f6596805b58c8ceb0ff3a97c2587", "XMD", 100000000},
		{"c01e7dc1d1aba995c19b257412340b057f8ad1482ccb6a9bb0adce61afbf05d4", "XMR", 1000000000},
		{"ef23e8d0b30df5b3d98d42c25a95335768bc1efb562c3a7229b34e2d8bceb304", "XYO", 1000000000},
		{"880ea0787f6c1555e59e3958a595086b7802fc7a38276bcd80d4525606557fbc", "ZIL", 1000000000},
		{"de395b1914718702687b477703bdd36e52119033a9037bb28f6b33a3d0c2f867", "ZRX", 1000000000},
		{"9f52e4a54759a6dbc8ddabe091ec0cc34139416baf8624781391fefe2f2d5860", "ZUM", 100000000},
		{"0000000000000000000000000000000000000000000000000000000000000004", "PRV", 1000000000},
		{"4d6b13becb876480de61b37622af248dfcffd41ddb174282994d558151e57976", "EOSBULL", 1000000000},
		{"1b694fc8b5825681dc9249cc566c6926e92130e3375b3b1e6a1e14acba9d38d5", "HEX", 100000000},
		{"86c45a9fdddc5546e3b4f09dba211b836aefc5d08ed22e7d33cff7f9b8b39c10", "NEO", 1000000000},
	}

	for index, token := range opts {
		remainder := index % 25
		if remainder == 0 && index != 0 {
			time.Sleep(1 * time.Minute)
		}
		exchangeRateResp, err := s.getDataFromApi(token.Symbol)

		if err == nil {
			pTokenModel := models.PToken{
				TokenID: token.TokenID,
				Name:    token.Symbol,
				Symbol:  token.Symbol,
				Decimal: token.Decimals,
				Price:   exchangeRateResp,
			}

			err = pTokenStore.StorePToken(&pTokenModel)
			if err != nil {
				fmt.Println("err:", err)
			}
		} else {
			fmt.Println("err:", err)
		}
	}
}

func (c *CronUpdatePrice) getDataFromApi(currency string) (float64, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v1/tools/price-conversion", nil)
	if err != nil {
		log.Print(err)
		return 0, nil
	}

	q := url.Values{}
	q.Add("symbol", currency)
	q.Add("amount", "1")
	q.Add("convert", "USD")
	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", "23445d58-34c4-4e81-8ee6-964693c71f2f")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request to server")
		return 0, nil
	}
	if resp.Status != "200 OK" {
		return 0, nil
	}
	respBody, _ := ioutil.ReadAll(resp.Body)

	var exchangeRateResp ExchangeRateResp
	if err := json.Unmarshal(respBody, &exchangeRateResp); err != nil {
		fmt.Println("error Unmarshal body ....", err)
		return 0, nil
	}

	return exchangeRateResp.Data.Quote.USD.Price, nil
}
