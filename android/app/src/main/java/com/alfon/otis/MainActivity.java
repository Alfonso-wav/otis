package com.alfon.otis;

import android.os.Bundle;

import androidx.core.splashscreen.SplashScreen;

import com.getcapacitor.BridgeActivity;

public class MainActivity extends BridgeActivity {
    @Override
    public void onCreate(Bundle savedInstanceState) {
        SplashScreen.installSplashScreen(this);
        registerPlugin(GoServerPlugin.class);
        super.onCreate(savedInstanceState);
        // Allow background music (intro) to autoplay on splash load without a
        // user gesture. WebView defaults to requiring a gesture, which would
        // keep intro.mp3 silent until the user taps the Snorlax.
        this.bridge.getWebView().getSettings().setMediaPlaybackRequiresUserGesture(false);
    }
}
