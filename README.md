# CS367 Lab & Homework - Student API (Individual Assignment)

**โปรเจกต์นี้เป็นผลงานเดี่ยว (Individual Assignment) ในรายวิชา CS367** เอกสารนี้รวบรวมขั้นตอนการทำงานจริง ตั้งแต่การวางโครงสร้าง ไปจนถึงการแก้ปัญหาระดับ System Environment และ Command Line ที่พบระหว่างการพัฒนา RESTful API ด้วยภาษา Go

## 1. การวางโครงสร้างระบบ (Layered Architecture)
* **สิ่งที่ทำ:** ออกแบบโค้ดโดยแยกสัดส่วนความรับผิดชอบชัดเจน ได้แก่ `models` (โครงสร้างข้อมูล), `config` (เชื่อมต่อฐานข้อมูล), `repositories` (จัดการ DB), `services` (จัดการ Business Logic), และ `handlers` (รับส่ง HTTP Request/Response)
* **ผลลัพธ์:** ทำให้โค้ดเป็นระเบียบ แต่เมื่อเริ่มรันเพื่อเชื่อมต่อฐานข้อมูลจริง กลับพบปัญหาในขั้นตอนถัดไป

## 2. ปัญหาโลกแตก `CGO_ENABLED` บน Windows
* **ปัญหาที่พบ:** เมื่อสั่ง `go run main.go` ระบบฟ้อง Error เกี่ยวกับ CGO สาเหตุเกิดจากไลบรารี SQLite มาตรฐาน (`go-sqlite3`) จำเป็นต้องใช้ C Compiler (เช่น GCC/MinGW) ในการทำงาน ซึ่งสภาพแวดล้อมบน Windows ไม่ได้ติดตั้งไว้
* **วิธีการแก้ไข:** ทำการเปลี่ยน Driver ของ Database ใหม่ โดยหันมาใช้ไลบรารีที่เป็น Pure Go แทน เพื่อให้ทำงานแบบ Cross-platform ได้ทันทีโดยไม่ง้อภาษา C
* **คำสั่งที่ใช้แก้ไข:**
  ```go
  // ในไฟล์ database.go เปลี่ยน Import เป็น:
  import "[github.com/glebarez/go-sqlite](https://github.com/glebarez/go-sqlite)"
จากนั้นรัน go mod tidy เพื่อโหลดไลบรารีใหม่ แล้วสั่ง go run main.go อีกครั้ง คราวนี้เซิร์ฟเวอร์รันผ่านและสร้างไฟล์ students.db ได้สำเร็จ (หน้าเว็บแสดงค่า null)

3. ศึกปะทะ Windows PowerShell และคำสั่ง curl
ปัญหาที่พบ: เมื่อต้องการทดสอบยิง API (POST Request) เพื่อเพิ่มข้อมูลนักเรียน ได้พยายามใช้คำสั่ง curl.exe ผ่าน PowerShell แต่พบ Error invalid character '\''' looking for beginning of object key string

สาเหตุ: PowerShell บน Windows มีปัญหาจุกจิกเรื่องการอ่านเครื่องหมายคำพูด (Single/Double Quotes) และการ Escape String ในก้อน JSON ทำให้เซิร์ฟเวอร์อ่านข้อมูลที่ส่งไปไม่ออก

4. ทางออกด้วย VS Code Extension และ SQL Query
วิธีการแก้ไข: เพื่อข้ามข้อจำกัดของ PowerShell จึงเปลี่ยนวิธีเพิ่มข้อมูลมาเป็นการจัดการ Database โดยตรงผ่านเครื่องมือของ VS Code

สิ่งที่ทำ: 1. ใช้ Extension SQLite เพื่อเปิดไฟล์ students.db
2. เปิดหน้าต่าง SQL Query และเขียนคำสั่งเพิ่มข้อมูลด้วยตนเอง:

SQL
INSERT INTO students (id, name, major, gpa) 
VALUES ('6609xxxx', 'Nong Puti', 'CS', 3.8);
สั่งรัน "Run Selected Query" เพื่ออัปเดตข้อมูลลงตาราง

5. การตรวจสอบผลลัพธ์สุดท้าย (The Moment of Truth)
สิ่งที่ทำ: รีเฟรชหน้าเว็บ http://localhost:8080/students อีกครั้งในขณะที่ Gin Server ยังทำงานอยู่

ผลลัพธ์ที่ได้: ข้อมูล JSON ของนักเรียนที่เพิ่ง Insert ลงไป ปรากฏขึ้นบนเบราว์เซอร์อย่างถูกต้องแทนที่คำว่า null

ข้อสรุป: เป็นการยืนยันว่า Layer ทั้งหมดของแอปพลิเคชัน (Handler -> Service -> Repository -> Database) เชื่อมต่อและทำงานร่วมกันได้อย่างสมบูรณ์แบบ 100%
