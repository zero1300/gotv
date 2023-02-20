package video

import (
	"fmt"
	"math"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func duration(url string) (time.Duration, error) {
	c := fmt.Sprintf("ffprobe -v error -show_entries format=duration -of default=noprint_wrappers=1:nokey=1 -i %s", url)
	cmd := exec.Command("bash", "-c", c)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		return time.Duration(0), err
	}
	o := strings.TrimSpace(string(out))
	fmt.Println(o)
	f64, _ := strconv.ParseFloat(o, 64)
	fp := f64 * math.Pow(1000.0, 3.0)
	td := time.Duration(int64(math.Round(fp)))
	return td, nil
}
