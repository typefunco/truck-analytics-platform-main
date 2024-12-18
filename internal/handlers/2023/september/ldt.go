package september

import (
	"context"
	"log/slog"
	"net/http"
	"truck-analytics-platform/internal/db"

	"github.com/gin-gonic/gin"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

func NineMonth2023Ldt(ctx *gin.Context) {
	type TruckAnalytics struct {
		RegionName  string `json:"region_name"`
		DONGFENG    int    `json:"dongfeng"`
		FOTON       int    `json:"foton"`
		GAZ         int    `json:"gaz"`
		ISUZU       int    `json:"isuzu"`
		JAC         int    `json:"jac"`
		KAMAZ       int    `json:"kamaz"`
		OTHER       int    `json:"other"`
		TotalMarket int    `json:"total"`
	}

	type TruckAnalyticsResponse struct {
		Data  *orderedmap.OrderedMap[string, []*TruckAnalytics] `json:"data"`
		Error string                                            `json:"error,omitempty"`
	}

	// SQL запрос для получения данных
	query := `
	WITH base_data AS (
		SELECT 
			"Federal_district",
			"Region",  
			"Brand",
			"City",
			SUM("Quantity") AS total_sales
		FROM ldt_3_5_12_truck_analytics_10_2023
		WHERE "Brand" IN ('DONGFENG', 'FOTON', 'GAZ', 'ISUZU', 'JAC', 'KAMAZ')
		AND "Month_of_registration" <= 9   -- Filter for months <= 9
		GROUP BY "Federal_district", "Region", "City", "Brand"
	),
	federal_totals AS (
		SELECT 
			"Federal_district",
			"Region",
			'TOTAL' AS "City",
			"Brand",
			SUM(total_sales) AS total_sales
		FROM base_data
		GROUP BY "Federal_district", "Region", "Brand"
	),
	other_brands AS (
		SELECT 
			"Federal_district",
			"Region",  
			'OTHER' AS "Brand",
			'OTHER' AS "City",
			SUM("Quantity") AS total_sales
		FROM ldt_3_5_12_truck_analytics_10_2023
		WHERE "Brand" NOT IN ('DONGFENG', 'FOTON', 'GAZ', 'ISUZU', 'JAC', 'KAMAZ')
		AND "Month_of_registration" <= 9  -- Filter for months <= 9
		GROUP BY "Federal_district", "Region"
	),
	combined_data AS (
		SELECT "Federal_district", "Region", "City", "Brand", total_sales
		FROM base_data
		UNION ALL
		SELECT "Federal_district", "Region", "City", "Brand", total_sales
		FROM federal_totals
		UNION ALL
		SELECT "Federal_district", "Region", "City", "Brand", total_sales
		FROM other_brands
	),
	pivoted_data AS (
		SELECT 
			"Region" AS Oblast_name,
			"Federal_district", 
			MAX(CASE WHEN "Brand" = 'DONGFENG' THEN total_sales END) AS DONGFENG,
			MAX(CASE WHEN "Brand" = 'FOTON' THEN total_sales END) AS FOTON,
			MAX(CASE WHEN "Brand" = 'GAZ' THEN total_sales END) AS GAZ,
			MAX(CASE WHEN "Brand" = 'ISUZU' THEN total_sales END) AS ISUZU,
			MAX(CASE WHEN "Brand" = 'JAC' THEN total_sales END) AS JAC,
			MAX(CASE WHEN "Brand" = 'KAMAZ' THEN total_sales END) AS KAMAZ,
			MAX(CASE WHEN "Brand" = 'OTHER' THEN total_sales END) AS OTHER,
			COALESCE(MAX(CASE WHEN "Brand" = 'DONGFENG' THEN total_sales END), 0) +
			COALESCE(MAX(CASE WHEN "Brand" = 'FOTON' THEN total_sales END), 0) +
			COALESCE(MAX(CASE WHEN "Brand" = 'GAZ' THEN total_sales END), 0) +
			COALESCE(MAX(CASE WHEN "Brand" = 'ISUZU' THEN total_sales END), 0) +
			COALESCE(MAX(CASE WHEN "Brand" = 'JAC' THEN total_sales END), 0) +
			COALESCE(MAX(CASE WHEN "Brand" = 'KAMAZ' THEN total_sales END), 0) +
			COALESCE(MAX(CASE WHEN "Brand" = 'OTHER' THEN total_sales END), 0) AS TOTAL
		FROM combined_data
		GROUP BY "Region", "Federal_district"
	),
	federal_totals_by_region AS (
		SELECT 
			"Federal_district",
			"Federal_district" AS Oblast_name,  -- Use the Federal_district name in the total row
			SUM(DONGFENG) AS DONGFENG,
			SUM(FOTON) AS FOTON,
			SUM(GAZ) AS GAZ,
			SUM(ISUZU) AS ISUZU,
			SUM(JAC) AS JAC,
			SUM(KAMAZ) AS KAMAZ,
			SUM(OTHER) AS OTHER,
			SUM(TOTAL) AS TOTAL
		FROM pivoted_data
		GROUP BY "Federal_district"
	),
	final_data AS (
		SELECT Oblast_name, "Federal_district", DONGFENG, FOTON, GAZ, ISUZU, JAC, KAMAZ, OTHER, TOTAL
		FROM pivoted_data
		UNION ALL
		SELECT Oblast_name, "Federal_district", DONGFENG, FOTON, GAZ, ISUZU, JAC, KAMAZ, OTHER, TOTAL
		FROM federal_totals_by_region
	)
	SELECT 
		"Federal_district",  -- Federal_district comes first
		Oblast_name,  -- Oblast_name comes second
		COALESCE(DONGFENG, 0) AS DONGFENG,
		COALESCE(FOTON, 0) AS FOTON,
		COALESCE(GAZ, 0) AS GAZ,
		COALESCE(ISUZU, 0) AS ISUZU,
		COALESCE(JAC, 0) AS JAC,
		COALESCE(KAMAZ, 0) AS KAMAZ,
		COALESCE(OTHER, 0) AS OTHER,
		COALESCE(TOTAL, 0) AS TOTAL
	FROM final_data
	ORDER BY 
		"Federal_district",  -- Sort by Federal_district
		CASE 
			WHEN Oblast_name = "Federal_district" THEN 1  
			ELSE 0
		END,
		Oblast_name;
    `

	// Соединение с базой данных
	db, err := db.Connect()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, TruckAnalyticsResponse{Error: "Can't connect to database"})
		return
	}

	rows, err := db.Query(context.Background(), query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, TruckAnalyticsResponse{Error: "Failed to execute query: " + err.Error()})
		return
	}
	defer rows.Close()

	// Мапа для группировки данных по федеральным округам
	dataByDistrict := orderedmap.New[string, []*TruckAnalytics]()

	// Define the custom order of districts
	customOrder := []string{
		"Central",
		"North West",
		"Volga",
		"North Caucasian",
		"South",
		"Ural",
		"Siberia",
		"Far East",
	}

	// Prepopulate the map with the desired order
	for _, district := range customOrder {
		dataByDistrict.Set(district, []*TruckAnalytics{})
	}

	// Обработка данных из результата SQL запроса
	for rows.Next() {
		var ta TruckAnalytics
		var federalDistrict string

		err := rows.Scan(
			&federalDistrict,
			&ta.RegionName,
			&ta.DONGFENG,
			&ta.FOTON,
			&ta.GAZ,
			&ta.ISUZU,
			&ta.JAC,
			&ta.KAMAZ,
			&ta.OTHER,
			&ta.TotalMarket,
		)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, TruckAnalyticsResponse{Error: "Failed to scan row: " + err.Error()})
			return
		}

		// Перевод федерального округа, если есть в маппинге
		if translatedDistrict, ok := districtTranslations[federalDistrict]; ok {
			federalDistrict = translatedDistrict
		}

		// Перевод региона, если есть в маппинге
		if translatedRegion, ok := regionTranslations[ta.RegionName]; ok {
			ta.RegionName = translatedRegion
		}

		// Рассчитываем общий рынок
		ta.TotalMarket = null2Zero(&ta.DONGFENG) + null2Zero(&ta.FOTON) + null2Zero(&ta.GAZ) + null2Zero(&ta.JAC) + null2Zero(&ta.ISUZU) + null2Zero(&ta.KAMAZ) + null2Zero(&ta.OTHER)

		// Добавляем данные о регионе в соответствующий федеральный округ
		if existing, exists := dataByDistrict.Get(federalDistrict); exists {
			dataByDistrict.Set(federalDistrict, append(existing, &ta))
		}
	}

	// Проверка на ошибки при итерации
	if err := rows.Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, TruckAnalyticsResponse{Error: "Error iterating over rows: " + err.Error()})
		return
	}

	// Отправка ответа
	ctx.JSON(http.StatusOK, TruckAnalyticsResponse{
		Data: dataByDistrict,
	})
}

