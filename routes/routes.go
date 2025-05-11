package routes

import (
	"github.com/gin-gonic/gin"
	"production/controllers"
	"production/middleware"
)

func RegisterRoutes(r *gin.Engine) {

	r.GET("/", controllers.ShowLoginPage)

	// Units
	r.GET("/units", middleware.Authorize("/units"), controllers.ListUnits)
	r.DELETE("/units/delete/:id", middleware.Authorize("/units/delete/:id"), controllers.DeleteUnit)
	r.POST("/units/add", middleware.Authorize("/units/add"), controllers.AddUnit)
	r.POST("/units/edit/:id", middleware.Authorize("/units/edit/:id"), controllers.UpdateUnit)
	r.GET("/units/get/:id", middleware.Authorize("/units/get/:id"), controllers.GetUnit)
	r.GET("/units/list", middleware.Authorize("/units/list"), controllers.GetUnitsList)

	// Raw materials
	r.GET("/raw-materials", middleware.Authorize("/raw-materials"), controllers.ListRawMaterials)
	r.DELETE("/raw-materials/delete/:id", middleware.Authorize("/raw-materials/delete/:id"), controllers.DeleteRawMaterial)
	r.POST("/raw-materials/add", middleware.Authorize("/raw-materials/add"), controllers.AddRawMaterial)
	r.POST("/raw-materials/edit/:id", middleware.Authorize("/raw-materials/edit/:id"), controllers.UpdateRawMaterial)
	r.GET("/raw-materials/get/:id", middleware.Authorize("/raw-materials/get/:id"), controllers.GetRawMaterial)
	r.GET("/raw-materials/list", middleware.Authorize("/raw-materials/list"), controllers.GetRawMaterialsList)

	// Finished goods
	r.GET("/finished-goods", middleware.Authorize("/finished-goods"), controllers.ListFinishedGoods)
	r.DELETE("/finished-goods/delete/:id", middleware.Authorize("/finished-goods/delete/:id"), controllers.DeleteFinishedGood)
	r.POST("/finished-goods/add", middleware.Authorize("/finished-goods/add"), controllers.AddFinishedGood)
	r.POST("/finished-goods/edit/:id", middleware.Authorize("/finished-goods/edit/:id"), controllers.UpdateFinishedGood)
	r.GET("/finished-goods/get/:id", middleware.Authorize("/finished-goods/get/:id"), controllers.GetFinishedGood)
	r.GET("/finished-goods/list", middleware.Authorize("/finished-goods/list"), controllers.GetFinishedGoodsList)

	// Ingredients
	r.GET("/ingredients", middleware.Authorize("/ingredients"), controllers.ListIngredients)
	r.DELETE("/ingredients/delete/:id", middleware.Authorize("/ingredients/delete/:id"), controllers.DeleteIngredient)
	r.POST("/ingredients/add", middleware.Authorize("/ingredients/add"), controllers.AddIngredient)
	r.POST("/ingredients/edit/:id", middleware.Authorize("/ingredients/edit/:id"), controllers.UpdateIngredient)
	r.GET("/ingredients/get/:id", middleware.Authorize("/ingredients/get/:id"), controllers.GetIngredient)
	r.GET("/ingredients/:product_id", middleware.Authorize("/ingredients/:product_id"), controllers.GetIngredientsByProduct)
	r.GET("/ingredients/used-raw-materials/:product_id", middleware.Authorize("/ingredients/used-raw-materials/:product_id"), controllers.GetUsedRawMaterialsByProduct)
	r.GET("/ingredients/list", middleware.Authorize("/ingredients/list"), controllers.GetIngredientsList)

	// Purchases
	r.GET("/raw-material-purchases", middleware.Authorize("/raw-material-purchases"), controllers.ListRawMaterialPurchases)
	r.DELETE("/purchases/delete/:id", middleware.Authorize("/purchases/delete/:id"), controllers.DeletePurchase)
	r.POST("/purchases/add", middleware.Authorize("/purchases/add"), controllers.AddPurchase)

	// Employees
	r.GET("/employees", middleware.Authorize("/employees"), controllers.ListEmployees)
	r.GET("/employees/list", middleware.Authorize("/employees/list"), controllers.GetEmployeesList)
	r.GET("/positions/list", middleware.Authorize("/positions/list"), controllers.GetAllPositions)
	r.GET("/employees/get/:id", middleware.Authorize("/employees/get/:id"), controllers.GetEmployee)
	r.POST("/employees/add", middleware.Authorize("/employees/add"), controllers.AddEmployee)
	r.GET("/employees/next-username", middleware.Authorize("/employees/next-username"), controllers.GetNextUsername)

	r.POST("/employees/edit/:id", middleware.Authorize("/employees/edit/:id"), controllers.UpdateEmployee)
	r.DELETE("/employees/delete/:id", middleware.Authorize("/employees/delete/:id"), controllers.DeleteEmployee)

	// Positions
	r.POST("/positions/add", middleware.Authorize("/positions/add"), controllers.AddPosition)
	r.POST("/positions/edit/:id", middleware.Authorize("/positions/edit/:id"), controllers.EditPosition)
	r.DELETE("/positions/delete/:id", middleware.Authorize("/positions/delete/:id"), controllers.DeletePosition)
	r.GET("/positions/get/:id", middleware.Authorize("/positions/get/:id"), controllers.GetPosition)
	r.GET("/positions", middleware.Authorize("/positions"), controllers.ListPositions)

	// Budget
	r.GET("/budget", middleware.Authorize("/budget"), controllers.BudgetList)
	r.GET("/budget/get-row/:id", middleware.Authorize("/budget/get-row/:id"), controllers.GetBudgetRow)
	r.GET("/budget/get", middleware.Authorize("/budget/get"), controllers.GetBudget)
	r.GET("/markup/get", middleware.Authorize("/markup/get"), controllers.GetMarkup)
	r.PUT("/budget/update/:id", middleware.Authorize("/budget/update/:id"), controllers.UpdateBudget)

	// Production
	r.GET("/production", middleware.Authorize("/production"), controllers.ListProductProduction)
	r.POST("/production/produce/:product_id", middleware.Authorize("/production/produce/:product_id"), controllers.ProduceProduct)

	// Sales
	r.GET("/sales", middleware.Authorize("/sales"), controllers.ListSales)
	r.POST("/sales/add", middleware.Authorize("/sales/add"), controllers.MakeSale)

	// Salaries
	r.GET("/salaries", middleware.Authorize("/salaries"), controllers.ShowSalariesPage)
	r.GET("/salaries/:year/:month", middleware.Authorize("/salaries/:year/:month"), controllers.GetSalaryByDate)
	r.POST("/salaries/calculate/:year/:month", middleware.Authorize("/salaries/calculate/:year/:month"), controllers.CalculateSalary)
	r.PUT("/salaries/edit/:id", middleware.Authorize("/salaries/edit/:id"), controllers.EditSalary)
	r.POST("/salaries/pay/:year/:month", middleware.Authorize("/salaries/pay/:year/:month"), controllers.PaySalaries)
	r.GET("/salaries/total-unpaid/:year/:month", middleware.Authorize("/salaries/total-unpaid/:year/:month"), controllers.GetUnpaidSalariesTotal)

	// Credits
	r.GET("/credits", middleware.Authorize("/credits"), controllers.ListCredits)
	r.POST("/credits/add", middleware.Authorize("/credits/add"), controllers.CreateCredit)
	r.GET("/credits/:id/payments", middleware.Authorize("/credits/:id/payments"), controllers.ShowPaymentsPage)
	r.POST("/credits/pay/:id", middleware.Authorize("/credits/pay/:id"), controllers.PayCredit)

	// Profile
	r.POST("/login", controllers.Login)
	r.POST("/change-password", controllers.ChangePassword)
	r.GET("/logout", controllers.Logout)
	r.GET("/home", middleware.Authorize("/home"), controllers.HomePage)
	r.GET("/user/permissions", controllers.GetUserPermissions)

	adminGroup := r.Group("/admin")
	{
		// Управление ролями
		adminGroup.GET("/roles", middleware.Authorize("/admin/roles"), controllers.AllPositions)
		adminGroup.GET("/roles/:id/permissions", middleware.Authorize("/admin/roles/:id/permissions"), controllers.GetPositionPermissions)
		adminGroup.PUT("/roles/:id/permissions/update", middleware.Authorize("/admin/roles/:id/permissions/update"), controllers.UpdatePositionPermissions) // Обновление разрешений роли

		// Управление пользователями
		adminGroup.GET("/users", middleware.Authorize("/admin/users"), controllers.AllUsers)
		adminGroup.GET("/users/:id/permissions", middleware.Authorize("/admin/users/:id/permissions"), controllers.GetEmployeePermissionsByID)
		adminGroup.PUT("/users/:id/permissions/update", middleware.Authorize("/admin/users/:id/permissions/update"), controllers.UpdateUserPermissions) // Обновление разрешений пользователя
		adminGroup.GET("/users/:id/role", middleware.Authorize("/admin/users/:id/role"), controllers.GetRoleByUserID)

		adminGroup.GET("/permissions", middleware.Authorize("/admin/permissions"), controllers.AllPermissions)
	}

	reportsGroup := r.Group("/reports")
	{
		reportsGroup.GET("", middleware.Authorize("/reports"), controllers.ShowReportsPage)
		reportsGroup.POST("/sales", middleware.Authorize("/reports/sales"), controllers.SalesReportHandler)
		reportsGroup.POST("/productions", middleware.Authorize("/reports/productions"), controllers.ProductionReportHandler)
		reportsGroup.POST("/purchases", middleware.Authorize("/reports/purchases"), controllers.PurchaseReportHandler)
		reportsGroup.POST("/salaries", middleware.Authorize("/reports/salaries"), controllers.SalaryReportHandler)
		reportsGroup.POST("/payments", middleware.Authorize("/reports/payments"), controllers.CreditReportHandler)
		reportsGroup.POST("/export", middleware.Authorize("/reports/export"), controllers.ExportReport)
	}

	tasksGroup := r.Group("/tasks")
	{
		tasksGroup.GET("/my", middleware.Authorize("/tasks/my"), controllers.GetMyTasks)
		tasksGroup.POST("", middleware.Authorize("/tasks"), controllers.AssignTask)
		tasksGroup.GET("/list", middleware.Authorize("/tasks"), controllers.GetTasks)
	}
}
