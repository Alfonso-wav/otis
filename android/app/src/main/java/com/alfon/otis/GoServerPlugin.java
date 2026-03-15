package com.alfon.otis;

import android.util.Log;

import com.getcapacitor.Plugin;
import com.getcapacitor.annotation.CapacitorPlugin;

import mobile.Mobile;

/**
 * Capacitor plugin that starts/stops the Go HTTP server.
 * The server runs on localhost:8080 and serves the REST API
 * that the frontend communicates with via fetch().
 */
@CapacitorPlugin(name = "GoServer")
public class GoServerPlugin extends Plugin {

    private static final String TAG = "GoServerPlugin";
    private static final int PORT = 8080;
    private boolean running = false;

    @Override
    public void load() {
        startServer();
    }

    private void startServer() {
        if (running) return;
        try {
            String dataDir = getContext().getFilesDir().getAbsolutePath();
            Mobile.start(PORT, dataDir);
            running = true;
            Log.i(TAG, "Go server started on port " + PORT + " with dataDir=" + dataDir);
        } catch (Exception e) {
            Log.e(TAG, "Failed to start Go server", e);
        }
    }

    @Override
    public void handleOnDestroy() {
        stopServer();
    }

    private void stopServer() {
        if (!running) return;
        try {
            Mobile.stop();
            running = false;
            Log.i(TAG, "Go server stopped");
        } catch (Exception e) {
            Log.e(TAG, "Failed to stop Go server", e);
        }
    }
}
