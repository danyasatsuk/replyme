package replyme

import (
	"embed"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

//go:embed locales/*.toml
var localeFiles embed.FS

var bundle *i18n.Bundle

var active *i18n.Localizer

func loadAllTOMLLocales(bundle *i18n.Bundle, fs embed.FS, tomlPath string) error {
	files, err := fs.ReadDir(tomlPath)
	if err != nil {
		return fmt.Errorf("cannot read embedded locales: %w", err)
	}

	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".toml") {
			continue
		}

		data, err := fs.ReadFile(filepath.Join(tomlPath, f.Name()))
		if err != nil {
			return fmt.Errorf("cannot read locale file %s: %w", f.Name(), err)
		}

		var raw struct {
			Messages []map[string]interface{} `toml:"message"`
		}
		if err := toml.Unmarshal(data, &raw); err != nil {
			return fmt.Errorf("failed to parse %s: %w", f.Name(), err)
		}

		langName := strings.TrimSuffix(f.Name(), filepath.Ext(f.Name()))
		tag, err := language.Parse(langName)
		if err != nil {
			log.Printf("⚠️ Invalid language tag in file %s: %v", f.Name(), err)
			continue
		}

		for _, msg := range raw.Messages {
			id, ok := msg["id"].(string)
			if !ok {
				log.Printf("⚠️ Messages without ID skipped: %+v", msg)
				continue
			}

			translation, _ := msg["translation"].(string)
			description, _ := msg["description"].(string)

			message := &i18n.Message{
				ID:          id,
				Description: description,
				Other:       translation,
			}

			if err := bundle.AddMessages(tag, message); err != nil {
				log.Printf("⚠️ couldn't add message %s: %v", id, err)
			}
		}
	}

	return nil
}

func i18nInit() error {
	bundle = i18n.NewBundle(language.Russian)

	if err := loadAllTOMLLocales(bundle, localeFiles, "locales"); err != nil {
		return err
	}

	langCode := detectSystemLanguageCode()

	active = i18n.NewLocalizer(bundle, langCode)
	return nil
}

func L(id string) string {
	l, err := active.Localize(&i18n.LocalizeConfig{MessageID: id})
	if err != nil {
		return "TRANSLATE NOT FOUND"
	}
	return l
}

func detectSystemLanguageCode() string {
	switch runtime.GOOS {
	case "windows":
		return getFromEnvFallback("en")
	case "darwin":
		lang := detectMacLang()
		if lang != "" {
			return lang
		}
		return getFromEnvFallback("en")
	default:
		return getFromEnvFallback("en")
	}
}

func getFromEnvFallback(fallback string) string {
	for _, env := range []string{"LC_ALL", "LANG", "LANGUAGE"} {
		val := os.Getenv(env)
		if val != "" {
			return normalizeLangCode(val)
		}
	}
	return fallback
}

func detectMacLang() string {
	out, err := exec.Command("defaults", "read", "-g", "AppleLocale").Output()
	if err != nil {
		return ""
	}
	return normalizeLangCode(string(out))
}

func normalizeLangCode(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	if i := strings.IndexAny(s, "_-."); i > 0 {
		return s[:i]
	}
	return s
}
