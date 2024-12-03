package september

import (
	"context"
	"log/slog"
	"net/http"
	"truck-analytics-platform/internal/db"

	"github.com/gin-gonic/gin"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

func NineMonth2024Ldt(ctx *gin.Context) {
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

	// Ответ с данными по анализу грузовиков
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
    FROM ldt_3_5_12_truck_analytics_10_2024
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
    FROM ldt_3_5_12_truck_analytics_10_2024
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
		slog.Warn("Can't connect to database")
		ctx.JSON(http.StatusInternalServerError, TruckAnalyticsResponse{Error: "Can't connect to database"})
		return
	}

	rows, err := db.Query(context.Background(), query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, TruckAnalyticsResponse{Error: "Failed to execute query: " + err.Error()})
		return
	}
	defer rows.Close()

	// Ordered map to maintain custom order for districts
	dataByDistrict := orderedmap.New[string, []*TruckAnalytics]()

	// Define the custom order of districts
	customOrder := []string{
		"Central",
		"North West",
		"Volga",
		"South",
		"Ural",
		"Siberia",
		"Far East",
	}

	// Prepopulate the map with the desired order
	for _, district := range customOrder {
		dataByDistrict.Set(district, []*TruckAnalytics{})
	}

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
		ta.TotalMarket = nullToZero(&ta.DONGFENG) + nullToZero(&ta.FOTON) + nullToZero(&ta.GAZ) + nullToZero(&ta.JAC) + nullToZero(&ta.ISUZU) + nullToZero(&ta.JAC) + nullToZero(&ta.KAMAZ) + nullToZero(&ta.OTHER)

		// Add data to the appropriate district
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

func NineMonth2024LDTTotal(ctx *gin.Context) {

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

	type Summary struct {
		DONGFENG    int `json:"dongfeng"`
		FOTON       int `json:"foton"`
		GAZ         int `json:"gaz"`
		ISUZU       int `json:"isuzu"`
		JAC         int `json:"jac"`
		KAMAZ       int `json:"kamaz"`
		OTHER       int `json:"other"`
		TotalMarket int `json:"total"`
	}

	type TruckAnalyticsResponse struct {
		Data  map[string][]TruckAnalytics `json:"data"`
		Error string                      `json:"error,omitempty"`
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
		FROM ldt_3_5_12_truck_analytics_10_2024
		WHERE 
			"Month_of_registration" <= 9
		GROUP BY "Federal_district"
		ORDER BY "Federal_district"
	`

	// Подключение к базе данных
	db, err := db.Connect()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, TruckAnalyticsResponse{Error: "Can't connect to database"})
		return
	}
	defer db.Close(context.Background())

	// Запрос к базе данных
	rows, err := db.Query(context.Background(), query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, TruckAnalyticsResponse{Error: "Failed to execute query: " + err.Error()})
		return
	}
	defer rows.Close()

	var data []TruckAnalytics
	var summary Summary

	// Обработка данных из результата SQL запроса
	for rows.Next() {
		var ta TruckAnalytics
		var federalDistrict string

		err := rows.Scan(
			&federalDistrict,
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

		// Суммирование данных для итогов
		summary.DONGFENG += ta.DONGFENG
		summary.FOTON += ta.FOTON
		summary.GAZ += ta.GAZ
		summary.ISUZU += ta.ISUZU
		summary.JAC += ta.JAC
		summary.KAMAZ += ta.KAMAZ
		summary.OTHER += ta.OTHER
		summary.TotalMarket += ta.TotalMarket

		// Добавляем данные о регионе в список
		ta.RegionName = federalDistrict
		data = append(data, ta)
	}

	// Проверка на ошибки при итерации
	if err := rows.Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, TruckAnalyticsResponse{Error: "Error iterating over rows: " + err.Error()})
		return
	}

	// Создаем карту для ответа
	response := make(map[string][]TruckAnalytics)

	// Добавляем итоговые данные в начало
	response["Summary"] = []TruckAnalytics{
		{
			RegionName:  "Summary",
			DONGFENG:    summary.DONGFENG,
			FOTON:       summary.FOTON,
			GAZ:         summary.GAZ,
			ISUZU:       summary.ISUZU,
			JAC:         summary.JAC,
			KAMAZ:       summary.KAMAZ,
			OTHER:       summary.OTHER,
			TotalMarket: summary.TotalMarket,
		},
	}

	// Добавляем данные по округам
	for _, ta := range data {
		if _, exists := response[ta.RegionName]; !exists {
			response[ta.RegionName] = []TruckAnalytics{}
		}
		response[ta.RegionName] = append(response[ta.RegionName], ta)
	}

	// Отправка ответа
	ctx.JSON(http.StatusOK, TruckAnalyticsResponse{
		Data: response,
	})
}
