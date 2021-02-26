package main

import (
    "fmt"
    "github.com/nsf/termbox-go"
    "math/rand"
    "os"
    "time"
)

//============================================================================

type box struct {
    Score int
    Size  int
    Cube  [][]int
}

type angle int

const (
    Rotate90  angle = 90
    Rotate180 angle = 180
    Rotate270 angle = 270
)

type direction int

const (
    Up direction = iota
    Down
    Left
    Right
)

//============================================================================

func CreateBox(size int) *box {
    box0 := box{}
    box0.Size = size
    box0.Cube = make([][]int, size)
    for i := range box0.Cube {
        box0.Cube[i] = make([]int, size)
    }
    return &box0
}

func (b *box) Print() {
    for i := range b.Cube {
        for j := range b.Cube[i] {
            fmt.Printf("%6d", b.Cube[i][j])
        }
        fmt.Printf("\n")
    }
}

func (b *box) Rotate(an angle) {
    var boxTemp *box = CreateBox(b.Size)
    var size int = b.Size - 1
    switch an {
    case Rotate90:
        for i := range b.Cube {
            for j := range b.Cube[i] {
                boxTemp.Cube[j][size-i] = b.Cube[i][j]
            }
        }
        b.Cube = boxTemp.Cube
    case Rotate180:
        for i := range b.Cube {
            for j := range b.Cube[i] {
                boxTemp.Cube[size-i][size-j] = b.Cube[i][j]
            }
        }
        b.Cube = boxTemp.Cube
    case Rotate270:
        for i := range b.Cube {
            for j := range b.Cube[i] {
                boxTemp.Cube[size-j][i] = b.Cube[i][j]
            }
        }
        b.Cube = boxTemp.Cube
    default:
        fmt.Println("err")
    }
}

func (b *box) CanMerge() bool {
    for i := range b.Cube {
        for j := range b.Cube[i] {
            if b.Cube[i][j] == 0 {
                return true
            } else if (i-1 >= 0 && b.Cube[i][j] == b.Cube[i-1][j]) || (j-1 >= 0 && b.Cube[i][j] == b.Cube[i][j-1]) {
                return true
            }
        }
    }
    return false
}

func (b *box) Merge(dir direction) bool {
    var merged bool = false
    switch dir {
    case Up:
        for i := 1; i < b.Size; i++ {
            for j := range b.Cube[i] {
                for t := i; t > 0; t-- {
                    if b.Cube[t][j] == b.Cube[t-1][j] {
                        merged = true
                        b.Cube[t-1][j] = b.Cube[t][j] << 1
                        b.Cube[t][j] = 0
                        
                    } else if b.Cube[t-1][j] == 0 {
                        b.Cube[t-1][j] = b.Cube[t][j]
                        b.Cube[t][j] = 0
                        merged = true
                    }
                }
            }
        }
    case Right:
        b.Rotate(Rotate270)
        merged = b.Merge(Up)
        b.Rotate(Rotate90)
    case Down:
        b.Rotate(Rotate180)
        merged = b.Merge(Up)
        b.Rotate(Rotate180)
    case Left:
        b.Rotate(Rotate90)
        merged = b.Merge(Up)
        b.Rotate(Rotate270)
    }
    
    return merged
}

func (b *box) RandGenerate(maxNum int) bool {
    var list = make([][]int, b.Size*b.Size)
    var size = 0
    for i := range b.Cube {
        for j := range b.Cube[i] {
            if b.Cube[i][j] == 0 {
                list[size] = []int{i, j}
                size++
            }
        }
    }
    
    if size == 0 {
        return false
    }
    
    randArr := RandArr(size, maxNum)
    
    for i := range randArr {
        point := randArr[i]
        b.Cube[list[point][0]][list[point][1]] = 2 << rand.Intn(2)
    }
    
    return true
}

func RandArr(cap int, count int) []int {
    arr := make([]int, cap)
    
    for i := 0; i < cap; i++ {
        arr[i] = i
    }
    
    rand.Seed(time.Now().Unix())
    rand.Shuffle(len(arr), func(i int, j int) {
        arr[i], arr[j] = arr[j], arr[i]
    })
    
    if cap >= count {
        return arr[:count]
    } else {
        return arr
    }
}

//============================================================================

func DrewStr(str string, x0, y0 int, fg, bg termbox.Attribute) error {
    for i, c := range str {
        termbox.SetCell(x0+i, y0-1, c, fg, bg)
    }
    return termbox.Flush()
}

