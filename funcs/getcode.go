package funcs

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/ethereum/go-ethereum/common"
	"github.com/nkbai/log"
	"golang.org/x/net/html"
)

func GetSourceCode(tokenAddr common.Address) (string, error) {
	doc, err := goquery.NewDocument(fmt.Sprintf("https://etherscan.io/address/%s", tokenAddr.String()))
	if err != nil {
		return "", err
	}
	var node *html.Node
	sel := doc.Find("#editor")
	if len(sel.Nodes) > 0 {
		node = sel.Get(0)
	} else {
		return "", errors.New("not found")
	}
	s := node.FirstChild.Data
	s = html.UnescapeString(s)
	return s, nil
}

/*
根据这样一个列表进行下载.
0x56ba2ee7890461f463f7be02aac3099f6d5811a8;BlockCAT (CAT)
0x0d262e5dc4a06a0f1c90ce79c7a60c09dfc884e4;J8T (J8T)
*/
func GetAllSourceCode(info string) {
	ss := strings.Split(info, "\n")
	for _, s := range ss {
		if len(s) < 43 {
			continue
		}
		s2 := strings.Split(s, ";")
		addr, name := s2[0], s2[1]
		log.Trace("addr=%s,name=%s\n", addr, name)
		code, err := GetSourceCode(common.HexToAddress(addr))
		if err != nil {
			log.Info("%s %s cannot get code %s", name, addr, err)
			continue
		}
		code += "\n//" + addr
		ioutil.WriteFile("./erc20/"+name+".sol", []byte(code), os.ModePerm)
	}
}
