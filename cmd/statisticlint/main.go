package main

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/appends"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/defers"
	"golang.org/x/tools/go/analysis/passes/directive"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/fieldalignment"
	"golang.org/x/tools/go/analysis/passes/findcall"
	"golang.org/x/tools/go/analysis/passes/framepointer"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/pkgfact"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/reflectvaluecompare"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/slog"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stdversion"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/timeformat"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"golang.org/x/tools/go/analysis/passes/usesgenerics"
	"honnef.co/go/tools/staticcheck"

	"github.com/rodeorm/shortener/pkg/osexitchecker"
)

/*
Statisticlint это инструмент для статистического анализа программы
*/
func main() {

	var checks []*analysis.Analyzer

	checks = addStandardAnalyzers(checks)
	checks = addStaticCheckAnalyzersSA(checks)
	checks = addStaticCheckAnalyzersOther(checks)

	checks = append(checks, osexitchecker.Analyzer)

	if len(checks) != 0 {
		multichecker.Main(
			checks...,
		)
	}
}

// addStaticCheckAnalyzers добавляет остальные анализаторы (кроме класса SA) пакета staticcheck.io
func addStaticCheckAnalyzersOther(analyzers []*analysis.Analyzer) []*analysis.Analyzer {

	for _, a := range staticcheck.Analyzers {
		switch a.Analyzer.Name {
		case "S1000":
			analyzers = append(analyzers, a.Analyzer)
		case "ST1000":
			analyzers = append(analyzers, a.Analyzer)
		case "QF1001":
			analyzers = append(analyzers, a.Analyzer)
		}
	}
	return analyzers
}

// addStaticCheckAnalyzersSA добавляет все анализаторы класса SA пакета staticcheck.io
func addStaticCheckAnalyzersSA(analyzers []*analysis.Analyzer) []*analysis.Analyzer {

	for _, a := range staticcheck.Analyzers {
		if a.Analyzer.Name[0] == 'S' && a.Analyzer.Name[1] == 'A' {
			analyzers = append(analyzers, a.Analyzer)
		}
	}
	return analyzers
}

