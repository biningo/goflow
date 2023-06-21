package goflow

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (g *Goflow) addStaticRoutes() *Goflow {
	g.router.Static("/css", g.Options.AssetBasePath+"css")
	g.router.Static("/dist", g.Options.AssetBasePath+"dist")
	g.router.Static("/src", g.Options.AssetBasePath+"src")
	g.router.LoadHTMLGlob(g.Options.AssetBasePath + "html/*.html.tmpl")
	return g
}

func (g *Goflow) addStreamRoute() *Goflow {
	g.router.GET("/stream", g.getJobRuns(g.Options.StreamJobRuns))
	return g
}

func (g *Goflow) addAPIRoutes() *Goflow {
	api := g.router.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			var msg struct {
				Health string `json:"health"`
			}
			msg.Health = "OK"
			c.JSON(http.StatusOK, msg)
		})

		api.GET("/jobs", func(c *gin.Context) {
			jobNames := make([]string, 0)
			for _, job := range g.Jobs {
				jobNames = append(jobNames, job().Name)
			}
			var msg struct {
				Jobs []string `json:"jobs"`
			}
			msg.Jobs = jobNames
			c.JSON(http.StatusOK, msg)
		})

		api.GET("/jobs/:name", func(c *gin.Context) {
			name := c.Param("name")
			jobFn, ok := g.Jobs[name]

			if ok {
				tasks := jobFn().Tasks
				taskNames := make([]string, 0)
				for _, task := range tasks {
					taskNames = append(taskNames, task.Name)
				}

				var msg struct {
					JobName   string   `json:"job"`
					TaskNames []string `json:"tasks"`
					Schedule  string   `json:"schedule"`
				}
				msg.JobName = name
				msg.TaskNames = taskNames
				msg.Schedule = g.Jobs[name]().Schedule

				c.JSON(http.StatusOK, msg)
			} else {
				c.JSON(http.StatusNotFound, "Not found")
			}
		})

		api.GET("/jobs/:name/dag", func(c *gin.Context) {
			name := c.Param("name")
			jobFn, ok := g.Jobs[name]

			if ok {
				c.JSON(http.StatusOK, jobFn().Dag)
			} else {
				c.JSON(http.StatusNotFound, "Not found")
			}
		})

		api.GET("/jobs/:name/isActive", func(c *gin.Context) {
			name := c.Param("name")
			jobFn, ok := g.Jobs[name]

			if ok {
				c.JSON(http.StatusOK, jobFn().Active)
			} else {
				c.JSON(http.StatusNotFound, "Not found")
			}
		})

		api.POST("/jobs/:name/submit", func(c *gin.Context) {
			name := c.Param("name")
			_, ok := g.Jobs[name]

			var msg struct {
				Job     string `json:"job"`
				Success bool   `json:"success"`
				ID      string `json:"id"`
			}
			msg.Job = name

			if ok {
				jobRun := g.runJob(name)
				msg.Success = true
				msg.ID = jobRun.name()
				c.JSON(http.StatusOK, msg)
			} else {
				msg.Success = false
				c.JSON(http.StatusNotFound, msg)
			}
		})

		api.POST("/jobs/:name/toggleActive", func(c *gin.Context) {
			name := c.Param("name")
			_, ok := g.Jobs[name]

			var msg struct {
				Job     string `json:"job"`
				Success bool   `json:"success"`
				Active  bool   `json:"active"`
			}
			msg.Job = name

			if ok {
				isActive, _ := g.toggleActive(name)
				msg.Success = true
				msg.Active = isActive
				c.JSON(http.StatusOK, msg)
			} else {
				msg.Success = false
				c.JSON(http.StatusNotFound, msg)
			}
		})
	}

	return g
}

func (g *Goflow) addUIRoutes() *Goflow {
	ui := g.router.Group("/ui")
	{
		ui.GET("/", func(c *gin.Context) {
			jobs := make([]*Job, 0)
			for _, job := range g.Jobs {
				jobs = append(jobs, job())
			}
			c.HTML(http.StatusOK, "index.html.tmpl", gin.H{"jobs": jobs})
		})

		ui.GET("/jobs/:name", func(c *gin.Context) {
			name := c.Param("name")
			jobFn, ok := g.Jobs[name]

			if ok {
				tasks := jobFn().Tasks
				taskNames := make([]string, 0)
				for _, task := range tasks {
					taskNames = append(taskNames, task.Name)
				}

				c.HTML(http.StatusOK, "job.html.tmpl", gin.H{
					"jobName":   name,
					"taskNames": taskNames,
					"schedule":  g.Jobs[name]().Schedule,
				})
			} else {
				c.String(http.StatusNotFound, "Not found")
			}
		})
	}

	g.router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/ui/")
	})

	return g
}
