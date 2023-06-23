package getreleasegithub

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/google/go-github/v53/github"
)

// // code based on https://gist.github.com/metal3d/002e4f0d8545f83c2ace

// // var assets, download bool

// func prepareRequest(url string) *http.Request {

// 	req, _ := http.NewRequest("GET", url, nil)
// 	req.Header.Add("User-Agent", "metal3d-go-client")
// 	return req
// }

// // Download resource from given url, write 1 in chan when finished
// func downloadResource(repo string, id float64, c chan int) {
// 	defer func() { c <- 1 }()
// 	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/assets/%.0f", repo, id)
// 	fmt.Printf("Start: %s\n", url)
// 	req := prepareRequest(url)

// 	req.Header.Add("Accept", "application/octet-stream")

// 	client := http.Client{}
// 	resp, _ := client.Do(req)

// 	disp := resp.Header.Get("Content-disposition")
// 	re := regexp.MustCompile(`filename=(.+)`)
// 	matches := re.FindAllStringSubmatch(disp, -1)

// 	if len(matches) == 0 || len(matches[0]) == 0 {
// 		log.Println("WTF: ", matches)
// 		log.Println(resp.Header)
// 		log.Println(req)
// 		return
// 	}

// 	disp = matches[0][1]

// 	f, err := os.OpenFile(disp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	b := make([]byte, 4096)
// 	var i int

// 	for err == nil {
// 		i, err = resp.Body.Read(b)
// 		f.Write(b[:i])
// 	}
// 	fmt.Printf("Finished: %s -> %s\n", url, disp)
// 	f.Close()
// }

func downloadUriToPath(filepath string, url string) (err error) {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func ensureOutputPath(outputFile string) error {
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		return nil
	}
	return fmt.Errorf("file %s already exists", outputFile)
}

func DownloadPublic(owner string, repo string, tag string, outputFile string) error {
	client := github.NewClient(nil)
	err := ensureOutputPath(outputFile)
	if err != nil {
		return err
	}
	uri, _, err := client.Repositories.GetArchiveLink(context.Background(), owner, repo, github.Tarball, &github.RepositoryContentGetOptions{Ref: tag}, true)
	if err != nil {
		return err
	}

	err = downloadUriToPath(outputFile, uri.String())
	if err != nil {
		return err
	}
	return nil
}

// 	// command to call
// 	command := "releases/latest"
// 	if len(tag) > 0 {
// 		command = fmt.Sprintf("releases/tags/%s", tag)
// 	}

// 	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", repo, command)

// 	// create a request with basic-auth
// 	req := prepareRequest(url)

// 	// Add required headers
// 	req.Header.Add("Accept", "application/vnd.github.v3.text-match+json")
// 	req.Header.Add("Accept", "application/vnd.github.moondragon+json")

// 	// call github
// 	client := http.Client{}
// 	resp, err := client.Do(req)

// 	if err != nil {
// 		return fmt.Errorf("error while making request: %v", err)
// 	}

// 	// status in <200 or >299
// 	if resp.StatusCode < 200 || resp.StatusCode > 299 {
// 		return fmt.Errorf("error: %d %s", resp.StatusCode, resp.Status)
// 	}

// 	bodyText, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return fmt.Errorf("error reading response: %v", err)
// 	}

// 	// prepare result
// 	result := make(map[string]interface{})
// 	json.Unmarshal(bodyText, &result)

// 	// print download url
// 	results := make([]interface{}, 0)
// 	// if !assets {
// 	// 	if download {
// 	// results = append(results, result["id"])
// 	// 	} else {
// 	// results = append(results, result["zipball_url"])
// 	// 	}
// 	// } else {
// 	for _, asset := range result["assets"].([]interface{}) {
// 		// if download {
// 		results = append(results, asset.(map[string]interface{})["id"])
// 		// } else {
// 		// 	results = append(results, asset.(map[string]interface{})["browser_download_url"])
// 		// }
// 		// }
// 	}

// 	// if !download {
// 	// 	// only print results
// 	// 	for _, res := range results {
// 	// 		fmt.Println(res)
// 	// 	}
// 	// } else {
// 	// Download results - parallel downloading, use channel to syncronize
// 	c := make(chan int)
// 	for _, res := range results {
// 		go downloadResource(repo, res.(float64), c)
// 	}
// 	// wait for downloads end
// 	for i := 0; i < len(results); i++ {
// 		<-c
// 	}
// 	// }
// 	return nil
// }
