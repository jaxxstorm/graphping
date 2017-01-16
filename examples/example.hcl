interval = 10
prefix = "graphping"


target_group "first" {
  interval = 2
  prefix = "search"
  target "google" {
    address = "www.google.co.uk"
  }
  target "bing" {
    address = "www.bing.com"
  }
}

target_group "second" {
  prefix = "uk"
  target "bbc" {
    address = "www.bbc.co.uk"
  }
}
