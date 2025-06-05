package crawl

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"

	"github.com/pkg/errors"
	"github.com/temoto/robotstxt"

	"github.com/glidea/zenfeed/pkg/util/text_convert"
)

type Crawler interface {
	Markdown(ctx context.Context, u string) ([]byte, error)
}

type local struct {
	hc *http.Client

	robotsDataCache sync.Map
}

func NewLocal() Crawler {
	return &local{
		hc: &http.Client{},
	}
}

func (c *local) Markdown(ctx context.Context, u string) ([]byte, error) {
	// Check if the page is allowed.
	if err := c.checkAllowed(ctx, u); err != nil {
		return nil, errors.Wrapf(err, "check robots.txt for %s", u)
	}

	// Prepare the request.
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "create request for %s", u)
	}
	req.Header.Set("User-Agent", userAgent)

	// Send the request.
	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "fetch %s", u)
	}
	defer func() { _ = resp.Body.Close() }()

	// Parse the response.
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("received non-200 status code %d from %s", resp.StatusCode, u)
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "read body from %s", u)
	}

	// Convert the body to markdown.
	mdBytes, err := textconvert.HTMLToMarkdown(bodyBytes)
	if err != nil {
		return nil, errors.Wrap(err, "convert html to markdown")
	}

	return mdBytes, nil
}

const userAgent = "ZenFeed"

func (c *local) checkAllowed(ctx context.Context, u string) error {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return errors.Wrapf(err, "parse url %s", u)
	}

	d, err := c.getRobotsData(ctx, parsedURL.Host)
	if err != nil {
		return errors.Wrapf(err, "check robots.txt for %s", parsedURL.Host)
	}
	if !d.TestAgent(parsedURL.Path, userAgent) {
		return errors.Errorf("disallowed by robots.txt for %s", u)
	}

	return nil
}

// getRobotsData fetches and parses robots.txt for a given host.
func (c *local) getRobotsData(ctx context.Context, host string) (*robotstxt.RobotsData, error) {
	// Check the cache.
	if data, found := c.robotsDataCache.Load(host); found {
		return data.(*robotstxt.RobotsData), nil
	}

	// Prepare the request.
	robotsURL := fmt.Sprintf("https://%s/robots.txt", host)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, robotsURL, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "create request for %s", robotsURL)
	}
	req.Header.Set("User-Agent", userAgent)

	// Send the request.
	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "fetch %s", robotsURL)
	}
	defer func() { _ = resp.Body.Close() }()

	// Parse the response.
	switch resp.StatusCode {
	case http.StatusOK:
		data, err := robotstxt.FromResponse(resp)
		if err != nil {
			return nil, errors.Wrapf(err, "parse robots.txt from %s", robotsURL)
		}
		c.robotsDataCache.Store(host, data)

		return data, nil

	case http.StatusNotFound:
		data := &robotstxt.RobotsData{}
		c.robotsDataCache.Store(host, data)

		return data, nil

	case http.StatusUnauthorized, http.StatusForbidden:
		return nil, errors.Errorf("access to %s denied (status %d)", robotsURL, resp.StatusCode)
	default:
		return nil, errors.Errorf("unexpected status %d fetching %s", resp.StatusCode, robotsURL)
	}
}

type jina struct {
	hc    *http.Client
	token string
}

func NewJina(token string) Crawler {
	return &jina{
		hc: &http.Client{},

		// If token is empty, will not affect to use, but rate limit will be lower.
		// See https://jina.ai/api-dashboard/rate-limit.
		token: token,
	}
}

func (c *jina) Markdown(ctx context.Context, u string) ([]byte, error) {
	proxyURL := "https://r.jina.ai/" + u
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, proxyURL, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "create request for %s", u)
	}

	req.Header.Set("X-Engine", "browser")
	req.Header.Set("X-Robots-Txt", userAgent)
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "fetch %s", proxyURL)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("received non-200 status code %d from %s", resp.StatusCode, proxyURL)
	}

	mdBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "read body from %s", proxyURL)
	}

	return mdBytes, nil
}
