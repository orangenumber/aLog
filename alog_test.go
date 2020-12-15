// (c) 2020 Gon Y Yi. <https://gonyyi.com/copyright.txt>

package alog_test

import (
	"bytes"
	"github.com/gonyyi/alog"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

/*
    =====================================================================================================================
    BENCHMARK RESULTS
    =====================================================================================================================

	v0.1.0
		| Test                         | NoOfRun | Speed      | Mem     | Alloc       |
		|:-----------------------------|:--------|:-----------|:--------|:------------|
		| Benchmark_ALog_Print-12      | 169,268 | 6659 ns/op | 0 B/op  | 0 allocs/op |
		| Benchmark_ALog_Printf-12     | 167,683 | 7141 ns/op | 8 B/op  | 0 allocs/op |
		| Benchmark_ALog_Output-12     | 172,851 | 6769 ns/op | 0 B/op  | 0 allocs/op |
		| Benchmark_Printj-12          | 160,222 | 7491 ns/op | 96 B/op | 2 allocs/op |

	v0.1.2
		| Test                         | NoOfRun | Speed      | Mem     | Alloc       |
		|:-----------------------------|:--------|:-----------|:--------|:------------|
		| Benchmark_ALog_Print-12      | 165,474 | 6277 ns/op | 0 B/op  | 0 allocs/op |
		| Benchmark_ALog_Printf-12     | 183,186 | 6454 ns/op | 8 B/op  | 0 allocs/op |
		| Benchmark_ALog_Output-12     | 161,769 | 6454 ns/op | 0 B/op  | 0 allocs/op |
		| Benchmark_Printj-12          | 174,801 | 7990 ns/op | 0 B/op  | 0 allocs/op |

	v0.2.1 (support buffering)
		| Test                         | NoOfRun | Speed      | Mem    | Alloc       |
		|:-----------------------------|:--------|:-----------|:-------|:------------|
		| Benchmark_ALog_Print-12      | 432720  | 2718 ns/op | 0 B/op | 0 allocs/op |
		| Benchmark_ALog_Printf-12     | 411183  | 2769 ns/op | 0 B/op | 0 allocs/op |
		| Benchmark_ALog_Printj-12     | 336685  | 3250 ns/op | 0 B/op | 0 allocs/op |
		| Benchmark_ALog_Printf_Buf-12 | 3485294 | 334 ns/op  | 0 B/op | 0 allocs/op |
		| Benchmark_ALog_Print_Buf-12  | 3682484 | 355 ns/op  | 0 B/op | 0 allocs/op |
		| Benchmark_ALog_Printj_Buf-12 | 2053513 | 595 ns/op  | 0 B/op | 0 allocs/op |

	v1.1.0 (support Level)
		| Test                         | NoOfRun | Speed      | Mem    | Alloc       |
		|:-----------------------------|:--------|:-----------|:-------|:------------|
		| Benchmark_ALog_Print-12      | 432720  | 2718 ns/op | 0 B/op | 0 allocs/op |
		| Benchmark_ALog_Printf-12     | 411183  | 2769 ns/op | 0 B/op | 0 allocs/op |
		| Benchmark_ALog_Printj-12     | 336685  | 3250 ns/op | 0 B/op | 0 allocs/op |
		| Benchmark_ALog_Printf_Buf-12 | 3485294 | 334 ns/op  | 0 B/op | 0 allocs/op |
		| Benchmark_ALog_Print_Buf-12  | 3682484 | 355 ns/op  | 0 B/op | 0 allocs/op |
		| Benchmark_ALog_Printj_Buf-12 | 2053513 | 595 ns/op  | 0 B/op | 0 allocs/op |

*/

// =====================================================================================================================
// TEST
// =====================================================================================================================
func Test_ALog(t *testing.T) {
	var b bytes.Buffer
	l := alog.New(&b, "", 0)
	l.Printf("test: %s", "alog")
	exp := "test: alog\n"
	if b.String() != exp {
		t.Fatalf("unexpected: exp=<%s>; act=<%s>", exp, b.String())
	}
}
func Test_ALog_Level_Check(t *testing.T) {
	var b bytes.Buffer
	l := alog.New(&b, "", 0)
	l.LvEnable(alog.ALL)
	test := func(s string) {
		println(s, "debug:", l.LvIsEnabled(alog.DEBUG))
		println(s, "info: ", l.LvIsEnabled(alog.INFO))
		println(s, "warn: ", l.LvIsEnabled(alog.WARN))
		println(s, "error:", l.LvIsEnabled(alog.ERROR))
		println(s, "fatal:", l.LvIsEnabled(alog.FATAL))
	}
	test("t1")
	l.LvDisable(alog.ERROR | alog.INFO)
	test("t2")
	l.LvEnable(alog.ERROR | alog.INFO)
	test("t3")
}
func Test_ALog_Level(t *testing.T) {
	var b bytes.Buffer
	l := alog.New(&b, "", 0)

	{
		testPrint := func(s string) {
			l.Printfl(alog.DEBUG, "debug: %s", s)
			l.Printfl(alog.INFO, "info: %s", s)
			l.Printfl(alog.WARN, "warn: %s", s)
			l.Printfl(alog.ERROR, "error: %s", s)
			l.Printfl(alog.FATAL, "fatal: %s", s)
		}

		testPrint("t1")
		l.LvDisable(alog.INFO | alog.WARN)

		testPrint("t2")
		l.LvEnable(alog.DEBUG | alog.WARN)
		testPrint("t3")

		exp := "info: t1\nwarn: t1\nerror: t1\nfatal: t1\nerror: t2\nfatal: t2\ndebug: t3\nwarn: t3\nerror: t3\nfatal: t3\n"
		if b.String() != exp {
			t.Fatalf("unexpected: exp=<%s>; act=<%s>", exp, b.String())
		}
	}
	{
		l.LvEnable(alog.WARN | alog.INFO | alog.ERROR | alog.FATAL)
		l.LvDisable(alog.DEBUG)
		b.Reset()
		testPrint := func(s string) {
			l.Printl(alog.DEBUG, "debug: "+s)
			l.Printl(alog.INFO, "info: "+s)
			l.Printl(alog.WARN, "warn: "+s)
			l.Printl(alog.ERROR, "error: "+s)
			l.Printl(alog.FATAL, "fatal: "+s)
		}

		testPrint("t1")
		l.LvDisable(alog.INFO | alog.WARN)
		testPrint("t2")
		l.LvEnable(alog.DEBUG | alog.WARN)
		testPrint("t3")

		exp := "info: t1\nwarn: t1\nerror: t1\nfatal: t1\nerror: t2\nfatal: t2\ndebug: t3\nwarn: t3\nerror: t3\nfatal: t3\n"
		if b.String() != exp {
			t.Fatalf("unexpected: exp=<%s>; act=<%s>", exp, b.String())
		}
	}
}
func Test_ALog_Close(t *testing.T) {
	tmpFile := "./tmp/alog.close.txt"
	// Create file
	{
		out, _ := os.Create(tmpFile)
		l := alog.New(out, "CLOSE|", alog.F_PREFIX)
		l.Print("test")
		l.Print("log")
		l.Close()
		l.Print("XX") // this shouldn't print
		for i := 0; i < 10; i++ {
			l.Print("xx", i)
		}
	}
	// Check file
	{
		fi, err := ioutil.ReadFile(tmpFile)
		if err != nil {
			t.Fatal(err.Error())
		}
		sfi := string(fi)
		sfi = strings.Replace(sfi, "\n", "\\n", -1)
		exp := "CLOSE|test\\nCLOSE|log\\n"

		if sfi != exp {
			t.Fatalf("unexpected: exp=<%s>; act=<%s>", exp, sfi)
		}
	}
}
func Test_ALog_Std(t *testing.T) {
	var b bytes.Buffer
	alog.SetOutput(&b)
	alog.SetFlag(0)
	alog.Printf("test: %s", "alog")

	exp := "test: alog\n"
	if b.String() != exp {
		t.Fatalf("unexpected: exp=<%s>; act=<%s>", exp, b.String())
	}
}

// =====================================================================================================================
// BENCHMARK
// =====================================================================================================================
func Benchmark_ALog_Print(b *testing.B) {
	b.StartTimer()
	out, _ := os.Create("./tmp/alog_print.txt")
	x := alog.New(out, "test ", alog.F_STD)
	for i := 0; i < b.N; i++ {
		x.Print("Print(): ", i, ", an", " ", "a", "w", 3, "s", "o", "m", 3)
	}
	x.Close()
	b.StopTimer()
	b.ReportAllocs()
}
func Benchmark_ALog_Print_Buf(b *testing.B) {
	b.StartTimer()
	out, _ := os.Create("./tmp/alog_print_buf.txt")
	x := alog.New(out, "test ", alog.F_STD|alog.F_USE_BUF_2K)
	for i := 0; i < b.N; i++ {
		x.Print("Print(): ", i, ", an", " ", "a", "w", 3, "s", "o", "m", 3)
	}
	x.Close()
	b.StopTimer()
	b.ReportAllocs()
}
func Benchmark_ALog_Printfl(b *testing.B) {
	b.StartTimer()
	out, _ := os.Create("./tmp/alog_printfl.txt")
	x := alog.New(out, "test ", alog.F_STD|alog.F_MICROSEC|alog.F_DATE)
	for i := 0; i < b.N; i++ {
		x.Printfl(alog.INFO, "sample with %d", i) // fmt.Fprintf() can't be easily optimized.. maybe need to write my own..
	}
	x.Close()
	b.StopTimer()
	b.ReportAllocs()
}
func Benchmark_ALog_Printf(b *testing.B) {
	b.StartTimer()
	out, _ := os.Create("./tmp/alog_printf.txt")
	x := alog.New(out, "test ", alog.F_STD|alog.F_MICROSEC|alog.F_DATE)

	for i := 0; i < b.N; i++ {
		x.Printf("sample with %d", i) // fmt.Fprintf() can't be easily optimized.. maybe need to write my own..
	}
	x.Close()
	b.StopTimer()
	b.ReportAllocs()
}
func Benchmark_ALog_Printf_Buf(b *testing.B) {
	b.StartTimer()
	out, _ := os.Create("./tmp/alog_printf_buf.txt")
	x := alog.New(out, "test ", alog.F_STD|alog.F_MICROSEC|alog.F_DATE|alog.F_USE_BUF_2K)

	for i := 0; i < b.N; i++ {
		x.Printf("sample with %d", i) // fmt.Fprintf() can't be easily optimized.. maybe need to write my own..
	}
	x.Close()
	b.StopTimer()
	b.ReportAllocs()
}
func Benchmark_ALog_Printj(b *testing.B) {
	out, _ := os.Create("./tmp/alog_printj.txt")
	x := alog.New(out, "jsonTest", alog.F_STD)
	a := struct {
		Name  string `json:"name"`
		City  string `json:"city"`
		Count int    `json:"cnt"`
	}{
		Name: "Gon",
		City: "Conway",
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		// a.Count = i
		a.Count = i
		x.Printj("log|", &a)
		// v0.1.1, 96 B/op, 2 allocs/op
		// v0.2.0, 0 B/op,  0 allocs/op
	}
	x.Close()
	b.StopTimer()
	b.ReportAllocs()
}
func Benchmark_ALog_Printj_Buf(b *testing.B) {
	out, _ := os.Create("./tmp/alog_printj_buf.txt")
	x := alog.New(out, "jsonTest", alog.F_STD|alog.F_USE_BUF_2K)
	a := struct {
		Name  string `json:"name"`
		City  string `json:"city"`
		Count int    `json:"cnt"`
	}{
		Name: "Gon",
		City: "Conway",
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		// a.Count = i
		a.Count = i
		x.Printj("log|", &a)
		// v0.1.1, 96 B/op, 2 allocs/op
		// v0.2.0, 0 B/op,  0 allocs/op
	}
	x.Close()
	b.StopTimer()
	b.ReportAllocs()
}

// =====================================================================================================================
// STANDARD BUILT-IN LOGGER
// =====================================================================================================================
// func Benchmark_Builtin_Logger_Printf(b *testing.B) {
// 	out, _ := os.Create("./tmp/builtin.printf.txt")
// 	x := log.New(out, "", log.Lmicroseconds)
//
// 	b.StartTimer()
// 	for i := 0; i < b.N; i++ {
// 		x.Printf("aaa: %d\n", i)
// 	}
// 	b.StopTimer()
// 	b.ReportAllocs()
// }
// func Benchmark_Builtin_Logger_Print(b *testing.B) {
// 	out, _ := os.Create("./tmp/builtin.print.txt")
// 	x := log.New(out, "", log.Lmicroseconds)
//
// 	b.StartTimer()
// 	for i := 0; i < b.N; i++ {
// 		x.Print(i)
// 	}
// 	b.StopTimer()
// 	b.ReportAllocs()
// }
