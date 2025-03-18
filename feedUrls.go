package main

var NewsFeeds = map[string][]string{
	"technology": {
		"https://www.wired.com/feed/rss",
		"https://www.theverge.com/rss/index.xml",
		"https://feeds.arstechnica.com/arstechnica/index",
	},
	"business": {
		"https://www.forbes.com/innovation/feed",
		"https://www.business-standard.com/rss/home_page_top_stories.rss",
	},
	"sports": {
		"https://www.espn.com/espn/rss/news",
		"https://www.sports-reference.com/blog/feed/",
		"https://api.foxsports.com/v1/rss?partnerKey=zBaFxRyGKCfxBagJG9b8pqLyndmvo7UU",
	},
	"entertainment": {
		"https://www.entertainment-focus.com/feed/",
		"https://deadline.com/feed/",
		"https://variety.com/feed/",
	},
	"science": {
		"https://www.livemint.com/rss/science",
		// "https://www.sciencedaily.com/rss/all.xml",
		"https://www.livescience.com/feeds/all",
		"https://scitechdaily.com/feed/",
	},
	"world": {
		"https://feeds.bbci.co.uk/news/world/rss.xml",
		"https://rss.nytimes.com/services/xml/rss/nyt/World.xml",
		"http://feeds.washingtonpost.com/rss/world",
	},
	"health": {
		"https://www.statnews.com/feed",
		"https://abcnews.go.com/abcnews/healthheadlines",
		"https://feeds.npr.org/103537970/rss.xml",
	},
	"ai": {
		"https://www.artificialintelligence-news.com/feed/",
		"https://www.unite.ai/feed/",
		"https://www.analyticsinsight.net/feed/",
	},
	"hollywood": {
		"https://feeds.feedburner.com/ndtvmovies-hollywood",
		"https://deadline.com/feed/",
		"https://variety.com/feed/",
		"https://www.hollywoodreporter.com/feed",
	},
	"defence": {
		"https://www.defensenews.com/arc/outboundfeeds/rss/",
		"https://breakingdefense.com/feed",
	},
	"politics": {
		"https://feeds.nbcnews.com/nbcnews/public/politics",
		"https://rss.politico.com/politics-news.xml",
		"https://thehill.com/homenews/feed/",
		"https://www.realclearpolitics.com/index.xml",
	},
	"automobile": {
		"https://www.caranddriver.com/rss/all.xml/",
		"https://www.autocar.co.uk/rss",
	},
	"space": {
		"https://www.space.com/feeds/all",
		"https://spacenews.com/feed/",
		"https://www.universetoday.com/feed/",
	},
	"economy": {
		"https://www.economist.com/finance-and-economics/rss.xml",
		"https://feeds.marketwatch.com/marketwatch/topstories/",
	},
	"bollywood": {
		"https://feeds.feedburner.com/ndtvmovies-bollywood",
		"https://www.bollywoodhungama.com/rss/news.xml",
		// "https://www.filmibeat.com/rss/feeds/bollywood-fb.xml",
		"https://timesofindia.indiatimes.com/rssfeeds/1081479906.cms",
	},
}

func GetFeedUrls(category string) []string {
	return NewsFeeds[category]
}
