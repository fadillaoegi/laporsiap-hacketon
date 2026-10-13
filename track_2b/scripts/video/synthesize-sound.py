"""Original stereo sound design, without samples or third-party music.

Soft tonal bed, filtered air transitions, and a D-major sonic signature.
Uses only Python's standard library. Outputs 48 kHz stereo PCM WAV.
"""
from array import array
import math
from pathlib import Path
import random
import sys
import wave

RATE = 48000
TAU = 2 * math.pi


def ease(x):
    x = max(0.0, min(1.0, x))
    return x * x * (3 - 2 * x)


def compose(mode, destination):
    duration = 9.0 if mode == "opening" else 8.0
    count = int(duration * RATE)
    left, right = array("d", [0.0]) * count, array("d", [0.0]) * count
    rng = random.Random(20261009 + (mode == "closing"))

    def add(i, value, pan):
        left[i] += value * math.sqrt((1 - pan) / 2)
        right[i] += value * math.sqrt((1 + pan) / 2)

    # A warm, restrained chord with slow stereo movement.
    for frequency, gain, phase in [(146.832, .045, 0), (220.0, .035, .8),
                                    (293.665, .025, 1.3), (369.994, .022, 2.2)]:
        for i in range(count):
            t = i / RATE
            env = ease(t / 1.5) * (1 - ease((t - (duration - 2.0)) / 1.8))
            wobble = 1 + .0012 * math.sin(TAU * .16 * t + phase)
            tone = math.sin(TAU * frequency * wobble * t + phase)
            tone += .12 * math.sin(TAU * frequency * 2 * t + phase)
            add(i, tone * gain * env, .26 * math.sin(.7 * t + phase))

    def whoosh(start, length, gain, pan):
        low = 0.0
        start_index = int(start * RATE)
        for j in range(int(length * RATE)):
            i = start_index + j
            if i >= count:
                break
            u = j / (length * RATE)
            low += (.035 + .12 * u) * (rng.uniform(-1, 1) - low)
            envelope = math.sin(math.pi * u) ** 2.5
            breath = low * .9 + .06 * math.sin(TAU * (180 + 70 * u) * j / RATE)
            add(i, gain * breath * envelope, pan * (1 - 2 * u))

    def bell(start, frequency, gain, pan, length=2.1):
        begin = int(start * RATE)
        for j in range(min(int(length * RATE), count - begin)):
            t = j / RATE
            attack = ease(t / .035)
            decay = math.exp(-t * 2.0) * (1 - ease((t - (length - .35)) / .35))
            # Clean fundamental and gentle glass harmonics.
            tone = math.sin(TAU * frequency * t)
            tone += .20 * math.sin(TAU * frequency * 2 * t + .1)
            tone += .06 * math.sin(TAU * frequency * 3 * t + .2)
            add(begin + j, gain * tone * attack * decay, pan)

    if mode == "opening":
        whoosh(.25, 1.20, .22, -.55)
        whoosh(1.35, 1.20, .20, .45)
        whoosh(2.80, 1.40, .16, -.35)
        for start, frequency, pan in [(1.65, 587.33, -.3), (2.00, 739.99, .0),
                                     (2.45, 880.0, .3), (4.25, 1174.66, .05)]:
            bell(start, frequency, .090 if start < 4 else .065, pan)
        bell(4.25, 293.665, .10, -.1, 2.8)
    else:
        whoosh(.10, 1.25, .19, -.4)
        whoosh(1.15, 1.20, .13, .35)
        for start, frequency, pan in [(1.25, 880.0, .25), (1.60, 739.99, -.1),
                                     (2.10, 587.33, .1), (2.75, 1174.66, .0)]:
            bell(start, frequency, .075, pan)
        bell(2.75, 293.665, .09, -.2, 3.0)

    peak = max(max(abs(x) for x in left), max(abs(x) for x in right))
    gain = .67 / peak
    pcm = array("h")
    for i in range(count):
        t = i / RATE
        edge = ease(t / .10) * (1 - ease((t - (duration - .65)) / .60))
        pcm.extend((round(left[i] * gain * edge * 32767),
                    round(right[i] * gain * edge * 32767)))
    if sys.byteorder != "little":
        pcm.byteswap()
    destination.parent.mkdir(parents=True, exist_ok=True)
    with wave.open(str(destination), "wb") as file:
        file.setnchannels(2)
        file.setsampwidth(2)
        file.setframerate(RATE)
        file.writeframes(pcm.tobytes())
    print(f"Created {destination}: {duration:.0f}s, stereo 48kHz; original synthesized sound")


if __name__ == "__main__":
    folder = Path(sys.argv[1] if len(sys.argv) > 1 else "tmp/video")
    for clip in ("opening", "closing"):
        compose(clip, folder / f"{clip}-sound-raw.wav")
