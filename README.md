# morse
**-- --- .-. ... .**  
Go package for string to Morse code converting.  

## Usage:
[DefafaultTranslator](https://github.com/please-the-turtle/morse/blob/master/defaultTranslator.go) translates string to morse phrase
```go
tr := morse.NewDefaultTranslator()
p := morse.Parse("Morse", tr)
```
[WavConverter](https://github.com/please-the-turtle/morse/blob/master/wave/wave.go) converts morse phrase to wav format
```go
tr := morse.NewDefaultTranslator()
p := morse.Parse("Morse", tr)
conv := wave.DefaultWavConverter()
wav, err := conv.Convert(p)
if err != nil {
	log.Println(err)
}

// Writing wave to a file
f, err := os.Create("morse_code.wav")
if err != nil {
	log.Println(err)
}
defer f.Close()
num, err := f.Write(wav)
if err != nil {
  fmt.Println(err)
}
```
