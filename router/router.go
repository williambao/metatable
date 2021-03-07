package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/williambao/metatable/router/middleware/header"
	"github.com/williambao/metatable/router/middleware/session"
	"github.com/williambao/metatable/server"
)

func Load(middleware ...gin.HandlerFunc) http.Handler {
	e := gin.New()
	e.Use(gin.Recovery())

	e.Use(header.NoCache)
	e.Use(header.Options)
	e.Use(middleware...)
	e.Use(session.SetUser())

	api := e.Group("/api")

	account := api.Group("/account")
	{
		account.POST("/login", server.Login)
		account.POST("/loginwx", server.LoginWX)
		account.POST("/register", server.Register)
		// account.GET("/request-sms/:phone", server.RequestSMS)

	}

	user := api.Group("/user")
	{
		user.Use(session.MustActiveUser())

		user.GET("", server.GetSelf)
		user.PUT("/:id", server.UpdateUser)
		user.GET("/:id/tables", server.GetUserTables)
		user.POST("/:id/update-account", server.UpdateAccount)
		user.POST("/:id/password", server.UpdatePassword)
	}

	users := api.Group("/users")
	{
		users.Use(session.MustAdmin())

	}

	tables := api.Group("/tables")
	{
		// tables.Use(session.MustActiveUser())

		tables.POST("", server.CreateTable)
		// tables.GET("/:tableId")

		table := tables.Group("/:tableId")
		{
			table.Use(session.SetTable())
			table.Use(session.SetTablePermission())

			table.GET("", session.MustTableViewDataPermission(), server.GetTable)
			table.PUT("", session.MustTableAdmin(), server.UpdateTable)
			table.DELETE("", session.MustTableAdmin(), server.DeleteTable)
			table.PUT("/column-orders", session.MustTableAdmin(), server.UpdateColumnOrders)

			subscribes := table.Group("/subscribes")
			{
				subscribes.GET("", server.GetTableSubscribes)
				subscribes.POST("", server.AddTableSubscribe)
				subscribes.DELETE("/:tableSubscribeId", server.DeleteTableSubscribe)
			}

			users := table.Group("/users")
			{
				users.GET("", session.MustTableViewDataPermission(), server.GetTableUsers)
				users.POST("", session.MustTableAdmin(), server.AddTableUser)
				users.PUT("/:tableUserId", session.MustTableAdmin(), server.UpdateTableUser)
				users.DELETE("/:tableUserId", session.MustTableAdmin(), server.DeleteTableUser)
			}

			// 编辑表格列, 只有表格管理员才可编辑列
			columns := table.Group("/columns")
			{
				columns.GET("", session.MustTableViewDataPermission(), server.GetTableColumns)
				columns.POST("", session.MustTableAdmin(), server.CreateTableColumn)
				columns.GET("/:columnId", session.MustTableAdmin(), server.GetColumn)
				columns.PUT("/:columnId", session.MustTableAdmin(), server.UpdateColumn)
				columns.DELETE("/:columnId", session.MustTableAdmin(), server.DeleteColumn)

				// columns.POST("/:columnId/options", session.MustTableAdmin(), server.CreateColumnOption)
				// columns.PUT("/:columnId/options/:optionId", session.MustTableAdmin(), server.UpdateColumnOption)

			}

			// 表格数据
			// 只有订阅了表格的人就能查看数据.
			// 新增,编辑, 删除则单独检查权限
			records := table.Group("/data")
			{
				records.GET("", session.MustTableViewDataPermission(), server.GetRecords)
				records.POST("", session.MustTableInsertDataPermission(), server.CreateRecord)
				records.DELETE("", session.MustTableDeleteDataPermission(), server.DeleteRecordByBatch)
				record := records.Group("/:recordId")
				{
					record.PUT("", session.MustTableEditDataPermission(), server.EditReocrdColumn)
					record.PATCH("", session.MustTableEditDataPermission(), server.PatchRecordColumn)
					record.DELETE("", session.MustTableDeleteDataPermission(), server.DeleteRecord)
				}
			}

			views := table.Group("/views")
			{
				views.POST("", session.MustTableEditDataPermission(), server.CreateTableView)
				view := views.Group("/:viewId")
				{
					view.PUT("", session.MustTableEditDataPermission(), server.UpdateTableView)
					view.GET("/data", session.MustTableViewDataPermission(), server.GetViewRecords)
				}

			}
		}

	}

	organizations := api.Group("/organizations")
	{
		organizations.POST("", session.MustActiveUser(), server.CreateOrganization)
		organizations.GET("", session.MustActiveUser(), server.GetUserOrganizations)
		organization := organizations.Group("/:organizationId")
		{
			organization.PUT("", session.MustActiveUser(), server.UpdateOrganization)
			organization.DELETE("", session.MustActiveUser(), server.DeleteOrganization)
			organization.GET("/members", session.MustActiveUser(), server.GetOrganizationUsers)
			organization.POST("/members", session.MustActiveUser(), server.AddOrganizationUser)
			organization.PUT("/members/:organizationUserId", session.MustActiveUser(), server.UpdateOrganizationUser)
			organization.DELETE("/members/:organizationUserId", session.MustActiveUser(), server.DeleteOrganizationUser)
		}
	}

	templates := api.Group("/templates")
	{
		templates.GET("", server.GetTemplates)
		templates.POST("", session.MustActiveUser(), server.CreateTemplate)
		template := templates.Group("/:templateId")
		{
			template.GET("", server.GetTemplate)
			template.PUT("", session.MustActiveUser(), server.UpdateTemplate)
			template.DELETE("", session.MustActiveUser(), server.DeleteTemplate)
			template.PUT("/column-orders", session.MustActiveUser(), server.UpdateTemplateColumnOrders)

			// 编辑表格列, 只有表格管理员才可编辑列
			templateColumns := template.Group("/columns")
			{
				templateColumns.GET("", server.GetTemplateColumns)
				templateColumns.POST("", session.MustActiveUser(), server.CreateTemplateColumn)
				// columns.GET("/:columnId", session.MustTableAdmin(), server.GetColumn)
				templateColumns.PUT("/:columnId", session.MustActiveUser(), server.UpdateTemplateColumn)
				templateColumns.DELETE("/:columnId", session.MustActiveUser(), server.DeleteTemplateColumn)

				// columns.POST("/:columnId/options", session.MustTableAdmin(), server.CreateColumnOption)
				// columns.PUT("/:columnId/options/:optionId", session.MustTableAdmin(), server.UpdateColumnOption)

			}
		}
	}

	files := api.Group("/files")
	{
		files.GET("/token", session.MustActiveUser(), server.GetFileUploadToken)
		files.DELETE("/:id", session.MustActiveUser(), server.DeleteFileById)
		files.POST("/callback", server.FileUploadCallback)
	}

	dicts := api.Group("/dict")
	{
		dicts.GET("/conditions", server.GetConditions)
	}
	return e
}
