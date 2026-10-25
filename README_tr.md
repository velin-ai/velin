# Velin

**Kendi AI Homelab'in.** Velin; kullanıcı hesapları, API anahtarları, sohbet geçmişi ve [Ollama](https://ollama.com) tabanlı model geçidi içeren, kendi sunucunuzda çalışan bir yapay zekâ sohbet backend'i ve web arayüzüdür.

[English](README.md) | Türkçe

## Özellikler

- Kısa ömürlü JWT erişim jetonları ve döndürülen refresh jetonlarıyla kullanıcı adı/parola kimlik doğrulaması
- API anahtarı yönetimi — anahtarlar yalnızca bir kez gösterilir ve SHA-256 özeti olarak saklanır
- PostgreSQL'de saklanan sohbetler ve mesajlar, kullanıcı bazlı sahiplik kontrolü
- Sohbetin son 20 mesajının bağlam olarak modele gönderilmesi
- Ollama üzerinden kendi sunucunuzda model çağrıları; harici yapay zekâ servisi gerekmez
- Go binary'si tarafından sunulan gömülü web arayüzü
- Uygulama içi dil seçiciyle İngilizce ve Türkçe arayüz

## Teknolojiler

| Katman | Teknoloji |
| --- | --- |
| Backend | Go 1.27, chi router |
| Veritabanı | PostgreSQL |
| Kimlik doğrulama | JWT (HS256), bcrypt, SHA-256 |
| Yapay zekâ | Ollama |
| Frontend | React 19, TypeScript, Vite |
| Yapılandırma | `go-keyring` ile işletim sistemi anahtar kasası |

## Gereksinimler

