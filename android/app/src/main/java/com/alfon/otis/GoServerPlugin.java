package com.alfon.otis;

import android.util.Log;

import com.getcapacitor.Plugin;
import com.getcapacitor.annotation.CapacitorPlugin;

import java.net.HttpURLConnection;
import java.net.URL;

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
    private static final int HEALTH_TIMEOUT_MS = 10000;
    private static final int HEALTH_INTERVAL_MS = 200;
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
            waitForServer();
        } catch (Exception e) {
            Log.e(TAG, "Failed to start Go server", e);
        }
    }

    private void waitForServer() {
        long deadline = System.currentTimeMillis() + HEALTH_TIMEOUT_MS;
        String healthUrl = "http://localhost:" + PORT + "/api/pokemon?offset=0&limit=1";
        while (System.currentTimeMillis() < deadline) {
            try {
                HttpURLConnection conn = (HttpURLConnection) new URL(healthUrl).openConnection();
                conn.setConnectTimeout(500);
                conn.setReadTimeout(500);
                conn.setRequestMethod("GET");
                int code = conn.getResponseCode();
                conn.disconnect();
                if (code == 200) {
                    Log.i(TAG, "Go server is ready (health check passed)");
                    return;
                }
            } catch (Exception ignored) {
                // Server not ready yet
            }
            try {
                Thread.sleep(HEALTH_INTERVAL_MS);
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
                return;
            }
        }
        Log.w(TAG, "Go server health check timed out after " + HEALTH_TIMEOUT_MS + "ms");
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
