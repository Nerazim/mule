package scraper

import (
  "fmt"
  "io/ioutil"
  "net/http"
  "regexp"
  "github.com/moovweb/gokogiri"
  "github.com/moovweb/gokogiri/xml"
)

func ScrapeAllTheThings(url string) {
  pageSource := retrievePageSource(url)

  doc, err := gokogiri.ParseHtml(pageSource)
  errorHandler(err)
  defer doc.Free()
  matches, err := doc.Search(".//*[@class='item-container clearfix match collapsed']")
  errorHandler(err)

  fmt.Println(parseMatches(matches))
}

func retrievePageSource(url string) []byte {
  resp, err := http.Get(url)
  errorHandler(err)
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  errorHandler(err)
  return body
}

func parseMatches(matches []xml.Node) []match {
  out := []match{}
  for i := range matches {
    temp := match{}
    players, err := matches[i].Search(".//*[@class='player-name']")
    errorHandler(err)
    if len(players) == 3 {
      temp.player1Name = players[0].Content()
      temp.player2Name = players[1].Content()
      temp.winrar= players[2].Content()
    }

    races, err := matches[i].Search(".//*[@class='race-icon']")
    errorHandler(err)
    if len(races) == 2 {
      temp.player1Race = parseRace(races[0].String())
      temp.player2Race = parseRace(races[1].String())
    }

    res, err := matches[i].Search(".//dd")
    errorHandler(err)
    if len(res) == 5 {
      temp.matchDate = res[4].Content()
      temp.teams = res[2].Content()
    }
    out = append(out, temp)
  }
  return out
}

func parseRace(html string) string {
  regex, err := regexp.Compile(`Protoss|Terran|Zerg`)
  errorHandler(err)
  return regex.FindString(html)
}
