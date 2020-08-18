package services

import (
	"fmt"
	"lxm-ency/utils"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	encyDomain        string        = "https://baike.baidu.com"
	encyURL           string        = "https://baike.baidu.com/item/%s"
	crawlerResultSpan time.Duration = time.Hour

	// SynpPath - path for article synopsis
	SynpPath string = "data/synps/"
)

// Encyclopedia crawler
type encyCrawler struct {
	keyword  string
	title    string
	subtitle string
	timgURL  string
	synp     string
	link     string
	cid      string
}

// Constructor
func newEncyCrawler(keyword string) (*encyCrawler, error) {
	var ec = new(encyCrawler)
	ec.keyword = keyword
	ec.title = fmt.Sprintf("Knowledge about %s", keyword)
	ec.subtitle = "Encyclopedia knowledge helps you know the world better."

	// crawl data
	err := ec.crawl()
	if err != nil {
		return ec, err
	}

	// initialize thumbnail
	err = ec.initTimg()
	if err != nil {
		return ec, err
	}

	// save card
	cid, err := utils.SaveCard(ec.title, ec.subtitle, ec.timgURL, ec.link)
	if err != nil {
		return ec, err
	}
	ec.cid = cid

	return ec, nil
}

func (ec *encyCrawler) initTimg() error {
	ec.timgURL = `https://ss0.bdstatic.com/70cFuHSh_Q1YnxGkpoWK1HF6hhy/it/u=1048407325,3557533849&fm=15&gp=0.jpg` //兜底图片
	return nil
}

func (ec *encyCrawler) crawl() error {
	// do not crawl if local data are available
	exist, _ := utils.FileExist(ec.getSynpPath())
	if exist {
		result, err := readResult(ec.getSynpPath())
		if err != nil {
			return err
		}
		ec.synp = result.Synp
		ec.link = result.Link
	} else {
		response, err := download(ec.getURL())
		if err != nil {
			return err
		}
		defer response.Body.Close()
		synp, link, err := ec.analyze(response)
		if err != nil {
			return err
		}
		writeResult(ec.getSynpPath(), synp, link)
		ec.synp = synp
		ec.link = link
		go utils.RemoveWithDelay(ec.getSynpPath(), crawlerResultSpan)
	}
	return nil
}

func (ec *encyCrawler) analyze(resp *http.Response) (string, string, error) {
	link := ec.getURL()
	body, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", "", err
	}

	// handle homonyms
	polysemy := false
	body.Find(".lemmaWgt-subLemmaListTitle").Each(func(index int, item *goquery.Selection) {
		polysemy = true
	})
	if polysemy {
		fmt.Println("Found multiple articles..")
		var nextURL string
		body.Find(".list-dot").Each(func(index int, item *goquery.Selection) {
			if nextURL != "" {
				return
			}
			item.Find("a").Each(func(index int, item *goquery.Selection) {
				if nextURL != "" {
					return
				}
				link, _ := item.Attr("href")
				nextURL = encyDomain + link
			})
		})
		resp, err = download(nextURL)
		if err != nil {
			return "", "", err
		}
		body, err = goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			return "", "", err
		}
		link = nextURL
	}

	// select synopsis using goquery
	var paras []string
	body.Find(".lemma-summary").Each(func(index int, item *goquery.Selection) {
		item.Find(".para").Each(func(index int, item *goquery.Selection) {
			formatText := format(item.Text())
			if formatText != "" {
				paras = append(paras, formatText)
			}
		})
	})
	if len(paras) == 0 {
		return "", "", encyHandleErr{"Lookup failed"}
	}
	return concat(paras), link, nil
}

func (ec *encyCrawler) getURL() string {
	return fmt.Sprintf(encyURL, url.QueryEscape(ec.keyword))
}

func (ec *encyCrawler) getSynpPath() string {
	return path.Join(utils.Cwd, SynpPath, "temp_"+ec.keyword+".json")
}

func download(link string) (*http.Response, error) {
	fmt.Print("Trying to access link: ", link)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", link, nil)

	// Simulate browser headers
	headers := map[string]string{
		"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.106 Safari/537.36",
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return resp, err
	}
	fmt.Printf(" (%s %s)\n", resp.Proto, resp.Status)

	return resp, nil
}

func format(text string) string {
	re := regexp.MustCompile(`\n\[[\s\S]+?\] \n`)
	text = re.ReplaceAllString(text, "")
	re = regexp.MustCompile(`\s+`)
	text = re.ReplaceAllString(text, "")
	return text
}

func concat(paras []string) string {
	fmt.Println(paras)
	var synp string
	for i := range paras {
		if i != 0 {
			synp += "\n\n"
		}
		synp += paras[i]
	}
	return synp
}
