package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/chromedp/chromedp"
)

func TakeScreenShot(ctx context.Context, url string) ([]byte, error) {
	ctx, cancel := chromedp.NewContext(ctx,
		chromedp.WithDebugf(log.Printf),
	)
	defer cancel()
	var buf []byte
	if err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(1*time.Second),
		chromedp.CaptureScreenshot(&buf),
	); err != nil {
		return nil, err
	}
	return buf, nil
}

func main() {
	// url := "https://ord.zuexeuz.net/preview/f59f1a5cf10c8b081f503b397e8e3813822c3ace68b7b601fea6ae3e8253c2fai0" // DOES NOT WORK
	url := "https://ord.zuexeuz.net/preview/c145fa1fc4c4b9e5cf7b839207502d28d31fd6cb3d0483470ac1c0d15ad24735i0"
	// url := "https://ord.zuexeuz.net/preview/821dcd5da03e0a6771ecac114ba1054690183f42b72ae8303e84260a654356c7i0" // WORKS
	data, err := TakeScreenShot(context.Background(), url)
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile("shot.png", data, 0o644); err != nil {
		panic(err)
	}
}
