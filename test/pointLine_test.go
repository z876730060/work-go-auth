package test

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"testing"
)

type Point struct {
	x float32
	y float32
}

type Line struct {
	p1  Point //基点
	p2  Point //终止点
	vec Point //向量
}

func TestGetPointLine(t *testing.T) {
	lineToPoint(Point{1, 0.5}, Point{2, 0}, Point{3, 3}, Point{2, 5}, Point{1, 4}, Point{2, 3}, Point{-1, 3})
}

func (p *Point) print() {
	fmt.Println("x:", p.x, "y:", p.y)
}

func (p *Point) String() string {
	return fmt.Sprintf("x:%f, y:%f", p.x, p.y)
}

func lineToPoint(p1 ...Point) {
	if len(p1) < 2 {
		fmt.Println("至少需要两个点")
		return
	}

	maxX, maxY, minX, minY := float32(0), float32(0), float32(0), float32(0)
	for i, p := range p1 {
		if i == 0 {
			maxX, maxY, minX, minY = p.x, p.y, p.x, p.y
		} else {
			if p.x > maxX {
				maxX = p.x
			}
			if p.y > maxY {
				maxY = p.y
			}
			if p.x < minX {
				minX = p.x
			}
			if p.y < minY {
				minY = p.y
			}
		}
	}
	fmt.Println("maxX:", maxX, "maxY:", maxY, "minX:", minX, "minY:", minY)

	var slope float32 = 0
	if maxX != minX {
		slope = float32(maxY-minY) / float32(maxX-minX)
	}

	if len(p1) == 2 {
		fmt.Println("两点无法完成封闭图形")
	}

	isOneLine := true
	for i, p := range p1 {
		if i == 0 || i == 1 {
			continue
		}
		if !checkIsOneLine(slope, p) {
			fmt.Println("不是一条线")
			isOneLine = false
			break
		}
	}
	if isOneLine {
		fmt.Println("是一条线")
		return
	}
	fmt.Println("继续封闭图形逻辑")

	var finalLines []Line
	for {
		lines := make([]Line, 0)
		findLine(&lines, p1)

		x, y := float32(0), float32(0)
		for _, line := range lines {
			x += line.vec.x
			y += line.vec.y
		}
		if x == 0 && y == 0 {
			isIntersect := false
			for i := 0; i < len(lines); i++ {
				for j := i + 1; j < len(lines); j++ {
					if i == j {
						continue
					}
					if checkIsIntersect(lines[i], lines[j]) {
						fmt.Println("线段相交")
						isIntersect = true
						break
					}

				}
			}
			if !isIntersect {
				fmt.Println("封闭图形")
				finalLines = lines
				break
			}
		}
	}
	// 打印最终的线
	for _, line := range finalLines {
		fmt.Println(line)
	}
}

func findLine(lines *[]Line, p1s []Point) {
	uesdIndex := make([]int, 1)
	currentIndex := 0

	for {
		index := rand.IntN(len(p1s))
		if slices.Contains(uesdIndex, index) {
			continue
		}
		uesdIndex = append(uesdIndex, index)
		*lines = append(*lines, Line{p1s[currentIndex], p1s[index], Point{p1s[index].x - p1s[currentIndex].x, p1s[index].y - p1s[currentIndex].y}})
		currentIndex = index
		if len(uesdIndex) == len(p1s) {
			// 已经是最后一个点
			*lines = append(*lines, Line{p1s[currentIndex], p1s[0], Point{p1s[0].x - p1s[currentIndex].x, p1s[0].y - p1s[currentIndex].y}})
			return
		}
	}
}

func checkIsOneLine(slope float32, point Point) bool {
	if point.y == point.x*slope {
		return true
	}
	return false
}

