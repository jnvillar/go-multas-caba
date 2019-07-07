package browser

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"sync"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

var captchas = []string{"RED CICLOVÍAS PROTEGIDAS", "MASCOTAS EN EL SUBTE"}

func getSumOfTransitFines(domain, captcha string) string {
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	var sumOfFines string
	err := chromedp.Run(ctx,
		chromedp.Navigate("https://www.buenosaires.gob.ar/consulta-de-infracciones"),
		chromedp.WaitVisible("edit-dominio", chromedp.ByID),
		chromedp.SendKeys("#edit-dominio", domain),
		chromedp.WaitVisible("input#brand_cap_answer"),
		chromedp.SendKeys("input#brand_cap_answer", captcha),
		chromedp.Click("button#edit-submit", chromedp.NodeVisible),
		chromedp.WaitReady("#actasComprobantes-view"),
		chromedp.EvaluateAsDevTools(`document.querySelector("h4 > strong").innerText`, &sumOfFines),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n error %v \n", err.Error())
	}
	return sumOfFines
}

func TransitFines(domain string) string {
	var resChan = make(chan string)
	var sumOfFines string
	var finished bool
	for _, captcha := range captchas {
		go func(captcha string) {
			var res string
			// retry three times with the same captcha
			for i := 0; i < 3 && res == "" && !finished; i++ {
				res = getSumOfTransitFines(domain, captcha)
			}
			resChan <- res
		}(captcha)
	}
	for i := 0; i < len(captchas) && sumOfFines == ""; i++ {
		sumOfFines = <-resChan
	}
	finished = true
	return sumOfFines
}

func fullScreenShot(route string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			// get layout metrics
			_, _, contentSize, err := page.GetLayoutMetrics().Do(ctx)
			if err != nil {
				return err
			}
			width, height := int64(math.Ceil(contentSize.Width)), int64(math.Ceil(contentSize.Height))
			// force viewport emulation
			err = emulation.SetDeviceMetricsOverride(width, height, 1, false).
				WithScreenOrientation(&emulation.ScreenOrientation{
					Type:  emulation.OrientationTypePortraitPrimary,
					Angle: 0,
				}).
				Do(ctx)
			if err != nil {
				return err
			}
			// capture screenshot
			res, err := page.CaptureScreenshot().
				WithQuality(50).
				WithClip(&page.Viewport{
					X:      contentSize.X,
					Y:      contentSize.Y,
					Width:  contentSize.Width,
					Height: contentSize.Height,
					Scale:  1,
				}).Do(ctx)
			if err != nil {
				return err
			}
			if err := ioutil.WriteFile(fmt.Sprintf("%s.png", route), res, 0644); err != nil {
				log.Fatal(err)
			}
			return nil
		}),
	}
}

func WaitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}
