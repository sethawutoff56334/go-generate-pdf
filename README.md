# Generate PDF Service

REST API สำหรับสร้างใบเสร็จรับเงิน ในรูปแบบ PDF จากข้อมูล JSON

## The Choice:

Go module: chromedp

## The Reason: chromedp

1. **Render ได้เหมือนเบราว์เซอร์จริง** — รองรับ CSS ครบถ้วน (Flexbox, Grid, Web Fonts) ไม่ต้องกังวลว่า layout จะเพี้ยน
2. **รองรับภาษาไทยได้ดี** — ใช้ font rendering ของ Chrome โดยตรง ไม่ต้องหา font เอง
3. **ใช้ HTML/CSS เป็น template** — ออกแบบหน้าตาเอกสารได้ง่าย แก้ไข style ได้สะดวก ไม่ต้องวาด PDF ทีละ element

## Prerequisites

- **Google Chrome** หรือ **Chromium** ติดตั้งอยู่ในเครื่อง (chromedp จะเรียกใช้ headless mode)

## How to run

### 1. Clone & Install dependencies

```bash
git clone https://github.com/sethawutoff56334/go-generate-pdf.git
go mod download
```

### 2. รันเซิร์ฟเวอร์

```bash
go run main.go
```

เซิร์ฟเวอร์จะเริ่มทำงานที่ port `8080`

### 3. เรียกใช้ API แล้วเปิดดู PDF

**ใช้ไฟล์ example.json:**

```bash
curl -X POST http://localhost:8080/generate-pdf \
  -H "Content-Type: application/json" \
  -d @example.json
```

**ใส่ JSON ตรงๆ ในคำสั่ง:**

```bash
curl -X POST http://localhost:8080/generate-pdf \
  -H "Content-Type: application/json" \
  -d '{
    "receipt_no": "R1234",
    "customer": {
      "first_name": "sethawut",
      "last_name": "pornsiripiyakul",
      "email": "sethawutoff@gmail.com",
      "phone_no": "0851138723"
    },
    "pricing_summary": {
      "sub_total": 1429,
      "discount": 150,
      "total": 1279
    },
    "product_list": [
      {"product_name": "มาม่ากระป๋อง", "qty": 3, "price_per_unit": 30},
      {"product_name": "น้ำดื่มสิงห์", "qty": 12, "price_per_unit": 7},
      {"product_name": "ข้าวสารหอมมะลิ 5 กก.", "qty": 1, "price_per_unit": 250},
      {"product_name": "นมสดโฟร์โมสต์", "qty": 6, "price_per_unit": 18},
      {"product_name": "ไข่ไก่แพ็ค 10 ฟอง", "qty": 2, "price_per_unit": 55},
      {"product_name": "น้ำมันพืชองุ่น 1 ลิตร", "qty": 1, "price_per_unit": 75},
      {"product_name": "ผงซักฟอกแอทแทค", "qty": 2, "price_per_unit": 120},
      {"product_name": "สบู่ลักส์", "qty": 4, "price_per_unit": 35},
      {"product_name": "ยาสีฟันคอลเกต", "qty": 1, "price_per_unit": 65},
      {"product_name": "กระดาษทิชชู่ 6 ม้วน", "qty": 3, "price_per_unit": 89}
    ]
  }' 
```