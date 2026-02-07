package handlers

import (
	"net/http"

	"example.com/student-api/models"
	"example.com/student-api/services"
	"github.com/gin-gonic/gin"
)

type StudentHandler struct {
	Service *services.StudentService
}

func (h *StudentHandler) GetStudents(c *gin.Context) {
	students, err := h.Service.GetStudents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, students)
}

func (h *StudentHandler) GetStudentByID(c *gin.Context) {
	id := c.Param("id")
	student, err := h.Service.GetStudentByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}
	c.JSON(http.StatusOK, student)
}

func (h *StudentHandler) CreateStudent(c *gin.Context) {
	var student models.Student
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.Service.CreateStudent(student); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, student)
}

// PUT /students/:id
func (h *StudentHandler) UpdateStudent(c *gin.Context) {
	id := c.Param("id")
	var student models.Student
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}
	if err := h.Service.UpdateStudent(id, student); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found or invalid data"})
		return
	}
	student.Id = id
	c.JSON(http.StatusOK, student)
}

// DELETE /students/:id
func (h *StudentHandler) DeleteStudent(c *gin.Context) {
	id := c.Param("id")
	if err := h.Service.DeleteStudent(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}
	c.Status(http.StatusNoContent)
}
