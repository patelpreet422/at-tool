package client

var Langs = map[string]string{
	"4001": "C (GCC 9.2.1)",
	"4002": "C (Clang 10.0.0)",
	"4003": "C++ (GCC 9.2.1)",
	"4004": "C++ (Clang 10.0.0)",
	"4005": "Java (OpenJDK 11.0.6)",
	"4006": "Python (3.8.2)",
	"4007": "Bash (5.0.11)",
	"4008": "bc (1.07.1)",
	"4009": "Awk (GNU Awk 4.1.4)",
	"4010": "C# (.NET Core 3.1.201)",
	"4011": "C# (Mono-mcs 6.8.0.105)",
	"4012": "C# (Mono-csc 3.5.0)",
	"4013": "Clojure (1.10.1.536)",
	"4014": "Crystal (0.33.0)",
	"4015": "D (DMD 2.091.0)",
	"4016": "D (GDC 9.2.1)",
	"4017": "D (LDC 1.20.1)",
	"4018": "Dart (2.7.2)",
	"4019": "dc (1.4.1)",
	"4020": "Erlang (22.3)",
	"4021": "Elixir (1.10.2)",
	"4022": "F# (.NET Core 3.1.201)",
	"4023": "F# (Mono 10.2.3)",
	"4024": "Forth (gforth 0.7.3)",
	"4025": "Fortran (GNU Fortran 9.2.1)",
	"4026": "Go (1.14.1)",
	"4027": "Haskell (GHC 8.8.3)",
	"4028": "Haxe (4.0.3); js",
	"4029": "Haxe (4.0.3); Java",
	"4030": "JavaScript (Node.js 12.16.1)",
	"4031": "Julia (1.4.0)",
	"4032": "Kotlin (1.3.71)",
	"4033": "Lua (Lua 5.3.5)",
	"4034": "Lua (LuaJIT 2.1.0)",
	"4035": "Dash (0.5.8)",
	"4036": "Nim (1.0.6)",
	"4037": "Objective-C (Clang 10.0.0)",
	"4038": "Common Lisp (SBCL 2.0.3)",
	"4039": "OCaml (4.10.0)",
	"4040": "Octave (5.2.0)",
	"4041": "Pascal (FPC 3.0.4)",
	"4042": "Perl (5.26.1)",
	"4043": "Raku (Rakudo 2020.02.1)",
	"4044": "PHP (7.4.4)",
	"4045": "Prolog (SWI-Prolog 8.0.3)",
	"4046": "PyPy2 (7.3.0)",
	"4047": "PyPy3 (7.3.0)",
	"4048": "Racket (7.6)",
	"4049": "Ruby (2.7.1)",
	"4050": "Rust (1.42.0)",
	"4051": "Scala (2.13.1)",
	"4052": "Java (OpenJDK 1.8.0)",
	"4053": "Scheme (Gauche 0.9.9)",
	"4054": "Standard ML (MLton 20130715)",
	"4055": "Swift (5.2.1)",
	"4056": "Text (cat 8.28)",
	"4057": "TypeScript (3.8)",
	"4058": "Visual Basic (.NET Core 3.1.101)",
	"4059": "Zsh (5.4.2)",
	"4060": "COBOL - Fixed (OpenCOBOL 1.1.0)",
	"4061": "COBOL - Free (OpenCOBOL 1.1.0)",
	"4062": "Brainfuck (bf 20041219)",
	"4063": "Ada2012 (GNAT 9.2.1)",
	"4064": "Unlambda (2.0.0)",
	"4065": "Cython (0.29.16)",
	"4066": "Sed (4.4)",
	"4067": "Vim (8.2.0460)",
}

// LangsExt language's ext
// TODO: Add extension for all the languages supported on atcoder below is copied as is from cf-tool project
var LangsExt = map[string]string{
	"GNU C11":               "c",
	"Clang++17 Diagnostics": "cpp",
	"GNU C++0x":             "cpp",
	"GNU C++":               "cpp",
	"GNU C++11":             "cpp",
	"GNU C++14":             "cpp",
	"GNU C++17":             "cpp",
	"MS C++":                "cpp",
	"MS C++ 2017":           "cpp",
	"Mono C#":               "cs",
	"D":                     "d",
	"Go":                    "go",
	"Haskell":               "hs",
	"Kotlin":                "kt",
	"Ocaml":                 "ml",
	"Delphi":                "pas",
	"FPC":                   "pas",
	"PascalABC.NET":         "pas",
	"Perl":                  "pl",
	"PHP":                   "php",
	"Python 2":              "py",
	"Python 3":              "py",
	"PyPy 2":                "py",
	"PyPy 3":                "py",
	"Ruby":                  "rb",
	"Rust":                  "rs",
	"JavaScript":            "js",
	"Node.js":               "js",
	"Q#":                    "qs",
	"Java":                  "java",
	"Java 6":                "java",
	"Java 7":                "java",
	"Java 8":                "java",
	"Java 9":                "java",
	"Java 10":               "java",
	"Java 11":               "java",
	"Tcl":                   "tcl",
	"F#":                    "fs",
	"Befunge":               "bf",
	"Pike":                  "pike",
	"Io":                    "io",
	"Factor":                "factor",
	"Cobol":                 "cbl",
	"Secret_171":            "secret_171",
	"Ada":                   "adb",
	"FALSE":                 "f",
	"":                      "txt",
}