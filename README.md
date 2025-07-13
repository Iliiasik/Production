# 🏭 Web application for sports equipment manufacturers

A specialized web-based management system for automating business processes in the production of sports equipment — including gym machines, balls, apparel, and more.

## 📌 Introduction

Modern information technologies are transforming traditional management methods, especially in the manufacturing industry. Sports equipment companies face numerous challenges — from organizing and controlling production to ensuring product quality and optimizing procurement and sales.

This project delivers a custom web application tailored for sports equipment manufacturers, aiming to automate core operations, accelerate decision-making, and boost overall business efficiency.

## ❗ Problem statement

Many manufacturers still rely on manual paperwork for managing operations, which leads to:

- ⏱️ High labor intensity and time consumption  
- 📉 Low reliability with risks of data loss  
- 🔍 Inconvenience in accessing and searching for information  

## 🎯 Project goals

The primary objectives of this system are:

- Optimize business processes and reduce operational costs  
- Minimize resource/material losses  
- Centralize data storage and streamline workflows  
- Improve internal analytics and decision-making

## 🚀 Key features

- 📦 Raw materials and inventory management (add, edit, write-off)  
- 👷 Employee database and payroll automation  
- 💰 Credit issuance and tracking  
- 🛠️ Production tracking and finished goods accounting  
- 📊 Sales and procurement control  
- 📈 Reports and analytics module

## 🖼️ Screenshots

> Interface language: **Russian**  
> Below are some screenshots of the application:

### 🔐 Authorization
<img width="2437" height="1263" alt="image" src="https://github.com/user-attachments/assets/dacad560-ea93-47ba-ba6c-a0134b568043" />

### 📊 Profile
<img width="974" height="502" alt="image" src="https://github.com/user-attachments/assets/75343a3e-6d59-477a-81ba-d4cbd6f7c76e" />

<img width="1686" height="1159" alt="image" src="https://github.com/user-attachments/assets/f4923c4a-333d-42de-9f16-55b82c530c8f" />

<img width="2504" height="1281" alt="image" src="https://github.com/user-attachments/assets/3184debc-0a89-45e6-ba03-0cbcb5e14576" />

### 👷 RBAC and ACL
<img width="974" height="490" alt="image" src="https://github.com/user-attachments/assets/360e1d60-34b1-4710-9e06-3f9711825826" />

<img width="974" height="509" alt="image" src="https://github.com/user-attachments/assets/d84eff6a-4c8c-4992-ae17-634176a02cdc" />

### 🛠️ Production module
<img width="974" height="708" alt="image" src="https://github.com/user-attachments/assets/4d76b0ab-e7f5-411c-9f76-3f819b4e0ec7" />

<img width="578" height="802" alt="image" src="https://github.com/user-attachments/assets/cac41bdd-b2b1-4a21-a219-f97a680f42f7" />

<img width="2354" height="1227" alt="image" src="https://github.com/user-attachments/assets/30e827c8-6960-4f87-bec7-e77d4e4ce10a" />

### 💳 Credits  
<img width="2538" height="1111" alt="image" src="https://github.com/user-attachments/assets/d65958ea-bb9a-400d-a3c6-c4f9e90d2d9d" />

### 💸 Credit payments
<img width="2513" height="1241" alt="image" src="https://github.com/user-attachments/assets/50e2359c-7cad-4ce7-a376-650a8690fdc2" />

### 💼 Salary calculation interface
<img width="2485" height="1191" alt="image" src="https://github.com/user-attachments/assets/e363b498-4354-4fee-9093-04214796cc02" />

<img width="2469" height="865" alt="image" src="https://github.com/user-attachments/assets/7d101118-86c0-407a-8a8a-7d7cfa844964" />

### 🏦 Budget overview
<img width="1640" height="858" alt="image" src="https://github.com/user-attachments/assets/668ae788-f5b7-4a9b-b370-00715ca95bbe" />

### 🧾 Reports module
<img width="974" height="450" alt="image" src="https://github.com/user-attachments/assets/25701412-56a0-459c-b17c-e90f4e8f50a3" />

<img width="974" height="487" alt="image" src="https://github.com/user-attachments/assets/507ad009-e518-466d-9e25-904d772e3d2e" />

### 🧾 Export
<img width="874" height="562" alt="image" src="https://github.com/user-attachments/assets/3723a5b3-1660-4ada-ae46-d867d4323c9e" />

<img width="974" height="454" alt="image" src="https://github.com/user-attachments/assets/de72d83c-802b-4ee3-8b8f-de226142bc13" />

<img width="463" height="659" alt="image" src="https://github.com/user-attachments/assets/7063f541-b0f4-4aa1-93d4-1bfe06eb7751" />

<img width="533" height="647" alt="image" src="https://github.com/user-attachments/assets/627495c5-051b-4668-847d-7a99f347d45d" />

## ⚙️ Tech stack

- **Backend**: Go (Gin), PostgreSQL,   
- **Frontend**: HTML, SCSS, Vanilla JavaScript  
- **Other**: JWT Auth, SweetAlert2, Role-Based Access (RBAC)

## 🧠 Business Logic & Data Layer

The core business processes are implemented using PostgreSQL functions, stored procedures, and triggers, ensuring high performance and data integrity at the database level.

### ⚙️ Key database mechanisms:
- **Stored procedures** for:
  - Automated salary calculation  
  - Credit payment allocation  
  - Report generation (e.g., sales, payroll, production)  
- **Functions** for on-demand calculations and summaries  
- **Triggers** for enforcing data integrity and automating workflow transitions. For example, a trigger is used to check budget availability before allowing certain operations, ensuring financial constraints are respected.

This architecture guarantees transactional accuracy, maintainability, and modularity by offloading critical business logic to the database layer.
