package screenshot

import (
	"context"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

// struct to define screenshot settings
type Settings struct {
	Url         string  // query url
	Width       int     // width of view port (default 1920)
	Height      int     // height of view port (default 1024)
	Scale       float64 // scale of final image (default 0.1, meaning that with default settings our image will be 192x102px)
	WaitForIdle bool    // wait for the page to be idle - wait for networkIdle event in Browser (default false)
	Wait        int     // ms to wait before rendering - only meaningful if WaitForIdle is false (default 10000, 10 seconds)
	Timeout     int     // ms to wait before returning timeout (default 60000, 60 seconds)
	Quality     int     // png quality
	Hash        string  // file hash for this setting to quickly find cached files
	MaxAge      int     // maximum age of cache file in seconds (!) before it gets reloaded (set to 0 to never renew files, default 2592000 = 30 days)
	Force       bool    // force update
}

// CreateScreenShot is the entry point to create screenshot
func CreateScreenShot(settings Settings) ([]byte, error) {
	// create context with timeout
	timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), time.Duration(settings.Timeout)*time.Millisecond)
	defer timeoutCancel()
	ctx, cancel := chromedp.NewContext(timeoutCtx)
	defer cancel()

	// start capture process using chromedp
	var buf []byte
	if err := chromedp.Run(ctx, runCreateScreenShot(settings, &buf)); err != nil {
		return nil, err
	}

	return buf, nil
}

// runCreateScreenShot is the actual runner for creating screenshots
func runCreateScreenShot(settings Settings, res *[]byte) chromedp.Tasks {
	// wait for networkIdle event in Browser
	if settings.WaitForIdle {
		return chromedp.Tasks{
			navigateAndWaitFor(settings.Url, "networkIdle"),
			createScreenShot(settings, res),
		}
	}

	// wait for a fixed number of milliseconds
	return chromedp.Tasks{
		chromedp.Navigate(settings.Url),
		chromedp.Sleep(time.Duration(settings.Wait) * time.Millisecond),
		createScreenShot(settings, res),
	}
}

// createScreenShot is the function that will create screenshots - saved into res as bytes
func createScreenShot(settings Settings, res *[]byte) chromedp.ActionFunc {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		// force viewport emulation
		err := emulation.SetDeviceMetricsOverride(int64(settings.Width), int64(settings.Height), 1, false).
			WithScreenOrientation(&emulation.ScreenOrientation{
				Type:  emulation.OrientationTypeLandscapePrimary,
				Angle: 0,
			}).
			Do(ctx)
		if err != nil {
			return err
		}

		// capture screenshot
		*res, err = page.CaptureScreenshot().
			WithQuality(int64(settings.Quality)).
			WithClip(&page.Viewport{
				X:      0,
				Y:      0,
				Width:  float64(settings.Width),
				Height: float64(settings.Height),
				Scale:  settings.Scale,
			}).Do(ctx)
		if err != nil {
			return err
		}
		return nil
	})
}

// taken from https://github.com/chromedp/chromedp/issues/431 - thanks go to https://github.com/wietsevenema
func navigateAndWaitFor(url string, eventName string) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		_, _, _, err := page.Navigate(url).Do(ctx)
		if err != nil {
			return err
		}

		return waitFor(ctx, eventName)
		return nil
	}
}

// waitFor blocks until eventName is received.
// Examples of events you can wait for:
//
//	init, DOMContentLoaded, firstPaint,
//	firstContentfulPaint, firstImagePaint,
//	firstMeaningfulPaintCandidate,
//	load, networkAlmostIdle, firstMeaningfulPaint, networkIdle
//
// This is not super reliable, I've already found incidental cases where
// networkIdle was sent before load. It's probably smart to see how
// puppeteer implements this exactly.
// taken from https://github.com/chromedp/chromedp/issues/431 - thanks go to https://github.com/wietsevenema
func waitFor(ctx context.Context, eventName string) error {
	ch := make(chan struct{})
	cctx, cancel := context.WithCancel(ctx)
	chromedp.ListenTarget(cctx, func(ev interface{}) {
		switch e := ev.(type) {
		case *page.EventLifecycleEvent:
			if e.Name == eventName {
				cancel()
				close(ch)
			}
		}
	})
	select {
	case <-ch:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}

}
