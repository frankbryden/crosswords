package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type WordList struct {
	words []string
}

type Move struct {
	x, y int
	d rune
}

func (wl WordList) MaxCharSurface() int{
	sum := 0
	for _, word := range wl.words {
		sum += len(word)
	}
	return sum
}

func (wl WordList) Copy() *WordList {
	newWl := make([]string, len(wl.words))
	copy(newWl, wl.words)
	return &WordList{newWl}
}

func (wl WordList) RemoveWord(i int) *WordList {
	newWl := make([]string, len(wl.words))
	copy(newWl, wl.words)
	newWl[i] = newWl[len(newWl) - 1]
	newWl = newWl[:len(newWl) - 1]
	return &WordList{newWl}
}


func (wl WordList) DetermineMinCrosswordSize() int {
	requiredSurface := wl.MaxCharSurface() //we need to make a square grid where the surface is larger than this val
	size := 0
	for (size*size) < requiredSurface {
		size++
	}
	return size
}

type Crossword struct {
	data [][]rune //represents the actual crossword - a grid of letters, with 0 meaning empty, and -1 for black block
	wl *WordList //list of words remaining to fill it
}

func (cw Crossword) isEmptyRow(i int) bool {
	for _, d := range cw.data[i] {
		if d != 0 {
			return false
		}
	}
	return true
}

func (cw Crossword) isEmptyCol(c int) bool {
	for i := 0; i < len(cw.data[0]); i++ {
		fmt.Println(c, i)
		if cw.data[i][c] != 0 {
			fmt.Println(cw.data[c], "is not empty")
			return false
		}
	}
	fmt.Println(cw.data[c], "is empty")
	return true
}

func inflateWordList(filename string) *WordList {
	dat, e := ioutil.ReadFile(filename)
	
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
	
	lines := strings.Split(string(dat), "\n")
	lines = lines[:len(lines) - 1]
	words := make([]string, len(lines))
	
	for i, word := range lines {
		words[i] = word
	}
	
	return &WordList{words}
}

func (cw Crossword) CanPlaceWord(word string, i, x, y int, dir rune) bool{
	fmt.Println(word, i, x, y, string(dir))
	velX, velY := 0, 0
	if dir == 'h' {
		velX = 1
		x -= i
	} else if dir == 'v'{
		velY = 1
		y -= i
		
	} else {
		fmt.Println("unknown direction in can place", string(dir), "with word", word)
	}
	if x < 0 || y < 0 || x + (velX * len(word)) >= len(cw.data[0]) || y + (velY * len(word)) >= len(cw.data){
		//fmt.Println("out of crossword bounds, on left or top")
		return false
	}
	if dir == 'h' {
		if x > 0 && cw.data[y][x - 1] != 0 {
			return false
		}
	} else if dir == 'v' {
		if y > 0 && cw.data[y - 1][x] != 0 {
			return false
		}
	}
	
	
	//fmt.Println("word will end at x", (x + (velX * len(word))), ", y", )
	//fmt.Println("attempting to place", word, "at", x, ",", y, "in direction", string(dir))
	count := 0
	for count < len(word){
		//Checks we need to make
		//place must be empty OR non empty AND char is what we need
		
		if cw.data[y][x] != 0 && string(cw.data[y][x]) != string(word[count]) {
			fmt.Println("failed because char already there and different")
			return false
		} else if cw.data[y][x] == 0 || string(cw.data[y][x]) != string(word[count]) {
			if velX == 1{
				//we are moving right, check chars above and below
				if y > 0 && cw.data[y - 1][x] != 0 {
					fmt.Println("going right and char above")
					return false
				} else if y < len(cw.data) - 1 &&cw.data[y + 1][x] != 0 {
					fmt.Println("going right and char below")
					return false
				}
			} else if velY == 1{
				//we are moving down, check chars left and right
				if x > 0 && cw.data[y][x - 1] != 0 {
					fmt.Println("going down and char to the left")
					return false
				} else if x < len(cw.data[0]) - 1 && cw.data[y][x + 1] != 0 {
					fmt.Println("going down and char to the right")
					return false
				}
			}
		}
		x += velX
		y += velY
		count++
	}
	if x > len(cw.data) - 1 || y > len(cw.data) - 1 {
		//fmt.Println("failed because word does not fit within boundaries of crossword [right or bottom]")
		return false
	}
	return true
}

func (cw Crossword) getPossibleMoves(word string) []Move{
	possibleMoves := make([]Move, 0)
	for x := 0; x < len(cw.data); x++ { //probably not the best practice, but can do as cw always square
		for y := 0; y < len(cw.data); y++ {
			curChar := cw.data[y][x]
			for i, c := range word {
				if c == curChar {
					//fmt.Println("")
					//we need to see if word fits here
					if cw.CanPlaceWord(word, i, x, y, 'v'){
						//fmt.Println("we can place vertically at", x, ",", y, "with i", i, "?")
						possibleMoves = append(possibleMoves, Move{x, y - i, 'v'})
						//fmt.Println(possibleMoves)
					} else if cw.CanPlaceWord(word, i, x, y, 'h'){
						//fmt.Println("we can place horizontally at", x, ",", y, "with i", i, "?")
						possibleMoves = append(possibleMoves, Move{x - i, y, 'h'})
						//fmt.Println(possibleMoves)
					}
				}
			}
		}
	}
	return possibleMoves
}


func makeEmptyCrossword(size int, wl *WordList) *Crossword{
	dat := make([][]rune, size)
	for i := range dat {
		dat[i] = make([]rune, size)
	}
	return &Crossword{dat, wl}
}

