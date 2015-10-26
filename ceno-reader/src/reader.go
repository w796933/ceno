package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	rss "github.com/jteeuwen/go-pkg-rss"
	"github.com/jteeuwen/go-pkg-xmlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nicksnyder/go-i18n/i18n"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"time"
)

/**
 * Handle the receipt of a new channel.
 * @param {*rss.Feed} feed - A pointer to the object representing the feed received from
 * @param {[]*rss.Channel} newChannels - An array of pointers to received channels
 */
func channelFeedHandler(feed *rss.Feed, newChannels []*rss.Channel) {
	return
}

/**
 * Handle the receipt of a new item.
 * @param {*rss.Feed} feed - A pointer to the object representing the feed received from
 * @param {*rss.Channel} channel - A pointer to the channel object the item was received from
 * @param {[]*rss.Item} newItems - An array of pointers to items received from the channel
 */
func itemFeedHandler(feed *rss.Feed, channel *rss.Channel, newItems []*rss.Item) {
	T, _ := i18n.Tfunc(os.Getenv(LANG_ENVVAR), DEFAULT_LANG)
	fmt.Println("Feed URL is", feed.Url)
	for _, item := range newItems {
		url := item.Links[0].Href
		bundleData, bundleStatus := GetBundle(url)
		if bundleStatus == Failure {
			fmt.Println(T("bundle_fail_err", map[string]string{
				"Url": url,
			}))
			continue
		}
		inserted := InsertFreenet(bundleData)
		if inserted == Success {
			saveErr := SaveItem(DBConnection, feed.Url, item)
			if saveErr != nil {
				fmt.Println(T("db_store_error_rdr", map[string]string{
					"Error": saveErr.Error(),
				}))
			} else {
				fmt.Println(T("insert_success_rdr", map[string]string{
					"Url": url,
				}))
			}
		} else {
			fmt.Println(T("insertion_fail_err"))
		}
	}
}

/**
 * Periodically polls an RSS or Atom feed for new items.
 * @param {string} URL - The address of the feed
 * @param {xmlx.CharsetFunc} charsetReader - A function for handling the charset of items
 */
func pollFeed(URL string, charsetReader xmlx.CharsetFunc) {
	// Poll every five seconds
	T, _ := i18n.Tfunc(os.Getenv(LANG_ENVVAR), DEFAULT_LANG)
	feed := rss.New(5, true, channelFeedHandler, itemFeedHandler)
	for {
		defer func() {
			r := recover()
			if r != nil {
				fmt.Println(T("feed_poll_err", map[string]string{
					"Url":   URL,
					"Error": "Panicked when fetching from feed",
				}))
			}
		}()
		if err := feed.Fetch(URL, charsetReader); err != nil {
			fmt.Println(T("feed_poll_err", map[string]string{
				"Url":   URL,
				"Error": err.Error(),
			}))
		}
		<-time.After(time.Duration(feed.SecondsTillUpdate() * 1e9))
	}
}

/**
 * Handle the following of a feed in a separate goroutine.
 * @param {chan Feed} requests - A channel through which descriptions of feeds to be followed are received
 */
func followFeeds(requests chan SaveFeedRequest) {
	T, _ := i18n.Tfunc(os.Getenv(LANG_ENVVAR), DEFAULT_LANG)
	for {
		request := <-requests
		feedInfo := request.FeedInfo
		fmt.Println("Got a request to handle a feed.")
		fmt.Println(feedInfo)
		saveErr := SaveFeed(DBConnection, feedInfo)
		if saveErr != nil {
			fmt.Println("Could not save")
			fmt.Println(saveErr)
			request.W.Write([]byte(T("db_store_error_rdr", map[string]interface{}{"Error": saveErr.Error()})))
			return
		} else {
			fmt.Println("Saved")
			request.W.Write([]byte(T("req_handle_success_rdr")))
		}
		if feedInfo.Charset == "" {
			go pollFeed(feedInfo.Url, nil)
		} else {
			charsetFn, found := CharsetReaders[feedInfo.Charset]
			if found {
				go pollFeed(feedInfo.Url, charsetFn)
			} else {
				go pollFeed(feedInfo.Url, nil)
			}
		}
	}
}

