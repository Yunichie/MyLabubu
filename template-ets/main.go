package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

type SumberDaya struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	Jenis string `json:"jenis"`
	Stok  int    `json:"stok"`
}

type Peminjaman struct {
	ID           int    `json:"id"`
	SumberDayaID int    `json:"sumber_daya_id"`
	Peminjam     string `json:"peminjam"`
	Tanggal      string `json:"tanggal"`
	SlotMulai    int    `json:"slot_mulai"`
	SlotSelesai  int    `json:"slot_selesai"`
	Jumlah       int    `json:"jumlah"`
	Status       string `json:"status"`
}

type sumberDayaInput struct {
	Nama  *string `json:"nama"`
	Jenis *string `json:"jenis"`
	Stok  *int    `json:"stok"`
}

type peminjamanInput struct {
	SumberDayaID *int    `json:"sumber_daya_id"`
	Peminjam     *string `json:"peminjam"`
	Tanggal      *string `json:"tanggal"`
	SlotMulai    *int    `json:"slot_mulai"`
	SlotSelesai  *int    `json:"slot_selesai"`
	Jumlah       *int    `json:"jumlah"`
}

type store struct {
	mu           sync.Mutex
	sumberDaya   []SumberDaya
	peminjaman   []Peminjaman
	nextSDID     int
	nextPinjamID int
}

func NewHandler() http.Handler {
	s := &store{nextSDID: 1, nextPinjamID: 1}
	mux := http.NewServeMux()

	// B1
	mux.HandleFunc("POST /sumber-daya", s.createSumberDaya)
	mux.HandleFunc("GET /sumber-daya", s.listSumberDaya)
	mux.HandleFunc("GET /sumber-daya/{id}", s.getSumberDaya)

	// B2, B3
	mux.HandleFunc("POST /peminjaman", s.createPeminjaman)
	mux.HandleFunc("GET /peminjaman", s.listPeminjaman)
	mux.HandleFunc("GET /peminjaman/{id}", s.getPeminjaman)

	return mux
}

func main() {
	log.Fatal(http.ListenAndServe(":8080", NewHandler()))
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func parseID(r *http.Request) (id int, ok bool) {
	r_id := r.PathValue("id")
	id, err := strconv.Atoi(r_id)
	if err != nil {
		return 0, false
	}
	return id, true
}

func tanggalValid(t string) bool {
	_, err := time.Parse("2006-01-02", t)
	return err == nil
}

func (s *store) createSumberDaya(w http.ResponseWriter, r *http.Request) {
	var in sumberDayaInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "body tidak valid")
		return
	}
	if in.Nama == nil || in.Jenis == nil || in.Stok == nil {
		writeError(w, http.StatusBadRequest, "nama, jenis, dan stok wajib diisi")
		return
	}
	nama := strings.TrimSpace(*in.Nama)
	if nama == "" {
		writeError(w, http.StatusBadRequest, "Nama tidak boleh kosong")
		return
	}
	if utf8.RuneCountInString(nama) > 50 {
		writeError(w, http.StatusBadRequest, "Nama tidak boleh lebih dari 50 karakter")
		return
	}
	jenis := strings.TrimSpace(*in.Jenis)
	if jenis != "ruang" && jenis != "alat" {
		writeError(w, http.StatusBadRequest, "Jenis hanya boleh alat atau ruang")
		return
	}
	stok := *in.Stok
	if stok < 1 {
		writeError(w, http.StatusBadRequest, "Stok harus minimal 1")
		return
	}

	sd := SumberDaya{Nama: nama, Jenis: jenis, Stok: stok}
	s.mu.Lock()
	defer s.mu.Unlock()

	sd.ID = s.nextSDID
	s.nextSDID++
	s.sumberDaya = append(s.sumberDaya, sd)
	writeJSON(w, http.StatusCreated, sd)
}

func (s *store) listSumberDaya(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	hasil := []SumberDaya{}
	hasil = append(hasil, s.sumberDaya...)
	writeJSON(w, http.StatusOK, hasil)
}

func (s *store) getSumberDaya(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "id harus angka")
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, v := range s.sumberDaya {
		if v.ID == id {
			writeJSON(w, http.StatusOK, v)
			return
		}
	}
	writeError(w, http.StatusNotFound, "sumber daya tidak ditemukan")
}

