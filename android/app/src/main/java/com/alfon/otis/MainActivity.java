package com.alfon.otis;

import android.os.Bundle;

import com.getcapacitor.BridgeActivity;

public class MainActivity extends BridgeActivity {
    @Override
    public void onCreate(Bundle savedInstanceState) {
        registerPlugin(GoServerPlugin.class);
        super.onCreate(savedInstanceState);
    }
}