// Помощь для обработки нулевых значений
func null2Zero(val *int) int {
	if val == nil {
		return 0
	}
	return *val
}

type DistrictData struct {
	RegionName  string `json:"region_name"`
	DONGFENG    *int   `json:"dongfeng"`
	FOTON       *int   `json:"foton"`
	GAZ         *int   `json:"gaz"`
	ISUZU       *int   `json:"isuzu"`
	JAC         *int   `json:"jac"`
	KAMAZ       *int   `json:"kamaz"`
	OTHER       *int   `json:"other"`
	TotalMarket int    `json:"total"`
}

func NineMonth2023LDTTotal(ctx *gin.Context) {
	// Структура для ответа
	type Response struct {
		Data  *orderedmap.OrderedMap[string, []DistrictData] `json:"data"`
		Error string                                         `json:"error,omitempty"`
	}

	// Мапа для перевода федеральных округов и регионов
	var districtTranslations = map[string]string{
		"Дальневосточный Федеральный Округ":   "Far East",
		"Приволжский Федеральный Округ":       "Volga",
		"Северо-Западный Федеральный Округ":   "North West",
		"Северо-Кавказский Федеральный Округ": "North Caucasian",
		"Сибирский Федеральный Округ":         "Siberia",
		"Уральский Федеральный Округ":         "Ural",
		"Центральный Федеральный Округ":       "Central",
		"Южный Федеральный Округ":             "South",
	}

	// SQL-запрос для получения данных по округам
	query := `
		SELECT 
			"Federal_district",
			COALESCE(SUM(CASE WHEN "Brand" = 'DONGFENG' THEN "Quantity" END), 0) AS DONGFENG,
			COALESCE(SUM(CASE WHEN "Brand" = 'FOTON' THEN "Quantity" END), 0) AS FOTON,
			COALESCE(SUM(CASE WHEN "Brand" = 'GAZ' THEN "Quantity" END), 0) AS GAZ,
			COALESCE(SUM(CASE WHEN "Brand" = 'ISUZU' THEN "Quantity" END), 0) AS ISUZU,
			COALESCE(SUM(CASE WHEN "Brand" = 'JAC' THEN "Quantity" END), 0) AS JAC,
			COALESCE(SUM(CASE WHEN "Brand" = 'KAMAZ' THEN "Quantity" END), 0) AS KAMAZ,
			COALESCE(SUM(CASE WHEN "Brand" NOT IN ('DONGFENG', 'FOTON', 'GAZ', 'ISUZU', 'JAC', 'KAMAZ') THEN "Quantity" END), 0) AS OTHER,
			COALESCE(SUM("Quantity"), 0) AS TOTAL
		FROM ldt_3_5_12_truck_analytics_10_2023
		WHERE 
			"Month_of_registration" <= 9
		GROUP BY "Federal_district"
		ORDER BY "Federal_district"
	`

	db, err := db.Connect()
	if err != nil {
		slog.Info("Can't connect to database:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Can't connect to database"})
		return
	}
	defer db.Close(context.Background())

	rows, err := db.Query(context.Background(), query)
	if err != nil {
		slog.Info("Failed to execute query:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute query"})
		return
	}
	defer rows.Close()

	dataByDistrict := orderedmap.New[string, []DistrictData]()
	customOrder := []string{
		"Summary",
		"Central",
		"North West",
		"Volga",
		"South",
		"North Caucasian",
		"Ural",
		"Siberia",
		"Far East",
	}

	for _, district := range customOrder {
		dataByDistrict.Set(district, []DistrictData{})
	}

	var summary DistrictData
	summary.RegionName = "Summary"
	initializeIntFields(&summary)

	for rows.Next() {
		var item DistrictData
		err := rows.Scan(
			&item.RegionName,
			&item.DONGFENG,
			&item.FOTON,
			&item.GAZ,
			&item.ISUZU,
			&item.JAC,
			&item.KAMAZ,
			&item.OTHER,
			&item.TotalMarket,
		)
		if err != nil {
			slog.Info("Failed to scan row:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row"})
			return
		}

		addToSummary(&summary, &item)

		translatedRegionName, exists := districtTranslations[item.RegionName]
		if !exists {
			translatedRegionName = item.RegionName
		}

		if existing, exists := dataByDistrict.Get(translatedRegionName); exists {
			item.RegionName = translatedRegionName
			dataByDistrict.Set(translatedRegionName, append(existing, item))
		}
	}

	if rows.Err() != nil {
		slog.Info("Failed to iterate over rows:", rows.Err())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to iterate over rows"})
		return
	}

	if summaryData, exists := dataByDistrict.Get("Summary"); exists {
		dataByDistrict.Set("Summary", append(summaryData, summary))
	}

	ctx.JSON(http.StatusOK, Response{
		Data: dataByDistrict,
	})
}

func initializeIntFields(data *DistrictData) {
	data.DONGFENG = new(int)
	data.FOTON = new(int)
	data.GAZ = new(int)
	data.ISUZU = new(int)
	data.JAC = new(int)
	data.KAMAZ = new(int)
	data.OTHER = new(int)
}

func addToSummary(summary *DistrictData, item *DistrictData) {
	addIntFields(summary.DONGFENG, item.DONGFENG)
	addIntFields(summary.FOTON, item.FOTON)
	addIntFields(summary.GAZ, item.GAZ)
	addIntFields(summary.ISUZU, item.ISUZU)
	addIntFields(summary.JAC, item.JAC)
	addIntFields(summary.KAMAZ, item.KAMAZ)
	addIntFields(summary.OTHER, item.OTHER)
	summary.TotalMarket += item.TotalMarket
}

func addIntFields(a, b *int) {
	if a != nil && b != nil {
		*a += *b
	}
}
