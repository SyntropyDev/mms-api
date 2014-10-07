package model_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"testing"

	"github.com/ChimeraCoder/anaconda"
	"github.com/SyntropyDev/mms-api/model"
	"github.com/huandu/facebook"
	"github.com/jteeuwen/go-pkg-rss"
)

const (
	feedATOM = "http://www.theverge.com/rss/frontpage"
	feedRSS  = "http://www.techmeme.com/feed.xml"
)

var (
	itemRSS = &feeder.Item{
		Title: `Oculus Reveals Its New "Crescent Bay" Developer Kit With 360-Degree Head Tracking And Headphones (Josh Constine/TechCrunch)`,
		Description: `<A HREF="http://techcrunch.com/2014/09/20/oculus-connect-announcements/"><IMG VSPACE="4" HSPACE="4" BORDER="0" ALIGN="RIGHT" SRC="http://www.techmeme.com/140920/i8.jpg"></A>
<P><A HREF="http://www.techmeme.com/140920/p8#a140920p8" TITLE="Techmeme permalink"><IMG WIDTH=11 HEIGHT=12 SRC="http://www.techmeme.com/img/pml.png" STYLE="border:none;padding:0;margin:0;"></A> Josh Constine / <A HREF="http://techcrunch.com/">TechCrunch</A>:<BR>
<SPAN STYLE="font-size:1.3em;"><B><A HREF="http://techcrunch.com/2014/09/20/oculus-connect-announcements/">Oculus Reveals Its New &ldquo;Crescent Bay&rdquo; Developer Kit With 360-Degree Head Tracking And Headphones</A></B></SPAN>&nbsp; &mdash;&nbsp; Oculus gave the world the first look at its new developer kit Crescent Bay today at Oculus' Connect conference, which you can watch live here.&nbsp; Crescent Bay has a faster frame rate &hellip; </P>`,
		Guid:    sPtr("http://www.techmeme.com/140920/p8#a140920p8"),
		PubDate: "Sat, 20 Sep 2014 13:30:06 -0400",
		Links: []*feeder.Link{
			{Href: "http://www.techmeme.com/140920/p8#a140920p8"},
		},
	}

	// &{Title: Links:[0xc2080e04c0]}

	itemATOM = &feeder.Item{
		Title:       `NYPD goes back to Twitter school after photo contest backfires`,
		Description: "",
		Guid:        nil,
		PubDate:     "2014-09-21T03:08:02-04:00",
		Id:          `http://www.theverge.com/2014/9/21/6663163/nypd-goes-back-to-twitter-school-after-photo-content-backfires`,
		Content: &feeder.Content{
			Type: "html",
			Text: `<img alt="" src="http://cdn2.vox-cdn.com/uploads/chorus_image/image/39185966/nypd-police-department-station-stock_1020.0_standard_800.0.jpg"/>
			<p>Police departments around North America have gradually, if warily, come to embrace social media as a way of communicating with their communities. But things can quickly go awry on Facebook and Twitter, a lesson the New York Police Department learned in April when a photo contest designed to showcase friendly photos of citizens with police officers was taken over by images of police <a href="http://www.theverge.com/2014/4/22/5641266/nypd-twitter-photo-contest-backfires">aggressively subduing people</a>. A month after the #myNYPD fiasco, top NYPD officials have been taking courses in using Twitter effectively, <a target="_blank" href="http://online.wsj.com/articles/officers-train-with-nypds-twitter-police-1411002492">the <i>Wall Street Journal</i> reports</a>.</p>
			<p><a href="http://www.theverge.com/2014/9/21/6663163/nypd-goes-back-to-twitter-school-after-photo-content-backfires">Continue reading&hellip;</a></p>`,
		},
		Links: []*feeder.Link{
			{Href: "http://www.theverge.com/2014/9/21/6672811/vr-typing-trainer-for-oculus-rift"},
		},
	}
)

func sPtr(s string) *string { return &s }

func TestRSS(t *testing.T) {
	m := &model.Member{ID: 1, Name: "Test"}
	f := &model.Feed{ID: 1, Type: string(model.FeedTypeRSS), Identifier: feedRSS}
	story := model.NewStoryRSS(m, f, itemRSS)

	if story.ImagesRaw != "http://www.techmeme.com/140920/i8.jpg" {
		t.Fatal("image not found")
	}
	if story.LinksRaw != "http://www.techmeme.com/140920/p8#a140920p8" {
		t.Fatal("links not found")
	}
	if story.Timestamp != 1411234206000 {
		t.Fatal("timestamp incorrect")
	}
	if story.SourceID != "http://www.techmeme.com/140920/p8#a140920p8" {
		t.Fatal("source id incorrect")
	}
}

func TestAtom(t *testing.T) {
	m := &model.Member{ID: 1, Name: "Test"}
	f := &model.Feed{ID: 1, Type: string(model.FeedTypeRSS), Identifier: feedATOM}
	story := model.NewStoryRSS(m, f, itemATOM)

	if story.ImagesRaw != "http://cdn2.vox-cdn.com/uploads/chorus_image/image/39185966/nypd-police-department-station-stock_1020.0_standard_800.0.jpg" {
		t.Fatal("image not found")
	}
	if story.LinksRaw != "http://www.theverge.com/2014/9/21/6672811/vr-typing-trainer-for-oculus-rift" {
		t.Fatal("links not found")
	}
	if story.Timestamp != 1411283282000 {
		t.Fatal("timestamp incorrect")
	}
	if story.SourceID != "http://www.theverge.com/2014/9/21/6663163/nypd-goes-back-to-twitter-school-after-photo-content-backfires" {
		t.Fatal("source id incorrect")
	}
}

func TestFacebook(t *testing.T) {
	initConfig()
	app := facebook.New(os.Getenv("facebookApiID"), os.Getenv("facebookAppSecret"))

	app.RedirectUri = "http://syntropy.io"
	session := app.Session(app.AppAccessToken())
	result, err := session.Api("/syntropydevelopment/posts", facebook.GET, nil)
	if err != nil {
		t.Fatal(err)
	}

	j, err := json.Marshal(result)
	if err != nil {
		t.Fatal(err)
	}

	type fbookResp struct {
		Data []struct {
			Description string
			Message     string
			Story       string
			Title       string
		}
	}

	resp := &fbookResp{}
	if err := json.Unmarshal(j, resp); err != nil {
		t.Fatal(err)
	}

	for _, data := range resp.Data {
		fmt.Println(data.Message)
	}
}

func initConfig() error {
	b, err := ioutil.ReadFile("../config.json")
	if err != nil {
		return err
	}

	config := map[string]string{}
	if err := json.Unmarshal(b, &config); err != nil {
		return err
	}

	for key, value := range config {
		os.Setenv(key, value)
	}
	return nil
}

func TestTwitter(t *testing.T) {
	b, err := ioutil.ReadFile("/Users/logan/go/src/github.com/SyntropyDev/mms-api/config/config.json")
	if err != nil {
		t.Fatal(err)
	}

	config := map[string]string{}
	if err := json.Unmarshal(b, &config); err != nil {
		t.Fatal(err)
	}

	anaconda.SetConsumerKey(config["twitterApiKey"])
	anaconda.SetConsumerSecret(config["twitterApiSecret"])
	api := anaconda.NewTwitterApi("", "")

	v := url.Values{}
	v.Set("screen_name", "dubvNOW")
	v.Set("include_rts", "false")

	tweets, err := api.GetUserTimeline(v)
	if err != nil {
		t.Fatal(err)
	}

	for _, tweet := range tweets {
		fmt.Println(tweet)
	}
}
