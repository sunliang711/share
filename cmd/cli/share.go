package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"share/types"
	"share/utils"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	buildstamp string
	githash    string
)

var (
	RemoteURL = "http://localhost:8081"
	PUSH_TEXT = "/api/v1/push_text"
	PULL      = "/api/v1/pull"
)

// 两个功能模块
// push 参数 -f: 上传文件 (默认从标准输入) -v (verbose) 返回错误写到stderr,否则stdout
// pull 参数： key
// 如果是文件： 默认以原来的文件名保存,可以使用-o参数覆盖
// 如果是文本:  默认输出到标准输出，可以指定-o参数指定文件

// 配置
// url
func main() {
	if len(os.Args) < 2 {
		usage := `Usage: %v CMD

CMD:
	push 		[-f filename] [-v|--verbose]
	pull 		<key> [-o output] [-v|--verbose]

	version|-V|--version
`
		fmt.Fprintf(os.Stderr, usage, os.Args[0])
		os.Exit(1)
	}
	// logrus.SetLevel(logrus.DebugLevel)

	var configFile string
	usr, err := user.Current()
	if err == nil {
		configFile = usr.HomeDir + "/.share.toml"
		logrus.Debugf("config file: %v", configFile)
	}

	if _, err := os.Stat(configFile); os.IsNotExist(err) {

	} else {
		logrus.Debugf("read config file: %v", configFile)
		// read config file
		viper.SetConfigFile(configFile)
		err = viper.ReadInConfig()
		if err != nil {
			utils.QuitMsg(fmt.Sprintf("Read config file error: %v", err))
		}
	}

	if viper.GetString("remote_url") != "" {
		RemoteURL = viper.GetString("remote_url")
	}

	logrus.Debugf("remote url: %v", RemoteURL)
	switch os.Args[1] {
	case "push":
		push()
	case "pull":
		pull()
	case "version", "--version", "-V":
		version()
	default:
	}

}

func push() {
	inputFile := pflag.StringP("file", "f", "", "input file")
	verbose := pflag.BoolP("verbose", "v", false, "verbose")
	pflag.Parse()
	if len(*inputFile) > 0 {
		fmt.Fprintf(os.Stderr, "%s", "TODO")
		os.Exit(1)
	} else {
		// read stdin
		var result []byte
		buf := make([]byte, 1024)
		if *verbose {
			logrus.Infof("Read from stdin...\n")
		}
		for {
			n, err := os.Stdin.Read(buf)
			if n > 0 {
				result = append(result, buf[:n]...)
			}
			if err == io.EOF {
				if *verbose {
					logrus.Infof("End\n")
				}
				break
			} else if err != nil {
				utils.QuitMsg(fmt.Sprintf("read stdin error: %v", err))
			}
		}

		data := types.Share{Content: string(result)}
		bs, err := json.Marshal(&data)
		if err != nil {
			utils.QuitMsg(fmt.Sprintf("Encode request body error: %v", err))
		}

		if *verbose {
			logrus.Infof("remote url: %v", RemoteURL)
		}
		req, err := http.NewRequest("POST", RemoteURL+PUSH_TEXT, bytes.NewBuffer(bs))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			utils.QuitMsg(fmt.Sprintf("POST error: %v", err))
		}
		defer resp.Body.Close()
		io.Copy(os.Stdout, resp.Body)
	}
}

func pull() {
	outputFile := pflag.StringP("output", "o", "", "output file")
	verbose := pflag.BoolP("verbose", "v", false, "verbose")
	pflag.Parse()

	if pflag.NArg() < 2 {
		fmt.Fprintf(os.Stderr, "pull need key\n")
		os.Exit(1)
	}
	key := pflag.Arg(1)
	if *verbose {
		logrus.Infof("key: %v", key)
		logrus.Infof("remote url: %v", RemoteURL)
	}
	resp, err := http.Get(fmt.Sprintf("%s%s/%s", RemoteURL, PULL, key))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	// io.Copy(os.Stdout, resp.Body)

	var r types.Share
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		utils.QuitMsg(fmt.Sprintf("Decode response error: %v\n", err))
	}
	if r.Code != types.OK {
		utils.QuitMsg(fmt.Sprintf("Error: %s\n", r.Msg))
	}

	switch r.Type {
	case types.TextType:
		if len(*outputFile) > 0 {
			if f, err := os.OpenFile(*outputFile, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0600); err != nil {
				utils.QuitMsg(fmt.Sprintf("Open file error: %v", err))
			} else {
				fmt.Fprintf(f, "%s", r.Content)
			}
		} else {
			fmt.Fprintf(os.Stdout, "%s", r.Content)
		}
	case types.FileType:
		utils.QuitMsg("FileType TODO\n")
	default:
		utils.QuitMsg("Invalid type\n")
	}

}

func version() {
	fmt.Fprintf(os.Stderr, "Build time: %s\ngit rev: %s\n", buildstamp, githash)
	os.Exit(0)
}