/**
 * Write a file listing items that have been inserted into Freenet
 * @param feedUrl - The URL of the feed from which the items were served
 * @param marshalledItems - The marshalled information about items to write
 */
func writeItemsFile(feedUrl string, marshalledItems []byte) error {
	filename := base64.StdEncoding.EncodeToString([]byte(feedUrl)) + ".json"
	location := path.Join(JSON_FILE_DIR, filename)
	return ioutil.WriteFile(location, marshalledItems, os.ModePerm)
}

/**
 * Handle requests to have a new RSS or Atom feed followed.
 * POST /follow {"url": string, "type": string, "charset": string}
 * @param {chan Feed} requests - A channel through which descriptions of feeds to be followed are received
 */
func followHandler(requests chan SaveFeedRequest) func(http.ResponseWriter, *http.Request) {
	T, _ := i18n.Tfunc(os.Getenv(LANG_ENVVAR), DEFAULT_LANG)
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Got request")
		if r.Method != "POST" {
			w.Write([]byte(T("method_not_impl_rdr")))
			return
		}
		feedInfo := Feed{}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&feedInfo); err != nil {
			fmt.Println("Error decoding JSON")
			fmt.Println(err)
			w.Write([]byte(T("invalid_follow_req_rdr")))
			return
		}
		requests <- SaveFeedRequest{feedInfo, w}
	}
}

/**
 * Handle requests to have an RSS or Atom feed unfollowed.
 * DELETE /unfollow {"url": string}
 */
func unfollowHandler(w http.ResponseWriter, r *http.Request) {
	T, _ := i18n.Tfunc(os.Getenv(LANG_ENVVAR), DEFAULT_LANG)
	if r.Method != "DELETE" {
		w.Write([]byte(T("method_not_impl_rdr")))
		return
	}
	deleteReq := DeleteFeedRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&deleteReq)
	if err != nil {
		fmt.Println("Error decoding JSON")
		fmt.Println(err)
		w.Write([]byte(T("invalid_unfollow_req_rdr")))
		return
	}
	deleteErr := DeleteFeed(DBConnection, deleteReq.Url)
	if deleteErr != nil {
		w.Write([]byte(T("feed_delete_err_rdr", map[string]interface{}{
			"Error": deleteErr.Error(),
		})))
	} else {
		w.Write([]byte(T("feed_delete_success_rdr")))
	}
}

/**
 * Handle a request to have JSON files describing feeds and articles generated and inserted into
 * the distributed store being used. Also creates files for distribution in json-files.
 */
func insertHandler(w http.ResponseWriter, r *http.Request) {
	feeds, feedErr := AllFeeds(DBConnection)
	if feedErr != nil {
		fmt.Println("Couldn't get feeds")
		fmt.Println(feedErr)
		return
	}
	writeFeedsErr := writeFeeds(feeds)
	if writeFeedsErr != nil {
		fmt.Println(writeFeedsErr)
		return
	}
	for _, feed := range feeds {
		items, itemsError := GetItems(DBConnection, feed.Url)
		if itemsError != nil {
			fmt.Println("Couldn't get items for " + feed.Url)
			fmt.Println(itemsError)
		} else {
			writeItemsErr := writeItems(feed.Url, items)
			if writeItemsErr != nil {
				fmt.Println("Could not write items for " + feed.Url)
			} else {
				fmt.Println("Success!")
			}
		}
	}
}

/**
 * Write information about feeds being followed.
 * Try to insert JSON containing this information to Freenet and, only if that succeeds,
 * write the JSON to a file for distribution with the client.
 * @param feeds - Information about the feeds being followed
 */
