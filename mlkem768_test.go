package mlkem768

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"math/big"
	"testing"
)

func TestFieldAdd(t *testing.T) {
	for a := fieldElement(0); a < q; a++ {
		for b := fieldElement(0); b < q; b++ {
			got := fieldAdd(a, b)
			exp := (a + b) % q
			if got != exp {
				t.Fatalf("%d + %d = %d, expected %d", a, b, got, exp)
			}
		}
	}
}

func TestFieldSub(t *testing.T) {
	for a := fieldElement(0); a < q; a++ {
		for b := fieldElement(0); b < q; b++ {
			got := fieldSub(a, b)
			exp := (a - b + q) % q
			if got != exp {
				t.Fatalf("%d - %d = %d, expected %d", a, b, got, exp)
			}
		}
	}
}

func TestFieldMul(t *testing.T) {
	for a := fieldElement(0); a < q; a++ {
		for b := fieldElement(0); b < q; b++ {
			got := fieldMul(a, b)
			exp := fieldElement((uint32(a) * uint32(b)) % q)
			if got != exp {
				t.Fatalf("%d * %d = %d, expected %d", a, b, got, exp)
			}
		}
	}
}

func TestDecompressCompress(t *testing.T) {
	for _, bits := range []uint8{1, 4, 10} {
		for a := uint16(0); a < 1<<bits; a++ {
			f := decompress(a, bits)
			if f >= q {
				t.Fatalf("decompress(%d, %d) = %d >= q", a, bits, f)
			}
			got := compress(f, bits)
			if got != a {
				t.Fatalf("compress(decompress(%d, %d), %d) = %d", a, bits, bits, got)
			}
		}

		for a := fieldElement(0); a < q; a++ {
			c := compress(a, bits)
			if c >= 1<<bits {
				t.Fatalf("compress(%d, %d) = %d >= 2^bits", a, bits, c)
			}
		}
	}
}

func BitRev7(n uint8) uint8 {
	if n>>7 != 0 {
		panic("not 7 bits")
	}
	var r uint8
	r |= n >> 6 & 0b0000_0001
	r |= n >> 4 & 0b0000_0010
	r |= n >> 2 & 0b0000_0100
	r |= n /**/ & 0b0000_1000
	r |= n << 2 & 0b0001_0000
	r |= n << 4 & 0b0010_0000
	r |= n << 6 & 0b0100_0000
	return r
}

func TestZetas(t *testing.T) {
	ζ := big.NewInt(17)
	q := big.NewInt(q)
	for k, zeta := range zetas {
		// ζ^BitRev7(k) mod q
		exp := new(big.Int).Exp(ζ, big.NewInt(int64(BitRev7(uint8(k)))), q)
		if big.NewInt(int64(zeta)).Cmp(exp) != 0 {
			t.Errorf("zetas[%d] = %v, expected %v", k, zeta, exp)
		}
	}
}

func TestGammas(t *testing.T) {
	ζ := big.NewInt(17)
	q := big.NewInt(q)
	for k, gamma := range gammas {
		// ζ^2BitRev7(i)+1
		exp := new(big.Int).Exp(ζ, big.NewInt(int64(BitRev7(uint8(k)))*2+1), q)
		if big.NewInt(int64(gamma)).Cmp(exp) != 0 {
			t.Errorf("gammas[%d] = %v, expected %v", k, gamma, exp)
		}
	}
}

