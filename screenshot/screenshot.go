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
	Url     string  // query url
	Width   int     // width of view port (default 1920)
	Height  int     // height of view port (default 1024)
	Scale   float64 // scale of final image (default 0.1, meaning that with default settings our image will be 192x102px)
	Wait    int     // ms to wait before rendering (default 10000, 10 seconds)
	Timeout int     // ms to wait before returning timeout (default 60000, 60 seconds)
	Quality int     // png quality
	Hash    string  // file hash for this setting to quickly find cached files
	MaxAge  int     // maximum age of cache file in seconds (!) before it gets reloaded (set to 0 to never renew files, default 2592000 = 30 days)
}

// entry point to create screenshot
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

func runCreateScreenShot(settings Settings, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(settings.Url),
		chromedp.Sleep(time.Duration(settings.Wait) * time.Millisecond),
		chromedp.ActionFunc(func(ctx context.Context) error {
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
		}),
	}
}
