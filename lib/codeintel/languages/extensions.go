package languages

import (
	"path/filepath"
	"strings"

	"github.com/go-enry/go-enry/v2"
)

// getLanguagesByAlias is a replacement for enry.GetLanguagesByAlias
// It supports languages that are missing in go-enry
func GetLanguageByAlias(alias string) (lang string, ok bool) {
	normalizedAlias := strings.ToLower(alias)
	if lang, ok = unsupportedByEnryAliasMap[normalizedAlias]; ok {
		return lang, true
	}

	return enry.GetLanguageByAlias(normalizedAlias)
}

// getLanguagesByExtension is a replacement for enry.GetLanguagesByExtension
// It supports languages that are missing in go-enry
func GetLanguageExtensions(alias string) []string {
	if lang, ok := unsupportedByEnryNameToExtensionMap[alias]; ok {
		return []string{lang}
	}

	return enry.GetLanguageExtensions(alias)
}

// getLanguagesByExtension is a replacement for enry.GetLanguagesByExtension
// to work around the following limitations:
//   - For some extensions which are overwhelmingly used by a certain file type
//     in practice, such as '.ts', '.md' and '.yaml', it returns ambiguous results.
//   - It does not provide any information about binary files.
//   - Some languages are not supported by enry yet (e.g. Pkl)
func getLanguagesByExtension(path string) (candidates []string, isLikelyBinaryFile bool) {
	ext := filepath.Ext(path)
	if ext == "" {
		return nil, false
	}
	if lang, ok := unsupportedByEnryExtensionToNameMap[ext]; ok {
		return []string{lang}, false
	}
	if _, ok := commonBinaryFileExtensions[ext[1:]]; ok {
		return nil, true
	}
	if lang, ok := overrideAmbiguousExtensionsMap[ext]; ok {
		return []string{lang}, false
	}
	return enry.GetLanguagesByExtension(path, nil, nil), false
}

var commonBinaryFileExtensions = func() map[string]struct{} {
	m := map[string]struct{}{}
	for _, s := range commonBinaryFileExtensionsList {
		m[s] = struct{}{}
	}
	return m
}()

var overrideAmbiguousExtensionsMap = map[string]string{
	// Ignoring the uncommon usage of '.cs' for Smalltalk.
	".cs": "C#",
	// The other languages are Filterscript, Forth, GLSL. Out of that,
	// Forth and GLSL commonly use other extensions. Ignore Filterscript
	// as it is niche.
	".fs": "F#",
	// Not considering "GCC Machine Description".
	".md": "Markdown",
	// The other main language using '.rs' is RenderScript, but that's deprecated.
	// See https://developer.android.com/guide/topics/renderscript/compute
	".rs": "Rust",
	// In i18n contexts, there are XML files with '.ts' and '.tsx' extensions,
	// but we ignore those for now to avoid penalizing the common case.
	".tsx": "TSX",
	".ts":  "TypeScript",
	// Ignoring "Adblock Filter List" and "Vim Help File".
	".txt": "Text",
	// Ignoring other variants of YAML.
	".yaml": "YAML",
	// ".yml" is not included here in parallel to ".yaml"
	// as it is the first extension for 'YAML' and not the first
	// for other variants of YAML, hence only 'YAML' is picked by enry.
}

var unsupportedByEnryExtensionToNameMap = map[string]string{
	// Pkl Configuration Language (https://pkl-lang.org/)
	// NOTE: Add to linguist on 6/7/24
	// can remove once go-enry package updates
	// to that linguist version
	".pkl": "Pkl",
	// Magik Language
	".magik": "Magik",
}

var unsupportedByEnryNameToExtensionMap = reverseMap(unsupportedByEnryExtensionToNameMap)

var unsupportedByEnryAliasMap = map[string]string{
	// Pkl Configuration Language (https://pkl-lang.org/)
	"pkl": "Pkl",
	// Magik Language
	"magik": "Magik",
}

func reverseMap(m map[string]string) map[string]string {
	n := make(map[string]string, len(m))
	for k, v := range m {
		n[v] = k
	}
	return n
}

