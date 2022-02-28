package main 

/* Workaround for build errs: go mod init main; go mod tidy 
 * https://golangbyexample.com/go-mod-tidy/
 * https://stackoverflow.com/questions/53837919/should-go-sum-file-be-checked-in-to-the-git-repository 
 * Cannot refer to unexported name loglevel.setPriorityString: 
 * https://www.sneppets.com/golang/cannot-refer-unexported-name-undefined-error-go/ 
 * https://stackoverflow.com/questions/51642410/chaincode-not-building-go-program-error-cannot-refer-to-unexported-name
 * fix: changed outdated method name in tutorial https://github.com/llimllib/loglevel/blob/master/log.go
 * Understanding go.mod and go.sum https://faun.pub/understanding-go-mod-and-go-sum-5fd7ec9bcc34
 */
import (
	"fmt"
	"net/http"
	"os"
	"golang.org/x/net/html" // helps to better parse html 
	"golang.org/x/net/html/atom" // to break down html down into atoms
	"strings"
	log "github.com/llimllib/loglevel" // log level package to help build a log for err handling
)

/* Program pull all links in html body and display to console and output to a console 
Can modify if you like to output to a file */

var maxDepth = 2

// define custom type in go 
type Link struct {
	url string
	text string
	depth int
}


type HttpError struct {
	original string
}


/* Purpose: read links
 * Input: pointer to resp <http.Response>; depth <int>
 * Output: []Link https://appliedgo.net/slices/
 */
 
func LinkReader(resp *http.Response, depth int) []Link {
	page := html.NewTokenizer(resp.Body) // allows us to parse html and create tokens
	links := []Link{}

	var start *html.Token 
	var text string 

	// assign each page to a token and then sift through token to pull out the links
	for {
		_= page.Next() // use anonymous char to make pages move forward 
		token := page.Token()

		// err handling
		if token.Type == html.ErrorToken {
			break
		}

		if start != nil && token.Type == html.TextToken {
			// print out the curr link 
			text = fmt.Sprintf("%S%S", text, token.Data)
		}

		// Atom.A = 0x1 integer code that maps to a specific HTML string
		// in this case <a> tag (which is what we are crawling) https://pkg.go.dev/golang.org/x/net/html/atom
		if token.DataAtom == atom.A {
			// switch depending on what type of token we're dealing with 
			switch token.Type  {
			case html.StartTagToken: 
				if len(token.Attr) > 0 {
					start = &token
				}
			case html.EndTagToken: 
				if start == nil {
					// end of token found but no start so log err 
					log.Warnf("Link End found without Start: %s", text)
					continue // dont break from switch
				}
				// token + string + int
				link := NewLink(*start, text, depth)
				if link.Valid() {
					links = append(links, link)
					log.Debugf("Link Found %v", link)

				}

				// reset start and text 
				start = nil 
				text = ""

			}
		}
	}

	log.Debug(links)
	return links
}

/* Purpose: Create new Link 
 * Input: tag<html.Token>; text<str>; depth<int>
 * Output: Link 
 */
func NewLink(tag html.Token, text string, depth int ) Link {

	// fixed: cannot use depth (type int) as type string in field value
	link := Link{text: strings.TrimSpace(text), depth: depth}

	// Scan tag.Attr for hrefs and assign to link.url 
	for i := range tag.Attr {
		if tag.Attr[i].Key == "href" {
			link.url = strings.TrimSpace(tag.Attr[i].Val)

		}
	}
	return link 
}


/* Purpose: format strings
 * Input: self<Link>
 * Output: string<String>
 */
func (self Link) String() string {
	spacer := strings.Repeat("\t", self.depth)
	// modify way the strings look, rather than regex
	return fmt.Sprintf("%s%s (%d) - %s", spacer, self.text, self.depth, self.url) 
}

/* Purpose: check if link is valid
 * Input: self<Link>
 * Output: bool 
 */
func (self Link) Valid() bool {
	if self.depth >= maxDepth {
		return false 
	}

	if len(self.text) == 0 {
		return false 
	}

	if len(self.url) == 0 || strings.Contains(strings.ToLower(self.url), "javascript") {
		return false
	}
	return true 
}

/* Purpose: deal with error by reseting to original self
 * Input: self<HttpError>
 * Output: string 
 */
func (self HttpError) Error() string {
	return self.original
}

/* Purpose: allows us to download http multiple times
 * Input: url<str>; depth<int>
 * Output: nil  
 */
func recurDownloader(url string, depth int) {
	// download from url
	page, err := downloader(url)
	if err != nil {
		log.Error(err)
		return
	}
	links := LinkReader(page, depth)

	// blank identifier: avoids having to declare all the variables for the returns values https://go.dev/doc/effective_go#blank
	for _, link := range links {
		fmt.Println(link)
		if depth + 1 < maxDepth {
			// recursive call, incrementing depth by 1 
			recurDownloader(link.url, depth + 1)
		}
	}
}

/* Purpose: allows us to download http 
 * Input: url<str>; depth<int>
 * Output: nil  
 */
 func downloader(url string) (resp *http.Response, err error) {
 	log.Debugf("Downloading %s", url)
 	resp, err = http.Get(url)
 	if err != nil {
 		log.Debugf("Error: %s", err)
 		return 
 	}

 	if resp.StatusCode > 299 {
 		err = HttpError { fmt.Sprintf( "Error (%d): %s", resp.StatusCode, url)}
 		log.Debug(err)
 		return
 	}
 	return
 }


func main() {
	// set output priority to debug level of info 
	// outdated in tutorial 
	// log.setPriorityString("info")
	log.SetPriorityString("info")
	
	// set prefix to str crawler 
	// outdated in tutorial 
	// log.setPrefix("crawler")
	log.SetPrefix("crawler")

	// os.Args allows us to call piece of http from console 
	log.Debug(os.Args)

	if len(os.Args) < 2 {
		log.Fatalln("Missing Url arg")
	}

	// download http multiple times 
	recurDownloader(os.Args[1], 0) 

}