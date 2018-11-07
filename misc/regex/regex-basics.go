package main

import (
	"fmt"
	"regexp"
	"strconv"
)

func main() {

	msg1 := "aaalisten tcp :8055: bind: permission denied."
	msg2 := "listen tcp :8055: bind: permission denied"

	pattern := regexp.MustCompile("listen tcp :\\d{1,4}: bind: permission denied")
	pattern2 := regexp.MustCompile("\\d{1,4}")

	fmt.Println(pattern.FindStringSubmatch(msg1))
	fmt.Println(pattern.MatchString(msg1))
	fmt.Println(pattern.MatchString(msg2))

	port, _ := strconv.Atoi(pattern2.FindString(msg2))
	fmt.Println(port)

}
