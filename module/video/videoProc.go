package video

import (
	"fmt"
	"math"
	"math/rand"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func duration(url string) (int64, error) {

	sysType := runtime.GOOS
	if sysType == "windows" {
		rand.Seed(2)
		i := rand.Intn(100)
		return int64(i), nil
	}

	c := fmt.Sprintf(`ffprobe -v error -show_entries format=duration -of default=noprint_wrappers=1:nokey=1 -i "%s"`, url)
	fmt.Println(c)
	cmd := exec.Command("bash", "-c", c)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return 0, err
	}
	o := strings.TrimSpace(string(out))
	f64, _ := strconv.ParseFloat(o, 64)
	i64 := int64(math.Round(f64))
	return i64, nil
}
