package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/chromedp/chromedp"
)

func TakeScreenShot(ctx context.Context, url string) ([]byte, error) {
	// options := []chromedp.ExecAllocatorOption{}
	// options = append(options, chromedp.DefaultExecAllocatorOptions[:]...)
	// options = append(options, chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"))
	// options = append(options, chromedp.Flag("disable-dev-shm-usage", true))
	// options = append(options, chromedp.Flag("ignore-certificate-errors", true)) // RIP shittyproxy.go
	// options = append(options, chromedp.WindowSize(250, 250))

	// Through trial and error we've landed on the following options to pass to chromedp/headless
	// chrome. This is a combination of chromedp's DefaultExecAllocatorOptions as well as various
	// options mentioned for tooling by the chrome team. See:
	// https://pkg.go.dev/github.com/chromedp/chromedp#pkg-variables
	// https://github.com/GoogleChrome/chrome-launcher/blob/main/docs/chrome-flags-for-tools.md
	chromeOptions := []chromedp.ExecAllocatorOption{
		// NOTE: Currently we're utilizing chromedp's "docker-headless-shell" project which provides a
		//       custom build of the "old" headless chrome. This is fine for now, but that may change
		//       as the chrome team's focus shifts to the new "native" headless mode.
		//       see: https://developer.chrome.com/articles/new-headless/
		chromedp.Headless,

		// Features to enable/disable (https://niek.github.io/chrome-features/)
		chromedp.Flag("enable-features", "NetworkService,NetworkServiceInProcess2"),
		// TODO: remove AvoidUnnecessaryBeforeUnloadCheckSync below
		// once crbug.com/1324138 is fixed and released.
		// AcceptCHFrame disabled because of crbug.com/1348106.
		chromedp.Flag("disable-features", "site-per-process,Translate,BlinkGenPropertyTrees,InterestFeedContentSuggestions,CalculateNativeWinOcclusion,DialMediaRouteProvider,OptimizationHints,MediaRouter,BackForwardCache,AcceptCHFrame,AvoidUnnecessaryBeforeUnloadCheckSync"),

		// Other flags from chromedp.DefaultExecAllocatorOptions
		// see: https://pkg.go.dev/github.com/chromedp/chromedp#pkg-variables
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Flag("disable-background-networking", true),
		chromedp.Flag("disable-background-timer-throttling", true),
		chromedp.Flag("disable-backgrounding-occluded-windows", true),
		chromedp.Flag("disable-breakpad", true),
		chromedp.Flag("disable-client-side-phishing-detection", true),
		chromedp.Flag("disable-default-apps", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-hang-monitor", true),
		chromedp.Flag("disable-ipc-flooding-protection", true),
		chromedp.Flag("disable-popup-blocking", true),
		chromedp.Flag("disable-prompt-on-repost", true),
		chromedp.Flag("disable-renderer-backgrounding", true),
		chromedp.Flag("disable-sync", true),
		chromedp.Flag("force-color-profile", "srgb"),
		chromedp.Flag("metrics-recording-only", true),
		chromedp.Flag("safebrowsing-disable-auto-update", true),
		chromedp.Flag("enable-automation", true),
		chromedp.Flag("password-store", "basic"),
		chromedp.Flag("use-mock-keychain", true),

		// Other flags mentioned for tooling
		// see: https://github.com/GoogleChrome/chrome-launcher/blob/main/docs/chrome-flags-for-tools.md
		chromedp.Flag("disable-component-extensions-with-background-pages", true),
		chromedp.Flag("autoplay-policy", "user-gesture-required"),
		chromedp.Flag("deny-permission-prompts", true),
		chromedp.Flag("disable-notifications", true),
		chromedp.Flag("block-new-web-contents", true),
		chromedp.Flag("noerrdialogs", true),
		chromedp.Flag("disable-component-update", true),
		chromedp.Flag("disable-domain-reliability", true),
		chromedp.Flag("no-default-browser-check", true),
		chromedp.Flag("no-pings", true),
		chromedp.Flag("no-service-autorun", true),
		chromedp.Flag("disable-crash-reporter", true),
		// https://github.com/GoogleChrome/chrome-launcher/blob/f1baa9af20db9c1e3c55187ba9e91c6f1cee99e9/src/flags.ts#LL60
		chromedp.Flag("force-fieldtrials", "*BackgroundTracing/default/"),

		// Flags that we believe help us avoid various issues
		chromedp.IgnoreCertErrors,
		// All current literature states that "DisableGPU" is no longer needed, but it seems to solve
		// some problems (and our workers don't have a GPU anyway).
		// e.g. https://github.com/chromedp/chromedp/issues/904
		chromedp.DisableGPU,
		// https://stackoverflow.com/questions/69037458/selenium-chromedriver-gives-initializesandbox-called-with-multiple-threads-in
		chromedp.Flag("disable-software-rasterizer", true),

		// Override the default user agent
		// chromedp.UserAgent(userAgent),
	}

	// options = append(options, chromedp.Flag("--disable-software-rasterizer", true))

	context, cancel := chromedp.NewExecAllocator(ctx, chromeOptions...)
	defer cancel()

	browserCtx, cancelBrowserCtx := chromedp.NewContext(context)
	defer cancelBrowserCtx()

	var filebyte []byte
	if err := chromedp.Run(browserCtx, chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.Sleep(1 * time.Second),
		chromedp.CaptureScreenshot(&filebyte),
	}); err != nil {
		return nil, err
	}
	return filebyte, nil
}

func main() {
	// url := "https://ord.zuexeuz.net/preview/fff38a9f9375e675001de89c11d281587b2796eadf56d1849b89739865f64498i0"
	url := "https://ord.zuexeuz.net/preview/f59f1a5cf10c8b081f503b397e8e3813822c3ace68b7b601fea6ae3e8253c2fai0"
	// url := "https://ord.zuexeuz.net/preview/821dcd5da03e0a6771ecac114ba1054690183f42b72ae8303e84260a654356c7i0"
	ctx := context.Background()

	data, err := TakeScreenShot(ctx, url)
	if err != nil {
		panic(err)
	}

	defer ctx.Done()

	pngFile, err := os.Create("./shot.png")
	if err != nil {
		panic(err)
	}

	defer pngFile.Close()

	pngFile.Write(data)
	fmt.Println("screen shot taken")
}