- Go 1.27 veya üzeri
- Node.js 20.19 veya üzeri ve npm
- PostgreSQL 14 veya üzeri
- Yerelde çalışan (veya yapılandırılan adresten erişilebilen) [Ollama](https://ollama.com)
- Çalışan bir işletim sistemi anahtar kasası: Linux'ta Secret Service, macOS'ta Keychain, Windows'ta Credential Manager

## Kurulum

### 1. Veritabanı

Boş bir PostgreSQL veritabanı oluşturun:

```sh
createdb velin
```

Tablolar uygulama başlarken otomatik oluşturulur.

### 2. Yapılandırma

Velin, gizli değerlerini işletim sistemi anahtar kasasından okur. Aşağıdaki komutlarla bir kez kaydedin. Her komut değeri yazıp çıkar.

```sh
go run ./cmd/velin -set "admin-parolaniz"
go run ./cmd/velin -db "postgres://kullanici:parola@localhost:5432/velin?sslmode=disable"
go run ./cmd/velin -jwt "uzun-rastgele-bir-sir"
```

Grafik arayüzü olmayan (headless) Linux sistemlerinde `gnome-keyring` gibi bir Secret Service sağlayıcısının çalıştığından ve kilidinin açık olduğundan emin olun.

### 3. Frontend'i derleyin

Go binary'si `ui/api/dist` klasörünü derleme zamanında gömer; bu yüzden önce frontend derlenmelidir:

```sh
cd ui/api
npm install
npm run build
cd ../..
```

`ui/api/dist` yoksa `go build` embed hatası verir.

### 4. Çalıştırın

```sh
go run ./cmd/velin
```

veya bir binary oluşturun:

```sh
go build -o velin ./cmd/velin
./velin
```

<http://localhost:8080> adresini açın.

### `go install` ile kurulum

`velin` binary'sini yayınlanan modülden doğrudan kurabilirsiniz:

```sh
go install github.com/velin-ai/velin/cmd/velin@latest
```

`go install` binary'yi `$(go env GOPATH)/bin` (veya `$GOBIN`) dizinine koyar. Bu dizin `PATH` içinde değilse ekleyin:

```sh
export PATH="$PATH:$(go env GOPATH)/bin"
```

Ardından Velin'i bir kez yapılandırıp diğer komutlar gibi çalıştırın:

```sh
velin -set "admin-parolaniz"
velin -db "postgres://kullanici:parola@localhost:5432/velin?sslmode=disable"
velin -jwt "uzun-rastgele-bir-sir"
velin
```

Web arayüzü derleme zamanında binary'ye gömülür. `go install ...@latest` yalnızca `ui/api/dist` yayınlanan modül sürümünde mevcutsa çalışır. Yerel bir kopyadan kurulum için önce frontend'i derleyip depodan kurun:

```sh
cd ui/api && npm run build && cd ../..
go install ./cmd/velin
```

### Geliştirme modu

API'yi ve Vite geliştirme sunucusunu yan yana çalıştırın:

```sh
# terminal 1
go run ./cmd/velin
```

```sh
# terminal 2
cd ui/api
npm install
npm run dev
```

Frontend <http://localhost:5173> üzerinde çalışır ve `/api` isteklerini `http://localhost:8080` adresine yönlendirir; böylece arayüz değişiklikleri için Go derlemesi gerekmez.

## Ollama

Sohbet istekleri `http://localhost:11434/api/chat` adresine gönderilir. Adres `internal/api/ai/ollama.go` içinde tanımlıdır; Ollama başka bir sunucuda çalışıyorsa burayı güncelleyin.

```sh
ollama serve
ollama pull llama3.1
```

Model adı her sohbet isteğiyle birlikte gönderilir; Ollama'da bulunan herhangi bir model kullanılabilir.

## Yapılandırma bayrakları

| Bayrak | Açıklama | Varsayılan |
| --- | --- | --- |
| `-set` | Admin parolasını anahtar kasasına kaydeder veya değiştirir | — |
| `-db` | PostgreSQL bağlantı adresini kaydeder veya değiştirir | — |
| `-jwt` | JWT imzalama sırrını kaydeder veya değiştirir | — |
| `-api-port` | API sunucusunun dinleyeceği adres | `:8080` |
| `-admin-port` | Admin servisi için ayrılmış | `:9090` |

Değerler anahtar kasasında `velin` servis adı altında saklanır.

## API

Tüm API rotaları `/api/v1` önekiyle başlar.

| Yöntem | Yol | Kimlik doğrulama | Açıklama |
| --- | --- | --- | --- |
| `POST` | `/auth/register` | Herkese açık | Hesap oluşturur |
| `POST` | `/auth/login` | Herkese açık | Giriş yapar; erişim ve refresh jetonları döner |
| `POST` | `/auth/refresh` | Herkese açık | Refresh jetonunu döndürür ve yeni erişim jetonu verir |
| `POST` | `/api-key` | JWT | API anahtarı oluşturur |
| `DELETE` | `/api-key/{id}` | JWT | API anahtarlarınızdan birini siler |
| `GET` | `/api-keys` | JWT | API anahtarlarınızı listeler |
| `POST` | `/chats` | JWT | Sohbet oluşturur |
| `GET` | `/chats` | JWT | Sohbetlerinizi listeler |
| `POST` | `/chats/{chat_id}/messages` | API anahtarı | Mesaj gönderir ve model yanıtını döner |
| `GET` | `/` | Herkese açık | Gömülü web arayüzünü sunar |

JWT korumalı rotalar şu başlığı bekler:

```http
Authorization: Bearer <erisim-jetonu>
```

Sohbet mesajı rotaları bir API anahtarı bekler:

```http
Authorization: Bearer vln_live_<anahtar>
```

### Örnek

```sh
curl -s http://localhost:8080/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"demo","password":"secret"}'
```

## Proje yapısı

```text
cmd/velin/                 Giriş noktası ve HTTP router
internal/admin/config/     İşletim sistemi anahtar kasası yapılandırması
internal/api/ai/           Ollama istemcisi
internal/api/handlers/     Auth, API anahtarı ve sohbet handler'ları
internal/api/middlewares/  JWT ve API anahtarı kimlik doğrulaması
internal/api/models/       İstek ve yanıt modelleri
internal/store/            PostgreSQL bağlantısı ve şema
ui/api/                    React + TypeScript frontend
ui/api/src/i18n/           İngilizce ve Türkçe çeviriler
assets/logo.svg            Marka logosu
```

## Güvenlik notları

- Parolalar bcrypt (cost 12) ile hash'lenir.
- Erişim jetonları HS256 ile imzalanır ve 30 dakika sonra sona erer.
- Refresh jetonları 30 gün geçerlidir, SHA-256 özeti olarak saklanır ve kullanıldığında döndürülür.
- API anahtarları yalnızca bir kez gösterilir ve SHA-256 özeti olarak saklanır.
- Sohbet ve mesaj erişimi kimliği doğrulanmış kullanıcıyla sınırlıdır.
- Frontend jetonları `localStorage` içinde tutar. Velin'i herkese açık bir ortamda yayınlayacaksanız HTTP-only cookie tabanlı bir oturum katmanı ve ek sıkılaştırmalar değerlendirin.

## Uluslararasılaştırma

Arayüz varsayılan olarak İngilizce açılır. Türkçe, **Dashboard → Settings → Language** bölümünden seçilebilir. Seçim `localStorage` içinde saklanır ve anında uygulanır.

Çeviri dosyaları:

- `ui/api/src/i18n/en.json`
- `ui/api/src/i18n/tr.json`

Yeni bir dil eklemek için bir JSON dosyası ekleyin, `ui/api/src/i18n/I18nProvider.tsx` ve `ui/api/src/i18n/context.ts` içinde kaydedin, mevcut sözlüklere dil etiketini ekleyin.

## Lisans

[MIT](LICENSE)
