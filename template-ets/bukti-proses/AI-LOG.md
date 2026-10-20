# AI-LOG

Catat setiap penggunaan AI yang menghasilkan kode. Jika tidak memakai AI sama sekali, tulis satu baris: **tidak memakai AI**.

Log yang jujur tidak mengurangi nilai, sebanyak apa pun AI dipakai. Yang mengurangi nilai adalah log yang tidak sesuai dengan kode, atau isi log yang tidak bisa dijelaskan.

## tanggalValid()
- Prompt: bagaimana cara memvalidasi tanggal sesuai requirement?
- Dipakai: fungsi tanggalValid() di line 95-98
- Diubah: di line 96, variable `time, err` saya ubah jadi `_, err` karena variable `time` tidak dipakai.

## Membuat skeleton/placeholder untuk fungsi yang harus diimplementasikan
- Prompt: Buatkan placeholder untuk fungsi yang harus diimplementasikan sesuai requirement
- Dipakai: fungsi B1, B2, dan B3
- Diubah: tidak ada. Hanya mengisi todo.

## Meminta contoh implementasi locking
- Prompt: berikan contoh implementasi lock dan unlock
- Dipakai: hasil kode contoh di B1
- Diubah: tidak ada

## Implementasi query param
- Prompt: Berikan contoh implementasi query param yang sesuai dengan requirement
- Dipakai: kode di line 177-180
- Diubah: tidak ada

## Implementasi adaBentrok()
- Prompt: Bagaimana cara mendeteksi bentrok yang sesuai dengan requirement?
- Dipakai: fungsi adaBentrok()
- Diubah: tidak ada

## Refaktor kode untuk menggunakan fungsi helper yang telah dibuat
- Prompt: Ubah kode yang belum memakai helper untuk memakai fungsi helper
- Dipakai: semua
- Diubah: tidak ada

## Mengecek kesesuaian dengan requirements
- Prompt: Cek apakah kode saya sudah sesuai dengan requirement
- Dipakai: semua
- Diubah: kode-kode yang perlu diperbaiki supaya sesuai
