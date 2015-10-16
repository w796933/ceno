package plugins.CENO.Bridge.Signaling;

import java.io.IOException;
import java.net.MalformedURLException;
import java.util.concurrent.TimeUnit;

import org.bouncycastle.crypto.AsymmetricCipherKeyPair;
import org.bouncycastle.crypto.params.RSAKeyParameters;

import plugins.CENO.CENOErrCode;
import plugins.CENO.CENOException;
import plugins.CENO.Bridge.CENOBridge;
import freenet.client.FetchException;
import freenet.client.FetchResult;
import freenet.client.InsertException;
import freenet.client.async.BaseClientPutter;
import freenet.client.async.ClientContext;
import freenet.client.async.ClientPutCallback;
import freenet.keys.FreenetURI;
import freenet.node.RequestClient;
import freenet.support.Logger;
import freenet.support.SimpleFieldSet;
import freenet.support.api.Bucket;
import freenet.support.io.ResumeFailedException;

public class ChannelMaker {
	private String bridgeInsertURI;
	private AsymmetricCipherKeyPair asymKeyPair;

	static final long KSK_POLLING_PAUSE = TimeUnit.MINUTES.toMillis(4);

	public ChannelMaker(String insertURI, AsymmetricCipherKeyPair asymKeyPair) throws CENOException {
		this.bridgeInsertURI = insertURI;
		this.asymKeyPair = asymKeyPair;

		Puzzle puzzle = generatePuzzle();
		try {
			ChannelMakerAnnouncer channelAnnouncer = new ChannelMakerAnnouncer(bridgeInsertURI, (RSAKeyParameters) asymKeyPair.getPublic(), puzzle.getQuestion());
			channelAnnouncer.doAnnounce();
		} catch (IOException e) {
			throw new CENOException(CENOErrCode.RR, "IOException in channel announcer page creation");
		} catch (InsertException e) {
			throw new CENOException(CENOErrCode.RR, "Could not insert channel announcer page");
		}

		try {
			ChannelMakerListener channelListener = new ChannelMakerListener(puzzle.getAnswer());
			Thread listenerThread = new Thread(channelListener);
			listenerThread.start();
		} catch (MalformedURLException e) {
			throw new CENOException(CENOErrCode.RR, "Could not start Channel Maker Listener thread.");
		}
	}

	private Puzzle generatePuzzle() {
		Puzzle puzzle = new Puzzle();
		return puzzle;
	}

	private class ChannelMakerAnnouncer {
		private String bridgeInsertURI;
		private RSAKeyParameters pubAsymKey;
		private String puzzleQuestion;


		public ChannelMakerAnnouncer(String bridgeInsertURI, RSAKeyParameters pubAsymKey, String puzzleQuestion) {
			this.bridgeInsertURI = bridgeInsertURI;
			this.pubAsymKey = pubAsymKey;
			this.puzzleQuestion = puzzleQuestion;
		}

		public void doAnnounce() throws IOException, InsertException {
			SimpleFieldSet sfs = new SimpleFieldSet(false, true);
			sfs.putOverwrite("asymkey.modulus", pubAsymKey.getModulus().toString(32));
			sfs.putOverwrite("asymkey.pubexponent", pubAsymKey.getExponent().toString(32));
			sfs.putOverwrite("question", puzzleQuestion);

			FreenetURI insertURIconfig = new FreenetURI(bridgeInsertURI);
			FreenetURI announcementURI = new FreenetURI("USK", CENOBridge.announcerPath, insertURIconfig.getRoutingKey(), insertURIconfig.getCryptoKey(), insertURIconfig.getExtra());

			CENOBridge.nodeInterface.insertFreesite(announcementURI, "default.html", sfs.toOrderedString(), new AnnouncementInsertionCB());
		}
	}

	private class ChannelMakerListener implements Runnable {
		private String puzzleAnswer;
		private FreenetURI channelMakingKSK;

		private volatile boolean continueLoop;

		public ChannelMakerListener(String puzzleAnswer) throws MalformedURLException {
			this.puzzleAnswer = puzzleAnswer;
			channelMakingKSK = new FreenetURI("KSK@" + puzzleAnswer);
			continueLoop = true;
		}

		@Override
		public void run() {
			try {
				while(continueLoop) {
					// TODO Poll "KSK@puzzleAnswer" and discover insertion SSKs by clients
					FetchResult kskContent = null;
					try {
						kskContent = CENOBridge.nodeInterface.fetchURI(channelMakingKSK);
					} catch (FetchException e) {
						// TODO Fine-grain log messages according to FetchException codes
						Logger.warning(ChannelMakerListener.class, "Exception while fetching KSK clients use for making channels");
					}

					if (kskContent != null && kskContent.getMimeType() == "text/html") {
						Logger.normal(ChannelMakerListener.class, "A client has posted information for establishing a signaling channel");
						SimpleFieldSet sfs = new SimpleFieldSet(kskContent.toString(), false, true, true);
						ChannelManager.getInstance().addChannel(sfs);
					}
					// Pause the looping thread
					Thread.sleep(KSK_POLLING_PAUSE);
				}
			} catch (InterruptedException e) {
				continueLoop = false;
			} catch (IOException e) {
				// TODO Log this
			}
		}

	}

	class Puzzle {
		private String question;
		private String answer;

		public Puzzle() {
			this.question = Long.toHexString(Double.doubleToLongBits(Math.random()));
			this.answer = this.question;
		}

		public String getQuestion() {
			return question;
		}

		public String getAnswer() {
			return answer;
		}
	}

	class AnnouncementInsertionCB implements ClientPutCallback {
		
		private FreenetURI uri;

		public AnnouncementInsertionCB() {
		}

		@Override
		public void onResume(ClientContext context) throws ResumeFailedException {
		}

		@Override
		public RequestClient getRequestClient() {
			return CENOBridge.nodeInterface.getRequestClient();
		}

		@Override
		public void onGeneratedURI(FreenetURI uri, BaseClientPutter state) {
			this.uri = uri;
		}

		@Override
		public void onGeneratedMetadata(Bucket metadata, BaseClientPutter state) {
		}

		@Override
		public void onFetchable(BaseClientPutter state) {
		}

		@Override
		public void onSuccess(BaseClientPutter state) {
			Logger.normal(ChannelMaker.class, "Successfully inserted Channel Maker Announcer page with URI: " + uri);
		}

		@Override
		public void onFailure(InsertException e, BaseClientPutter state) {			
		}

	}
}
