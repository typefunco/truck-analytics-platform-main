package september

import (
	"context"
	"log/slog"
	"net/http"
	"truck-analytics-platform/internal/db"

	"github.com/gin-gonic/gin"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

func NineMonth2024Mdt(ctx *gin.Context) {
	type TruckAnalytics struct {
		RegionName  string `json:"region_name"`
		DONGFENG    int    `json:"dongfeng"`
		FOTON       int    `json:"foton"`
		HOWO        int    `json:"howo"`
		JAC         int    `json:"jac"`
		KAMAZ       int    `json:"kamaz"`
		URAL        int    `json:"ural"`
		DAEWOO      int    `json:"daewoo"`
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
        SUM(CAST("Quantity" AS INTEGER)) AS total_sales -- Cast Quantity to INTEGER
    FROM mdt_12_18_truck_analytics_10_2023
    WHERE "Brand" IN ('DONGFENG', 'FOTON', 'HOWO', 'JAC', 'KAMAZ', 'URAL', 'DAEWOO')
    AND "Month_of_registration" <= 9  -- Filter for months <= 9
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
        SUM(CAST("Quantity" AS INTEGER)) AS total_sales -- Cast Quantity to INTEGER
    FROM mdt_12_18_truck_analytics_10_2024
    WHERE "Brand" NOT IN ('DONGFENG', 'FOTON', 'HOWO', 'JAC', 'KAMAZ', 'URAL', 'DAEWOO')
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
        MAX(CASE WHEN "Brand" = 'HOWO' THEN total_sales END) AS HOWO,
        MAX(CASE WHEN "Brand" = 'JAC' THEN total_sales END) AS JAC,
        MAX(CASE WHEN "Brand" = 'KAMAZ' THEN total_sales END) AS KAMAZ,
        MAX(CASE WHEN "Brand" = 'URAL' THEN total_sales END) AS URAL,
        MAX(CASE WHEN "Brand" = 'DAEWOO' THEN total_sales END) AS DAEWOO,
        MAX(CASE WHEN "Brand" = 'OTHER' THEN total_sales END) AS OTHER,
        COALESCE(MAX(CASE WHEN "Brand" = 'DONGFENG' THEN total_sales END), 0) +
        COALESCE(MAX(CASE WHEN "Brand" = 'FOTON' THEN total_sales END), 0) +
        COALESCE(MAX(CASE WHEN "Brand" = 'HOWO' THEN total_sales END), 0) +
        COALESCE(MAX(CASE WHEN "Brand" = 'JAC' THEN total_sales END), 0) +
        COALESCE(MAX(CASE WHEN "Brand" = 'KAMAZ' THEN total_sales END), 0) +
        COALESCE(MAX(CASE WHEN "Brand" = 'URAL' THEN total_sales END), 0) +
        COALESCE(MAX(CASE WHEN "Brand" = 'DAEWOO' THEN total_sales END), 0) +
        COALESCE(MAX(CASE WHEN "Brand" = 'OTHER' THEN total_sales END), 0) AS TOTAL
    FROM combined_data
    GROUP BY "Region", "Federal_district"
),
federal_totals_by_region AS (
    SELECT 
        "Federal_district",
        'TOTAL' AS Oblast_name,  
        SUM(DONGFENG) AS DONGFENG,
        SUM(FOTON) AS FOTON,
        SUM(HOWO) AS HOWO,
        SUM(JAC) AS JAC,
        SUM(KAMAZ) AS KAMAZ,
        SUM(URAL) AS URAL,
        SUM(DAEWOO) AS DAEWOO,
        SUM(OTHER) AS OTHER,
        SUM(TOTAL) AS TOTAL
    FROM pivoted_data
    GROUP BY "Federal_district"
),
final_data AS (
    SELECT Oblast_name, "Federal_district", DONGFENG, FOTON, HOWO, JAC, KAMAZ, URAL, DAEWOO, OTHER, TOTAL
    FROM pivoted_data
    UNION ALL
    SELECT Oblast_name, "Federal_district", DONGFENG, FOTON, HOWO, JAC, KAMAZ, URAL, DAEWOO, OTHER, TOTAL
    FROM federal_totals_by_region
)
SELECT 
    "Federal_district",  -- Changed column order here
    CASE 
        WHEN Oblast_name = 'TOTAL' THEN "Federal_district"  
        ELSE Oblast_name 
    END AS Oblast_name,
    COALESCE(DONGFENG, 0) AS DONGFENG,
    COALESCE(FOTON, 0) AS FOTON,
    COALESCE(HOWO, 0) AS HOWO,
    COALESCE(JAC, 0) AS JAC,
    COALESCE(KAMAZ, 0) AS KAMAZ,
    COALESCE(URAL, 0) AS URAL,
    COALESCE(DAEWOO, 0) AS DAEWOO,
    COALESCE(OTHER, 0) AS OTHER,
    COALESCE(TOTAL, 0) AS TOTAL
FROM final_data
ORDER BY 
    "Federal_district",  
    CASE 
        WHEN Oblast_name = 'TOTAL' THEN 1  
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
		"North Caucasian",
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
			&ta.HOWO,
			&ta.JAC,
			&ta.KAMAZ,
			&ta.URAL,
			&ta.DAEWOO,
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
		ta.TotalMarket = nullToZero(&ta.DONGFENG) + nullToZero(&ta.FOTON) + nullToZero(&ta.HOWO) + nullToZero(&ta.JAC) + nullToZero(&ta.KAMAZ) + nullToZero(&ta.JAC) + nullToZero(&ta.URAL) + nullToZero(&ta.DAEWOO) + nullToZero(&ta.OTHER)

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

type DistrictDataMDT struct {
	RegionName  string `json:"region_name"`
	DONGFENG    *int   `json:"dongfeng"`
	FOTON       *int   `json:"foton"`
	HOWO        *int   `json:"howo"`
	JAC         *int   `json:"jac"`
	KAMAZ       *int   `json:"kamaz"`
	URAL        *int   `json:"ural"`
	DAEWOO      *int   `json:"daewoo"`
	OTHER       *int   `json:"other"`
	TotalMarket int    `json:"total"`
}

func NineMonth2024MDTTotal(ctx *gin.Context) {

	type Response struct {
		Data  *orderedmap.OrderedMap[string, []DistrictDataMDT] `json:"data"`
		Error string                                            `json:"error,omitempty"`
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
        COALESCE(SUM(CASE WHEN "Brand" = 'DONGFENG' THEN CAST("Quantity" AS INTEGER) END), 0) AS DONGFENG,
        COALESCE(SUM(CASE WHEN "Brand" = 'FOTON' THEN CAST("Quantity" AS INTEGER) END), 0) AS FOTON,
        COALESCE(SUM(CASE WHEN "Brand" = 'HOWO' THEN CAST("Quantity" AS INTEGER) END), 0) AS HOWO,
        COALESCE(SUM(CASE WHEN "Brand" = 'JAC' THEN CAST("Quantity" AS INTEGER) END), 0) AS JAC,
        COALESCE(SUM(CASE WHEN "Brand" = 'KAMAZ' THEN CAST("Quantity" AS INTEGER) END), 0) AS KAMAZ,
        COALESCE(SUM(CASE WHEN "Brand" = 'URAL' THEN CAST("Quantity" AS INTEGER) END), 0) AS URAL,
        COALESCE(SUM(CASE WHEN "Brand" = 'DAEWOO' THEN CAST("Quantity" AS INTEGER) END), 0) AS DAEWOO,
        COALESCE(SUM(CASE WHEN "Brand" NOT IN ('DONGFENG', 'FOTON', 'HOWO', 'JAC', 'KAMAZ', 'URAL', 'DAEWOO') THEN CAST("Quantity" AS INTEGER) END), 0) AS OTHER,
        COALESCE(SUM(CAST("Quantity" AS INTEGER)), 0) AS TOTAL
    FROM mdt_12_18_truck_analytics_10_2024
    WHERE 
        "Month_of_registration" <= 9
    GROUP BY "Federal_district"
    ORDER BY "Federal_district";
	`

	// Подключение к базе данных
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

	dataByDistrict := orderedmap.New[string, []DistrictDataMDT]()
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
		dataByDistrict.Set(district, []DistrictDataMDT{})
	}

	var summary DistrictDataMDT
	summary.RegionName = "Summary"
	initializeIntFieldsMDT(&summary)

	// Обработка данных из результата SQL запроса
	for rows.Next() {
		var item DistrictDataMDT

		err := rows.Scan(
			&item.RegionName,
			&item.DONGFENG,
			&item.FOTON,
			&item.HOWO,
			&item.JAC,
			&item.KAMAZ,
			&item.URAL,
			&item.DAEWOO,
			&item.OTHER,
			&item.TotalMarket,
		)
		if err != nil {
			slog.Info("Failed to scan row:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row"})
			return
		}

		addToSummaryMDT(&summary, &item)

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

func initializeIntFieldsMDT(data *DistrictDataMDT) {
	data.DONGFENG = new(int)
	data.FOTON = new(int)
	data.HOWO = new(int)
	data.JAC = new(int)
	data.KAMAZ = new(int)
	data.DAEWOO = new(int)
	data.URAL = new(int)
	data.OTHER = new(int)
}

func addToSummaryMDT(summary *DistrictDataMDT, item *DistrictDataMDT) {
	addIntFields(summary.DONGFENG, item.DONGFENG)
	addIntFields(summary.FOTON, item.FOTON)
	addIntFields(summary.URAL, item.URAL)
	addIntFields(summary.HOWO, item.HOWO)
	addIntFields(summary.DAEWOO, item.DAEWOO)
	addIntFields(summary.JAC, item.JAC)
	addIntFields(summary.KAMAZ, item.KAMAZ)
	addIntFields(summary.OTHER, item.OTHER)
	summary.TotalMarket += item.TotalMarket
}
