# Censorship.NO! Threat Model

## Contents

  1. [Application Description](#1-application-description)
  2. [Security Objectives](#2-security-objectives)
  3. [Assumptions](#3-assumptions)
  4. [Threats](#4-threats)
  5. [Security Audit Results](#5-security-audit-results-as-of-ceno-v051)


## 1. Application Description

CENO (Censorship.NO!) is an innovative approach to censorship circumvention,
based on P2P storage networks. Users do not need to know a friendly proxy server
in the uncensored zone to bypass local filtering. CENO maintains strong privacy
and anonymity features as well as offering users plausible deniability in an
emergency situation. CENO is built in advance of aggressive Internet filtering
and the establishment of national intranets to fence off citizens from the
wicked Web.

Main objective of the CENO project is to deliver in a peer-to-peer fashion
static content that otherwise would not be available because of Internet
censorship. It is therefore a great tool for reading the news – a selection of
news feeds has already been inserted and getting updated on a daily basis,
accessible via the "CENO Portal". Furthermore, users can anonymously request
specific nodes that have access to the uncensored Web (from now on "Bridges") to
fetch a URL, prepare a bundle for it and insert it in the distributed storage,
so as to be retrieved by the user that requested it. Future requests for the
same URL will be handled directly via the cache without the need for user-Bridge
communication, even if the website is no more live. No knowledge of the global
network topology is required in order to retrieve bundles or send a message to a
Bridge. In cases of nationwide Internet throttling, content will remain
available to the peers given that a copy of that bundle is cached in the
in-country network of peers. When using CENO, you will not be able to visit
dynamic websites, log into websites, stream media or read and send emails.


<!-- ## 2. Architecture Overview

#### Insertion Authorities

##### Master Bridges
##### RSS-insertion Bridges
##### Signal Bridges

##### Bridge Agents
###### Bundle Server
###### Bundle Inserter
###### Request Receiver
###### RSS-Reader

##### Backbone Nodes

### CENO Users

##### Client Agents
###### CENO Router Browser extensions
###### CENO Client Proxy
###### Local Cache Server
###### Request Sender -->


## 2. Security Objectives

#### Anonymity

CENO promises strong user anonymity guarantees, inherited by using Freenet as
the underlying storage and communication medium. Nodes that are using CENO are
indistinguishable from the rest of the Freenet nodes and form part of the same
global network of peers. Our user-Bridge signaling mechanism leaks no metadata
(apart from the intention of a node to establish a secure channel with the
Bridge), and the bridge nodes cannot know who is requesting a URL.

It is known that sophisticated attacks to networks such as Freenet could expose
the identity of users, as we are describing in the "Threats" section.


#### A secure communication mechanism

CENO users establish a signaling channel with their Bridge of preference. The
mechanism in operation provides the following features:

  * Confidentiality
  * Integrity
  * Causality Preservation
  * Sender and, once the channel has been established, Recipient Anonymity
  * No shared secrets needed
  * No service provider required
  * Asynchronicity
  * Established channels spam-resistance

Users can always drop a channel and establish a new one. Spammers cannot
interfere with already established channels, but can only spam their own
channel. Spamming the channel establishment mechanism might temporarily prevent
new users from establishing a new channel with a Bridge, but will not result in
Denial of Service of the rest of CENO functionality. Notably, the design of our
signaling mechanism requires no prior authentication or manual interaction by
the users, is significantly efficient compared to similar ones, and can scale
horizontally so as to handle increasing demand.


#### Transport Privacy

CENO traffic among peers is end-to-end encrypted and appears on traffic analysis
as random UDP noise. The ephemeral encryption keys are kept safely by the peers
and no non-global entities can discover which sets of nodes are communicating
with each other. Global adversaries will be able to tell with whom a certain
node is exchanging encrypted messages, but will not be able to find out the
actual content. In few words, even though adversaries (such as your Internet
Service Provider) might see that you are using a censorship circumvention tool,
it will be extremely difficult for them to know what you are using it for.


#### Plausible deniability

When using CENO, you are contributing to the whole of Freenet by sharing part of
your hard drive and by storing encrypted chunks of files in it. This is how
content remains available, once inserted in the network. The majority of this
content is not related to CENO, and decryption keys are not included in these
chunks. Unlike torrents, you cannot know or control what is stored in your
machine, and you don't need to directly connect with nodes that host a specific
Freenet file. Therefore, it would be impossible for users to know what kind of
parts of files they are sharing, given the immense number of files that have
been inserted in the distributed cache.


#### Resistance against active network interference

CENO will not establish any connections with servers or proxies, but will use
the Freenet's ones with other peers in order to retrieve content and forward
requests for URLs to a CENO Insertion Authority. Freenet traffic (encrypted
end-to-end and over UDP) looks like random noise to Deep Packet Inspection, so
it is difficult for adversaries to create specific rules for dropping Freenet
connections without affecting other services. It is worth mentioning that
adversaries that control the network transport link might be able to block
connections to the Freenet seed nodes. This could imply that certain nodes will
not be able to find other peers in order to request content, unless the users
know other people who are already using CENO/Freenet and can help them become
part of the global network.

In scenarios where connections with peers in other countries are throttled (e.g.
because of a national firewall), but those within the country are left intact,
CENO users will still be able to communicate with each other and, given that one
of them has cached a specific resource, she could serve it upon request to the
rest of the in-country network and pass it along as her node is replicating its
cached content. This concept can also be extended to mesh networks that are
using independent network infrastructure. The CENO team imagines that in such
situations, when a small set of nodes can connect with the global network, CENO
users will still be able to access the portal content and request new bundles.


## 3. Assumptions

#### User Behavior

  * Users do not manually save on their hard drives content from CENO Portal
    bundles or other files retrieved from Freenet.
  * Users always browse CENO via the browser (Firefox or Chrome) window that
    opens when they start CENO and that is using the customized profile and
    add-ons.
  * Users do not: alter the settings of the preconfigured browser profile;
    disable the CENO add-on; or enable other add-ons that could interfere with
    CENO.
  * Users do not make requests for privacy sensitive URLs, for example URLs that
    include a username in the GET parameters.
  * Users do not manually change the Insertion Authority they want to use
    and select a malicious one.
  * Users do not add as friends nodes operated by entities they do not know or
    do not trust.
  * Users do not manually degrade the default Freenet/CENOBox security options.
  * Users do not manually set Freenet node logging level to a level that retains
    private information, such as the URLs they have been requesting.
  * Users do not run software that could interact with Freenet (via the Freenet
    Client Protocol or the Web interface).


#### Freenet Nodes

  * Freenet nodes do not drop requests for lookups or insertions.
  * Freenet peers respond honestly when they are asked whether they store a
    specific chunk.
  * Freenet nodes do not log requests of other peers.
  * Freenet peers do not probe their neighbors for specific chunks.
  * Freenet seed nodes are operating and reachable.
  * No adversary controls a large part of the Freenet network.
  * No malevolent nodes flood their neighbors' data stores.

We recommend you to refer to the [Freenet Threat
model](https://wiki.freenetproject.org/Threat_model) in its project wiki (also
[Major attacks](https://wiki.freenetproject.org/Major_attacks)
and
[Potential attackers](https://wiki.freenetproject.org/Potential_attackers)).
We shall not consider Freenet vulnerabilities (such as Sybil attacks) as part of
our threat model from now on.

At any rate, CENO protocol has been designed in such a way that future CENO
versions can be easily ported to use other storage or transport networks, by
replacing the related agents.


#### Insertion Authority/Bridge Maintainers

  * The Bundle Server is configured to route all requests via an anonymization
    network.
  * Bridges do not trust compromised SSL Certificates.
  * IA Maintainers do not run software that could interfere with CENO agents, or
    Freenet, or expose that they are hosting a Bridge.


#### General Assumptions

  * Users' and Insertion Authority Maintainers' systems are not compromised or
    affected by malware that could interfere with CENO.
  * Users and Bridges are given enough bandwidth and storage resources to
    effectively participate in the network.
  * No security vulnerabilities exist in the programming languages, frameworks
    or libraries in use.
  * Cryptography works and there are no operating quantum computers.


## 4. Threats

#### i. Malicious Insertion Authorities

CENO users have access via the CENO Portal and CENO Browser window only to
bundles and feeds that have been inserted by the Insertion Authority (IA) they
trust.  CENOBox comes preconfigured with an Insertion Authority the CENO team is
maintaining, but users can manually decide to use another Bridge, both for
reading the portal news as well as for requesting new URLs, by setting an
alternative valid Insertion Authority bridgeKey in the `.CENO/client.properties`
file.

When selecting to retrieve bundles from independent Insertion Authorities, users
could be exposed to tampered or potentially "harmful" content. Malicious IAs
could censor content, alter it, deny to serve a resource (i.e. to insert a
bundle for a URL), or not serve specific established signaling channels.
Uncloaking users by injecting code in the bundles can be considered impossible,
since CENO Local Cache Server strips any scripts from bundles before serving
them to the users. In the worst case scenario, a malicious AI controlling a
significant amount of nodes in the Freenet network could try to combine the URLs
requested by a channel with a de-anonymization attack, therefore correlating the
requested URLs with the IP that is making insertions in an established channel.

This is a medium-risk issue with low exploitability, that affects only users
who have manually selected to use a non-default IA.


#### ii. CENO Insertion Authority private key compromisation

For content to be discoverable and authenticatable, Insertion Authorities
use a private SSK – [Signed Subspace Key](https://wiki.freenetproject.org/SSK)
– for inserting bundles, portal feeds, information for establishing secure
signaling channels, etc. This key is included in the `.CENO/bridge.properties`
file in Signal and Master/RSS-Insertion Bridges.

If an adversary gets access to that private key, she will be able to insert
content that would be discoverable with the Insertion Authority's public
bridgeKey, without CENO users being able to tell the difference. Furthermore,
if the adversary inserts newer editions of pre-inserted bundles, lookups from
users in the distributed cache will return the latest, tampered content.
Adversaries could then behave as described in the "Malicious Insertion
Authorities" threat, but also invest on the technique and cause Denial of
Service of the Insertion Authority.

We categorize this issue among the high-risk ones, with low exploitability. At
the moment there is no implemented mechanism for recovering from such an
incident. The CENO team is planning to work on a status page that users' clients
will automatically poll to check whether an Insertion Authority has been
compromised. That page will be inserted by the Insertion Authority owner under
the same SSK Freenet key.

<!--- ##### ix. Control over Insertion Authority's key area

//Investigate whether SSKs are all stored in the same area -->


#### iii. Established signaling channels compromisation

Established channels are stored on Signal Bridges at a database in the same
directory as the Freenet files and CENO agents. Obviously there is no mapping
between a channel and an IP, but if an adversary gains access to a channel's
database file, she will be able to get a history of the URLs requested by each
user. Given that the list of URLs include user-specific attributes, or can be
matched with a targeted user's behavior, the adversary could extract
information regarding the way a specific individual has been using CENO.

We consider this threat to be characterized by high risk, but by low
exploitability.


#### iv. CENO agents and Freenet node exposure via requests to their RESTful API

CENO Agents are built in a micro-services fashion and are communicating with
each other over an HTTP RESTful API. It is worth mentioning that agents accept
requests only from localhost connections, yet there is no provision for blocking
requests from non-CENO services running on the users' or IA's machines. Email
clients that parse HTML and Javascript emails, as well as browsers or software
specifically developed for this reason, could identify and fingerprint a running
CENO agent. In such a case, out-of-band mechanisms, for example requests to a
remote machine, could easily leak the IP of a CENO user or bridge node
maintainer. In addition, malicious software could change the configuration
options of the Freenet node, or actively interfere in other ways with the CENO
agents so as to jeopardize the anonymity of users.

This is a high-risk and easily exploitable issue the CENO developers aim towards
mitigating as soon as possible. In the meantime, we strongly recommend that IA
maintainers use their bridge machines only for that purpose and avoid installing
software that is not related to CENO, and that CENO users avoid running any
other programs while using CENO. If supported by the operating system, CENO
should be run by a separate user.


#### v. CENO Signaling channel establishment puzzle slots flooding

In order to establish a secure signaling channel, Signal Bridges publish their
public RSA key along with a quiz whose solution leads to a set of keys, from now
on slots, that everybody can use to insert a page in Freenet. Clients generate a
new channel, encrypt it so that it is readable only by the Signal Bridge and
insert it into one of those slots. Eventually the Signal Bridge will poll the
slots, discover and decrypt the request of a client and start accepting requests
on that channel.

However, in distributed storage networks such as Freenet, slots can be used only
once and, once consumed, Signal Bridges will have to publish a new puzzle. A
malicious adversary could keep inserting nonce and consuming the slots as soon
as published, therefore causing a Denial of Service on the channel establishment
process. Previously established channels will not be affected by the slots
flooding attack.

This is a medium-risk vulnerability (it does not affect the security of users or
Insertion Authorities), but an easily exploitable one. Future CENO versions will
require users to solve a CAPTCHA the first time they are communicating with
an Insertion Authority, in order to stop spambots.


#### vi. CENO User machine compromisation / physical seizure

Given that a user has not manually downloaded content from Freenet or saved CENO
bundles, all chunks stored in the user's hard disk are encrypted. Even though it
will be obvious that the owner of the machine had been using CENO, the adversary
will need to have the decryption keys in order to prove that the machine owner
was hosting particular content (see Security Objective "Plausible deniability").

Compromisation of the `.CENO/client.properties` configuration file will leak the
private signaling channel the machine owner had established with an Insertion
Authority. The adversary could then use the signaling channel key to look up
previous requests by the owner and thus find out what URLs she had been
visiting. In addition, the exposure of the Insertion Authority bridgeKey would
indicate what kind of news sources the machine owner could have been reading.

If the adversary gains access to a CENO node in operation, she will also be
able to discover the Freenet files (including CENO bundles/URLs) the machine
owner had been accessing in that session.

We consider physical seizure of CENO client nodes a high-impact, low-risk issue.


#### vii. Freenet traffic throttling

Alarmed by the ever-increasing capacity of global adversaries to actively
perform network interference, we are considering Freenet traffic throttling a
possible scenario. It is known that Freenet users within China have been having
problems with connecting with nodes outside of the country.

Efforts have been made towards adding support for obfuscating Freenet nodes
traffic, but this is an ongoing process and needs the support of the pluggable
transport implementors groups. At the moment, there does not exist a project
that has implemented UDP traffic obfuscation.

This is a high-risk issue which is exploitable only by adversaries who have
access to the network infrastructure, such as ISPs. It is important to mention
that usually adversaries throttle only outgoing connections and do not mess with
in-country node-to-node traffic. In such cases, CENO will still be partially
functional and can serve as a tool for sharing information among peers (see
Security Objective "Resistance to active network interference").


##### viii. Insertion Authority Network Interference

Bundle Server instances running on Insertion Authority Bridge and RSS-Inserter
nodes need to fetch content from the uncensored Web. Such requests are not
protected by CENO and hence could expose the IP addresses of those nodes. In
order to exploit that, a malicious entity could request from an IA to create a
bundle for a website that she owns, logging the IPs of the machines that are
accessing it in the background. Once the IP address of an IA node is revealed,
the adversary could perform other type of attacks in order to bring down or gain
access to that node. Even worse, a global adversary could tamper with the
network traffic of de-anonymized IA nodes and insert altered content in CENO.

Insertion Authority maintainers are strongly advised to route all Bundle Server
traffic via an anonymization network, such as Tor or I2P. In addition, we
recommend the use of [HTTPS-Proxy](https://github.com/equalitie/HTTPS-Proxy)
that integrates the
[HTTPS-Everywhere](https://github.com/EFForg/https-everywhere) ruleset, in order
to upgrade HTTP requests to HTTPS ones whenever possible.

This is a medium-risk issue with low exploitability.


## 5. Security Audit results (as of CENO v0.5.1)

CENO v0.5.1 has been reviewed by NCC Group and results have been published [here].
The agents that were audited are:

  * CENO Client Proxy
  * CENO Bundle Server
  * CENO Freenet plugins (`CENO.jar` and `CENOBridge.jar`)
  * CENO Browser Extensions

##### NCC-CENO-006: Bundler Server Does Not Validate TLS Certificate
This issue was **fixed** at CENO v0.6-rc, by explicitly setting `strictSSL` to
true, as suggested by the auditors group.

##### NCC-CENO-015: Client Downgrades HTTPS Connection Attempts to HTTP
In order to mitigate this issue, the CENO team has decided to imply the
HTTPS-Everywhere ruleset to all outgoing Bundle Server requests. We consider
this "best-effort" upgrade to HTTPS the best thing to do, given that it should
not be the first user who decides whether the HTTP over the HTTPS version is
inserted in the distributed storage, and in order to ensure the genuinity of the
content inserted. This is **work in progress**.

##### NCC-CENO-017: CENO Request Scheme Obfuscates URLs and Flattens Web Origins
The CENO team believes that by dropping support for Javascript and by providing
a hardened browser profile for Firefox, this issue has been **mitigated**.

##### NCC-CENO-018: JavaScript Can Unmask Users Via Multiple Mechanisms
Similarly to the previous issue identified by the NCC Group, by applying
readability mode on the bundles inserted by Bridges and by stripping any
Javascript at the Lookup Cache Servers, we consider this issue **resolved**.

The proposed solution of inserting assets individually rather than as part of
the bundle is something that could be considered as a future improvement.

##### NCC-CENO-019: Exposed Freenet Plugin Enable XSS In Freenet Origin and User Identification
This issue, overlapping with threat iv. identified above, is still a **critical
vulnerability**.

##### NCC-CENO-021: CENO Client Daemon Is Exposed and Does Not Act as an HTTP Proxy
This issue remains a **critical vulnerability** and will be addressed by using
application secrets, in the same way as with NCC-CENO-019.

##### NCC-CENO-013: Bridge Server Cannot Recover From Compromise
This issue remains a **vulnerability**. The CENO team describes their plans to
mitigate it in threat ii. above.

##### NCC-CENO-002: CENOBox Configures Freenet To Run In Opennet Mode
Unfortunately there is no alternative way for shipping CENOBox to new Freenet
users. In the documentation it is demonstrated how one can add people she
personally knows and trusts as Freenet friends and eventually switch to
"Darknet" mode. Users are encouraged to do so in order to better protect
themselves from de-anonymization attacks. CENO team **will not address** this
issue.

##### NCC-CENO-001: Distribution of Untrusted and Unsigned Mozilla Firefox Add-On
CENO team is now distributing a signed version of their Mozilla Firefox add-on
included in the CENOBox, along with a profile that automatically loads it.
Therefore, the issue is considered **resolved**. However, the team has decided
not to upload the add-on to the addons.mozilla.org site (AMO), since it could
interfere with their normal browsing experience, for example by downgrading
HTTPS to HTTP requests.

##### NCC-CENO-003: CENO Bridge Uses Static Source Port
Bridge maintainers are encouraged to proxy Bundler Server traffic via an
anonymization network (specifically Tor), in which case the static source port
will not be exposed. We consider this issue **mitigated**.

##### NCC-CENO-005: Extension Toggle Switch Implementation is Unsafe
CENO team has removed the Toggle Switch, therefore this issue is considered
**resolved**.

##### NCC-CENO-008: Bridge Server Private Files Accessible To Other Users
CENOBridge distributable explicitly sets `bridge.properties` file permissions and
makes it unreadable by other users. In addition, Bridge maintainers are
encouraged to avoid running non-related services on a bridge machine and are
instructed to check the access permissions of such sensitive files. We consider
this issue **resolved**.

##### NCC-CENO-009: Firefox Add-Ons Can Modify CENO Add-On
The CENO team has followed the auditor's suggestion and now ships a
pre-configured profile that loads the Firefox CENO add-on, while blocking access
to the Mozilla Add-On site. Therefore the issue is **resolved**.

##### NCC-CENO-007: Sensitive Information Stored In Bridge Server logs
The CENO team has followed the auditor's suggestions and we consider this issue
**resolved**.

##### NCC-CENO-014: CENO Freemail Account Accessible To All Users
This issue is **not relevant** now that CENO is using the decentralized
signaling mechanism.

##### NCC-CENO-010: Hard-Coded Bridge in Application Source Code
Users can now configure their nodes to use alternative Insertion Authorities.
However, there is no mechanism for announcing compromised bridge nodes yet.
Therefore this issue is **partially mitigated**.

##### NCC-CENO-011: Missing Recommendations For Bridge Server Hardening
The CENO team has composed an instruction hardening guide for bridge nodes
maintainers which can be found in [doc/secureBridge.md](secureBridge.md) and
[ceno-freenet/INSTALL.Bridge.md](../ceno-freenet/INSTALL.Bridge.md), therefore
the issue is **resolved**.

##### NCC-CENO-016: USK URL Generation Does Not Escape Base64-Encoded URL Data
All CENO agents are now using URL-safe base64 encoding (rfc4648-sec5), therefore
the issue is **resolved**.


<hr>
Threat Model version 1.0  
CENO version 1.0.0  
ceno@equalit.ie