func TestKeyGenVector(t *testing.T) {
	// From the October 2023 version of PQC Intermediate Values available at
	// https://csrc.nist.gov/Projects/post-quantum-cryptography/post-quantum-cryptography-standardization/example-files
	// Key Generation -- ML-KEM-768.txt

	// Note that d == z in the vectors, which is unfortunate because—aside from
	// being confusing, as this would NOT be possible in practice—it makes it
	// impossible to detect e.g. a mistake swapping the two.
	d, _ := hex.DecodeString("92AC7D1F83BAFAE6EE86FE00F95D813375772434860F5FF7D54FFC37399BC4CC")
	z, _ := hex.DecodeString("92AC7D1F83BAFAE6EE86FE00F95D813375772434860F5FF7D54FFC37399BC4CC")
	ekExp, _ := hex.DecodeString("D2E69A05534A7232C5F1B766E93A5EE2EA1B26E860A3441ADEA91EDB782CABC8A5D011A21BC388E7F486F0B7993079AE3F1A7C85D27D0F492184D59062142B76A43734A90D556A95DC483DD82104ED58CA1571C39685827951434CC1001AA4C813261E4F93028E14CD08F768A454310C3B010C83B74D04A57BB977B3D8BCF3AAA78CA12B78F010D95134928A5E5D96A029B442A41888038B29C2F122B0B6B3AF121AEA29A05553BDF1DB607AFB17001860AF1823BCF03DB3B441DA163A28C523A5FB4669A64234A4BCD1217FF2635BD97680FF938DBCF10E9532A9A79A5B073A9E8DB2123D210FAEA200B664838E80071F2BA254AAC890A46E28EC342D92812B01593071657E7A3A4A75CB3D5279CE88405AC5ADACB2051E022EE0AC9BBFE32DEF98667ED347ADCB3930F3CAD031391B709A4E61B8DD4B3FB741B5BD60BF304015EE7546A24B59EADCA137C7125074726B7686EC551B7BC26BBDB20FC3783534E34EE1F1BC6B77AB49A6667846975778C3C536830450A3FA910259722F3F806E6EB4B9346763FEF0922BC4B6EB3826AFF24EADC6CF6E477C2E055CFB7A90A55C06D0B2A2F5116069E64A5B5078C0577BC8E7900EA71C341C02AD854EA5A01AF2A605CB2068D52438CDDC60B03882CC024D13045F2BA6B0F446AAA5958760617945371FD78C28A40677A6E72F513B9E0667A9BAF446C1BA931BA81834234792A2A2B2B3701F31B7CF467C80F1981141BB457793E1307091C48B5914646A60CE1A301543779D7C3342AD179796C2C440D99DF9D41B52E32625A82AA5F579A9920BFFBA964FA70DB259C85E68C813817B1347BF19814DA5E9364A4645E621923D955C211A55D355C816DA04730AA324085E622B51D6109B49F673ADD00E414755C8024AA0164F24556DED963D61143856CB4FF0567E3320730DBCBF12F66E2B70B20054A6DEA42614B50EF72B156F5149FC263DD7E039C55A3EE9827DF92C565D24C55E0A81C6494695344D948748AFBA9F762C0EA90BB724897902000775613949602C48C78A9440678C24086D326D79643BAF7036C66C7E026AAEFDA2807A60BD7FC91363BB0234A590984AA011F11D40268218A1588377B3D7671B8B99789919B86EE82B18EC22D4E80A1F27853D889419D460DEF7567AA4567969C43048C32B8462A9C9386EB3152A6976AA783CDD1A8C57A9B6BBD837A00624B58B4BA3DBB63BB8200E7BC88881BEBDA925BCA028E291AA1C22539CD04F90090D7F74108C32B8022C1591C881E76304E2408190E20F09A54FC23420E2620E9D87A3108A94FEEA72D5AB7FCFB972E6561B1A7B062F1A682E020AA2562812B296547B917824CDB88C582B5A6890177BC70C91ACAC9ABE290AEB2C34A7E2368955CB456A345368ABE3B91B47FC30B0233A09BA79FB11238AC508CCE61095F854C23204A8D36BFC2C6E05A72AF5244B17C12101E01451570EB110567E850E79C000142441FE4160027545F6290E85451B80234A9406C390B0CEA3C8335D4C6F8550B544C9343E61BA1C8489D1B0399739168AF740A481B0F5C3372530CA06B508ECE838AB78BEE1E597A9B14F6AEC7A3BD1AA8D10BAC23B9802902CD529AB6EF54DB3110CFB561E7E6948E65281250416C349C8100B3B4D3D0F62ACAD8D161175B134F7564937CD")
	dkExp, _ := hex.DecodeString("19D74AD5472A8B2BAAD2A56702C9B3B5510EF3924858061D57F90DD9A1A01FEC2F57C51A888805341B617C515539597750835C3ED7A033B039D72491332C5DF4A69B6DF26171877AD1E50AC50100BE4728786685DA7A739E843FF0D45922D7281E210D5E82B944652F4862CFB3D902DE60AFD0A164471B26144A1D7A38096503095911762EBA7962C4511D05A128F2781ECB3D1F5BB1244237611ABAB924991F8A2732E27032357920F197C7692D60A9444472258CB457C1B71B77995469F3A962F3ABA6699614FCCCEA741E21C600C4357BBFAB452927C3D441BF8ED73152F75C08F540E186ACCA3326F422C84B988D77E61AE61859CF8541F89209E4983040C5617654808852B649B899A399AEC2C8BBA8A542F345ABF2813F65E9A791D32CC2D76026FB8D0C94B657489ABB487DA4A2C0E3868D3CF47F1CBB2FA79C53CFF6264777C09B177C91315484D2B30B0CA21F55ADD23C57E1911C3F086BCAD21798486EB47B7C58577381C09F5252582D1B27A7D5B8E060CE78209CC82BAE4DA606800C8DB1268F7AD2B793A44F34612CCEA31CE7D796A65A2691D61500625F83E7BE57077EE9C1B8C1CAA137CC4B6573308C19668B24B01E966903ABBCB79B67BE0A3E3E058AADA189B9EA80359AC26F4C5C53735FE4FC35247337760CCA3529B8D266BB6C48010654CDBC5A3E9757524675ABC413130CC2701F28933EABB8392B0D6D059CFC3A30326C4FCC810B37A4748C1C53928A4913E48B186697162C33FFFB06DD5161C8639DB195C6CA64829B2B3A2E4C9683B66DF7FB1909904E00020DBA134E02A168D76AC076BB77D4DC8496B4BBE7B4690BA29B62A91ABE72BEF323A44C8903E482B60D99BA61D1BBCF9CB9673534C1D647662374EE2C7C5F0081BAD149F44206717684D9746B2048633AF7A68C6865FB590358D8CF821458369B0C31EB597CF5BE78EB480EA04E35FACC380372C8C0A04DE276B1A72121E596CBB25EF7536AD3804184A87BDFB5A769160BFBB0CA3C360790E5562BB78EFE0069C77483AD35CAC237C61DE78A7DB46FC917124CA17510DB7DA218890F448EF6318613A1C97C928E2B7B6A54617BCCB6CDF278AE542B56AD7BB5ECD8C46A66C4FA0950CE41352CB85711890458F299BF40BA6FF2C0713862268B5F08E49845B09443997AB29A62073C0D9818C020167D4749231C059E6F483F976817C90C20A9C937079C2D4BE30DA974A97E4BC53ED96A55169F4A23A3EA24BD8E01B8FAEB95D4E53FFFECB60802C388A40F4660540B1B1F8176C9811BB26A683CA789564A2940FCEB2CE6A92A1EE45EE4C31857C9B9B8B56A79D95A46CB393A31A2737BAFEA6C81066A672B34C10AA98957C91766B730036A56D940AA4EBCB758B08351E2C4FD19453BF3A6292A993D67C7ECC72F42F782E9EBAA1A8B3B0F567AB39421F6A67A6B8410FD94A721D365F1639E9DDABFD0A6CE1A4605BD2B1C9B977BD1EA32867368D6E639D019AC101853BC153C86F85280FC763BA24FB57A296CB12D32E08AB32C551D5A45A4A28F9ADC28F7A2900E25A40B5190B22AB19DFB246F42B24F97CCA9B09BEAD246E1734F446677B38B7522B780727C117440C9F1A024520C141A69CDD2E69A05534A7232C5F1B766E93A5EE2EA1B26E860A3441ADEA91EDB782CABC8A5D011A21BC388E7F486F0B7993079AE3F1A7C85D27D0F492184D59062142B76A43734A90D556A95DC483DD82104ED58CA1571C39685827951434CC1001AA4C813261E4F93028E14CD08F768A454310C3B010C83B74D04A57BB977B3D8BCF3AAA78CA12B78F010D95134928A5E5D96A029B442A41888038B29C2F122B0B6B3AF121AEA29A05553BDF1DB607AFB17001860AF1823BCF03DB3B441DA163A28C523A5FB4669A64234A4BCD1217FF2635BD97680FF938DBCF10E9532A9A79A5B073A9E8DB2123D210FAEA200B664838E80071F2BA254AAC890A46E28EC342D92812B01593071657E7A3A4A75CB3D5279CE88405AC5ADACB2051E022EE0AC9BBFE32DEF98667ED347ADCB3930F3CAD031391B709A4E61B8DD4B3FB741B5BD60BF304015EE7546A24B59EADCA137C7125074726B7686EC551B7BC26BBDB20FC3783534E34EE1F1BC6B77AB49A6667846975778C3C536830450A3FA910259722F3F806E6EB4B9346763FEF0922BC4B6EB3826AFF24EADC6CF6E477C2E055CFB7A90A55C06D0B2A2F5116069E64A5B5078C0577BC8E7900EA71C341C02AD854EA5A01AF2A605CB2068D52438CDDC60B03882CC024D13045F2BA6B0F446AAA5958760617945371FD78C28A40677A6E72F513B9E0667A9BAF446C1BA931BA81834234792A2A2B2B3701F31B7CF467C80F1981141BB457793E1307091C48B5914646A60CE1A301543779D7C3342AD179796C2C440D99DF9D41B52E32625A82AA5F579A9920BFFBA964FA70DB259C85E68C813817B1347BF19814DA5E9364A4645E621923D955C211A55D355C816DA04730AA324085E622B51D6109B49F673ADD00E414755C8024AA0164F24556DED963D61143856CB4FF0567E3320730DBCBF12F66E2B70B20054A6DEA42614B50EF72B156F5149FC263DD7E039C55A3EE9827DF92C565D24C55E0A81C6494695344D948748AFBA9F762C0EA90BB724897902000775613949602C48C78A9440678C24086D326D79643BAF7036C66C7E026AAEFDA2807A60BD7FC91363BB0234A590984AA011F11D40268218A1588377B3D7671B8B99789919B86EE82B18EC22D4E80A1F27853D889419D460DEF7567AA4567969C43048C32B8462A9C9386EB3152A6976AA783CDD1A8C57A9B6BBD837A00624B58B4BA3DBB63BB8200E7BC88881BEBDA925BCA028E291AA1C22539CD04F90090D7F74108C32B8022C1591C881E76304E2408190E20F09A54FC23420E2620E9D87A3108A94FEEA72D5AB7FCFB972E6561B1A7B062F1A682E020AA2562812B296547B917824CDB88C582B5A6890177BC70C91ACAC9ABE290AEB2C34A7E2368955CB456A345368ABE3B91B47FC30B0233A09BA79FB11238AC508CCE61095F854C23204A8D36BFC2C6E05A72AF5244B17C12101E01451570EB110567E850E79C000142441FE4160027545F6290E85451B80234A9406C390B0CEA3C8335D4C6F8550B544C9343E61BA1C8489D1B0399739168AF740A481B0F5C3372530CA06B508ECE838AB78BEE1E597A9B14F6AEC7A3BD1AA8D10BAC23B9802902CD529AB6EF54DB3110CFB561E7E6948E65281250416C349C8100B3B4D3D0F62ACAD8D161175B134F7564937CDECE9E246AAD11021A67B20EB8F7765AC2823A9D18C93EC282D6DBC53CD6DF57592AC7D1F83BAFAE6EE86FE00F95D813375772434860F5FF7D54FFC37399BC4CC")

	ek, dk := kemKeyGen(d, z)
	if !bytes.Equal(ek, ekExp) {
		t.Errorf("ek: got %x, expected %x", ek, ekExp)
	}
	if !bytes.Equal(dk, dkExp) {
		t.Errorf("dk: got %x, expected %x", dk, dkExp)
	}
}

