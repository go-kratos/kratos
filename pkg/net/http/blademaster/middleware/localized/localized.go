package localized

import (
	"net/url"
	"sort"
	"strconv"
	"strings"

	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
)

type language struct {
	Name   string
	Weight float64
}

// detectLocalizedWeb detect localized from Accept-Language Header.
func detectLocalizedWeb(ctx *bm.Context) (langs []string) {
	parseLang := func(s string) (language, error) {
		seps := strings.SplitN(s, ";", 2)
		lang := language{Name: seps[0]}
		if len(seps) == 1 {
			lang.Weight = 1
			return lang, nil
		}
		params, err := url.ParseQuery(seps[1])
		if err != nil {
			return lang, err
		}
		lang.Weight, err = strconv.ParseFloat(params.Get("q"), 32)
		return lang, err
	}
	items := strings.Split(ctx.Request.Header.Get("Accept-Language"), ",")
	if items[0] == "" {
		return
	}
	// three language is most common accept language send by browser
	languages := make([]language, 0, len(items))
	for _, s := range items {
		l, err := parseLang(s)
		if err != nil {
			continue
		}
		languages = append(languages, l)
	}
	sort.Slice(languages, func(i, j int) bool {
		return languages[i].Weight > languages[j].Weight
	})
	for i := range languages {
		langs = append(langs, languages[i].Name)
	}
	return
}

// detectLocalizediOS iOS send locale info in form format as {language}_{region},
// e.g. en_us zh_cn zh_tw zh_mo zh_hk ja_jp.
//
// the language used following the specification defined at https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes
//
// the region used following the specification default at https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2
// NOTE: region use lower case not a upper case in ISO_3166-1
func detectLocalizediOS(ctx *bm.Context) (langs []string) {
	if err := ctx.Request.ParseForm(); err != nil {
		return
	}
	locale := ctx.Request.Form.Get("locale")
	if locale == "" {
		return
	}
	seqs := strings.SplitN(locale, "_", 2)
	if len(seqs) != 2 {
		// ignored invalid locale
		return
	}
	language := seqs[0]
	region := strings.ToUpper(seqs[1])
	return []string{language + "-" + region}
}

// detectLocalizedAndroid Android sed locale info in form as same as iOS
// format as Android Code: `Locale.getDefault().getLanguage+"_"+Locale.getDefault().getCountry`
// see: https://developer.android.com/reference/java/util/Locale
// e.g. en_US zh_CN
//
// getCountry: Returns the country/region code for this locale, which should either be the empty string,
// an uppercase ISO 3166 2-letter code, or a UN M.49 3-digit code.
//
// getLanguage: Returns the language code of this Locale following ISO 639.
func detectLocalizedAndroid(ctx *bm.Context) (langs []string) {
	if err := ctx.Request.ParseForm(); err != nil {
		return
	}
	locale := ctx.Request.Form.Get("locale")
	if locale == "" {
		return
	}
	// simple replace '_' to '-' convert to bcp47 format
	return []string{strings.Replace(locale, "_", "-", -1)}
}

// Localized detect locale from http header
// see: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Accept-Language
func Localized(ctx *bm.Context) {
	detectFn := detectLocalizedWeb
	val := metadata.Value(ctx, metadata.Device)
	if val != nil {
		if device, ok := val.(*bm.Device); ok && device.RawMobiApp != "" {
			switch {
			case device.IsIOS():
				detectFn = detectLocalizediOS
			case device.IsAndroid():
				detectFn = detectLocalizedAndroid
			}
		}
	}
	md, ok := metadata.FromContext(ctx)
	if ok {
		md[metadata.Locale] = detectFn(ctx)
	}
}
