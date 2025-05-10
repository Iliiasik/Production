package database

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"production/models"
)

var DB *gorm.DB

func InitDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")

	if user == "" || password == "" || dbName == "" || host == "" || port == "" {
		log.Fatal("One or more environment variables are missing")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbName, port,
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	log.Println("Connection success")

	if DB == nil {
		log.Fatal("Database connection is nil")
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Error getting DB instance: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}
	log.Println("Database ping success")

	err = DB.AutoMigrate(
		&models.Unit{},
		&models.RawMaterial{},
		&models.FinishedGood{},
		&models.Position{},
		&models.Employee{},
		&models.SalaryRecord{},
		&models.Ingredient{},
		&models.Budget{},
		&models.RawMaterialPurchase{},
		&models.ProductSale{},
		&models.ProductProduction{},
		&models.Credit{},
		&models.CreditPayment{},
		&models.Permission{},
		&models.UserPermission{},
		&models.PositionPermission{},
	)
	if err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}

	log.Println("Migrations completed successfully")

	seedPermissions()

	// Привязка всех разрешений к роли "Админ"
	seedAdminPositionPermissions()
}
func seedPermissions() {
	permissions := []models.Permission{
		// Profile
		{Name: "/home", Description: "Личный кабинет пользователя", Category: "Профиль", VisibleToUser: true},

		// Inventory
		{Name: "/units", Description: "Просмотр единиц измерения", Category: "Инвентарь", VisibleToUser: true},
		{Name: "/units/delete/:id", Description: "Удаление единицы измерения", Category: "Инвентарь", VisibleToUser: true},
		{Name: "/units/add", Description: "Добавление единицы измерения", Category: "Инвентарь", VisibleToUser: true},
		{Name: "/units/edit/:id", Description: "Редактирование единицы измерения", Category: "Инвентарь", VisibleToUser: true},
		{Name: "/units/get/:id", Description: "Получение конкретной единицы измерения", Category: "Инвентарь", VisibleToUser: false},
		{Name: "/units/list", Description: "Список единиц измерения", Category: "Инвентарь", VisibleToUser: false},

		{Name: "/raw-materials", Description: "Просмотр сырья", Category: "Инвентарь", VisibleToUser: true},
		{Name: "/raw-materials/delete/:id", Description: "Удаление сырья", Category: "Инвентарь", VisibleToUser: true},
		{Name: "/raw-materials/add", Description: "Добавление сырья", Category: "Инвентарь", VisibleToUser: true},
		{Name: "/raw-materials/edit/:id", Description: "Редактирование сырья", Category: "Инвентарь", VisibleToUser: true},
		{Name: "/raw-materials/get/:id", Description: "Получение конкретного сырья", Category: "Инвентарь", VisibleToUser: false},
		{Name: "/raw-materials/list", Description: "Список сырья", Category: "Инвентарь", VisibleToUser: false},

		{Name: "/finished-goods", Description: "Просмотр готовой продукции", Category: "Инвентарь", VisibleToUser: true},
		{Name: "/finished-goods/delete/:id", Description: "Удаление готовой продукции", Category: "Инвентарь", VisibleToUser: true},
		{Name: "/finished-goods/add", Description: "Добавление готовой продукции", Category: "Инвентарь", VisibleToUser: true},
		{Name: "/finished-goods/edit/:id", Description: "Редактирование готовой продукции", Category: "Инвентарь", VisibleToUser: true},
		{Name: "/finished-goods/get/:id", Description: "Получение конкретной готовой продукции", Category: "Инвентарь", VisibleToUser: false},
		{Name: "/finished-goods/list", Description: "Список готовой продукции", Category: "Инвентарь", VisibleToUser: false},

		{Name: "/ingredients", Description: "Просмотр ингредиентов", Category: "Инвентарь", VisibleToUser: true},
		{Name: "/ingredients/delete/:id", Description: "Удаление ингредиента", Category: "Инвентарь", VisibleToUser: true},
		{Name: "/ingredients/add", Description: "Добавление ингредиента", Category: "Инвентарь", VisibleToUser: true},
		{Name: "/ingredients/edit/:id", Description: "Редактирование ингредиента", Category: "Инвентарь", VisibleToUser: true},
		{Name: "/ingredients/get/:id", Description: "Получение конкретного ингредиента", Category: "Инвентарь", VisibleToUser: false},
		{Name: "/ingredients/:product_id", Description: "Просмотр ингредиентов по продукту", Category: "Инвентарь", VisibleToUser: false},
		{Name: "/ingredients/used-raw-materials/:product_id", Description: "Просмотр используемого сырья по продукту", Category: "Инвентарь", VisibleToUser: false},
		{Name: "/ingredients/list", Description: "Список ингредиентов", Category: "Инвентарь", VisibleToUser: false},

		// Transactions
		{Name: "/raw-material-purchases", Description: "Просмотр закупок сырья", Category: "Транзакции", VisibleToUser: true},
		{Name: "/purchases/delete/:id", Description: "Удаление закупки", Category: "Транзакции", VisibleToUser: true},
		{Name: "/purchases/add", Description: "Добавление закупки", Category: "Транзакции", VisibleToUser: true},

		{Name: "/production", Description: "Просмотр производства", Category: "Транзакции", VisibleToUser: true},
		{Name: "/production/produce/:product_id", Description: "Производство продукции", Category: "Транзакции", VisibleToUser: true},

		{Name: "/sales", Description: "Просмотр продаж", Category: "Транзакции", VisibleToUser: true},
		{Name: "/sales/add", Description: "Добавление продажи", Category: "Транзакции", VisibleToUser: true},

		// Finance
		{Name: "/budget", Description: "Просмотр бюджета", Category: "Финансы", VisibleToUser: true},
		{Name: "/budget/get-row/:id", Description: "Получение строки бюджета", Category: "Финансы", VisibleToUser: false},
		{Name: "/budget/get", Description: "Получение данных бюджета", Category: "Финансы", VisibleToUser: false},
		{Name: "/markup/get", Description: "Получение наценки", Category: "Финансы", VisibleToUser: false},
		{Name: "/budget/update/:id", Description: "Редактирование бюджета", Category: "Финансы", VisibleToUser: true},

		{Name: "/credits", Description: "Просмотр кредитов", Category: "Финансы", VisibleToUser: true},
		{Name: "/credits/add", Description: "Добавление кредита", Category: "Финансы", VisibleToUser: true},
		{Name: "/credits/:id/payments", Description: "Просмотр выплат по кредиту", Category: "Финансы", VisibleToUser: true},
		{Name: "/credits/pay/:id", Description: "Погашение кредита", Category: "Финансы", VisibleToUser: true},

		// Employees
		{Name: "/employees", Description: "Просмотр сотрудников", Category: "Сотрудники", VisibleToUser: true},
		{Name: "/employees/list", Description: "Список сотрудников", Category: "Сотрудники", VisibleToUser: false},
		{Name: "/positions/list", Description: "Список должностей", Category: "Сотрудники", VisibleToUser: false},
		{Name: "/employees/get/:id", Description: "Получение данных сотрудника", Category: "Сотрудники", VisibleToUser: false},
		{Name: "/employees/add", Description: "Добавление сотрудника", Category: "Сотрудники", VisibleToUser: true},
		{Name: "/employees/next-username", Description: "Получение номера логина", Category: "Сотрудники", VisibleToUser: false},
		{Name: "/employees/edit/:id", Description: "Редактирование сотрудника", Category: "Сотрудники", VisibleToUser: true},
		{Name: "/employees/delete/:id", Description: "Удаление сотрудника", Category: "Сотрудники", VisibleToUser: true},

		{Name: "/positions/add", Description: "Добавление должности", Category: "Сотрудники", VisibleToUser: true},
		{Name: "/positions/edit/:id", Description: "Редактирование должности", Category: "Сотрудники", VisibleToUser: true},
		{Name: "/positions/delete/:id", Description: "Удаление должности", Category: "Сотрудники", VisibleToUser: true},
		{Name: "/positions/get/:id", Description: "Получение данных должности", Category: "Сотрудники", VisibleToUser: false},
		{Name: "/positions", Description: "Просмотр списка должностей", Category: "Сотрудники", VisibleToUser: true},

		{Name: "/salaries", Description: "Просмотр зарплат", Category: "Сотрудники", VisibleToUser: true},
		{Name: "/salaries/:year/:month", Description: "Получение зарплат за дату", Category: "Сотрудники", VisibleToUser: false},
		{Name: "/salaries/calculate/:year/:month", Description: "Расчёт зарплаты", Category: "Сотрудники", VisibleToUser: true},
		{Name: "/salaries/edit/:id", Description: "Редактирование записи зарплаты", Category: "Сотрудники", VisibleToUser: true},
		{Name: "/salaries/pay/:year/:month", Description: "Выдача зарплаты", Category: "Сотрудники", VisibleToUser: true},
		{Name: "/salaries/total-unpaid/:year/:month", Description: "Получение суммы невыплаченных зарплат", Category: "Сотрудники", VisibleToUser: false},

		{Name: "/admin/roles", Description: "Управление ролями", Category: "Админ", VisibleToUser: false},
		{Name: "/admin/roles/:id/permissions", Description: "Просмотр и управление разрешениями роли", Category: "Админ", VisibleToUser: false},
		{Name: "/admin/roles/:id/permissions/update", Description: "Обновление разрешений роли", Category: "Админ", VisibleToUser: false},

		{Name: "/admin/users", Description: "Управление пользователями", Category: "Админ", VisibleToUser: false},
		{Name: "/admin/users/:id/permissions", Description: "Просмотр и управление разрешениями пользователя", Category: "Админ", VisibleToUser: false},
		{Name: "/admin/users/:id/permissions/update", Description: "Обновление разрешений пользователя", Category: "Админ", VisibleToUser: false},
		{Name: "/admin/users/:id/role", Description: "Просмотр роли пользователя", Category: "Админ", VisibleToUser: false},

		{Name: "/admin/permissions", Description: "Управление разрешениями", Category: "Админ", VisibleToUser: false},

		{Name: "/reports", Description: "Страница отчетов", Category: "Отчеты", VisibleToUser: false},
		{Name: "/reports/sales", Description: "Отчеты о продажах", Category: "Отчеты", VisibleToUser: true},
		{Name: "/reports/productions", Description: "Отчеты по производству", Category: "Отчеты", VisibleToUser: true},
		{Name: "/reports/purchases", Description: "Отчеты о закупках", Category: "Отчеты", VisibleToUser: true},
		{Name: "/reports/salaries", Description: "Отчеты о зарплатах", Category: "Отчеты", VisibleToUser: true},
		{Name: "/reports/payments", Description: "Отчеты по кредитным выплатам", Category: "Отчеты", VisibleToUser: true},
		{Name: "/reports/export", Description: "Экспорт отчетов", Category: "Отчеты", VisibleToUser: true},
	}

	for _, p := range permissions {
		var existing models.Permission
		err := DB.Where("name = ?", p.Name).First(&existing).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := DB.Create(&p).Error; err != nil {
					log.Fatalf("Не удалось вставить permission %q: %v", p.Name, err)
				}
			} else {
				log.Fatalf("Ошибка проверки permission %q: %v", p.Name, err)
			}
		}
	}

	log.Println("seedPermissions: permissions table is up to date")
}

func seedAdminPositionPermissions() {
	// ищем позицию
	var adminPos models.Position
	if err := DB.Where("name = ?", "Админ").First(&adminPos).Error; err != nil {
		log.Printf("Position 'Админ' not found, skip seeding admin permissions: %v", err)
		return
	}

	// получаем все permissions
	var allPerms []models.Permission
	if err := DB.Find(&allPerms).Error; err != nil {
		log.Fatalf("failed to load permissions: %v", err)
	}

	for _, perm := range allPerms {
		var pp models.PositionPermission
		// проверяем, не создана ли уже связка
		err := DB.
			Where("position_id = ? AND permission_id = ?", adminPos.ID, perm.ID).
			First(&pp).Error

		if err == gorm.ErrRecordNotFound {
			pp = models.PositionPermission{
				PositionID:   adminPos.ID,
				PermissionID: perm.ID,
			}
			if err := DB.Create(&pp).Error; err != nil {
				log.Fatalf("failed to seed position_permission for admin: %v", err)
			}
		} else if err != nil {
			log.Fatalf("error checking position_permission: %v", err)
		}
	}
	log.Println("Admin position_permissions seeding done")
}