// Source: https://github.com/sindresorhus/binary-extensions/blob/main/binary-extensions.json
// License: https://github.com/sindresorhus/binary-extensions/blob/main/license
// Replace the contents with
// curl -L https://raw.githubusercontent.com/sindresorhus/binary-extensions/main/binary-extensions.json | jq '.[]' | awk  '{print $1 ","}'
//
// Not adding a leading '.' here to make it easier to update/compare the list.
var commonBinaryFileExtensionsList = []string{
	"3dm",
	"3ds",
	"3g2",
	"3gp",
	"7z",
	"a",
	"aac",
	"adp",
	"ai",
	"aif",
	"aiff",
	"alz",
	"ape",
	"apk",
	"appimage",
	"ar",
	"arj",
	"asf",
	"au",
	"avi",
	"bak",
	"baml",
	"bh",
	"bin",
	"bk",
	"bmp",
	"btif",
	"bz2",
	"bzip2",
	"cab",
	"caf",
	"cgm",
	"class",
	"cmx",
	"cpio",
	"cr2",
	"cur",
	"dat",
	"dcm",
	"deb",
	"dex",
	"djvu",
	"dll",
	"dmg",
	"dng",
	"doc",
	"docm",
	"docx",
	"dot",
	"dotm",
	"dra",
	"DS_Store",
	"dsk",
	"dts",
	"dtshd",
	"dvb",
	"dwg",
	"dxf",
	"ecelp4800",
	"ecelp7470",
	"ecelp9600",
	"egg",
	"eol",
	"eot",
	"epub",
	"exe",
	"f4v",
	"fbs",
	"fh",
	"fla",
	"flac",
	"flatpak",
	"fli",
	"flv",
	"fpx",
	"fst",
	"fvt",
	"g3",
	"gh",
	"gif",
	"graffle",
	"gz",
	"gzip",
	"h261",
	"h263",
	"h264",
	"icns",
	"ico",
	"ief",
	"img",
	"ipa",
	"iso",
	"jar",
	"jpeg",
	"jpg",
	"jpgv",
	"jpm",
	"jxr",
	"key",
	"ktx",
	"lha",
	"lib",
	"lvp",
	"lz",
	"lzh",
	"lzma",
	"lzo",
	"m3u",
	"m4a",
	"m4v",
	"mar",
	"mdi",
	"mht",
	"mid",
	"midi",
	"mj2",
	"mka",
	"mkv",
	"mmr",
	"mng",
	"mobi",
	"mov",
	"movie",
	"mp3",
	"mp4",
	"mp4a",
	"mpeg",
	"mpg",
	"mpga",
	"mxu",
	"nef",
	"npx",
	"numbers",
	"nupkg",
	"o",
	"odp",
	"ods",
	"odt",
	"oga",
	"ogg",
	"ogv",
	"otf",
	"ott",
	"pages",
	"pbm",
	"pcx",
	"pdb",
	"pdf",
	"pea",
	"pgm",
	"pic",
	"png",
	"pnm",
	"pot",
	"potm",
	"potx",
	"ppa",
	"ppam",
	"ppm",
	"pps",
	"ppsm",
	"ppsx",
	"ppt",
	"pptm",
	"pptx",
	"psd",
	"pya",
	"pyc",
	"pyo",
	"pyv",
	"qt",
	"rar",
	"ras",
	"raw",
	"resources",
	"rgb",
	"rip",
	"rlc",
	"rmf",
	"rmvb",
	"rpm",
	"rtf",
	"rz",
	"s3m",
	"s7z",
	"scpt",
	"sgi",
	"shar",
	"snap",
	"sil",
	"sketch",
	"slk",
	"smv",
	"snk",
	"so",
	"stl",
	"suo",
	"sub",
	"swf",
	"tar",
	"tbz",
	"tbz2",
	"tga",
	"tgz",
	"thmx",
	"tif",
	"tiff",
	"tlz",
	"ttc",
	"ttf",
	"txz",
	"udf",
	"uvh",
	"uvi",
	"uvm",
	"uvp",
	"uvs",
	"uvu",
	"viv",
	"vob",
	"war",
	"wav",
	"wax",
	"wbmp",
	"wdp",
	"weba",
	"webm",
	"webp",
	"whl",
	"wim",
	"wm",
	"wma",
	"wmv",
	"wmx",
	"woff",
	"woff2",
	"wrm",
	"wvx",
	"xbm",
	"xif",
	"xla",
	"xlam",
	"xls",
	"xlsb",
	"xlsm",
	"xlsx",
	"xlt",
	"xltm",
	"xltx",
	"xm",
	"xmind",
	"xpi",
	"xpm",
	"xwd",
	"xz",
	"z",
	"zip",
	"zipx",
}