func (b *box) DrewCubeBorder(x0, y0 int, fg, bg termbox.Attribute) error {
    for y := 0; y <= b.Size; y++ {
        for x := 0; x < 5*b.Size; x++ {
            termbox.SetCell(x0+x, y0+y*2, '-', fg, bg)
        }
        for x := 0; x <= len(b.Cube); x++ {
            termbox.SetCell(x0+x*5, y0+y*2, '+', fg, bg)
            
            if y < len(b.Cube) {
                termbox.SetCell(x0+x*5, y0+y*2+1, '|', fg, bg)
            }
        }
    }
    return termbox.Flush()
}

func (b *box) DrewCubeValue(x0, y0 int, fg, bg termbox.Attribute) error {
    for y := 0; y < b.Size; y++ {
        for x := 0; x < b.Size; x++ {
            if v:= b.Cube[y][x]; v > 0 {
                DrewStr(fmt.Sprint(v), x0+x*5, y0+y*2, fg, bg)
            }
           
        }
    }
    return termbox.Flush()
}

func (b *box) Drew(ox, oy int) error {
    fg := termbox.ColorWhite
    bg := termbox.ColorBlack
    _ = termbox.Clear(fg, bg)
    
    strScore := "SCORE: " + fmt.Sprint(b.Score)
    DrewStr(strScore, 2, 2, termbox.ColorRed, termbox.ColorYellow)
    
    return termbox.Flush()
}

func (b *box) Drew2(ox, oy int) error {
    fg := termbox.ColorWhite
    bg := termbox.ColorBlack
    termbox.Clear(fg, bg)
    
    str := " SCORE: " + fmt.Sprint(b.Score)
    for n, c := range str {
        termbox.SetCell(ox+n, oy-1, c, fg, bg)
    }
    str = "ESC:exit " + "Enter:replay"
    for n, c := range str {
        termbox.SetCell(ox+n, oy-2, c, fg, bg)
    }
    str = "Play with arrow key!"
    for n, c := range str {
        termbox.SetCell(ox+n, oy-3, c, fg, bg)
    }
    fg = termbox.ColorGreen
    //bg = termbox.ColorGreen
    
    for y := 0; y <= len(b.Cube); y++ {
        for x := 0; x < 5*len(b.Cube); x++ {
            termbox.SetCell(ox+x, oy+y*2, '-', fg, bg)
        }
        for x := 0; x <= len(b.Cube); x++ {
            termbox.SetCell(ox+x*5, oy+y*2, '+', fg, bg)
            
            if y < len(b.Cube) {
                termbox.SetCell(ox+x*5, oy+y*2+1, '|', fg, bg)
            }
        }
    }
    fg = termbox.ColorYellow
    bg = termbox.ColorBlack
    for i := range b.Cube {
        for j := range b.Cube[i] {
            if b.Cube[i][j] > 0 {
                str := fmt.Sprint(b.Cube[i][j])
                for n, char := range str {
                    termbox.SetCell(ox+j*5+1+n, oy+i*2+1, char, fg, bg)
                }
            }
        }
    }
    return termbox.Flush()
}

//============================================================================
func main() {
    b := CreateBox(4)
    b.RandGenerate(3)
    
    err := termbox.Init()
    if err != nil {
        panic(err)
    }
    
    defer termbox.Close()
   
    rand.Seed(time.Now().UnixNano())
    for  {
    
   
    termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
    //DrewStr("12345667", 2, 2, termbox.ColorRed, termbox.ColorYellow)
    
    b.DrewCubeBorder(5, 5, termbox.ColorYellow, termbox.ColorBlack)
    b.DrewCubeValue(5+1, 5+2, termbox.ColorWhite, termbox.ColorBlack)
    
    
    
    
    ev := termbox.PollEvent()
    
    switch ev.Type {
    case termbox.EventKey:
        switch ev.Key {
        case termbox.KeyArrowUp:
        
        case termbox.KeyArrowDown:
        case termbox.KeyArrowLeft:
        case termbox.KeyArrowRight:
        case termbox.KeyEsc, termbox.KeyEnter:
        default:
            //t.Print(0, 3)
        }
    
    case termbox.EventResize:
    case termbox.EventError:
        panic(ev.Err)
    }
    }
}

func main2() {
    fmt.Println("上：W    下：S    左:A     右：D")
    b := CreateBox(5)
    b.RandGenerate(3)
    b.Print()
    
    var c int16
    merged := false
    score := 0
    for {
        if !b.CanMerge() {
            fmt.Println("over:", score)
            os.Exit(0)
        }
        
        fmt.Scanf("%c", &c)
        
        switch c {
        case 'w':
            merged = b.Merge(Up)
        case 's':
            merged = b.Merge(Down)
        case 'a':
            merged = b.Merge(Left)
        case 'd':
            merged = b.Merge(Right)
        default:
            continue
        }
        
        if merged {
            score += 2
            fmt.Println("========================================: ", score)
            b.RandGenerate(2)
            b.Print()
        }
        
    }
}
