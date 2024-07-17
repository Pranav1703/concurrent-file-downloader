package main

import (
    "fmt"
    "io"
    "net/http"
    "os"
    "sync"
)

type Download struct {
	Url string
	FileName string
}

func main() {
	downloadLinks := []Download{
		{Url:"https://jsonplaceholder.typicode.com/posts/1",FileName:"Test-download-1.txt"},  //add download links and filename for it here
		{Url: "https://randomuser.me/api/",FileName: "RandomUserData.txt"},
	}

	var wg sync.WaitGroup;
	resultCh := make(chan string,len(downloadLinks))
	errorCh := make(chan error,len(downloadLinks))

	for _,download := range(downloadLinks){
		wg.Add(1)

		go func(download Download) {
			defer wg.Done()

			err := DownloadFile(download.Url,download.FileName)
			if err!=nil {
				errorCh <- err
			}else{
				resultCh <- fmt.Sprintf("Downloaded: %s",download.FileName)
			}

		}(download)
	}

	go func() {
		wg.Wait()
		close(resultCh)
		close(errorCh)
	}()
	
	for result := range resultCh{
		fmt.Println(result)
	}

	for err := range errorCh {
		fmt.Println(err)
	}

}

func DownloadFile(url string, FileName string) error{
	res,err := http.Get(url)
	if err!=nil {
		fmt.Println(err)
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err!=nil {
		fmt.Println(err)
	}

	return os.WriteFile(FileName,data,0664)

}