func TestEncapsVector(t *testing.T) {
	// From the October 2023 version of PQC Intermediate Values available at
	// https://csrc.nist.gov/Projects/post-quantum-cryptography/post-quantum-cryptography-standardization/example-files
	// Encapsulation -- ML-KEM-768.txt
	ek, _ := hex.DecodeString("1456A2EE8C3556054ABC79B4882C3190E5CA726AB402E5B09728C0F4F79C9FC2ADD828ABE432B1501B60F46CCBC86A3378C34895708A13671B20B389479AAA01C69D6B3B7D07D1C3AB54B91C580F5A336B30069A4F134FFD3764CE73A047E2844771742BF4710B972D4F6590A1C53A975368C271B670F1A4036441054A66E8815997512288552FD7149FFB705AAE133F8414060D0092FA8A1627D78AB2ABC6696288BAF5C60EF370827A7EFA72AE5C6741A5DA043D5940F121485372A98F472D60F05F74D95F01A1991E73A3E0A9536467A4738AB4CF385BA772827EB8CC058B3572E40B598444C181C7F6D9B760A7B907092E9C3351EA234E4449BD9B61A134654E2DA191FF0793961569D3594448BBC2586999A6671EFCA957F3A6699A4A1B2F4707ABA0B2DB20114FE68A4E2815AF3AAC4B8C6BE5648C50CC35C27C57288028D361708D302EEBB860BEE691F656A2550CB321E9293D7516C599817B766BA928B108779A1C8712E74C76841AC58B8C515BF4749BF715984445B2B53063384001E55F68867B1AF46CA70CA8EA74172DB80B5218BDE4F00A0E658DB5A18D94E1427AF7AE358CCEB238772FCC83F10828A4A367D42C4CB6933FDD1C1C7B86AD8B009657A96222D7BA92F527AF877970A83247F47A23FC2285118B57717715204674DA9C94B62BC7838CF87200156B26BA4671159931C49322D80671A0F332EAA2BBF893BE408B9EAC6A505483AA9075BD1368B51F99211F480A9C542A75B5BE08E43ADAF301DD729A85954010E64892A2AA4F15C0BD70B3D856494FF9BA0FE4CE12991CA06B5E3D0B2AF1F797B7A2B760910AE9F833D0D4267A58052C2990F161B886E251711C09D085C3D958B144192C9CC3224A460715B6784EB0B26F237187507D85C5110ACC71CE47198F254553356DAB448C38D243A7C02BE40C908C828D05C081DFAB8FC6B5CFE7D56E7317157DC053B2B3489986B081288871818585E09931095E3274A084115BE276438254A796270A7B4306F08B98D9C2AAECF7065E74446B7C696DBAAF8B4625A10B07827B4A8BABAB09B64AE1C375BB785441F319FB9AC2F14C95FFB252ABBB809C6909CD97706E40691CBA61C9252BD38A04311CA5BB2CA79578347505D0888851E082648BD003BE97C0F8F66759EC96A96A081C6822C4510559537042FC15F069A649B74A10961B354A1F625B04E25B293CF65FB4F53A80CC733D7A175775BF8A9ABB9201620E83A7F3E724D1287DBC44BDD5D85FC71545A927BEEDE537A7768735CC1486C7C3F31104DB67343F435D2D45554BAAC9CDB5822E8422AE8321C78ABE9F261FD4810A79E33E94E63B3341872C92253521997C084FBC060B8B125CCC88AC85AC5FE3168ACB059B3F119C4E050A20732F501BB9B3E687C846B5C2653F8886373E1004A2AB8D1BB970A7E571D8A46EE81B782F26942DD394FDD9A5E4C5631D985528604B1CC976275B6AC8A67CEEC10FFACBBA3D3BB141321DFC3C9231FC96E448B9AB847021E2C8D90C6BCAF2B1240783B62C79DEDC072A5763E660AF2C27C3F0C3C09207CAD990BB41A7BFCEC99F51596A0E83778F85C006AC6D1FE981B4C4BA1CB575A7D07AE2D31BA760095F74BC163841CF8FF77F894ABC6D261ED87A4530363B949C4AD24EFB3A56809478DDA2")
	msg, _ := hex.DecodeString("40BE9DCAC16E9CA73D49D0C83F9D3D89BB71574A4219A0F393DFECE2988394C4")
	ctExp, _ := hex.DecodeString("778D6B03791ACAF56CAAFCC78CEE5CBCA1DE8737E9C7FF4AE5F384D344E08223C74C824CB5848520517C7F0EA0645EB6F889517AE5216B0CF41DDC3F0D1DF9BC6E4DECB236A5EA8B214F64266D3CDE08E0CB00E5D91F586706B1EE533D20476F4423B78F916B1726EEEA959FFB9AC634D04A94D09923CB0D4E730CCA4144E7C4884921652DA4928C68E644F673CFC57D3E87CF5BE581A89F9CB8F0FCE2782D681E5CE88AF58458C3D63D807572DE5AA8E1FAF2DCD14EDB7349565B7D3271DDBEB0B6CC7AFE08635784311159733C46E5FDC5E0CD36CE5685ACFB1AFE50ABB46F447521E60D9C8F0E4CA28C190ABB40C365F412471E95A8EA396D4BD8070EEB1F02B07C825367AA1EC0F10C3862416BB21AD6CA748A86E9829EFC1A0499093C85176D37F574C75CF5EDFA8D920D3268CB34C6A4BB0002869BC05D7C8FCC0658D4A01EACD74557A37D98A763074752DFDD6429881CAFF577D3A048031BD52C4E9726398590F9519FD59405D6B3C307AFCB168A985785D954A6D1DC1EA92E1EB6F946A4D99DD6CA307ABFD8362FABA98BB264C69C5F555D60883CC56019FEB4E8000C48B7E68CD667F00B5250CEF293A4A9E778726E62F120361E21AB3140464CDC6ABDE9EA05198D8B3BB671B9111A2F317582847CA5015664F22CDB08C143187BDE2129B54F34160295D75FE9A494FD7E67AAA76B57AAFFD89D01A71DF5C8158620298D582BBEFA6D09AC412A99AA3BE9C383504948C43DD5AF4127B1435804F44BAFA142BFC2A95D95FB2EF0641ABE71064DE51D6B9EC50857B8EEF7F48036313D0E936763B8F7BDE69B064DD5761D80EA6F1A8B37565753C579BBB895EFB9FCB3FC5FA3362E3774F0F77140B973CAE587BAD2F3B566A9C25A969347E5C54F87F1105E9C074867D94077CCAE3ABEA54520EDB51D9DAABE7848E78FDF66E07E2E22B30251931E890BAF1F5E177D4D9CEC9E4969481FD7C1335A0ED5879F34EF4BB4F66C28803CEA162BA461506D52EB3AE16951922B06825186C3D4CE1B51F3C92F3C52F2D04D1F13B2B17C9EEB882CCE0EB88B7EA9A1CE4E37415CC84C7BC436A4628386CC77D9AFD207911BD9BFD8A7FA05C275BE0C4C6A8FC0A61BDA1D67AE33B5310BE1290DC71C1418EB5744BF2842C1652173A49A692E71FE43258A205B3CAAB90C0304A51E77D01B404A01FAE2F83AB80C5DBF6CF518C001F46A633FA169B1BDB77A9D0B1E0C007835C09F6ABBA96F3F53564DA508EE8861A483A81749D4A44672B1EF1605F29D168B74B736B4F13501D7AD1213118A7832E666A50BE8010D54322A526CF7A4E543A79D0D98E004FBEC76EA3F7E887BDBAF50DADFDDDF3FFECF6D3F77EA4B9B16DC754F4A68E5EF32F6A137E7C9E3C3E8C2E236C7EBC45D46EC1677A5A8BB2668443B0BE8693DC257F13D8B9A90100B92B4D1761B819673832C32020671BFB3D0220A363E4BED6D649D3F7368CFE081E196A43D4708798E31BB2A2F61824674ABA2FC9DCD05DB84B8627AE11488886F921BC79AE1FD03")
	kExp, _ := hex.DecodeString("616E0B753A3B7F40FEF9A389F58F16BFBB04622941D2464BDAE767820DFAC38E")

	ct, k, err := kemEncaps(ek, msg)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(ct, ctExp) {
		t.Errorf("ct: got %x, expected %x", ct, ctExp)
	}
	if !bytes.Equal(k, kExp) {
		t.Errorf("k: got %x, expected %x", k, kExp)
	}
}

