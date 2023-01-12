package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"net"
	"net/http"
	"io/ioutil"
	"github.com/cheggaaa/pb"


)

const applicationName string = "Marceline"
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
	dumpPath        string
	outputName      string
	ipAddr          string
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
	flag.String("o", "", "Path to store dump files")
	flag.String("f", "", "Name for output file")
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
        dumpPath = viper.GetString("o")
	outputName = viper.GetString("f")

}




func elastic() {

	fmt.Print("Enter an IP address: ")
	fmt.Scan(&ipAddr)

	ip := net.ParseIP(ipAddr)
	if ip == nil {
		fmt.Println("Invalid IP address")
		return
	}

	ipv4 := ip.To4()
	if ipv4 == nil {
		fmt.Println("Invalid IPv4 address")
		return
	}

	ipRange := &net.IPNet{
		IP:   ipv4,
		Mask: net.CIDRMask(24, 32),
	}
	fmt.Println("IP range:", ipRange)

	fmt.Println("Starting Scan for Hunt Elastic.Please Wait until Scan Complete")
	ips := ipRange.IP.Mask(ipRange.Mask)
	total := 0
	for ; ipRange.Contains(ips); inc(ips) {
		total++
	}
	bar := pb.StartNew(total)
	var elasticIPs []string
	for ip := ipRange.IP.Mask(ipRange.Mask); ipRange.Contains(ip); inc(ip) {
		// Check if port 9200 is open
		_, err := net.DialTimeout("tcp", fmt.Sprintf("%s:9200", ip), 1*time.Second)
		if err == nil {
			elasticIPs = append(elasticIPs, ip.String())

		}
		bar.Increment()
	}
	bar.Finish()
	fmt.Println("Scan finished. IPs that use Elastic are:")
	for _, ip := range elasticIPs {
		fmt.Println(ip)
	}

}

