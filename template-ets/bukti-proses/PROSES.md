# PROSES

## 1. Persona dan masalahnya
Persona D: panitia acara. Satu acara butuh beberapa sumber daya sekaligus, misalnya aula, proyektor, dan sound system; kalau salah satunya tidak tersedia, acaranya batal, dan panitia tidak mau ruangnya sudah terpesan tetapi alatnya tidak ada. Fitur khasnya adalah peminjaman paket: semua sumber daya dipesan dalam satu permintaan, dan hasilnya berhasil semua atau gagal semua.

## 2. Rancangan struktur data
- `SumberDaya` (id, nama, jenis, stok) dan `Peminjaman` (id, sumber_daya_id, peminjam, tanggal, slot_mulai, slot_selesai, jumlah, status), dengan tag JSON yang persis sama dengan kontrak.
- `sumberDayaInput` dan `peminjamanInput` untuk body request. Field-nya pointer (`*string`, `*int`) supaya bisa membedakan field yang tidak dikirim (`nil`) dengan nilai `0` atau `""`, karena semua field wajib ada.
- `store` berisi:
  - `sumberDaya []SumberDaya` dan `peminjaman []Peminjaman`: slice, karena daftar harus selalu terurut menurut id naik, dan dengan `append` urutannya otomatis terjaga.
  - `nextSDID` dan `nextPinjamID`: counter id, dimulai dari 1 dan bertambah satu per data baru.
  - `mu sync.Mutex`: satu kunci untuk semua data bersama.
- Setiap `NewHandler()` membuat `store` baru yang kosong.
- Perubahan di tengah jalan: awalnya validasi peminjaman ditulis langsung di handler `createPeminjaman`, lalu saya pindahkan ke fungsi `validasiPeminjaman` supaya bisa dipakai juga oleh PUT.
- Untuk persona D direncanakan field `paket_id` di `Peminjaman` dan counter `nextPaketID` di `store`.

## 3. Satu keputusan desain
Cek bentrok (`adaBentrok`) dilakukan per slot, bukan hanya membandingkan rentang. Untuk setiap slot dari `slot_mulai` sampai `slot_selesai`, saya jumlahkan unit yang dipakai peminjaman berstatus aktif dengan sumber daya dan tanggal yang sama, lalu tolak (409) kalau jumlah itu ditambah jumlah baru melebihi stok.

Pilihan lain yang saya pertimbangkan adalah cukup mengecek apakah rentang beririsan dengan peminjaman lain. Cara itu salah untuk alat yang stoknya lebih dari 1, karena dua peminjaman boleh beririsan selama total unitnya tidak melebihi stok. Dengan cek per slot, irisan sebagian (misalnya 3-5 dan 4-4) dan rentang yang hanya bersebelahan (1-2 dan 3-4, tidak bentrok) otomatis tertangani dengan benar.

## 4. Bug atau kendala (minimal dua)
### Bug 1
- Gejala: `POST /sumber-daya` dengan `jenis` "kursi" atau `stok` 0 membalas 400, tetapi sumber daya itu tetap tersimpan dan muncul di `GET /sumber-daya`.
- Cara menemukan: menguji `POST /sumber-daya` dengan data tidak valid, lalu memanggil `GET /sumber-daya` dan melihat data yang seharusnya ditolak ikut tersimpan.
- Penyebab dan solusi: saya lupa menulis `return` setelah `writeError()`. `writeError()` hanya menulis respons error, tidak menghentikan fungsi, sehingga kode tetap lanjut sampai `append` ke slice. Solusinya, menambahkan `return` setelah setiap `writeError()`.

### Bug 2
- Gejala: nama sumber daya berisi 30 huruf `é` ditolak dengan pesan "Nama tidak boleh lebih dari 50 karakter", padahal hanya 30 karakter.
- Cara menemukan: menguji `POST /sumber-daya` dengan nama yang memakai huruf non-ASCII.
- Penyebab dan solusi: saya memakai `len(nama)`, yang menghitung byte, bukan karakter. Huruf `é` butuh 2 byte di UTF-8, jadi 30 karakter dihitung 60. Solusinya, memakai `utf8.RuneCountInString()` (package `unicode/utf8`) yang menghitung per karakter (rune), untuk `nama` dan juga `peminjam`.

## 5. Change request
Belum selesai