func (cw *Crossword) trim(){
	//there should never be empty rows in the middle.
	//ie. when scanning downards, as soon as we meet a character
	//no more empty rows till the bottom ones
	
	//remove empty cols starting from left
	lowerX := 0
	for cw.isEmptyCol(lowerX) && lowerX < len(cw.data[0]) - 1{
		lowerX++
	}
	
	//remove empty cols starting from right
	upperX := len(cw.data[0]) - 1
	for cw.isEmptyCol(upperX) && upperX > 0 {
		upperX--
	}
	upperX++
	
	//remove empty rows starting from top
	lowerY := 0
	for cw.isEmptyRow(lowerY) && lowerY < len(cw.data) - 1 {
		lowerY++
	}
	
	//remove empty rows starting from bottom
	upperY := len(cw.data) - 1
	for cw.isEmptyRow(upperY) && upperY > 0 {
		upperY--
	}
	upperY++
	
	//now trim according to the values of lowerX, upperX, lowerY, upperY
	fmt.Println(lowerX, upperX, lowerY, upperY)
	cw.Print()
	newData := make([][]rune, upperY - lowerY)
	copy(newData, cw.data[lowerY:upperY])
	
	cw.data = newData
	fmt.Println("LEN OF NEW DATA", len(newData), "and cw", len(cw.data))
	fmt.Println("wooorking")
	cw.Print()
	//cw.data = cw.data[lowerY:upperY]
	for i := range cw.data {
		newRow := make([]rune, upperX - lowerX)
		copy(newRow, cw.data[i][lowerX:upperX])
		cw.data[i] = newRow
	}
	cw.Print()
	
}

func (cw Crossword) render(filename string){
	var b strings.Builder
	
	fmt.Println(len(cw.data))
	
	cw.trim()
	
	fmt.Println(len(cw.data))
	
	b.WriteString(`
	<html>
		<head>
		<style>
			table {
				width: 50%;
				height: 200px;
				border: 2px solid black;
			}
			
			td {
				border: 2px solid black;
				border-collapse: collapse;
				text-align:center;
			}
			
			.blackSquare {
				background: black;
			}
			
			.content {
				width:70px;
				height: 70px;
				font-size: 24;
				line-height: 3;
			}
			 
		</style>
		</head>
		<body>
		<table>`)
	for y := 0; y < len(cw.data); y++ {
		b.WriteString("<tr>")
		for x := 0; x < len(cw.data[y]); x++ {
			val := cw.data[y][x]
			if val == 0 {
				b.WriteString("<td><div class='content blackSquare'></div></td>")
			} else {
				b.WriteString("<td><div class='content'>" + string(val) + "</div></td>")
			}
			
		}
		b.WriteString("</tr>")
	}
	b.WriteString(`</table>
		</body>
	</html>`)
	ioutil.WriteFile(filename, []byte(b.String()), 0644)
}

func (cw Crossword) fill() {
	fringe := make([]*Crossword, 0)
	//start with each word placed randomly, then expand search tree
	for i := range cw.wl.words {
		//create a crossword with word placed and wordList truncated of word
		newWl := cw.wl.RemoveWord(i)
		newCw := cw.Copy()
		newCw.wl = newWl
		newCw.place(cw.wl.words[i], 1, len(cw.data)/2, 'h')
		fringe = append(fringe, newCw)
	}
	fmt.Println("starting fringe: size", len(fringe))
	for len(fringe) > 0 {
		curCw := fringe[len(fringe) - 1]
		fmt.Println("\nWorking on the following crossword")
		curCw.Print()
		fringe = fringe[:len(fringe) - 1]
		//we pick a word from the word list, see if there is a common letter. if there is, see if fits.
		word := curCw.wl.words[0]
		fmt.Println("and trying to fit word", word)
		moves := curCw.getPossibleMoves(word)
		//fmt.Println("found", len(moves), "possible moves")
		for _, move := range moves {
			newWl := curCw.wl.RemoveWord(0)
			newCw := curCw.Copy()
			newCw.wl = newWl
			//fmt.Println("going to perform move", move)
			newCw.place(word, move.x, move.y, move.d)
			if len(newCw.wl.words) == 0{
				fmt.Println("we have a solution")
				newCw.Print()
				newCw.render("cw.html")
				os.Exit(1)
			} else {
				fmt.Println("words left", newCw.wl.words, len(newCw.wl.words), newCw.wl.words[0])
			}
			//newCw.Print()
			fringe = append(fringe, newCw)
		}
		
	}
}

func (cw Crossword) place(word string, x int, y int, direction rune) {
	velX, velY := 0, 0
	if direction == 'v'{
		velY = 1
	} else if direction == 'h'{
		velX = 1
	} else {
		fmt.Println("unknown direction in place'", direction, "' with word", word)
		os.Exit(1)
	}
	count := 0
	for count < len(word){
		cw.data[y][x] = rune(word[count])
		x += velX
		y += velY
		count++
	}
}

func (cw Crossword) Copy() *Crossword{
	newData := make([][]rune, len(cw.data))
	for i := range cw.data{
		newData[i] = make([]rune, len(cw.data))
		copy(newData[i], cw.data[i])
	}
	newWl := cw.wl.Copy()
	return &Crossword{newData, newWl}
}

func (cw Crossword) Print() {
	for _, row := range cw.data {
		for _, col := range row {
			if col != 0 {
				fmt.Print(string(col))
			} else {
				fmt.Print("_")
			}
			
		}
		fmt.Println()
	}
}

func main(){
	wl := inflateWordList("words.txt")
	fmt.Println(wl, wl.words, len(wl.words))
	minSize := wl.DetermineMinCrosswordSize()
	cw := makeEmptyCrossword(2*minSize, wl)
	fmt.Println(wl.MaxCharSurface())
	cw.Print()
	cw.fill()
}