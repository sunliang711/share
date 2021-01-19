package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"share/types"
	"share/utils"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

var (
	buildstamp string
	githash    string
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
		usage := `
Usage: %v CMD

CMD:
	push [-f filename]
	pull <key> [-o output]
	version|--v|--version
`
		fmt.Fprintf(os.Stderr, usage, os.Args[0])
		os.Exit(1)
	}

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

var (
	remoteURL = "http://localhost:8081"
	PUSH_TEXT = "/api/v1/push_text"
	PULL      = "/api/v1/pull"
)

func push() {
	inputFile := pflag.StringP("file", "f", "", "input file")
	pflag.Parse()
	if len(*inputFile) > 0 {

	} else {
		// read stdin
		var lines []string
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			logrus.Infof("line: %v", line)
			lines = append(lines, line)
		}
		oneLine := strings.Join(lines, "")
		data := types.Share{Content: oneLine}
		bs, err := json.Marshal(&data)
		if err != nil {
			panic(err)
		}

		req, err := http.NewRequest("POST", remoteURL+PUSH_TEXT, bytes.NewBuffer(bs))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		io.Copy(os.Stdout, resp.Body)

	}
}

func pull() {
	outputFile := pflag.StringP("output", "o", "", "output file")
	pflag.Parse()

	if pflag.NArg() < 2 {
		fmt.Fprintf(os.Stderr, "pull need key\n")
		os.Exit(1)
	}
	key := pflag.Arg(1)
	resp, err := http.Get(fmt.Sprintf("%s%s/%s", remoteURL, PULL, key))
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