func checkIsIntersect(line1, line2 Line) bool {
	if line1.vec.x == 0 || line1.vec.y == 0 || line2.vec.x == 0 || line2.vec.y == 0 {
		// 有一条线是水平或垂直线
		if line1.vec.x == 0 {
			// 垂线
			if line2.vec.x == 0 {
				// 另一条线也是垂线
				if line1.p1.x == line2.p1.x {
					if (line1.p1.y >= line2.p1.y && line1.p1.y <= line2.p2.y) || (line1.p1.y >= line2.p2.y && line1.p1.y <= line2.p1.y) ||
						(line1.p2.y >= line2.p1.y && line1.p2.y <= line2.p2.y) || (line1.p2.y >= line2.p2.y && line1.p2.y <= line2.p1.y) {
						return true
					}
				} else {
					// 两条线在不同垂直线上
					return false
				}
			} else {
				if (line1.p1.x >= line2.p1.x && line1.p1.x <= line2.p2.x) || (line1.p1.x >= line2.p2.x && line1.p1.x <= line2.p1.x) {
					return true
				}
			}
		} else if line1.vec.y == 0 {
			if line2.vec.y == 0 {
				// 两条线都是水平，判断y坐标是否相等
				if line1.p1.y == line2.p1.y {
					// 继续判断两天线段是否有重叠
					if (line1.p1.x >= line2.p1.x && line1.p1.x <= line2.p2.x) || (line1.p1.x >= line2.p2.x && line1.p1.x <= line2.p1.x) ||
						(line1.p2.x >= line2.p1.x && line1.p2.x <= line2.p2.x) || (line1.p2.x >= line2.p2.x && line1.p2.x <= line2.p1.x) {
						return true
					}
				} else {
					// 两条线在不同水平线上
					return false
				}
			} else {
				// 另一条线不是水平或垂直线,判断y坐标是否在水平线上
				if (line1.p1.y > line2.p1.y && line1.p1.y < line2.p2.y) || (line1.p1.y > line2.p2.y && line1.p1.y < line2.p1.y) {
					return true
				}
			}
		} else if line2.vec.x == 0 {
			// 第一条线已经不是水平或垂直线了,只要第二条线的x坐标在第一条线的x坐标范围内,就相交
			if (line2.p1.x > line1.p1.x && line2.p1.x < line1.p2.x) || (line2.p1.x > line1.p2.x && line2.p1.x < line1.p1.x) {
				return true
			}
		} else if line2.vec.y == 0 {
			// 另一条线是水平或垂直线,只要第二条线的y坐标在第二条线的y坐标范围内,就相交
			if (line2.p1.y > line1.p1.y && line2.p1.y < line1.p2.y) || (line2.p1.y > line1.p2.y && line2.p1.y < line1.p1.y) {
				return true
			}
		}
		return false
	} else {
		// 计算两条线的斜率
		slope1 := line1.vec.y / line1.vec.x
		slope2 := line2.vec.y / line2.vec.x

		// 如果斜率相同，并且存在点相同, 则相交
		if slope1 == slope2 && (line1.p1 == line2.p1 || line1.p1 == line2.p2 || line1.p2 == line2.p1 || line1.p2 == line2.p2) {
			return true
		}

		//TODO 处理垂直线
		k := line1.p1.y - line1.p1.x*slope1
		k1 := line2.p1.y - line2.p1.x*slope2

		x := (k1 - k) / (slope1 - slope2)
		y := slope1*x + k

		// 判断相交的点必须在两条线段的范围内, 这两条线必须是端点连接
		if ((x >= line1.p1.x && x <= line1.p2.x) || (x >= line1.p2.x && x <= line1.p1.x)) &&
			((y >= line1.p1.y && y <= line1.p2.y) || (y >= line1.p2.y && y <= line1.p1.y)) &&
			((x >= line2.p1.x && x <= line2.p2.x) || (x >= line2.p2.x && x <= line2.p1.x)) &&
			((y >= line2.p1.y && y <= line2.p2.y) || (y >= line2.p2.y && y <= line2.p1.y)) {
			if (x == line1.p1.x && y == line1.p1.y) || (x == line1.p2.x && y == line1.p2.y) ||
				(x == line2.p1.x && y == line2.p1.y) || (x == line2.p2.x && y == line2.p2.y) {
				return false
			}
			fmt.Println("这两个线有交点")
			return true
		}
		return false
	}
}
