package logging

import (
	seelog "github.com/cihub/seelog"
	// "errors"
	// "io"
	"fmt"
)

var Logger seelog.LoggerInterface

func loadAppConf() {
	// appConf := `<seelog />`
	appConf := `
<seelog minlevel="warn">
    <outputs formatid="common">
        <rollingfile type="size" filename="/Users/winniewang/source/go-workspace/logs/roll.log" maxsize="100000" maxrolls="5"/>
        <filter levels="critical">
            <file path="/Users/winniewang/source/go-workspace/logs/critical.log" formatid="critical"/>
            <smtp formatid="criticalemail" senderaddress="zhou.zhefu@gmail.com" sendername="Zhou Zhefu" hostname="smtp.gmail.com" hostport="587" username="zhou.zhefu" password="5642326">
                <recipient address="zhou.zhefu@gmail.com"/>
            </smtp>
        </filter>
    </outputs>
    <formats>
        <format id="common" format="%Date/%Time [%LEV] %Msg%n" />
        <format id="critical" format="%File %FullPath %Func %Msg%n" />
        <format id="criticalemail" format="Critical error on our server!\n    %Time %Date %RelFile %Func %Msg \nSent by Seelog"/>
    </formats>
</seelog>
	`

	logger, err := seelog.LoggerFromConfigAsBytes([]byte(appConf))
	if err != nil {
		fmt.Println(err)
		return
	}

	UseLogger(logger)
}

func init() {
	loadAppConf()
}

func DisableLogger() {
	Logger = seelog.Disabled
}

func UseLogger(logger seelog.LoggerInterface) {
	Logger = logger
}