func (s *store) createPeminjaman(w http.ResponseWriter, r *http.Request) {
	var in peminjamanInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "body tidak valid")
		return
	}
	if in.SumberDayaID == nil || in.Peminjam == nil || in.Tanggal == nil ||
		in.SlotMulai == nil || in.SlotSelesai == nil || in.Jumlah == nil {
		writeError(w, http.StatusBadRequest, "semua field wajib diisi")
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	peminjam := strings.TrimSpace(*in.Peminjam)
	if peminjam == "" {
		writeError(w, http.StatusBadRequest, "Peminjam kosong")
		return
	}
	if utf8.RuneCountInString(peminjam) > 50 {
		writeError(w, http.StatusBadRequest, "Peminjam maksimal 50 karakter")
		return
	}
	tanggal := strings.TrimSpace(*in.Tanggal)
	if !tanggalValid(tanggal) {
		writeError(w, http.StatusBadRequest, "Tanggal tidak valid")
		return
	}
	slot_mulai := *in.SlotMulai
	slot_selesai := *in.SlotSelesai
	if slot_mulai < 1 || slot_mulai > 8 || slot_selesai < 1 || slot_selesai > 8 || slot_mulai > slot_selesai {
		writeError(w, http.StatusBadRequest, "Slot tidak valid")
		return
	}
	jumlah := *in.Jumlah
	if jumlah < 1 {
		writeError(w, http.StatusBadRequest, "Jumlah tidak valid")
		return
	}

	stok := -1
	for _, sd := range s.sumberDaya {
		if sd.ID == *in.SumberDayaID {
			stok = sd.Stok
			break
		}
	}
	if stok == -1 {
		writeError(w, http.StatusBadRequest, "sumber daya tidak ditemukan")
		return
	}
	if jumlah > stok {
		writeError(w, http.StatusBadRequest, "jumlah melebihi stok")
		return
	}

	p := Peminjaman{
		SumberDayaID: *in.SumberDayaID, Peminjam: peminjam, Tanggal: tanggal,
		SlotMulai: slot_mulai, SlotSelesai: slot_selesai, Jumlah: jumlah, Status: "aktif",
	}
	if s.adaBentrok(p, stok, 0) {
		writeError(w, http.StatusConflict, "slot sudah penuh")
		return
	}

	p.ID = s.nextPinjamID
	s.nextPinjamID++
	s.peminjaman = append(s.peminjaman, p)
	writeJSON(w, http.StatusCreated, p)
}

func (s *store) listPeminjaman(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filterSD := q.Get("sumber_daya_id")
	filterTanggal := q.Get("tanggal")
	filterStatus := q.Get("status")

	sdID := 0
	if filterSD != "" {
		n, err := strconv.Atoi(filterSD)
		if err != nil {
			writeError(w, http.StatusBadRequest, "sumber_daya_id harus angka")
			return
		}
		sdID = n
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	hasil := []Peminjaman{}
	for _, p := range s.peminjaman {
		if filterSD != "" && p.SumberDayaID != sdID {
			continue
		}
		if filterTanggal != "" && p.Tanggal != filterTanggal {
			continue
		}
		if filterStatus != "" && p.Status != filterStatus {
			continue
		}
		hasil = append(hasil, p)
	}
	writeJSON(w, http.StatusOK, hasil)
}

func (s *store) getPeminjaman(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "id harus angka")
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, v := range s.peminjaman {
		if v.ID == id {
			writeJSON(w, http.StatusOK, v)
			return
		}
	}
	writeError(w, http.StatusNotFound, "peminjaman tidak ditemukan")
}

func (s *store) adaBentrok(p Peminjaman, stok int, abaikanID int) bool {
	for slot := p.SlotMulai; slot <= p.SlotSelesai; slot++ {
		terpakai := 0
		for _, lain := range s.peminjaman {
			if lain.ID == abaikanID || lain.Status != "aktif" {
				continue
			}
			if lain.SumberDayaID == p.SumberDayaID && lain.Tanggal == p.Tanggal &&
				lain.SlotMulai <= slot && slot <= lain.SlotSelesai {
				terpakai += lain.Jumlah
			}
		}
		if terpakai+p.Jumlah > stok {
			return true
		}
	}
	return false
}
