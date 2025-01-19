package routes 

import (
	"Diplomski/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func addMaterial(context *gin.Context) {
	var newMaterial models.Materijal
	err := context.ShouldBindJSON(&newMaterial)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not parse data."})
		return
	}

	err = newMaterial.InsertMaterial()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not add material."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "New material added!", "Material": newMaterial})
}

func getMaterials(context *gin.Context){
	materials, err := models.GetAllMaterials()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch materials."})
		return
	}
	context.JSON(http.StatusOK, materials)
}

func getMaterial(context *gin.Context) {
	materialId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not parse material id."})
		return
	}

	material, err := models.GetMaterialByID(materialId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch material."})
		return
	}

	context.JSON(http.StatusOK, material)

}

func updateMaterial(context *gin.Context) {
	materialId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not parse material id."})
		return
	}

	var updatedMaterial models.Materijal
	err = context.ShouldBindJSON(&updatedMaterial)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not parse requested data."})
		return
	}

	updatedMaterial.Materijal_ID = materialId
	err = updatedMaterial.Update()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not update material."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Material updated successfully!"})
}

func deleteMaterial(context *gin.Context) {
	materialId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not get material id."})
		return
	}

	material, err := models.GetMaterialByID(materialId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch material."})
		return
	}

	err = material.Delete()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not delete material."})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Delete completed successfully."})

}