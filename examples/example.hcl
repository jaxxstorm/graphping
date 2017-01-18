interval = 10
prefix = "graphping"


target_group "search_enginers" {
  interval = 2
  prefix = "search"
  target "google" {
    address = "www.google.co.uk"
  }
  target "bing" {
    address = "www.bing.com"
  }
}

target_group "news_sites" {
  prefix = "uk"
  target "bbc" {
    address = "www.bbc.co.uk"
  }
}
