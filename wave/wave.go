package wave

import (
	"errors"
	"math"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/wav"
	"github.com/orcaman/writerseeker"
	"github.com/please-the-turtle/morse"
)

// WavConverter converts morse.Phrase to wav format.
type WavConverter struct {
	// Frequency of output wav.
	Frequency float64
	// Length of the dot sound.
	DotLen time.Duration
}

// Returns WavConverter with default configuratiions.
func DefaultWavConverter() *WavConverter {
	return &WavConverter{
		Frequency: 800.0,
		DotLen:    80 * time.Millisecond,
	}
}

// Converts Morse phrase to wav.
func (c *WavConverter) Convert(p morse.Phrase) ([]byte, error) {
	s := sineWave{
		sampleRate: 44100,
		freq:       c.Frequency,
		phase:      0,
	}

	format := beep.Format{
		SampleRate:  s.sampleRate,
		NumChannels: 2,
		Precision:   2,
	}

	w := &writerseeker.WriterSeeker{}

	buf := beep.NewBuffer(format)
	for _, r := range p {
		switch r {
		case '.':
			buf.Append(c.dot(&s))
		case '-':
			buf.Append(c.dash(&s))
		case '/':
			buf.Append(c.wordGap(&s))
		case ' ':
			buf.Append(c.letterGap(&s))
		}
		buf.Append(c.gap(&s))
	}

	err := wav.Encode(w, buf.Streamer(0, buf.Len()), format)
	if err != nil {
		err = errors.New("Morse to wav encoding: " + err.Error())
		return nil, err
	}

	r := w.BytesReader()
	result := make([]byte, r.Size())
	_, err = r.Read(result)
	if err != nil {
		err = errors.New("Morse to wav encoding: " + err.Error())
		return nil, err
	}

	return result, nil
}

func (w *sineWave) silence(d time.Duration) beep.Streamer {
	return beep.Silence(w.sampleRate.N(d))
}

func (c *WavConverter) dot(s *sineWave) beep.Streamer {
	return s.generate(c.DotLen)
}

func (c *WavConverter) dash(s *sineWave) beep.Streamer {
	return s.generate(3 * c.DotLen)
}

func (c *WavConverter) gap(s *sineWave) beep.Streamer {
	return s.silence(c.DotLen)
}

func (c *WavConverter) letterGap(s *sineWave) beep.Streamer {
	return s.silence(3 * c.DotLen)
}

func (c *WavConverter) wordGap(s *sineWave) beep.Streamer {
	return s.silence(c.DotLen)
}

type sineWave struct {
	sampleRate beep.SampleRate
	freq       float64
	phase      float64
}

func (w *sineWave) generate(d time.Duration) beep.Streamer {
	dt := w.freq / float64(w.sampleRate)
	num := w.sampleRate.N(d)

	return beep.StreamerFunc(
		func(samples [][2]float64) (n int, ok bool) {
			if num == 0 {
				return 0, false
			}

			if 0 < num && num < len(samples) {
				samples = samples[:num]
			}

			for i := range samples {
				v := math.Sin(w.phase * 2.0 * math.Pi) // period of the wave is thus defined as: 2 * PI.
				samples[i] = [2]float64{v, v}
				_, w.phase = math.Modf(w.phase + dt)
			}

			if num > 0 {
				num -= len(samples)
			}

			return len(samples), true
		})
}
