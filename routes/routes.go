package routes

import (
	"github.com/gin-gonic/gin"
	"production/controllers"
	"production/ws"
)

func RegisterRoutes(r *gin.Engine) {

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "home.html", nil)
	})

	// Units
	r.GET("/units", controllers.ListUnits)
	r.DELETE("/units/delete/:id", controllers.DeleteUnit)
	r.POST("/units/add", controllers.AddUnit)
	r.POST("/units/edit/:id", controllers.UpdateUnit)
	r.GET("/units/get/:id", controllers.GetUnit)
	r.GET("/units/list", controllers.GetUnitsList)

	//Raw materials
	r.GET("/raw-materials", controllers.ListRawMaterials)
	r.DELETE("/raw-materials/delete/:id", controllers.DeleteRawMaterial)
	r.POST("/raw-materials/add", controllers.AddRawMaterial)
	r.POST("/raw-materials/edit/:id", controllers.UpdateRawMaterial)
	r.GET("/raw-materials/get/:id", controllers.GetRawMaterial)
	r.GET("/raw-materials/list", controllers.GetRawMaterialsList)

	// Finished goods
	r.GET("/finished-goods", controllers.ListFinishedGoods)
	r.DELETE("/finished-goods/delete/:id", controllers.DeleteFinishedGood)
	r.POST("/finished-goods/add", controllers.AddFinishedGood)
	r.POST("/finished-goods/edit/:id", controllers.UpdateFinishedGood)
	r.GET("/finished-goods/get/:id", controllers.GetFinishedGood)
	r.GET("/finished-goods/list", controllers.GetFinishedGoodsList)

	// Ingredients
	r.GET("/ingredients", controllers.ListIngredients)
	r.DELETE("/ingredients/delete/:id", controllers.DeleteIngredient)
	r.POST("/ingredients/add", controllers.AddIngredient)
	r.POST("/ingredients/edit/:id", controllers.UpdateIngredient)
	r.GET("/ingredients/get/:id", controllers.GetIngredient)
	r.GET("/ingredients/:product_id", controllers.GetIngredientsByProduct)
	r.GET("/ingredients/used-raw-materials/:product_id", controllers.GetUsedRawMaterialsByProduct)
	r.GET("/ingredients/list", controllers.GetIngredientsList)

	// Purchases - table
	r.GET("/raw-material-purchases", controllers.ListRawMaterialPurchases)
	r.DELETE("/purchases/delete/:id", controllers.DeletePurchase)
	r.POST("/purchases/add", controllers.AddPurchase)

	// Employees
	r.GET("/employees", controllers.ListEmployees)
	r.GET("/employees/list", controllers.GetEmployeesList)
	r.GET("/positions/list", controllers.GetAllPositions)
	r.GET("/employees/get/:id", controllers.GetEmployee)
	r.POST("/employees/add", controllers.AddEmployee)
	r.POST("/employees/edit/:id", controllers.UpdateEmployee)
	r.DELETE("/employees/delete/:id", controllers.DeleteEmployee)

	// Positions
	r.POST("/positions/add", controllers.AddPosition)
	r.POST("/positions/edit/:id", controllers.EditPosition)
	r.DELETE("/positions/delete/:id", controllers.DeletePosition)
	r.GET("/positions/get/:id", controllers.GetPosition)
	r.GET("/positions", controllers.ListPositions)

	// Budget
	r.GET("/budget", controllers.BudgetList)
	r.GET("/budget/get-row/:id", controllers.GetBudgetRow)
	r.GET("/budget/get", controllers.GetBudget)
	r.GET("/markup/get", controllers.GetMarkup)
	r.PUT("/budget/update/:id", controllers.UpdateBudget)

	// Production
	r.GET("/production", controllers.ListProductProduction)
	r.POST("/production/produce/:product_id", controllers.ProduceProduct)

	// Sales
	r.GET("/sales", controllers.ListSales)
	r.POST("/sales/add", controllers.MakeSale)

	// Salaries
	r.GET("/salaries", controllers.ShowSalariesPage)
	r.GET("/salaries/:year/:month", controllers.GetSalaryByDate)
	r.POST("/salaries/calculate/:year/:month", controllers.CalculateSalary)
	r.PUT("/salaries/edit/:id", controllers.EditSalary)
	r.POST("/salaries/pay/:year/:month", controllers.PaySalaries)
	r.GET("/salaries/total-unpaid/:year/:month", controllers.GetUnpaidSalariesTotal)

	// WebSocket
	r.GET("/ws", gin.WrapF(ws.HandleWebSocket))
}
