import type { CapacitorConfig } from "@capacitor/cli";

const config: CapacitorConfig = {
  appId: "com.alfon.otis",
  appName: "Otis Pokédex",
  webDir: "dist",
  android: {
    path: "../android",
  },
  server: {
    androidScheme: "http",
  },
};

export default config;