func TestPKEDecryptVector(t *testing.T) {
	// From the October 2023 version of PQC Intermediate Values available at
	// https://csrc.nist.gov/Projects/post-quantum-cryptography/post-quantum-cryptography-standardization/example-files
	// Decapsulation -- ML-KEM-768.txt
	dk, _ := hex.DecodeString("3456859BF707E672AC712B7E70F5427574597502B81DE8931C92A9C0D22A8E1773CB87472205A31C32206BA4BCF42259533CB3A19C0200860244A6C3F6921845B0A05850187A4310B3D5223AAAA0C79B9BBCFCCB3F751214EB0CFAC1A29ED8848A5A49BA84BA68E6B6F5057D493105FF38A9F44B4E7F6CBE7D216408F7B48605B270B253B001A5401C0C9127CC185B1B0CF92B99FBA0D95A295F873515520C86321B8C966C837AAB34B2BFFAB2A2A4301B356B26CDC4563802901B4762F284281A382E5F762BEF47B519A81A108657EBE962BE120B5FB3B9ED338CCF47B3A03952A16633F6E6B534E6B63D05706EFA0F94C03A2B856AE551422F9011F2589A41B96A2CD213C6999B09E91FF423CB106A1A920B84B811469497154223987F005C72F8AF388B090C639F8C774FC5A294C74A212C91A86C328AEBEA558AB43F8B873534FA2EF9E66CEF3C52CD471AB78375E745B9D0AA65D2278B9275AE5348B16CF62AC8065734E4BD77B80CCF897605EB76F485AF8A0B466557A83C0292CCF903EE7AA57C3B51AD660189B86139E380425B31A92689DF2431BFA7B69EAB1727451B29DA8B8BF851E1BC2D3A63134CA9663C57AEC6985CEBD56DB0447B136B017A974761C3C67D33772F9964E5434D643504332A3027294A078C599CB29163109CE3B56CE698B4D3F59E2956A1F03A4B955593F2D2457FFAAE9624A0711045B3F55292F20CC9D0CD791A21597B0F2CD980F3510F0B0239022000D735586EE6A73F3A3DCBD6BD1A85C86512ABF3C51CE00A0331F65360462C022329597A81C3F92FC17938C9138F4111387979C28F0334F90119221374DAB045929B49E43A9646A243F4464DAF811AB00630C75961BCD4AF5D99115A3749191BA8FD41CE0B3C89A695B4BB85064FD3AF95C9B4AEE09AC7B0CC69ECA36A004B6CD662A6D32795053EF0A03ADA3B98BFE3B46A79723E3A45AB3C31950669AD77072062CC3B504DF1334FD6909EAC7915F1D5AD16639F5FB564416454259134D565882CB381CBA58B76880767B50AC1B85795D7268433B371230ED4C72F99AB1AD1E595A459CF0A2334AA1463ADE4BDC9249605381857BB98095B41132946CA2457DFAA9149582AA19927B63689E2929AA41027BEF4921970BAD4A55490D91ABE251DEF4552CA88034106A02CE4B058F8B59624B67E063BF178B015E4281EB114A2BC2454943A4B4647122C42CBEA4E94154FD3E4B791F6290B782994206853D67000A633F320A8A374CA5D4038F9CA4244DCB02E9A84E1F7C8A821132B32B9A840557B34780665301724BA2606681D945E34D7CF941B8963CAA1001A491B8B2E43570E9AB95C0A57C503F0AB960B4856D0251574710FE5CB474284FC1049AA2A7B03694A1C763E99DAC6AD0BA8038B138A64432E349116A031E8C792781751BA473CBDF55720005ABDAA13D50182F0E633776BB0675C40472BAD1F9672769183D0CCC810BC25A8573220569F6AC4BAC22A1354D8B36C0580D0E5299E629C506CC7655546FF27810C97B51BA056BBF86ED9CB7C0A537F72D0CF9AD2C231E29EBF553F613CBB15B3721A20077E505FD390CB19F6488A107DEE1CAC58AB7034BA690300219595B3695C1234E8B57E33C8D3A048454A616DF3C9B56A6FF2026AF997725FC95579043BAE9399B6790D637B4FA820B0B2D2CAB607BAF6A372734C31EE0026F3C076D14A8E3EE66AAD8BBBCCEB9DC70C7B6BB0BB76C200C231601CA0873EC8710F4B18D57290B033727C601EDB71C2B0F0C21D553E0E7A4F77716839C7C8448ABB9F66A54E8A4B08A79D9A392CA1270031388BAD56217E32AEF55411974906A245C00712B3CBB1170685193FE25ACD7AC13D32073F3879A5D78375F0052CF79175BAB46D22370597BD06789EDD0711CC4243507A02B4FAADBB62250CC997AE0327AEB00DEB529192A64B1096A86B19674D0B0AF05C4AAE178C2C9A6442E94ED0A56033A11EE42632C0B4AA51D42150790F41062B77253C25BA4DE559761F0A90068389728BC977F70CF7BCCFBD883DF13C79F5F2C34312CB1D5A55D78C1B242096A8C0593CFB2753460BD30ABA306C74173995748385D00B3670E61324D87DE8A14450DC493768777FF0CE6810937A711229561A5EF2BB69861074E00BD93266E4B86269E18EEA2CAACB60A1358636CD7A7CA6BB682130241784B101EA5BFD6C3A07158621614736F6996D5A4E14963A12D836E533A0C8912DB7E11685A4A53D8285F08750DFF66DA27C23B97542DEFB99E470ACD5E647C940CB57301B43CC3E68E64E28B06770695EF609265E06C60F22CB875849E62BAB88CC10ECF622C379CB54F13D8B2BAC902B9AB02BB330B45AC8B741C2647AC45B5BF48A6D3FE039986CC940C60A94E66CF644531016A5272450824314B5662A0A909ABFB46FD27BAED3ABA8259361596882B08B2AC7233930FC3786738ED2F81EE638C45C3B9CFD1951DB5BCC1445C2C1625D57D57B53904B6A1AB681580755E89FA79775A657CD62B4426304BC0C711E2807A2C9E852D4B4359EE6B53E4675F523C90782572DC7368FB400C328C70FC846B5E98A4330BBB627BDD784B4DAF0B1F645944942B4C2B6225C8B31E989545522BA6F10396034CB1CA745977844D570894C611A5608A757416D6DE59963C32798C493EFD2264C231910E9A30090CA7B5384F231B89BA68A238190EF1A2A43CB01703470A0F061A70738944BCD9B7004F24797AECB88B1091CFED0590B0415453C39B6EC45B66305FAEA6B55A4B7967505FE3862A267ADBFE05B9181A06501893391650EAAA4A6D16853349276F98E0F44CD726615C61C16713094D8AB093CAC71F2803E7D39109EF5009C9C2CDAF7B7A6B37A33A49881F4BB5D7245A14C5042280C76A84E63F49D0D619D46D723BAA747A3BA90A6FB637A9A1DC02268FD5C043D18CBA1528AC8E225C1F923D1CC84F2E78E25DC3CCE9353C9DAC2AD726A79F64940801DD5701EFBDCB80A98A25993CD7F80591320B63172718647B976A98A771686F0120A053B0C4474604305890FECAF23475DDCC11BC08A9C5F592ABB1A153DB1B883C0507EB68F78E0A14DEBBFEEC621E10A69B6DAAFAA916B539533E508007C4188CE05C862D101D4DB1DF3C4502B8C8AE1457488A36EAD2665BFACB321760281DB9CA72C7614363404A0A8EABC058A23A346875FA96BB18AC2CCF093B8A855673811CED47CBE1EE81D2CF07E43FC4872090853743108865F02C5612AA87166707EE90FFD5B8021F0AA016E5DBCD91F57B3562D3A2BCFA20A4C03010B8AA144E6482804B474FEC1F5E138BE632A3B9C82483DC6890A13B1E8EE6AF714EC5EFAC3B1976B29DADB605B14D3732B5DE118596516858117E2634C4EA0CC")
	ct, _ := hex.DecodeString("DFA6B9D72A63B420B89DDE50F7E0D56ECF876BFEF991FCE91C8D286FA6EABAC1730FD87741FE4AD717B282A21E235A55C3757D88D4CE62F414EB77EB9D357EE29D00087BF8110E5BBBC7C90419072EAE044BF7E183D43A94B2632AA14649619B70649521BC19370942EF70F36C34C8C23591EE0CA71A12D279E0F52D39ED0F913F8C262621FB242E680DEB307B0749C6B393A8EF66F8B04AAFA877B951AB93F598B4B2FAB04F88AC803984FF37E3FE74F3A616D5314EB3A826F874F8ECD3A5647D04942A57EFC09638470DC0A9DF40B317571D3984A78CF7D11751090722B3059E07591CC4A2ED9BA0DCE99BE9E5EE5DB8D698CDEB5814759BA977C90079CF2AFDE478069C513A60091A3A5D0111E22DE06CB145C14E22A214CB278C8152B0681BCAFF54D552B54A671C0DFEF775E7C54FEFC4853868C955971ABDAC2A76292CCCD4FD1C706B7D3614159673E9D7B29A2D3F63363129E7A21E803A460F2714E3E25922780AF38257CD1495ACD1E01980638DF58A153DAB07EFB5C7E78ADACF631956D69CCDA070459568BD9D11A2934BCF1643BC99468238910B1F742EBB3C03D39FD45CFB85BA309E29DD9B5CD560819EC729FCAC8B9D725E3E8ABEDE4B5298A8658EE3F781B0CE683CBB7335CD57EFE2204A8F197446D7314CDBF4C5D08CCC41F80857CC9571FBFB906060F7E17C8CEF0F274AFF83E393B15F2F9589A13AF4BC78E16CDDE62361D63B8DC903B70C01A43419CD2052150BD28719F61FF31F4A9BEC4DDBCEC1F8FB2EFBF37DFFFA4C7FECA8CE6D626BFDA16EE708D9206814A2EF988525615D4AC9BE608C4B03ABEE95B32A5DB74A96119A7E159AF99CD98E88EAF09F0D780E7C7E814B8E88B4F4E15FA54995D0ECBAD3EF046A4947F3E8B9E744241489B806FE9401E78BAFC8E882E9D6D0700F720C0024E7DA49061C5D18A62074040ABC0003200ED465231797930A2E2AA501F64862DDA13014A99F9D3270AA907EEB3FDBFF291600DF1F6B39684B11E396B70D86F90492E82B09BA25607B0C286FBC070182AC76FA7C859AAFEA87016AED22C3605A2789A1D439FD8D933342DAB745A3E550E7D77C01A6234BDA7D6BB19D495E6560FCE8396FC3C6E088ED60F5F2771416EA3BE5BE472B6404906C91E71D9A8672F390083655AB7D0EC6EDFE86789CE20BE2EA90CA5CC31416FB24CBAF94DA1468FE696BCDF5247CF117CBE9334076CA6896B2F6A016B1F7C73728807898D8B199756C2B0AA2457E1B4F7754C4576CE5645614EA15C1AE28B094EB217C7A7A41239576CBDA380EE68783432730AD5EBE7F51D6BE7FB02AB37BE0C96AAC9F3C790A18D159E6BABA71EC88C110FD84C336DF630F271CF79328B6C879DF7CDE0F70712220B1FBB9ACB48248D91F0E2B6E3BE40C2B221E626E7E330D9D83CC0668F7308591E14C7D72B841A6F05F3FDC139EECC1536765650B55A9CEC6BBF54CCEC5C3AC9A0E39F48F237BD4C660CB1A8D250BB6C8C010FEC34CC3D91599271C7531330F12A3E44FAFD905D2C6")
	kExp, _ := hex.DecodeString("BD7256B242F404869D662F80BF677A16C0C6FC1568CCA5B64582A01A6A142D71")

	k, err := kemDecaps(dk, ct)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(k, kExp) {
		t.Errorf("k: got %x, expected %x", k, kExp)
	}
}