func writeFeeds(feeds []Feed) error {
	T, _ := i18n.Tfunc(os.Getenv(LANG_ENVVAR), DEFAULT_LANG)
	marshalledFeeds, marshalError := json.Marshal(map[string]interface{}{
		"version": 1.0,
		"feeds":   feeds,
	})
	if marshalError != nil {
		fmt.Println("Couldn't marshal array of feeds")
		fmt.Println(marshalError)
		return marshalError
	}
	// The bundle inserter expects the "bundle" field to be a string,
	// so the JSON decoder used to decode the data below would return an error
	// if it tried to also treat the bundle as the JSON that it is and decode that too.
	bundleData, _ := json.Marshal(map[string]string{
		"url":     FeedsListIdentifier,
		"created": time.Now().Format(time.UnixDate),
		"bundle":  fmt.Sprintf("%s", string(marshalledFeeds)),
	})
	feedsInsertedStatus := InsertFreenet(bundleData)
	if feedsInsertedStatus == Success {
		// We don't want to write the data that was sent to the BI. Just the feeds stuff.
		feedWriteErr := ioutil.WriteFile(FeedsJsonFile, marshalledFeeds, os.ModePerm)
		if feedWriteErr != nil {
			fmt.Println("Couldn't write " + FeedsJsonFile)
			fmt.Println(feedWriteErr)
			return feedWriteErr
		}
	} else {
		fmt.Println("Failed to insert feeds list into Freenet")
		return errors.New(T("insertion_fail_err"))
	}
	return nil
}

/**
 * Write information about items received from a feed being followed.
 * Try to insert JSON containing this information into Freenet and, only if that succeeds,
 * write teh JSON to a file for distribution with the client.
 * @param feedUrl - The URL of the feed from which the items were received
 * @param items - Information about the items being followed
 */
func writeItems(feedUrl string, items []Item) error {
	T, _ := i18n.Tfunc(os.Getenv(LANG_ENVVAR), DEFAULT_LANG)
	marshalled, marshalErr := json.Marshal(map[string]interface{}{
		"version": 1.0,
		"items":   items,
	})
	if marshalErr != nil {
		fmt.Println("Couldn't marshal items for " + feedUrl)
		fmt.Println(marshalErr)
		return marshalErr
	}
	// The bundle inserter expects the "bundle" field to be a string,
	// so the JSON decoder used to decode the data below would return an error
	// if it tried to also treat the bundle as the JSON that it is and decode that too.
	bundleData, _ := json.Marshal(map[string]string{
		"url":     feedUrl,
		"created": time.Now().Format(time.UnixDate),
		"bundle":  fmt.Sprintf("%s", string(marshalled)),
	})
	insertStatus := InsertFreenet(bundleData)
	if insertStatus == Success {
		writeErr := writeItemsFile(feedUrl, marshalled)
		// We don't want to write the data that was sent to the bundle inserter, just the items stuff.
		if writeErr != nil {
			fmt.Println("Couldn't write item")
			fmt.Println(writeErr)
			return writeErr
		}
	} else {
		fmt.Println("Could not insert items into Freenet")
		return errors.New(T("insertion_fail_err"))
	}
	return nil
}

/**
 * TODO - Periodically delete items from the DB that we won't see again
 */

func main() {
	// Configure the i18n library to use the preferred language set in the CENOLANG environment variable
	setLanguage := os.Getenv("CENOLANG")
	if setLanguage == "" {
		os.Setenv("CENOLANG", "en-us")
		setLanguage = "en-us"
	}
	i18n.MustLoadTranslationFile("./translations/" + setLanguage + ".all.json")
	T, _ := i18n.Tfunc(setLanguage, "en-us")
	// Check that the configuration supplied has valid fields, or panic
	conf, err := ReadConfigFile(CONFIG_FILE)
	if err != nil {
		panic(T("no_config_rdr", map[string]interface{}{"Location": CONFIG_FILE}))
	} else if !ValidConfiguration(conf) {
		panic(T("invalid_config_rdr"))
	} else {
		Configuration = conf
	}
	// Establish a connection to the database
	var dbErr error
	DBConnection, dbErr = InitDBConnection(DB_FILENAME)
	if dbErr != nil {
		panic(T("database_init_error_rdr", map[string]interface{}{"Error": dbErr.Error()}))
	}
	// Set up the HTTP server to listen for requests for new feeds to read
	requestNewFollow := make(chan SaveFeedRequest)
	go followFeeds(requestNewFollow)
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/follow", followHandler(requestNewFollow))
	http.HandleFunc("/unfollow", unfollowHandler)
	http.HandleFunc("/insert", insertHandler)
	fmt.Println(T("listening_msg_rdr", map[string]interface{}{"Port": Configuration.PortNumber}))
	if err := http.ListenAndServe(Configuration.PortNumber, nil); err != nil {
		panic(err)
	}
}