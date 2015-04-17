package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
)

const (
	SRC_IMG      = "images/numbers.png"
	CHARS_LENGTH = 12
	STARTX       = 0
	STARTY       = 0
	SPRITE_WIDTH = 300
)

type (
	Chars   map[int]draw.Image
	counter map[string][]int
)

func (c counter) Inc(id string) {
	if c[id] == nil {
		c[id] = make([]int, 3)
	}

	if c[id][2] == 9 {
		if c[id][1] == 9 {
			c[id][0] += 1
			c[id][1] = 0
		} else {
			c[id][1] += 1
		}
		c[id][2] = 0
	} else {
		c[id][2] += 1
	}
}

var (
	counts counter = make(counter)
	chars  Chars   = make(Chars)
)

func handleError(err error) {
	if err != nil {
		fmt.Println("ERROR:", err)
	}
}

func main() {
	// Read in source image
	spritesFile, err := os.Open(SRC_IMG)
	handleError(err)
	// decode image.
	sprites, err := png.Decode(spritesFile)
	handleError(err)

	var x int = STARTX
	var y int = STARTY

	for i := 0; i < CHARS_LENGTH; i++ {
		chars[i] = image.NewRGBA(image.Rect(0, 0, 100, 100))
		draw.Draw(
			chars[i],
			image.Rect(0, 0, 100, 100),
			sprites,
			image.Pt(x, y),
			draw.Src,
		)

		if x += 100; x >= SPRITE_WIDTH {
			x = 0
			y += 100
		}
	}

	mux := http.NewServeMux()
	mux.Handle("/counter/", counterHandler())

	log.Fatal(http.ListenAndServe(":8080", mux))
}

func counterHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		re := regexp.MustCompile("/counter/(\\w*)")
		matches := re.FindStringSubmatch(r.URL.Path)
		if len(matches) == 0 {
			fmt.Println("Error: not enough matching")
			return
		}

		var identifier string
		if len(matches) > 1 {
			identifier = matches[1]
		}
		switch r.Method {
		case "GET":
			counts.Inc(identifier)
			current := counts[identifier]
			theImage := image.NewRGBA(image.Rect(0, 0, 300, 100))
			for i, n := range current {
				r := image.Rect(100*i, 0, 100*i+100, 100)
				char := chars[SpriteIndex(n)]
				draw.DrawMask(
					theImage,
					r,
					&image.Uniform{
						color.RGBA{uint8(rand.Intn(255)), uint8(rand.Intn(255)), uint8(rand.Intn(255)), 255},
					},
					image.ZP,
					char,
					image.ZP,
					draw.Src,
				)
				// draw.Draw(theImage, r, char, image.ZP, draw.Src)
			}
			png.Encode(w, theImage)
		case "DELETE":
			counts[identifier] = nil
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
}

func SpriteIndex(n int) int {
	switch {
	case n == 0:
		return 10
	case n <= 9:
		return n - 1
	default:
		return -500
	}
}