var sinkElement fieldElement

func BenchmarkSampleNTT(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sinkElement ^= sampleNTT(bytes.Repeat([]byte("A"), 32), '4', '2')[0]
	}
}

var sink byte

func BenchmarkKeyGen(b *testing.B) {
	d := make([]byte, 32)
	rand.Read(d)
	z := make([]byte, 32)
	rand.Read(z)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ek, dk := kemKeyGen(d, z)
		sink ^= ek[0] ^ dk[0]
	}
}

func BenchmarkEncaps(b *testing.B) {
	d := make([]byte, 32)
	rand.Read(d)
	z := make([]byte, 32)
	rand.Read(z)
	m := make([]byte, 32)
	rand.Read(m)
	ek, _ := kemKeyGen(d, z)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c, K, err := kemEncaps(ek, m)
		if err != nil {
			b.Fatal(err)
		}
		sink ^= c[0] ^ K[0]
	}
}

func BenchmarkDecaps(b *testing.B) {
	d := make([]byte, 32)
	rand.Read(d)
	z := make([]byte, 32)
	rand.Read(z)
	m := make([]byte, 32)
	rand.Read(m)
	ek, dk := kemKeyGen(d, z)
	c, _, err := kemEncaps(ek, m)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		K, err := kemDecaps(dk, c)
		if err != nil {
			b.Fatal(err)
		}
		sink ^= K[0]
	}
}

func BenchmarkRoundTrip(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ek, dk := GenerateKey()
		c, Ke, err := Encapsulate(ek)
		if err != nil {
			b.Fatal(err)
		}
		Kd, err := Decapsulate(dk, c)
		if err != nil {
			b.Fatal(err)
		}
		if !bytes.Equal(Ke, Kd) {
			b.Fail()
		}
	}
}
