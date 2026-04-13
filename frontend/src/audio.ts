// Background music module. Two-audio pool for crossfade between tab tracks.
// Volume + mute persist in localStorage and apply live.

export type TrackId = "intro" | "pokedex" | "explore" | "simulations";

const SRC: Record<TrackId, string> = {
  intro: "/audio/intro.mp3",
  pokedex: "/audio/pokedex.mp3",
  explore: "/audio/explore.mp3",
  simulations: "/audio/simulations.mp3",
};

const FADE_MS = 600;
const VOLUME_KEY = "music-volume";
const MUTED_KEY = "music-muted";

let currentAudio: HTMLAudioElement | null = null;
let currentTrack: TrackId | null = null;
let userVolume = readVolume();
let muted = readMuted();

function clamp(n: number, lo: number, hi: number): number {
  return Math.max(lo, Math.min(hi, n));
}

function readVolume(): number {
  const raw = localStorage.getItem(VOLUME_KEY);
  if (raw === null) return 0.5;
  const n = parseFloat(raw);
  return Number.isFinite(n) ? clamp(n, 0, 1) : 0.5;
}

function readMuted(): boolean {
  return localStorage.getItem(MUTED_KEY) === "true";
}

function persistVolume(): void {
  localStorage.setItem(VOLUME_KEY, String(userVolume));
}

function persistMuted(): void {
  localStorage.setItem(MUTED_KEY, String(muted));
}

function fadeTo(el: HTMLAudioElement, target: number, ms: number, onDone?: () => void): void {
  const start = el.volume;
  const t0 = performance.now();
  function step(now: number): void {
    const p = Math.min(1, (now - t0) / ms);
    el.volume = clamp(start + (target - start) * p, 0, 1);
    if (p < 1) requestAnimationFrame(step);
    else onDone?.();
  }
  requestAnimationFrame(step);
}

function fadeOutAndStop(el: HTMLAudioElement, ms: number): void {
  fadeTo(el, 0, ms, () => {
    el.pause();
    el.src = "";
  });
}

/** Play given track with crossfade. Safe to call on user gesture. */
export async function playTrack(id: TrackId): Promise<void> {
  if (currentTrack === id && currentAudio && !currentAudio.paused) return;

  const next = new Audio(SRC[id]);
  next.loop = true;
  next.volume = 0;

  try {
    await next.play();
  } catch {
    // Autoplay blocked — bail without mutating state so caller can retry later.
    return;
  }

  const targetVol = muted ? 0 : userVolume;
  fadeTo(next, targetVol, FADE_MS);

  if (currentAudio) fadeOutAndStop(currentAudio, FADE_MS);

  currentAudio = next;
  currentTrack = id;
}

export function stop(): void {
  if (currentAudio) {
    fadeOutAndStop(currentAudio, FADE_MS);
    currentAudio = null;
    currentTrack = null;
  }
}

export function setVolume(v: number): void {
  userVolume = clamp(v, 0, 1);
  if (currentAudio && !muted) currentAudio.volume = userVolume;
  persistVolume();
}

export function setMuted(b: boolean): void {
  muted = b;
  if (currentAudio) currentAudio.volume = b ? 0 : userVolume;
  persistMuted();
}

export function getVolume(): number {
  return userVolume;
}

export function getMuted(): boolean {
  return muted;
}

export function isTrackPlaying(id: TrackId): boolean {
  return currentTrack === id && currentAudio !== null && !currentAudio.paused;
}