// increment the last octet of an IP address
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
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
		file, err := os.Create(dumpPath + "/" + outputName)
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
			file, err := os.Create(dumpPath + "/" + outputName)
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
	        elastic()
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
	  -o                           Path to store dump files
	  -f                           Name for output file
	  
	  `
	fmt.Println(applicationName + " " + applicationVersion)
	fmt.Println(message)
}

func banner() {
	banner := `


                                                                                                                  ...                                                                                                                                    
                                                                                                                     ...                                                                                                                                  
                                                                                                                       ...                                                                                                                                
                                                                                                             ..         ....                                                                                                                              
                                                                                                            ....          .....                                                                                                                           
                                                                                                            .....           .....                                                                                                                         
                                                                                                            ......             ....                                                                                                                       
                                                                                                           ..:d:...              ....                                                                                                                     
                                                                                                          ...oKOl'...              .....                                                                                                                  
                                                                                                          ..;OKKKd,...               .....                                                                                                                
                                                                                                         ..'oKKKKKOc'...               .....                                                                                                              
                                                                                                         ..:OKKKKKK0x;....               .....                                                                                                            
                                                                                                        ..'dKKKKKKKKK0o,....               .....                                                                                                          
                                                                                                        ..c0KKKKKKKKKKKOl'....               .....                                                                                                        
                                                                                                       ..,kKKKKKKKKKKKKKKkc'....               .....                                                                                                      
                                                                                                      ..'oKKKKKKKKKKKKKKKKKkc;'...               .....                                                                                                    
                                                                                                     ...c0KKKKKKKKKKKKKKKKKKK0kc'...               .....                                                                                                  
                                                                                                     ..,kKKKKKKKKKKKKKKKKKKKKKKKkl,...               .....                                                                                                
                                                                                                    ..'dKKKKKKKKKKKKKKKKKKKKKKKKKKOl,...               .....                                                                                              
                                                                                                   ...c0KKKKKKKKKKKKKKKKKKKKKKKKKKKKOo;...               ....                                                                                             
                                                                                                   ..;OKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKK0d;...              .....                                                                                           
                                                                                                  ..,xKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKK0d:...              .....                                                                                         
                                                                                                 ...oKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKK0x:'..             .....                                                                                        
                                                                                                ...c0KKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKK0x:'..             ....                                                                                       
                                                                                                ..;kKKKKKKKK000KKKKKKKKKKKKKKKKKKKKKKKKKK000KK0x:'..            ....                                                                                      
                                                                                               ..'dKKKKKKKOc'.'o0KKKKKKKKKKKKKKKKKKKKKKOc'.'lOKK0x:...           ....                                                                                     
                                                                                              ...l0KKKKKKKo.   .xXKKKKKKKKKKKKKKKKKKKKKo.   .dXKKK0x;..           ....                                                                                    
                                                                                              ..:OKKKKKKKKOc'.,o0KKKKKKKKKKKKKKKKKKKKKKOc'.'lOKKKKKK0o,..          ....                                                                  .                
                                                                                             ..,xKKKKKKKKKKK00KK0OOO0KKKKKKKKKKKKKKKKKK00000KKKKKKKKKKOl'.          ....                                                                ..                
                  .                                                                         ...oKKKKKKKKKKKKKKOc'...;loodkO0000OOkxooc,'';o0KKKKKKKKKKKKx;..         ....                                                               ..                
                  ..                                                                        ..:OKKKKKKKKKKKKKKl......o: ..''''''...,d, ....xXKKKKKKKKKKKKOl..         ...                                                               ...               
                  ..                                                                       ..,xKKKKKKKKKKKKKKKx'.......''...........'.... ,OKKKKKKKKKKKKKKKd,.         ...                                                              ...               
                  ..                                                                      ...lKKKKKKKKKKKKKKKKKk:..,,,;:::;'.............:kKKKKKKKKKKKKKKKKKx,.        ...                                                              ...               
                  ..                                                                      ..;OKKKKKKKKKKKKKKKKKKKOo:;,''',,;,,........;lkKKKKKKKKKKKKKKKKKKKKd'.        ...                                                             ...               
                  ...                                                                    ..'dKKKKKKKKKKKKKKKKKKKKKKK0kdolc:::::;:cloxOKKKKKKKKKKKKKKKKKKKKKKKo...       ...                                                             ...               
                  ...                                                                    ..:OKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKk.  ..      ...                                                            ....               
                  ...                                                                   ...c0KKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKO,   ...      ..                                                            ....               
                  ....                                                                  ....l0KKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKO;     ..      ..                                                            ....               
                   ...                                                                 ...  .c0KKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKk,      ...     ..                                                           ....                
                   ....                                                                ...    ;OKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKd.       ...     .                                                            ....                
                   .....                                                               ...     .d0KKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKO:.        ...                                                                 .....                
                    ....                                                              ...        ;xKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKOl.           ..                                                                .....                 
                     ....                                                             ...         .:kKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKOl'             ..                                                                .....                 
                     .....                                                            ...           .;x0KKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKkc.              ...                                                               .....                  
                      .....                                                           ...              ,d0KKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKK0x:.                ...                                                              .....                   
                       .....                                                          ...                'lOKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKK0d,                   ...                                                             ......                   
                        .....                                                         ...                  .cOKKKKKKKKKKKKKKKKKKKKKKKKKKK0o'                     ...                                                             .....                    
                         .....                                                        ...                    .l0KKKKKKKKKKKKKKKKKKKKKKKKd'                      ...                                                             .....                     
                          .....                                                       ...                      ;OKKKKKKKKKKKKKKKKKKKKK0l.                       ...                                                            .....                      
                           .....                                                      ...                       ;0KKKKKKKKKKKKKKKKKKKKo.                       ...                                                            .....                       
                            .....                                                     ...                       .dKKKKKKKKKKKKKKKK0kOO;                        ...                                                           .....                        
                              ....                                                     ..                        lKKKKKKKKKKKKKKKKk;:x,                       ...                                                           ....                          
                               ....                                                    ..                       .dXKKKKKKKKKKKKKKKKOO0:                       ...                                                          ....                           
                                ....                                                   ...                      ,OKKKKKKKKKKKKKKKKKxlxd.                     ...                                                          ....                            
                                  ...                                                   ..                     .xKKKKKKKKKKKKKKKKKKxcd0l.                   ....                                                          ...                             
                                   ...                                                  ...                   .dKKKKKKKKKKKKKKKKKKKKKKK0l.                  ...                                                          ...                              
                                    ...                                                  ..                  'xKKKKKKKKKKKKKKKKKKKKKKKKK0o.                ...                                                          ...                               
                                      ..                                                 ...               .:OKKKKKKKKKKKKKKKKKKKKKKKKKKKKk:.              ...                                                         ...                                
                                       ..                                                ...              .dKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKK0,             ...                                                          ..                                 
                                        ..                                                ...           . ,OXKKKKKKKKKKKKKKKKKKKKKKKKKKKKKXx..'.         ....                                                         ..                                  
                                         .                                                ...         .';..oKKKKKKKKKKKKKKKKKKKKKKKKKKKKKK0:.':;,.       ...                                                          .                                   
                                                                                           ...      .,;;;,.,OKKKKKKKKKKKKKKKKKKKKKKKKKKKKKd..;;;;;,'.   ...                                                                                               
                                                                                           ...   .';;;;;;;'.c0KKKKKKKKKKKKKKKKKKKKKKKKKKKk'.,:;;;;;;;,.....                                                                                               
                                                                                            ....,;;;;;;;;;;..oKKKKKKKKKKKKKKKKKKKKKKKKKKO;.,;;;;;;;;;;,...                                                                                                
                                                                                             ..';;;;;;;;;;;;..oKKKKKKKKKKKKKKKKKKKKKKKKO;.';;;;;;;;;;;...                                                                                                 
                                                                                             ...,;;;;;;;;;;;;..l0KKKKKKKKKKKKKKKKKKKKKk,.';;;;;;;;;;;'...                                                                                                 
                                                                                              ..';;;;;;;;;;;;;..;kKKKKKKKKKKKKKKKKKK0o'.,;;;;;;;;;;;,...                                                                                                  
                                                                                              ...,;;;;;;;;;;;;;,..ckKKKKKKKKKKKKKK0d;.';;;;;;;;;;;;;'..                                                                                                   
                                                                                               ...;;;;;;;;;;;;;;;'..;ok0KKKKKKK0xl,.';;;;;;;;;;;;;;,...                                                                      
                                                                                               ...,;;;;;;;;;;;;;;;;,...;ldxxdoc,..';;;;;;;;;;;;;;;;'..                                                                        


`
	fmt.Println(banner)
}

