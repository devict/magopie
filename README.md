# ![Magopie][magopietext] [![GoDoc][godoc-badge]][godoc] [![Build Status][travis-badge]][travis]
![Magopie Mascot][mascot]

Your Personal Torrent Search Engine. *(It's a silent "O".)*


## What's with the name?
Magopie is pronounced just like "magpie" the bird. We chose the name because
magpies are known for collecting things (especially if they're shiny). We added
the "silent o" mostly for fun but also to differentiate from an existing Python
project named `magpie`.

## Use Cases
1. I often want to search for some torrent but don't want to be subjected to
   NSFW ads that are prevalent on many torrent sites. Magopie proxies my query
   to multiple upstream torrent search engines and collates the results into a
   basic collection.
2. I want to remotely start a torrent on a home server or seedbox. Magopie does
   not participate in any peer-to-peer file sharing itself; it merely downloads
   a `.torrent` file to a preconfigured directory on the server. Another
   program such as [Transmission][transmission] needs to be configured to watch
   that same directory for new `.torrent` files.
3. I want to do all of that from my phone.

## Notable Features
* Android client using Go bindings via [gomobile][gomobile].
* Downloaded `.torrent` files are saved to disk through an [afero][afero]
  filesystem interface providing for future extensibility to other remote afero
  implementations.
* Torrents are gathered by parsing Kick Ass Torrents XML feeds and by scraping
  the search page of The Pirate Bay using [goquery][goquery].
* Searches against upstream sites are performed concurrently.

## Installation
TODO

## License
Magopie is licensed with the [MIT license](LICENSE).

[godoc]: https://godoc.org/github.com/gophergala2016/magopie "GoDoc"
[godoc-badge]: https://godoc.org/github.com/gophergala2016/magopie?status.svg "GoDoc Badge"
[travis]: https://travis-ci.org/gophergala2016/magopie "Travis CI"
[travis-badge]: https://travis-ci.org/gophergala2016/magopie.svg?branch=master
[transmission]: http://www.transmissionbt.com/ "Transmission"
[gomobile]: https://github.com/golang/mobile "gomobile"
[afero]: https://github.com/spf13/afero "Afero"
[goquery]: https://github.com/PuerkitoBio/goquery "goquery"
[mascot]: https://raw.githubusercontent.com/gophergala2016/magopie/master/magopie.png "Magopie Mascot"
[magopietext]: https://raw.githubusercontent.com/gophergala2016/magopie/master/magopie_logo_text.png "Magopie"
