# <img src="https://raw.githubusercontent.com/devict/magopie/master/logo.png" style="background-color: transparent!" alt="Magopie" title="Magopie" /><br />[![GoDoc][godoc-badge]][godoc] [![Build Status][travis-badge]][travis] [![Build Status][goreport-badge]][goreport]

Your Personal Torrent Search Engine.  
*(It's a silent "O".)*

This project was an [honorable mention][gala-blog] in the 2016 Gopher Gala.

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

## Demo
![Demo][demogif]

Note: we were tunneling through ngrok for the demo to work around a local
networking problem.

## Notable Features
* Android client using Go bindings via [gomobile][gomobile].
* Downloaded `.torrent` files are saved to disk through an [afero][afero]
  filesystem interface providing for future extensibility to other remote afero
  implementations.
* Torrents are gathered by parsing Kick Ass Torrents XML feeds and by scraping
  the search page of The Pirate Bay using [goquery][goquery].
* Searches against upstream sites are performed concurrently.

## Future Plans
* Serve magopie over TLS with automatic integration with Let's Encrypt.
* Improve/replace the authentication mechanism.
* Create an iOS app using the gomobile bindings.

## Development
Go dependencies are vendored under `vendor/`. If you need to add or update a
dep you should use the [govendor tool][govendor].

Development currently targets Go version 1.6 but will probably work on 1.5 with
the vendor experiment enabled.

### Android
To work on the Android app do the following:

* Install [gomobile][gomobile].
* Install [Android Studio][android] and SDK level 22.
* Connect a device and ensure you can see it with `adb devices`.
* In Android Studio set your `GOPATH`. First open the settings dialog at
  `Android Studio > Preferences` on Mac OS X or under `File > Settings` on
  Linux and Windows. From there go to `Build, Execution, Deployment > Path
  Variables` and add the variable `GOPATH` that points to your normal `$GOPATH`
  such as `/Users/jwalker/go`.
* Open the project `$GOPATH/src/github.com/devict/magopie/cmd/android/magopie`.
* Run the `app` module.

## Licenses
Magopie is licensed with the [MIT license](LICENSE).

Margo, the Magopie mascot, is a derivative work of the Go gopher, and thus, is licensed under 
the Creative Commons 3.0 Attributions license.

The Go gopher was designed by Renee French. http://reneefrench.blogspot.com/
The design is licensed under the Creative Commons 3.0 Attributions license.
Read this article for more details: https://blog.golang.org/gopher

[godoc]: https://godoc.org/github.com/devict/magopie "GoDoc"
[godoc-badge]: https://godoc.org/github.com/devict/magopie?status.svg "GoDoc Badge"
[travis]: https://travis-ci.org/devict/magopie "Travis CI"
[travis-badge]: https://travis-ci.org/devict/magopie.svg?branch=master
[goreport]: http://goreportcard.com/report/devict/magopie "Go Report Card"
[goreport-badge]: http://goreportcard.com/badge/devict/magopie "Go Report Card Badge"
[transmission]: http://www.transmissionbt.com/ "Transmission"
[gomobile]: https://github.com/golang/mobile "gomobile"
[android]: http://developer.android.com/sdk/index.html "Android Studio"
[afero]: https://github.com/spf13/afero "Afero"
[goquery]: https://github.com/PuerkitoBio/goquery "goquery"
[demogif]: http://i.imgur.com/cLshfTl.gif "Demo"
[govendor]: https://github.com/kardianos/govendor "govendor"
[gala-blog]: http://gophergala.com/blog/gopher/gala/2016/02/05/winners-2016/ "Gopher Gala blog"
