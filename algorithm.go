package main

import (
	"fmt"
	"sort"
)

// 左右两个t1和t2 如果`t2.Px-(t1.Px+t1.Width)<delta`被认为是在一个组中
// 上下两个t1和t2 如果`t2.Px-t1.Px<delta`被认为是在一个组中
const (
	delta float32 = 8
)

var (
	lessFunc func(t1, t2 Target) bool
)

// Group 同种类型的目标(网卡/光口)距离比较近的进行分组
// 例如: [n]Group n组,每组Row行Col列,起始点(MinX,MinY)
type Group struct {
	MinX float32
	MinY float32
	Row  int
	Col  int
}

func (g *Group) round() Group {
	g.MinX = round(g.MinX, 1)
	g.MinY = round(g.MinY, 1)
	return *g
}

// TargetSlice implements sort.Interface
type TargetSlice []Target

func (a TargetSlice) Len() int {
	return len(a)
}

func (a TargetSlice) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a TargetSlice) Less(i, j int) bool {
	return lessFunc(a[i], a[j])
}

func lessX(t1, t2 Target) bool { return t1.Px < t2.Px }

func lessY(t1, t2 Target) bool { return t1.Py < t2.Py }

func lessXY(t1, t2 Target) bool {
	if dEqual(t1.Px, t2.Px) {
		return t1.Py < t2.Py
	}
	return t1.Px < t2.Px
}

// SortByX 按X排序
func (a TargetSlice) SortByX() {
	lessFunc = lessX
	sort.Sort(a)
}

// SortByY 按Y排序
func (a TargetSlice) SortByY() {
	lessFunc = lessY
	sort.Sort(a)
}

// Sort 按XY排序
func (a TargetSlice) Sort() {
	lessFunc = lessXY
	sort.Sort(a)
}

func dEqual(x, y float32) bool {
	return (x-y)*(x-y) < 2*delta*delta
}

// round(123.456, 10) = 120
// round(123.456, .1) = 123.5
// round(123.456, 5) = 125
func round(number, accuracy float32) float32 {
	return float32(int(number/accuracy+.5)) * accuracy
}

func groupTargets(ts TargetSlice) []Group {
	if len(ts) == 0 {
		return make([]Group, 0)
	}
	ts.adjust()
	var AvgWidth, AvgHeight float32
	for _, t := range ts {
		AvgWidth += t.Width
		AvgHeight += t.Height
	}
	AvgWidth /= float32(len(ts))
	AvgHeight /= float32(len(ts))
	ret := make([]Group, 0)
	var tmp *Group
	for index, t := range ts {
		if index == 0 {
			tmp = &Group{t.Px, t.Py, 1, 1}
			continue
		}
		if t.Px-ts[index-1].Px-AvgWidth > delta {
			ret = append(ret, tmp.round())
			tmp = &Group{t.Px, t.Py, 1, 1}
			continue
		}
		if dEqual(t.Px, tmp.MinX) {
			if t.Px < tmp.MinX {
				tmp.MinX = t.Px
			}
			n := int((t.Py-tmp.MinY)/AvgHeight + 1.5)
			if n > tmp.Row {
				tmp.Row = n
			}
		}
		if dEqual(t.Py, tmp.MinY) {
			if t.Py < tmp.MinY {
				tmp.MinY = t.Py
			}
			n := int((t.Px-tmp.MinX)/AvgWidth + 1.5)
			if n > tmp.Col {
				tmp.Col = n
			}
		}
	}
	ret = append(ret, tmp.round())
	loger <- fmt.Sprintf("分组结果: %#v\n", ret)
	return ret
}

// 调整位置,横纵坐标相近的点自动对齐
func (a TargetSlice) adjust() {
	var i, j int
	var sum float32
	// a[...].Px 对齐
	a.SortByX()
	for j = range a {
		if !dEqual(a[i].Px, a[j].Px) {
			sum = 0
			for _, tmp := range a[i:j] {
				sum += tmp.Px
			}
			sum /= float32(j - i)
			for index := range a[i:j] {
				a[index+i].Px = sum
			}
			i = j
		}
	}
	sum = 0
	for _, tmp := range a[i:] {
		sum += tmp.Px
	}
	sum /= float32(j - i + 1)
	for index := range a[i:] {
		a[index+i].Px = sum
	}
	// a[...].Py 对齐
	i = 0
	a.SortByY()
	for j = range a {
		if !dEqual(a[i].Py, a[j].Py) {
			sum = 0
			for _, tmp := range a[i:j] {
				sum += tmp.Py
			}
			sum /= float32(j - i)
			for index := range a[i:j] {
				a[index+i].Py = sum
			}
			i = j
		}
	}
	sum = 0
	for _, tmp := range a[i:] {
		sum += tmp.Py
	}
	sum /= float32(j - i + 1)
	for index := range a[i:] {
		a[index+i].Py = sum
	}
	a.Sort()
	// loger <- fmt.Sprintf("调整目标位置后: %#v\n", a)
}

func serial(col, start int, width float32) string {
	str := fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%.1f" height="16" font-size="8">`, width*float32(col))
	fifthW := width / 5
	for i := 0; i < col; i++ {
		n := start + i*2
		str += fmt.Sprintf(`<g transform="translate(%.1f)"><g transform="translate(%.1f)"><text x="2" y="8">%d</text><text x="2" y="16">▲</text></g><g transform="translate(%.1f)"><text x="2" y="8">%d</text><text x="2" y="16">▼</text></g></g>`,
			float32(i)*width, fifthW, n, fifthW*3, n+1)
	}
	str += `</svg>`
	return str
}

func init() {
	loger = make(chan string, 10) // 得尽早执行
	lessFunc = lessXY
}
