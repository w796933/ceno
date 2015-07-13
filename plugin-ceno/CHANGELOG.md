# Change Log

## Unreleased
### Added  
  - LCS will try to fetch the latest version, using -1 as the suggested version in the ULPR request
  - Support for every type of URL
  - Automated the process of releasing the CENOBox

### Changed
  - USK calculation is using the base64 encoded version of the whole URL, instead of the domain
  - Updated LCS and RS error codes and default messages
  - RR base64 encodes the URL before sending it to BS


## [v0.3.0-rc1] - 2015-06-15
### Added
  - Improved CC-LCS integration and error handling
  - Support for requesting too big bundles
  - Improve FetchException handling at ULPR and local cache lookups

### Changed
  - Increased insertions priority to Interactive
  - Update BundleInserter table only after a bundle insertion for a URL has successfully started
  - LCS Base64 decodes the URL parameter from CC requests
  - Rename ulprStatus to ULPRStatus

### Fixed
  - Fix the values being inserted at BundleInserter's insertTable

### Removed
  - Remove content filtering from ULPRs



## [v0.3.0-rc0] - 2015-06-11
### Added
  - Improved CC-LCS integration and error handling
  - Improved URL validation
  - Asynchronous insertions of bundles using Manifests
  - Clear freemail messages being sent/received before unloading plugin
  - Use timer within ULPRManager before sending a freemail request
  - Automate the freemail accprops copying for CENO client identity
  - Handle Permanent Redirect exceptions at ULPR and local cache lookups

### Changed
  - Base64 encode of URLs for inter-agent message exchange

### Fixed
  - Issues with Freemails being discarded
  - Issue with inserting bundles for URLs that do not have an extraPath



## [v0.2.1] - 2015-07-15
### Added
  - Basic support for errors and polished integration of LCS with CC
  - Support for FCP messaging with the WOT plugin
  - Support for CENO specific error codes and exceptions

### Changed
  - Updated the WOT identities used for RS->RR message exchange

### Removed
  - Removed the lookup handler from the bridge



## [v0.2.0] - 2015-04-19
### Added
  - Request for bundles from the bridge over Freemail
  - Support HTML and JSON clients
  - Updated to be compliant with the new CENO communication protocol.
  - Created Client plugin that implements FredPluginHTTP
  - Functionality for local cache lookups
  - Added method for making lookup requests to the freenet node's local cache
  - Serve static templates to HTML clients

### Changed
  - Split client and bridge plugins
  - Use ULPRs (Ultra Light Passive Requests) for lookups in the distributed cache
  - Extracted URLtoUSK class



## [v0.1.0] - 2015-03-06
### Added
  - URL to USK translation
  - Insertion and lookup functionality
  - Configuration from a local file
  - Communication with the bridge using a RESTful API.



[Unreleased]: https://github.com/equalitie/ceno/compare/v0.3.0-rc1...HEAD
[v0.3.0-rc1]: https://github.com/equalitie/ceno/compare/v0.3.0-rc0...v0.3.0-rc1
[v0.3.0-rc0]: https://github.com/equalitie/ceno/compare/v0.2.1...v0.3.0-rc0
[v0.2.1]: https://github.com/equalitie/ceno/compare/v0.2.0...v0.2.1
[v0.2.0]: https://github.com/equalitie/ceno/compare/v0.1.0...v0.2.0
[v0.1.0]: https://github.com/equalitie/ceno/compare/48c7c207...v0.1.0