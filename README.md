Tender Management Backend
Bu loyiha Tender Management Backend xizmatini amalga oshiradi. Loyihada backend dasturlash orqali tenderlarni boshqarish tizimini yaratish uchun mo'ljallangan imkoniyatlar mavjud.

📋 Xususiyatlar
Tenderlarni boshqarish: Tenderlarni yaratish, o‘zgartirish, o‘chirish va ko‘rish funksiyalari.
Foydalanuvchi boshqaruvi: Foydalanuvchi autentifikatsiyasi va ruxsatlar tizimi.
Ma'lumotlar bazasi: Ma'lumotlarni saqlash uchun optimallashtirilgan model va migratsiyalar.
RESTful API: Tender va foydalanuvchilar bilan ishlash uchun API xizmatlari.
JSON formatida javoblar: Xatolar va muvaffaqiyatli natijalar uchun moslashuvchan javoblar.
📁 Loyihaning Strukturasi
graphql
Копировать код
Tender_Management_Backend/
├── cmd/                  # Asosiy dastur ishga tushirish kodi
├── internal/             # Ichki logika (servislar, modellar, handlerlar)
│   ├── handlers/         # HTTP endpointlar uchun handlerlar
│   ├── models/           # Ma'lumotlar modeli va strukturalar
│   ├── services/         # Biznes logikasi uchun xizmatlar
├── configs/              # Konfiguratsiya fayllari
├── migrations/           # Ma'lumotlar bazasi migratsiyalari
├── docs/                 # API uchun hujjatlar
└── README.md             # Loyihaning asosiy hujjati
🔧 O'rnatish va Ishga Tushirish
Kod bazasini yuklab oling:

bash
Копировать код
git clone https://github.com/username/Tender_Management_Backend.git
cd Tender_Management_Backend
Kerakli kutubxonalarni o‘rnatish:

bash
Копировать код
go mod tidy
Ma'lumotlar bazasi sozlamalari: configs/ papkasidagi config.yaml faylini moslashtiring.

Migratsiyalarni ishga tushiring:

bash
Копировать код
go run cmd/migrate.go
Xizmatni ishga tushiring:

bash
Копировать код
go run cmd/main.go
📖 API Huqumatlari
API hujjatlari Swagger orqali taqdim etilgan.
Ishga tushirilgandan so‘ng quyidagi URL manzilga kiring:
bash
Копировать код
http://localhost:8080/swagger/index.html
⚙️ Texnologiyalar
Backend: Go (Golang)
Ma'lumotlar bazasi: PostgreSQL
API: RESTful
Migratsiyalar: SQL bo‘yicha avtomatik migratsiyalar
JSON: Moslashuvchan formatda ma'lumot almashinuvi
🛠 Hissa Qo'shish
Fork qiling va loyihani yuklab oling.
O‘zgarishlar kiriting va test qiling.
Pull Request yarating.
