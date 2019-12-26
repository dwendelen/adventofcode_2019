package main

import (
	"../../util"
	"fmt"
	"log"
	"strconv"
)

const (
	NORTH = "north"
	SOUTH = "south"
	EAST  = "east"
	WEST  = "west"

	ASTRO  = "astrolabe"
	MANI   = "manifold"
	MOUSE  = "mouse"
	MUG    = "mug"
	SAND   = "sand"
	BROCH  = "space law space brochure"
	WREATH = "wreath"

	CHECK = "Security Checkpoint"
)

var items = [...]string{ASTRO, MANI, MOUSE, MUG, SAND, BROCH, WREATH}

func main() {
	program := load("day25/input.txt")

	robot := newRobot(program)
	robot.run()
	robot.waitForCommand()

	//Take everything south
	robot.goDir(SOUTH)
	robot.take(BROCH)
	robot.goDir(SOUTH)
	robot.take(MOUSE)
	robot.goDir(SOUTH)
	robot.take(ASTRO)
	robot.goDir(SOUTH)
	robot.take(MUG)
	robot.goDir(NORTH)
	robot.goDir(NORTH)
	robot.goDir(WEST)
	robot.goDir(NORTH)
	robot.goDir(NORTH)
	robot.take(WREATH)
	robot.goDir(SOUTH)
	robot.goDir(SOUTH)
	robot.goDir(EAST)
	robot.goDir(NORTH)
	// Back at sick bay
	robot.goDir(WEST)
	robot.take(SAND)
	robot.goDir(NORTH)
	robot.take(MANI)
	robot.goDir(SOUTH)
	robot.goDir(WEST)
	room := robot.goDir(WEST)
	if room != CHECK {
		log.Fatalln("Oops wrong room")
	}
	robot.command("inv")
	robot.crackCode()
}

type Robot struct {
	program       map[int64]int64
	input, output chan int64
}

func newRobot(program map[int64]int64) *Robot {
	return &Robot{program, make(chan int64), make(chan int64)}
}

func (r *Robot) crackCode() {
	code := 0
	room := CHECK
	for room == CHECK {
		r.dropAll()
		subcode := code
		for i := 0; i < len(items); i++ {
			if subcode%2 == 1 {
				r.take(items[i])
			}
			subcode = subcode / 2
		}
		room = r.goDir(WEST)
		code++
	}

	fmt.Println("Code:", code)
	subcode := code
	for i := 0; i < len(items); i++ {
		if subcode%2 == 1 {
			fmt.Println(items[i])
		}
		subcode = subcode / 2
	}
}

func (r *Robot) take(item string) {
	r.command("take " + item)
}

func (r *Robot) dropAll() {
	for _, item := range items {
		r.command("drop " + item)
	}
}

func (r *Robot) goDir(dir string) string {
	return r.command(dir)
}

func (r *Robot) command(cmd string) string {
	fmt.Println("<<<", cmd)
	for _, char := range cmd {
		r.input <- int64(char)
	}
	r.input <- '\n'
	return r.waitForCommand()
}

func (r *Robot) run() {
	go run(r.program, r.input, r.output)
}

var buffer = make([]byte, 0)

func (r *Robot) waitForCommand() string {
	var room string
	for {
		for {
			i := <-r.output
			if i > 255 {
				log.Fatalln("Can't handle this")
			}
			if i == '\n' {
				break
			}
			buffer = append(buffer, byte(i))
		}
		asString := string(buffer)
		buffer = buffer[:0]
		fmt.Println(">>>", asString)
		if asString == "Command?" {
			break
		} else if len(asString) > 6 && asString[:2] == "==" {
			room = asString[3 : len(asString)-3]
		}
	}
	return room
}

func load(file string) map[int64]int64 {
	values := make(map[int64]int64)
	i := int64(0)

	util.ReadCommaSeparated(file, func(in string) error {
		asInt, err := strconv.Atoi(in)
		if err != nil {
			return err
		}
		values[i] = int64(asInt)
		i++
		return nil
	})

	return values
}

func run(program map[int64]int64, input chan int64, output chan int64) {
	values := make(map[int64]int64)
	for k, v := range program {
		values[k] = v
	}

	pos := int64(0)
	offset := int64(0)

	for {
		mode3, mode2, mode1, code := parse(readValue(values, pos))

		switch code {
		case 99: // halt
			close(output)
			return
		case 1: // add
			val1 := read(values, offset, readValue(values, pos+1), mode1)
			val2 := read(values, offset, readValue(values, pos+2), mode2)
			write(values, offset, readValue(values, pos+3), mode3, val1+val2)
			pos += 4
		case 2: // multiply
			val1 := read(values, offset, readValue(values, pos+1), mode1)
			val2 := read(values, offset, readValue(values, pos+2), mode2)
			write(values, offset, readValue(values, pos+3), mode3, val1*val2)
			pos += 4
		case 3: // input
			val := <-input
			write(values, offset, readValue(values, pos+1), mode1, val)
			pos += 2
		case 4: // output
			val := read(values, offset, readValue(values, pos+1), mode1)
			output <- val
			pos += 2
		case 5: // jump-if-true
			cond := read(values, offset, readValue(values, pos+1), mode1)
			addr := read(values, offset, readValue(values, pos+2), mode2)
			if cond != 0 {
				pos = addr
			} else {
				pos += 3
			}
		case 6: // jump-if-false
			cond := read(values, offset, readValue(values, pos+1), mode1)
			addr := read(values, offset, readValue(values, pos+2), mode2)
			if cond == 0 {
				pos = addr
			} else {
				pos += 3
			}
		case 7: // less than
			val1 := read(values, offset, readValue(values, pos+1), mode1)
			val2 := read(values, offset, readValue(values, pos+2), mode2)
			var res int64
			if val1 < val2 {
				res = 1
			} else {
				res = 0
			}
			write(values, offset, readValue(values, pos+3), mode3, res)
			pos += 4
		case 8: // equals
			val1 := read(values, offset, readValue(values, pos+1), mode1)
			val2 := read(values, offset, readValue(values, pos+2), mode2)
			var res int64
			if val1 == val2 {
				res = 1
			} else {
				res = 0
			}
			write(values, offset, readValue(values, pos+3), mode3, res)
			pos += 4
		case 9: //adjust offset
			val := read(values, offset, readValue(values, pos+1), mode1)
			offset += val
			pos += 2
		default:
			log.Fatal("Unknown op code " + string(code))
		}
	}
}

func readValue(values map[int64]int64, nb int64) int64 {
	res, ok := values[nb]
	if ok {
		return res
	} else {
		return 0
	}
}

func read(values map[int64]int64, offset int64, nb int64, mode int64) int64 {
	switch mode {
	case 0:
		return readValue(values, nb)
	case 1:
		return nb
	case 2:
		return readValue(values, nb+offset)
	default:
		log.Fatal("Illegal mode " + string(mode))
		return 0
	}
}

func write(values map[int64]int64, offset int64, nb int64, mode int64, val int64) {
	switch mode {
	case 0:
		values[nb] = val
	case 1:
		log.Fatal("You can not use immediate mode with a store")
	case 2:
		values[nb+offset] = val
	default:
		log.Fatal("Illegal mode " + string(mode))
	}
}

func parse(instruction int64) (mode3, mode2, mode1, code int64) {
	if instruction >= 100000 {
		log.Fatal("Illegal instruction " + string(instruction))
	}
	code = instruction % 100
	rest := instruction / 100
	mode1 = rest % 10
	rest = rest / 10
	mode2 = rest % 10
	rest = rest / 10
	mode3 = rest % 10

	return
}
