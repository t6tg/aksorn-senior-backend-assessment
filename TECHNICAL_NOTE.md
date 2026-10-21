# Technical Note

## 1. Design decision

Describe one decision in your implementation, why you chose it, and its tradeoff.

คำนวณ progress ตอนอ่านจาก attempt แทนที่จะเก็บตารางสรุป best score ทำให้ PostgreSQL เป็น source of truth และ คำนวนข้อมูลจาก Fixed course ถ้าบน production จริงการคำนวนโดยการไม่ เก็บ data แบบนี้จะทำให้การคำนวนนั้นทำซ้ำทุกครั้งที่มีการเรียก ควรมีตารางเก็บ progress และการผ่านของข้อมูล

## 2. Design questions

Brief answers are sufficient. No additional implementation is required.

### Retried requests

How would you enforce globally unique attempt IDs, return the original result for an identical replay, and return 409 for different content under the same ID? Consider simultaneous requests with the same ID.

ใส่ unique บน attempt_id ใน DB ถ้าซ้ำและเนื้อหาเหมือนเดิมให้คืนผลเดิม ถ้าต่างให้ตอบ 409 ส่วนกรณีมาพร้อมกัน unique index จะให้ผ่านแค่ตัวเดียว

### Concurrent requests

Two attempts with distinct IDs and scores of 90 and 70 arrive simultaneously for the same learner and lesson.

Can your implementation lose the higher score? What could a simultaneous progress read observe? What should a new read return after both writes have committed? Is any protection needed?

คะแนนสูงกว่าไม่หาย เพราะเก็บทุก attempt แล้วค่อยหาค่าสูงสุดตอนอ่าน เลยไม่ต้องป้องกันอะไรเพิ่ม

### Redis

Suppose attempt history grows and progress reads become slow.

Would you use Redis? What would you store, how would it stay consistent with PostgreSQL, and what would happen if Redis were unavailable or its data were lost? An alternative approach is welcome.

ลอง optimize query SQL ก่อน ถ้ายังช้าค่อย cache ข้อมูลใน Redis โดยให้ Postgres เป็นข้อมูลหลัก ถ้า Redis พังก็อ่านจาก Postgres แทน

## 3. Limitations and production readiness

List any limitations or unfinished work, and what you would verify before production use.

ตอนนี้ยังไม่กัน attempt_id ซ้ำกัน และอ่านข้อมูลจาก progress โดยดึง attempt ทั้งหมดมาคิด ถ้าข้อมูลใหญ่ขึ้นจะช้า ควรเพิ่ม unique constraint, และทำ benchmark

## 4. Time and AI tools

Time spent: 1 hr.

AI tools used: Github Copilot
