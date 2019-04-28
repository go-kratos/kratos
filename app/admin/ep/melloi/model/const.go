package model

// const const
const (
	PROTOCOL_HTTP     = 0
	PROTOCOL_GRPC     = 1
	PROTOCOL_SCENE    = 2
	HTTP_SCRIPT_TYPE  = 1
	SCENE_SCRIPT_TYPE = 2

	//HeaderStart header start
	HeaderStart = "<elementProp name=\"\" elementType=\"Header\"><stringProp name=\"Header.name\">"
	//HeaderMid  header mid
	HeaderMid = "</stringProp><stringProp name=\"Header.value\">"
	//HeaderEnd header end
	HeaderEnd = "</stringProp></elementProp>"
	//ElementPropName element propname
	ElementPropName = "<elementProp name=\""

	//HTTPArgument http arg
	HTTPArgument = "\" elementType=\"HTTPArgument\">"

	//HTTPArgumentEncode http arg code
	HTTPArgumentEncode = "<boolProp name=\"HTTPArgument.always_encode\">false</boolProp>"

	//ArgumentStart arg start
	ArgumentStart = "<stringProp name=\"Argument.value\">"
	//ArgumentMid arg mid
	ArgumentMid = "</stringProp><stringProp name=\"Argument.metadata\">=</stringProp><boolProp name=\"HTTPArgument.use_equals\"" +
		">true</boolProp><stringProp name=\"Argument.name\">"
	//ArgumentEnd arg end
	ArgumentEnd = "</stringProp></elementProp>"

	//AsyncInfo async info
	AsyncInfo = "<boolProp name=\"asyncCall\">true</boolProp>"

	//MultipartName Multipart Name
	MultipartName = "<elementProp name=\"HTTPsampler.Files\" elementType=\"HTTPFileArgs\"><collectionProp name=\"HTTPFileArgs.files\"><elementProp name=\""

	//MultipartFilePath MultipartFile Path
	MultipartFilePath = "\" elementType=\"HTTPFileArg\"><stringProp name=\"File.path\">"

	//MultipartFilePathd MultipartFile Pathd
	MultipartFilePathd = "</stringProp><stringProp name=\"File.paramname\">"

	//MultipartMime type
	MultipartMimetype = "</stringProp><stringProp name=\"File.mimetype\">"

	//Multipart End
	MultipartEnd = "</stringProp></elementProp></collectionProp></elementProp>"

	//Assertion start
	AssertionStart = "<stringProp name=\"927604211\">"

	//Assertion End
	AssertionEnd = "</stringProp>"

	//ConstTimer const timer
	ConstTimer = "<ConstantTimer guiclass=\"ConstantTimerGui\" testclass=\"ConstantTimer\" testname=\"固定定时器\" enabled=\"true\">" +
		"<stringProp name=\"ConstantTimer.delay\">1000</stringProp>" +
		"</ConstantTimer>" +
		"<hashTree/>"

	//Randomtimer Random Timer
	RandomTimer = "<GaussianRandomTimer guiclass=\"GaussianRandomTimerGui\" testclass=\"GaussianRandomTimer\" testname=\"高斯随机定时器\" enabled=\"true\">" +
		"<stringProp name=\"ConstantTimer.delay\">1000</stringProp>" +
		"<stringProp name=\"RandomTimer.range\">500</stringProp>" +
		"</GaussianRandomTimer>" +
		"<hashTree/>"
	//
)
