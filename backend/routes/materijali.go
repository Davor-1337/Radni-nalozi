package routes 

import (
	"Diplomski/models"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"Diplomski/database"
	"database/sql"
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

func getMaterialCategoryStats(context *gin.Context) {

	materials, err := models.GetMaterialUsedByCategory()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch stats for materials."})
		return
	}
	context.JSON(http.StatusOK, materials)
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

func filterMaterials(c *gin.Context) {
    var payload struct {
        Search string `json:"search"`
    }
    if err := c.ShouldBindJSON(&payload); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message":"Invalid input","error":err.Error()})
        return
    }

    type Materijal struct {
        Materijal_ID       int64   `json:"Materijal_ID"`
        NazivMaterijala    string  `json:"NazivMaterijala"`
        Kategorija         string  `json:"Kategorija"`
        Cijena             float64 `json:"Cijena"`
        KolicinaUSkladistu int64   `json:"KolicinaUSkladistu"`
    }

    query := `
      SELECT Materijal_ID, NazivMaterijala, Kategorija, Cijena, KolicinaUSkladistu
        FROM Materijal
       WHERE 1=1`
    args := []interface{}{}

    if payload.Search != "" {
        query += `
          AND (
            LOWER(NazivMaterijala) LIKE LOWER(@Search) OR
            LOWER(Kategorija)      LIKE LOWER(@Search)
          )`
        args = append(args, sql.Named("Search", "%"+payload.Search+"%"))
    }

    rows, err := database.DB.Query(query, args...)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message":"Query error","error":err.Error()})
        return
    }
    defer rows.Close()

    var list []Materijal
    for rows.Next() {
        var m Materijal
        if err := rows.Scan(&m.Materijal_ID, &m.NazivMaterijala, &m.Kategorija, &m.Cijena, &m.KolicinaUSkladistu); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"message":"Scan error","error":err.Error()})
            return
        }
        list = append(list, m)
    }
    c.JSON(http.StatusOK, gin.H{"materials": list})
}