// addStandardAnalyzers добавляет стандартные статические анализаторы пакета golang.org/x/tools/go/analysis/passes;
func addStandardAnalyzers(analyzers []*analysis.Analyzer) []*analysis.Analyzer {
	analyzers = append(analyzers, appends.Analyzer)             //  обнаруживает, если в append есть только одна переменная.
	analyzers = append(analyzers, asmdecl.Analyzer)             //  сообщает о несоответствиях между ассемблерными файлами и объявлениями Go.
	analyzers = append(analyzers, assign.Analyzer)              //  обнаруживает бесполезные присваивания.
	analyzers = append(analyzers, atomic.Analyzer)              //  проверяет общие ошибки при использовании пакета sync/atomic.
	analyzers = append(analyzers, atomicalign.Analyzer)         //  проверяет неповоротные аргументы для функций sync/atomic.
	analyzers = append(analyzers, bools.Analyzer)               //  обнаруживает распространенные ошибки с логическими операторами.
	analyzers = append(analyzers, buildssa.Analyzer)            //  строит представление SSA для бесошибочного пакета и возвращает совокупность всех функций внутри него.
	analyzers = append(analyzers, buildtag.Analyzer)            //  проверяет теги сборки.
	analyzers = append(analyzers, cgocall.Analyzer)             //  обнаруживает некоторые нарушения правил передачи указателей cgo.
	analyzers = append(analyzers, composite.Analyzer)           //  проверяет неключевые составные литералы.
	analyzers = append(analyzers, copylock.Analyzer)            //  проверяет, были ли блокировки ошибочно переданы по значению.
	analyzers = append(analyzers, ctrlflow.Analyzer)            //  предоставляет синтаксическую контрольную-flow граф (CFG) для тела функции.
	analyzers = append(analyzers, deepequalerrors.Analyzer)     //  проверяет использование reflect.DeepEqual с ошибочными значениями.
	analyzers = append(analyzers, defers.Analyzer)              //  проверяет распространенные ошибки в defer-выражениях.
	analyzers = append(analyzers, directive.Analyzer)           //  проверяет известные директивы инструментов Go.
	analyzers = append(analyzers, errorsas.Analyzer)            //  проверяет, что второй аргумент для errors.As является указателем на тип, реализующий интерфейс error.
	analyzers = append(analyzers, fieldalignment.Analyzer)      //  обнаруживает структуры, которые будут использовать меньше памяти, если их поля будут отсортированы.
	analyzers = append(analyzers, findcall.Analyzer)            //  служит тривиальным примером и тестом API анализа.
	analyzers = append(analyzers, framepointer.Analyzer)        //  сообщает о коде ассемблера, который испортит указатель кадра до его сохранения.
	analyzers = append(analyzers, httpresponse.Analyzer)        //  проверяет наличие ошибок при использовании HTTP-ответов.
	analyzers = append(analyzers, ifaceassert.Analyzer)         //  помечает невозможные утверждения типов интерфейсов.
	analyzers = append(analyzers, inspect.Analyzer)             //  предоставляет инспектор AST для синтаксических деревьев пакета.
	analyzers = append(analyzers, loopclosure.Analyzer)         //  проверяет ссылки на переменные внешнего цикла из вложенных функций.
	analyzers = append(analyzers, lostcancel.Analyzer)          //  проверяет неиспользование функции отмены контекста.
	analyzers = append(analyzers, nilfunc.Analyzer)             //  проверяет бесполезные сравнения с nil.
	analyzers = append(analyzers, nilness.Analyzer)             //  проверяет граф управления потоком функции SSA и сообщает об ошибках, таких как разыменование нулевых указателей и неэффективные нулевые сравнения.
	analyzers = append(analyzers, pkgfact.Analyzer)             //  является демонстрацией и тестом механизма фактов пакета.
	analyzers = append(analyzers, printf.Analyzer)              //  проверяет совместимость форматированных строк и аргументов Printf.
	analyzers = append(analyzers, reflectvaluecompare.Analyzer) //  проверяет на случайное использование == или reflect.DeepEqual для сравнения значений reflect.Value.
	analyzers = append(analyzers, shadow.Analyzer)              //  проверяет затененные переменные.
	analyzers = append(analyzers, shift.Analyzer)               //  проверяет смещения, превышающие ширину целого числа.
	analyzers = append(analyzers, sigchanyzer.Analyzer)         //  обнаруживает неправильное использование небуферизованного сигнала как аргумента для signal.Notify.
	analyzers = append(analyzers, slog.Analyzer)                //  проверяет несоответствия между парами ключ-значение в вызовах log/slog.
	analyzers = append(analyzers, sortslice.Analyzer)           //  проверяет вызовы sort.Slice, которые не используют тип среза в качестве первого аргумента.
	analyzers = append(analyzers, stdmethods.Analyzer)          //  проверяет опечатки в сигнатурах методов, аналогичных известным интерфейсам.
	analyzers = append(analyzers, stdversion.Analyzer)          //  сообщает о использовании символов стандартной библиотеки, которые являются "слишком новыми" для версии Go, действующей в ссылающемся файле.
	analyzers = append(analyzers, stringintconv.Analyzer)       //  помечает преобразования типов из целых чисел в строки.
	analyzers = append(analyzers, structtag.Analyzer)           //  проверяет корректность тегов полей структур.
	analyzers = append(analyzers, testinggoroutine.Analyzer)    //  обнаруживает вызовы Fatal из тестовой горутины.
	analyzers = append(analyzers, tests.Analyzer)               //  проверяет распространенные ошибки использования тестов и примеров.
	analyzers = append(analyzers, timeformat.Analyzer)          //  проверяет вызовы time.Format или time.Parse с неверным форматом.
	analyzers = append(analyzers, unmarshal.Analyzer)           //  проверяет передачу типов, не являющихся указателями или интерфейсами, в функции unmarshal и decode.
	analyzers = append(analyzers, unreachable.Analyzer)         //  проверяет наличие недостижимого кода.
	analyzers = append(analyzers, unsafeptr.Analyzer)           //  проверяет недопустимые преобразования uintptr в unsafe.Pointer.
	analyzers = append(analyzers, unusedresult.Analyzer)        //  проверяет неиспользуемые результаты вызовов некоторых чистых функций.
	analyzers = append(analyzers, unusedwrite.Analyzer)         //  проверяет на наличие неиспользуемых записей в элементах структуры или массива.
	analyzers = append(analyzers, usesgenerics.Analyzer)        //  проверяет использование возможностей обобщений, добавленных в Go 1.18.
	return analyzers
}
