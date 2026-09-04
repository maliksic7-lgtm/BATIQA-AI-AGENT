// Command eval runs the AI pipeline against 30 guest utterances and writes
// docs/AI_EVAL_REPORT.md with real accuracy numbers (used in the pitch deck).
package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"batiqa-ai/internal/service/ai"
)

type tc struct {
	msg        string
	wantIntent string
	wantAction string
}

var cases = []tc{
	// Indonesian - hotel info
	{"jam sarapan sampai jam berapa?", ai.IntentBreakfastInformation, ai.ActionAnswer},
	{"wifi passwordnya apa?", ai.IntentWifiInformation, ai.ActionAnswer},
	{"check out jam berapa ya?", ai.IntentCheckoutInformation, ai.ActionAnswer},
	{"kolam renang buka jam berapa?", ai.IntentFacilityInformation, ai.ActionAnswer},
	{"ada gym nggak?", ai.IntentFacilityInformation, ai.ActionAnswer},
	{"restoran hotel buka jam berapa?", ai.IntentRestaurantInformation, ai.ActionAnswer},
	{"boleh bawa hewan peliharaan?", ai.IntentHotelPolicy, ai.ActionAnswer},
	// Indonesian - housekeeping
	{"tolong antar 2 handuk ke kamar", ai.IntentTowelRequest, ai.ActionCreateTicket},
	{"kamar saya belum dibersihkan", ai.IntentRoomCleaningRequest, ai.ActionCreateTicket},
	{"minta tambahan sabun sama shampoo dong", ai.IntentAmenityRequest, ai.ActionCreateTicket},
	{"besok tolong bersihkan kamar pagi ya", ai.IntentRoomCleaningRequest, ai.ActionCreateTicket},
	// Indonesian - engineering
	{"ac kamar ga dingin panas banget", ai.IntentACProblem, ai.ActionCreateTicket},
	{"AC rusak total mati!", ai.IntentACProblem, ai.ActionCreateTicket},
	{"tv nya ga nyala nih", ai.IntentTVProblem, ai.ActionCreateTicket},
	{"wifi sering putus putus", ai.IntentWifiProblem, ai.ActionCreateTicket},
	{"lampu kamar mati semua", ai.IntentLightProblem, ai.ActionCreateTicket},
	{"shower airnya kecil sekali", ai.IntentShowerProblem, ai.ActionCreateTicket},
	{"kloset mampet tolong", ai.IntentPlumbingProblem, ai.ActionCreateTicket},
	// Slang / informal / typo
	{"b0s antar handuk dunk 2 biji", ai.IntentTowelRequest, ai.ActionCreateTicket},
	{"ac q koq gak dingin yaa", ai.IntentACProblem, ai.ActionCreateTicket},
	{"mau tanya, mandi airnya anget ga ya?", ai.IntentShowerProblem, ai.ActionCreateTicket},
	// Recommendations
	{"rekomendasi tempat makan enak budget 50 ribu", ai.IntentRestaurantRecommendation, ai.ActionAnswer},
	{"ada mall terdekat?", ai.IntentShoppingRecommendation, ai.ActionAnswer},
	{"tempat wisata sekitar sini apa aja?", ai.IntentTourismRecommendation, ai.ActionAnswer},
	// English
	{"What time is breakfast?", ai.IntentBreakfastInformation, ai.ActionAnswer},
	{"The AC is not working at all", ai.IntentACProblem, ai.ActionCreateTicket},
	{"Could I get two more towels please?", ai.IntentTowelRequest, ai.ActionCreateTicket},
	{"Any nice restaurant nearby?", ai.IntentRestaurantRecommendation, ai.ActionAnswer},
	// General
	{"halo, selamat pagi!", ai.IntentGreeting, ai.ActionAnswer},
	{"makasih banyak ya", ai.IntentThankYou, ai.ActionAnswer},
}

func main() {
	svc := ai.NewService()

	pass := 0
	var b strings.Builder
	b.WriteString("# Laporan Evaluasi AI — BATIQA Assistant\n\n")
	b.WriteString(fmt.Sprintf("> Dibuat otomatis: %s | Provider aktif: **%s** | Model: `%s`\n\n",
		time.Now().Format("2006-01-02 15:04"), providerName(), strings.TrimSpace(os.Getenv("GEMINI_MODEL"))))

	b.WriteString("| # | Ucapan Tamu | Intent Diharapkan | Hasil | Action | Status |\n|---|---|---|---|---|---|\n")
	for i, c := range cases {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		res, err := svc.Process(ctx, ai.Request{SessionID: "eval", RoomNumber: "305", Message: c.msg})
		cancel()
		gotIntent, gotAction := "-", "-"
		if err == nil && res != nil {
			gotIntent = res.Intent
			gotAction = res.Action.Type
		}
		ok := gotIntent == c.wantIntent && gotAction == c.wantAction
		if ok {
			pass++
		}
		status := "PASS"
		if !ok {
			status = fmt.Sprintf("FAIL (dapat %s/%s)", gotIntent, gotAction)
		}
		b.WriteString(fmt.Sprintf("| %d | %s | %s | %s | %s | %s |\n",
			i+1, c.msg, c.wantIntent, gotIntent, gotAction, status))
	}

	pct := float64(pass) / float64(len(cases)) * 100
	b.WriteString("\n## Ringkasan\n\n")
	fmt.Fprintf(&b, "- **Akurasi gabungan (intent+action): %.1f%%** (%d/%d)\n", pct, pass, len(cases))
	b.WriteString("- Metodologi: setiap ucapan diproses end-to-end melalui pipeline penuh\n")
	b.WriteString("  (language → intent → entity → routing → action) dengan validasi backend,\n")
	b.WriteString("  mencakup bahasa formal/santai/typo, ID & EN.\n")

	os.WriteFile("docs/AI_EVAL_REPORT.md", []byte(b.String()), 0644)
	fmt.Printf("Eval done: %d/%d (%.1f%%) — laporan tersimpan di docs/AI_EVAL_REPORT.md\n", pass, len(cases), pct)
}

func providerName() string {
	if os.Getenv("AI_PROVIDER") == "mock" || (os.Getenv("AI_PROVIDER") == "" && os.Getenv("GEMINI_API_KEY") == "") {
		return "rule-based mock"
	}
	return "Gemini"
}
