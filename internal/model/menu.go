package model

// MenuItem is a food or drink offered for room-service ordering. Price is in
// Indonesian Rupiah (integer). Category is FOOD or DRINK.
type MenuItem struct {
	Name     string `json:"name"`
	Category string `json:"category"` // FOOD | DRINK
	Price    int    `json:"price"`
}

// RoomServiceMenu is the FRESQA Bistro room-service catalog served to guests
// and used to validate restaurant orders. DEMO pricing — verify before production.
var RoomServiceMenu = []MenuItem{
	// Makanan (FOOD)
	{"Gulai Ikan Patin", "FOOD", 95000},
	{"Ayam Panggang Madu", "FOOD", 85000},
	{"Sate Padang", "FOOD", 65000},
	{"Nasi Goreng", "FOOD", 55000},
	{"Nasi Uduk", "FOOD", 45000},
	{"Bubur Ayam", "FOOD", 40000},
	{"Omelet Telur", "FOOD", 35000},
	{"Roti Bakar", "FOOD", 30000},
	// Minuman (DRINK)
	{"Kopi", "DRINK", 25000},
	{"Teh", "DRINK", 20000},
	{"Jus Buah", "DRINK", 30000},
	{"Air Mineral", "DRINK", 15000},
}

// IsValidMenuItem reports whether a name exists in the room-service menu.
func IsValidMenuItem(name string) bool {
	for _, m := range RoomServiceMenu {
		if m.Name == name {
			return true
		}
	}
	return false
}

// MenuItemByName returns the catalog entry for a name (ok=false if unknown).
func MenuItemByName(name string) (MenuItem, bool) {
	for _, m := range RoomServiceMenu {
		if m.Name == name {
			return m, true
		}
	}
	return MenuItem{}, false
}
