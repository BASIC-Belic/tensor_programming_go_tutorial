/** Tutorial Playlist: https://www.youtube.com/watch?v=uCR_A-Bphl0&list=PLJbE2Yu2zumCe9cO3SIyragJ8pLmVv0z9
 * 
 * Other articles I consulted: 
 * https://github.com/golang/sublime-config/blob/master/docs/user.md
 * https://www.wolfe.id.au/2015/03/05/using-sublime-text-for-go-development/
 * https://github.com/golang/sublime-build/blob/master/docs/configuration.md
 * https://www.alexedwards.net/blog/streamline-your-sublime-text-and-go-workflow
 * https://jisaacks.github.io/GitGutter/install/
 * https://stackoverflow.com/questions/64448560/golang-package-is-not-in-goroot-usr-local-go-src-packagename
**/ 
package main 
// // core IO lib in go, allows us to print, take in input 
// import "fmt" 

import (
	"fmt" // core IO lib in go, allows us to print, take in input
	"flag" // allows you to implement command line flags 
	"net/http"
	"log"
)

func main() {
	fmt.Println("Enter your name: ")
	// var then variable name then variable type
	var input string 
	fmt.Scanln(&input) /** pass reference of input so that we can collect 
	the input f2rom the user until they press enter key **/
	fmt.Println("Hello, %s!", input)

	/** Build a server to allow us to serve static html files **/

	// creating strng that gets passed to command line
	port := flag.String("p", "8000", "port") 
	dir := flag.String("d", ".", "dir") // . bc index.html is in current dir 
	flag.Parse() // will parse those two flags 

	/** Creating http request **/
	// handle takes in pattern string, root pattern in this case & Handler
	// pass server directory which is http dir and pass that the pointer for dir 
	http.Handle("/", http.FileServer(http.Dir(*dir))) 
	//print log function
	log.Printf("Serving %s on Http port: %s\n", *dir, *port)
	// fatal similar to print but followed by call to os.exit
	/** http listen and serve listens on tcp network address 
	 * and then serves the handlerto handle requests on incoming connections 
	 * need define where we want it to listen with nil as second handler**/
	log.Fatal(http.ListenAndServe(":" + *port, nil))

	/** RUN and open http://localhost:8000/ in browser **/
	
}